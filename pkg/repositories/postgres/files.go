package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresFileRepository struct {
	session *sql.DB
}

func NewPostgresFileRepository(session *sql.DB) *PostgresFileRepository {
	return &PostgresFileRepository{session: session}
}

func (p *PostgresFileRepository) CreateIndices(ctx context.Context) {
}

func (p *PostgresFileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	query := "SELECT id, bucket_id, checksum, content_type, created_at, file_size, file_size_limit, finalized, metadata, name, path, updated_at FROM file WHERE id = $1"
	row := p.session.QueryRowContext(ctx, query, id)
	var file models.File
	err := row.Scan(&file.ID, &file.BucketID, &file.Checksum, &file.ContentType, &file.CreatedAt, &file.FileSize, &file.FileSizeLimit, &file.Finalized, &file.Metadata, &file.Name, &file.Path, &file.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	refCount, err := p.GetFileReferenceCount(ctx, file.ID)
	if err != nil {
		log.Printf("Warning: Failed to get reference count for file %s: %v", file.ID, err)
		file.References = 0
	} else {
		file.References = refCount
	}

	return &file, nil
}

func (p *PostgresFileRepository) GetFileByName(ctx context.Context, name string) (*models.File, error) {
	query := "SELECT id, bucket_id, checksum, content_type, created_at, file_size, file_size_limit, finalized, metadata, name, path, updated_at FROM file WHERE name = $1"
	row := p.session.QueryRowContext(ctx, query, name)
	var file models.File
	err := row.Scan(&file.ID, &file.BucketID, &file.Checksum, &file.ContentType, &file.CreatedAt, &file.FileSize, &file.FileSizeLimit, &file.Finalized, &file.Metadata, &file.Name, &file.Path, &file.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	refCount, err := p.GetFileReferenceCount(ctx, file.ID)
	if err != nil {
		log.Printf("Warning: Failed to get reference count for file %s: %v", file.ID, err)
		file.References = 0
	} else {
		file.References = refCount
	}

	return &file, nil
}

func (p *PostgresFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	now := time.Now().UTC()
	file.CreatedAt = now
	file.UpdatedAt = now
	file.ID = file.Name
	query := `INSERT INTO file (id, bucket_id, name, file_size, file_size_limit, finalized, content_type, checksum, metadata, path, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := p.session.ExecContext(ctx, query, file.ID, file.BucketID, file.Name, file.FileSize, file.FileSizeLimit, file.Finalized, file.ContentType, file.Checksum, file.Metadata, file.Path, file.CreatedAt, file.UpdatedAt)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}

	// Initialize file counter
	counterQuery := `INSERT INTO file_counter (id, ref) VALUES ($1, 1) ON CONFLICT (id) DO UPDATE SET ref = file_counter.ref + 1`
	_, err = p.session.ExecContext(ctx, counterQuery, file.ID)
	if err != nil {
		log.Printf("Error updating file counter: %v", err)
		return err
	}

	return nil
}

func (p *PostgresFileRepository) UpdateFile(ctx context.Context, file *models.File) error {
	file.UpdatedAt = time.Now().UTC()
	query := `UPDATE file SET bucket_id = $1, finalized = $2, name = $3, file_size = $4, file_size_limit = $5, content_type = $6, checksum = $7, metadata = $8, path = $9, updated_at = $10 WHERE id = $11`
	_, err := p.session.ExecContext(ctx, query, file.BucketID, file.Finalized, file.Name, file.FileSize, file.FileSizeLimit, file.ContentType, file.Checksum, file.Metadata, file.Path, file.UpdatedAt, file.ID)
	if err != nil {
		log.Printf("Error updating file: %v", err)
		return err
	}
	return nil
}

func (p *PostgresFileRepository) DeleteFile(ctx context.Context, id string) error {
	query := "DELETE FROM file WHERE id = $1"
	_, err := p.session.ExecContext(ctx, query, id)
	return err
}

func (p *PostgresFileRepository) ListFiles(ctx context.Context, bucketID string) ([]*models.File, error) {
	query := "SELECT id, bucket_id, checksum, content_type, created_at, file_size, file_size_limit, finalized, metadata, name, path, updated_at FROM file WHERE bucket_id = $1"
	rows, err := p.session.QueryContext(ctx, query, bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.BucketID, &file.Checksum, &file.ContentType, &file.CreatedAt, &file.FileSize, &file.FileSizeLimit, &file.Finalized, &file.Metadata, &file.Name, &file.Path, &file.UpdatedAt); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Populate reference counts
	for _, file := range files {
		refCount, err := p.GetFileReferenceCount(ctx, file.ID)
		if err != nil {
			log.Printf("Warning: Failed to get reference count for file %s: %v", file.ID, err)
			file.References = 0
		} else {
			file.References = refCount
		}
	}

	return files, nil
}

func (p *PostgresFileRepository) AtomicIncrement(ctx context.Context, id string) error {
	query := "INSERT INTO file_counter (id, ref) VALUES ($1, 1) ON CONFLICT (id) DO UPDATE SET ref = file_counter.ref + 1"
	_, err := p.session.ExecContext(ctx, query, id)
	return err
}

func (p *PostgresFileRepository) AtomicDecrement(ctx context.Context, id string) error {
	query := "UPDATE file_counter SET ref = ref - 1 WHERE id = $1"
	_, err := p.session.ExecContext(ctx, query, id)
	return err
}

func (p *PostgresFileRepository) GetFileReferenceCount(ctx context.Context, fileID string) (int64, error) {
	var refCount int64
	query := "SELECT ref FROM file_counter WHERE id = $1"
	err := p.session.QueryRowContext(ctx, query, fileID).Scan(&refCount)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return refCount, nil
}
