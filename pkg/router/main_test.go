package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/argon-chat/KineticaFS/pkg/repositories/postgres"
	"github.com/argon-chat/KineticaFS/pkg/repositories/scylla"

	// "github.com/argon-chat/KineticaFS/pkg/router"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	minioCtr "github.com/testcontainers/testcontainers-go/modules/minio"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	scylladbCtr "github.com/testcontainers/testcontainers-go/modules/scylladb"
)

var (
	sharedAppPostgres *testApp
	sharedAppScylla   *testApp
	minioEndpoint     string
	minioAccessKey    string
	minioSecretKey    string
	minioBucketName   string
	minioContainer    *minioCtr.MinioContainer
	minioClient       *minio.Client
	regionsConfigPath string
)

type testApp struct {
	repo   *repositories.ApplicationRepository
	server *httptest.Server
	dbType string
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Channel to collect errors from parallel initialization
	type initResult struct {
		name string
		err  error
	}
	results := make(chan initResult, 3)

	// Initialize PostgreSQL testcontainer
	var postgresContainer *pg.PostgresContainer
	go func() {
		var pgURL string
		if os.Getenv("TEST_DATABASE_URL") != "" {
			pgURL = os.Getenv("TEST_DATABASE_URL")
		} else {
			var err error
			postgresContainer, err = pg.Run(ctx,
				"postgres:18-alpine",
				pg.WithDatabase("kineticafs_test"),
				pg.WithUsername("postgres"),
				pg.WithPassword("postgres"),
				pg.BasicWaitStrategies(),
			)
			if err != nil {
				results <- initResult{"postgres", fmt.Errorf("failed to start postgres container: %w", err)}
				return
			}

			connString, err := postgresContainer.ConnectionString(ctx)
			if err != nil {
				results <- initResult{"postgres", fmt.Errorf("failed to get postgres connection string: %w", err)}
				return
			}
			pgURL = connString
		}

		app, err := setupTestAppWithDB(pgURL, "postgres")
		if err != nil {
			results <- initResult{"postgres", err}
			return
		}
		sharedAppPostgres = app
		results <- initResult{"postgres", nil}
	}()

	// Initialize ScyllaDB testcontainer
	var scyllaContainer *scylladbCtr.Container
	go func() {
		var scyllaURL string
		if os.Getenv("TEST_SCYLLA_URL") != "" {
			scyllaURL = os.Getenv("TEST_SCYLLA_URL")
		} else {
			var err error
			scyllaContainer, err = scylladbCtr.Run(ctx,
				"scylladb/scylla:latest",
			)
			if err != nil {
				results <- initResult{"scylla", fmt.Errorf("failed to start scylla container: %w", err)}
				return
			}

			connString, err := scyllaContainer.NonShardAwareConnectionHost(ctx)
			if err != nil {
				results <- initResult{"scylla", fmt.Errorf("failed to get scylla connection string: %w", err)}
				return
			}
			scyllaURL = connString
		}

		// Set keyspace for ScyllaDB
		viper.Set("scylla_keyspace", "kineticafs_test")

		app, err := setupTestAppWithDB(scyllaURL, "scylla")
		if err != nil {
			results <- initResult{"scylla", err}
			return
		}
		sharedAppScylla = app
		results <- initResult{"scylla", nil}
	}()

	// Initialize MinIO testcontainer
	go func() {
		var err error
		if os.Getenv("TEST_MINIO_URL") != "" {
			minioEndpoint = os.Getenv("TEST_MINIO_URL")
			minioAccessKey = os.Getenv("TEST_MINIO_ACCESS_KEY")
			minioSecretKey = os.Getenv("TEST_MINIO_SECRET_KEY")
		} else {
			minioContainer, err = minioCtr.Run(ctx,
				"minio/minio:latest",
				minioCtr.WithUsername("minioadmin"),
				minioCtr.WithPassword("minioadmin"),
			)
			if err != nil {
				results <- initResult{"minio", fmt.Errorf("failed to start minio container: %w", err)}
				return
			}

			endpoint, err := minioContainer.ConnectionString(ctx)
			if err != nil {
				results <- initResult{"minio", fmt.Errorf("failed to get minio connection string: %w", err)}
				return
			}
			minioEndpoint = endpoint
			minioAccessKey = "minioadmin"
			minioSecretKey = "minioadmin"
		}

		// Generate unique bucket name for this test run
		minioBucketName = fmt.Sprintf("kineticafs-test-%d", time.Now().Unix())

		// Create MinIO client
		minioClient, err = minio.New(minioEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
			Secure: false,
		})
		if err != nil {
			results <- initResult{"minio", fmt.Errorf("failed to create minio client: %w", err)}
			return
		}

		// Create MinIO bucket
		err = minioClient.MakeBucket(ctx, minioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			results <- initResult{"minio", fmt.Errorf("failed to create minio bucket: %w", err)}
			return
		}
		log.Printf("Created MinIO bucket: %s", minioBucketName)

		results <- initResult{"minio", nil}
	}()

	// Collect results from parallel initialization
	initErrors := make(map[string]error)
	for i := 0; i < 3; i++ {
		result := <-results
		if result.err != nil {
			initErrors[result.name] = result.err
			log.Printf("Failed to initialize %s: %v", result.name, result.err)
		} else {
			log.Printf("Successfully initialized %s", result.name)
		}
	}

	// Check if we have at least PostgreSQL
	if initErrors["postgres"] != nil {
		log.Printf("FATAL: PostgreSQL initialization failed: %v", initErrors["postgres"])
		os.Exit(1)
	}

	// ScyllaDB is optional for now
	if initErrors["scylla"] != nil {
		log.Printf("WARNING: ScyllaDB tests will be skipped: %v", initErrors["scylla"])
	}

	if initErrors["minio"] != nil {
		log.Printf("FATAL: MinIO initialization failed: %v", initErrors["minio"])
		os.Exit(1)
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup
	if sharedAppPostgres != nil {
		sharedAppPostgres.server.Close()
		sharedAppPostgres.repo.Close()
	}
	if sharedAppScylla != nil {
		sharedAppScylla.server.Close()
		sharedAppScylla.repo.Close()
	}
	if minioClient != nil && minioBucketName != "" {
		log.Printf("Cleaning up MinIO bucket: %s", minioBucketName)
		cleanTestS3Bucket(context.Background(), minioClient, minioBucketName)
		minioClient.RemoveBucket(context.Background(), minioBucketName)
	}
	if minioContainer != nil {
		log.Printf("Terminating MinIO container")
		minioContainer.Terminate(ctx)
	}
	if postgresContainer != nil {
		postgresContainer.Terminate(ctx)
	}
	if scyllaContainer != nil {
		log.Printf("Terminating ScyllaDB container")
		scyllaContainer.Terminate(ctx)
	}
	if regionsConfigPath != "" {
		os.Remove(regionsConfigPath)
	}

	os.Exit(exitCode)
}

