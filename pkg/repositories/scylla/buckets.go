package scylla

import (
	"log"

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
	indexQueries := []string{}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}
func (s *ScyllaBucketRepository) GetBucketByID(id string) (*models.Bucket, error) {
	panic("implement me")
}

func (s *ScyllaBucketRepository) GetBucketByName(name string) (*models.Bucket, error) {
	panic("implement me")
}

func (s *ScyllaBucketRepository) CreateBucket(bucket *models.Bucket) error {
	panic("implement me")
}

func (s *ScyllaBucketRepository) UpdateBucket(bucket *models.Bucket) error {
	panic("implement me")
}

func (s *ScyllaBucketRepository) DeleteBucket(id string) error {
	panic("implement me")
}

func (s *ScyllaBucketRepository) ListBuckets() ([]*models.Bucket, error) {
	panic("implement me")
}
