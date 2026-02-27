package repository

import (
	"context"

	"github.com/shadowpr1est/Knigapoisk/internal/domain/entity"
)

type GenreRepository interface {
	Create(ctx context.Context, genre *entity.Genre) error
	GetByID(ctx context.Context, id int64) (*entity.Genre, error)
	List(ctx context.Context) ([]*entity.Genre, error)
	GetBooksByGenreID(ctx context.Context, id int64) ([]*entity.Genre, error)
}
