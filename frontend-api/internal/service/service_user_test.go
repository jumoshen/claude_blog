package service

import (
	"testing"
)

func TestBlogUserInfo_Fields(t *testing.T) {
	user := BlogUserInfo{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Nickname: "Test User",
		Avatar:   "https://example.com/avatar.png",
	}

	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", user.Email)
	}
	if user.Nickname != "Test User" {
		t.Errorf("Expected Nickname 'Test User', got '%s'", user.Nickname)
	}
	if user.Avatar != "https://example.com/avatar.png" {
		t.Errorf("Expected Avatar URL, got '%s'", user.Avatar)
	}
}

func TestBlogUserInfo_EmptyAvatar(t *testing.T) {
	user := BlogUserInfo{
		Username: "user1",
		Email:    "user@example.com",
		Nickname: "User",
		Avatar:   "",
	}

	if user.Avatar != "" {
		t.Errorf("Expected empty Avatar, got '%s'", user.Avatar)
	}
}
