package service

import (
	"strings"
	"testing"
)

func TestExtractTOC_SingleHeading(t *testing.T) {
	svc := &Service{}
	html := `<h2 id="intro">简介</h2><p>内容</p>`
	toc := svc.ExtractTOC(html)

	if len(toc) != 1 {
		t.Fatalf("Expected 1 TOC item, got %d", len(toc))
	}
	if toc[0].Level != 2 {
		t.Errorf("Expected Level 2, got %d", toc[0].Level)
	}
	if toc[0].Text != "简介" {
		t.Errorf("Expected Text '简介', got '%s'", toc[0].Text)
	}
	if toc[0].ID != "intro" {
		t.Errorf("Expected ID 'intro', got '%s'", toc[0].ID)
	}
}

func TestExtractTOC_MultipleHeadings(t *testing.T) {
	svc := &Service{}
	html := `
		<h2 id="intro">简介</h2>
		<h2 id="install">安装</h2>
		<h3 id="windows">Windows</h3>
		<h3 id="macos">macOS</h3>
		<h1 id="top">顶部标题</h1>
	`
	toc := svc.ExtractTOC(html)

	if len(toc) != 5 {
		t.Fatalf("Expected 5 TOC items, got %d", len(toc))
	}

	// Check first item
	if toc[0].Level != 2 || toc[0].Text != "简介" || toc[0].ID != "intro" {
		t.Errorf("First item mismatch: %+v", toc[0])
	}

	// Check h1
	if toc[4].Level != 1 || toc[4].Text != "顶部标题" || toc[4].ID != "top" {
		t.Errorf("Last item mismatch: %+v", toc[4])
	}
}

func TestExtractTOC_NoHeadings(t *testing.T) {
	svc := &Service{}
	html := `<p>Just a paragraph</p><p>Another paragraph</p>`
	toc := svc.ExtractTOC(html)

	if len(toc) != 0 {
		t.Errorf("Expected 0 TOC items, got %d", len(toc))
	}
}

func TestExtractTOC_EmptyString(t *testing.T) {
	svc := &Service{}
	toc := svc.ExtractTOC("")

	if len(toc) != 0 {
		t.Errorf("Expected 0 TOC items, got %d", len(toc))
	}
}

func TestTocItem_Structure(t *testing.T) {
	item := TocItem{
		Level: 2,
		Text:  "测试标题",
		ID:    "test-title",
	}

	if item.Level != 2 {
		t.Errorf("Expected Level 2, got %d", item.Level)
	}
	if item.Text != "测试标题" {
		t.Errorf("Expected Text '测试标题', got '%s'", item.Text)
	}
	if item.ID != "test-title" {
		t.Errorf("Expected ID 'test-title', got '%s'", item.ID)
	}

	// Simple check - verify fields exist
	result := item.Level > 0 && strings.Contains("测试标题", item.Text)
	if !result {
		t.Error("TocItem fields not properly set")
	}
}
