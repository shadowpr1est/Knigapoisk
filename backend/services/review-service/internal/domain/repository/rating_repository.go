package repository

import (
	"context"

	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/entity"
)

type RatingRepository interface {
	Create(ctx context.Context, rating *entity.BookRating) error
	GetByBookID(ctx context.Context, bookID int64) (*entity.BookRating, error)
	Update(ctx context.Context, rating *entity.BookRating) error
}

