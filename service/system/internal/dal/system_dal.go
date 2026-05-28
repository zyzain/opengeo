package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/service/system/internal/domain/model"
)

// SystemConfigRepository 系统配置仓储
type SystemConfigRepository struct {
	db *gorm.DB
}

// NewSystemConfigRepository 创建系统配置仓储
func NewSystemConfigRepository(db *gorm.DB) *SystemConfigRepository {
	return &SystemConfigRepository{db: db}
}

// GetByKey 根据Key获取配置
func (r *SystemConfigRepository) GetByKey(ctx context.Context, key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	if err := r.db.WithContext(ctx).Where("config_key = ?", key).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("config not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return &config, nil
}

// Set 设置配置
func (r *SystemConfigRepository) Set(ctx context.Context, key, value, configType, description string, isPublic bool) error {
	config := &model.SystemConfig{
		ConfigKey:   key,
		ConfigValue: value,
		ConfigType:  configType,
		Description: description,
		IsPublic:    isPublic,
		UpdatedAt:   time.Now(),
	}

	if err := r.db.WithContext(ctx).
		Where("config_key = ?", key).
		Assign(map[string]interface{}{
			"config_value": value,
			"config_type":  configType,
			"description":  description,
			"is_public":    isPublic,
			"updated_at":   time.Now(),
		}).
		FirstOrCreate(config).Error; err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}
	return nil
}

// List 列出配置
func (r *SystemConfigRepository) List(ctx context.Context, isPublic *bool, page, pageSize int) ([]*model.SystemConfig, int32, error) {
	var configs []*model.SystemConfig
	var total int64

	query := r.db.WithContext(ctx).Model(&model.SystemConfig{})

	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count configs: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("config_key ASC").Find(&configs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list configs: %w", err)
	}

	return configs, int32(total), nil
}

// Delete 删除配置
func (r *SystemConfigRepository) Delete(ctx context.Context, key string) error {
	if err := r.db.WithContext(ctx).Where("config_key = ?", key).Delete(&model.SystemConfig{}).Error; err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}
	return nil
}

// PluginRepository 插件仓储
type PluginRepository struct {
	db *gorm.DB
}

// NewPluginRepository 创建插件仓储
func NewPluginRepository(db *gorm.DB) *PluginRepository {
	return &PluginRepository{db: db}
}

// Create 创建插件
func (r *PluginRepository) Create(ctx context.Context, plugin *model.Plugin) error {
	if err := r.db.WithContext(ctx).Create(plugin).Error; err != nil {
		return fmt.Errorf("failed to create plugin: %w", err)
	}
	return nil
}

// GetByID 根据ID获取插件
func (r *PluginRepository) GetByID(ctx context.Context, id int64) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := r.db.WithContext(ctx).First(&plugin, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("plugin not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get plugin: %w", err)
	}
	return &plugin, nil
}

// GetByName 根据名称获取插件
func (r *PluginRepository) GetByName(ctx context.Context, name string) (*model.Plugin, error) {
	var plugin model.Plugin
	if err := r.db.WithContext(ctx).Where("plugin_name = ?", name).First(&plugin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("plugin not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get plugin: %w", err)
	}
	return &plugin, nil
}

// List 列出插件
func (r *PluginRepository) List(ctx context.Context, pluginType string, isEnabled *bool, page, pageSize int) ([]*model.Plugin, int32, error) {
	var plugins []*model.Plugin
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Plugin{})

	if pluginType != "" {
		query = query.Where("plugin_type = ?", pluginType)
	}
	if isEnabled != nil {
		query = query.Where("is_enabled = ?", *isEnabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count plugins: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("plugin_name ASC").Find(&plugins).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list plugins: %w", err)
	}

	return plugins, int32(total), nil
}

// Update 更新插件
func (r *PluginRepository) Update(ctx context.Context, plugin *model.Plugin) error {
	if err := r.db.WithContext(ctx).Save(plugin).Error; err != nil {
		return fmt.Errorf("failed to update plugin: %w", err)
	}
	return nil
}

// Delete 删除插件
func (r *PluginRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Plugin{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete plugin: %w", err)
	}
	return nil
}

// WebhookRepository Webhook仓储
type WebhookRepository struct {
	db *gorm.DB
}

// NewWebhookRepository 创建Webhook仓储
func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// Create 创建Webhook
func (r *WebhookRepository) Create(ctx context.Context, webhook *model.Webhook) error {
	if err := r.db.WithContext(ctx).Create(webhook).Error; err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}
	return nil
}

