package models

type TokenType int8

const (
	AdminToken TokenType = 1 << 0
	UserToken  TokenType = 1 << 1
)

type ServiceToken struct {
	ApplicationModel
	Name      string    `json:"name" binding:"required" gorm:"uniqueIndex"`
	AccessKey string    `json:"access_key" binding:"required"`
	TokenType TokenType `json:"token_type"`
}

func (st ServiceToken) GetID() string {
	return st.ID
}
