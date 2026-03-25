package jwt

import (
	"testing"
	"time"
)

func TestJWT_GenerateAndValidate(t *testing.T) {
	jwtUtil := NewWithConfig("test-secret", 3600, "test-issuer", nil)

	// Generate token
	token, jti, err := jwtUtil.GenerateToken(123, "testuser", "Test User", "http://avatar.url", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}

	if jti == "" {
		t.Fatal("JTI should not be empty")
	}

	// Validate token
	claims, err := jwtUtil.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 123 {
		t.Errorf("Expected UserID 123, got %d", claims.UserID)
	}

	if claims.Login != "testuser" {
		t.Errorf("Expected Login 'testuser', got '%s'", claims.Login)
	}

	if claims.Name != "Test User" {
		t.Errorf("Expected Name 'Test User', got '%s'", claims.Name)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", claims.Email)
	}

	if claims.ID != jti {
		t.Errorf("Expected JTI '%s', got '%s'", jti, claims.ID)
	}

	t.Logf("Token validated successfully: %+v", claims)
}

func TestJWT_InvalidToken(t *testing.T) {
	jwtUtil := NewWithConfig("test-secret", 3600, "test-issuer", nil)

	_, err := jwtUtil.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Expected error for invalid token")
	}

	t.Logf("Invalid token correctly rejected: %v", err)
}

func TestJWT_WrongSecret(t *testing.T) {
	jwtUtil1 := NewWithConfig("secret-1", 3600, "test-issuer", nil)
	jwtUtil2 := NewWithConfig("secret-2", 3600, "test-issuer", nil)

	// Generate with secret 1
	token, _, err := jwtUtil1.GenerateToken(123, "testuser", "Test User", "http://avatar.url", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with secret 2
	_, err = jwtUtil2.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error when validating with wrong secret")
	}

	t.Logf("Wrong secret correctly rejected: %v", err)
}

func TestJWT_Expiration(t *testing.T) {
	// Create JWT with very short expiration (1 second)
	jwtUtil := NewWithConfig("test-secret", 1, "test-issuer", nil)

	token, _, err := jwtUtil.GenerateToken(123, "testuser", "Test User", "http://avatar.url", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Try to validate expired token
	_, err = jwtUtil.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}

	t.Logf("Expired token correctly rejected: %v", err)
}
