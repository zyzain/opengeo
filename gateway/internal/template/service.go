package template

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

func (s *Service) Create(ctx context.Context, userID int64, name, description, templateType, content, tags string, isPublic bool) (map[string]interface{}, error) {
	t := &Template{
		UserID:       userID,
		Name:         name,
		Description:  description,
		TemplateType: templateType,
		Content:      content,
		Tags:         tags,
		IsPublic:     isPublic,
	}
	if err := s.store.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create template failed: %w", err)
	}
	return map[string]interface{}{
		"id":            t.ID,
		"name":          t.Name,
		"description":   t.Description,
		"template_type": t.TemplateType,
		"content":       t.Content,
		"tags":          t.Tags,
		"is_public":     t.IsPublic,
		"usage_count":   t.UsageCount,
		"rating":        t.Rating,
		"created_at":    t.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) Get(ctx context.Context, id int64) (map[string]interface{}, error) {
	t, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}
	return map[string]interface{}{
		"id":            t.ID,
		"name":          t.Name,
		"description":   t.Description,
		"template_type": t.TemplateType,
		"content":       t.Content,
		"tags":          t.Tags,
		"is_public":     t.IsPublic,
		"usage_count":   t.UsageCount,
		"rating":        t.Rating,
		"created_at":    t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":    t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) Update(ctx context.Context, id int64, name, description, templateType, content, tags string) (map[string]interface{}, error) {
	t, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}
	if name != "" {
		t.Name = name
	}
	if description != "" {
		t.Description = description
	}
	if templateType != "" {
		t.TemplateType = templateType
	}
	if content != "" {
		t.Content = content
	}
	if tags != "" {
		t.Tags = tags
	}
	if err := s.store.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("update template failed: %w", err)
	}
	return map[string]interface{}{
		"id":         t.ID,
		"name":       t.Name,
		"updated_at": t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, templateType string, page, pageSize int) (map[string]interface{}, error) {
	items, total, err := s.store.List(ctx, templateType, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list templates failed: %w", err)
	}

	var list []map[string]interface{}
	for _, t := range items {
		list = append(list, map[string]interface{}{
			"id":            t.ID,
			"name":          t.Name,
			"description":   t.Description,
			"template_type": t.TemplateType,
			"content":       t.Content,
			"tags":          t.Tags,
			"is_public":     t.IsPublic,
			"usage_count":   t.UsageCount,
			"rating":        t.Rating,
			"created_at":    t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return map[string]interface{}{
		"items": list,
		"total": total,
	}, nil
}

func (s *Service) IncrementUsage(ctx context.Context, id int64) error {
	return s.store.IncrementUsage(ctx, id)
}
