package entity

import "time"

type Progress struct {
	ID         int64
	UserID     int64
	BookID     int64
	FileID     int64
	Page       int
	Percentage float64
	UpdatedAt  time.Time
}
