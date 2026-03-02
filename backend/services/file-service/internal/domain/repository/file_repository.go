package repository

import (
	"context"

	"github.com/shadowpr1est/knigapoisk-file-service/internal/domain/entity"
)

type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
	GetByID(ctx context.Context, id int64) (*entity.File, error)
	GetByBookIDAndFormat(ctx context.Context, bookID int64, format entity.FileFormat) (*entity.File, error)
	Delete(ctx context.Context, id int64) error
}

