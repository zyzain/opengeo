package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"opengeo/service/system/internal/dal"
	"opengeo/service/system/internal/domain/model"
)

// DeploymentMode 部署模式
type DeploymentMode string

const (
	DeploymentPrivate  DeploymentMode = "private"  // 私有化单机部署
	DeploymentSaaS     DeploymentMode = "saas"     // SaaS多租户云端
)

// TenantService 多租户服务
type TenantService struct {
	tenantRepo *dal.TenantRepository
	configRepo *dal.SystemConfigRepository
}

// NewTenantService 创建多租户服务
func NewTenantService(tenantRepo *dal.TenantRepository, configRepo *dal.SystemConfigRepository) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		configRepo: configRepo,
	}
}

// ==================== 部署模式管理 ====================

// GetDeploymentMode 获取当前部署模式
func (s *TenantService) GetDeploymentMode(ctx context.Context) (DeploymentMode, error) {
	config, err := s.configRepo.GetByKey(ctx, "deployment_mode")
	if err != nil {
		// 默认私有化部署
		return DeploymentPrivate, nil
	}

	switch config.ConfigValue {
	case "saas":
		return DeploymentSaaS, nil
	default:
		return DeploymentPrivate, nil
	}
}

// SetDeploymentMode 设置部署模式
func (s *TenantService) SetDeploymentMode(ctx context.Context, mode DeploymentMode) error {
	return s.configRepo.Set(ctx, "deployment_mode", string(mode), "string", "系统部署模式: private/saas", true)
}

