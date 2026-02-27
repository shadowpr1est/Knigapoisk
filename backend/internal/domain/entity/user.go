package entity

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
