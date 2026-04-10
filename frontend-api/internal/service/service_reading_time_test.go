package service

import (
	"testing"
)

func TestCalculateReadingTime_ChineseContent(t *testing.T) {
	// 1200个中文字符 = 3分钟
	content := ""
	for i := 0; i < 1200; i++ {
		content += "中"
	}
	readingTime := CalculateReadingTime(content)
	if readingTime != 3 {
		t.Errorf("Expected 3 minutes for 1200 Chinese chars, got %d", readingTime)
	}
}

func TestCalculateReadingTime_EnglishContent(t *testing.T) {
	// 400个英文单词 = 2分钟
	content := ""
	for i := 0; i < 400; i++ {
		if i > 0 {
			content += " "
		}
		content += "word"
	}
	readingTime := CalculateReadingTime(content)
	if readingTime != 2 {
		t.Errorf("Expected 2 minutes for 400 English words, got %d", readingTime)
	}
}

func TestCalculateReadingTime_MixedContent(t *testing.T) {
	// 400中文 + 200英文 = 1 + 1 = 2分钟
	content := ""
	for i := 0; i < 400; i++ {
		content += "中"
	}
	for i := 0; i < 200; i++ {
		if i > 0 {
			content += " "
		}
		content += "word"
	}
	readingTime := CalculateReadingTime(content)
	if readingTime != 2 {
		t.Errorf("Expected 2 minutes for mixed content, got %d", readingTime)
	}
}

func TestCalculateReadingTime_ShortContent(t *testing.T) {
	// 100个字符 = 不到1分钟，但应该返回1
	content := "这是一些中文字符"
	readingTime := CalculateReadingTime(content)
	if readingTime < 1 {
		t.Errorf("Expected at least 1 minute for short content, got %d", readingTime)
	}
}

func TestCalculateReadingTime_EmptyContent(t *testing.T) {
	readingTime := CalculateReadingTime("")
	if readingTime < 1 {
		t.Errorf("Expected at least 1 minute for empty content, got %d", readingTime)
	}
}

func TestCalculateReadingTime_WithHTML(t *testing.T) {
	// HTML标签应该被移除
	content := "<h1>标题</h1><p>这是正文内容</p>"
	readingTime := CalculateReadingTime(content)
	if readingTime < 1 {
		t.Errorf("Expected at least 1 minute, got %d", readingTime)
	}
}

func TestPostInfo_ReadingTime(t *testing.T) {
	info := PostInfo{
		Title:       "Test Post",
		Slug:        "test-post",
		ReadingTime: 5,
	}

	if info.ReadingTime != 5 {
		t.Errorf("Expected ReadingTime 5, got %d", info.ReadingTime)
	}
}
