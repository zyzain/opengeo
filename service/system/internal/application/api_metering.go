package application

import (
	"context"
	"errors"
	"time"
)

// APIPlan API 套餐
type APIPlan struct {
	ID           int64
	Name         string
	Code         string
	Description  string
	MonthlyQuota int32
	PriceCents   int64
	OveragePrice int64
	Features     []string
}

// TenantSubscription 租户订阅
type TenantSubscription struct {
	ID         int64
	TenantID   int64
	PlanID     int64
	Plan       *APIPlan
	Status     string
	StartsAt   time.Time
	EndsAt     *time.Time
	APIKey     string
	APISecret  string
}

// APIUsage API 用量
type APIUsage struct {
	TenantID      int64
	PlanCode      string
	MonthlyQuota  int32
	UsedThisMonth int32
	Remaining     int32
	QuotaResetAt  time.Time
	UsageByType   map[string]int32
	CostThisMonth int64
}

// BillingRecord 计费记录
type BillingRecord struct {
	ID          int64
	TenantID    int64
	APIType     string
	Endpoint    string
	RequestID   string
	TokensInput int64
	TokensOutput int64
	CostCents   int64
	IsOverage   bool
	Status      string
	CreatedAt   time.Time
}

// MeteringService 计量服务
type MeteringService struct {
	subscriptionRepo SubscriptionRepository
	usageRepo        UsageRepository
	billingRepo      BillingRepository
}

// SubscriptionRepository 订阅仓储接口
type SubscriptionRepository interface {
	FindByTenantID(ctx context.Context, tenantID int64) (*TenantSubscription, error)
	Save(ctx context.Context, sub *TenantSubscription) error
}

// UsageRepository 用量仓储接口
type UsageRepository interface {
	GetUsage(ctx context.Context, tenantID int64) (*APIUsage, error)
	IncrementUsage(ctx context.Context, tenantID int64, apiType string, count int32) error
	ResetUsage(ctx context.Context, tenantID int64) error
}

// BillingRepository 计费仓储接口
type BillingRepository interface {
	Save(ctx context.Context, record *BillingRecord) error
	List(ctx context.Context, tenantID int64, page, pageSize int32) ([]*BillingRecord, int32, error)
}

// NewMeteringService 创建计量服务
func NewMeteringService(
	subRepo SubscriptionRepository,
	usageRepo UsageRepository,
	billingRepo BillingRepository,
) *MeteringService {
	return &MeteringService{
		subscriptionRepo: subRepo,
		usageRepo:        usageRepo,
		billingRepo:      billingRepo,
	}
}

// CheckQuotaRequest 检查配额请求
type CheckQuotaRequest struct {
	TenantID int64
	APIType  string
}

// CheckQuotaResponse 检查配额响应
type CheckQuotaResponse struct {
	Allowed     bool
	Remaining   int32
	ResetAt     time.Time
	IsOverage   bool
	OverageCost int64
}

// CheckQuota 检查配额
func (s *MeteringService) CheckQuota(ctx context.Context, req *CheckQuotaRequest) (*CheckQuotaResponse, error) {
	// 获取订阅信息
	sub, err := s.subscriptionRepo.FindByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("no active subscription")
	}

	// 获取用量
	usage, err := s.usageRepo.GetUsage(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}

	// 检查是否超配额
	remaining := sub.Plan.MonthlyQuota - usage.UsedThisMonth
	if remaining <= 0 {
		// 计算超额费用
		overageCost := sub.Plan.OveragePrice
		return &CheckQuotaResponse{
			Allowed:     true, // 允许超额调用，但需要计费
			Remaining:   0,
			ResetAt:     usage.QuotaResetAt,
			IsOverage:   true,
			OverageCost: overageCost,
		}, nil
	}

	return &CheckQuotaResponse{
		Allowed:   true,
		Remaining: remaining,
		ResetAt:   usage.QuotaResetAt,
	}, nil
}

// RecordAPIUsageRequest 记录 API 用量请求
type RecordAPIUsageRequest struct {
	TenantID     int64
	APIType      string
	Endpoint     string
	RequestID    string
	TokensInput  int64
	TokensOutput int64
	CostCents    int64
}

// RecordAPIUsage 记录 API 用量
func (s *MeteringService) RecordAPIUsage(ctx context.Context, req *RecordAPIUsageRequest) error {
	// 增加用量计数
	if err := s.usageRepo.IncrementUsage(ctx, req.TenantID, req.APIType, 1); err != nil {
		return err
	}

	// 保存计费记录
	record := &BillingRecord{
		TenantID:     req.TenantID,
		APIType:      req.APIType,
		Endpoint:     req.Endpoint,
		RequestID:    req.RequestID,
		TokensInput:  req.TokensInput,
		TokensOutput: req.TokensOutput,
		CostCents:    req.CostCents,
		Status:       "success",
		CreatedAt:    time.Now(),
	}

	return s.billingRepo.Save(ctx, record)
}

// GetUsage 获取用量
func (s *MeteringService) GetUsage(ctx context.Context, tenantID int64) (*APIUsage, error) {
	return s.usageRepo.GetUsage(ctx, tenantID)
}

// GetBillingRecords 获取计费记录
func (s *MeteringService) GetBillingRecords(ctx context.Context, tenantID int64, page, pageSize int32) ([]*BillingRecord, int32, error) {
	return s.billingRepo.List(ctx, tenantID, page, pageSize)
}
