package repository

import (
	"context"

	"github.com/shadowpr1est/Knigapoisk/internal/domain/entity"
)

type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) error
	GetByID(ctx context.Context, id int64) (*entity.Book, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Book, error)
	Update(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int64, error)
}
