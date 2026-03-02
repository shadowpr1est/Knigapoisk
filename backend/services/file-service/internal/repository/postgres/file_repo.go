package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-file-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-file-service/internal/domain/repository"
)

type FileRepo struct {
	db *sqlx.DB
}

func NewFileRepo(db *sqlx.DB) repository.FileRepository {
	return &FileRepo{db: db}
}

func (r *FileRepo) Create(ctx context.Context, file *entity.File) error {
	query := `
		INSERT INTO files (book_id, format, storage_key, size_bytes, uploaded_at, checksum)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`
	return r.db.QueryRowxContext(
		ctx,
		query,
		file.BookID,
		file.Format,
		file.StorageKey,
		file.SizeBytes,
		file.UploadedAt,
		file.Checksum,
	).Scan(&file.ID)
}

func (r *FileRepo) GetByID(ctx context.Context, id int64) (*entity.File, error) {
	var f entity.File
	query := `
		SELECT id, book_id, format, storage_key, size_bytes, uploaded_at, checksum
		FROM files
		WHERE id = $1
	`
	if err := r.db.GetContext(ctx, &f, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *FileRepo) GetByBookIDAndFormat(ctx context.Context, bookID int64, format entity.FileFormat) (*entity.File, error) {
	var f entity.File
	query := `
		SELECT id, book_id, format, storage_key, size_bytes, uploaded_at, checksum
		FROM files
		WHERE book_id = $1 AND format = $2
	`
	if err := r.db.GetContext(ctx, &f, query, bookID, format); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (r *FileRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM files WHERE id = $1`, id)
	return err
}

