package models

import "time"

type ServiceToken struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	AccessKey string    `json:"access_key" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type IServiceTokenRepository interface {
	GetServiceTokenById(id string) (*ServiceToken, error)
	GetServiceTokenByName(name string) (*ServiceToken, error)
	CreateServiceToken(token *ServiceToken) error
	UpdateServiceToken(token *ServiceToken) error
	RevokeServiceToken(id string) error
}
