package repositories

import (
	"fmt"
	"reflect"

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

func (ar *ApplicationRepository) Migrate() error {
	switch ar.dbType {
	case "scylla":
		return ar.migrateDatabase(migrateScyllaModel, ar.executeScyllaQuery)
	case "postgres":
		return ar.migrateDatabase(migratePostgresModel, ar.executePostgresQuery)
	default:
		return fmt.Errorf("unsupported database type: %s", ar.dbType)
	}
}

func (ar *ApplicationRepository) migrateDatabase(
	modelMigrator func(*migration) string,
	queryExecutor func(string) error,
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
		if err := queryExecutor(query); err != nil {
			return fmt.Errorf("failed to migrate table %s: %w", tableName, err)
		}
	}
	return nil
}

// Database-specific query executors
func (ar *ApplicationRepository) executeScyllaQuery(query string) error {
	return ar.db.(*scylla.ScyllaConnection).Session.Query(query).Exec()
}

func (ar *ApplicationRepository) executePostgresQuery(query string) error {
	_, err := ar.db.(*postgres.PostgresConnection).DB.Exec(query)
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
	return ar, nil
}
