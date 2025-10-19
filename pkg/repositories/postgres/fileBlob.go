package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

// ddl
// table fileblob
// (
//     createdat timestamp,
//     updatedat timestamp,
//     fileid    text,
//     id        text not null
//         primary key
// )

type PostgresFileBlobRepository struct {
	session *sql.DB
}

func NewPostgresFileBlobRepository(session *sql.DB) *PostgresFileBlobRepository {
	return &PostgresFileBlobRepository{session: session}
}

func (p *PostgresFileBlobRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{
		"create index if not exists fileblob_fileid_idx on fileblob (fileid)",
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
