package entity

import "time"

type Review struct {
	ID        int64
	UserID    int64
	BookID    int64
	Rating    int16
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
