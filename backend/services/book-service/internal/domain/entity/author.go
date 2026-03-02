package entity

import "time"

type Author struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Biography string    `db:"biography"`
	BornYear  int       `db:"born_year"`
	CreatedAt time.Time `db:"created_at"`
}

