package repository

import (
	"context"

	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/entity"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *entity.Review) error
	GetByID(ctx context.Context, id int64) (*entity.Review, error)
	GetByUserIDAndBookID(ctx context.Context, userID, bookID int64) (*entity.Review, error)
	ListByBookID(ctx context.Context, bookID int64, limit, offset int) ([]*entity.Review, error)
	Update(ctx context.Context, review *entity.Review) error
	Delete(ctx context.Context, id int64) error
	DeleteByBookID(ctx context.Context, bookID int64) error
}

