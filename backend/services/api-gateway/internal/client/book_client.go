package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	bookpb "github.com/shadowpr1est/knigapoisk-book-service/api/proto"
)

type BookClient interface {
	GetBook(ctx context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error)
	ListBooks(ctx context.Context, req *bookpb.ListBooksRequest) (*bookpb.ListBooksResponse, error)
	CreateBook(ctx context.Context, req *bookpb.CreateBookRequest) (*bookpb.CreateBookResponse, error)
	UpdateBook(ctx context.Context, req *bookpb.UpdateBookRequest) (*bookpb.UpdateBookResponse, error)
	DeleteBook(ctx context.Context, req *bookpb.DeleteBookRequest) (*bookpb.DeleteBookResponse, error)
	GetBooksByAuthor(ctx context.Context, req *bookpb.GetBooksByAuthorRequest) (*bookpb.GetBooksByAuthorResponse, error)
	GetBooksByGenre(ctx context.Context, req *bookpb.GetBooksByGenreRequest) (*bookpb.GetBooksByGenreResponse, error)
}

type bookClient struct {
	cc bookpb.BookServiceClient
}

func NewBookClient(addr string) (BookClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &bookClient{cc: bookpb.NewBookServiceClient(conn)}, nil
}

func (c *bookClient) GetBook(ctx context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
	return c.cc.GetBook(ctx, req)
}

func (c *bookClient) ListBooks(ctx context.Context, req *bookpb.ListBooksRequest) (*bookpb.ListBooksResponse, error) {
	return c.cc.ListBooks(ctx, req)
}

func (c *bookClient) CreateBook(ctx context.Context, req *bookpb.CreateBookRequest) (*bookpb.CreateBookResponse, error) {
	return c.cc.CreateBook(ctx, req)
}

func (c *bookClient) UpdateBook(ctx context.Context, req *bookpb.UpdateBookRequest) (*bookpb.UpdateBookResponse, error) {
	return c.cc.UpdateBook(ctx, req)
}

func (c *bookClient) DeleteBook(ctx context.Context, req *bookpb.DeleteBookRequest) (*bookpb.DeleteBookResponse, error) {
	return c.cc.DeleteBook(ctx, req)
}

func (c *bookClient) GetBooksByAuthor(ctx context.Context, req *bookpb.GetBooksByAuthorRequest) (*bookpb.GetBooksByAuthorResponse, error) {
	return c.cc.GetBooksByAuthor(ctx, req)
}

func (c *bookClient) GetBooksByGenre(ctx context.Context, req *bookpb.GetBooksByGenreRequest) (*bookpb.GetBooksByGenreResponse, error) {
	return c.cc.GetBooksByGenre(ctx, req)
}

