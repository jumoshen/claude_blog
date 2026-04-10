package service

import (
	"testing"

	"markdown-blog/internal/model"
)

func TestCategoryInfo_Fields(t *testing.T) {
	info := CategoryInfo{
		Name:        "Go",
		Slug:        "go",
		Description: "Go语言",
		Sort:        10,
		PostCount:   5,
	}

	if info.Name != "Go" {
		t.Errorf("Expected Name 'Go', got '%s'", info.Name)
	}
	if info.Slug != "go" {
		t.Errorf("Expected Slug 'go', got '%s'", info.Slug)
	}
	if info.Description != "Go语言" {
		t.Errorf("Expected Description 'Go语言', got '%s'", info.Description)
	}
	if info.Sort != 10 {
		t.Errorf("Expected Sort 10, got %d", info.Sort)
	}
	if info.PostCount != 5 {
		t.Errorf("Expected PostCount 5, got %d", info.PostCount)
	}
}

func TestCategory_ModelConversion(t *testing.T) {
	category := &model.Category{
		Name:        "Python",
		Slug:        "python",
		Description: "Python语言",
		Sort:        20,
	}

	info := CategoryInfo{
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Sort:        category.Sort,
	}

	if info.Name != "Python" {
		t.Errorf("Expected Name 'Python', got '%s'", info.Name)
	}
	if info.Slug != "python" {
		t.Errorf("Expected Slug 'python', got '%s'", info.Slug)
	}
	if info.Description != "Python语言" {
		t.Errorf("Expected Description 'Python语言', got '%s'", info.Description)
	}
	if info.Sort != 20 {
		t.Errorf("Expected Sort 20, got %d", info.Sort)
	}
}