// ==================== 租户管理 ====================

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*model.Tenant, error) {
	// 检查部署模式
	mode, _ := s.GetDeploymentMode(ctx)
	if mode == DeploymentPrivate {
		// 私有化模式下只允许一个默认租户
		existing, _, _ := s.tenantRepo.List(ctx, 1, 1)
		if len(existing) > 0 {
			return nil, fmt.Errorf("private mode only allows one tenant")
		}
	}

	// 构建配额配置
	quotaConfig := &TenantQuota{
		MaxUsers:       100,
		MaxAccounts:    50,
		MaxContents:    1000,
		MaxChannels:    10,
		MaxPublishPerDay: 500,
		MaxStorageMB:   1024,
		EnabledFeatures: []string{"content", "publish", "monitor"},
	}
	if req.QuotaConfig != nil {
		quotaConfig = req.QuotaConfig
	}

	quotaJSON, _ := json.Marshal(quotaConfig)

	tenant := &model.Tenant{
		TenantName:   req.Name,
		TenantCode:   req.Code,
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
		Status:       1,
		QuotaConfig:  string(quotaJSON),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

// GetTenant 获取租户
func (s *TenantService) GetTenant(ctx context.Context, tenantID int64) (*model.Tenant, error) {
	return s.tenantRepo.GetByID(ctx, tenantID)
}

// ListTenants 列出租户
func (s *TenantService) ListTenants(ctx context.Context, page, pageSize int) ([]*model.Tenant, int32, error) {
	return s.tenantRepo.List(ctx, page, pageSize)
}

// UpdateTenant 更新租户
func (s *TenantService) UpdateTenant(ctx context.Context, tenantID int64, req *UpdateTenantRequest) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	if req.ContactName != "" {
		tenant.ContactName = req.ContactName
	}
	if req.ContactEmail != "" {
		tenant.ContactEmail = req.ContactEmail
	}
	if req.QuotaConfig != nil {
		quotaJSON, _ := json.Marshal(req.QuotaConfig)
		tenant.QuotaConfig = string(quotaJSON)
	}
	tenant.UpdatedAt = time.Now()

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return tenant, nil
}

// DeleteTenant 删除租户
func (s *TenantService) DeleteTenant(ctx context.Context, tenantID int64) error {
	return s.tenantRepo.Delete(ctx, tenantID)
}

// EnableTenant 启用租户
func (s *TenantService) EnableTenant(ctx context.Context, tenantID int64) error {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	tenant.Status = 1
	tenant.UpdatedAt = time.Now()
	return s.tenantRepo.Update(ctx, tenant)
}

// DisableTenant 禁用租户
func (s *TenantService) DisableTenant(ctx context.Context, tenantID int64) error {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	tenant.Status = 0
	tenant.UpdatedAt = time.Now()
	return s.tenantRepo.Update(ctx, tenant)
}

// ==================== 配额管理 ====================

// TenantQuota 租户配额
type TenantQuota struct {
	MaxUsers         int      `json:"max_users"`
	MaxAccounts      int      `json:"max_accounts"`
	MaxContents      int      `json:"max_contents"`
	MaxChannels      int      `json:"max_channels"`
	MaxPublishPerDay int      `json:"max_publish_per_day"`
	MaxStorageMB     int      `json:"max_storage_mb"`
	EnabledFeatures  []string `json:"enabled_features"`
	RateLimit        int      `json:"rate_limit"` // API限流 QPS
}

// GetTenantQuota 获取租户配额
func (s *TenantService) GetTenantQuota(ctx context.Context, tenantID int64) (*TenantQuota, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	quota := &TenantQuota{
		MaxUsers:         100,
		MaxAccounts:      50,
		MaxContents:      1000,
		MaxChannels:      10,
		MaxPublishPerDay: 500,
		MaxStorageMB:     1024,
		EnabledFeatures:  []string{"content", "publish", "monitor"},
		RateLimit:        100,
	}

	if tenant.QuotaConfig != "" {
		json.Unmarshal([]byte(tenant.QuotaConfig), quota)
	}

	return quota, nil
}

// CheckQuotaLimit 检查配额限制
func (s *TenantService) CheckQuotaLimit(ctx context.Context, tenantID int64, resource string, currentCount int) (bool, error) {
	quota, err := s.GetTenantQuota(ctx, tenantID)
	if err != nil {
		return false, err
	}

	var limit int
	switch resource {
	case "users":
		limit = quota.MaxUsers
	case "accounts":
		limit = quota.MaxAccounts
	case "contents":
		limit = quota.MaxContents
	case "channels":
		limit = quota.MaxChannels
	default:
		return true, nil
	}

	return currentCount < limit, nil
}

// HasFeature 检查租户是否启用某功能
func (s *TenantService) HasFeature(ctx context.Context, tenantID int64, feature string) (bool, error) {
	quota, err := s.GetTenantQuota(ctx, tenantID)
	if err != nil {
		return false, err
	}

	for _, f := range quota.EnabledFeatures {
		if f == feature || f == "*" {
			return true, nil
		}
	}

	return false, nil
}

// ==================== 租户数据隔离 ====================

// TenantContext 租户上下文
type TenantContext struct {
	TenantID   int64  `json:"tenant_id"`
	TenantCode string `json:"tenant_code"`
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
}

// WithTenant 将租户信息注入上下文
func WithTenant(ctx context.Context, tc *TenantContext) context.Context {
	return context.WithValue(ctx, "tenant", tc)
}

// GetTenantFromContext 从上下文获取租户信息
func GetTenantFromContext(ctx context.Context) *TenantContext {
	if tc, ok := ctx.Value("tenant").(*TenantContext); ok {
		return tc
	}
	return nil
}

// ==================== 请求/响应模型 ====================

type CreateTenantRequest struct {
	Name         string       `json:"name"`
	Code         string       `json:"code"`
	ContactName  string       `json:"contact_name"`
	ContactEmail string       `json:"contact_email"`
	QuotaConfig  *TenantQuota `json:"quota_config"`
}

type UpdateTenantRequest struct {
	ContactName  string       `json:"contact_name"`
	ContactEmail string       `json:"contact_email"`
	QuotaConfig  *TenantQuota `json:"quota_config"`
}
