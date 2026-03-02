package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/client"
	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/delivery/http/middleware"
	readingpb "github.com/shadowpr1est/knigapoisk-reading-service/api/proto"
)

type ReadingHandler struct {
	readingClient client.ReadingClient
	validate      *validator.Validate
}

func NewReadingHandler(readingClient client.ReadingClient) *ReadingHandler {
	return &ReadingHandler{
		readingClient: readingClient,
		validate:      validator.New(),
	}
}

type saveProgressRequest struct {
	BookID     int64   `json:"book_id" validate:"required"`
	FileID     int64   `json:"file_id" validate:"required"`
	Page       int32   `json:"page" validate:"required"`
	Percentage float64 `json:"percentage" validate:"required"`
}

type addBookmarkRequest struct {
	BookID int64  `json:"book_id" validate:"required"`
	Page   int32  `json:"page" validate:"required"`
	Note   string `json:"note"`
}

func (h *ReadingHandler) SaveProgress(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	var req saveProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	_, err := h.readingClient.SaveProgress(c.Request.Context(), &readingpb.SaveProgressRequest{
		UserId:     userID,
		BookId:     req.BookID,
		FileId:     req.FileID,
		Page:       req.Page,
		Percentage: req.Percentage,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "reading service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ReadingHandler) GetProgress(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	bookID, err := strconv.ParseInt(c.Param("bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bookID"})
		return
	}
	resp, err := h.readingClient.GetProgress(c.Request.Context(), &readingpb.GetProgressRequest{
		UserId: userID,
		BookId: bookID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "reading service error"})
		return
	}
	c.JSON(http.StatusOK, resp.Progress)
}

func (h *ReadingHandler) AddBookmark(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	var req addBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	resp, err := h.readingClient.AddBookmark(c.Request.Context(), &readingpb.AddBookmarkRequest{
		UserId: userID,
		BookId: req.BookID,
		Page:   req.Page,
		Note:   req.Note,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "reading service error"})
		return
	}
	c.JSON(http.StatusOK, resp.Bookmark)
}

func (h *ReadingHandler) GetBookmarks(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	bookID, err := strconv.ParseInt(c.Param("bookID"), 10, 64)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid bookID"})
		return
	}
	resp, err := h.readingClient.GetBookmarks(c.Request.Context(), &readingpb.GetBookmarksRequest{
		UserId: userID,
		BookId: bookID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "reading service error"})
		return
	}
	c.JSON(http.StatusOK, resp.Bookmarks)
}

func (h *ReadingHandler) DeleteBookmark(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserIDKey)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	_, err = h.readingClient.DeleteBookmark(c.Request.Context(), &readingpb.DeleteBookmarkRequest{
		UserId:     userID,
		BookmarkId: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "reading service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

