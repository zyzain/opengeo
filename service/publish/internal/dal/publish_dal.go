package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/service/publish/internal/domain/model"
)

// PublishTaskRepository 发布任务仓储
type PublishTaskRepository struct {
	db *gorm.DB
}

// NewPublishTaskRepository 创建发布任务仓储
func NewPublishTaskRepository(db *gorm.DB) *PublishTaskRepository {
	return &PublishTaskRepository{db: db}
}

// Create 创建发布任务
func (r *PublishTaskRepository) Create(ctx context.Context, task *model.PublishTask) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return fmt.Errorf("failed to create publish task: %w", err)
	}
	return nil
}

// GetByID 根据ID获取发布任务
func (r *PublishTaskRepository) GetByID(ctx context.Context, id int64) (*model.PublishTask, error) {
	var task model.PublishTask
	if err := r.db.WithContext(ctx).First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("publish task not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get publish task: %w", err)
	}
	return &task, nil
}

// Update 更新发布任务
func (r *PublishTaskRepository) Update(ctx context.Context, task *model.PublishTask) error {
	if err := r.db.WithContext(ctx).Save(task).Error; err != nil {
		return fmt.Errorf("failed to update publish task: %w", err)
	}
	return nil
}

// Delete 删除发布任务
func (r *PublishTaskRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.PublishTask{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete publish task: %w", err)
	}
	return nil
}

// List 列出发布任务
func (r *PublishTaskRepository) List(ctx context.Context, userID int64, status model.PublishStatus, page, pageSize int) ([]*model.PublishTask, int32, error) {
	var tasks []*model.PublishTask
	var total int64

	query := r.db.WithContext(ctx).Model(&model.PublishTask{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count publish tasks: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list publish tasks: %w", err)
	}

	return tasks, int32(total), nil
}

// GetPendingTasks 获取待处理任务
func (r *PublishTaskRepository) GetPendingTasks(ctx context.Context, limit int) ([]*model.PublishTask, error) {
	var tasks []*model.PublishTask

	err := r.db.WithContext(ctx).
		Where("status = ? AND (scheduled_time IS NULL OR scheduled_time <= ?)", model.PublishStatusPending, time.Now()).
		Limit(limit).
		Order("created_at ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending tasks: %w", err)
	}

	return tasks, nil
}

// ChannelRepository 渠道仓储
type ChannelRepository struct {
	db *gorm.DB
}

// NewChannelRepository 创建渠道仓储
func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

// Create 创建渠道
func (r *ChannelRepository) Create(ctx context.Context, channel *model.Channel) error {
	if err := r.db.WithContext(ctx).Create(channel).Error; err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	return nil
}

// GetByID 根据ID获取渠道
func (r *ChannelRepository) GetByID(ctx context.Context, id int64) (*model.Channel, error) {
	var channel model.Channel
	if err := r.db.WithContext(ctx).First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("channel not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	return &channel, nil
}

// Update 更新渠道
func (r *ChannelRepository) Update(ctx context.Context, channel *model.Channel) error {
	if err := r.db.WithContext(ctx).Save(channel).Error; err != nil {
		return fmt.Errorf("failed to update channel: %w", err)
	}
	return nil
}

// Delete 删除渠道
func (r *ChannelRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Channel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}
	return nil
}

// ListByUser 列出用户的渠道
func (r *ChannelRepository) ListByUser(ctx context.Context, userID int64) ([]*model.Channel, error) {
	var channels []*model.Channel

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(100).
		Find(&channels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	return channels, nil
}

// ListEnabled 列出启用的渠道
func (r *ChannelRepository) ListEnabled(ctx context.Context, userID int64) ([]*model.Channel, error) {
	var channels []*model.Channel

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_enabled = ?", userID, true).
		Order("created_at DESC").
		Limit(100).
		Find(&channels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	return channels, nil
}