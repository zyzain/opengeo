package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/service/publish/internal/domain/model"
)

// PlatformRepository 平台仓储
type PlatformRepository struct {
	db *gorm.DB
}

// NewPlatformRepository 创建平台仓储
func NewPlatformRepository(db *gorm.DB) *PlatformRepository {
	return &PlatformRepository{db: db}
}

// Create 创建平台
func (r *PlatformRepository) Create(ctx context.Context, platform *model.Platform) error {
	if err := r.db.WithContext(ctx).Create(platform).Error; err != nil {
		return fmt.Errorf("failed to create platform: %w", err)
	}
	return nil
}

// GetByID 根据ID获取平台
func (r *PlatformRepository) GetByID(ctx context.Context, id int64) (*model.Platform, error) {
	var platform model.Platform
	if err := r.db.WithContext(ctx).First(&platform, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("platform not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get platform: %w", err)
	}
	return &platform, nil
}

// GetByCode 根据代码获取平台
func (r *PlatformRepository) GetByCode(ctx context.Context, code string) (*model.Platform, error) {
	var platform model.Platform
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&platform).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("platform not found: %s", code)
		}
		return nil, fmt.Errorf("failed to get platform: %w", err)
	}
	return &platform, nil
}

// Update 更新平台
func (r *PlatformRepository) Update(ctx context.Context, platform *model.Platform) error {
	platform.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(platform).Error; err != nil {
		return fmt.Errorf("failed to update platform: %w", err)
	}
	return nil
}

// Delete 删除平台
func (r *PlatformRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Platform{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete platform: %w", err)
	}
	return nil
}

// List 列出平台
func (r *PlatformRepository) List(ctx context.Context, isEnabled *bool, page, pageSize int) ([]*model.Platform, int32, error) {
	var platforms []*model.Platform
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Platform{})

	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count platforms: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort_order DESC, id ASC").Find(&platforms).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list platforms: %w", err)
	}

	return platforms, int32(total), nil
}

// ListAll 列出所有启用的平台
func (r *PlatformRepository) ListAll(ctx context.Context) ([]*model.Platform, error) {
	var platforms []*model.Platform
	if err := r.db.WithContext(ctx).Where("is_enabled = ?", true).Order("sort_order DESC, id ASC").Find(&platforms).Error; err != nil {
		return nil, fmt.Errorf("failed to list platforms: %w", err)
	}
	return platforms, nil
}

// UpdateStatus 更新平台状态
func (r *PlatformRepository) UpdateStatus(ctx context.Context, id int64, isEnabled bool) error {
	if err := r.db.WithContext(ctx).Model(&model.Platform{}).Where("id = ?", id).Update("is_enabled", isEnabled).Error; err != nil {
		return fmt.Errorf("failed to update platform status: %w", err)
	}
	return nil
}

// Count 统计平台数量
func (r *PlatformRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Platform{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count platforms: %w", err)
	}
	return count, nil
}