// GetByID 根据ID获取Webhook
func (r *WebhookRepository) GetByID(ctx context.Context, id int64) (*model.Webhook, error) {
	var webhook model.Webhook
	if err := r.db.WithContext(ctx).First(&webhook, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("webhook not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}
	return &webhook, nil
}

// List 列出Webhook
func (r *WebhookRepository) List(ctx context.Context, userID int64, page, pageSize int) ([]*model.Webhook, int32, error) {
	var webhooks []*model.Webhook
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Webhook{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count webhooks: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&webhooks).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list webhooks: %w", err)
	}

	return webhooks, int32(total), nil
}

// Update 更新Webhook
func (r *WebhookRepository) Update(ctx context.Context, webhook *model.Webhook) error {
	if err := r.db.WithContext(ctx).Save(webhook).Error; err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}
	return nil
}

// Delete 删除Webhook
func (r *WebhookRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Webhook{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	return nil
}

// CreateEvent 创建Webhook事件
func (r *WebhookRepository) CreateEvent(ctx context.Context, event *model.WebhookEvent) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}
	return nil
}

// GetEvents 获取Webhook事件
func (r *WebhookRepository) GetEvents(ctx context.Context, webhookID int64, page, pageSize int) ([]*model.WebhookEvent, int32, error) {
	var events []*model.WebhookEvent
	var total int64

	query := r.db.WithContext(ctx).Model(&model.WebhookEvent{}).Where("webhook_id = ?", webhookID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("triggered_at DESC").Find(&events).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}

	return events, int32(total), nil
}

// TenantRepository 租户仓储
type TenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository 创建租户仓储
func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create 创建租户
func (r *TenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	if err := r.db.WithContext(ctx).Create(tenant).Error; err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}

// GetByID 根据ID获取租户
func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	return &tenant, nil
}

// GetByCode 根据Code获取租户
func (r *TenantRepository) GetByCode(ctx context.Context, code string) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.WithContext(ctx).Where("tenant_code = ?", code).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found: %s", code)
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	return &tenant, nil
}

// List 列出租户
func (r *TenantRepository) List(ctx context.Context, page, pageSize int) ([]*model.Tenant, int32, error) {
	var tenants []*model.Tenant
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Tenant{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tenants: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list tenants: %w", err)
	}

	return tenants, int32(total), nil
}

// Update 更新租户
func (r *TenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	if err := r.db.WithContext(ctx).Save(tenant).Error; err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}
	return nil
}

// Delete 删除租户
func (r *TenantRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Tenant{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}

// TranslationRepository 翻译仓储
type TranslationRepository struct {
	db *gorm.DB
}

// NewTranslationRepository 创建翻译仓储
func NewTranslationRepository(db *gorm.DB) *TranslationRepository {
	return &TranslationRepository{db: db}
}

// Get 获取翻译
func (r *TranslationRepository) Get(ctx context.Context, locale, key string) (string, error) {
	var translation model.Translation
	if err := r.db.WithContext(ctx).Where("locale = ? AND key = ?", locale, key).First(&translation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return key, nil // 返回key作为默认值
		}
		return "", fmt.Errorf("failed to get translation: %w", err)
	}
	return translation.Value, nil
}

// Set 设置翻译
func (r *TranslationRepository) Set(ctx context.Context, locale, key, value string) error {
	translation := &model.Translation{
		Locale:    locale,
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).
		Where("locale = ? AND key = ?", locale, key).
		Assign(map[string]interface{}{
			"value":      value,
			"updated_at": time.Now(),
		}).
		FirstOrCreate(translation).Error; err != nil {
		return fmt.Errorf("failed to set translation: %w", err)
	}
	return nil
}

// List 列出翻译
func (r *TranslationRepository) List(ctx context.Context, locale string) (map[string]string, error) {
	var translations []*model.Translation

	query := r.db.WithContext(ctx).Where("locale = ?", locale)
	if err := query.Limit(1000).Find(&translations).Error; err != nil {
		return nil, fmt.Errorf("failed to list translations: %w", err)
	}

	result := make(map[string]string)
	for _, t := range translations {
		result[t.Key] = t.Value
	}

	return result, nil
}

// AuditLogRepository 审计日志仓储
type AuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository 创建审计日志仓储
func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create 创建日志
func (r *AuditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

// List 列出日志
func (r *AuditLogRepository) List(ctx context.Context, userID int64, action string, page, pageSize int) ([]*model.AuditLog, int32, error) {
	var logs []*model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count logs: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list logs: %w", err)
	}

	return logs, int32(total), nil
}