func setupTestAppWithDB(connectionString, dbType string) (*testApp, error) {
	viper.Set("database", dbType)
	viper.Set(dbType, connectionString)
	viper.Set("migrate", true)
	viper.Set("migration_path", "../../migrations")
	viper.Set("cors-allowed-origins", "*")
	viper.Set("cors-allowed-headers", "Origin,Content-Type,Accept,Authorization,X-API-Token")

	repo, err := repositories.NewApplicationRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// Run migrations
	ctx := context.Background()
	if err := repo.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize repositories
	repo.InitializeRepo(ctx, repo)

	// Create test server
	r := NewRouter(repo, 0) // Port 0 for random port
	server := httptest.NewServer(r.Handler())

	return &testApp{
		repo:   repo,
		server: server,
		dbType: dbType,
	}, nil
}

func cleanTestDB(t *testing.T, app *testApp) {
	ctx := context.Background()

	switch app.dbType {
	case "postgres":
		// Get the database connection
		db := app.repo.GetDB().(*postgres.PostgresConnection).DB

		// Truncate all tables in reverse dependency order
		tables := []string{"file_blob", "file", "bucket", "service_token"}
		for _, table := range tables {
			query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
			if _, err := db.ExecContext(ctx, query); err != nil {
				if !strings.Contains(err.Error(), "does not exist") {
					t.Fatalf("Failed to truncate %s: %v", table, err)
				}
			}
		}
	case "scylla":
		// Get the ScyllaDB session
		session := app.repo.GetDB().(*scylla.ScyllaConnection).Session

		// Truncate all tables in reverse dependency order
		tables := []string{"FileBlob", "File", "Bucket", "ServiceToken"}
		for _, table := range tables {
			query := fmt.Sprintf("TRUNCATE %s", table)
			if err := session.Query(query).Exec(); err != nil {
				if !strings.Contains(err.Error(), "doesn't exist") {
					t.Fatalf("Failed to truncate %s: %v", table, err)
				}
			}
		}
	}

	// Clean S3 bucket
	if minioClient != nil && minioBucketName != "" {
		cleanTestS3Bucket(ctx, minioClient, minioBucketName)
	}
}

func setupTestApp(t *testing.T, backend string) *testApp {
	var app *testApp

	switch backend {
	case "postgres":
		if sharedAppPostgres == nil {
			t.Fatal("Shared PostgreSQL app not initialized")
		}
		app = sharedAppPostgres
	case "scylla":
		if sharedAppScylla == nil {
			t.Skip("ScyllaDB not available")
		}
		app = sharedAppScylla
	default:
		t.Fatalf("Unknown backend: %s", backend)
	}

	cleanTestDB(t, app)
	return app
}

