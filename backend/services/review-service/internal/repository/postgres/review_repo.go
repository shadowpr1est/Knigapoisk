package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/repository"
)

type ReviewRepo struct {
	db *sqlx.DB
}

func NewReviewRepo(db *sqlx.DB) repository.ReviewRepository {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(ctx context.Context, review *entity.Review) error {
	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	query := `
		INSERT INTO reviews (user_id, book_id, rating, text, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`
	return r.db.QueryRowxContext(
		ctx,
		query,
		review.UserID,
		review.BookID,
		review.Rating,
		review.Text,
		review.CreatedAt,
		review.UpdatedAt,
	).Scan(&review.ID)
}

func (r *ReviewRepo) GetByID(ctx context.Context, id int64) (*entity.Review, error) {
	var rv entity.Review
	query := `
		SELECT id, user_id, book_id, rating, text, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`
	if err := r.db.GetContext(ctx, &rv, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rv, nil
}

func (r *ReviewRepo) GetByUserIDAndBookID(ctx context.Context, userID, bookID int64) (*entity.Review, error) {
	var rv entity.Review
	query := `
		SELECT id, user_id, book_id, rating, text, created_at, updated_at
		FROM reviews
		WHERE user_id = $1 AND book_id = $2
	`
	if err := r.db.GetContext(ctx, &rv, query, userID, bookID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rv, nil
}

func (r *ReviewRepo) ListByBookID(ctx context.Context, bookID int64, limit, offset int) ([]*entity.Review, error) {
	var list []*entity.Review
	query := `
		SELECT id, user_id, book_id, rating, text, created_at, updated_at
		FROM reviews
		WHERE book_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &list, query, bookID, limit, offset); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ReviewRepo) Update(ctx context.Context, review *entity.Review) error {
	review.UpdatedAt = time.Now()
	query := `
		UPDATE reviews
		SET rating = $1,
		    text = $2,
		    updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		review.Rating,
		review.Text,
		review.UpdatedAt,
		review.ID,
	)
	return err
}

func (r *ReviewRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	return err
}

func (r *ReviewRepo) DeleteByBookID(ctx context.Context, bookID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE book_id = $1`, bookID)
	return err
}

