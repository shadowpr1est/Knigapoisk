package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-book-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-book-service/internal/domain/repository"
)

type BookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) repository.BookRepository {
	return &BookRepo{db: db}
}

func (r *BookRepo) Create(ctx context.Context, book *entity.Book) error {
	query := `
		INSERT INTO books (title, description, language, published_at, cover_url, file_id, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(
		ctx,
		query,
		book.Title,
		book.Description,
		book.Language,
		book.PublishedAt,
		book.CoverURL,
		book.FileID,
		book.Status,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
}

func (r *BookRepo) GetByID(ctx context.Context, id int64) (*entity.Book, error) {
	var b entity.Book
	query := `
		SELECT id, title, description, language, published_at, cover_url, file_id, status, created_at, updated_at
		FROM books
		WHERE id = $1
	`
	if err := r.db.GetContext(ctx, &b, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (r *BookRepo) List(ctx context.Context, limit, offset int) ([]*entity.Book, error) {
	var books []*entity.Book
	query := `
		SELECT id, title, description, language, published_at, cover_url, file_id, status, created_at, updated_at
		FROM books
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	if err := r.db.SelectContext(ctx, &books, query, limit, offset); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepo) Update(ctx context.Context, book *entity.Book) error {
	query := `
		UPDATE books
		SET title = $1,
		    description = $2,
		    language = $3,
		    published_at = $4,
		    cover_url = $5,
		    file_id = $6,
		    status = $7,
		    updated_at = now()
		WHERE id = $8
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		book.Title,
		book.Description,
		book.Language,
		book.PublishedAt,
		book.CoverURL,
		book.FileID,
		book.Status,
		book.ID,
	)
	return err
}

func (r *BookRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM books WHERE id = $1`, id)
	return err
}

func (r *BookRepo) Count(ctx context.Context) (int64, error) {
	var cnt int64
	if err := r.db.GetContext(ctx, &cnt, `SELECT COUNT(*) FROM books`); err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *BookRepo) GetByAuthorID(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Book, int64, error) {
	var books []*entity.Book
	query := `
		SELECT b.id, b.title, b.description, b.language, b.published_at, b.cover_url,
		       b.file_id, b.status, b.created_at, b.updated_at
		FROM books b
		JOIN book_authors ba ON ba.book_id = b.id
		WHERE ba.author_id = $1
		ORDER BY b.id
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &books, query, authorID, limit, offset); err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.db.GetContext(ctx, &total, `
		SELECT COUNT(*)
		FROM books b
		JOIN book_authors ba ON ba.book_id = b.id
		WHERE ba.author_id = $1
	`, authorID); err != nil {
		return nil, 0, err
	}
	return books, total, nil
}

func (r *BookRepo) GetByGenreID(ctx context.Context, genreID int64, limit, offset int) ([]*entity.Book, int64, error) {
	var books []*entity.Book
	query := `
		SELECT b.id, b.title, b.description, b.language, b.published_at, b.cover_url,
		       b.file_id, b.status, b.created_at, b.updated_at
		FROM books b
		JOIN book_genres bg ON bg.book_id = b.id
		WHERE bg.genre_id = $1
		ORDER BY b.id
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &books, query, genreID, limit, offset); err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.db.GetContext(ctx, &total, `
		SELECT COUNT(*)
		FROM books b
		JOIN book_genres bg ON bg.book_id = b.id
		WHERE bg.genre_id = $1
	`, genreID); err != nil {
		return nil, 0, err
	}
	return books, total, nil
}

