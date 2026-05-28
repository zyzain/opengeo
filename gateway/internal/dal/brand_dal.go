package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/gateway/internal/model"
)

// BrandRepository 品牌仓储
type BrandRepository struct {
	db *gorm.DB
}

// NewBrandRepository 创建品牌仓储
func NewBrandRepository(db *gorm.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

// List 列出品牌
func (r *BrandRepository) List(ctx context.Context, tenantID int64, filter *BrandFilter) ([]*model.Brand, int32, error) {
	var brands []*model.Brand
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Brand{}).Where("tenant_id = ? AND status != 3", tenantID)

	if filter.Keyword != "" {
		query = query.Where("name LIKE ? OR slug LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if filter.Industry != "" {
		query = query.Where("industry = ?", filter.Industry)
	}
	if filter.Status > 0 {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count brands: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Order("created_at DESC").Find(&brands).Error; err != nil {
		return nil, 0, fmt.Errorf("list brands: %w", err)
	}

	return brands, int32(total), nil
}

// FindByID 根据ID查找品牌
func (r *BrandRepository) FindByID(ctx context.Context, tenantID, id int64) (*model.Brand, error) {
	var brand model.Brand
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ? AND status != 3", tenantID, id).First(&brand).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find brand: %w", err)
	}
	return &brand, nil
}

// Create 创建品牌
func (r *BrandRepository) Create(ctx context.Context, tenantID int64, req *BrandCreateRequest) (*model.Brand, error) {
	brand := &model.Brand{
		TenantID:     tenantID,
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		LogoURL:      req.LogoURL,
		Website:      req.Website,
		Industry:     req.Industry,
		FoundedYear:  req.FoundedYear,
		Headquarters: req.Headquarters,
		Status:       1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(brand).Error; err != nil {
		return nil, fmt.Errorf("create brand: %w", err)
	}

	return brand, nil
}

// Update 更新品牌
func (r *BrandRepository) Update(ctx context.Context, tenantID, id int64, req *BrandUpdateRequest) (*model.Brand, error) {
	brand, err := r.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, fmt.Errorf("brand not found")
	}

	if req.Name != nil {
		brand.Name = *req.Name
	}
	if req.Description != nil {
		brand.Description = *req.Description
	}
	if req.LogoURL != nil {
		brand.LogoURL = *req.LogoURL
	}
	if req.Website != nil {
		brand.Website = *req.Website
	}
	if req.Industry != nil {
		brand.Industry = *req.Industry
	}
	if req.FoundedYear != nil {
		brand.FoundedYear = *req.FoundedYear
	}
	if req.Headquarters != nil {
		brand.Headquarters = *req.Headquarters
	}
	if req.Status != nil {
		brand.Status = *req.Status
	}
	brand.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(brand).Error; err != nil {
		return nil, fmt.Errorf("update brand: %w", err)
	}

	return brand, nil
}

// Delete 删除品牌（软删除）
func (r *BrandRepository) Delete(ctx context.Context, tenantID, id int64) error {
	if err := r.db.WithContext(ctx).Model(&model.Brand{}).Where("tenant_id = ? AND id = ?", tenantID, id).Update("status", 3).Error; err != nil {
		return fmt.Errorf("delete brand: %w", err)
	}
	return nil
}

// GetMetadata 获取品牌元数据
func (r *BrandRepository) GetMetadata(ctx context.Context, tenantID, brandID int64) (*model.BrandMetadata, error) {
	var metadata model.BrandMetadata
	if err := r.db.WithContext(ctx).Where("brand_id = ?", brandID).First(&metadata).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get metadata: %w", err)
	}
	return &metadata, nil
}

