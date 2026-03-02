package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/client"
	authpb "github.com/shadowpr1est/knigapoisk-auth-service/api/proto"
)

type AuthHandler struct {
	authClient client.AuthClient
	validate   *validator.Validate
}

func NewAuthHandler(authClient client.AuthClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		validate:   validator.New(),
	}
}

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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

	resp, err := h.authClient.Register(c.Request.Context(), &authpb.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "auth service error"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
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

	resp, err := h.authClient.Login(c.Request.Context(), &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.authClient.Refresh(c.Request.Context(), &authpb.RefreshRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	_, err := h.authClient.Logout(c.Request.Context(), &authpb.LogoutRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "auth service error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

