package repositories

import "fmt"

type ApplicationRepository struct {
	db     any
	dbType string

	ServiceTokens IServiceTokenRepository
	Buckets       IBucketRepository
	Files         IFileRepository
}

func NewApplicationRepository(dbType, connectionString string) (*ApplicationRepository, error) {
	switch dbType {
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
