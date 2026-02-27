package repository

import (
	"context"

	"github.com/shadowpr1est/Knigapoisk/internal/domain/entity"
)

type ProgressRepository interface {
	Create(ctx context.Context, progress *entity.Progress) error
	GetByUserIDAndBookID(ctx context.Context, userID, bookID int64) (*entity.Progress, error)
	Update(ctx context.Context, progress *entity.Progress) error
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.Progress, error)
	DeleteByBookID(ctx context.Context, bookID int64) error
}