// UpdateMetadata 更新品牌元数据
func (r *BrandRepository) UpdateMetadata(ctx context.Context, tenantID, brandID int64, req *MetadataUpdateRequest) (*model.BrandMetadata, error) {
	metadata, err := r.GetMetadata(ctx, tenantID, brandID)
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = &model.BrandMetadata{
			BrandID:    brandID,
			CreatedAt:  time.Now(),
		}
	}

	// 更新字段
	if req.VIProfile != nil {
		// 将 interface{} 转换为 JSON 字符串
		metadata.VIProfile = fmt.Sprintf("%v", req.VIProfile)
	}
	if req.ToneProfile != nil {
		metadata.ToneProfile = fmt.Sprintf("%v", req.ToneProfile)
	}
	if req.AudienceProfiles != nil {
		metadata.AudienceProfiles = fmt.Sprintf("%v", req.AudienceProfiles)
	}
	if req.CompetitorList != nil {
		metadata.CompetitorList = fmt.Sprintf("%v", req.CompetitorList)
	}
	if req.BrandValues != nil {
		metadata.BrandValues = req.BrandValues
	}
	if req.UniqueSellingPoints != nil {
		metadata.UniqueSellingPoints = req.UniqueSellingPoints
	}
	metadata.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(metadata).Error; err != nil {
		return nil, fmt.Errorf("update metadata: %w", err)
	}

	return metadata, nil
}

// ListGlossary 列出术语表
func (r *BrandRepository) ListGlossary(ctx context.Context, tenantID, brandID int64, filter *GlossaryFilter) ([]*model.GlossaryEntry, int32, error) {
	var entries []*model.GlossaryEntry
	var total int64

	query := r.db.WithContext(ctx).Model(&model.GlossaryEntry{}).Where("brand_id = ?", brandID)

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Keyword != "" {
		query = query.Where("term LIKE ? OR definition LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count glossary entries: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Order("created_at DESC").Find(&entries).Error; err != nil {
		return nil, 0, fmt.Errorf("list glossary entries: %w", err)
	}

	return entries, int32(total), nil
}

// CreateGlossaryEntry 创建术语
func (r *BrandRepository) CreateGlossaryEntry(ctx context.Context, tenantID, brandID int64, req *GlossaryCreateRequest) (*model.GlossaryEntry, error) {
	entry := &model.GlossaryEntry{
		BrandID:     brandID,
		Term:        req.Term,
		Definition:  req.Definition,
		Category:    req.Category,
		Context:     req.Context,
		IsForbidden: req.IsForbidden,
		IsPreferred: req.IsPreferred,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 将 aliases 转换为 JSON 字符串
	if req.Aliases != nil {
		entry.Aliases = req.Aliases
	}

	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return nil, fmt.Errorf("create glossary entry: %w", err)
	}

	return entry, nil
}

// UpdateGlossaryEntry 更新术语
func (r *BrandRepository) UpdateGlossaryEntry(ctx context.Context, tenantID, entryID int64, req *GlossaryUpdateRequest) (*model.GlossaryEntry, error) {
	var entry model.GlossaryEntry
	if err := r.db.WithContext(ctx).First(&entry, entryID).Error; err != nil {
		return nil, fmt.Errorf("glossary entry not found")
	}

	if req.Term != nil {
		entry.Term = *req.Term
	}
	if req.Definition != nil {
		entry.Definition = *req.Definition
	}
	if req.Category != nil {
		entry.Category = *req.Category
	}
	if req.Aliases != nil {
		entry.Aliases = req.Aliases
	}
	if req.Context != nil {
		entry.Context = *req.Context
	}
	if req.IsForbidden != nil {
		entry.IsForbidden = *req.IsForbidden
	}
	if req.IsPreferred != nil {
		entry.IsPreferred = *req.IsPreferred
	}
	entry.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&entry).Error; err != nil {
		return nil, fmt.Errorf("update glossary entry: %w", err)
	}

	return &entry, nil
}

// DeleteGlossaryEntry 删除术语
func (r *BrandRepository) DeleteGlossaryEntry(ctx context.Context, tenantID, entryID int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.GlossaryEntry{}, entryID).Error; err != nil {
		return fmt.Errorf("delete glossary entry: %w", err)
	}
	return nil
}

