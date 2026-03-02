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

	authpb "github.com/shadowpr1est/knigapoisk-auth-service/api/proto"
	"github.com/shadowpr1est/knigapoisk-auth-service/config"
	authgrpc "github.com/shadowpr1est/knigapoisk-auth-service/internal/delivery/grpc"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/jwt"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/logger"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/repository/postgres"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/usecase/auth"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	baseLogger := logger.NewLogger()
	logg := logger.WithService(baseLogger, "auth-service")
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

	userRepo := postgres.NewUserRepo(db)
	tokenRepo := postgres.NewTokenRepo(db)
	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTTTL)
	authUseCase := auth.NewAuthUseCase(userRepo, tokenRepo, jwtManager)

	serverLogger := logg.With(zap.String("component", "grpc-server"))

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authgrpc.LoggingInterceptor(serverLogger),
			authgrpc.RecoveryInterceptor(serverLogger),
		),
	)

	authServer := authgrpc.NewAuthServer(authUseCase, serverLogger)
	authpb.RegisterAuthServiceServer(grpcServer, authServer)
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

