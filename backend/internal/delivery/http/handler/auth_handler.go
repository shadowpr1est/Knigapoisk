package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shadowpr1est/Knigapoisk/internal/usecase/auth"
)

type AuthHandler struct {
	authUseCase auth.UseCase
	validate    *validator.Validate
}

func NewAuthHandler(authUseCase auth.UseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	output, err := h.authUseCase.Register(c.Request.Context(), auth.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case auth.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{Error: "user already exists"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		}
		return
	}
	c.SetCookie("refresh_token", output.RefreshToken, 30*24*60*60, "/", "", false, true)
	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	output, err := h.authUseCase.Login(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid login or password"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		}
		return
	}
	c.SetCookie("refresh_token", output.RefreshToken, 30*24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "refresh token not found"})
		return
	}
	output, err := h.authUseCase.Refresh(c.Request.Context(), auth.RefreshInput{
		RefreshToken: refreshToken,
	})
	if err != nil {
		switch err {
		case auth.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid refresh token"})
		case auth.ErrTokenExpired:
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "refresh token expired"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		}
		return
	}
	c.SetCookie("refresh_token", output.RefreshToken, 30*24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "refresh token not found"})
		return
	}
	if err := h.authUseCase.Logout(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	if err := h.authUseCase.LogoutAll(c.Request.Context(), userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out from all devices"})
}
