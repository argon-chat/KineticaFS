package scylla

import (
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
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

func (s *ScyllaBucketRepository) CreateIndices() {
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS bucket_name_idx ON bucket (name)",
		"CREATE INDEX IF NOT EXISTS bucket_region_idx ON bucket (region)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}
func (s *ScyllaBucketRepository) GetBucketByID(id string) (*models.Bucket, error) {
	query := s.session.Query("select id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat from bucket where id = ?", id)
	var bucket models.Bucket
	var storageType int8
	var useSSL bool
	if err := query.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &useSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	bucket.StorageType = models.StorageType(storageType)
	bucket.UseSSL = useSSL
	return &bucket, nil
}

func (s *ScyllaBucketRepository) GetBucketByName(name string) (*models.Bucket, error) {
	query := s.session.Query("select id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat from bucket where name = ?", name)
	var bucket models.Bucket
	var storageType int8
	var useSSL bool
	if err := query.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &useSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	bucket.StorageType = models.StorageType(storageType)
	bucket.UseSSL = useSSL
	return &bucket, nil
}

func (s *ScyllaBucketRepository) CreateBucket(bucket *models.Bucket) error {
	now := time.Now()
	bucket.CreatedAt = now
	bucket.UpdatedAt = now
	query := s.session.Query(
		"insert into bucket (id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		bucket.ID, bucket.Name, bucket.Region, bucket.Endpoint, bucket.S3Provider, bucket.AccessKey, bucket.SecretKey, bucket.StorageType, bucket.UseSSL, bucket.CustomConfig, bucket.CreatedAt, bucket.UpdatedAt)
	return query.Exec()
}

func (s *ScyllaBucketRepository) UpdateBucket(bucket *models.Bucket) error {
	query := s.session.Query(
		"update bucket set name = ?, region = ?, endpoint = ?, s3provider = ?, accesskey = ?, secretkey = ?, storagetype = ?, usessl = ?, customconfig = ?, updatedat = ? where id = ?",
		bucket.Name, bucket.Region, bucket.Endpoint, bucket.S3Provider, bucket.AccessKey, bucket.SecretKey, bucket.StorageType, bucket.UseSSL, bucket.CustomConfig, time.Now(), bucket.ID)
	return query.Exec()
}

func (s *ScyllaBucketRepository) DeleteBucket(id string) error {
	return s.session.Query("delete from bucket where id = ?", id).Exec()
}

func (s *ScyllaBucketRepository) ListBuckets() ([]*models.Bucket, error) {
	var buckets []*models.Bucket
	iter := s.session.Query("select id, name, region, endpoint, s3provider, accesskey, secretkey, storagetype, usessl, customconfig, createdat, updatedat from bucket").Iter()
	var bucket models.Bucket
	var storageType int8
	var useSSL bool
	for iter.Scan(&bucket.ID, &bucket.Name, &bucket.Region, &bucket.Endpoint, &bucket.S3Provider, &bucket.AccessKey, &bucket.SecretKey, &storageType, &useSSL, &bucket.CustomConfig, &bucket.CreatedAt, &bucket.UpdatedAt) {
		bucket.StorageType = models.StorageType(storageType)
		bucket.UseSSL = useSSL
		newBucket := &models.Bucket{
			ApplicationModel: bucket.ApplicationModel,
			Name:             bucket.Name,
			Region:           bucket.Region,
			Endpoint:         bucket.Endpoint,
			S3Provider:       bucket.S3Provider,
			AccessKey:        bucket.AccessKey,
			SecretKey:        bucket.SecretKey,
			StorageType:      bucket.StorageType,
			UseSSL:           bucket.UseSSL,
			CustomConfig:     bucket.CustomConfig,
		}
		buckets = append(buckets, newBucket)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return buckets, nil
}
