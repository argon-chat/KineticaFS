package models

type ServiceToken struct {
	ApplicationModel
	Name      string `json:"name" binding:"required" gorm:"uniqueIndex"`
	AccessKey string `json:"access_key" binding:"required"`
	TokenType int8   `json:"token_type"`
}

func (st ServiceToken) GetID() string {
	return st.ID
}
