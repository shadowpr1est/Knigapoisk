package http

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/delivery/http/handler"
	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/delivery/http/middleware"
)

type Router struct {
	engine         *gin.Engine
	authHandler    *handler.AuthHandler
	bookHandler    *handler.BookHandler
	fileHandler    *handler.FileHandler
	readingHandler *handler.ReadingHandler
	reviewHandler  *handler.ReviewHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *handler.AuthHandler,
	bookHandler *handler.BookHandler,
	fileHandler *handler.FileHandler,
	readingHandler *handler.ReadingHandler,
	reviewHandler *handler.ReviewHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	r := gin.New()
	r.Use(gin.Recovery())

	router := &Router{
		engine:         r,
		authHandler:    authHandler,
		bookHandler:    bookHandler,
		fileHandler:    fileHandler,
		readingHandler: readingHandler,
		reviewHandler:  reviewHandler,
		authMiddleware: authMiddleware,
	}
	router.registerRoutes()
	return router
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) registerRoutes() {
	r.engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.engine.GET("/openapi.yaml", func(c *gin.Context) {
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.File(filepath.Join(".", "openapi.yaml"))
	})
	r.engine.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.engine.Group("/api/v1")

	api.POST("/auth/register", r.authHandler.Register)
	api.POST("/auth/login", r.authHandler.Login)
	api.POST("/auth/refresh", r.authHandler.Refresh)

	// public endpoints
	api.GET("/books", r.bookHandler.ListBooks)
	api.GET("/books/:id", r.bookHandler.GetBook)
	api.GET("/reviews/:bookID", r.reviewHandler.GetReviews)
	api.GET("/reviews/:bookID/rating", r.reviewHandler.GetRating)

	private := api.Group("")
	private.Use(r.authMiddleware.Handler())

	private.POST("/auth/logout", r.authHandler.Logout)
	// private.POST("/auth/logout-all", ...) // можно добавить позднее

	private.POST("/books", r.bookHandler.CreateBook)
	private.PUT("/books/:id", r.bookHandler.UpdateBook)
	private.DELETE("/books/:id", r.bookHandler.DeleteBook)
	private.GET("/authors/:id/books", r.bookHandler.GetBooksByAuthor)
	private.GET("/genres/:id/books", r.bookHandler.GetBooksByGenre)

	private.POST("/files/upload", r.fileHandler.UploadFile)
	private.GET("/files/:id", r.fileHandler.GetFile)

	private.GET("/reading/progress/:bookID", r.readingHandler.GetProgress)
	private.POST("/reading/progress", r.readingHandler.SaveProgress)
	private.GET("/reading/bookmarks/:bookID", r.readingHandler.GetBookmarks)
	private.POST("/reading/bookmarks", r.readingHandler.AddBookmark)
	private.DELETE("/reading/bookmarks/:id", r.readingHandler.DeleteBookmark)

	private.POST("/reviews", r.reviewHandler.CreateReview)
	private.PUT("/reviews/:id", r.reviewHandler.UpdateReview)
	private.DELETE("/reviews/:id", r.reviewHandler.DeleteReview)
}

