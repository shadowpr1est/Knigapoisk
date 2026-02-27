package entity

import "time"

type BookStatus string

const (
	BookStatusActive BookStatus = "active"
	BookStatusHidden BookStatus = "hidden"
)

type Book struct {
	ID          int64
	Title       string
	Description string
	Language    string
	PublishedAt time.Time
	CoverUrl    string
	FileID      int64
	Status      BookStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Authors []Author
	Genres  []Genre
}
