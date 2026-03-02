package entity

import "time"

type BookStatus string

const (
	BookStatusActive BookStatus = "active"
	BookStatusHidden BookStatus = "hidden"
)

type Book struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Language    string     `db:"language"`
	PublishedAt time.Time  `db:"published_at"`
	CoverURL    string     `db:"cover_url"`
	FileID      int64      `db:"file_id"`
	Status      BookStatus `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`

	Authors []Author
	Genres  []Genre
}

