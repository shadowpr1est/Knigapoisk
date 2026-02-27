package repository

import (
	"context"

	"github.com/shadowpr1est/Knigapoisk/internal/domain/entity"
)

type AuthorRepository interface {
	Create(ctx context.Context, author *entity.Author) error
	GetByID(ctx context.Context, id int64) (*entity.Author, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Author, error)
	Update(ctx context.Context, author *entity.Author) error
	Delete(ctx context.Context, id int64) error

	GetBooksByAuthorID(ctx context.Context, id int64) ([]*entity.Book, error)
}
