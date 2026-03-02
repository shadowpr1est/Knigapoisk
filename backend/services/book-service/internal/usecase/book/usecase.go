package book

import (
	"context"
	"errors"
	"time"

	"github.com/shadowpr1est/knigapoisk-book-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-book-service/internal/domain/repository"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type UseCase interface {
	CreateBook(ctx context.Context, input CreateBookInput) (CreateBookOutput, error)
	GetBook(ctx context.Context, id int64) (GetBookOutput, error)
	ListBooks(ctx context.Context, limit, offset int) (ListBooksOutput, error)
	UpdateBook(ctx context.Context, input UpdateBookInput) (UpdateBookOutput, error)
	DeleteBook(ctx context.Context, id int64) error
	GetBooksByAuthor(ctx context.Context, authorID int64, limit, offset int) (ListBooksOutput, error)
	GetBooksByGenre(ctx context.Context, genreID int64, limit, offset int) (ListBooksOutput, error)
}

type BookUseCase struct {
	bookRepo   repository.BookRepository
	authorRepo repository.AuthorRepository
	genreRepo  repository.GenreRepository
}

func NewBookUseCase(
	bookRepo repository.BookRepository,
	authorRepo repository.AuthorRepository,
	genreRepo repository.GenreRepository,
) UseCase {
	return &BookUseCase{
		bookRepo:   bookRepo,
		authorRepo: authorRepo,
		genreRepo:  genreRepo,
	}
}

type CreateBookInput struct {
	Title       string
	Description string
	Language    string
	PublishedAt time.Time
	CoverURL    string
	FileID      int64
	Status      entity.BookStatus
	AuthorIDs   []int64
	GenreIDs    []int64
}

type CreateBookOutput struct {
	ID int64
}

type GetBookOutput struct {
	Book *entity.Book
}

type ListBooksOutput struct {
	Books []*entity.Book
	Total int64
}

type UpdateBookInput struct {
	ID          int64
	Title       string
	Description string
	Language    string
	PublishedAt time.Time
	CoverURL    string
	FileID      int64
	Status      entity.BookStatus
	AuthorIDs   []int64
	GenreIDs    []int64
}

type UpdateBookOutput struct {
	Book *entity.Book
}

func (u *BookUseCase) CreateBook(ctx context.Context, input CreateBookInput) (CreateBookOutput, error) {
	now := time.Now()
	book := &entity.Book{
		Title:       input.Title,
		Description: input.Description,
		Language:    input.Language,
		PublishedAt: input.PublishedAt,
		CoverURL:    input.CoverURL,
		FileID:      input.FileID,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := u.bookRepo.Create(ctx, book); err != nil {
		return CreateBookOutput{}, err
	}
	return CreateBookOutput{ID: book.ID}, nil
}

func (u *BookUseCase) GetBook(ctx context.Context, id int64) (GetBookOutput, error) {
	book, err := u.bookRepo.GetByID(ctx, id)
	if err != nil {
		return GetBookOutput{}, err
	}
	if book == nil {
		return GetBookOutput{}, ErrBookNotFound
	}
	return GetBookOutput{Book: book}, nil
}

func (u *BookUseCase) ListBooks(ctx context.Context, limit, offset int) (ListBooksOutput, error) {
	books, err := u.bookRepo.List(ctx, limit, offset)
	if err != nil {
		return ListBooksOutput{}, err
	}
	total, err := u.bookRepo.Count(ctx)
	if err != nil {
		return ListBooksOutput{}, err
	}
	return ListBooksOutput{
		Books: books,
		Total: total,
	}, nil
}

func (u *BookUseCase) UpdateBook(ctx context.Context, input UpdateBookInput) (UpdateBookOutput, error) {
	book, err := u.bookRepo.GetByID(ctx, input.ID)
	if err != nil {
		return UpdateBookOutput{}, err
	}
	if book == nil {
		return UpdateBookOutput{}, ErrBookNotFound
	}

	book.Title = input.Title
	book.Description = input.Description
	book.Language = input.Language
	book.PublishedAt = input.PublishedAt
	book.CoverURL = input.CoverURL
	book.FileID = input.FileID
	book.Status = input.Status
	book.UpdatedAt = time.Now()

	if err := u.bookRepo.Update(ctx, book); err != nil {
		return UpdateBookOutput{}, err
	}
	return UpdateBookOutput{Book: book}, nil
}

func (u *BookUseCase) DeleteBook(ctx context.Context, id int64) error {
	return u.bookRepo.Delete(ctx, id)
}

func (u *BookUseCase) GetBooksByAuthor(ctx context.Context, authorID int64, limit, offset int) (ListBooksOutput, error) {
	books, total, err := u.bookRepo.GetByAuthorID(ctx, authorID, limit, offset)
	if err != nil {
		return ListBooksOutput{}, err
	}
	return ListBooksOutput{
		Books: books,
		Total: total,
	}, nil
}

func (u *BookUseCase) GetBooksByGenre(ctx context.Context, genreID int64, limit, offset int) (ListBooksOutput, error) {
	books, total, err := u.bookRepo.GetByGenreID(ctx, genreID, limit, offset)
	if err != nil {
		return ListBooksOutput{}, err
	}
	return ListBooksOutput{
		Books: books,
		Total: total,
	}, nil
}

