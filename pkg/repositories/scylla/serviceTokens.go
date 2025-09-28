package scylla

import (
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

func (s *ScyllaServiceTokenRepository) GetServiceTokenById(id string) (*models.ServiceToken, error) {
	panic("implement me")
}

func (s *ScyllaServiceTokenRepository) GetServiceTokenByName(name string) (*models.ServiceToken, error) {
	query := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where name = ?", name)
	var token models.ServiceToken
	var tokenType int8
	if err := query.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt); err != nil {
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
