package scylla

import (
	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gocql/gocql"
)

type ScyllaServiceTokenRepository struct {
	session *gocql.Session
}

func (s ScyllaServiceTokenRepository) GetServiceTokenById(id string) (*models.ServiceToken, error) {
	panic("implement me")
}

func (s ScyllaServiceTokenRepository) GetServiceTokenByName(name string) (*models.ServiceToken, error) {
	panic("implement me")
}

func (s ScyllaServiceTokenRepository) CreateServiceToken(token *models.ServiceToken) error {
	panic("implement me")
}

func (s ScyllaServiceTokenRepository) UpdateServiceToken(token *models.ServiceToken) error {
	panic("implement me")
}

func (s ScyllaServiceTokenRepository) RevokeServiceToken(id string) error {
	panic("implement me")
}
