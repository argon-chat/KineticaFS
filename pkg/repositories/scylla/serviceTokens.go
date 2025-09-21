package scylla

import (
	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
)

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
	panic("implement me")
}

func (s *ScyllaServiceTokenRepository) CreateServiceToken(token *models.ServiceToken) error {
	panic("implement me")
}

func (s *ScyllaServiceTokenRepository) RevokeServiceToken(id string) error {
	panic("implement me")
}
