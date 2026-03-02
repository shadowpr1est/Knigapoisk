package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/client"
	bookpb "github.com/shadowpr1est/knigapoisk-book-service/api/proto"
)

type BookHandler struct {
	bookClient client.BookClient
	validate   *validator.Validate
}

func NewBookHandler(bookClient client.BookClient) *BookHandler {
	return &BookHandler{
		bookClient: bookClient,
		validate:   validator.New(),
	}
}

type createBookRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Language    string  `json:"language"`
	PublishedAt string  `json:"published_at"`
	CoverURL    string  `json:"cover_url"`
	FileID      int64   `json:"file_id"`
	Status      string  `json:"status" validate:"required,oneof=active hidden"`
	AuthorIDs   []int64 `json:"author_ids"`
	GenreIDs    []int64 `json:"genre_ids"`
}

type updateBookRequest = createBookRequest

func (h *BookHandler) ListBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.bookClient.ListBooks(c.Request.Context(), &bookpb.ListBooksRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "book service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": resp.Books,
		"total": resp.Total,
	})
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	resp, err := h.bookClient.GetBook(c.Request.Context(), &bookpb.GetBookRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "book not found"})
		return
	}
	c.JSON(http.StatusOK, resp.Book)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	publishedAt, _ := time.Parse(time.RFC3339, req.PublishedAt)
	resp, err := h.bookClient.CreateBook(c.Request.Context(), &bookpb.CreateBookRequest{
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
		PublishedAt: publishedAt.Format(time.RFC3339),
		CoverUrl:    req.CoverURL,
		FileId:      req.FileID,
		Status:      toProtoBookStatus(req.Status),
		AuthorIds:   req.AuthorIDs,
		GenreIds:    req.GenreIDs,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "book service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": resp.Id})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	var req updateBookRequest
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
	publishedAt, _ := time.Parse(time.RFC3339, req.PublishedAt)
	resp, err := h.bookClient.UpdateBook(c.Request.Context(), &bookpb.UpdateBookRequest{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
		PublishedAt: publishedAt.Format(time.RFC3339),
		CoverUrl:    req.CoverURL,
		FileId:      req.FileID,
		Status:      toProtoBookStatus(req.Status),
		AuthorIds:   req.AuthorIDs,
		GenreIds:    req.GenreIDs,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "book service error"})
		return
	}
	c.JSON(http.StatusOK, resp.Book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	_, err = h.bookClient.DeleteBook(c.Request.Context(), &bookpb.DeleteBookRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "book service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *BookHandler) GetBooksByAuthor(c *gin.Context) {
	authorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || authorID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.bookClient.GetBooksByAuthor(c.Request.Context(), &bookpb.GetBooksByAuthorRequest{
		AuthorId: authorID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "book service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": resp.Books,
		"total": resp.Total,
	})
}

func (h *BookHandler) GetBooksByGenre(c *gin.Context) {
	genreID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || genreID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.bookClient.GetBooksByGenre(c.Request.Context(), &bookpb.GetBooksByGenreRequest{
		GenreId: genreID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "book service error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": resp.Books,
		"total": resp.Total,
	})
}

func toProtoBookStatus(s string) bookpb.BookStatus {
	switch s {
	case "active":
		return bookpb.BookStatus_BOOK_STATUS_ACTIVE
	case "hidden":
		return bookpb.BookStatus_BOOK_STATUS_HIDDEN
	default:
		return bookpb.BookStatus_BOOK_STATUS_UNSPECIFIED
	}
}

