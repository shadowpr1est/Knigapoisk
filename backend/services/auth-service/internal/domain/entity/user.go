package entity

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

