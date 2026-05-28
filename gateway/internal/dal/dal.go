package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/gateway/internal/model"
)

// BrowserFingerprintRepository 浏览器指纹仓储
type BrowserFingerprintRepository struct {
	db *gorm.DB
}

// NewBrowserFingerprintRepository 创建浏览器指纹仓储
func NewBrowserFingerprintRepository(db *gorm.DB) *BrowserFingerprintRepository {
	return &BrowserFingerprintRepository{db: db}
}

// Create 创建指纹
func (r *BrowserFingerprintRepository) Create(ctx context.Context, fp *model.BrowserFingerprint) error {
	fp.CreatedAt = time.Now()
	fp.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Create(fp).Error; err != nil {
		return fmt.Errorf("create fingerprint: %w", err)
	}
	return nil
}

// GetByID 根据ID获取指纹
func (r *BrowserFingerprintRepository) GetByID(ctx context.Context, id int64) (*model.BrowserFingerprint, error) {
	var fp model.BrowserFingerprint
	if err := r.db.WithContext(ctx).First(&fp, id).Error; err != nil {
		return nil, fmt.Errorf("fingerprint not found")
	}
	return &fp, nil
}

// Update 更新指纹
func (r *BrowserFingerprintRepository) Update(ctx context.Context, fp *model.BrowserFingerprint) error {
	fp.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(fp).Error; err != nil {
		return fmt.Errorf("update fingerprint: %w", err)
	}
	return nil
}

// Delete 删除指纹
func (r *BrowserFingerprintRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.BrowserFingerprint{}, id).Error; err != nil {
		return fmt.Errorf("delete fingerprint: %w", err)
	}
	return nil
}

// List 列出指纹
func (r *BrowserFingerprintRepository) List(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.BrowserFingerprint, int32, error) {
	var fps []*model.BrowserFingerprint
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BrowserFingerprint{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count fingerprints: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&fps).Error; err != nil {
		return nil, 0, fmt.Errorf("list fingerprints: %w", err)
	}

	return fps, int32(total), nil
}

// ToggleStatus 切换状态
func (r *BrowserFingerprintRepository) ToggleStatus(ctx context.Context, id int64) (*model.BrowserFingerprint, error) {
	fp, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if fp.Status == "active" {
		fp.Status = "inactive"
	} else {
		fp.Status = "active"
	}
	fp.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(fp).Error; err != nil {
		return nil, fmt.Errorf("toggle fingerprint status: %w", err)
	}

	return fp, nil
}

// ProxyIPRepository 代理IP仓储
type ProxyIPRepository struct {
	db *gorm.DB
}

// NewProxyIPRepository 创建代理IP仓储
func NewProxyIPRepository(db *gorm.DB) *ProxyIPRepository {
	return &ProxyIPRepository{db: db}
}

// Create 创建代理IP
func (r *ProxyIPRepository) Create(ctx context.Context, proxy *model.ProxyIP) error {
	proxy.CreatedAt = time.Now()
	proxy.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Create(proxy).Error; err != nil {
		return fmt.Errorf("create proxy: %w", err)
	}
	return nil
}

// GetByID 根据ID获取代理IP
func (r *ProxyIPRepository) GetByID(ctx context.Context, id int64) (*model.ProxyIP, error) {
	var proxy model.ProxyIP
	if err := r.db.WithContext(ctx).First(&proxy, id).Error; err != nil {
		return nil, fmt.Errorf("proxy not found")
	}
	return &proxy, nil
}

// Update 更新代理IP
func (r *ProxyIPRepository) Update(ctx context.Context, proxy *model.ProxyIP) error {
	proxy.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(proxy).Error; err != nil {
		return fmt.Errorf("update proxy: %w", err)
	}
	return nil
}

// Delete 删除代理IP
func (r *ProxyIPRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.ProxyIP{}, id).Error; err != nil {
		return fmt.Errorf("delete proxy: %w", err)
	}
	return nil
}

// List 列出代理IP
func (r *ProxyIPRepository) List(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.ProxyIP, int32, error) {
	var proxies []*model.ProxyIP
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ProxyIP{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count proxies: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&proxies).Error; err != nil {
		return nil, 0, fmt.Errorf("list proxies: %w", err)
	}

	return proxies, int32(total), nil
}

// ToggleStatus 切换状态
func (r *ProxyIPRepository) ToggleStatus(ctx context.Context, id int64) (*model.ProxyIP, error) {
	proxy, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if proxy.Status == "active" {
		proxy.Status = "inactive"
	} else {
		proxy.Status = "active"
	}
	proxy.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(proxy).Error; err != nil {
		return nil, fmt.Errorf("toggle proxy status: %w", err)
	}

	return proxy, nil
}

// ContentTemplateRepository 内容模板仓储
type ContentTemplateRepository struct {
	db *gorm.DB
}

// NewContentTemplateRepository 创建内容模板仓储
func NewContentTemplateRepository(db *gorm.DB) *ContentTemplateRepository {
	return &ContentTemplateRepository{db: db}
}

// Create 创建模板
func (r *ContentTemplateRepository) Create(ctx context.Context, tpl *model.ContentTemplate) error {
	tpl.CreatedAt = time.Now()
	tpl.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Create(tpl).Error; err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}

// GetByID 根据ID获取模板
func (r *ContentTemplateRepository) GetByID(ctx context.Context, id int64) (*model.ContentTemplate, error) {
	var tpl model.ContentTemplate
	if err := r.db.WithContext(ctx).First(&tpl, id).Error; err != nil {
		return nil, fmt.Errorf("template not found")
	}
	return &tpl, nil
}

