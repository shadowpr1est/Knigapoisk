package entity

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}
