package postgres

import (
	"database/sql"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresFileRepostiory struct {
	session *sql.DB
}

func (p PostgresFileRepostiory) GetFileByID(id string) (*models.File, error) {
	panic("implement me")
}

func (p PostgresFileRepostiory) GetFileByName(bucketID, name string) (*models.File, error) {
	panic("implement me")
}

func (p PostgresFileRepostiory) CreateFile(file *models.File) error {
	panic("implement me")
}

func (p PostgresFileRepostiory) UpdateFile(file *models.File) error {
	panic("implement me")
}

func (p PostgresFileRepostiory) DeleteFile(id string) error {
	panic("implement me")
}

func (p PostgresFileRepostiory) ListFiles(bucketID string) ([]*models.File, error) {
	panic("implement me")
}
