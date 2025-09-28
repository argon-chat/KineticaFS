package postgres

import (
	"database/sql"
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresFileRepository struct {
	session *sql.DB
}

func NewPostgresFileRepository(session *sql.DB) *PostgresFileRepository {
	return &PostgresFileRepository{session: session}
}

func (s *PostgresFileRepository) CreateIndices() {
	indexQueries := []string{}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if _, err := s.session.Exec(indexQuery); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
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
