package port

import (
	"context"

	"opengeo/service/tenant/internal/domain"
)

// TenantRepository 租户持久化接口（Domain 层定义，Adapter 层实现）
type TenantRepository interface {
	// Save 保存租户（创建或更新）
	Save(ctx context.Context, tenant *domain.Tenant) error

	// FindByID 根据 ID 查找租户
	FindByID(ctx context.Context, id int64) (*domain.Tenant, error)

	// FindBySlug 根据标识查找租户
	FindBySlug(ctx context.Context, slug string) (*domain.Tenant, error)

	// FindByDomain 根据域名查找租户
	FindByDomain(ctx context.Context, domain string) (*domain.Tenant, error)

	// List 列出租户
	List(ctx context.Context, filter *TenantFilter) ([]*domain.Tenant, int32, error)

	// Delete 删除租户（软删除）
	Delete(ctx context.Context, id int64) error

	// ExistsBySlug 检查标识是否存在
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	// ExistsByDomain 检查域名是否存在
	ExistsByDomain(ctx context.Context, domain string) (bool, error)
}

// TenantFilter 租户查询过滤器
type TenantFilter struct {
	Status  *domain.TenantStatus
	Plan    *domain.TenantPlan
	Keyword string
	Page    int32
	PageSize int32
	SortBy  string
	SortOrder string
}

// TenantQuotaRepository 租户配额仓储接口
type TenantQuotaRepository interface {
	// GetQuota 获取租户配额
	GetQuota(ctx context.Context, tenantID int64) (*domain.TenantQuota, error)

	// UpdateQuota 更新租户配额
	UpdateQuota(ctx context.Context, tenantID int64, quota *domain.TenantQuota) error

	// IncrementBrandCount 增加品牌计数
	IncrementBrandCount(ctx context.Context, tenantID int64) error

	// DecrementBrandCount 减少品牌计数
	DecrementBrandCount(ctx context.Context, tenantID int64) error

	// IncrementUserCount 增加用户计数
	IncrementUserCount(ctx context.Context, tenantID int64) error

	// DecrementUserCount 减少用户计数
	DecrementUserCount(ctx context.Context, tenantID int64) error

	// IncrementAPIUsed 增加 API 使用计数
	IncrementAPIUsed(ctx context.Context, tenantID int64, count int32) error

	// ResetAPIUsage 重置 API 使用计数
	ResetAPIUsage(ctx context.Context, tenantID int64) error
}

// TenantUsageRepository 租户用量统计仓储接口
type TenantUsageRepository interface {
	// GetUsage 获取租户用量统计
	GetUsage(ctx context.Context, tenantID int64, periodStart, periodEnd string) (*domain.TenantUsage, error)

	// RecordAPIUsage 记录 API 使用
	RecordAPIUsage(ctx context.Context, tenantID int64, apiType string, tokensUsed int64, costCents int64) error
}
