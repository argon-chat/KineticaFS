package scylla

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// DDL
// bucket
// (
//     id           text primary key,
//     accesskey    text,
//     createdat    timestamp,
//     customconfig text,
//     endpoint     text,
//     name         text,
//     region       text,
//     s3provider   text,
//     secretkey    text,
//     storagetype  int,
//     updatedat    timestamp,
//     usessl       boolean
// )

type ScyllaBucketRepository struct {
	session *gocql.Session
}

func NewScyllaBucketRepository(session *gocql.Session) *ScyllaBucketRepository {
	return &ScyllaBucketRepository{session: session}
}

func (s *ScyllaBucketRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS bucket_name_idx ON bucket (name)",
		"CREATE INDEX IF NOT EXISTS bucket_region_idx ON bucket (region)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).WithContext(ctx).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}
func (s *ScyllaBucketRepository) GetBucketByID(ctx context.Context, id string) (*models.Bucket, error) {
	query := s.session.Query("select id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat from bucket where id = ?", id).
		WithContext(ctx)
	var bucket models.Bucket
	var storageType int8
	if err := query.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &bucket.UseSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	bucket.StorageType = models.StorageType(storageType)
	return &bucket, nil
}

func (s *ScyllaBucketRepository) GetBucketByName(ctx context.Context, name string) (*models.Bucket, error) {
	query := s.session.Query("select id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat from bucket where name = ?", name).
		WithContext(ctx)
	var bucket models.Bucket
	var storageType int8
	if err := query.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &bucket.UseSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	bucket.StorageType = models.StorageType(storageType)
	return &bucket, nil
}

func (s *ScyllaBucketRepository) CreateBucket(ctx context.Context, bucket *models.Bucket) error {
	now := time.Now()
	bucket.CreatedAt = now
	bucket.UpdatedAt = now
	bucket.ID = uuid.NewString()
	query := s.session.Query(
		"insert into bucket (id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		bucket.ID, bucket.Name, bucket.Region, bucket.Endpoint, bucket.S3Provider, bucket.AccessKey, bucket.SecretKey, bucket.StorageType, bucket.UseSSL, bucket.CustomConfig, bucket.CreatedAt, bucket.UpdatedAt).
		WithContext(ctx)
	return query.Exec()
}

func (s *ScyllaBucketRepository) UpdateBucket(ctx context.Context, bucket *models.Bucket) error {
	query := s.session.Query(
		"update bucket set name = ?, region = ?, endpoint = ?, s3provider = ?, accesskey = ?, secretkey = ?, storagetype = ?, usessl = ?, customconfig = ?, updatedat = ? where id = ?",
		bucket.Name, bucket.Region, bucket.Endpoint, bucket.S3Provider, bucket.AccessKey, bucket.SecretKey, bucket.StorageType, bucket.UseSSL, bucket.CustomConfig, time.Now(), bucket.ID).
		WithContext(ctx)
	return query.Exec()
}

func (s *ScyllaBucketRepository) DeleteBucket(ctx context.Context, id string) error {
	return s.session.Query("delete from bucket where id = ?", id).WithContext(ctx).Exec()
}

func (s *ScyllaBucketRepository) ListBuckets(ctx context.Context) ([]*models.Bucket, error) {
	iter := s.session.Query("select id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat from bucket").WithContext(ctx).Iter()

	estimatedSize := iter.NumRows()
	buckets := make([]*models.Bucket, 0, estimatedSize)
	scanner := iter.Scanner()

	for scanner.Next() {
		bucket := &models.Bucket{}
		var storageType int8

		if err := scanner.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &bucket.UseSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt); err != nil {
			return nil, err
		}

		bucket.StorageType = models.StorageType(storageType)
		buckets = append(buckets, bucket)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return buckets, nil
}
