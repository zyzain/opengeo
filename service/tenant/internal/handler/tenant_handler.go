package handler

import (
	"context"

	"opengeo/service/tenant/internal/application"
	"opengeo/service/tenant/internal/domain"
)

// TenantHandler 租户 RPC Handler
type TenantHandler struct {
	service *application.TenantService
}

// NewTenantHandler 创建租户 Handler
func NewTenantHandler(service *application.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

// CreateTenant 创建租户
func (h *TenantHandler) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error) {
	result, err := h.service.CreateTenant(ctx, &application.CreateTenantRequest{
		Name:          req.Name,
		Slug:          req.Slug,
		Domain:        req.Domain,
		Plan:          domain.TenantPlan(req.Plan),
		AdminEmail:    req.AdminEmail,
		AdminUsername: req.AdminUsername,
		AdminPassword: req.AdminPassword,
	})
	if err != nil {
		return nil, err
	}

	return toTenantResponse(result), nil
}

// GetTenant 获取租户
func (h *TenantHandler) GetTenant(ctx context.Context, req *GetTenantRequest) (*TenantResponse, error) {
	result, err := h.service.GetTenant(ctx, req.TenantId)
	if err != nil {
		return nil, err
	}

	return toTenantResponse(result), nil
}

// UpdateTenant 更新租户
func (h *TenantHandler) UpdateTenant(ctx context.Context, req *UpdateTenantRequest) (*TenantResponse, error) {
	updateReq := &application.UpdateTenantRequest{
		ID: req.TenantId,
	}

	if req.Name != "" {
		updateReq.Name = &req.Name
	}
	if req.Domain != "" {
		updateReq.Domain = &req.Domain
	}
	if req.LogoUrl != "" {
		updateReq.LogoUrl = &req.LogoUrl
	}
	if req.Status > 0 {
		status := domain.TenantStatus(req.Status)
		updateReq.Status = &status
	}
	if req.Settings != nil {
		updateReq.Settings = req.Settings
	}

	result, err := h.service.UpdateTenant(ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return toTenantResponse(result), nil
}

// DeleteTenant 删除租户
func (h *TenantHandler) DeleteTenant(ctx context.Context, req *DeleteTenantRequest) (*DeleteTenantResponse, error) {
	err := h.service.DeleteTenant(ctx, req.TenantId)
	if err != nil {
		return nil, err
	}

	return &DeleteTenantResponse{Success: true}, nil
}

// ListTenants 列出租户
func (h *TenantHandler) ListTenants(ctx context.Context, req *ListTenantsRequest) (*ListTenantsResponse, error) {
	listReq := &application.ListTenantsRequest{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	if req.Status > 0 {
		status := domain.TenantStatus(req.Status)
		listReq.Status = &status
	}
	if req.Plan > 0 {
		plan := domain.TenantPlan(req.Plan)
		listReq.Plan = &plan
	}

	result, err := h.service.ListTenants(ctx, listReq)
	if err != nil {
		return nil, err
	}

	tenants := make([]*TenantResponse, len(result.Tenants))
	for i, t := range result.Tenants {
		tenants[i] = toTenantResponse(t)
	}

	return &ListTenantsResponse{
		Tenants:    tenants,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// GetTenantQuota 获取租户配额
func (h *TenantHandler) GetTenantQuota(ctx context.Context, req *GetTenantQuotaRequest) (*TenantQuotaResponse, error) {
	result, err := h.service.GetTenantQuota(ctx, req.TenantId)
	if err != nil {
		return nil, err
	}

	return &TenantQuotaResponse{
		TenantId:     result.TenantID,
		BrandLimit:   result.BrandLimit,
		BrandCount:   result.BrandCount,
		UserLimit:    result.UserLimit,
		UserCount:    result.UserCount,
		StorageLimit: result.StorageLimit,
		StorageUsed:  result.StorageUsed,
		ApiQuota:     result.APIQuota,
		ApiUsed:      result.APIUsed,
	}, nil
}

// toTenantResponse 转换为响应
func toTenantResponse(t *domain.Tenant) *TenantResponse {
	return &TenantResponse{
		Id:           t.ID,
		Name:         t.Name,
		Slug:         t.Slug,
		Domain:       t.Domain,
		LogoUrl:      t.LogoURL,
		Plan:         int32(t.Plan),
		Status:       int32(t.Status),
		BrandLimit:   t.BrandLimit,
		UserLimit:    t.UserLimit,
		StorageLimit: t.StorageLimit,
		ApiQuota:     t.APIQuota,
		ApiUsed:      t.APIUsed,
		Settings:     t.Settings,
		CreatedAt:    t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// 请求/响应结构体
type CreateTenantRequest struct {
	Name          string
	Slug          string
	Domain        string
	Plan          int32
	AdminEmail    string
	AdminUsername string
	AdminPassword string
}

type GetTenantRequest struct {
	TenantId int64
}

type UpdateTenantRequest struct {
	TenantId int64
	Name     string
	Domain   string
	LogoUrl  string
	Status   int32
	Settings map[string]string
}

type DeleteTenantRequest struct {
	TenantId int64
}

type DeleteTenantResponse struct {
	Success bool
}

type ListTenantsRequest struct {
	Status   int32
	Plan     int32
	Keyword  string
	Page     int32
	PageSize int32
}

type ListTenantsResponse struct {
	Tenants    []*TenantResponse
	Total      int32
	Page       int32
	PageSize   int32
	TotalPages int32
}

type GetTenantQuotaRequest struct {
	TenantId int64
}

type TenantResponse struct {
	Id           int64
	Name         string
	Slug         string
	Domain       string
	LogoUrl      string
	Plan         int32
	Status       int32
	BrandLimit   int32
	UserLimit    int32
	StorageLimit int64
	ApiQuota     int32
	ApiUsed      int32
	Settings     map[string]string
	CreatedAt    string
	UpdatedAt    string
}

type TenantQuotaResponse struct {
	TenantId     int64
	BrandLimit   int32
	BrandCount   int32
	UserLimit    int32
	UserCount    int32
	StorageLimit int64
	StorageUsed  int64
	ApiQuota     int32
	ApiUsed      int32
}
