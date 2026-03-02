package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-auth-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/domain/repository"
)

type TokenRepo struct {
	db *sqlx.DB
}

func NewTokenRepo(db *sqlx.DB) repository.TokenRepository {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) Create(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens(user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	).Scan(
		&token.ID,
		&token.CreatedAt,
	)
}

func (r *TokenRepo) GetByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	err := r.db.GetContext(ctx, &token, query, tokenHash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepo) Revoke(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE token_hash = $1
	`
	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *TokenRepo) RevokeAllByUserID(ctx context.Context, userID int64) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *TokenRepo) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE expires_at < $1
	`
	_, err := r.db.ExecContext(ctx, query, time.Now())
	return err
}

