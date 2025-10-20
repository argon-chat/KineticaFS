package scylla

import (
	"context"
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
)

// DDL
// table file
// (
//     id          text primary key,
//     bucketid    text,
//     checksum    text,
//     contenttype text,
//     createdat   timestamp,
//     filesize    int,
//     metadata    text,
//     name        text,
//     path        text,
//     updatedat   timestamp
// )

type ScyllaFileRepository struct {
	session *gocql.Session
}

func NewScyllaFileRepository(session *gocql.Session) *ScyllaFileRepository {
	return &ScyllaFileRepository{session: session}
}

func (s *ScyllaFileRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).WithContext(ctx).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (s *ScyllaFileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	query := "SELECT id, bucketid, checksum, contenttype, createdat, filesize, metadata, name, path, updatedat FROM file WHERE id = ?"
	row := s.session.Query(query, id).WithContext(ctx)
	var file models.File
	err := row.Scan(&file.ID, &file.BucketID, &file.Checksum, &file.ContentType, &file.CreatedAt, &file.FileSize, &file.Metadata, &file.Name, &file.Path, &file.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *ScyllaFileRepository) GetFileByName(ctx context.Context, bucketID, name string) (*models.File, error) {
	panic("implement me")
}

func (s *ScyllaFileRepository) CreateFile(ctx context.Context, file *models.File) error {
	file.CreatedAt = time.Now().UTC()
	file.UpdatedAt = file.CreatedAt
	file.ID = file.Name
	query := `INSERT INTO file (id, bucketid, name, filesize, contenttype, createdat, updatedat) VALUES (?, ?, ?, ?, ?, ?, ?)`
	if err := s.session.Query(query, file.ID, file.BucketID, file.Name, file.FileSize, file.ContentType, file.CreatedAt, file.UpdatedAt).WithContext(ctx).Exec(); err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	return nil
}

func (s *ScyllaFileRepository) UpdateFile(ctx context.Context, file *models.File) error {
	panic("implement me")
}

func (s *ScyllaFileRepository) DeleteFile(ctx context.Context, id string) error {
	panic("implement me")
}

func (s *ScyllaFileRepository) ListFiles(ctx context.Context, bucketID string) ([]*models.File, error) {
	panic("implement me")
}
