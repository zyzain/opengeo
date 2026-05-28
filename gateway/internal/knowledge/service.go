package knowledge

import (
	"context"
	"fmt"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, userID int64, entityName, entityType, entityData, authorityLinks string) (map[string]interface{}, error) {
	e := &Entity{
		UserID:         userID,
		EntityName:     entityName,
		EntityType:     entityType,
		EntityData:     entityData,
		AuthorityLinks: authorityLinks,
	}
	if err := s.store.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("create entity failed: %w", err)
	}
	return entityToMap(e), nil
}

func (s *Service) Get(ctx context.Context, id int64) (map[string]interface{}, error) {
	e, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("entity not found: %w", err)
	}
	return entityToMap(e), nil
}

func (s *Service) Update(ctx context.Context, id int64, entityName, entityType, entityData, authorityLinks string) (map[string]interface{}, error) {
	e, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("entity not found: %w", err)
	}
	e.EntityName = entityName
	e.EntityType = entityType
	e.EntityData = entityData
	e.AuthorityLinks = authorityLinks
	if err := s.store.Update(ctx, e); err != nil {
		return nil, fmt.Errorf("update entity failed: %w", err)
	}
	return entityToMap(e), nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, userID int64, entityType string, page, pageSize int) (map[string]interface{}, error) {
	items, total, err := s.store.List(ctx, userID, entityType, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list entities failed: %w", err)
	}

	var list []map[string]interface{}
	for _, e := range items {
		list = append(list, entityToMap(&e))
	}

	return map[string]interface{}{
		"items": list,
		"total": total,
	}, nil
}

func (s *Service) Search(ctx context.Context, userID int64, keyword string) (map[string]interface{}, error) {
	items, err := s.store.SearchByName(ctx, userID, keyword, 20)
	if err != nil {
		return nil, fmt.Errorf("search entities failed: %w", err)
	}

	var list []map[string]interface{}
	for _, e := range items {
		list = append(list, entityToMap(&e))
	}

	return map[string]interface{}{
		"items": list,
	}, nil
}

func entityToMap(e *Entity) map[string]interface{} {
	return map[string]interface{}{
		"id":              e.ID,
		"user_id":         e.UserID,
		"entity_name":     e.EntityName,
		"entity_type":     e.EntityType,
		"entity_data":     e.EntityData,
		"authority_links": e.AuthorityLinks,
		"content_count":   e.ContentCount,
		"created_at":      e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":      e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
