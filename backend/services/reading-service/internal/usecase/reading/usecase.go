package reading

import (
	"context"
	"time"

	"github.com/shadowpr1est/knigapoisk-reading-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-reading-service/internal/domain/repository"
)

type UseCase interface {
	SaveProgress(ctx context.Context, input SaveProgressInput) error
	GetProgress(ctx context.Context, userID, bookID int64) (GetProgressOutput, error)
	AddBookmark(ctx context.Context, input AddBookmarkInput) (AddBookmarkOutput, error)
	GetBookmarks(ctx context.Context, userID, bookID int64) (GetBookmarksOutput, error)
	DeleteBookmark(ctx context.Context, userID, bookmarkID int64) error
}

type ReadingUseCase struct {
	progressRepo repository.ProgressRepository
	bookmarkRepo repository.BookmarkRepository
}

func NewReadingUseCase(
	progressRepo repository.ProgressRepository,
	bookmarkRepo repository.BookmarkRepository,
) UseCase {
	return &ReadingUseCase{
		progressRepo: progressRepo,
		bookmarkRepo: bookmarkRepo,
	}
}

type SaveProgressInput struct {
	UserID     int64
	BookID     int64
	FileID     int64
	Page       int
	Percentage float64
}

type GetProgressOutput struct {
	Progress *entity.ReadingProgress
}

type AddBookmarkInput struct {
	UserID int64
	BookID int64
	Page   int
	Note   string
}

type AddBookmarkOutput struct {
	Bookmark *entity.Bookmark
}

type GetBookmarksOutput struct {
	Bookmarks []*entity.Bookmark
}

func (u *ReadingUseCase) SaveProgress(ctx context.Context, input SaveProgressInput) error {
	existing, err := u.progressRepo.GetByUserIDAndBookID(ctx, input.UserID, input.BookID)
	if err != nil {
		return err
	}
	now := time.Now()
	if existing == nil {
		p := &entity.ReadingProgress{
			UserID:     input.UserID,
			BookID:     input.BookID,
			FileID:     input.FileID,
			Page:       input.Page,
			Percentage: input.Percentage,
			UpdatedAt:  now,
		}
		return u.progressRepo.Create(ctx, p)
	}
	existing.FileID = input.FileID
	existing.Page = input.Page
	existing.Percentage = input.Percentage
	existing.UpdatedAt = now
	return u.progressRepo.Update(ctx, existing)
}

func (u *ReadingUseCase) GetProgress(ctx context.Context, userID, bookID int64) (GetProgressOutput, error) {
	p, err := u.progressRepo.GetByUserIDAndBookID(ctx, userID, bookID)
	if err != nil {
		return GetProgressOutput{}, err
	}
	return GetProgressOutput{Progress: p}, nil
}

func (u *ReadingUseCase) AddBookmark(ctx context.Context, input AddBookmarkInput) (AddBookmarkOutput, error) {
	b := &entity.Bookmark{
		UserID: input.UserID,
		BookID: input.BookID,
		Page:   input.Page,
		Note:   input.Note,
	}
	if err := u.bookmarkRepo.Create(ctx, b); err != nil {
		return AddBookmarkOutput{}, err
	}
	return AddBookmarkOutput{Bookmark: b}, nil
}

func (u *ReadingUseCase) GetBookmarks(ctx context.Context, userID, bookID int64) (GetBookmarksOutput, error) {
	items, err := u.bookmarkRepo.ListByUserIDAndBookID(ctx, userID, bookID)
	if err != nil {
		return GetBookmarksOutput{}, err
	}
	return GetBookmarksOutput{Bookmarks: items}, nil
}

func (u *ReadingUseCase) DeleteBookmark(ctx context.Context, userID, bookmarkID int64) error {
	b, err := u.bookmarkRepo.GetByID(ctx, bookmarkID)
	if err != nil {
		return err
	}
	if b == nil {
		return nil
	}
	if b.UserID != userID {
		return nil
	}
	return u.bookmarkRepo.Delete(ctx, bookmarkID)
}

