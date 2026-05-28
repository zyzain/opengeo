package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// ListBrands 列出品牌
func (h *Handler) ListBrands(ctx context.Context, c *app.RequestContext) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	industry := c.Query("industry")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	// 调用品牌服务
	brands, total, err := h.brandRepo.List(ctx, 1, &BrandFilter{ // TODO: 从 JWT 获取 tenant_id
		Keyword:  keyword,
		Industry: industry,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to list brands",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"brands":    brands,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetBrand 获取品牌
func (h *Handler) GetBrand(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	brand, err := h.brandRepo.FindByID(ctx, 1, id) // TODO: 从 JWT 获取 tenant_id
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to get brand",
		})
		return
	}

	if brand == nil {
		c.JSON(consts.StatusNotFound, utils.H{
			"error": "Brand not found",
		})
		return
	}

	c.JSON(consts.StatusOK, brand)
}

// CreateBrand 创建品牌
func (h *Handler) CreateBrand(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Name         string `json:"name"`
		Slug         string `json:"slug"`
		Description  string `json:"description"`
		LogoURL      string `json:"logo_url"`
		Website      string `json:"website"`
		Industry     string `json:"industry"`
		FoundedYear  int32  `json:"founded_year"`
		Headquarters string `json:"headquarters"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid request body",
		})
		return
	}

	// 验证必填字段
	if req.Name == "" || req.Slug == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Name and slug are required",
		})
		return
	}

	// 创建品牌
	brand, err := h.brandRepo.Create(ctx, 1, &BrandCreateRequest{ // TODO: 从 JWT 获取 tenant_id
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  req.Description,
		LogoURL:      req.LogoURL,
		Website:      req.Website,
		Industry:     req.Industry,
		FoundedYear:  req.FoundedYear,
		Headquarters: req.Headquarters,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to create brand",
		})
		return
	}

	c.JSON(consts.StatusCreated, brand)
}

// UpdateBrand 更新品牌
func (h *Handler) UpdateBrand(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	var req struct {
		Name         *string `json:"name"`
		Description  *string `json:"description"`
		LogoURL      *string `json:"logo_url"`
		Website      *string `json:"website"`
		Industry     *string `json:"industry"`
		FoundedYear  *int32  `json:"founded_year"`
		Headquarters *string `json:"headquarters"`
		Status       *int32  `json:"status"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid request body",
		})
		return
	}

	// 更新品牌
	brand, err := h.brandRepo.Update(ctx, 1, id, &BrandUpdateRequest{ // TODO: 从 JWT 获取 tenant_id
		Name:         req.Name,
		Description:  req.Description,
		LogoURL:      req.LogoURL,
		Website:      req.Website,
		Industry:     req.Industry,
		FoundedYear:  req.FoundedYear,
		Headquarters: req.Headquarters,
		Status:       req.Status,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to update brand",
		})
		return
	}

	c.JSON(consts.StatusOK, brand)
}

// DeleteBrand 删除品牌
func (h *Handler) DeleteBrand(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	err = h.brandRepo.Delete(ctx, 1, id) // TODO: 从 JWT 获取 tenant_id
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to delete brand",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"success": true,
	})
}

// GetBrandMetadata 获取品牌元数据
func (h *Handler) GetBrandMetadata(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	metadata, err := h.brandRepo.GetMetadata(ctx, 1, id) // TODO: 从 JWT 获取 tenant_id
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to get brand metadata",
		})
		return
	}

	c.JSON(consts.StatusOK, metadata)
}

// UpdateBrandMetadata 更新品牌元数据
func (h *Handler) UpdateBrandMetadata(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	var req struct {
		VIProfile            interface{} `json:"vi_profile"`
		ToneProfile          interface{} `json:"tone_profile"`
		AudienceProfiles     interface{} `json:"audience_profiles"`
		CompetitorList       interface{} `json:"competitor_list"`
		BrandValues          []string    `json:"brand_values"`
		UniqueSellingPoints  []string    `json:"unique_selling_points"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid request body",
		})
		return
	}

	metadata, err := h.brandRepo.UpdateMetadata(ctx, 1, id, &MetadataUpdateRequest{ // TODO: 从 JWT 获取 tenant_id
		VIProfile:            req.VIProfile,
		ToneProfile:          req.ToneProfile,
		AudienceProfiles:     req.AudienceProfiles,
		CompetitorList:       req.CompetitorList,
		BrandValues:          req.BrandValues,
		UniqueSellingPoints:  req.UniqueSellingPoints,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to update brand metadata",
		})
		return
	}

	c.JSON(consts.StatusOK, metadata)
}

// ListGlossaryEntries 列出术语表
func (h *Handler) ListGlossaryEntries(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	entries, total, err := h.brandRepo.ListGlossary(ctx, 1, id, &GlossaryFilter{ // TODO: 从 JWT 获取 tenant_id
		Category: category,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to list glossary entries",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"entries": entries,
		"total":   total,
		"page":    page,
		"page_size": pageSize,
	})
}

// CreateGlossaryEntry 创建术语
func (h *Handler) CreateGlossaryEntry(ctx context.Context, c *app.RequestContext) {
	brandID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	var req struct {
		Term        string   `json:"term"`
		Definition  string   `json:"definition"`
		Category    string   `json:"category"`
		Aliases     []string `json:"aliases"`
		Context     string   `json:"context"`
		IsForbidden bool     `json:"is_forbidden"`
		IsPreferred bool     `json:"is_preferred"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid request body",
		})
		return
	}

	// 验证必填字段
	if req.Term == "" || req.Definition == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Term and definition are required",
		})
		return
	}

	entry, err := h.brandRepo.CreateGlossaryEntry(ctx, 1, brandID, &GlossaryCreateRequest{ // TODO: 从 JWT 获取 tenant_id
		Term:        req.Term,
		Definition:  req.Definition,
		Category:    req.Category,
		Aliases:     req.Aliases,
		Context:     req.Context,
		IsForbidden: req.IsForbidden,
		IsPreferred: req.IsPreferred,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to create glossary entry",
		})
		return
	}

	c.JSON(consts.StatusCreated, entry)
}

