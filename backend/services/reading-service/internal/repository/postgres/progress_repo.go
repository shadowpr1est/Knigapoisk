package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-reading-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-reading-service/internal/domain/repository"
)

type ProgressRepo struct {
	db *sqlx.DB
}

func NewProgressRepo(db *sqlx.DB) repository.ProgressRepository {
	return &ProgressRepo{db: db}
}

func (r *ProgressRepo) Create(ctx context.Context, p *entity.ReadingProgress) error {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}
	query := `
		INSERT INTO reading_progress (user_id, book_id, file_id, page, percentage, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`
	return r.db.QueryRowxContext(
		ctx,
		query,
		p.UserID,
		p.BookID,
		p.FileID,
		p.Page,
		p.Percentage,
		p.UpdatedAt,
	).Scan(&p.ID)
}

func (r *ProgressRepo) GetByUserIDAndBookID(ctx context.Context, userID, bookID int64) (*entity.ReadingProgress, error) {
	var p entity.ReadingProgress
	query := `
		SELECT id, user_id, book_id, file_id, page, percentage, updated_at
		FROM reading_progress
		WHERE user_id = $1 AND book_id = $2
	`
	if err := r.db.GetContext(ctx, &p, query, userID, bookID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProgressRepo) Update(ctx context.Context, p *entity.ReadingProgress) error {
	p.UpdatedAt = time.Now()
	query := `
		UPDATE reading_progress
		SET file_id = $1,
		    page = $2,
		    percentage = $3,
		    updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		p.FileID,
		p.Page,
		p.Percentage,
		p.UpdatedAt,
		p.ID,
	)
	return err
}

func (r *ProgressRepo) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.ReadingProgress, error) {
	var items []*entity.ReadingProgress
	query := `
		SELECT id, user_id, book_id, file_id, page, percentage, updated_at
		FROM reading_progress
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &items, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProgressRepo) DeleteByBookID(ctx context.Context, bookID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reading_progress WHERE book_id = $1`, bookID)
	return err
}

