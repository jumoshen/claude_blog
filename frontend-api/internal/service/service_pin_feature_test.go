package service

import (
	"testing"
)

func TestPostInfo_PinFeatureFields(t *testing.T) {
	info := PostInfo{
		Title:      "Pinned Post",
		Slug:       "pinned-post",
		IsPinned:   true,
		IsFeatured: true,
	}

	if !info.IsPinned {
		t.Error("Expected IsPinned to be true")
	}
	if !info.IsFeatured {
		t.Error("Expected IsFeatured to be true")
	}
	if info.Title != "Pinned Post" {
		t.Errorf("Expected Title 'Pinned Post', got '%s'", info.Title)
	}
}

func TestPostInfo_NotPinnedNotFeatured(t *testing.T) {
	info := PostInfo{
		Title:      "Regular Post",
		Slug:       "regular-post",
		IsPinned:   false,
		IsFeatured: false,
	}

	if info.IsPinned {
		t.Error("Expected IsPinned to be false")
	}
	if info.IsFeatured {
		t.Error("Expected IsFeatured to be false")
	}
}
