package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-reading-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-reading-service/internal/domain/repository"
)

type BookmarkRepo struct {
	db *sqlx.DB
}

func NewBookmarkRepo(db *sqlx.DB) repository.BookmarkRepository {
	return &BookmarkRepo{db: db}
}

func (r *BookmarkRepo) Create(ctx context.Context, b *entity.Bookmark) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	query := `
		INSERT INTO bookmarks (user_id, book_id, page, note, created_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`
	return r.db.QueryRowxContext(
		ctx,
		query,
		b.UserID,
		b.BookID,
		b.Page,
		b.Note,
		b.CreatedAt,
	).Scan(&b.ID)
}

func (r *BookmarkRepo) GetByID(ctx context.Context, id int64) (*entity.Bookmark, error) {
	var b entity.Bookmark
	query := `
		SELECT id, user_id, book_id, page, note, created_at
		FROM bookmarks
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

func (r *BookmarkRepo) ListByUserIDAndBookID(ctx context.Context, userID, bookID int64) ([]*entity.Bookmark, error) {
	var items []*entity.Bookmark
	query := `
		SELECT id, user_id, book_id, page, note, created_at
		FROM bookmarks
		WHERE user_id = $1 AND book_id = $2
		ORDER BY page
	`
	if err := r.db.SelectContext(ctx, &items, query, userID, bookID); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *BookmarkRepo) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.Bookmark, error) {
	var items []*entity.Bookmark
	query := `
		SELECT id, user_id, book_id, page, note, created_at
		FROM bookmarks
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	if err := r.db.SelectContext(ctx, &items, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *BookmarkRepo) Update(ctx context.Context, b *entity.Bookmark) error {
	query := `
		UPDATE bookmarks
		SET page = $1,
		    note = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, b.Page, b.Note, b.ID)
	return err
}

func (r *BookmarkRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM bookmarks WHERE id = $1`, id)
	return err
}

func (r *BookmarkRepo) DeleteByBookID(ctx context.Context, bookID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM bookmarks WHERE book_id = $1`, bookID)
	return err
}

