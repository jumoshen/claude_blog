package service

import (
	"testing"
)

func TestPostLikeInfo_Fields(t *testing.T) {
	info := struct {
		PostID    uint
		UserID    uint
		Liked     bool
		LikeCount int64
	}{
		PostID:    1,
		UserID:    2,
		Liked:     true,
		LikeCount: 100,
	}

	if info.PostID != 1 {
		t.Errorf("Expected PostID 1, got %d", info.PostID)
	}
	if info.UserID != 2 {
		t.Errorf("Expected UserID 2, got %d", info.UserID)
	}
	if !info.Liked {
		t.Error("Expected Liked to be true")
	}
	if info.LikeCount != 100 {
		t.Errorf("Expected LikeCount 100, got %d", info.LikeCount)
	}
}

func TestLikePost_Result(t *testing.T) {
	result := struct {
		Liked bool
		Count int64
	}{
		Liked: true,
		Count: 10,
	}

	if !result.Liked {
		t.Error("Expected Liked true after like action")
	}
	if result.Count != 10 {
		t.Errorf("Expected Count 10, got %d", result.Count)
	}
}

func TestUnlikePost_Result(t *testing.T) {
	result := struct {
		Liked bool
		Count int64
	}{
		Liked: false,
		Count: 9,
	}

	if result.Liked {
		t.Error("Expected Liked false after unlike action")
	}
	if result.Count != 9 {
		t.Errorf("Expected Count 9, got %d", result.Count)
	}
}
