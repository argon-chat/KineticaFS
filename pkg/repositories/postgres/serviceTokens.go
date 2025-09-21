package postgres

import (
	"database/sql"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresServiceTokenRepository struct {
	session *sql.DB
}

func NewPostgresServiceTokenRepository(session *sql.DB) *PostgresServiceTokenRepository {
	return &PostgresServiceTokenRepository{session: session}
}

func (p *PostgresServiceTokenRepository) GetServiceTokenById(id string) (*models.ServiceToken, error) {
	panic("implement me")
}

func (p *PostgresServiceTokenRepository) GetServiceTokenByName(name string) (*models.ServiceToken, error) {
	panic("implement me")
}

func (p *PostgresServiceTokenRepository) CreateServiceToken(token *models.ServiceToken) error {
	panic("implement me")
}

func (p *PostgresServiceTokenRepository) RevokeServiceToken(id string) error {
	panic("implement me")
}
