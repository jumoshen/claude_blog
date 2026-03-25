package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/pkg/jwt"
	"markdown-blog/internal/pkg/response"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix       = "Bearer "
	UserContextKey     = "user"
)

func AuthMiddleware(jwtUtil *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set(UserContextKey, claims)
		c.Next()
	}
}

func OptionalAuthMiddleware(jwtUtil *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserContextKey, claims)
		c.Next()
	}
}

func GetUserClaims(c *gin.Context) *jwt.Claims {
	value, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	claims, ok := value.(*jwt.Claims)
	if !ok {
		return nil
	}
	return claims
}
