package postgres

import (
	"database/sql"
	"log"

	"github.com/argon-chat/KineticaFS/pkg/models"
)

type PostgresServiceTokenRepository struct {
	session *sql.DB
}

func NewPostgresServiceTokenRepository(session *sql.DB) *PostgresServiceTokenRepository {
	return &PostgresServiceTokenRepository{session: session}
}

func (s *PostgresServiceTokenRepository) CreateIndices() {
	indexQueries := []string{}
	for _, indexQuery := range indexQueries {
		log.Printf("Executing index creation query: %s", indexQuery)
		if _, err := s.session.Exec(indexQuery); err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
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
