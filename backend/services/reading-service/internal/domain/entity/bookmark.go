package entity

import "time"

type Bookmark struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	BookID    int64     `db:"book_id"`
	Page      int       `db:"page"`
	Note      string    `db:"note"`
	CreatedAt time.Time `db:"created_at"`
}

