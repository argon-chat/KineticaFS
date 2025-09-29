package scylla

import (
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
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

func (s *ScyllaServiceTokenRepository) CreateIndices() {
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS servicetoken_name_idx ON servicetoken (name)",
		"CREATE INDEX IF NOT EXISTS servicetoken_accesskey_idx ON servicetoken (accesskey)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if err := s.session.Query(indexQuery).Exec(); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (s *ScyllaServiceTokenRepository) GetServiceTokenByAccessKey(accessKey string) (*models.ServiceToken, error) {
	query := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where accesskey = ?", accessKey)
	var token models.ServiceToken
	var tokenType int8
	if err := query.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *ScyllaServiceTokenRepository) GetAllServiceTokens() ([]*models.ServiceToken, error) {
	var tokens []*models.ServiceToken
	iter := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken").Iter()
	var token models.ServiceToken
	var tokenType int8
	for iter.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt) {
		token.TokenType = models.TokenType(tokenType)
		newToken := &models.ServiceToken{
			Name:      token.Name,
			AccessKey: token.AccessKey,
			TokenType: token.TokenType,
		}
		newToken.ID = token.ID
		newToken.CreatedAt = token.CreatedAt
		newToken.UpdatedAt = token.UpdatedAt
		tokens = append(tokens, newToken)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *ScyllaServiceTokenRepository) GetServiceTokenById(id string) (*models.ServiceToken, error) {
	query := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where id = ?", id)
	var token models.ServiceToken
	var tokenType int8
	if err := query.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *ScyllaServiceTokenRepository) GetServiceTokenByName(name string) (*models.ServiceToken, error) {
	query := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where name = ?", name)
	var token models.ServiceToken
	var tokenType int8
	if err := query.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *ScyllaServiceTokenRepository) CreateServiceToken(token *models.ServiceToken) error {
	query := s.session.Query(
		"insert into servicetoken (id, name, accesskey, tokentype, createdat, updatedat) values (?, ?, ?, ?, ?, ?)",
		token.ID, token.Name, token.AccessKey, int8(token.TokenType), token.CreatedAt, token.UpdatedAt)
	return query.Exec()
}

func (s *ScyllaServiceTokenRepository) RevokeServiceToken(id string) error {
	panic("implement me")
}
