package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/shadowpr1est/knigapoisk-api-gateway/config"
	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/client"
	httpdelivery "github.com/shadowpr1est/knigapoisk-api-gateway/internal/delivery/http"
	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/delivery/http/handler"
	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/delivery/http/middleware"
	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/logger"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	baseLogger := logger.NewLogger()
	logg := logger.WithService(baseLogger, "api-gateway")
	defer func() { _ = logg.Sync() }()

	logg.Info("starting api-gateway")

	authClient, err := client.NewAuthClient(cfg.AuthServiceAddr)
	if err != nil {
		logg.Fatal("failed to create auth client", zap.Error(err))
	}
	bookClient, err := client.NewBookClient(cfg.BookServiceAddr)
	if err != nil {
		logg.Fatal("failed to create book client", zap.Error(err))
	}
	fileClient, err := client.NewFileClient(cfg.FileServiceAddr)
	if err != nil {
		logg.Fatal("failed to create file client", zap.Error(err))
	}
	readingClient, err := client.NewReadingClient(cfg.ReadingServiceAddr)
	if err != nil {
		logg.Fatal("failed to create reading client", zap.Error(err))
	}
	reviewClient, err := client.NewReviewClient(cfg.ReviewServiceAddr)
	if err != nil {
		logg.Fatal("failed to create review client", zap.Error(err))
	}

	authHandler := handler.NewAuthHandler(authClient)
	bookHandler := handler.NewBookHandler(bookClient)
	fileHandler := handler.NewFileHandler(fileClient)
	readingHandler := handler.NewReadingHandler(readingClient)
	reviewHandler := handler.NewReviewHandler(reviewClient)

	authMiddleware := middleware.NewAuthMiddleware(authClient)

	router := httpdelivery.NewRouter(
		authHandler,
		bookHandler,
		fileHandler,
		readingHandler,
		reviewHandler,
		authMiddleware,
	)

	server := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router.Engine(),
	}

	go func() {
		logg.Info("http server listening", zap.String("addr", cfg.HTTPPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("http server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logg.Info("shutdown signal received", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logg.Error("http server shutdown error", zap.Error(err))
	}

	logg.Info("api-gateway stopped")
}

