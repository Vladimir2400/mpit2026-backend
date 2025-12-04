package middleware

import (
	"net/http"
	"strings"

	"github.com/gdugdh24/mpit2026-backend/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authUseCase *auth.VKAuthUseCase
}

func NewAuthMiddleware(authUseCase *auth.VKAuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

// RequireAuth is a middleware that validates JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Verify token
		userID, err := m.authUseCase.VerifyToken(c.Request.Context(), token)
		if err != nil {
			statusCode := http.StatusUnauthorized
			message := "invalid token"

			switch err.Error() {
			case "session not found":
				message = "session not found"
			case "session expired":
				message = "session expired"
			}

			c.JSON(statusCode, gin.H{
				"error": message,
			})
			c.Abort()
			return
		}

		// Set user_id in context for handlers
		c.Set("user_id", userID)
		c.Set("token", token)

		c.Next()
	}
}

// OptionalAuth is a middleware that validates token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		userID, err := m.authUseCase.VerifyToken(c.Request.Context(), token)
		if err == nil {
			c.Set("user_id", userID)
			c.Set("token", token)
		}

		c.Next()
	}
}
