package postgres

import (
	"database/sql"
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

// DDL
// servicetoken
// (
//     name      text,
//     accesskey text,
//     tokentype text,
//     id        text not null
//         primary key,
//     createdat timestamp,
//     updatedat timestamp
// )

type PostgresServiceTokenRepository struct {
	session *sql.DB
}

func NewPostgresServiceTokenRepository(session *sql.DB) *PostgresServiceTokenRepository {
	return &PostgresServiceTokenRepository{session: session}
}

func (s *PostgresServiceTokenRepository) CreateIndices() {
	indexQueries := []string{
		"create index if not exists servicetoken_name_idx on servicetoken (name)",
		"create index if not exists servicetoken_accesskey_idx on servicetoken (accesskey)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if _, err := s.session.Exec(indexQuery); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (p *PostgresServiceTokenRepository) GetServiceTokenById(id string) (*models.ServiceToken, error) {
	row := p.session.QueryRow("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where id = $1", id)
	var token models.ServiceToken
	var tokenType int8
	err := row.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *PostgresServiceTokenRepository) GetAllServiceTokens() ([]*models.ServiceToken, error) {
	rows, err := s.session.Query("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tokens []*models.ServiceToken
	for rows.Next() {
		var token models.ServiceToken
		var tokenType int8
		err := rows.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt)
		if err != nil {
			return nil, err
		}
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}
func (s *PostgresServiceTokenRepository) GetServiceTokenByAccessKey(accessKey string) (*models.ServiceToken, error) {
	row := s.session.QueryRow("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where accesskey = $1", accessKey)
	var token models.ServiceToken
	var tokenType int8
	err := row.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (p *PostgresServiceTokenRepository) GetServiceTokenByName(name string) (*models.ServiceToken, error) {
	row := p.session.QueryRow("select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where name = $1", name)
	var token models.ServiceToken
	var tokenType int8
	err := row.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (p *PostgresServiceTokenRepository) CreateServiceToken(token *models.ServiceToken) error {
	_, err := p.session.Exec(
		"insert into servicetoken (id, name, accesskey, tokentype, createdat, updatedat) values ($1, $2, $3, $4, $5, $6)",
		token.ID, token.Name, token.AccessKey, int8(token.TokenType), token.CreatedAt, token.UpdatedAt)
	return err
}

func (p *PostgresServiceTokenRepository) RevokeServiceToken(id string) error {
	_, err := p.session.Exec("delete from servicetoken where id = $1", id)
	return err
}
