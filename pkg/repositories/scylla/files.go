package scylla

import (
	"context"
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
)

type ScyllaFileRepository struct {
	session *gocql.Session
}

func NewScyllaFileRepository(session *gocql.Session) *ScyllaFileRepository {
	return &ScyllaFileRepository{session: session}
}

func (s *ScyllaFileRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).WithContext(ctx).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (s *ScyllaFileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	panic("implement me")
}

func (s *ScyllaFileRepository) GetFileByName(ctx context.Context, bucketID, name string) (*models.File, error) {
	panic("implement me")
}

func (s *ScyllaFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	panic("implement me")
}

func (s *ScyllaFileRepository) UpdateFile(ctx context.Context, file *models.File) error {
	panic("implement me")
}

func (s *ScyllaFileRepository) DeleteFile(ctx context.Context, id string) error {
	panic("implement me")
}

func (s *ScyllaFileRepository) ListFiles(ctx context.Context, bucketID string) ([]*models.File, error) {
	panic("implement me")
}
