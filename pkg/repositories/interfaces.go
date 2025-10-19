package repositories

import (
	"context"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type IRepository interface {
	CreateIndices(ctx context.Context)
}

type IBucketRepository interface {
	IRepository
	GetBucketByID(ctx context.Context, id string) (*models.Bucket, error)
	GetBucketByName(ctx context.Context, name string) (*models.Bucket, error)
	CreateBucket(ctx context.Context, bucket *models.Bucket) error
	UpdateBucket(ctx context.Context, bucket *models.Bucket) error
	DeleteBucket(ctx context.Context, id string) error
	ListBuckets(ctx context.Context) ([]*models.Bucket, error)
}

type IFileRepository interface {
	IRepository
	GetFileByID(ctx context.Context, id string) (*models.File, error)
	GetFileByName(ctx context.Context, bucketID, name string) (*models.File, error)
	CreateFile(ctx context.Context, file *models.File) error
	UpdateFile(ctx context.Context, file *models.File) error
	DeleteFile(ctx context.Context, id string) error
	ListFiles(ctx context.Context, bucketID string) ([]*models.File, error)
}

type IServiceTokenRepository interface {
	IRepository
	GetAllServiceTokens(ctx context.Context) ([]*models.ServiceToken, error)
	GetServiceTokenById(ctx context.Context, id string) (*models.ServiceToken, error)
	GetServiceTokenByAccessKey(ctx context.Context, accessKey string) (*models.ServiceToken, error)
	GetServiceTokenByName(ctx context.Context, name string) (*models.ServiceToken, error)
	CreateServiceToken(ctx context.Context, token *models.ServiceToken) error
	RevokeServiceToken(ctx context.Context, id string) error
}

type IFileBlobRepository interface {
	IRepository
	CreateFileBlob(ctx context.Context, blob *models.FileBlob) (*models.FileBlob, error)
	GetFileBlobByID(ctx context.Context, id string) (*models.FileBlob, error)
	DeleteFileBlobByID(ctx context.Context, id string) error
}
