package repository

import (
	"context"

	"github.com/shadowpr1est/knigapoisk-book-service/internal/domain/entity"
)

type GenreRepository interface {
	Create(ctx context.Context, genre *entity.Genre) error
	GetByID(ctx context.Context, id int64) (*entity.Genre, error)
	List(ctx context.Context) ([]*entity.Genre, error)
	GetBooksByGenreID(ctx context.Context, id int64, limit, offset int) ([]*entity.Book, int64, error)
}

