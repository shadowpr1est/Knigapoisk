package entity

import "time"

type ReadingProgress struct {
	ID         int64     `db:"id"`
	UserID     int64     `db:"user_id"`
	BookID     int64     `db:"book_id"`
	FileID     int64     `db:"file_id"`
	Page       int       `db:"page"`
	Percentage float64   `db:"percentage"`
	UpdatedAt  time.Time `db:"updated_at"`
}

