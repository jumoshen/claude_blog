package service

import (
	"testing"
)

func TestPostNavigation_WithPrevAndNext(t *testing.T) {
	nav := PostNavigation{
		Prev: &PostNavItem{Slug: "prev-post", Title: "Previous Post"},
		Next: &PostNavItem{Slug: "next-post", Title: "Next Post"},
	}

	if nav.Prev == nil {
		t.Error("Expected Prev to be set")
	}
	if nav.Next == nil {
		t.Error("Expected Next to be set")
	}
	if nav.Prev.Slug != "prev-post" {
		t.Errorf("Expected Prev.Slug 'prev-post', got '%s'", nav.Prev.Slug)
	}
	if nav.Next.Slug != "next-post" {
		t.Errorf("Expected Next.Slug 'next-post', got '%s'", nav.Next.Slug)
	}
}

func TestPostNavigation_OnlyPrev(t *testing.T) {
	nav := PostNavigation{
		Prev: &PostNavItem{Slug: "first-post", Title: "First Post"},
		Next: nil,
	}

	if nav.Prev == nil {
		t.Error("Expected Prev to be set")
	}
	if nav.Next != nil {
		t.Error("Expected Next to be nil")
	}
}

func TestPostNavigation_OnlyNext(t *testing.T) {
	nav := PostNavigation{
		Prev: nil,
		Next: &PostNavItem{Slug: "last-post", Title: "Last Post"},
	}

	if nav.Prev != nil {
		t.Error("Expected Prev to be nil")
	}
	if nav.Next == nil {
		t.Error("Expected Next to be set")
	}
}

func TestPostNavigation_Empty(t *testing.T) {
	nav := PostNavigation{
		Prev: nil,
		Next: nil,
	}

	if nav.Prev != nil {
		t.Error("Expected Prev to be nil")
	}
	if nav.Next != nil {
		t.Error("Expected Next to be nil")
	}
}

func TestPostNavItem_Fields(t *testing.T) {
	item := PostNavItem{
		Slug:  "test-slug",
		Title: "Test Title",
	}

	if item.Slug != "test-slug" {
		t.Errorf("Expected Slug 'test-slug', got '%s'", item.Slug)
	}
	if item.Title != "Test Title" {
		t.Errorf("Expected Title 'Test Title', got '%s'", item.Title)
	}
}
