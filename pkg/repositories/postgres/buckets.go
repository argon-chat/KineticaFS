package postgres

import (
	"database/sql"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresBucketRepository struct {
	session *sql.DB
}

func NewPostgresBucketRepository(session *sql.DB) *PostgresBucketRepository {
	return &PostgresBucketRepository{session: session}
}

func (p *PostgresBucketRepository) GetBucketByID(id string) (*models.Bucket, error) {
	panic("implement me")
}

func (p *PostgresBucketRepository) GetBucketByName(name string) (*models.Bucket, error) {
	panic("implement me")
}

func (p *PostgresBucketRepository) CreateBucket(bucket *models.Bucket) error {
	panic("implement me")
}

func (p *PostgresBucketRepository) UpdateBucket(bucket *models.Bucket) error {
	panic("implement me")
}

func (p *PostgresBucketRepository) DeleteBucket(id string) error {
	panic("implement me")
}

func (p *PostgresBucketRepository) ListBuckets() ([]*models.Bucket, error) {
	panic("implement me")
}
