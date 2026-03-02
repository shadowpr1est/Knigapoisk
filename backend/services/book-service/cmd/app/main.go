package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	bookpb "github.com/shadowpr1est/knigapoisk-book-service/api/proto"
	"github.com/shadowpr1est/knigapoisk-book-service/config"
	bookgrpc "github.com/shadowpr1est/knigapoisk-book-service/internal/delivery/grpc"
	"github.com/shadowpr1est/knigapoisk-book-service/internal/logger"
	"github.com/shadowpr1est/knigapoisk-book-service/internal/repository/postgres"
	"github.com/shadowpr1est/knigapoisk-book-service/internal/usecase/book"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	baseLogger := logger.NewLogger()
	logg := logger.WithService(baseLogger, "book-service")
	defer func() { _ = logg.Sync() }()

	logg.Info("starting service")

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN)
	if err != nil {
		logg.Fatal("failed to connect to db", zap.Error(err))
	}
	logg.Info("connected to db")

	if err := runMigrations(logg, db); err != nil {
		logg.Fatal("failed to run migrations", zap.Error(err))
	}

	bookRepo := postgres.NewBookRepo(db)
	uc := book.NewBookUseCase(bookRepo, nil, nil)

	serverLogger := logg.With(zap.String("component", "grpc-server"))

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			bookgrpc.LoggingInterceptor(serverLogger),
			bookgrpc.RecoveryInterceptor(serverLogger),
		),
	)

	bookServer := bookgrpc.NewBookServer(uc, serverLogger)
	bookpb.RegisterBookServiceServer(grpcServer, bookServer)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logg.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		logg.Info("grpc server listening", zap.String("addr", cfg.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			logg.Fatal("grpc server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logg.Info("shutdown signal received", zap.String("signal", sig.String()))

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		logg.Warn("force stopping grpc server")
		grpcServer.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdownDB(ctx, db); err != nil {
		logg.Error("db shutdown error", zap.Error(err))
	}

	logg.Info("service stopped")
}

func runMigrations(logger *zap.Logger, db *sqlx.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	logger.Info("running migrations", zap.String("dir", "./migrations"))
	if err := goose.Up(db.DB, "./migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func shutdownDB(ctx context.Context, db *sqlx.DB) error {
	ch := make(chan error, 1)
	go func() {
		ch <- db.Close()
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return fmt.Errorf("db close timeout: %w", ctx.Err())
	}
}

