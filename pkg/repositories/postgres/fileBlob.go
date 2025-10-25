package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresFileBlobRepository struct {
	session *sql.DB
}

func NewPostgresFileBlobRepository(session *sql.DB) *PostgresFileBlobRepository {
	return &PostgresFileBlobRepository{session: session}
}

func (p *PostgresFileBlobRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{
		"create index if not exists file_blob_file_id_idx on file_blob (file_id)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if _, err := p.session.ExecContext(ctx, indexQuery); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (p *PostgresFileBlobRepository) CreateFileBlob(ctx context.Context, blob *models.FileBlob) (*models.FileBlob, error) {
	return nil, nil
}

func (p *PostgresFileBlobRepository) GetFileBlobByID(ctx context.Context, id string) (*models.FileBlob, error) {
	return nil, nil
}

func (p *PostgresFileBlobRepository) DeleteFileBlobByID(ctx context.Context, id string) error {
	return nil
}
