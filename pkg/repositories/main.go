package repositories

import (
	"fmt"

	"github.com/spf13/viper"
)

type ApplicationRepository struct {
	db     any
	dbType string

	ServiceTokens IServiceTokenRepository
	Buckets       IBucketRepository
	Files         IFileRepository
}

func NewApplicationRepository() (*ApplicationRepository, error) {
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
	return nil
}

func newPostgresRepository(connectionString string) (*ApplicationRepository, error) {
	return nil, nil
}

func newScyllaRepository(connectionString string) (*ApplicationRepository, error) {
	return nil, nil
}
