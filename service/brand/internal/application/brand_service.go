package application

import (
	"context"
	"encoding/json"
	"time"

	"opengeo/service/brand/internal/domain"
	"opengeo/service/brand/internal/port"
)

// BrandService 品牌应用服务
type BrandService struct {
	brandRepo    port.BrandRepository
	metadataRepo port.BrandMetadataRepository
	glossaryRepo port.GlossaryRepository
	snapshotRepo port.SnapshotRepository
}

// NewBrandService 创建品牌服务
func NewBrandService(
	brandRepo port.BrandRepository,
	metadataRepo port.BrandMetadataRepository,
	glossaryRepo port.GlossaryRepository,
	snapshotRepo port.SnapshotRepository,
) *BrandService {
	return &BrandService{
		brandRepo:    brandRepo,
		metadataRepo: metadataRepo,
		glossaryRepo: glossaryRepo,
		snapshotRepo: snapshotRepo,
	}
}

// ========== 品牌管理 ==========

// CreateBrandRequest 创建品牌请求
type CreateBrandRequest struct {
	TenantID     int64
	Name         string
	Slug         string
	Description  string
	LogoURL      string
	Website      string
	Industry     string
	FoundedYear  int32
	Headquarters string
}

// CreateBrand 创建品牌
func (s *BrandService) CreateBrand(ctx context.Context, req *CreateBrandRequest) (*domain.Brand, error) {
	// 检查标识是否已存在
	exists, err := s.brandRepo.ExistsBySlug(ctx, req.TenantID, req.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrBrandSlugExists
	}

	// 创建品牌实体
	brand := &domain.Brand{
		TenantID:     req.TenantID,
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		LogoURL:      req.LogoURL,
		Website:      req.Website,
		Industry:     req.Industry,
		FoundedYear:  req.FoundedYear,
		Headquarters: req.Headquarters,
		Status:       domain.BrandStatusActive,
		Settings:     make(map[string]string),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 验证实体
	if err := brand.Validate(); err != nil {
		return nil, err
	}

	// 保存品牌
	if err := s.brandRepo.Save(ctx, brand); err != nil {
		return nil, err
	}

	// 创建默认元数据
	metadata := &domain.BrandMetadata{
		BrandID: brand.ID,
		VIProfile: domain.VIProfile{},
		ToneProfile: domain.ToneProfile{},
		AudienceProfiles: []domain.AudienceProfile{},
		BrandValues: []string{},
		UniqueSellingPoints: []string{},
		SchemaVersion: "1.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = s.metadataRepo.Save(ctx, metadata)

	return brand, nil
}

// GetBrand 获取品牌信息
func (s *BrandService) GetBrand(ctx context.Context, tenantID, id int64) (*domain.Brand, error) {
	brand, err := s.brandRepo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, domain.ErrBrandNotFound
	}
	return brand, nil
}

// UpdateBrandRequest 更新品牌请求
type UpdateBrandRequest struct {
	TenantID     int64
	ID           int64
	Name         *string
	Description  *string
	LogoURL      *string
	Website      *string
	Industry     *string
	FoundedYear  *int32
	Headquarters *string
	Status       *domain.BrandStatus
	Settings     map[string]string
}

// UpdateBrand 更新品牌信息
func (s *BrandService) UpdateBrand(ctx context.Context, req *UpdateBrandRequest) (*domain.Brand, error) {
	brand, err := s.brandRepo.FindByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, domain.ErrBrandNotFound
	}

	// 更新字段
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
	if req.Settings != nil {
		brand.Settings = req.Settings
	}
	brand.UpdatedAt = time.Now()

	// 验证实体
	if err := brand.Validate(); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.brandRepo.Save(ctx, brand); err != nil {
		return nil, err
	}

	return brand, nil
}

// DeleteBrand 删除品牌
func (s *BrandService) DeleteBrand(ctx context.Context, tenantID, id int64) error {
	brand, err := s.brandRepo.FindByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if brand == nil {
		return domain.ErrBrandNotFound
	}

	return s.brandRepo.Delete(ctx, tenantID, id)
}

// ListBrandsRequest 列出品牌请求
type ListBrandsRequest struct {
	TenantID int64
	Status   *domain.BrandStatus
	Industry string
	Keyword  string
	Page     int32
	PageSize int32
}

// ListBrandsResponse 列出品牌响应
type ListBrandsResponse struct {
	Brands     []*domain.Brand
	Total      int32
	Page       int32
	PageSize   int32
	TotalPages int32
}

