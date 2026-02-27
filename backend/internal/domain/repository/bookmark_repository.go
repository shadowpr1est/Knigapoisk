package repository

import (
	"context"

	"github.com/shadowpr1est/Knigapoisk/internal/domain/entity"
)

type BookmarkRepository interface {
	Create(ctx context.Context, bookmark *entity.Bookmark) error
	GetByID(ctx context.Context, id int64) (*entity.Bookmark, error)
	ListByUserIDAndBookID(ctx context.Context, userID, bookID int64) ([]*entity.Bookmark, error)
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.Bookmark, error)
	Update(ctx context.Context, bookmark *entity.Bookmark) error
	Delete(ctx context.Context, id int64) error
	DeleteByBookID(ctx context.Context, bookID int64) error
}
