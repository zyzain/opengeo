package mysql

import (
	"context"
	"database/sql"
	"time"

	"opengeo/service/tenant/internal/domain"
)

// TenantQuotaRepository MySQL 租户配额仓储实现
type TenantQuotaRepository struct {
	db *sql.DB
}

// NewTenantQuotaRepository 创建租户配额仓储
func NewTenantQuotaRepository(db *sql.DB) *TenantQuotaRepository {
	return &TenantQuotaRepository{db: db}
}

// GetQuota 获取租户配额
func (r *TenantQuotaRepository) GetQuota(ctx context.Context, tenantID int64) (*domain.TenantQuota, error) {
	query := `
		SELECT t.id, t.brand_limit, COUNT(b.id) as brand_count,
		       t.user_limit, COUNT(u.id) as user_count,
		       t.storage_limit, 0 as storage_used,
		       t.api_quota, t.api_used, t.quota_reset_at
		FROM tenants t
		LEFT JOIN brands b ON b.tenant_id = t.id AND b.status != 3
		LEFT JOIN users u ON u.tenant_id = t.id AND u.status != 0
		WHERE t.id = ?
		GROUP BY t.id
	`

	quota := &domain.TenantQuota{}
	var quotaResetAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&quota.TenantID,
		&quota.BrandLimit,
		&quota.BrandCount,
		&quota.UserLimit,
		&quota.UserCount,
		&quota.StorageLimit,
		&quota.StorageUsed,
		&quota.APIQuota,
		&quota.APIUsed,
		&quotaResetAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if quotaResetAt.Valid {
		quota.QuotaResetAt = &quotaResetAt.Time
	}

	return quota, nil
}

// UpdateQuota 更新租户配额
func (r *TenantQuotaRepository) UpdateQuota(ctx context.Context, tenantID int64, quota *domain.TenantQuota) error {
	query := `
		UPDATE tenants
		SET brand_limit = ?, user_limit = ?, storage_limit = ?, api_quota = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		quota.BrandLimit,
		quota.UserLimit,
		quota.StorageLimit,
		quota.APIQuota,
		time.Now(),
		tenantID,
	)
	return err
}

// IncrementBrandCount 增加品牌计数（通过触发器或应用层实现）
func (r *TenantQuotaRepository) IncrementBrandCount(ctx context.Context, tenantID int64) error {
	// 实际实现中，品牌计数是通过 COUNT 查询动态计算的
	// 这里可以添加缓存或计数器表来优化性能
	return nil
}

// DecrementBrandCount 减少品牌计数
func (r *TenantQuotaRepository) DecrementBrandCount(ctx context.Context, tenantID int64) error {
	return nil
}

// IncrementUserCount 增加用户计数
func (r *TenantQuotaRepository) IncrementUserCount(ctx context.Context, tenantID int64) error {
	return nil
}

// DecrementUserCount 减少用户计数
func (r *TenantQuotaRepository) DecrementUserCount(ctx context.Context, tenantID int64) error {
	return nil
}

// IncrementAPIUsed 增加 API 使用计数
func (r *TenantQuotaRepository) IncrementAPIUsed(ctx context.Context, tenantID int64, count int32) error {
	query := `UPDATE tenants SET api_used = api_used + ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, count, time.Now(), tenantID)
	return err
}

// ResetAPIUsage 重置 API 使用计数
func (r *TenantQuotaRepository) ResetAPIUsage(ctx context.Context, tenantID int64) error {
	nextMonth := time.Now().AddDate(0, 1, 0)
	nextMonth = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	query := `UPDATE tenants SET api_used = 0, quota_reset_at = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, nextMonth, time.Now(), tenantID)
	return err
}
