package models

import "time"

type User struct {
	ID          string `json:"id" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	OtpSecret   string
	OtpLastUsed time.Time
	MfaEnabled  bool `json:"mfa_enabled" binding:"required"`
}

type UserRepository interface {
	GetUserById(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id string) error
}
