package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/repository"
)

type RatingRepo struct {
	db *sqlx.DB
}

func NewRatingRepo(db *sqlx.DB) repository.RatingRepository {
	return &RatingRepo{db: db}
}

func (r *RatingRepo) Create(ctx context.Context, rating *entity.BookRating) error {
	query := `
		INSERT INTO book_ratings (book_id, rating_sum, rating_count, average)
		VALUES ($1,$2,$3,$4)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		rating.BookID,
		rating.RatingSum,
		rating.RatingCount,
		rating.Average,
	)
	return err
}

func (r *RatingRepo) GetByBookID(ctx context.Context, bookID int64) (*entity.BookRating, error) {
	var br entity.BookRating
	query := `
		SELECT book_id, rating_sum, rating_count, average
		FROM book_ratings
		WHERE book_id = $1
	`
	if err := r.db.GetContext(ctx, &br, query, bookID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &br, nil
}

func (r *RatingRepo) Update(ctx context.Context, rating *entity.BookRating) error {
	query := `
		UPDATE book_ratings
		SET rating_sum = $1,
		    rating_count = $2,
		    average = $3
		WHERE book_id = $4
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		rating.RatingSum,
		rating.RatingCount,
		rating.Average,
		rating.BookID,
	)
	return err
}

