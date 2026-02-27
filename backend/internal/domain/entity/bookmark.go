package entity

import "time"

type Bookmark struct {
	ID        int64
	UserID    int64
	BookID    int64
	Page      int
	Note      string
	CreatedAt time.Time
}
