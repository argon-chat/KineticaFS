package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/google/uuid"
)

// DDL
// servicetoken
// (
//     name      text,
//     accesskey text,
//     tokentype integer,
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

func (s *PostgresServiceTokenRepository) CreateIndices(ctx context.Context) {
	indexQueries := []string{
		"create index if not exists servicetoken_name_idx on servicetoken (name)",
		"create index if not exists servicetoken_accesskey_idx on servicetoken (accesskey)",
	}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if _, err := s.session.ExecContext(ctx, indexQuery); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func (p *PostgresServiceTokenRepository) GetServiceTokenById(ctx context.Context, id string) (*models.ServiceToken, error) {
	row := p.session.QueryRowContext(ctx, "select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where id = $1", id)
	var token models.ServiceToken
	var tokenType int8
	err := row.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (s *PostgresServiceTokenRepository) GetAllServiceTokens(ctx context.Context) ([]*models.ServiceToken, error) {
	rows, err := s.session.QueryContext(ctx, "select id, name, accesskey, tokentype, createdat, updatedat from servicetoken")
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
		tokens = append(tokens, &token)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}
func (s *PostgresServiceTokenRepository) GetServiceTokenByAccessKey(ctx context.Context, accessKey string) (*models.ServiceToken, error) {
	row := s.session.QueryRowContext(ctx, "select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where accesskey = $1", accessKey)
	var token models.ServiceToken
	var tokenType int8
	err := row.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (p *PostgresServiceTokenRepository) GetServiceTokenByName(ctx context.Context, name string) (*models.ServiceToken, error) {
	row := p.session.QueryRowContext(ctx, "select id, name, accesskey, tokentype, createdat, updatedat from servicetoken where name = $1", name)
	var token models.ServiceToken
	var tokenType int8
	err := row.Scan(&token.ID, &token.Name, &token.AccessKey, &tokenType, &token.CreatedAt, &token.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	token.TokenType = models.TokenType(tokenType)
	return &token, nil
}

func (p *PostgresServiceTokenRepository) CreateServiceToken(ctx context.Context, token *models.ServiceToken) error {
	token.ID = uuid.NewString()
	token.CreatedAt = time.Now()
	token.UpdatedAt = token.CreatedAt
	_, err := p.session.ExecContext(
		ctx,
		"insert into servicetoken (id, name, accesskey, tokentype, createdat, updatedat) values ($1, $2, $3, $4, $5, $6)",
		token.ID, token.Name, token.AccessKey, int8(token.TokenType), token.CreatedAt, token.UpdatedAt)
	return err
}

func (p *PostgresServiceTokenRepository) RevokeServiceToken(ctx context.Context, id string) error {
	_, err := p.session.ExecContext(ctx, "delete from servicetoken where id = $1", id)
	return err
}
