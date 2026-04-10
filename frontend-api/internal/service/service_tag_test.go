package service

import (
	"testing"

	"markdown-blog/internal/model"
)

func TestTagInfo_Fields(t *testing.T) {
	info := TagInfo{
		Name:      "Go",
		Slug:      "go",
		Color:     "#00FF00",
		PostCount: 10,
	}

	if info.Name != "Go" {
		t.Errorf("Expected Name 'Go', got '%s'", info.Name)
	}
	if info.Slug != "go" {
		t.Errorf("Expected Slug 'go', got '%s'", info.Slug)
	}
	if info.Color != "#00FF00" {
		t.Errorf("Expected Color '#00FF00', got '%s'", info.Color)
	}
	if info.PostCount != 10 {
		t.Errorf("Expected PostCount 10, got %d", info.PostCount)
	}
}

func TestTag_ModelConversion(t *testing.T) {
	tag := &model.Tag{
		Name:  "Python",
		Slug:  "python",
		Color: "#3572A5",
	}

	info := TagInfo{
		Name:  tag.Name,
		Slug:  tag.Slug,
		Color: tag.Color,
	}

	if info.Name != "Python" {
		t.Errorf("Expected Name 'Python', got '%s'", info.Name)
	}
	if info.Slug != "python" {
		t.Errorf("Expected Slug 'python', got '%s'", info.Slug)
	}
	if info.Color != "#3572A5" {
		t.Errorf("Expected Color '#3572A5', got '%s'", info.Color)
	}
}
