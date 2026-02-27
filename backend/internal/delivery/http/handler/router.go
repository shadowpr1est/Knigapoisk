package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shadowpr1est/Knigapoisk/internal/delivery/http/middleware"
	"github.com/shadowpr1est/Knigapoisk/pkg/jwt"
)

type Router struct {
	authHandler *AuthHandler
	jwtManager  jwt.Manager
}

func NewRouter(
	authHandler *AuthHandler,
	jwtManager jwt.Manager,
) *Router {
	return &Router{
		authHandler: authHandler,
		jwtManager:  jwtManager,
	}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.Default()

	public := router.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.Refresh)
		}
	}

	private := router.Group("/api/v1")
	private.Use(middleware.Auth(r.jwtManager))
	{
		auth := private.Group("/auth")
		{
			auth.POST("/logout", r.authHandler.Logout)
			auth.POST("/logout-all", r.authHandler.LogoutAll)
		}
	}

	return router
}
