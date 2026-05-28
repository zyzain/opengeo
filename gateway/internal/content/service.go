package content

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

func (s *Service) Create(ctx context.Context, userID int64, title, body, contentType, schemaMarkup string) (map[string]interface{}, error) {
	c := &Content{
		UserID:       userID,
		Title:        title,
		Body:         body,
		ContentType:  contentType,
		SchemaMarkup: schemaMarkup,
		Status:       0,
	}
	if err := s.store.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create content failed: %w", err)
	}
	return map[string]interface{}{
		"id":           c.ID,
		"user_id":      c.UserID,
		"title":        c.Title,
		"body":         c.Body,
		"content_type": c.ContentType,
		"status":       c.Status,
		"created_at":   c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) Get(ctx context.Context, id int64) (map[string]interface{}, error) {
	c, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}
	return map[string]interface{}{
		"id":            c.ID,
		"user_id":       c.UserID,
		"title":         c.Title,
		"body":          c.Body,
		"content_type":  c.ContentType,
		"schema_markup": c.SchemaMarkup,
		"status":        c.Status,
		"created_at":    c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":    c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) Update(ctx context.Context, id int64, title, body, schemaMarkup string) (map[string]interface{}, error) {
	c, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}
	c.Title = title
	c.Body = body
	c.SchemaMarkup = schemaMarkup
	if err := s.store.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("update content failed: %w", err)
	}
	return map[string]interface{}{
		"id":         c.ID,
		"title":      c.Title,
		"updated_at": c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, userID int64, page, pageSize, status int, contentType string) (map[string]interface{}, error) {
	items, total, err := s.store.List(ctx, userID, status, contentType, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list contents failed: %w", err)
	}

	type contentItem struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"user_id"`
		Title       string `json:"title"`
		Body        string `json:"body"`
		ContentType string `json:"content_type"`
		Status      int32  `json:"status"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	var list []contentItem
	for _, c := range items {
		list = append(list, contentItem{
			ID:          c.ID,
			UserID:      c.UserID,
			Title:       c.Title,
			Body:        c.Body,
			ContentType: c.ContentType,
			Status:      c.Status,
			CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return map[string]interface{}{
		"items": list,
		"total": total,
	}, nil
}
