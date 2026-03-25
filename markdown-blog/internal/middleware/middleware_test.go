package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"markdown-blog/internal/pkg/jwt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter(jwtUtil *jwt.JWT) *gin.Engine {
	r := gin.New()
	r.Use(AuthMiddleware(jwtUtil))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return r
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	jwtUtil := jwt.NewWithConfig("test-secret", 3600, "test-issuer", nil)
	router := setupTestRouter(jwtUtil)

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	jwtUtil := jwt.NewWithConfig("test-secret", 3600, "test-issuer", nil)
	router := setupTestRouter(jwtUtil)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic invalid")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtUtil := jwt.NewWithConfig("test-secret", 3600, "test-issuer", nil)
	router := setupTestRouter(jwtUtil)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtUtil := jwt.NewWithConfig("test-secret", 3600, "test-issuer", nil)
	token, _, err := jwtUtil.GenerateToken(123, "testuser", "Test User", "http://avatar.url", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	router := setupTestRouter(jwtUtil)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	jwtUtil1 := jwt.NewWithConfig("secret-1", 3600, "test-issuer", nil)
	jwtUtil2 := jwt.NewWithConfig("secret-2", 3600, "test-issuer", nil)

	token, _, err := jwtUtil1.GenerateToken(123, "testuser", "Test User", "http://avatar.url", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	router := setupTestRouter(jwtUtil2)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOptionalAuthMiddleware_NoHeader(t *testing.T) {
	jwtUtil := jwt.NewWithConfig("test-secret", 3600, "test-issuer", nil)

	r := gin.New()
	r.Use(OptionalAuthMiddleware(jwtUtil))
	r.GET("/test", func(c *gin.Context) {
		claims := GetUserClaims(c)
		if claims == nil {
			c.JSON(200, gin.H{"authenticated": false})
		} else {
			c.JSON(200, gin.H{"authenticated": true})
		}
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOptionalAuthMiddleware_InvalidToken(t *testing.T) {
	jwtUtil := jwt.NewWithConfig("test-secret", 3600, "test-issuer", nil)

	r := gin.New()
	r.Use(OptionalAuthMiddleware(jwtUtil))
	r.GET("/test", func(c *gin.Context) {
		claims := GetUserClaims(c)
		if claims == nil {
			c.JSON(200, gin.H{"authenticated": false})
		} else {
			c.JSON(200, gin.H{"authenticated": true})
		}
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOptionalAuthMiddleware_ValidToken(t *testing.T) {
	jwtUtil := jwt.NewWithConfig("test-secret", 3600, "test-issuer", nil)
	token, _, err := jwtUtil.GenerateToken(123, "testuser", "Test User", "http://avatar.url", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	r := gin.New()
	r.Use(OptionalAuthMiddleware(jwtUtil))
	r.GET("/test", func(c *gin.Context) {
		claims := GetUserClaims(c)
		if claims == nil {
			c.JSON(200, gin.H{"authenticated": false})
		} else {
			c.JSON(200, gin.H{"authenticated": true, "user_id": claims.UserID})
		}
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetUserClaims_NoClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	claims := GetUserClaims(c)
	if claims != nil {
		t.Error("Expected nil claims when none set")
	}
}

func TestGetUserClaims_WithClaims(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedClaims := &jwt.Claims{
		UserID: 123,
		Login:  "testuser",
		Name:   "Test User",
	}
	c.Set(UserContextKey, expectedClaims)

	claims := GetUserClaims(c)
	if claims == nil {
		t.Fatal("Expected claims to be returned")
	}
	if claims.UserID != 123 {
		t.Errorf("Expected UserID 123, got %d", claims.UserID)
	}
	if claims.Login != "testuser" {
		t.Errorf("Expected Login 'testuser', got '%s'", claims.Login)
	}
}

func TestGetUserClaims_WrongType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set wrong type
	c.Set(UserContextKey, "not a claims pointer")

	claims := GetUserClaims(c)
	if claims != nil {
		t.Error("Expected nil claims when wrong type set")
	}
}

func TestConstants(t *testing.T) {
	if AuthorizationHeader != "Authorization" {
		t.Errorf("Expected AuthorizationHeader to be 'Authorization', got '%s'", AuthorizationHeader)
	}
	if BearerPrefix != "Bearer " {
		t.Errorf("Expected BearerPrefix to be 'Bearer ', got '%s'", BearerPrefix)
	}
	if UserContextKey != "user" {
		t.Errorf("Expected UserContextKey to be 'user', got '%s'", UserContextKey)
	}
}
