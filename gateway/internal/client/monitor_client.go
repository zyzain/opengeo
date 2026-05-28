package client

import (
	"context"

	"gorm.io/gorm"
)

type MonitorClient struct {
	db *gorm.DB
}

func NewMonitorClient(db *gorm.DB) *MonitorClient {
	return &MonitorClient{db: db}
}

func (c *MonitorClient) GetAICitations(ctx context.Context, contentID int64, aiModel string, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("ai_citations")
	if contentID > 0 {
		query = query.Where("content_id = ?", contentID)
	}
	if aiModel != "" {
		query = query.Where("ai_model = ?", aiModel)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("tracked_at DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *MonitorClient) GetSourceScores(ctx context.Context, channelID, accountID int64) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("source_scores")
	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if accountID > 0 {
		query = query.Where("account_id = ?", accountID)
	}

	query.Count(&total)
	query.Order("score DESC").Find(&items)

	return map[string]interface{}{
		"items": items,
		"total": total,
	}, nil
}

func (c *MonitorClient) GetCompetitorMonitors(ctx context.Context, userID int64, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("competitor_monitors")
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

func (c *MonitorClient) GetROIMetrics(ctx context.Context, contentID, channelID int64, startDate, endDate string) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("roi_metrics")
	if contentID > 0 {
		query = query.Where("content_id = ?", contentID)
	}
	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if startDate != "" {
		query = query.Where("recorded_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("recorded_at <= ?", endDate)
	}

	query.Count(&total)
	query.Order("recorded_at DESC").Find(&items)

	return map[string]interface{}{
		"items": items,
		"total": total,
	}, nil
}

func (c *MonitorClient) GenerateSuggestions(ctx context.Context, contentID int64) (map[string]interface{}, error) {
	var items []map[string]interface{}
	c.db.WithContext(ctx).Table("optimization_suggestions").Where("content_id = ?", contentID).Order("priority DESC").Find(&items)

	return map[string]interface{}{
		"items": items,
		"total": len(items),
	}, nil
}
