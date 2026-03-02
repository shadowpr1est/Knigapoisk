package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	readingpb "github.com/shadowpr1est/knigapoisk-reading-service/api/proto"
	"github.com/shadowpr1est/knigapoisk-reading-service/internal/usecase/reading"
)

type ReadingServer struct {
	readingpb.UnimplementedReadingServiceServer
	useCase reading.UseCase
	logger  *zap.Logger
}

func NewReadingServer(uc reading.UseCase, logger *zap.Logger) *ReadingServer {
	return &ReadingServer{
		useCase: uc,
		logger:  logger,
	}
}

func (s *ReadingServer) SaveProgress(ctx context.Context, req *readingpb.SaveProgressRequest) (*readingpb.SaveProgressResponse, error) {
	err := s.useCase.SaveProgress(ctx, reading.SaveProgressInput{
		UserID:     req.UserId,
		BookID:     req.BookId,
		FileID:     req.FileId,
		Page:       int(req.Page),
		Percentage: req.Percentage,
	})
	if err != nil {
		s.logger.Error("SaveProgress error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &readingpb.SaveProgressResponse{Success: true}, nil
}

func (s *ReadingServer) GetProgress(ctx context.Context, req *readingpb.GetProgressRequest) (*readingpb.GetProgressResponse, error) {
	out, err := s.useCase.GetProgress(ctx, req.UserId, req.BookId)
	if err != nil {
		s.logger.Error("GetProgress error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if out.Progress == nil {
		return &readingpb.GetProgressResponse{}, nil
	}
	p := out.Progress
	return &readingpb.GetProgressResponse{
		Progress: &readingpb.ReadingProgress{
			Id:         p.ID,
			UserId:     p.UserID,
			BookId:     p.BookID,
			FileId:     p.FileID,
			Page:       int32(p.Page),
			Percentage: p.Percentage,
			UpdatedAt:  p.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *ReadingServer) AddBookmark(ctx context.Context, req *readingpb.AddBookmarkRequest) (*readingpb.AddBookmarkResponse, error) {
	out, err := s.useCase.AddBookmark(ctx, reading.AddBookmarkInput{
		UserID: req.UserId,
		BookID: req.BookId,
		Page:   int(req.Page),
		Note:   req.Note,
	})
	if err != nil {
		s.logger.Error("AddBookmark error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	b := out.Bookmark
	return &readingpb.AddBookmarkResponse{
		Bookmark: &readingpb.Bookmark{
			Id:        b.ID,
			UserId:    b.UserID,
			BookId:    b.BookID,
			Page:      int32(b.Page),
			Note:      b.Note,
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *ReadingServer) GetBookmarks(ctx context.Context, req *readingpb.GetBookmarksRequest) (*readingpb.GetBookmarksResponse, error) {
	out, err := s.useCase.GetBookmarks(ctx, req.UserId, req.BookId)
	if err != nil {
		s.logger.Error("GetBookmarks error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	resp := &readingpb.GetBookmarksResponse{
		Bookmarks: make([]*readingpb.Bookmark, 0, len(out.Bookmarks)),
	}
	for _, b := range out.Bookmarks {
		resp.Bookmarks = append(resp.Bookmarks, &readingpb.Bookmark{
			Id:        b.ID,
			UserId:    b.UserID,
			BookId:    b.BookID,
			Page:      int32(b.Page),
			Note:      b.Note,
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}

func (s *ReadingServer) DeleteBookmark(ctx context.Context, req *readingpb.DeleteBookmarkRequest) (*readingpb.DeleteBookmarkResponse, error) {
	if err := s.useCase.DeleteBookmark(ctx, req.UserId, req.BookmarkId); err != nil {
		s.logger.Error("DeleteBookmark error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &readingpb.DeleteBookmarkResponse{Success: true}, nil
}

func (s *ReadingServer) Health(ctx context.Context, _ *readingpb.HealthRequest) (*readingpb.HealthResponse, error) {
	return &readingpb.HealthResponse{Status: "ok"}, nil
}

