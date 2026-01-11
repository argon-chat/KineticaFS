package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/google/uuid"
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
	blob.ID = uuid.NewString()
	now := time.Now().UTC()
	blob.CreatedAt = now
	blob.UpdatedAt = now
	query := "INSERT INTO file_blob (id, created_at, file_id, updated_at) VALUES ($1, $2, $3, $4)"
	_, err := p.session.ExecContext(ctx, query, blob.ID, blob.CreatedAt, blob.FileID, blob.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func (p *PostgresFileBlobRepository) GetFileBlobByID(ctx context.Context, id string) (*models.FileBlob, error) {
	query := "SELECT id, created_at, file_id, updated_at FROM file_blob WHERE id = $1"
	row := p.session.QueryRowContext(ctx, query, id)
	var blob models.FileBlob
	err := row.Scan(&blob.ID, &blob.CreatedAt, &blob.FileID, &blob.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("file blob not found")
	}
	if err != nil {
		return nil, err
	}
	return &blob, nil
}

func (p *PostgresFileBlobRepository) DeleteFileBlobByID(ctx context.Context, id string) error {
	query := "DELETE FROM file_blob WHERE id = $1"
	_, err := p.session.ExecContext(ctx, query, id)
	return err
}

func (p *PostgresFileBlobRepository) GetAllFileBlobs(ctx context.Context) ([]*models.FileBlob, error) {
	query := "SELECT id, created_at, file_id, updated_at FROM file_blob"
	rows, err := p.session.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blobs []*models.FileBlob
	for rows.Next() {
		var blob models.FileBlob
		if err := rows.Scan(&blob.ID, &blob.CreatedAt, &blob.FileID, &blob.UpdatedAt); err != nil {
			return nil, err
		}
		blobs = append(blobs, &blob)
	}
	return blobs, rows.Err()
}
