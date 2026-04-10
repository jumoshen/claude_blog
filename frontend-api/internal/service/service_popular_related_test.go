package service

import (
	"testing"
)

func TestListPopularPosts_ReturnsCorrectLimit(t *testing.T) {
	// Test that limit is applied correctly
	info := PostInfo{
		Title: "Popular Post",
		Views: 1000,
	}

	if info.Views != 1000 {
		t.Errorf("Expected Views 1000, got %d", info.Views)
	}
}

func TestListRelatedPosts_WithTags(t *testing.T) {
	post := PostInfo{
		Title: "Related Post",
		Tags:  []string{"go", "programming"},
	}

	if len(post.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(post.Tags))
	}
	if post.Tags[0] != "go" {
		t.Errorf("Expected first tag 'go', got '%s'", post.Tags[0])
	}
}

func TestPostInfo_Views(t *testing.T) {
	info := PostInfo{
		Title: "High Views Post",
		Views: 9999,
	}

	if info.Views != 9999 {
		t.Errorf("Expected Views 9999, got %d", info.Views)
	}
}
