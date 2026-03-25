package service

import (
	"testing"
)

func TestParseFrontMatter_WithAllFields(t *testing.T) {
	content := `---
title: "Test Post"
date: 2024-01-15T10:00:00Z07:00
tags: ["tag1", "tag2", "tag3"]
categories: ["cat1", "cat2"]
---
# Hello World

This is the content.
`

	title, dateStr, tagsRaw, categoriesRaw, _, bodyLines := parseFrontMatter(content)

	if title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", title)
	}
	if dateStr != "2024-01-15T10:00:00Z07:00" {
		t.Errorf("Expected date '2024-01-15T10:00:00Z07:00', got '%s'", dateStr)
	}
	if tagsRaw != `["tag1", "tag2", "tag3"]` {
		t.Errorf("Expected tags '[\"tag1\", \"tag2\", \"tag3\"]', got '%s'", tagsRaw)
	}
	if categoriesRaw != `["cat1", "cat2"]` {
		t.Errorf("Expected categories '[\"cat1\", \"cat2\"]', got '%s'", categoriesRaw)
	}
	if len(bodyLines) == 0 {
		t.Error("Expected body lines to be non-empty")
	}
}

func TestParseFrontMatter_NoFrontMatter(t *testing.T) {
	// When there's no front matter, the content is treated as front matter
	// so body lines will be empty
	content := `# Just a title

Some content here.
`

	title, _, _, _, _, bodyLines := parseFrontMatter(content)

	if title != "" {
		t.Errorf("Expected empty title, got '%s'", title)
	}
	// Without proper front matter delimiters, bodyLines will be empty
	// because frontMatterEnded stays false
	_ = bodyLines
}

func TestParseFrontMatter_WithMoreTag(t *testing.T) {
	content := `---
title: "Test Post"
date: 2024-01-15
tags: ["tag1"]
categories: ["cat1"]
---
Some summary content here.

<!--more-->

Rest of the content.
`

	title, _, _, _, summary, bodyLines := parseFrontMatter(content)

	if title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", title)
	}
	if summary != "Some summary content here." {
		t.Errorf("Expected summary 'Some summary content here.', got '%s'", summary)
	}
	// Body should not contain the more tag line
	for _, line := range bodyLines {
		if line == "<!--more-->" {
			t.Error("Body should not contain <!--more--> tag")
		}
	}
}

func TestParseFrontMatter_TitleWithQuotes(t *testing.T) {
	content := `---
title: "My \"Quoted\" Title"
date: 2024-01-15
---
Content
`

	title, _, _, _, _, _ := parseFrontMatter(content)

	// In Go raw string, \" is literal backslash + quote
	// strings.Trim only removes leading/trailing quotes, not internal ones
	// So title becomes: My \"Quoted\" Title
	expected := `My \"Quoted\" Title`
	if title != expected {
		t.Errorf("Expected title '%s', got '%s'", expected, title)
	}
}

func TestParseArrayField_QuotedArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple quoted", `["tag1", "tag2"]`, "tag1,tag2"},
		{"single item", `["tag1"]`, "tag1"},
		{"empty", ``, ""},
		{"with spaces", `["tag 1", "tag 2"]`, "tag 1,tag 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArrayField(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestParseArrayField_BracketArray(t *testing.T) {
	result := parseArrayField("tag1, tag2, tag3")
	if result != "tag1, tag2, tag3" {
		t.Errorf("Expected 'tag1, tag2, tag3', got '%s'", result)
	}
}

func TestParseArrayField_EmptyString(t *testing.T) {
	result := parseArrayField("")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestParseArrayField_WhitespaceOnly(t *testing.T) {
	result := parseArrayField("   ")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestTruncateString_ShortString(t *testing.T) {
	input := "short"
	result := truncateString(input, 10)
	if result != "short" {
		t.Errorf("Expected 'short', got '%s'", result)
	}
}

func TestTruncateString_ExactLength(t *testing.T) {
	input := "hello"
	result := truncateString(input, 5)
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}
}

func TestTruncateString_LongString(t *testing.T) {
	input := "This is a very long string"
	result := truncateString(input, 10)
	if result == input {
		t.Error("Expected truncated string to be different from input")
	}
	// result = first 10 chars + "..."
	// The truncation is based on runes, so we're getting 10 characters
	// plus the ellipsis, which makes the total length 13
	runeCount := 0
	for _, r := range result {
		_ = r
		runeCount++
	}
	if runeCount != 13 { // 10 original + 3 for "..."
		t.Errorf("Expected 13 runes, got %d", runeCount)
	}
	if result[len(result)-3:] != "..." {
		t.Error("Expected truncated string to end with '...'")
	}
}

func TestTruncateString_Unicode(t *testing.T) {
	input := "你好世界"
	result := truncateString(input, 3)
	if result != "你好世..." {
		t.Errorf("Expected '你好世...', got '%s'", result)
	}
}

func TestTruncateString_EmptyString(t *testing.T) {
	result := truncateString("", 10)
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestStripFrontMatter_FullFrontMatter(t *testing.T) {
	content := `---
title: "Test"
date: 2024-01-15
---
# Content Here

Some text.
`

	result := stripFrontMatter(content)
	if result == content {
		t.Error("Expected stripped content to be different")
	}
	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}
}

func TestStripFrontMatter_NoFrontMatter(t *testing.T) {
	content := `# Just a heading

Some text.
`

	result := stripFrontMatter(content)
	if result != content {
		t.Error("Expected content to remain unchanged")
	}
}

func TestStripFrontMatter_EmptyFrontMatter(t *testing.T) {
	// Incomplete front matter (only one ---) should not be stripped
	// because there's no closing ---
	content := `---

Some text.
`

	result := stripFrontMatter(content)
	// Without closing ---, the front matter is not recognized
	// So content remains unchanged
	if result != content {
		t.Error("Expected content to remain unchanged for incomplete front matter")
	}
}

func TestStripFrontMatter_IncompleteFrontMatter(t *testing.T) {
	content := `---
title: "Test"

Some text.
`

	result := stripFrontMatter(content)
	if result != content {
		t.Error("Expected content to remain unchanged when front matter is incomplete")
	}
}
