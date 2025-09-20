package repositories

import "github.com/argon-chat/KineticaFS/pkg/models"

type IBucketRepository interface {
	GetBucketByID(id string) (*models.Bucket, error)
	GetBucketByName(name string) (*models.Bucket, error)
	CreateBucket(bucket *models.Bucket) error
	UpdateBucket(bucket *models.Bucket) error
	DeleteBucket(id string) error
	ListBuckets() ([]*models.Bucket, error)
}

type IFileRepository interface {
	GetFileByID(id string) (*models.File, error)
	GetFileByName(bucketID, name string) (*models.File, error)
	CreateFile(file *models.File) error
	UpdateFile(file *models.File) error
	DeleteFile(id string) error
	ListFiles(bucketID string) ([]*models.File, error)
}

type IServiceTokenRepository interface {
	GetServiceTokenById(id string) (*models.ServiceToken, error)
	GetServiceTokenByName(name string) (*models.ServiceToken, error)
	CreateServiceToken(token *models.ServiceToken) error
	RevokeServiceToken(id string) error
}
