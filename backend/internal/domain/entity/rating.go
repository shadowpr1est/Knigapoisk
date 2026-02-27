package entity

type Rating struct {
	BookID      int64
	RatingSum   int64
	RatingCount int64
	Average     float64
}
