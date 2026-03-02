package repository

import (
	"context"

	"github.com/shadowpr1est/knigapoisk-reading-service/internal/domain/entity"
)

type ProgressRepository interface {
	Create(ctx context.Context, progress *entity.ReadingProgress) error
	GetByUserIDAndBookID(ctx context.Context, userID, bookID int64) (*entity.ReadingProgress, error)
	Update(ctx context.Context, progress *entity.ReadingProgress) error
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.ReadingProgress, error)
	DeleteByBookID(ctx context.Context, bookID int64) error
}

