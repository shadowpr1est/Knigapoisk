package entity

type BookRating struct {
	BookID      int64   `db:"book_id"`
	RatingSum   int64   `db:"rating_sum"`
	RatingCount int64   `db:"rating_count"`
	Average     float64 `db:"average"`
}