// ListBrands 列出品牌
func (s *BrandService) ListBrands(ctx context.Context, req *ListBrandsRequest) (*ListBrandsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	filter := &port.BrandFilter{
		Status:   req.Status,
		Industry: req.Industry,
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	brands, total, err := s.brandRepo.List(ctx, req.TenantID, filter)
	if err != nil {
		return nil, err
	}

	totalPages := (total + int32(req.PageSize) - 1) / int32(req.PageSize)

	return &ListBrandsResponse{
		Brands:     brands,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ========== 品牌元数据 ==========

// GetBrandMetadata 获取品牌元数据
func (s *BrandService) GetBrandMetadata(ctx context.Context, tenantID, brandID int64) (*domain.BrandMetadata, error) {
	metadata, err := s.metadataRepo.FindByBrandID(ctx, tenantID, brandID)
	if err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, domain.ErrBrandNotFound
	}
	return metadata, nil
}

// UpdateBrandMetadataRequest 更新品牌元数据请求
type UpdateBrandMetadataRequest struct {
	TenantID             int64
	BrandID              int64
	VIProfile            *domain.VIProfile
	ToneProfile          *domain.ToneProfile
	AudienceProfiles     []domain.AudienceProfile
	CompetitorList       []domain.CompetitorInfo
	BrandValues          []string
	UniqueSellingPoints  []string
}

// UpdateBrandMetadata 更新品牌元数据
func (s *BrandService) UpdateBrandMetadata(ctx context.Context, req *UpdateBrandMetadataRequest) (*domain.BrandMetadata, error) {
	metadata, err := s.metadataRepo.FindByBrandID(ctx, req.TenantID, req.BrandID)
	if err != nil {
		return nil, err
	}
	if metadata == nil {
		// 创建新元数据
		metadata = &domain.BrandMetadata{
			BrandID: req.BrandID,
			SchemaVersion: "1.0",
			CreatedAt: time.Now(),
		}
	}

	// 更新字段
	if req.VIProfile != nil {
		metadata.VIProfile = *req.VIProfile
	}
	if req.ToneProfile != nil {
		metadata.ToneProfile = *req.ToneProfile
	}
	if req.AudienceProfiles != nil {
		metadata.AudienceProfiles = req.AudienceProfiles
	}
	if req.CompetitorList != nil {
		metadata.CompetitorList = req.CompetitorList
	}
	if req.BrandValues != nil {
		metadata.BrandValues = req.BrandValues
	}
	if req.UniqueSellingPoints != nil {
		metadata.UniqueSellingPoints = req.UniqueSellingPoints
	}
	metadata.UpdatedAt = time.Now()

	// 保存元数据
	if err := s.metadataRepo.Save(ctx, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

// ========== 品牌术语表 ==========

// CreateGlossaryEntryRequest 创建术语请求
type CreateGlossaryEntryRequest struct {
	TenantID   int64
	BrandID    int64
	Term       string
	Definition string
	Category   string
	Aliases    []string
	Context    string
	IsForbidden bool
}

// CreateGlossaryEntry 创建术语
func (s *BrandService) CreateGlossaryEntry(ctx context.Context, req *CreateGlossaryEntryRequest) (*domain.GlossaryEntry, error) {
	entry := &domain.GlossaryEntry{
		BrandID:     req.BrandID,
		Term:        req.Term,
		Definition:  req.Definition,
		Category:    req.Category,
		Aliases:     req.Aliases,
		Context:     req.Context,
		IsForbidden: req.IsForbidden,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.glossaryRepo.Save(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// GetGlossaryEntry 获取术语
func (s *BrandService) GetGlossaryEntry(ctx context.Context, tenantID, id int64) (*domain.GlossaryEntry, error) {
	entry, err := s.glossaryRepo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, domain.ErrGlossaryNotFound
	}
	return entry, nil
}

// UpdateGlossaryEntryRequest 更新术语请求
type UpdateGlossaryEntryRequest struct {
	TenantID    int64
	ID          int64
	Term        *string
	Definition  *string
	Category    *string
	Aliases     []string
	Context     *string
	IsForbidden *bool
	IsPreferred *bool
}

// UpdateGlossaryEntry 更新术语
func (s *BrandService) UpdateGlossaryEntry(ctx context.Context, req *UpdateGlossaryEntryRequest) (*domain.GlossaryEntry, error) {
	entry, err := s.glossaryRepo.FindByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, domain.ErrGlossaryNotFound
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

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.glossaryRepo.Save(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// DeleteGlossaryEntry 删除术语
func (s *BrandService) DeleteGlossaryEntry(ctx context.Context, tenantID, id int64) error {
	entry, err := s.glossaryRepo.FindByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if entry == nil {
		return domain.ErrGlossaryNotFound
	}

	return s.glossaryRepo.Delete(ctx, tenantID, id)
}

// ListGlossaryEntriesRequest 列出术语请求
type ListGlossaryEntriesRequest struct {
	TenantID    int64
	BrandID     int64
	Category    string
	Keyword     string
	IsForbidden *bool
	Page        int32
	PageSize    int32
}

// ListGlossaryEntriesResponse 列出术语响应
type ListGlossaryEntriesResponse struct {
	Entries    []*domain.GlossaryEntry
	Total      int32
	Page       int32
	PageSize   int32
	TotalPages int32
}

// ListGlossaryEntries 列出术语
func (s *BrandService) ListGlossaryEntries(ctx context.Context, req *ListGlossaryEntriesRequest) (*ListGlossaryEntriesResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	filter := &port.GlossaryFilter{
		Category:    req.Category,
		Keyword:     req.Keyword,
		IsForbidden: req.IsForbidden,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	entries, total, err := s.glossaryRepo.List(ctx, req.TenantID, req.BrandID, filter)
	if err != nil {
		return nil, err
	}

	totalPages := (total + int32(req.PageSize) - 1) / int32(req.PageSize)

	return &ListGlossaryEntriesResponse{
		Entries:    entries,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// BulkImportGlossaryRequest 批量导入术语请求
type BulkImportGlossaryRequest struct {
	TenantID         int64
	BrandID          int64
	Entries          []*domain.GlossaryEntry
	OverwriteExisting bool
}

// BulkImportGlossaryResponse 批量导入术语响应
type BulkImportGlossaryResponse struct {
	ImportedCount int32
	SkippedCount  int32
	ErrorCount    int32
}

// BulkImportGlossary 批量导入术语
func (s *BrandService) BulkImportGlossary(ctx context.Context, req *BulkImportGlossaryRequest) (*BulkImportGlossaryResponse, error) {
	// 设置品牌 ID
	for _, entry := range req.Entries {
		entry.BrandID = req.BrandID
		entry.CreatedAt = time.Now()
		entry.UpdatedAt = time.Now()
	}

	imported, skipped, err := s.glossaryRepo.BulkCreate(ctx, req.Entries)
	if err != nil {
		return nil, err
	}

	return &BulkImportGlossaryResponse{
		ImportedCount: imported,
		SkippedCount:  skipped,
	}, nil
}

// ========== 品牌快照 ==========

// CreateBrandSnapshotRequest 创建品牌快照请求
type CreateBrandSnapshotRequest struct {
	TenantID  int64
	BrandID   int64
	Version   string
	ChangeLog string
	CreatedBy int64
}

// CreateBrandSnapshot 创建品牌快照
func (s *BrandService) CreateBrandSnapshot(ctx context.Context, req *CreateBrandSnapshotRequest) (*domain.BrandSnapshot, error) {
	// 获取品牌信息
	brand, err := s.brandRepo.FindByID(ctx, req.TenantID, req.BrandID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, domain.ErrBrandNotFound
	}

	// 获取品牌元数据
	metadata, _ := s.metadataRepo.FindByBrandID(ctx, req.TenantID, req.BrandID)

	// 构建快照数据
	snapshotData := map[string]interface{}{
		"brand":    brand,
		"metadata": metadata,
	}
	dataBytes, err := json.Marshal(snapshotData)
	if err != nil {
		return nil, err
	}

	snapshot := &domain.BrandSnapshot{
		BrandID:      req.BrandID,
		Version:      req.Version,
		SnapshotData: string(dataBytes),
		ChangeLog:    req.ChangeLog,
		CreatedBy:    req.CreatedBy,
		CreatedAt:    time.Now(),
	}

	if err := s.snapshotRepo.Save(ctx, snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

// ListBrandSnapshotsRequest 列出品牌快照请求
type ListBrandSnapshotsRequest struct {
	TenantID int64
	BrandID  int64
	Page     int32
	PageSize int32
}

// ListBrandSnapshotsResponse 列出品牌快照响应
type ListBrandSnapshotsResponse struct {
	Snapshots  []*domain.BrandSnapshot
	Total      int32
	Page       int32
	PageSize   int32
	TotalPages int32
}

// ListBrandSnapshots 列出品牌快照
func (s *BrandService) ListBrandSnapshots(ctx context.Context, req *ListBrandSnapshotsRequest) (*ListBrandSnapshotsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	snapshots, total, err := s.snapshotRepo.List(ctx, req.TenantID, req.BrandID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	totalPages := (total + int32(req.PageSize) - 1) / int32(req.PageSize)

	return &ListBrandSnapshotsResponse{
		Snapshots:  snapshots,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
