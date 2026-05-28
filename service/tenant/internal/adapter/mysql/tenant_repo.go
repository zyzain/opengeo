package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"opengeo/service/tenant/internal/domain"
	"opengeo/service/tenant/internal/port"
)

// TenantRepository MySQL 租户仓储实现
type TenantRepository struct {
	db *sql.DB
}

// NewTenantRepository 创建租户仓储
func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Save 保存租户（创建或更新）
func (r *TenantRepository) Save(ctx context.Context, tenant *domain.Tenant) error {
	if tenant.ID == 0 {
		return r.create(ctx, tenant)
	}
	return r.update(ctx, tenant)
}

// create 创建租户
func (r *TenantRepository) create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (name, slug, domain, logo_url, plan, status, brand_limit, user_limit, storage_limit, api_quota, api_used, settings, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	settingsJSON, err := marshalJSON(tenant.Settings)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query,
		tenant.Name,
		tenant.Slug,
		tenant.Domain,
		tenant.LogoURL,
		tenant.Plan,
		tenant.Status,
		tenant.BrandLimit,
		tenant.UserLimit,
		tenant.StorageLimit,
		tenant.APIQuota,
		tenant.APIUsed,
		settingsJSON,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	tenant.ID = id

	return nil
}

// update 更新租户
func (r *TenantRepository) update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants
		SET name = ?, slug = ?, domain = ?, logo_url = ?, plan = ?, status = ?,
		    brand_limit = ?, user_limit = ?, storage_limit = ?, api_quota = ?, api_used = ?,
		    settings = ?, updated_at = ?
		WHERE id = ?
	`
	tenant.UpdatedAt = time.Now()

	settingsJSON, err := marshalJSON(tenant.Settings)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query,
		tenant.Name,
		tenant.Slug,
		tenant.Domain,
		tenant.LogoURL,
		tenant.Plan,
		tenant.Status,
		tenant.BrandLimit,
		tenant.UserLimit,
		tenant.StorageLimit,
		tenant.APIQuota,
		tenant.APIUsed,
		settingsJSON,
		tenant.UpdatedAt,
		tenant.ID,
	)
	return err
}

// FindByID 根据 ID 查找租户
func (r *TenantRepository) FindByID(ctx context.Context, id int64) (*domain.Tenant, error) {
	query := `
		SELECT id, name, slug, domain, logo_url, plan, status, brand_limit, user_limit, storage_limit,
		       api_quota, api_used, quota_reset_at, settings, created_at, updated_at
		FROM tenants
		WHERE id = ? AND status != 4
	`
	return r.findOne(ctx, query, id)
}

// FindBySlug 根据标识查找租户
func (r *TenantRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, slug, domain, logo_url, plan, status, brand_limit, user_limit, storage_limit,
		       api_quota, api_used, quota_reset_at, settings, created_at, updated_at
		FROM tenants
		WHERE slug = ? AND status != 4
	`
	return r.findOne(ctx, query, slug)
}

// FindByDomain 根据域名查找租户
func (r *TenantRepository) FindByDomain(ctx context.Context, domain string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, slug, domain, logo_url, plan, status, brand_limit, user_limit, storage_limit,
		       api_quota, api_used, quota_reset_at, settings, created_at, updated_at
		FROM tenants
		WHERE domain = ? AND status != 4
	`
	return r.findOne(ctx, query, domain)
}

// findOne 查找单个租户
func (r *TenantRepository) findOne(ctx context.Context, query string, args ...interface{}) (*domain.Tenant, error) {
	tenant := &domain.Tenant{}
	var settingsJSON sql.NullString
	var quotaResetAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Slug,
		&tenant.Domain,
		&tenant.LogoURL,
		&tenant.Plan,
		&tenant.Status,
		&tenant.BrandLimit,
		&tenant.UserLimit,
		&tenant.StorageLimit,
		&tenant.APIQuota,
		&tenant.APIUsed,
		&quotaResetAt,
		&settingsJSON,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if quotaResetAt.Valid {
		tenant.QuotaResetAt = &quotaResetAt.Time
	}
	if settingsJSON.Valid {
		tenant.Settings, err = unmarshalJSON(settingsJSON.String)
		if err != nil {
			return nil, err
		}
	} else {
		tenant.Settings = make(map[string]string)
	}

	return tenant, nil
}

// List 列出租户
func (r *TenantRepository) List(ctx context.Context, filter *port.TenantFilter) ([]*domain.Tenant, int32, error) {
	// 构建查询条件
	where := "WHERE status != 4"
	args := []interface{}{}

	if filter.Status != nil {
		where += " AND status = ?"
		args = append(args, *filter.Status)
	}
	if filter.Plan != nil {
		where += " AND plan = ?"
		args = append(args, *filter.Plan)
	}
	if filter.Keyword != "" {
		where += " AND (name LIKE ? OR slug LIKE ? OR domain LIKE ?)"
		keyword := "%" + filter.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tenants %s", where)
	var total int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	orderBy := "ORDER BY created_at DESC"
	if filter.SortBy != "" {
		orderBy = fmt.Sprintf("ORDER BY %s %s", filter.SortBy, filter.SortOrder)
	}

	offset := (filter.Page - 1) * filter.PageSize
	query := fmt.Sprintf(`
		SELECT id, name, slug, domain, logo_url, plan, status, brand_limit, user_limit, storage_limit,
		       api_quota, api_used, quota_reset_at, settings, created_at, updated_at
		FROM tenants %s %s LIMIT ? OFFSET ?
	`, where, orderBy)

	args = append(args, filter.PageSize, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tenants := make([]*domain.Tenant, 0)
	for rows.Next() {
		tenant := &domain.Tenant{}
		var settingsJSON sql.NullString
		var quotaResetAt sql.NullTime

		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Slug,
			&tenant.Domain,
			&tenant.LogoURL,
			&tenant.Plan,
			&tenant.Status,
			&tenant.BrandLimit,
			&tenant.UserLimit,
			&tenant.StorageLimit,
			&tenant.APIQuota,
			&tenant.APIUsed,
			&quotaResetAt,
			&settingsJSON,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if quotaResetAt.Valid {
			tenant.QuotaResetAt = &quotaResetAt.Time
		}
		if settingsJSON.Valid {
			tenant.Settings, err = unmarshalJSON(settingsJSON.String)
			if err != nil {
				return nil, 0, err
			}
		} else {
			tenant.Settings = make(map[string]string)
		}

		tenants = append(tenants, tenant)
	}

	return tenants, total, nil
}

// Delete 删除租户（软删除）
func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE tenants SET status = 4, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// ExistsBySlug 检查标识是否存在
func (r *TenantRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tenants WHERE slug = ? AND status != 4)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&exists)
	return exists, err
}

// ExistsByDomain 检查域名是否存在
func (r *TenantRepository) ExistsByDomain(ctx context.Context, domain string) (bool, error) {
	if domain == "" {
		return false, nil
	}
	query := `SELECT EXISTS(SELECT 1 FROM tenants WHERE domain = ? AND status != 4)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, domain).Scan(&exists)
	return exists, err
}
