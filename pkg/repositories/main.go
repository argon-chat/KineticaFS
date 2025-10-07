package repositories

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

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
	db     IDatabase
	dbType string

	ServiceTokens IServiceTokenRepository
	Buckets       IBucketRepository
	Files         IFileRepository
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
		return ar.migrateDatabase(ctx, migrateScyllaModel, ar.executeScyllaQuery)
	case "postgres":
		return ar.migrateDatabase(ctx, migratePostgresModel, ar.executePostgresQuery)
	default:
		return fmt.Errorf("unsupported database type: %s", ar.dbType)
	}
}

func (ar *ApplicationRepository) migrateDatabase(
	ctx context.Context,
	modelMigrator func(*migration) string,
	queryExecutor func(context.Context, string) error,
) error {
	for _, e := range migrationTypes {
		tableName := fmt.Sprintf("%T", e)
		model := migration{
			TableName: tableName,
			Fields:    make(map[string]string),
		}
		t := reflect.TypeOf(e)
		for _, i := range reflect.VisibleFields(t) {
			model.Fields[i.Name] = i.Type.Name()
		}
		query := modelMigrator(&model)
		if err := queryExecutor(ctx, query); err != nil {
			return fmt.Errorf("failed to migrate table %s: %w", tableName, err)
		}
	}
	return nil
}

// Database-specific query executors
func (ar *ApplicationRepository) executeScyllaQuery(ctx context.Context, query string) error {
	return ar.db.(*scylla.ScyllaConnection).Session.Query(query).WithContext(ctx).Exec()
}

func (ar *ApplicationRepository) executePostgresQuery(ctx context.Context, query string) error {
	_, err := ar.db.(*postgres.PostgresConnection).DB.ExecContext(ctx, query)
	return err
}

func newPostgresRepository(connectionString string) (*ApplicationRepository, error) {
	repository, err := postgres.NewPostgresConnection(connectionString)
	if err != nil {
		return nil, err
	}

	ar := &ApplicationRepository{
		db:     repository,
		dbType: "postgres",

		ServiceTokens: postgres.NewPostgresServiceTokenRepository(repository.DB),
		Buckets:       postgres.NewPostgresBucketRepository(repository.DB),
		Files:         postgres.NewPostgresFileRepository(repository.DB),
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
		db:     repository,
		dbType: "scylla",

		ServiceTokens: scylla.NewScyllaServiceTokenRepository(repository.Session),
		Buckets:       scylla.NewScyllaBucketRepository(repository.Session),
		Files:         scylla.NewScyllaFileRepository(repository.Session),
	}
	log.Printf("Scylla repository created: %+v", ar)
	return ar, nil
}

func (r *ApplicationRepository) InitializeRepo(ctx context.Context, repo *ApplicationRepository) {
	repo.ServiceTokens.CreateIndices(ctx)
	repo.Buckets.CreateIndices(ctx)
	repo.Files.CreateIndices(ctx)
}
