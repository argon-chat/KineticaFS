package scylla

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// DDL
// servicetoken
// (
//     id        text primary key,
//     accesskey text,
//     createdat timestamp,
//     name      text,
//     tokentype int,
//     updatedat timestamp
// )

type ScyllaServiceTokenRepository struct {
	session *gocql.Session
}

func NewScyllaServiceTokenRepository(session *gocql.Session) *ScyllaServiceTokenRepository {
	return &ScyllaServiceTokenRepository{session: session}
}

func (s *ScyllaServiceTokenRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS servicetoken_name_idx ON servicetoken (name)",
		"CREATE INDEX IF NOT EXISTS servicetoken_accesskey_idx ON servicetoken (accesskey)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).WithContext(ctx).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (s *ScyllaServiceTokenRepository) GetServiceTokenByAccessKey(ctx context.Context, accessKey string) (*models.ServiceToken, error) {
	query := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where accesskey = ?", accessKey).
		WithContext(ctx)
	var token models.ServiceToken
	var tokenType int8
	if err := query.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *ScyllaServiceTokenRepository) GetAllServiceTokens(ctx context.Context) ([]*models.ServiceToken, error) {
	iter := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken").
		WithContext(ctx).
		Iter()

	estimatedSize := iter.NumRows()
	tokens := make([]*models.ServiceToken, 0, estimatedSize)
	for {
		token := &models.ServiceToken{}
		var tokenType int8
		if !iter.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt) {
			break
		}
		token.TokenType = models.TokenType(tokenType)
		tokens = append(tokens, token)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *ScyllaServiceTokenRepository) GetServiceTokenById(ctx context.Context, id string) (*models.ServiceToken, error) {
	query := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where id = ?", id).
		WithContext(ctx)
	var token models.ServiceToken
	var tokenType int8
	if err := query.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *ScyllaServiceTokenRepository) GetServiceTokenByName(ctx context.Context, name string) (*models.ServiceToken, error) {
	query := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where name = ?", name).WithContext(ctx)
	var token models.ServiceToken
	var tokenType int8
	if err := query.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *ScyllaServiceTokenRepository) CreateServiceToken(ctx context.Context, token *models.ServiceToken) error {
	token.ID = uuid.NewString()
	token.CreatedAt = time.Now()
	token.UpdatedAt = token.CreatedAt
	query := s.session.Query(
		"insert into servicetoken (id, name, accesskey, tokentype, createdat, updatedat) values (?, ?, ?, ?, ?, ?)",
		token.ID, token.Name, token.AccessKey, int8(token.TokenType), token.CreatedAt, token.UpdatedAt).WithContext(ctx)
	return query.Exec()
}

func (s *ScyllaServiceTokenRepository) RevokeServiceToken(ctx context.Context, id string) error {
	return s.session.Query("delete from servicetoken where id = ?", id).WithContext(ctx).Exec()
}
