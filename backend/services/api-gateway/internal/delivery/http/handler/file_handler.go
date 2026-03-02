package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/client"
	filepb "github.com/shadowpr1est/knigapoisk-file-service/api/proto"
)

type FileHandler struct {
	fileClient client.FileClient
}

func NewFileHandler(fileClient client.FileClient) *FileHandler {
	return &FileHandler{fileClient: fileClient}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	bookIDStr := c.PostForm("book_id")
	formatStr := c.PostForm("format")
	if bookIDStr == "" || formatStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "book_id and format required"})
		return
	}
	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil || bookID <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid book_id"})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "file required"})
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to open file"})
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to read file"})
		return
	}
	resp, err := h.fileClient.UploadFile(c.Request.Context(), &filepb.UploadFileRequest{
		BookId: bookID,
		Format: toProtoFileFormat(formatStr),
		Data:   data,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "file service error"})
		return
	}
	c.JSON(http.StatusOK, resp.File)
}

func (h *FileHandler) GetFile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid id"})
		return
	}
	resp, err := h.fileClient.GetFile(c.Request.Context(), &filepb.GetFileRequest{Id: id})
	if err != nil || resp.File == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "file not found"})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\"file\"")
	_, _ = c.Writer.Write(resp.Data)
}

func toProtoFileFormat(s string) filepb.FileFormat {
	switch s {
	case "pdf":
		return filepb.FileFormat_FILE_FORMAT_PDF
	case "epub":
		return filepb.FileFormat_FILE_FORMAT_EPUB
	default:
		return filepb.FileFormat_FILE_FORMAT_UNSPECIFIED
	}
}

