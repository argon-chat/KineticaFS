package postgres

import (
	"database/sql"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresFileRepository struct {
	session *sql.DB
}

func NewPostgresFileRepository(session *sql.DB) *PostgresFileRepository {
	return &PostgresFileRepository{session: session}
}

func (p *PostgresFileRepository) GetFileByID(id string) (*models.File, error) {
	panic("implement me")
}

func (p *PostgresFileRepository) GetFileByName(bucketID, name string) (*models.File, error) {
	panic("implement me")
}

func (p *PostgresFileRepository) CreateFile(file *models.File) error {
	panic("implement me")
}

func (p *PostgresFileRepository) UpdateFile(file *models.File) error {
	panic("implement me")
}

func (p *PostgresFileRepository) DeleteFile(id string) error {
	panic("implement me")
}

func (p *PostgresFileRepository) ListFiles(bucketID string) ([]*models.File, error) {
	panic("implement me")
}
