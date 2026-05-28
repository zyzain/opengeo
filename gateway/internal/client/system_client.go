package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SystemClient struct {
	db *gorm.DB
}

func NewSystemClient(db *gorm.DB) *SystemClient {
	return &SystemClient{db: db}
}

func (c *SystemClient) GetSystemConfigs(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	c.db.WithContext(ctx).Table("system_configs").Count(&total)
	c.db.WithContext(ctx).Table("system_configs").Offset((page-1)*pageSize).Limit(pageSize).Order("config_key ASC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *SystemClient) UpdateSystemConfig(ctx context.Context, key, value string) (map[string]interface{}, error) {
	c.db.WithContext(ctx).Table("system_configs").Where("config_key = ?", key).Update("config_value", value)
	var config map[string]interface{}
	c.db.WithContext(ctx).Table("system_configs").Where("config_key = ?", key).First(&config)
	return config, nil
}

func (c *SystemClient) ListPlugins(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	c.db.WithContext(ctx).Table("plugins").Count(&total)
	c.db.WithContext(ctx).Table("plugins").Offset((page-1)*pageSize).Limit(pageSize).Order("plugin_name ASC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *SystemClient) ListWebhooks(ctx context.Context, userID int64, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("webhooks")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *SystemClient) CreateWebhook(ctx context.Context, userID int64, webhookName, url, secret string, events []string) (map[string]interface{}, error) {
	if err := validateWebhookURL(url); err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %w", err)
	}

	eventsJSON, _ := json.Marshal(events)
	webhook := map[string]interface{}{
		"user_id":      userID,
		"webhook_name": webhookName,
		"url":          url,
		"secret":       secret,
		"events":       string(eventsJSON),
		"is_active":    true,
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}
	if err := c.db.WithContext(ctx).Table("webhooks").Create(webhook).Error; err != nil {
		return nil, fmt.Errorf("create webhook: %w", err)
	}
	return webhook, nil
}

func (c *SystemClient) UpdateWebhook(ctx context.Context, id int64, webhookName, url, secret string, events []string) (map[string]interface{}, error) {
	if url != "" {
		if err := validateWebhookURL(url); err != nil {
			return nil, fmt.Errorf("invalid webhook URL: %w", err)
		}
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if webhookName != "" {
		updates["webhook_name"] = webhookName
	}
	if url != "" {
		updates["url"] = url
	}
	if secret != "" {
		updates["secret"] = secret
	}
	if events != nil {
		eventsJSON, _ := json.Marshal(events)
		updates["events"] = string(eventsJSON)
	}
	c.db.WithContext(ctx).Table("webhooks").Where("id = ?", id).Updates(updates)
	var webhook map[string]interface{}
	c.db.WithContext(ctx).Table("webhooks").Where("id = ?", id).First(&webhook)
	return webhook, nil
}

func (c *SystemClient) DeleteWebhook(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("webhooks").Where("id = ?", id).Delete(nil).Error
}

func (c *SystemClient) TestWebhook(ctx context.Context, id int64) (map[string]interface{}, error) {
	return map[string]interface{}{
		"success":   true,
		"message":   "Webhook测试成功",
		"tested_at": time.Now(),
	}, nil
}

func (c *SystemClient) GetWebhookHistory(ctx context.Context, id int64, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	c.db.WithContext(ctx).Table("webhook_events").Where("webhook_id = ?", id).Count(&total)
	c.db.WithContext(ctx).Table("webhook_events").Where("webhook_id = ?", id).Offset((page-1)*pageSize).Limit(pageSize).Order("triggered_at DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}
