package domain

import (
	"errors"
	"time"
)

// TenantStatus 租户状态
type TenantStatus int32

const (
	TenantStatusUnspecified TenantStatus = 0
	TenantStatusActive      TenantStatus = 1
	TenantStatusDisabled    TenantStatus = 2
	TenantStatusExpired     TenantStatus = 3
	TenantStatusArchived    TenantStatus = 4
)

// TenantPlan 租户套餐
type TenantPlan int32

const (
	TenantPlanUnspecified TenantPlan = 0
	TenantPlanFree        TenantPlan = 1
	TenantPlanStarter     TenantPlan = 2
	TenantPlanPro         TenantPlan = 3
	TenantPlanEnterprise  TenantPlan = 4
)

// Tenant 租户聚合根
type Tenant struct {
	ID           int64
	Name         string
	Slug         string
	Domain       string
	LogoURL      string
	Plan         TenantPlan
	Status       TenantStatus
	BrandLimit   int32
	UserLimit    int32
	StorageLimit int64
	APIQuota     int32
	APIUsed      int32
	QuotaResetAt *time.Time
	Settings     map[string]string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TenantQuota 租户配额值对象
type TenantQuota struct {
	TenantID     int64
	BrandLimit   int32
	BrandCount   int32
	UserLimit    int32
	UserCount    int32
	StorageLimit int64
	StorageUsed  int64
	APIQuota     int32
	APIUsed      int32
	QuotaResetAt *time.Time
}

// TenantUsage 租户用量值对象
type TenantUsage struct {
	TenantID            int64
	PeriodStart         time.Time
	PeriodEnd           time.Time
	ContentCount        int32
	PublishCount        int32
	AIOptimizationCount int32
	APICallCount        int32
	APITokensUsed       int64
	APICostCents        int64
	APIUsageByType      map[string]int32
}

// 领域错误
var (
	ErrTenantNotFound      = errors.New("tenant not found")
	ErrTenantSlugExists    = errors.New("tenant slug already exists")
	ErrTenantDomainExists  = errors.New("tenant domain already exists")
	ErrTenantQuotaExceeded = errors.New("tenant quota exceeded")
	ErrTenantInactive      = errors.New("tenant is inactive")
)

// Validate 验证租户实体
func (t *Tenant) Validate() error {
	if t.Name == "" {
		return errors.New("tenant name is required")
	}
	if t.Slug == "" {
		return errors.New("tenant slug is required")
	}
	if t.BrandLimit < 0 {
		return errors.New("brand limit must be non-negative")
	}
	if t.UserLimit < 0 {
		return errors.New("user limit must be non-negative")
	}
	if t.StorageLimit < 0 {
		return errors.New("storage limit must be non-negative")
	}
	if t.APIQuota < 0 {
		return errors.New("API quota must be non-negative")
	}
	return nil
}

// IsActive 检查租户是否活跃
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

// IsQuotaExceeded 检查 API 配额是否超限
func (t *Tenant) IsQuotaExceeded() bool {
	return t.APIUsed >= t.APIQuota
}

// IncrementAPIUsed 增加 API 使用计数
func (t *Tenant) IncrementAPIUsed(count int32) {
	t.APIUsed += count
}

// ResetAPIUsage 重置 API 使用计数（每月初）
func (t *Tenant) ResetAPIUsage() {
	t.APIUsed = 0
	now := time.Now()
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	t.QuotaResetAt = &nextMonth
}