// UpdateGlossaryEntry 更新术语
func (h *Handler) UpdateGlossaryEntry(ctx context.Context, c *app.RequestContext) {
	entryID, err := strconv.ParseInt(c.Param("entry_id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid entry ID",
		})
		return
	}

	var req struct {
		Term        *string  `json:"term"`
		Definition  *string  `json:"definition"`
		Category    *string  `json:"category"`
		Aliases     []string `json:"aliases"`
		Context     *string  `json:"context"`
		IsForbidden *bool    `json:"is_forbidden"`
		IsPreferred *bool    `json:"is_preferred"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid request body",
		})
		return
	}

	entry, err := h.brandRepo.UpdateGlossaryEntry(ctx, 1, entryID, &GlossaryUpdateRequest{ // TODO: 从 JWT 获取 tenant_id
		Term:        req.Term,
		Definition:  req.Definition,
		Category:    req.Category,
		Aliases:     req.Aliases,
		Context:     req.Context,
		IsForbidden: req.IsForbidden,
		IsPreferred: req.IsPreferred,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to update glossary entry",
		})
		return
	}

	c.JSON(consts.StatusOK, entry)
}

// DeleteGlossaryEntry 删除术语
func (h *Handler) DeleteGlossaryEntry(ctx context.Context, c *app.RequestContext) {
	entryID, err := strconv.ParseInt(c.Param("entry_id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid entry ID",
		})
		return
	}

	err = h.brandRepo.DeleteGlossaryEntry(ctx, 1, entryID) // TODO: 从 JWT 获取 tenant_id
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to delete glossary entry",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"success": true,
	})
}

// BulkImportGlossary 批量导入术语
func (h *Handler) BulkImportGlossary(ctx context.Context, c *app.RequestContext) {
	brandID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	var req struct {
		Entries []struct {
			Term        string   `json:"term"`
			Definition  string   `json:"definition"`
			Category    string   `json:"category"`
			Aliases     []string `json:"aliases"`
			Context     string   `json:"context"`
			IsForbidden bool     `json:"is_forbidden"`
		} `json:"entries"`
		OverwriteExisting bool `json:"overwrite_existing"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid request body",
		})
		return
	}

	result, err := h.brandRepo.BulkImportGlossary(ctx, 1, brandID, req.Entries, req.OverwriteExisting) // TODO: 从 JWT 获取 tenant_id
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to bulk import glossary",
		})
		return
	}

	c.JSON(consts.StatusOK, result)
}

// ListBrandSnapshots 列出品牌快照
func (h *Handler) ListBrandSnapshots(ctx context.Context, c *app.RequestContext) {
	brandID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	snapshots, total, err := h.brandRepo.ListSnapshots(ctx, 1, brandID, page, pageSize) // TODO: 从 JWT 获取 tenant_id
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to list brand snapshots",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"snapshots": snapshots,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateBrandSnapshot 创建品牌快照
func (h *Handler) CreateBrandSnapshot(ctx context.Context, c *app.RequestContext) {
	brandID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid brand ID",
		})
		return
	}

	var req struct {
		Version   string `json:"version"`
		ChangeLog string `json:"change_log"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Invalid request body",
		})
		return
	}

	snapshot, err := h.brandRepo.CreateSnapshot(ctx, 1, brandID, &SnapshotCreateRequest{ // TODO: 从 JWT 获取 tenant_id
		Version:   req.Version,
		ChangeLog: req.ChangeLog,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "Failed to create brand snapshot",
		})
		return
	}

	c.JSON(consts.StatusCreated, snapshot)
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
