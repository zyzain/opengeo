package content

import (
	"testing"
	"time"
)

func TestContent_Fields(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	c := Content{
		ID:           1,
		UserID:       100,
		Title:        "GEO优化指南",
		Body:         "这是一篇关于GEO优化的文章",
		ContentType:  "article",
		SchemaMarkup: `{"@type":"Article"}`,
		Status:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if c.ID != 1 {
		t.Errorf("expected ID=1, got %d", c.ID)
	}
	if c.UserID != 100 {
		t.Errorf("expected UserID=100, got %d", c.UserID)
	}
	if c.Title != "GEO优化指南" {
		t.Errorf("expected Title='GEO优化指南', got '%s'", c.Title)
	}
	if c.ContentType != "article" {
		t.Errorf("expected ContentType='article', got '%s'", c.ContentType)
	}
	if c.Status != 1 {
		t.Errorf("expected Status=1, got %d", c.Status)
	}
}

func TestContent_DefaultValues(t *testing.T) {
	c := Content{}
	if c.Status != 0 {
		t.Errorf("expected default Status=0, got %d", c.Status)
	}
	if c.ID != 0 {
		t.Errorf("expected default ID=0, got %d", c.ID)
	}
}

func BenchmarkContentCreation(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		_ = Content{
			ID:          int64(i),
			UserID:      100,
			Title:       "测试标题",
			Body:        "测试内容",
			ContentType: "article",
			Status:      0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}
}
