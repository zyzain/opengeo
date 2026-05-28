package client

import (
	"context"

	"opengeo/gateway/internal/knowledge"
)

type KnowledgeClient struct {
	svc *knowledge.Service
}

func NewKnowledgeClient(svc *knowledge.Service) *KnowledgeClient {
	return &KnowledgeClient{svc: svc}
}

func (c *KnowledgeClient) CreateEntity(ctx context.Context, userID int64, entityName, entityType, entityData, authorityLinks string) (map[string]interface{}, error) {
	return c.svc.Create(ctx, userID, entityName, entityType, entityData, authorityLinks)
}

func (c *KnowledgeClient) GetEntity(ctx context.Context, id int64) (map[string]interface{}, error) {
	return c.svc.Get(ctx, id)
}

func (c *KnowledgeClient) UpdateEntity(ctx context.Context, id int64, entityName, entityType, entityData, authorityLinks string) (map[string]interface{}, error) {
	return c.svc.Update(ctx, id, entityName, entityType, entityData, authorityLinks)
}

func (c *KnowledgeClient) DeleteEntity(ctx context.Context, id int64) error {
	return c.svc.Delete(ctx, id)
}

func (c *KnowledgeClient) ListEntities(ctx context.Context, userID int64, entityType string, page, pageSize int) (map[string]interface{}, error) {
	return c.svc.List(ctx, userID, entityType, page, pageSize)
}

func (c *KnowledgeClient) SearchEntities(ctx context.Context, userID int64, keyword string) (map[string]interface{}, error) {
	return c.svc.Search(ctx, userID, keyword)
}
