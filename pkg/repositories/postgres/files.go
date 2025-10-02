package postgres

import (
	"context"
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

func (s *PostgresFileRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if _, err := s.session.ExecContext(ctx, indexQuery); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}
func (p *PostgresFileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	panic("implement me")
}

func (p *PostgresFileRepository) GetFileByName(ctx context.Context, bucketID, name string) (*models.File, error) {
	panic("implement me")
}

func (p *PostgresFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	panic("implement me")
}

func (p *PostgresFileRepository) UpdateFile(ctx context.Context, file *models.File) error {
	panic("implement me")
}

func (p *PostgresFileRepository) DeleteFile(ctx context.Context, id string) error {
	panic("implement me")
}

func (p *PostgresFileRepository) ListFiles(ctx context.Context, bucketID string) ([]*models.File, error) {
	panic("implement me")
}
