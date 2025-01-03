package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	// Add fields for token validation, etc.
	validAPIKeys map[string]bool
}

func NewAuthMiddleware() *AuthMiddleware {
	// In a real application, these would come from a database or environment variables
	// This is just a placeholder for development
	validKeys := map[string]bool{
		"development-token": true,
		"test-token":        true,
	}

	return &AuthMiddleware{
		validAPIKeys: validKeys,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No authorization header provided",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate the token
		if !m.validateToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store user info in context if needed
		// c.Set("user_id", "some-user-id")

		c.Next()
	}
}

func (m *AuthMiddleware) validateToken(token string) bool {
	// For development/testing, just check if it's in our valid keys map
	return m.validAPIKeys[token]

	// In production, you would:
	// 1. Validate JWT token
	// 2. Check token expiration
	// 3. Verify token signature
	// 4. Check token revocation status
	// return jwt.ValidateToken(token)
}
