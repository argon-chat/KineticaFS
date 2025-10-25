package scylla

import (
	"context"
	"log"
	"time"

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
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS file_bucket_id_idx ON file (bucket_id)",
		"CREATE INDEX IF NOT EXISTS file_name_idx ON file (name)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).WithContext(ctx).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (s *ScyllaFileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	query := "SELECT id, bucket_id, checksum, content_type, created_at, file_size, file_size_limit, finalized, metadata, name, path, updated_at FROM file WHERE id = ?"
	row := s.session.Query(query, id).WithContext(ctx)
	var file models.File
	err := row.Scan(&file.ID, &file.BucketID, &file.Checksum, &file.ContentType, &file.CreatedAt, &file.FileSize, &file.FileSizeLimit, &file.Finalized, &file.Metadata, &file.Name, &file.Path, &file.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *ScyllaFileRepository) GetFileByName(ctx context.Context, bucketID, name string) (*models.File, error) {
	query := "SELECT id, bucket_id, checksum, content_type, created_at, file_size, file_size_limit, finalized, metadata, name, path, updated_at FROM file WHERE bucket_id = ? AND name = ? ALLOW FILTERING"
	row := s.session.Query(query, bucketID, name).WithContext(ctx)
	var file models.File
	err := row.Scan(&file.ID, &file.BucketID, &file.Checksum, &file.ContentType, &file.CreatedAt, &file.FileSize, &file.FileSizeLimit, &file.Finalized, &file.Metadata, &file.Name, &file.Path, &file.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *ScyllaFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	file.CreatedAt = time.Now().UTC()
	file.UpdatedAt = file.CreatedAt
	file.ID = file.Name
	query := `INSERT INTO file (id, bucket_id, name, file_size, file_size_limit, finalized, content_type, checksum, metadata, path, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if err := s.session.Query(query, file.ID, file.BucketID, file.Name, file.FileSize, file.FileSizeLimit, file.Finalized, file.ContentType, file.Checksum, file.Metadata, file.Path, file.CreatedAt, file.UpdatedAt).WithContext(ctx).Exec(); err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	return nil
}

func (s *ScyllaFileRepository) UpdateFile(ctx context.Context, file *models.File) error {
	file.UpdatedAt = time.Now().UTC()
	query := `UPDATE file SET bucket_id = ?, finalized = ?, name = ?, file_size = ?, file_size_limit = ?, content_type = ?, checksum = ?, metadata = ?, path = ?, updated_at = ? WHERE id = ?`
	if err := s.session.Query(query, file.BucketID, file.Finalized, file.Name, file.FileSize, file.FileSizeLimit, file.ContentType, file.Checksum, file.Metadata, file.Path, file.UpdatedAt, file.ID).WithContext(ctx).Exec(); err != nil {
		log.Printf("Error updating file: %v", err)
		return err
	}
	return nil
}

func (s *ScyllaFileRepository) DeleteFile(ctx context.Context, id string) error {
	panic("implement me")
}

func (s *ScyllaFileRepository) ListFiles(ctx context.Context, bucketID string) ([]*models.File, error) {
	panic("implement me")
}
