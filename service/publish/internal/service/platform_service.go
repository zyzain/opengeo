package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"opengeo/service/publish/internal/dal"
	"opengeo/service/publish/internal/domain/model"
)

// PlatformService 平台管理服务
type PlatformService struct {
	platformRepo *dal.PlatformRepository
}

// NewPlatformService 创建平台管理服务
func NewPlatformService(platformRepo *dal.PlatformRepository) *PlatformService {
	return &PlatformService{platformRepo: platformRepo}
}

// CreatePlatform 创建平台
func (s *PlatformService) CreatePlatform(ctx context.Context, req *CreatePlatformRequest) (*model.Platform, error) {
	// 检查代码是否重复
	existing, _ := s.platformRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("platform code already exists: %s", req.Code)
	}

	// 序列化配置模板
	configSchema := ""
	if req.ConfigTemplate != nil {
		data, _ := json.Marshal(req.ConfigTemplate)
		configSchema = string(data)
	}

	platform := &model.Platform{
		Code:         req.Code,
		Name:         req.Name,
		Icon:         req.Icon,
		Color:        req.Color,
		Description:  req.Description,
		ConfigSchema: configSchema,
		IsEnabled:    true,
		SortOrder:    req.SortOrder,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.platformRepo.Create(ctx, platform); err != nil {
		return nil, fmt.Errorf("failed to create platform: %w", err)
	}

	return platform, nil
}

// GetPlatform 获取平台
func (s *PlatformService) GetPlatform(ctx context.Context, id int64) (*model.Platform, error) {
	return s.platformRepo.GetByID(ctx, id)
}

// UpdatePlatform 更新平台
func (s *PlatformService) UpdatePlatform(ctx context.Context, id int64, req *UpdatePlatformRequest) (*model.Platform, error) {
	platform, err := s.platformRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		platform.Name = req.Name
	}
	if req.Icon != "" {
		platform.Icon = req.Icon
	}
	if req.Color != "" {
		platform.Color = req.Color
	}
	if req.Description != "" {
		platform.Description = req.Description
	}
	if req.ConfigTemplate != nil {
		data, _ := json.Marshal(req.ConfigTemplate)
		platform.ConfigSchema = string(data)
	}
	if req.SortOrder != 0 {
		platform.SortOrder = req.SortOrder
	}
	platform.UpdatedAt = time.Now()

	if err := s.platformRepo.Update(ctx, platform); err != nil {
		return nil, fmt.Errorf("failed to update platform: %w", err)
	}

	return platform, nil
}

// DeletePlatform 删除平台
func (s *PlatformService) DeletePlatform(ctx context.Context, id int64) error {
	return s.platformRepo.Delete(ctx, id)
}

// ListPlatforms 列出平台
func (s *PlatformService) ListPlatforms(ctx context.Context, isEnabled *bool, page, pageSize int) ([]*model.Platform, int32, error) {
	return s.platformRepo.List(ctx, isEnabled, page, pageSize)
}

// ListAllPlatforms 列出所有启用平台
func (s *PlatformService) ListAllPlatforms(ctx context.Context) ([]*model.Platform, error) {
	return s.platformRepo.ListAll(ctx)
}

// EnablePlatform 启用平台
func (s *PlatformService) EnablePlatform(ctx context.Context, id int64) error {
	return s.platformRepo.UpdateStatus(ctx, id, true)
}

// DisablePlatform 禁用平台
func (s *PlatformService) DisablePlatform(ctx context.Context, id int64) error {
	return s.platformRepo.UpdateStatus(ctx, id, false)
}

// SeedDefaultPlatforms 初始化默认平台
func (s *PlatformService) SeedDefaultPlatforms(ctx context.Context) error {
	count, err := s.platformRepo.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	for _, p := range model.DefaultPlatforms {
		configSchema := ""
		if p.ConfigTemplate != nil {
			data, _ := json.Marshal(p.ConfigTemplate)
			configSchema = string(data)
		}

		platform := &model.Platform{
			Code:         p.Code,
			Name:         p.Name,
			Icon:         p.Icon,
			Color:        p.Color,
			Description:  p.Description,
			ConfigSchema: configSchema,
			IsEnabled:    p.IsEnabled,
			SortOrder:    p.SortOrder,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.platformRepo.Create(ctx, platform); err != nil {
			continue
		}
	}

	return nil
}

// ==================== 请求/响应模型 ====================

type CreatePlatformRequest struct {
	Code           string                      `json:"code"`
	Name           string                      `json:"name"`
	Icon           string                      `json:"icon"`
	Color          string                      `json:"color"`
	Description    string                      `json:"description"`
	ConfigTemplate []model.PlatformConfigTemplate `json:"config_template"`
	SortOrder      int32                       `json:"sort_order"`
}

type UpdatePlatformRequest struct {
	Name           string                      `json:"name"`
	Icon           string                      `json:"icon"`
	Color          string                      `json:"color"`
	Description    string                      `json:"description"`
	ConfigTemplate []model.PlatformConfigTemplate `json:"config_template"`
	SortOrder      int32                       `json:"sort_order"`
}
