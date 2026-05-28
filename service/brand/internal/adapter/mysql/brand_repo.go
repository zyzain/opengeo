package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"opengeo/service/brand/internal/domain"
	"opengeo/service/brand/internal/port"
)

// BrandRepository MySQL 品牌仓储实现
type BrandRepository struct {
	db *sql.DB
}

// NewBrandRepository 创建品牌仓储
func NewBrandRepository(db *sql.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

// Save 保存品牌
func (r *BrandRepository) Save(ctx context.Context, brand *domain.Brand) error {
	if brand.ID == 0 {
		return r.create(ctx, brand)
	}
	return r.update(ctx, brand)
}

func (r *BrandRepository) create(ctx context.Context, brand *domain.Brand) error {
	query := `
		INSERT INTO brands (tenant_id, name, slug, description, logo_url, website, industry, founded_year, headquarters, status, settings, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	brand.CreatedAt = now
	brand.UpdatedAt = now

	settingsJSON, _ := json.Marshal(brand.Settings)

	result, err := r.db.ExecContext(ctx, query,
		brand.TenantID, brand.Name, brand.Slug, brand.Description, brand.LogoURL,
		brand.Website, brand.Industry, brand.FoundedYear, brand.Headquarters,
		brand.Status, string(settingsJSON), brand.CreatedAt, brand.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	brand.ID = id
	return nil
}

func (r *BrandRepository) update(ctx context.Context, brand *domain.Brand) error {
	query := `
		UPDATE brands SET name=?, slug=?, description=?, logo_url=?, website=?, industry=?,
		founded_year=?, headquarters=?, status=?, settings=?, updated_at=? WHERE id=? AND tenant_id=?
	`
	brand.UpdatedAt = time.Now()
	settingsJSON, _ := json.Marshal(brand.Settings)

	_, err := r.db.ExecContext(ctx, query,
		brand.Name, brand.Slug, brand.Description, brand.LogoURL, brand.Website,
		brand.Industry, brand.FoundedYear, brand.Headquarters, brand.Status,
		string(settingsJSON), brand.UpdatedAt, brand.ID, brand.TenantID,
	)
	return err
}

// FindByID 根据 ID 查找品牌
func (r *BrandRepository) FindByID(ctx context.Context, tenantID, id int64) (*domain.Brand, error) {
	query := `SELECT id, tenant_id, name, slug, description, logo_url, website, industry, founded_year, headquarters, status, settings, created_at, updated_at FROM brands WHERE id=? AND tenant_id=? AND status!=3`
	return r.findOne(ctx, query, id, tenantID)
}

// FindBySlug 根据标识查找品牌
func (r *BrandRepository) FindBySlug(ctx context.Context, tenantID int64, slug string) (*domain.Brand, error) {
	query := `SELECT id, tenant_id, name, slug, description, logo_url, website, industry, founded_year, headquarters, status, settings, created_at, updated_at FROM brands WHERE slug=? AND tenant_id=? AND status!=3`
	return r.findOne(ctx, query, slug, tenantID)
}

func (r *BrandRepository) findOne(ctx context.Context, query string, args ...interface{}) (*domain.Brand, error) {
	brand := &domain.Brand{}
	var settingsJSON sql.NullString

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&brand.ID, &brand.TenantID, &brand.Name, &brand.Slug, &brand.Description,
		&brand.LogoURL, &brand.Website, &brand.Industry, &brand.FoundedYear,
		&brand.Headquarters, &brand.Status, &settingsJSON, &brand.CreatedAt, &brand.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if settingsJSON.Valid {
		json.Unmarshal([]byte(settingsJSON.String), &brand.Settings)
	}
	if brand.Settings == nil {
		brand.Settings = make(map[string]string)
	}
	return brand, nil
}

// List 列出品牌
func (r *BrandRepository) List(ctx context.Context, tenantID int64, filter *port.BrandFilter) ([]*domain.Brand, int32, error) {
	where := "WHERE tenant_id=? AND status!=3"
	args := []interface{}{tenantID}

	if filter.Status != nil {
		where += " AND status=?"
		args = append(args, *filter.Status)
	}
	if filter.Industry != "" {
		where += " AND industry=?"
		args = append(args, filter.Industry)
	}
	if filter.Keyword != "" {
		where += " AND (name LIKE ? OR description LIKE ?)"
		kw := "%" + filter.Keyword + "%"
		args = append(args, kw, kw)
	}

	var total int32
	r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM brands %s", where), args...).Scan(&total)

	orderBy := "ORDER BY created_at DESC"
	offset := (filter.Page - 1) * filter.PageSize
	query := fmt.Sprintf(`SELECT id, tenant_id, name, slug, description, logo_url, website, industry, founded_year, headquarters, status, settings, created_at, updated_at FROM brands %s %s LIMIT ? OFFSET ?`, where, orderBy)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	brands := make([]*domain.Brand, 0)
	for rows.Next() {
		brand := &domain.Brand{}
		var settingsJSON sql.NullString
		rows.Scan(&brand.ID, &brand.TenantID, &brand.Name, &brand.Slug, &brand.Description,
			&brand.LogoURL, &brand.Website, &brand.Industry, &brand.FoundedYear,
			&brand.Headquarters, &brand.Status, &settingsJSON, &brand.CreatedAt, &brand.UpdatedAt)
		if settingsJSON.Valid {
			json.Unmarshal([]byte(settingsJSON.String), &brand.Settings)
		}
		if brand.Settings == nil {
			brand.Settings = make(map[string]string)
		}
		brands = append(brands, brand)
	}
	return brands, total, nil
}

// Delete 删除品牌
func (r *BrandRepository) Delete(ctx context.Context, tenantID, id int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE brands SET status=3, updated_at=? WHERE id=? AND tenant_id=?", time.Now(), id, tenantID)
	return err
}

// ExistsBySlug 检查标识是否存在
func (r *BrandRepository) ExistsBySlug(ctx context.Context, tenantID int64, slug string) (bool, error) {
	var exists bool
	r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM brands WHERE slug=? AND tenant_id=? AND status!=3)", slug, tenantID).Scan(&exists)
	return exists, nil
}

// CountByTenant 统计租户下的品牌数量
func (r *BrandRepository) CountByTenant(ctx context.Context, tenantID int64) (int32, error) {
	var count int32
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM brands WHERE tenant_id=? AND status!=3", tenantID).Scan(&count)
	return count, nil
}
