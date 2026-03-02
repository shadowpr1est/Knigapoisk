package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	reviewpb "github.com/shadowpr1est/knigapoisk-review-service/api/proto"
	"github.com/shadowpr1est/knigapoisk-review-service/internal/usecase/review"
)

type ReviewServer struct {
	reviewpb.UnimplementedReviewServiceServer
	useCase review.UseCase
	logger  *zap.Logger
}

func NewReviewServer(uc review.UseCase, logger *zap.Logger) *ReviewServer {
	return &ReviewServer{
		useCase: uc,
		logger:  logger,
	}
}

func (s *ReviewServer) CreateReview(ctx context.Context, req *reviewpb.CreateReviewRequest) (*reviewpb.CreateReviewResponse, error) {
	out, err := s.useCase.CreateReview(ctx, review.CreateReviewInput{
		UserID: req.UserId,
		BookID: req.BookId,
		Rating: int16(req.Rating),
		Text:   req.Text,
	})
	if err != nil {
		switch err {
		case review.ErrForbidden:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			s.logger.Error("CreateReview error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	r := out.Review
	return &reviewpb.CreateReviewResponse{
		Review: &reviewpb.Review{
			Id:        r.ID,
			UserId:    r.UserID,
			BookId:    r.BookID,
			Rating:    int32(r.Rating),
			Text:      r.Text,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
			UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *ReviewServer) GetReviews(ctx context.Context, req *reviewpb.GetReviewsRequest) (*reviewpb.GetReviewsResponse, error) {
	out, err := s.useCase.GetReviews(ctx, req.BookId, int(req.Limit), int(req.Offset))
	if err != nil {
		s.logger.Error("GetReviews error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	resp := &reviewpb.GetReviewsResponse{
		Reviews: make([]*reviewpb.Review, 0, len(out.Reviews)),
		Total:   out.Total,
	}
	for _, r := range out.Reviews {
		resp.Reviews = append(resp.Reviews, &reviewpb.Review{
			Id:        r.ID,
			UserId:    r.UserID,
			BookId:    r.BookID,
			Rating:    int32(r.Rating),
			Text:      r.Text,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
			UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}

func (s *ReviewServer) UpdateReview(ctx context.Context, req *reviewpb.UpdateReviewRequest) (*reviewpb.UpdateReviewResponse, error) {
	out, err := s.useCase.UpdateReview(ctx, review.UpdateReviewInput{
		ID:     req.Id,
		UserID: req.UserId,
		Rating: int16(req.Rating),
		Text:   req.Text,
	})
	if err != nil {
		switch err {
		case review.ErrReviewNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case review.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			s.logger.Error("UpdateReview error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	r := out.Review
	return &reviewpb.UpdateReviewResponse{
		Review: &reviewpb.Review{
			Id:        r.ID,
			UserId:    r.UserID,
			BookId:    r.BookID,
			Rating:    int32(r.Rating),
			Text:      r.Text,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
			UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *ReviewServer) DeleteReview(ctx context.Context, req *reviewpb.DeleteReviewRequest) (*reviewpb.DeleteReviewResponse, error) {
	if err := s.useCase.DeleteReview(ctx, req.UserId, req.Id); err != nil {
		switch err {
		case review.ErrForbidden:
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			s.logger.Error("DeleteReview error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &reviewpb.DeleteReviewResponse{Success: true}, nil
}

func (s *ReviewServer) GetRating(ctx context.Context, req *reviewpb.GetRatingRequest) (*reviewpb.GetRatingResponse, error) {
	out, err := s.useCase.GetRating(ctx, req.BookId)
	if err != nil {
		s.logger.Error("GetRating error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if out.Rating == nil {
		return &reviewpb.GetRatingResponse{}, nil
	}
	r := out.Rating
	return &reviewpb.GetRatingResponse{
		Rating: &reviewpb.BookRating{
			BookId:      r.BookID,
			RatingSum:   r.RatingSum,
			RatingCount: r.RatingCount,
			Average:     r.Average,
		},
	}, nil
}

func (s *ReviewServer) Health(ctx context.Context, _ *reviewpb.HealthRequest) (*reviewpb.HealthResponse, error) {
	return &reviewpb.HealthResponse{Status: "ok"}, nil
}

