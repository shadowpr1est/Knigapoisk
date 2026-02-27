package repository

import (
	"context"

	"github.com/shadowpr1est/Knigapoisk/internal/domain/entity"
)

type RatingRepository interface {
	Create(ctx context.Context, rating *entity.Rating) error
	Update(ctx context.Context, rating *entity.Rating) error
	GetByID(ctx context.Context, bookID int64) (*entity.Rating, error)
}
