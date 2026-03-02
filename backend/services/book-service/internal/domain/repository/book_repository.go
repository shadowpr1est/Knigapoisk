package repository

import (
	"context"

	"github.com/shadowpr1est/knigapoisk-book-service/internal/domain/entity"
)

type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) error
	GetByID(ctx context.Context, id int64) (*entity.Book, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Book, error)
	Update(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int64, error)

	GetByAuthorID(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Book, int64, error)
	GetByGenreID(ctx context.Context, genreID int64, limit, offset int) ([]*entity.Book, int64, error)
}

