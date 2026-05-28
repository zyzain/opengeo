package knowledge

import (
	"testing"
	"time"
)

func TestEntityToMap(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	e := &Entity{
		ID:             1,
		UserID:         100,
		EntityName:     "OpenGEO",
		EntityType:     "brand",
		EntityData:     `{"description":"test"}`,
		AuthorityLinks: `["https://example.com"]`,
		ContentCount:   5,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	m := entityToMap(e)

	if m["id"] != int64(1) {
		t.Errorf("expected id=1, got %v", m["id"])
	}
	if m["user_id"] != int64(100) {
		t.Errorf("expected user_id=100, got %v", m["user_id"])
	}
	if m["entity_name"] != "OpenGEO" {
		t.Errorf("expected entity_name=OpenGEO, got %v", m["entity_name"])
	}
	if m["entity_type"] != "brand" {
		t.Errorf("expected entity_type=brand, got %v", m["entity_type"])
	}
	if m["content_count"] != int32(5) {
		t.Errorf("expected content_count=5, got %v (type %T)", m["content_count"], m["content_count"])
	}
}

func TestEntity_Fields(t *testing.T) {
	e := Entity{
		ID:             10,
		UserID:         200,
		EntityName:     "DeepSeek",
		EntityType:     "product",
		EntityData:     `{"company":"深度求索"}`,
		AuthorityLinks: `[]`,
		ContentCount:   3,
	}

	if e.ID != 10 {
		t.Errorf("expected ID=10, got %d", e.ID)
	}
	if e.EntityName != "DeepSeek" {
		t.Errorf("expected EntityName=DeepSeek, got %s", e.EntityName)
	}
	if e.EntityType != "product" {
		t.Errorf("expected EntityType=product, got %s", e.EntityType)
	}
}

func TestEntity_TableName(t *testing.T) {
	e := Entity{}
	if e.TableName() != "entities" {
		t.Errorf("expected table name 'entities', got '%s'", e.TableName())
	}
}

func BenchmarkEntityToMap(b *testing.B) {
	now := time.Now()
	e := &Entity{
		ID:             1,
		UserID:         100,
		EntityName:     "OpenGEO",
		EntityType:     "brand",
		EntityData:     `{"description":"test"}`,
		AuthorityLinks: `["https://example.com"]`,
		ContentCount:   5,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityToMap(e)
	}
}
