package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bookpb "github.com/shadowpr1est/knigapoisk-book-service/api/proto"
	"github.com/shadowpr1est/knigapoisk-book-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-book-service/internal/usecase/book"
)

type BookServer struct {
	bookpb.UnimplementedBookServiceServer
	useCase book.UseCase
	logger  *zap.Logger
}

func NewBookServer(uc book.UseCase, logger *zap.Logger) *BookServer {
	return &BookServer{
		useCase: uc,
		logger:  logger,
	}
}

func (s *BookServer) GetBook(ctx context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
	out, err := s.useCase.GetBook(ctx, req.Id)
	if err != nil {
		if err == book.ErrBookNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		s.logger.Error("GetBook error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &bookpb.GetBookResponse{
		Book: toProtoBook(out.Book),
	}, nil
}

func (s *BookServer) ListBooks(ctx context.Context, req *bookpb.ListBooksRequest) (*bookpb.ListBooksResponse, error) {
	out, err := s.useCase.ListBooks(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("ListBooks error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	resp := &bookpb.ListBooksResponse{
		Books: make([]*bookpb.Book, 0, len(out.Books)),
		Total: out.Total,
	}
	for _, b := range out.Books {
		resp.Books = append(resp.Books, toProtoBook(b))
	}
	return resp, nil
}

func (s *BookServer) CreateBook(ctx context.Context, req *bookpb.CreateBookRequest) (*bookpb.CreateBookResponse, error) {
	publishedAt, _ := time.Parse(time.RFC3339, req.PublishedAt)
	input := book.CreateBookInput{
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
		PublishedAt: publishedAt,
		CoverURL:    req.CoverUrl,
		FileID:      req.FileId,
		Status:      fromProtoStatus(req.Status),
		AuthorIDs:   req.AuthorIds,
		GenreIDs:    req.GenreIds,
	}
	out, err := s.useCase.CreateBook(ctx, input)
	if err != nil {
		s.logger.Error("CreateBook error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &bookpb.CreateBookResponse{Id: out.ID}, nil
}

func (s *BookServer) UpdateBook(ctx context.Context, req *bookpb.UpdateBookRequest) (*bookpb.UpdateBookResponse, error) {
	publishedAt, _ := time.Parse(time.RFC3339, req.PublishedAt)
	input := book.UpdateBookInput{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
		PublishedAt: publishedAt,
		CoverURL:    req.CoverUrl,
		FileID:      req.FileId,
		Status:      fromProtoStatus(req.Status),
		AuthorIDs:   req.AuthorIds,
		GenreIDs:    req.GenreIds,
	}
	out, err := s.useCase.UpdateBook(ctx, input)
	if err != nil {
		if err == book.ErrBookNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		s.logger.Error("UpdateBook error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &bookpb.UpdateBookResponse{
		Book: toProtoBook(out.Book),
	}, nil
}

func (s *BookServer) DeleteBook(ctx context.Context, req *bookpb.DeleteBookRequest) (*bookpb.DeleteBookResponse, error) {
	if err := s.useCase.DeleteBook(ctx, req.Id); err != nil {
		s.logger.Error("DeleteBook error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &bookpb.DeleteBookResponse{Success: true}, nil
}

func (s *BookServer) GetBooksByAuthor(ctx context.Context, req *bookpb.GetBooksByAuthorRequest) (*bookpb.GetBooksByAuthorResponse, error) {
	out, err := s.useCase.GetBooksByAuthor(ctx, req.AuthorId, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("GetBooksByAuthor error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	resp := &bookpb.GetBooksByAuthorResponse{
		Books: make([]*bookpb.Book, 0, len(out.Books)),
		Total: out.Total,
	}
	for _, b := range out.Books {
		resp.Books = append(resp.Books, toProtoBook(b))
	}
	return resp, nil
}

func (s *BookServer) GetBooksByGenre(ctx context.Context, req *bookpb.GetBooksByGenreRequest) (*bookpb.GetBooksByGenreResponse, error) {
	out, err := s.useCase.GetBooksByGenre(ctx, req.GenreId, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("GetBooksByGenre error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	resp := &bookpb.GetBooksByGenreResponse{
		Books: make([]*bookpb.Book, 0, len(out.Books)),
		Total: out.Total,
	}
	for _, b := range out.Books {
		resp.Books = append(resp.Books, toProtoBook(b))
	}
	return resp, nil
}

func (s *BookServer) Health(ctx context.Context, _ *bookpb.HealthRequest) (*bookpb.HealthResponse, error) {
	return &bookpb.HealthResponse{Status: "ok"}, nil
}

func toProtoBook(b *entity.Book) *bookpb.Book {
	if b == nil {
		return nil
	}
	return &bookpb.Book{
		Id:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Language:    b.Language,
		PublishedAt: b.PublishedAt.Format(time.RFC3339),
		CoverUrl:    b.CoverURL,
		FileId:      b.FileID,
		Status:      toProtoStatus(b.Status),
		CreatedAt:   b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   b.UpdatedAt.Format(time.RFC3339),
	}
}

func toProtoStatus(s entity.BookStatus) bookpb.BookStatus {
	switch s {
	case entity.BookStatusActive:
		return bookpb.BookStatus_BOOK_STATUS_ACTIVE
	case entity.BookStatusHidden:
		return bookpb.BookStatus_BOOK_STATUS_HIDDEN
	default:
		return bookpb.BookStatus_BOOK_STATUS_UNSPECIFIED
	}
}

func fromProtoStatus(s bookpb.BookStatus) entity.BookStatus {
	switch s {
	case bookpb.BookStatus_BOOK_STATUS_ACTIVE:
		return entity.BookStatusActive
	case bookpb.BookStatus_BOOK_STATUS_HIDDEN:
		return entity.BookStatusHidden
	default:
		return ""
	}
}

