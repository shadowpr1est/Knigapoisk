package repository

import (
	"context"

	"github.com/shadowpr1est/Knigapoisk/internal/domain/entity"
)

type TokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllByUserID(ctx context.Context, userID int64) error
	DeleteExpired(ctx context.Context) error
}