func getOrCreateAdminToken(t *testing.T, app *testApp) string {
	// Check if this is first run
	resp := makeRequest(t, app, "GET", "/api/v1/st/first-run", "", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var firstRunResp struct {
		FirstRun bool `json:"first_run"`
	}
	err := json.NewDecoder(resp.Body).Decode(&firstRunResp)
	require.NoError(t, err)
	resp.Body.Close()

	if firstRunResp.FirstRun {
		// Create admin token via bootstrap
		resp = makeRequest(t, app, "POST", "/api/v1/st/bootstrap", "", nil)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var token models.ServiceToken
		err = json.NewDecoder(resp.Body).Decode(&token)
		require.NoError(t, err)
		resp.Body.Close()
		require.NotEmpty(t, token.AccessKey)

		return token.AccessKey
	}

	// Admin token already exists - this shouldn't happen with clean DB
	// but if it does, we can't retrieve the original access key
	t.Fatal("Database should be clean before getting admin token")
	return ""
}

func createTestServiceToken(t *testing.T, app *testApp, adminToken, name string) models.ServiceToken {
	body := map[string]interface{}{
		"name": name,
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	resp := makeAuthRequest(t, app, "POST", "/api/v1/st", adminToken, bytes.NewReader(bodyBytes))

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Create token failed: %s", string(bodyBytes))
	}
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var token models.ServiceToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	require.NoError(t, err)
	resp.Body.Close()

	return token
}

func createTestBucket(t *testing.T, app *testApp, adminToken, name, region string) models.Bucket {
	body := map[string]interface{}{
		"name":         name,
		"region":       region,
		"endpoint":     minioEndpoint,
		"access_key":   minioAccessKey,
		"secret_key":   minioSecretKey,
		"use_ssl":      false,
		"s3_provider":  "minio",
		"storage_type": 0,
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	resp := makeAuthRequest(t, app, "POST", "/api/v1/bucket", adminToken, bytes.NewReader(bodyBytes))

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Create bucket failed: %s", string(bodyBytes))
	}
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var bucket models.Bucket
	err = json.NewDecoder(resp.Body).Decode(&bucket)
	require.NoError(t, err)
	resp.Body.Close()

	return bucket
}

func makeRequest(t *testing.T, app *testApp, method, path, token string, body io.Reader) *http.Response {
	url := app.server.URL + path
	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	if token != "" {
		req.Header.Set("x-api-token", token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func makeAuthRequest(t *testing.T, app *testApp, method, path, token string, body io.Reader) *http.Response {
	resp := makeRequest(t, app, method, path, token, body)

	// Log response body for non-2xx responses
	if resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			t.Logf("Request %s %s failed with status %d: %s", method, path, resp.StatusCode, string(bodyBytes))
			// Create a new reader so the caller can still read the body
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}

	return resp
}

// cleanTestS3Bucket removes all objects from the test S3 bucket
func cleanTestS3Bucket(ctx context.Context, client *minio.Client, bucketName string) {
	objectsCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			log.Printf("Error listing objects in bucket %s: %v", bucketName, object.Err)
			continue
		}
		if err := client.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{}); err != nil {
			log.Printf("Error removing object %s from bucket %s: %v", object.Key, bucketName, err)
		}
	}
}

// listS3Objects returns a list of all object keys in the specified bucket
func listS3Objects(ctx context.Context, client *minio.Client, bucketName string) []string {
	var keys []string
	objectsCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			log.Printf("Error listing objects in bucket %s: %v", bucketName, object.Err)
			continue
		}
		keys = append(keys, object.Key)
	}
	return keys
}

// setupRegionsConfig creates a temporary regions.json file with test buckets
func setupRegionsConfig(t *testing.T, app *testApp, adminToken string) (string, *models.Bucket, *models.Bucket, *models.Bucket, func()) {
	ctx := context.Background()

	// Create test buckets in database
	bucket1 := createTestBucket(t, app, adminToken, "test-bucket-ru-1", "ru-1")
	bucket2 := createTestBucket(t, app, adminToken, "test-bucket-ru-2", "ru-1")
	bucket3 := createTestBucket(t, app, adminToken, "test-bucket-us-1", "us-east-1")

	// Create regions configuration
	regionsConfig := map[string]interface{}{
		"ru-1": map[string]interface{}{
			"id": 1,
			"buckets": []map[string]interface{}{
				{
					"id":       1,
					"bucketId": bucket1.ID,
				},
				{
					"id":       2,
					"bucketId": bucket2.ID,
				},
			},
		},
		"us-east-1": map[string]interface{}{
			"id": 2,
			"buckets": []map[string]interface{}{
				{
					"id":       1,
					"bucketId": bucket3.ID,
				},
			},
		},
	}

	// Write to temporary file
	tmpFile, err := os.CreateTemp("", "regions-*.json")
	require.NoError(t, err)
	defer tmpFile.Close()

	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(regionsConfig)
	require.NoError(t, err)

	regionsPath := tmpFile.Name()
	viper.Set("region", regionsPath)

	// Return cleanup function
	cleanup := func() {
		os.Remove(regionsPath)
		// Delete buckets from database
		app.repo.Buckets.DeleteBucket(ctx, bucket1.ID)
		app.repo.Buckets.DeleteBucket(ctx, bucket2.ID)
		app.repo.Buckets.DeleteBucket(ctx, bucket3.ID)
	}

	return regionsPath, &bucket1, &bucket2, &bucket3, cleanup
}
