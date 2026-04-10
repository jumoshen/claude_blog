package service

import (
	"testing"
	"time"
)

func TestPostInfo_ScheduledAt(t *testing.T) {
	scheduledTime := time.Now().Add(24 * time.Hour)
	info := PostInfo{
		Title:       "Scheduled Post",
		Slug:        "scheduled-post",
		ScheduledAt: &scheduledTime,
	}

	if info.ScheduledAt == nil {
		t.Error("Expected ScheduledAt to be set")
	}
	if info.Title != "Scheduled Post" {
		t.Errorf("Expected Title 'Scheduled Post', got '%s'", info.Title)
	}
}

func TestPostInfo_NoScheduledAt(t *testing.T) {
	info := PostInfo{
		Title:       "Immediate Post",
		Slug:        "immediate-post",
		ScheduledAt: nil,
	}

	if info.ScheduledAt != nil {
		t.Error("Expected ScheduledAt to be nil")
	}
}
