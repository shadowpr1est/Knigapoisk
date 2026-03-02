package entity

import "time"

type Review struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	BookID    int64     `db:"book_id"`
	Rating    int16     `db:"rating"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

