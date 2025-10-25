package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	scylladb "github.com/golang-migrate/migrate/v4/database/cassandra"
	pgsql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/argon-chat/KineticaFS/pkg/repositories/postgres"
	"github.com/argon-chat/KineticaFS/pkg/repositories/scylla"
	"github.com/spf13/viper"
)

var migrationTypes []models.ApplicationRecord

type migration struct {
	TableName string
	Fields    map[string]string
}
type IDatabase interface {
	Connect() error
	Close() error
}

type ApplicationRepository struct {
	db               IDatabase
	dbType           string
	connectionString string

	ServiceTokens IServiceTokenRepository
	Buckets       IBucketRepository
	Files         IFileRepository
	FileBlobs     IFileBlobRepository
}

func (a *ApplicationRepository) Close() error {
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("close %s: %w", a.dbType, err)
	}
	return nil
}

func NewApplicationRepository() (*ApplicationRepository, error) {
	migrationTypes = []models.ApplicationRecord{
		models.ServiceToken{},
		models.Bucket{},
		models.File{},
		models.FileBlob{},
	}
	dbType := viper.GetString("database")
	if dbType == "" {
		return nil, fmt.Errorf("database type is not set")
	}
	connectionString := viper.GetString(dbType)
	if connectionString == "" {
		return nil, fmt.Errorf("connection string for %s is not set", dbType)
	}
	switch dbType {
	case "scylla":
		return newScyllaRepository(connectionString)
	case "postgres":
		return newPostgresRepository(connectionString)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

type migrator struct {
	repo *ApplicationRepository
}

func (m *migrator) Run(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := m.repo.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}
	return nil
}

func NewMigrator(repo *ApplicationRepository) *migrator {
	return &migrator{
		repo: repo,
	}
}

func (ar *ApplicationRepository) Migrate(ctx context.Context) error {
	switch ar.dbType {
	case "scylla":
		return ar.migrateDatabase(ctx)
	case "postgres":
		return ar.migrateDatabase(ctx)
	default:
		return fmt.Errorf("unsupported database type: %s", ar.dbType)
	}
}

func (ar *ApplicationRepository) migrateDatabase(ctx context.Context) error {
	path := viper.GetString("migration_path")
	if path == "" {
		path = "./migrations"
	}
	migrationPath := fmt.Sprintf("file://%s/%s", path, ar.dbType)
	log.Println(migrationPath)
	var scyllaDriver *gocql.Session
	var postgresDriver *sql.DB
	var m *migrate.Migrate
	var driver database.Driver
	var err error
	switch ar.dbType {
	case "scylla":
		scyllaDriver = ar.db.(*scylla.ScyllaConnection).Session
		driver, err = scylladb.WithInstance(scyllaDriver, &scylladb.Config{
			KeyspaceName: ar.db.(*scylla.ScyllaConnection).Keyspace,
		})
	case "postgres":
		postgresDriver = ar.db.(*postgres.PostgresConnection).DB
		driver, err = pgsql.WithInstance(postgresDriver, &pgsql.Config{
			DatabaseName: ar.db.(*postgres.PostgresConnection).DatabaseName,
		})
	}

	if err != nil {
		return fmt.Errorf("create %s migration driver error: %w", ar.dbType, err)
	}
	m, err = migrate.NewWithDatabaseInstance(
		migrationPath,
		ar.dbType, driver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func newPostgresRepository(connectionString string) (*ApplicationRepository, error) {
	repository, err := postgres.NewPostgresConnection(connectionString)
	if err != nil {
		return nil, err
	}

	ar := &ApplicationRepository{
		db:               repository,
		dbType:           "postgres",
		connectionString: connectionString,

		ServiceTokens: postgres.NewPostgresServiceTokenRepository(repository.DB),
		Buckets:       postgres.NewPostgresBucketRepository(repository.DB),
		Files:         postgres.NewPostgresFileRepository(repository.DB),
		FileBlobs:     postgres.NewPostgresFileBlobRepository(repository.DB),
	}
	log.Printf("Postgres repository created: %+v", ar)
	return ar, nil
}

func newScyllaRepository(connectionString string) (*ApplicationRepository, error) {
	repository, err := scylla.NewScyllaConnection(connectionString)
	if err != nil {
		return nil, err
	}

	ar := &ApplicationRepository{
		db:               repository,
		dbType:           "scylla",
		connectionString: connectionString,

		ServiceTokens: scylla.NewScyllaServiceTokenRepository(repository.Session),
		Buckets:       scylla.NewScyllaBucketRepository(repository.Session),
		Files:         scylla.NewScyllaFileRepository(repository.Session),
		FileBlobs:     scylla.NewScyllaFileBlobRepository(repository.Session),
	}
	log.Printf("Scylla repository created: %+v", ar)
	return ar, nil
}

func (r *ApplicationRepository) InitializeRepo(ctx context.Context, repo *ApplicationRepository) {
	r.ServiceTokens.CreateIndices(ctx)
	r.Buckets.CreateIndices(ctx)
	r.Files.CreateIndices(ctx)
	r.FileBlobs.CreateIndices(ctx)
}
