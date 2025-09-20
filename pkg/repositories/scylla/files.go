package scylla

import (
	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
)

type ScyllaFileRepository struct {
	session *gocql.Session
}

func (s ScyllaFileRepository) GetFileByID(id string) (*models.File, error) {
	panic("implement me")
}

func (s ScyllaFileRepository) GetFileByName(bucketID, name string) (*models.File, error) {
	panic("implement me")
}

func (s ScyllaFileRepository) CreateFile(file *models.File) error {
	panic("implement me")
}

func (s ScyllaFileRepository) UpdateFile(file *models.File) error {
	panic("implement me")
}

func (s ScyllaFileRepository) DeleteFile(id string) error {
	panic("implement me")
}

func (s ScyllaFileRepository) ListFiles(bucketID string) ([]*models.File, error) {
	panic("implement me")
}
