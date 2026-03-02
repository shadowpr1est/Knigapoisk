package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shadowpr1est/knigapoisk-api-gateway/internal/client"
	authpb "github.com/shadowpr1est/knigapoisk-auth-service/api/proto"
)

const (
	ContextUserIDKey = "userID"
	ContextRoleKey   = "role"
)

type AuthMiddleware struct {
	authClient client.AuthClient
}

func NewAuthMiddleware(authClient client.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}
		token := parts[1]

		resp, err := m.authClient.Validate(c.Request.Context(), &authpb.ValidateRequest{
			AccessToken: token,
		})
		if err != nil || !resp.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(ContextUserIDKey, resp.UserId)
		c.Set(ContextRoleKey, resp.Role)
		c.Request.Header.Set("X-User-ID", strconv.FormatInt(resp.UserId, 10))

		c.Next()
	}
}

