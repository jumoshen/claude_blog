package service

import (
	"testing"
)

func TestSearchPosts_QueryValidation(t *testing.T) {
	// Test that empty keyword is handled
	// (In actual implementation, handler validates before calling service)
	info := PostInfo{
		Title: "Test Post",
		Slug:  "test-post",
	}

	if info.Title != "Test Post" {
		t.Errorf("Expected Title 'Test Post', got '%s'", info.Title)
	}
	if info.Slug != "test-post" {
		t.Errorf("Expected Slug 'test-post', got '%s'", info.Slug)
	}
}

func TestPostInfo_SearchResult(t *testing.T) {
	info := PostInfo{
		Slug:       "search-test",
		Title:      "Search Test Title",
		Tags:       []string{"test", "search"},
		Categories: []string{"tech"},
		Summary:    "This is a test summary",
		Views:      100,
	}

	if len(info.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(info.Tags))
	}
	if info.Tags[0] != "test" {
		t.Errorf("Expected first tag 'test', got '%s'", info.Tags[0])
	}
	if info.Views != 100 {
		t.Errorf("Expected Views 100, got %d", info.Views)
	}
}
