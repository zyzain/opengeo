package port

import (
	"context"

	"opengeo/service/brand/internal/domain"
)

// BrandRepository 品牌仓储接口
type BrandRepository interface {
	// Save 保存品牌（创建或更新）
	Save(ctx context.Context, brand *domain.Brand) error

	// FindByID 根据 ID 查找品牌
	FindByID(ctx context.Context, tenantID, id int64) (*domain.Brand, error)

	// FindBySlug 根据标识查找品牌
	FindBySlug(ctx context.Context, tenantID int64, slug string) (*domain.Brand, error)

	// List 列出品牌
	List(ctx context.Context, tenantID int64, filter *BrandFilter) ([]*domain.Brand, int32, error)

	// Delete 删除品牌（软删除）
	Delete(ctx context.Context, tenantID, id int64) error

	// ExistsBySlug 检查标识是否存在
	ExistsBySlug(ctx context.Context, tenantID int64, slug string) (bool, error)

	// CountByTenant 统计租户下的品牌数量
	CountByTenant(ctx context.Context, tenantID int64) (int32, error)
}

// BrandFilter 品牌查询过滤器
type BrandFilter struct {
	Status   *domain.BrandStatus
	Industry string
	Keyword  string
	Page     int32
	PageSize int32
	SortBy   string
	SortOrder string
}

// BrandMetadataRepository 品牌元数据仓储接口
type BrandMetadataRepository interface {
	// Save 保存品牌元数据
	Save(ctx context.Context, metadata *domain.BrandMetadata) error

	// FindByBrandID 根据品牌 ID 查找元数据
	FindByBrandID(ctx context.Context, tenantID, brandID int64) (*domain.BrandMetadata, error)

	// Delete 删除品牌元数据
	Delete(ctx context.Context, tenantID, brandID int64) error
}

// GlossaryRepository 术语表仓储接口
type GlossaryRepository interface {
	// Save 保存术语条目
	Save(ctx context.Context, entry *domain.GlossaryEntry) error

	// FindByID 根据 ID 查找术语条目
	FindByID(ctx context.Context, tenantID, id int64) (*domain.GlossaryEntry, error)

	// List 列出术语条目
	List(ctx context.Context, tenantID, brandID int64, filter *GlossaryFilter) ([]*domain.GlossaryEntry, int32, error)

	// Delete 删除术语条目
	Delete(ctx context.Context, tenantID, id int64) error

	// BulkCreate 批量创建术语条目
	BulkCreate(ctx context.Context, entries []*domain.GlossaryEntry) (int32, int32, error)

	// ExistsByTerm 检查术语是否存在
	ExistsByTerm(ctx context.Context, tenantID, brandID int64, term string) (bool, error)
}

// GlossaryFilter 术语查询过滤器
type GlossaryFilter struct {
	Category    string
	Keyword     string
	IsForbidden *bool
	Page        int32
	PageSize    int32
}

// SnapshotRepository 品牌快照仓储接口
type SnapshotRepository interface {
	// Save 保存快照
	Save(ctx context.Context, snapshot *domain.BrandSnapshot) error

	// FindByID 根据 ID 查找快照
	FindByID(ctx context.Context, tenantID, id int64) (*domain.BrandSnapshot, error)

	// List 列出快照
	List(ctx context.Context, tenantID, brandID int64, page, pageSize int32) ([]*domain.BrandSnapshot, int32, error)
}
