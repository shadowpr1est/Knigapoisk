package review

import (
	"context"
	"errors"

	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-review-service/internal/domain/repository"
)

var (
	ErrReviewNotFound = errors.New("review not found")
	ErrForbidden      = errors.New("forbidden")
)

type UseCase interface {
	CreateReview(ctx context.Context, input CreateReviewInput) (CreateReviewOutput, error)
	GetReviews(ctx context.Context, bookID int64, limit, offset int) (GetReviewsOutput, error)
	UpdateReview(ctx context.Context, input UpdateReviewInput) (UpdateReviewOutput, error)
	DeleteReview(ctx context.Context, userID, reviewID int64) error
	GetRating(ctx context.Context, bookID int64) (GetRatingOutput, error)
}

type ReviewUseCase struct {
	reviewRepo repository.ReviewRepository
	ratingRepo repository.RatingRepository
}

func NewReviewUseCase(
	reviewRepo repository.ReviewRepository,
	ratingRepo repository.RatingRepository,
) UseCase {
	return &ReviewUseCase{
		reviewRepo: reviewRepo,
		ratingRepo: ratingRepo,
	}
}

type CreateReviewInput struct {
	UserID int64
	BookID int64
	Rating int16
	Text   string
}

type CreateReviewOutput struct {
	Review *entity.Review
}

type GetReviewsOutput struct {
	Reviews []*entity.Review
	Total   int64
}

type UpdateReviewInput struct {
	ID     int64
	UserID int64
	Rating int16
	Text   string
}

type UpdateReviewOutput struct {
	Review *entity.Review
}

type GetRatingOutput struct {
	Rating *entity.BookRating
}

func (u *ReviewUseCase) CreateReview(ctx context.Context, input CreateReviewInput) (CreateReviewOutput, error) {
	existing, err := u.reviewRepo.GetByUserIDAndBookID(ctx, input.UserID, input.BookID)
	if err != nil {
		return CreateReviewOutput{}, err
	}
	if existing != nil {
		return CreateReviewOutput{}, ErrForbidden
	}

	rv := &entity.Review{
		UserID: input.UserID,
		BookID: input.BookID,
		Rating: input.Rating,
		Text:   input.Text,
	}
	if err := u.reviewRepo.Create(ctx, rv); err != nil {
		return CreateReviewOutput{}, err
	}

	if err := u.updateRatingOnCreate(ctx, rv); err != nil {
		return CreateReviewOutput{}, err
	}

	return CreateReviewOutput{Review: rv}, nil
}

func (u *ReviewUseCase) GetReviews(ctx context.Context, bookID int64, limit, offset int) (GetReviewsOutput, error) {
	list, err := u.reviewRepo.ListByBookID(ctx, bookID, limit, offset)
	if err != nil {
		return GetReviewsOutput{}, err
	}
	var total int64
	for range list {
		total++
	}
	return GetReviewsOutput{
		Reviews: list,
		Total:   total,
	}, nil
}

func (u *ReviewUseCase) UpdateReview(ctx context.Context, input UpdateReviewInput) (UpdateReviewOutput, error) {
	rv, err := u.reviewRepo.GetByID(ctx, input.ID)
	if err != nil {
		return UpdateReviewOutput{}, err
	}
	if rv == nil {
		return UpdateReviewOutput{}, ErrReviewNotFound
	}
	if rv.UserID != input.UserID {
		return UpdateReviewOutput{}, ErrForbidden
	}

	oldRating := rv.Rating
	rv.Rating = input.Rating
	rv.Text = input.Text

	if err := u.reviewRepo.Update(ctx, rv); err != nil {
		return UpdateReviewOutput{}, err
	}
	if err := u.updateRatingOnUpdate(ctx, rv.BookID, oldRating, rv.Rating); err != nil {
		return UpdateReviewOutput{}, err
	}

	return UpdateReviewOutput{Review: rv}, nil
}

func (u *ReviewUseCase) DeleteReview(ctx context.Context, userID, reviewID int64) error {
	rv, err := u.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return err
	}
	if rv == nil {
		return nil
	}
	if rv.UserID != userID {
		return ErrForbidden
	}
	if err := u.reviewRepo.Delete(ctx, reviewID); err != nil {
		return err
	}
	return u.updateRatingOnDelete(ctx, rv.BookID, rv.Rating)
}

func (u *ReviewUseCase) GetRating(ctx context.Context, bookID int64) (GetRatingOutput, error) {
	r, err := u.ratingRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return GetRatingOutput{}, err
	}
	return GetRatingOutput{Rating: r}, nil
}

func (u *ReviewUseCase) updateRatingOnCreate(ctx context.Context, rv *entity.Review) error {
	r, err := u.ratingRepo.GetByBookID(ctx, rv.BookID)
	if err != nil {
		return err
	}
	if r == nil {
		r = &entity.BookRating{
			BookID:      rv.BookID,
			RatingSum:   int64(rv.Rating),
			RatingCount: 1,
			Average:     float64(rv.Rating),
		}
		return u.ratingRepo.Create(ctx, r)
	}
	r.RatingSum += int64(rv.Rating)
	r.RatingCount++
	r.Average = float64(r.RatingSum) / float64(r.RatingCount)
	return u.ratingRepo.Update(ctx, r)
}

func (u *ReviewUseCase) updateRatingOnUpdate(ctx context.Context, bookID int64, oldRating, newRating int16) error {
	if oldRating == newRating {
		return nil
	}
	r, err := u.ratingRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return err
	}
	if r == nil {
		return nil
	}
	r.RatingSum += int64(newRating - oldRating)
	r.Average = float64(r.RatingSum) / float64(r.RatingCount)
	return u.ratingRepo.Update(ctx, r)
}

func (u *ReviewUseCase) updateRatingOnDelete(ctx context.Context, bookID int64, rating int16) error {
	r, err := u.ratingRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return err
	}
	if r == nil {
		return nil
	}
	r.RatingSum -= int64(rating)
	r.RatingCount--
	if r.RatingCount <= 0 {
		r.RatingSum = 0
		r.RatingCount = 0
		r.Average = 0
	} else {
		r.Average = float64(r.RatingSum) / float64(r.RatingCount)
	}
	return u.ratingRepo.Update(ctx, r)
}

