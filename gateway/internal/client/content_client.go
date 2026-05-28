package client

import (
	"context"

	"opengeo/gateway/internal/content"
)

type ContentClient struct {
	svc *content.Service
}

func NewContentClient(svc *content.Service) *ContentClient {
	return &ContentClient{svc: svc}
}

func (c *ContentClient) CreateContent(ctx context.Context, userID int64, title, body, contentType, schemaMarkup string) (map[string]interface{}, error) {
	return c.svc.Create(ctx, userID, title, body, contentType, schemaMarkup)
}

func (c *ContentClient) GetContent(ctx context.Context, id int64) (map[string]interface{}, error) {
	return c.svc.Get(ctx, id)
}

func (c *ContentClient) UpdateContent(ctx context.Context, id int64, title, body, schemaMarkup string) (map[string]interface{}, error) {
	return c.svc.Update(ctx, id, title, body, schemaMarkup)
}

func (c *ContentClient) DeleteContent(ctx context.Context, id int64) error {
	return c.svc.Delete(ctx, id)
}

func (c *ContentClient) ListContents(ctx context.Context, userID int64, page, pageSize, status int, contentType string) (map[string]interface{}, error) {
	return c.svc.List(ctx, userID, page, pageSize, status, contentType)
}

func (c *ContentClient) OptimizeContent(ctx context.Context, contentID int64, optimizationType string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"success":           true,
		"optimized_content": "优化后的内容",
		"score":             85.5,
		"details":           "内容结构良好",
	}, nil
}

func (c *ContentClient) PublishContent(ctx context.Context, contentID int64, channelIDs []int64) (map[string]interface{}, error) {
	return map[string]interface{}{
		"success":  true,
		"task_ids": []int64{1, 2, 3},
		"message":  "发布任务已创建",
	}, nil
}

func (c *ContentClient) CheckCompliance(ctx context.Context, contentID int64) (map[string]interface{}, error) {
	return map[string]interface{}{
		"content_id": contentID,
		"compliant":  true,
		"score":      95,
		"issues":     []interface{}{},
		"report":     "内容合规",
	}, nil
}
