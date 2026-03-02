package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	reviewpb "github.com/shadowpr1est/knigapoisk-review-service/api/proto"
)

type ReviewClient interface {
	CreateReview(ctx context.Context, req *reviewpb.CreateReviewRequest) (*reviewpb.CreateReviewResponse, error)
	GetReviews(ctx context.Context, req *reviewpb.GetReviewsRequest) (*reviewpb.GetReviewsResponse, error)
	UpdateReview(ctx context.Context, req *reviewpb.UpdateReviewRequest) (*reviewpb.UpdateReviewResponse, error)
	DeleteReview(ctx context.Context, req *reviewpb.DeleteReviewRequest) (*reviewpb.DeleteReviewResponse, error)
	GetRating(ctx context.Context, req *reviewpb.GetRatingRequest) (*reviewpb.GetRatingResponse, error)
}

type reviewClient struct {
	cc reviewpb.ReviewServiceClient
}

func NewReviewClient(addr string) (ReviewClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &reviewClient{cc: reviewpb.NewReviewServiceClient(conn)}, nil
}

func (c *reviewClient) CreateReview(ctx context.Context, req *reviewpb.CreateReviewRequest) (*reviewpb.CreateReviewResponse, error) {
	return c.cc.CreateReview(ctx, req)
}

func (c *reviewClient) GetReviews(ctx context.Context, req *reviewpb.GetReviewsRequest) (*reviewpb.GetReviewsResponse, error) {
	return c.cc.GetReviews(ctx, req)
}

func (c *reviewClient) UpdateReview(ctx context.Context, req *reviewpb.UpdateReviewRequest) (*reviewpb.UpdateReviewResponse, error) {
	return c.cc.UpdateReview(ctx, req)
}

func (c *reviewClient) DeleteReview(ctx context.Context, req *reviewpb.DeleteReviewRequest) (*reviewpb.DeleteReviewResponse, error) {
	return c.cc.DeleteReview(ctx, req)
}

func (c *reviewClient) GetRating(ctx context.Context, req *reviewpb.GetRatingRequest) (*reviewpb.GetRatingResponse, error) {
	return c.cc.GetRating(ctx, req)
}

