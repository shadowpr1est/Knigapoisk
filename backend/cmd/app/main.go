package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/shadowpr1est/Knigapoisk/config"
	"github.com/shadowpr1est/Knigapoisk/internal/delivery/http/handler"
	"github.com/shadowpr1est/Knigapoisk/internal/repository/postgres"
	"github.com/shadowpr1est/Knigapoisk/internal/usecase/auth"
	"github.com/shadowpr1est/Knigapoisk/pkg/jwt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found:", err)
	}
	cfg := config.Load()

	db, err := postgres.NewDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	jwtManager := jwt.NewManager(cfg.JWT.SecretKey, cfg.JWT.TTL)

	userRepo := postgres.NewUserRepo(db)
	tokenRepo := postgres.NewTokenRepo(db)

	authUseCase := auth.NewAuthUseCase(userRepo, tokenRepo, jwtManager)

	authHandler := handler.NewAuthHandler(authUseCase)

	router := handler.NewRouter(authHandler, jwtManager)

	if err := router.InitRoutes().Run(cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
