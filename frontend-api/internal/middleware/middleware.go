package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/model"
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

// VisitLogger 访问日志中间件
func VisitLogger(createVisit func(*model.Visit) error, getUserID func(*gin.Context) int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 只记录成功的 GET 请求
		if c.Request.Method == "GET" && c.Writer.Status() == 200 {
			slug := c.Param("slug")
			if slug != "" {
				visit := &model.Visit{
					PostSlug:  slug,
					UserID:    getUserID(c),
					IP:        c.ClientIP(),
					UserAgent: c.Request.UserAgent(),
				}
				// 异步记录访问
				go func() {
					createVisit(visit)
				}()
			}
		}
	}
}

// RequirePermission 权限检查中间件
func RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetUserClaims(c)
		if claims == nil {
			response.Unauthorized(c, "Not logged in")
			c.Abort()
			return
		}

		// 从 context 获取用户权限
		permissions, exists := c.Get("permissions")
		if !exists {
			// 如果没有权限信息，默认放行（兼容旧逻辑）
			c.Next()
			return
		}

		perms, ok := permissions.([]string)
		if !ok {
			c.Next()
			return
		}

		// 检查是否有 * 权限（超级管理员）
		for _, p := range perms {
			if p == "*" {
				c.Next()
				return
			}
			if p == requiredPermission {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Permission denied")
		c.Abort()
	}
}