// Update 更新模板
func (r *ContentTemplateRepository) Update(ctx context.Context, tpl *model.ContentTemplate) error {
	tpl.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(tpl).Error; err != nil {
		return fmt.Errorf("update template: %w", err)
	}
	return nil
}

// Delete 删除模板
func (r *ContentTemplateRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.ContentTemplate{}, id).Error; err != nil {
		return fmt.Errorf("delete template: %w", err)
	}
	return nil
}

// List 列出模板
func (r *ContentTemplateRepository) List(ctx context.Context, userID int64, templateType string, page, pageSize int) ([]*model.ContentTemplate, int32, error) {
	var tpls []*model.ContentTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ContentTemplate{})
	if userID > 0 {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	}
	if templateType != "" {
		query = query.Where("template_type = ?", templateType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count templates: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tpls).Error; err != nil {
		return nil, 0, fmt.Errorf("list templates: %w", err)
	}

	return tpls, int32(total), nil
}

// IncrementUsage 增加使用次数
func (r *ContentTemplateRepository) IncrementUsage(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.ContentTemplate{}).Where("id = ?", id).
		Update("usage_count", gorm.Expr("usage_count + 1")).Error
}

// StaggerStrategyRepository 错峰策略仓储
type StaggerStrategyRepository struct {
	db *gorm.DB
}

// NewStaggerStrategyRepository 创建错峰策略仓储
func NewStaggerStrategyRepository(db *gorm.DB) *StaggerStrategyRepository {
	return &StaggerStrategyRepository{db: db}
}

// Create 创建策略
func (r *StaggerStrategyRepository) Create(ctx context.Context, strategy *model.StaggerStrategy) error {
	strategy.CreatedAt = time.Now()
	strategy.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Create(strategy).Error; err != nil {
		return fmt.Errorf("create stagger strategy: %w", err)
	}
	return nil
}

// GetByID 根据ID获取策略
func (r *StaggerStrategyRepository) GetByID(ctx context.Context, id int64) (*model.StaggerStrategy, error) {
	var strategy model.StaggerStrategy
	if err := r.db.WithContext(ctx).First(&strategy, id).Error; err != nil {
		return nil, fmt.Errorf("stagger strategy not found")
	}
	return &strategy, nil
}

// Update 更新策略
func (r *StaggerStrategyRepository) Update(ctx context.Context, strategy *model.StaggerStrategy) error {
	strategy.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(strategy).Error; err != nil {
		return fmt.Errorf("update stagger strategy: %w", err)
	}
	return nil
}

// Delete 删除策略
func (r *StaggerStrategyRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.StaggerStrategy{}, id).Error; err != nil {
		return fmt.Errorf("delete stagger strategy: %w", err)
	}
	return nil
}

// List 列出策略
func (r *StaggerStrategyRepository) List(ctx context.Context, userID int64, status string, page, pageSize int) ([]*model.StaggerStrategy, int32, error) {
	var strategies []*model.StaggerStrategy
	var total int64

	query := r.db.WithContext(ctx).Model(&model.StaggerStrategy{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count stagger strategies: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&strategies).Error; err != nil {
		return nil, 0, fmt.Errorf("list stagger strategies: %w", err)
	}

	return strategies, int32(total), nil
}

// ToggleStatus 切换状态
func (r *StaggerStrategyRepository) ToggleStatus(ctx context.Context, id int64) (*model.StaggerStrategy, error) {
	strategy, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strategy.Status == "active" {
		strategy.Status = "inactive"
	} else {
		strategy.Status = "active"
	}
	strategy.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(strategy).Error; err != nil {
		return nil, fmt.Errorf("toggle stagger strategy status: %w", err)
	}

	return strategy, nil
}

// StaggerConfigRepository 错峰配置仓储
type StaggerConfigRepository struct {
	db *gorm.DB
}

// NewStaggerConfigRepository 创建错峰配置仓储
func NewStaggerConfigRepository(db *gorm.DB) *StaggerConfigRepository {
	return &StaggerConfigRepository{db: db}
}

// GetByUserID 根据用户ID获取配置
func (r *StaggerConfigRepository) GetByUserID(ctx context.Context, userID int64) (*model.StaggerConfig, error) {
	var config model.StaggerConfig
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&config).Error; err != nil {
		// 返回默认配置
		return &model.StaggerConfig{
			UserID:           userID,
			MinInterval:      5,
			MaxInterval:      15,
			RandomRange:      30,
			BatchSize:        10,
			CooldownAfter:    50,
			CooldownDuration: 30,
		}, nil
	}
	return &config, nil
}

// Save 保存配置
func (r *StaggerConfigRepository) Save(ctx context.Context, config *model.StaggerConfig) error {
	config.UpdatedAt = time.Now()

	// 尝试更新，如果不存在则创建
	var existing model.StaggerConfig
	if err := r.db.WithContext(ctx).Where("user_id = ?", config.UserID).First(&existing).Error; err != nil {
		config.CreatedAt = time.Now()
		return r.db.WithContext(ctx).Create(config).Error
	}

	return r.db.WithContext(ctx).Model(&existing).Updates(map[string]interface{}{
		"min_interval":       config.MinInterval,
		"max_interval":       config.MaxInterval,
		"random_range":       config.RandomRange,
		"batch_size":         config.BatchSize,
		"cooldown_after":     config.CooldownAfter,
		"cooldown_duration":  config.CooldownDuration,
		"updated_at":         time.Now(),
	}).Error
}
