package repositories

type IBucketRepository interface {
	GetBucketByID(id string) (*Bucket, error)
	GetBucketByName(name string) (*Bucket, error)
	CreateBucket(bucket *Bucket) error
	UpdateBucket(bucket *Bucket) error
	DeleteBucket(id string) error
	ListBuckets() ([]*Bucket, error)
}

type IFileRepository interface {
	GetFileByID(id string) (*File, error)
	GetFileByName(bucketID, name string) (*File, error)
	CreateFile(file *File) error
	UpdateFile(file *File) error
	DeleteFile(id string) error
	ListFiles(bucketID string) ([]*File, error)
}

type IServiceTokenRepository interface {
	GetServiceTokenById(id string) (*ServiceToken, error)
	GetServiceTokenByName(name string) (*ServiceToken, error)
	CreateServiceToken(token *ServiceToken) error
	UpdateServiceToken(token *ServiceToken) error
	RevokeServiceToken(id string) error
}
