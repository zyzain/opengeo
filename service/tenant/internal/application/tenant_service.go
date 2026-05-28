package application

import (
	"context"
	"time"

	"opengeo/service/tenant/internal/domain"
	"opengeo/service/tenant/internal/port"
)

// TenantService 租户应用服务
type TenantService struct {
	repo      port.TenantRepository
	quotaRepo port.TenantQuotaRepository
	usageRepo port.TenantUsageRepository
}

// NewTenantService 创建租户服务
func NewTenantService(
	repo port.TenantRepository,
	quotaRepo port.TenantQuotaRepository,
	usageRepo port.TenantUsageRepository,
) *TenantService {
	return &TenantService{
		repo:      repo,
		quotaRepo: quotaRepo,
		usageRepo: usageRepo,
	}
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Name         string
	Slug         string
	Domain       string
	Plan         domain.TenantPlan
	AdminEmail   string
	AdminUsername string
	AdminPassword string
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*domain.Tenant, error) {
	// 检查标识是否已存在
	exists, err := s.repo.ExistsBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTenantSlugExists
	}

	// 检查域名是否已存在（如果提供了域名）
	if req.Domain != "" {
		exists, err = s.repo.ExistsByDomain(ctx, req.Domain)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrTenantDomainExists
		}
	}

	// 根据套餐设置默认配额
	brandLimit, userLimit, storageLimit, apiQuota := getDefaultQuota(req.Plan)

	// 创建租户实体
	tenant := &domain.Tenant{
		Name:         req.Name,
		Slug:         req.Slug,
		Domain:       req.Domain,
		Plan:         req.Plan,
		Status:       domain.TenantStatusActive,
		BrandLimit:   brandLimit,
		UserLimit:    userLimit,
		StorageLimit: storageLimit,
		APIQuota:     apiQuota,
		APIUsed:      0,
		Settings:     make(map[string]string),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 验证实体
	if err := tenant.Validate(); err != nil {
		return nil, err
	}

	// 保存租户
	if err := s.repo.Save(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// GetTenant 获取租户信息
func (s *TenantService) GetTenant(ctx context.Context, id int64) (*domain.Tenant, error) {
	tenant, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, domain.ErrTenantNotFound
	}
	return tenant, nil
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	ID      int64
	Name    *string
	Domain  *string
	LogoURL *string
	Status  *domain.TenantStatus
	Settings map[string]string
}

// UpdateTenant 更新租户信息
func (s *TenantService) UpdateTenant(ctx context.Context, req *UpdateTenantRequest) (*domain.Tenant, error) {
	tenant, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, domain.ErrTenantNotFound
	}

	// 更新字段
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Domain != nil {
		// 检查域名是否已被其他租户使用
		if *req.Domain != "" && *req.Domain != tenant.Domain {
			exists, err := s.repo.ExistsByDomain(ctx, *req.Domain)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, domain.ErrTenantDomainExists
			}
		}
		tenant.Domain = *req.Domain
	}
	if req.LogoURL != nil {
		tenant.LogoURL = *req.LogoURL
	}
	if req.Status != nil {
		tenant.Status = *req.Status
	}
	if req.Settings != nil {
		tenant.Settings = req.Settings
	}
	tenant.UpdatedAt = time.Now()

	// 验证实体
	if err := tenant.Validate(); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.repo.Save(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// DeleteTenant 删除租户
func (s *TenantService) DeleteTenant(ctx context.Context, id int64) error {
	tenant, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if tenant == nil {
		return domain.ErrTenantNotFound
	}

	return s.repo.Delete(ctx, id)
}

// ListTenantsRequest 列出租户请求
type ListTenantsRequest struct {
	Status   *domain.TenantStatus
	Plan     *domain.TenantPlan
	Keyword  string
	Page     int32
	PageSize int32
}

// ListTenantsResponse 列出租户响应
type ListTenantsResponse struct {
	Tenants    []*domain.Tenant
	Total      int32
	Page       int32
	PageSize   int32
	TotalPages int32
}

// ListTenants 列出租户
func (s *TenantService) ListTenants(ctx context.Context, req *ListTenantsRequest) (*ListTenantsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	filter := &port.TenantFilter{
		Status:   req.Status,
		Plan:     req.Plan,
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	tenants, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := (total + int32(req.PageSize) - 1) / int32(req.PageSize)

	return &ListTenantsResponse{
		Tenants:    tenants,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetTenantQuota 获取租户配额
func (s *TenantService) GetTenantQuota(ctx context.Context, tenantID int64) (*domain.TenantQuota, error) {
	// 检查租户是否存在
	tenant, err := s.repo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, domain.ErrTenantNotFound
	}

	return s.quotaRepo.GetQuota(ctx, tenantID)
}

// UpdateTenantQuotaRequest 更新租户配额请求
type UpdateTenantQuotaRequest struct {
	TenantID     int64
	BrandLimit   *int32
	UserLimit    *int32
	StorageLimit *int64
	APIQuota     *int32
}

// UpdateTenantQuota 更新租户配额
func (s *TenantService) UpdateTenantQuota(ctx context.Context, req *UpdateTenantQuotaRequest) (*domain.TenantQuota, error) {
	// 检查租户是否存在
	tenant, err := s.repo.FindByID(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, domain.ErrTenantNotFound
	}

	// 获取当前配额
	quota, err := s.quotaRepo.GetQuota(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}

	// 更新配额字段
	if req.BrandLimit != nil {
		quota.BrandLimit = *req.BrandLimit
		tenant.BrandLimit = *req.BrandLimit
	}
	if req.UserLimit != nil {
		quota.UserLimit = *req.UserLimit
		tenant.UserLimit = *req.UserLimit
	}
	if req.StorageLimit != nil {
		quota.StorageLimit = *req.StorageLimit
		tenant.StorageLimit = *req.StorageLimit
	}
	if req.APIQuota != nil {
		quota.APIQuota = *req.APIQuota
		tenant.APIQuota = *req.APIQuota
	}

	// 保存更新
	if err := s.quotaRepo.UpdateQuota(ctx, req.TenantID, quota); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, tenant); err != nil {
		return nil, err
	}

	return quota, nil
}

// GetTenantUsage 获取租户用量统计
func (s *TenantService) GetTenantUsage(ctx context.Context, tenantID int64, periodStart, periodEnd string) (*domain.TenantUsage, error) {
	// 检查租户是否存在
	tenant, err := s.repo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, domain.ErrTenantNotFound
	}

	return s.usageRepo.GetUsage(ctx, tenantID, periodStart, periodEnd)
}

// getDefaultQuota 根据套餐获取默认配额
func getDefaultQuota(plan domain.TenantPlan) (brandLimit, userLimit int32, storageLimit int64, apiQuota int32) {
	switch plan {
	case domain.TenantPlanFree:
		return 5, 10, 1073741824, 100        // 5 品牌, 10 用户, 1GB, 100 API/月
	case domain.TenantPlanStarter:
		return 20, 50, 10737418240, 1000     // 20 品牌, 50 用户, 10GB, 1000 API/月
	case domain.TenantPlanPro:
		return 100, 200, 107374182400, 10000 // 100 品牌, 200 用户, 100GB, 10000 API/月
	case domain.TenantPlanEnterprise:
		return 1000, 2000, 1099511627776, 100000 // 1000 品牌, 2000 用户, 1TB, 100000 API/月
	default:
		return 5, 10, 1073741824, 100
	}
}
