package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"opengeo/service/content/internal/domain/model"
	"opengeo/service/content/internal/port/outbound"
)

// ContentMySQLRepository MySQL内容仓储实现（出站适配器）
type ContentMySQLRepository struct {
	db *gorm.DB
}

// NewContentMySQLRepository 创建MySQL内容仓储
func NewContentMySQLRepository(db *gorm.DB) *ContentMySQLRepository {
	return &ContentMySQLRepository{db: db}
}

// Create 创建内容
func (r *ContentMySQLRepository) Create(ctx context.Context, content *model.Content) error {
	if err := r.db.WithContext(ctx).Create(content).Error; err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}
	return nil
}

// GetByID 根据ID获取内容
func (r *ContentMySQLRepository) GetByID(ctx context.Context, id int64) (*model.Content, error) {
	var content model.Content
	if err := r.db.WithContext(ctx).First(&content, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("content not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}
	return &content, nil
}

// Update 更新内容
func (r *ContentMySQLRepository) Update(ctx context.Context, content *model.Content) error {
	if err := r.db.WithContext(ctx).Save(content).Error; err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}
	return nil
}

// Delete 删除内容
func (r *ContentMySQLRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Content{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}
	return nil
}

// List 列出内容
func (r *ContentMySQLRepository) List(ctx context.Context, filter *outbound.ContentFilter) ([]*model.Content, int32, error) {
	var contents []*model.Content
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Content{})

	// 应用过滤条件
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Status >= 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ContentType != "" {
		query = query.Where("content_type = ?", filter.ContentType)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count contents: %w", err)
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(int(offset)).Limit(int(filter.PageSize)).Order("created_at DESC").Find(&contents).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list contents: %w", err)
	}

	return contents, int32(total), nil
}

// CreateVersion 创建内容版本
func (r *ContentMySQLRepository) CreateVersion(ctx context.Context, version *model.ContentVersion) error {
	if err := r.db.WithContext(ctx).Create(version).Error; err != nil {
		return fmt.Errorf("failed to create version: %w", err)
	}
	return nil
}

// GetVersionByID 根据ID获取内容版本
func (r *ContentMySQLRepository) GetVersionByID(ctx context.Context, id int64) (*model.ContentVersion, error) {
	var version model.ContentVersion
	if err := r.db.WithContext(ctx).First(&version, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("version not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	return &version, nil
}

// ListVersions 列出内容版本
func (r *ContentMySQLRepository) ListVersions(ctx context.Context, contentID int64, page, pageSize int32) ([]*model.ContentVersion, int32, error) {
	var versions []*model.ContentVersion
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ContentVersion{}).Where("content_id = ?", contentID)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count versions: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("version DESC").Find(&versions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list versions: %w", err)
	}

	return versions, int32(total), nil
}

// ContentTemplateMySQLRepository MySQL内容模板仓储实现（出站适配器）
type ContentTemplateMySQLRepository struct {
	db *gorm.DB
}

// NewContentTemplateMySQLRepository 创建MySQL内容模板仓储
func NewContentTemplateMySQLRepository(db *gorm.DB) *ContentTemplateMySQLRepository {
	return &ContentTemplateMySQLRepository{db: db}
}

// Create 创建模板
func (r *ContentTemplateMySQLRepository) Create(ctx context.Context, template *model.ContentTemplate) error {
	if err := r.db.WithContext(ctx).Create(template).Error; err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}
	return nil
}

// GetByID 根据ID获取模板
func (r *ContentTemplateMySQLRepository) GetByID(ctx context.Context, id int64) (*model.ContentTemplate, error) {
	var template model.ContentTemplate
	if err := r.db.WithContext(ctx).First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("template not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return &template, nil
}

// List 列出模板
func (r *ContentTemplateMySQLRepository) List(ctx context.Context, filter *outbound.TemplateFilter) ([]*model.ContentTemplate, int32, error) {
	var templates []*model.ContentTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ContentTemplate{})

	// 应用过滤条件
	if filter.UserID > 0 {
		query = query.Where("user_id = ? OR is_public = ?", filter.UserID, true)
	}
	if filter.TemplateType != "" {
		query = query.Where("template_type = ?", filter.TemplateType)
	}
	if filter.IsPublic {
		query = query.Where("is_public = ?", true)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count templates: %w", err)
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(int(offset)).Limit(int(filter.PageSize)).Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list templates: %w", err)
	}

	return templates, int32(total), nil
}

// Update 更新模板
func (r *ContentTemplateMySQLRepository) Update(ctx context.Context, template *model.ContentTemplate) error {
	if err := r.db.WithContext(ctx).Save(template).Error; err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	return nil
}

// Delete 删除模板
func (r *ContentTemplateMySQLRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.ContentTemplate{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	return nil
}

// KnowledgeEntityMySQLRepository MySQL知识实体仓储实现（出站适配器）
type KnowledgeEntityMySQLRepository struct {
	db *gorm.DB
}

// NewKnowledgeEntityMySQLRepository 创建MySQL知识实体仓储
func NewKnowledgeEntityMySQLRepository(db *gorm.DB) *KnowledgeEntityMySQLRepository {
	return &KnowledgeEntityMySQLRepository{db: db}
}

// Create 创建实体
func (r *KnowledgeEntityMySQLRepository) Create(ctx context.Context, entity *model.KnowledgeEntity) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}
	return nil
}

// GetByID 根据ID获取实体
func (r *KnowledgeEntityMySQLRepository) GetByID(ctx context.Context, id int64) (*model.KnowledgeEntity, error) {
	var entity model.KnowledgeEntity
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("entity not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}
	return &entity, nil
}

// List 列出实体
func (r *KnowledgeEntityMySQLRepository) List(ctx context.Context, filter *outbound.EntityFilter) ([]*model.KnowledgeEntity, int32, error) {
	var entities []*model.KnowledgeEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&model.KnowledgeEntity{})

	// 应用过滤条件
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count entities: %w", err)
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(int(offset)).Limit(int(filter.PageSize)).Order("created_at DESC").Find(&entities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list entities: %w", err)
	}

	return entities, int32(total), nil
}

// Update 更新实体
func (r *KnowledgeEntityMySQLRepository) Update(ctx context.Context, entity *model.KnowledgeEntity) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	return nil
}

// Delete 删除实体
func (r *KnowledgeEntityMySQLRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.KnowledgeEntity{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}
	return nil
}

// LinkToContent 关联实体到内容
func (r *KnowledgeEntityMySQLRepository) LinkToContent(ctx context.Context, contentID, entityID int64) error {
	// 创建关联记录
	sql := "INSERT INTO content_entities (content_id, entity_id) VALUES (?, ?)"
	if err := r.db.WithContext(ctx).Exec(sql, contentID, entityID).Error; err != nil {
		return fmt.Errorf("failed to link entity to content: %w", err)
	}
	return nil
}

// GetContentEntities 获取内容关联的实体
func (r *KnowledgeEntityMySQLRepository) GetContentEntities(ctx context.Context, contentID int64) ([]*model.KnowledgeEntity, error) {
	var entities []*model.KnowledgeEntity

	err := r.db.WithContext(ctx).
		Joins("JOIN content_entities ON knowledge_entities.id = content_entities.entity_id").
		Where("content_entities.content_id = ?", contentID).
		Find(&entities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get content entities: %w", err)
	}

	return entities, nil
}