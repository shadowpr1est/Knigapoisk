package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	readingpb "github.com/shadowpr1est/knigapoisk-reading-service/api/proto"
)

type ReadingClient interface {
	SaveProgress(ctx context.Context, req *readingpb.SaveProgressRequest) (*readingpb.SaveProgressResponse, error)
	GetProgress(ctx context.Context, req *readingpb.GetProgressRequest) (*readingpb.GetProgressResponse, error)
	AddBookmark(ctx context.Context, req *readingpb.AddBookmarkRequest) (*readingpb.AddBookmarkResponse, error)
	GetBookmarks(ctx context.Context, req *readingpb.GetBookmarksRequest) (*readingpb.GetBookmarksResponse, error)
	DeleteBookmark(ctx context.Context, req *readingpb.DeleteBookmarkRequest) (*readingpb.DeleteBookmarkResponse, error)
}

type readingClient struct {
	cc readingpb.ReadingServiceClient
}

func NewReadingClient(addr string) (ReadingClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &readingClient{cc: readingpb.NewReadingServiceClient(conn)}, nil
}

func (c *readingClient) SaveProgress(ctx context.Context, req *readingpb.SaveProgressRequest) (*readingpb.SaveProgressResponse, error) {
	return c.cc.SaveProgress(ctx, req)
}

func (c *readingClient) GetProgress(ctx context.Context, req *readingpb.GetProgressRequest) (*readingpb.GetProgressResponse, error) {
	return c.cc.GetProgress(ctx, req)
}

func (c *readingClient) AddBookmark(ctx context.Context, req *readingpb.AddBookmarkRequest) (*readingpb.AddBookmarkResponse, error) {
	return c.cc.AddBookmark(ctx, req)
}

func (c *readingClient) GetBookmarks(ctx context.Context, req *readingpb.GetBookmarksRequest) (*readingpb.GetBookmarksResponse, error) {
	return c.cc.GetBookmarks(ctx, req)
}

func (c *readingClient) DeleteBookmark(ctx context.Context, req *readingpb.DeleteBookmarkRequest) (*readingpb.DeleteBookmarkResponse, error) {
	return c.cc.DeleteBookmark(ctx, req)
}