// BulkImportGlossary 批量导入术语
func (r *BrandRepository) BulkImportGlossary(ctx context.Context, tenantID, brandID int64, entries []interface{}, overwrite bool) (map[string]int32, error) {
	result := map[string]int32{
		"imported": 0,
		"skipped":  0,
		"errors":   0,
	}

	for _, item := range entries {
		entryMap, ok := item.(map[string]interface{})
		if !ok {
			result["errors"]++
			continue
		}

		term, _ := entryMap["term"].(string)
		definition, _ := entryMap["definition"].(string)

		if term == "" || definition == "" {
			result["errors"]++
			continue
		}

		// 检查是否已存在
		var count int64
		r.db.WithContext(ctx).Model(&model.GlossaryEntry{}).Where("brand_id = ? AND term = ?", brandID, term).Count(&count)

		if count > 0 && !overwrite {
			result["skipped"]++
			continue
		}

		entry := &model.GlossaryEntry{
			BrandID:    brandID,
			Term:       term,
			Definition: definition,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if category, ok := entryMap["category"].(string); ok {
			entry.Category = category
		}
		if context, ok := entryMap["context"].(string); ok {
			entry.Context = context
		}
		if isForbidden, ok := entryMap["is_forbidden"].(bool); ok {
			entry.IsForbidden = isForbidden
		}

		if err := r.db.WithContext(ctx).Save(entry).Error; err != nil {
			result["errors"]++
			continue
		}

		result["imported"]++
	}

	return result, nil
}

// ListSnapshots 列出品牌快照
func (r *BrandRepository) ListSnapshots(ctx context.Context, tenantID, brandID int64, page, pageSize int) ([]*model.BrandSnapshot, int32, error) {
	var snapshots []*model.BrandSnapshot
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BrandSnapshot{}).Where("brand_id = ?", brandID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count snapshots: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&snapshots).Error; err != nil {
		return nil, 0, fmt.Errorf("list snapshots: %w", err)
	}

	return snapshots, int32(total), nil
}

// CreateSnapshot 创建品牌快照
func (r *BrandRepository) CreateSnapshot(ctx context.Context, tenantID, brandID int64, req *SnapshotCreateRequest) (*model.BrandSnapshot, error) {
	// 获取品牌信息
	brand, err := r.FindByID(ctx, tenantID, brandID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, fmt.Errorf("brand not found")
	}

	// 获取元数据
	metadata, _ := r.GetMetadata(ctx, tenantID, brandID)

	// 构建快照数据
	snapshotData := fmt.Sprintf(`{"brand": %v, "metadata": %v}`, brand, metadata)

	snapshot := &model.BrandSnapshot{
		BrandID:      brandID,
		Version:      req.Version,
		SnapshotData: snapshotData,
		ChangeLog:    req.ChangeLog,
		CreatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(snapshot).Error; err != nil {
		return nil, fmt.Errorf("create snapshot: %w", err)
	}

	return snapshot, nil
}

// 类型定义

type BrandFilter struct {
	Keyword  string
	Industry string
	Status   int
	Page     int
	PageSize int
}

type BrandCreateRequest struct {
	Name         string
	Slug         string
	Description  string
	LogoURL      string
	Website      string
	Industry     string
	FoundedYear  int32
	Headquarters string
}

type BrandUpdateRequest struct {
	Name         *string
	Description  *string
	LogoURL      *string
	Website      *string
	Industry     *string
	FoundedYear  *int32
	Headquarters *string
	Status       *int32
}

type MetadataUpdateRequest struct {
	VIProfile            interface{}
	ToneProfile          interface{}
	AudienceProfiles     interface{}
	CompetitorList       interface{}
	BrandValues          []string
	UniqueSellingPoints  []string
}

type GlossaryFilter struct {
	Category string
	Keyword  string
	Page     int
	PageSize int
}

type GlossaryCreateRequest struct {
	Term        string
	Definition  string
	Category    string
	Aliases     []string
	Context     string
	IsForbidden bool
	IsPreferred bool
}

type GlossaryUpdateRequest struct {
	Term        *string
	Definition  *string
	Category    *string
	Aliases     []string
	Context     *string
	IsForbidden *bool
	IsPreferred *bool
}

type SnapshotCreateRequest struct {
	Version   string
	ChangeLog string
}
