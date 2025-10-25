package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/google/uuid"
)

type PostgresBucketRepository struct {
	session *sql.DB
}

func NewPostgresBucketRepository(session *sql.DB) *PostgresBucketRepository {
	return &PostgresBucketRepository{session: session}
}

func (s *PostgresBucketRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{
		"create index if not exists bucket_name_idx on bucket (name)",
		"create index if not exists bucket_region_idx on bucket (region)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if _, err := s.session.ExecContext(ctx, indexQuery); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}
func (p *PostgresBucketRepository) GetBucketByID(ctx context.Context, id string) (*models.Bucket, error) {
	row := p.session.QueryRowContext(ctx, "select id, name, region, endpoint, s3_provider, access_key, secret_key, storage_type, use_ssl, custom_config, created_at, updated_at from bucket where id = $1", id)
	var bucket models.Bucket
	var storageType int8
	err := row.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &bucket.UseSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	bucket.StorageType = models.StorageType(storageType)
	return &bucket, nil
}

func (p *PostgresBucketRepository) GetBucketByName(ctx context.Context, name string) (*models.Bucket, error) {
	row := p.session.QueryRowContext(ctx, "select id, name, region, endpoint, s3_provider, access_key, secret_key, storage_type, use_ssl, custom_config, created_at, updated_at from bucket where name = $1", name)
	var bucket models.Bucket
	var storageType int8
	err := row.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &bucket.UseSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	bucket.StorageType = models.StorageType(storageType)
	return &bucket, nil
}

func (p *PostgresBucketRepository) CreateBucket(ctx context.Context, bucket *models.Bucket) error {
	now := time.Now()
	bucket.CreatedAt = now
	bucket.UpdatedAt = now
	bucket.ID = uuid.NewString()
	_, err := p.session.ExecContext(
		ctx,
		"insert into bucket (id, name, region, endpoint, s3_provider, access_key, secret_key, storage_type, use_ssl, custom_config, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		bucket.ID, bucket.Name, bucket.Region, bucket.Endpoint, bucket.S3Provider, bucket.AccessKey, bucket.SecretKey, bucket.StorageType, bucket.UseSSL, bucket.CustomConfig, bucket.CreatedAt, bucket.UpdatedAt)
	return err
}

func (p *PostgresBucketRepository) UpdateBucket(ctx context.Context, bucket *models.Bucket) error {
	_, err := p.session.ExecContext(
		ctx,
		"update bucket set name = $1, region = $2, endpoint = $3, s3_provider = $4, access_key = $5, secret_key = $6, storage_type = $7, use_ssl = $8, custom_config = $9, updated_at = $10 where id = $11",
		bucket.Name, bucket.Region, bucket.Endpoint, bucket.S3Provider, bucket.AccessKey, bucket.SecretKey, bucket.StorageType, bucket.UseSSL, bucket.CustomConfig, bucket.UpdatedAt, bucket.ID)
	return err
}

func (p *PostgresBucketRepository) DeleteBucket(ctx context.Context, id string) error {
	_, err := p.session.ExecContext(ctx, "delete from bucket where id = $1", id)
	return err
}

func (p *PostgresBucketRepository) ListBuckets(ctx context.Context) ([]*models.Bucket, error) {
	rows, err := p.session.QueryContext(ctx, "select id, name, region, endpoint, s3_provider, access_key, secret_key, storage_type, use_ssl, custom_config, created_at, updated_at from bucket")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var buckets []*models.Bucket
	for rows.Next() {
		bucket := &models.Bucket{}
		var storageType int8
		err := rows.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &bucket.UseSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bucket.StorageType = models.StorageType(storageType)
		buckets = append(buckets, bucket)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return buckets, nil
}
