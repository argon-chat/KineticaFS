package scylla

import (
	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
)

type ScyllaBucketRepository struct {
	session *gocql.Session
}

func (s ScyllaBucketRepository) GetBucketByID(id string) (*models.Bucket, error) {
	panic("implement me")
}

func (s ScyllaBucketRepository) GetBucketByName(name string) (*models.Bucket, error) {
	panic("implement me")
}

func (s ScyllaBucketRepository) CreateBucket(bucket *models.Bucket) error {
	panic("implement me")
}

func (s ScyllaBucketRepository) UpdateBucket(bucket *models.Bucket) error {
	panic("implement me")
}

func (s ScyllaBucketRepository) DeleteBucket(id string) error {
	panic("implement me")
}

func (s ScyllaBucketRepository) ListBuckets() ([]*models.Bucket, error) {
	panic("implement me")
}
