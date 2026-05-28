package client

import (
	"context"
)

type TemplateClient struct {
	svc interface {
		Create(ctx context.Context, userID int64, name, description, templateType, content, tags string, isPublic bool) (map[string]interface{}, error)
		Get(ctx context.Context, id int64) (map[string]interface{}, error)
		Update(ctx context.Context, id int64, name, description, templateType, content, tags string) (map[string]interface{}, error)
		Delete(ctx context.Context, id int64) error
		List(ctx context.Context, templateType string, page, pageSize int) (map[string]interface{}, error)
	}
}

func NewTemplateClient(svc interface {
	Create(ctx context.Context, userID int64, name, description, templateType, content, tags string, isPublic bool) (map[string]interface{}, error)
	Get(ctx context.Context, id int64) (map[string]interface{}, error)
	Update(ctx context.Context, id int64, name, description, templateType, content, tags string) (map[string]interface{}, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, templateType string, page, pageSize int) (map[string]interface{}, error)
}) *TemplateClient {
	return &TemplateClient{svc: svc}
}

func (c *TemplateClient) Create(ctx context.Context, userID int64, name, description, templateType, content, tags string, isPublic bool) (map[string]interface{}, error) {
	return c.svc.Create(ctx, userID, name, description, templateType, content, tags, isPublic)
}

func (c *TemplateClient) Get(ctx context.Context, id int64) (map[string]interface{}, error) {
	return c.svc.Get(ctx, id)
}

func (c *TemplateClient) Update(ctx context.Context, id int64, name, description, templateType, content, tags string) (map[string]interface{}, error) {
	return c.svc.Update(ctx, id, name, description, templateType, content, tags)
}

func (c *TemplateClient) Delete(ctx context.Context, id int64) error {
	return c.svc.Delete(ctx, id)
}

func (c *TemplateClient) List(ctx context.Context, templateType string, page, pageSize int) (map[string]interface{}, error) {
	return c.svc.List(ctx, templateType, page, pageSize)
}
