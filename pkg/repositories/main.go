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
		return ar.migrateScylla()
	case "postgres":
		return ar.migratePostgres()
	default:
		return fmt.Errorf("unsupported database type: %s", ar.dbType)
	}
}

func (ar *ApplicationRepository) migrateScylla() error {
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
		query := migrateScyllaModel(&model)
		if err := ar.db.(*scylla.ScyllaConnection).Session.Query(query).Exec(); err != nil {
			return fmt.Errorf("failed to migrate table %s: %w", tableName, err)
		}
	}
	return nil
}

func (ar *ApplicationRepository) migratePostgres() error {
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
		query := migratePostgresModel(&model)
		if _, err := ar.db.(*postgres.PostgresConnection).DB.Exec(query); err != nil {
			return fmt.Errorf("failed to migrate table %s: %w", tableName, err)
		}
	}
	return nil
}

func newPostgresRepository(connectionString string) (*ApplicationRepository, error) {
	repository, err := postgres.NewPostgresConnection(connectionString)
	if err != nil {
		return nil, err
	}

	ar := &ApplicationRepository{
		db:     repository,
		dbType: "postgres",

		ServiceTokens: nil,
		Buckets:       nil,
		Files:         nil,
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

		ServiceTokens: nil,
		Buckets:       nil,
		Files:         nil,
	}
	return ar, nil
}
