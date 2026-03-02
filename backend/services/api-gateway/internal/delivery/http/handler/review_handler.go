package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/client"
	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/delivery/http/middleware"
	reviewpb "github.com/shadowpr1est/knigapoisk-review-service/api/proto"
)

type ReviewHandler struct {
	reviewClient client.ReviewClient
	validate     *validator.Validate
}

func NewReviewHandler(reviewClient client.ReviewClient) *ReviewHandler {
	return &ReviewHandler{
		reviewClient: reviewClient,
		validate:     validator.New(),
	}
}

type createReviewRequest struct {
	BookID int64  `json:"book_id" validate:"required"`
	Rating int32  `json:"rating" validate:"required,min=1,max=5"`
	Text   string `json:"text"`
}

type updateReviewRequest struct {
	Rating int32  `json:"rating" validate:"required,min=1,max=5"`
	Text   string `json:"text"`
}

func (h *ReviewHandler) GetReviews(c *gin.Context) {
	bookID, err := strconv.ParseInt(c.Param("bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bookID"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.reviewClient.GetReviews(c.Request.Context(), &reviewpb.GetReviewsRequest{
		BookId: bookID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "review service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": resp.Reviews,
		"total": resp.Total,
	})
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.reviewClient.CreateReview(c.Request.Context(), &reviewpb.CreateReviewRequest{
		UserId: userID,
		BookId: req.BookID,
		Rating: req.Rating,
		Text:   req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "review service error"})
		return
	}
	c.JSON(http.StatusOK, resp.Review)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	var req updateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	resp, err := h.reviewClient.UpdateReview(c.Request.Context(), &reviewpb.UpdateReviewRequest{
		Id:     id,
		UserId: userID,
		Rating: req.Rating,
		Text:   req.Text,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "review service error"})
		return
	}
	c.JSON(http.StatusOK, resp.Review)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	_, err = h.reviewClient.DeleteReview(c.Request.Context(), &reviewpb.DeleteReviewRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "review service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ReviewHandler) GetRating(c *gin.Context) {
	bookID, err := strconv.ParseInt(c.Param("bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bookID"})
		return
	}
	resp, err := h.reviewClient.GetRating(c.Request.Context(), &reviewpb.GetRatingRequest{
		BookId: bookID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "review service error"})
		return
	}
	c.JSON(http.StatusOK, resp.Rating)
}

