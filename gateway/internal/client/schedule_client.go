package client

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ScheduleClient struct {
	db *gorm.DB
}

func NewScheduleClient(db *gorm.DB) *ScheduleClient {
	return &ScheduleClient{db: db}
}

func (c *ScheduleClient) CreateSchedule(ctx context.Context, userID int64, scheduleName, scheduleType, cronExpression, config string) (map[string]interface{}, error) {
	schedule := map[string]interface{}{
		"user_id":         userID,
		"schedule_name":   scheduleName,
		"schedule_type":   scheduleType,
		"cron_expression": cronExpression,
		"config":          config,
		"is_enabled":      true,
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}
	if err := c.db.WithContext(ctx).Table("schedules").Create(schedule).Error; err != nil {
		return nil, fmt.Errorf("create schedule: %w", err)
	}
	return schedule, nil
}

func (c *ScheduleClient) GetSchedule(ctx context.Context, id int64) (map[string]interface{}, error) {
	var schedule map[string]interface{}
	if err := c.db.WithContext(ctx).Table("schedules").Where("id = ?", id).First(&schedule).Error; err != nil {
		return nil, fmt.Errorf("schedule not found")
	}
	return schedule, nil
}

func (c *ScheduleClient) UpdateSchedule(ctx context.Context, id int64, scheduleName, cronExpression, config string) (map[string]interface{}, error) {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if scheduleName != "" {
		updates["schedule_name"] = scheduleName
	}
	if cronExpression != "" {
		updates["cron_expression"] = cronExpression
	}
	if config != "" {
		updates["config"] = config
	}
	c.db.WithContext(ctx).Table("schedules").Where("id = ?", id).Updates(updates)
	var schedule map[string]interface{}
	c.db.WithContext(ctx).Table("schedules").Where("id = ?", id).First(&schedule)
	return schedule, nil
}

func (c *ScheduleClient) DeleteSchedule(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("schedules").Where("id = ?", id).Delete(nil).Error
}

func (c *ScheduleClient) ListSchedules(ctx context.Context, userID int64, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("schedules")
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

func (c *ScheduleClient) EnableSchedule(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("schedules").Where("id = ?", id).Update("is_enabled", true).Error
}

func (c *ScheduleClient) DisableSchedule(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("schedules").Where("id = ?", id).Update("is_enabled", false).Error
}
