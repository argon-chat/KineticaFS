package scylla

import (
	"context"
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
)

// ddl
// table fileblob
// (
//     id        text primary key,
//     createdat timestamp,
//     fileid    text,
//     updatedat timestamp
// )

type ScyllaFileBlobRepository struct {
	session *gocql.Session
}

func NewScyllaFileBlobRepository(session *gocql.Session) *ScyllaFileBlobRepository {
	return &ScyllaFileBlobRepository{session: session}
}

func (s *ScyllaFileBlobRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS fileblob_fileid_idx ON fileblob (fileid)",
		"ALTER TABLE fileblob WITH default_time_to_live = 600",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).WithContext(ctx).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (s *ScyllaFileBlobRepository) CreateFileBlob(ctx context.Context, blob *models.FileBlob) (*models.FileBlob, error) {
	blob.ID = gocql.TimeUUID().String()
	query := "INSERT INTO fileblob (id, createdat, fileid, updatedat) VALUES (?, ?, ?, ?)"
	if err := s.session.Query(query, blob.ID, blob.CreatedAt, blob.FileID, blob.UpdatedAt).WithContext(ctx).Exec(); err != nil {
		return nil, err
	}
	return blob, nil
}

func (s *ScyllaFileBlobRepository) GetFileBlobByID(ctx context.Context, id string) (*models.FileBlob, error) {
	query := "SELECT id, createdat, fileid, hex, updatedat FROM fileblob WHERE id = ?"
	row := s.session.Query(query, id).WithContext(ctx)
	var blob models.FileBlob
	err := row.Scan(&blob.ID, &blob.CreatedAt, &blob.FileID, &blob.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &blob, nil
}

func (s *ScyllaFileBlobRepository) DeleteFileBlobByID(ctx context.Context, id string) error {
	query := "DELETE FROM fileblob WHERE id = ?"
	return s.session.Query(query, id).WithContext(ctx).Exec()
}
