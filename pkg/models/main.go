package models

import "time"

type ApplicationRecord interface {
	GetID() string
}

type ApplicationModel struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
