package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"opengeo/service/system/internal/dal"
	"opengeo/service/system/internal/domain/model"
)

// SystemService 系统服务
type SystemService struct {
	configRepo     *dal.SystemConfigRepository
	pluginRepo     *dal.PluginRepository
	webhookRepo    *dal.WebhookRepository
	tenantRepo     *dal.TenantRepository
	translationRepo *dal.TranslationRepository
	auditLogRepo   *dal.AuditLogRepository
}

// NewSystemService 创建系统服务
func NewSystemService(
	configRepo *dal.SystemConfigRepository,
	pluginRepo *dal.PluginRepository,
	webhookRepo *dal.WebhookRepository,
	tenantRepo *dal.TenantRepository,
	translationRepo *dal.TranslationRepository,
	auditLogRepo *dal.AuditLogRepository,
) *SystemService {
	return &SystemService{
		configRepo:      configRepo,
		pluginRepo:      pluginRepo,
		webhookRepo:     webhookRepo,
		tenantRepo:      tenantRepo,
		translationRepo: translationRepo,
		auditLogRepo:    auditLogRepo,
	}
}

// ==================== 系统配置 ====================

// GetConfig 获取配置
func (s *SystemService) GetConfig(ctx context.Context, key string) (*model.SystemConfig, error) {
	config, err := s.configRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return config, nil
}

// SetConfig 设置配置
func (s *SystemService) SetConfig(ctx context.Context, key, value, configType, description string, isPublic bool) error {
	if err := s.configRepo.Set(ctx, key, value, configType, description, isPublic); err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}
	return nil
}

// ListConfigs 列出配置
func (s *SystemService) ListConfigs(ctx context.Context, isPublic *bool, page, pageSize int) ([]*model.SystemConfig, int32, error) {
	configs, total, err := s.configRepo.List(ctx, isPublic, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list configs: %w", err)
	}
	return configs, total, nil
}

// DeleteConfig 删除配置
func (s *SystemService) DeleteConfig(ctx context.Context, key string) error {
	if err := s.configRepo.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}
	return nil
}

// ==================== 插件管理 ====================

// CreatePlugin 创建插件
func (s *SystemService) CreatePlugin(ctx context.Context, name, pluginType, description, version, author, configSchema string) (*model.Plugin, error) {
	plugin := &model.Plugin{
		PluginName:   name,
		PluginType:   pluginType,
		Description:  description,
		Version:      version,
		Author:       author,
		ConfigSchema: configSchema,
		IsEnabled:    true,
		InstalledAt:  time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.pluginRepo.Create(ctx, plugin); err != nil {
		return nil, fmt.Errorf("failed to create plugin: %w", err)
	}

	return plugin, nil
}

// GetPlugin 获取插件
func (s *SystemService) GetPlugin(ctx context.Context, pluginID int64) (*model.Plugin, error) {
	plugin, err := s.pluginRepo.GetByID(ctx, pluginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin: %w", err)
	}
	return plugin, nil
}

// ListPlugins 列出插件
func (s *SystemService) ListPlugins(ctx context.Context, pluginType string, isEnabled *bool, page, pageSize int) ([]*model.Plugin, int32, error) {
	plugins, total, err := s.pluginRepo.List(ctx, pluginType, isEnabled, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list plugins: %w", err)
	}
	return plugins, total, nil
}

// EnablePlugin 启用插件
func (s *SystemService) EnablePlugin(ctx context.Context, pluginID int64) error {
	plugin, err := s.pluginRepo.GetByID(ctx, pluginID)
	if err != nil {
		return fmt.Errorf("failed to get plugin: %w", err)
	}

	plugin.IsEnabled = true
	plugin.UpdatedAt = time.Now()

	if err := s.pluginRepo.Update(ctx, plugin); err != nil {
		return fmt.Errorf("failed to enable plugin: %w", err)
	}

	return nil
}

// DisablePlugin 禁用插件
func (s *SystemService) DisablePlugin(ctx context.Context, pluginID int64) error {
	plugin, err := s.pluginRepo.GetByID(ctx, pluginID)
	if err != nil {
		return fmt.Errorf("failed to get plugin: %w", err)
	}

	plugin.IsEnabled = false
	plugin.UpdatedAt = time.Now()

	if err := s.pluginRepo.Update(ctx, plugin); err != nil {
		return fmt.Errorf("failed to disable plugin: %w", err)
	}

	return nil
}

// DeletePlugin 删除插件
func (s *SystemService) DeletePlugin(ctx context.Context, pluginID int64) error {
	if err := s.pluginRepo.Delete(ctx, pluginID); err != nil {
		return fmt.Errorf("failed to delete plugin: %w", err)
	}
	return nil
}

// ==================== Webhook管理 ====================

// CreateWebhook 创建Webhook
func (s *SystemService) CreateWebhook(ctx context.Context, userID int64, name, url, secret string, events []string) (*model.Webhook, error) {
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal events: %w", err)
	}

	webhook := &model.Webhook{
		UserID:      userID,
		WebhookName: name,
		URL:         url,
		Secret:      secret,
		Events:      string(eventsJSON),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.webhookRepo.Create(ctx, webhook); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return webhook, nil
}

// GetWebhook 获取Webhook
func (s *SystemService) GetWebhook(ctx context.Context, webhookID int64) (*model.Webhook, error) {
	webhook, err := s.webhookRepo.GetByID(ctx, webhookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}
	return webhook, nil
}

// ListWebhooks 列出Webhook
func (s *SystemService) ListWebhooks(ctx context.Context, userID int64, page, pageSize int) ([]*model.Webhook, int32, error) {
	webhooks, total, err := s.webhookRepo.List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list webhooks: %w", err)
	}
	return webhooks, total, nil
}

// UpdateWebhook 更新Webhook
func (s *SystemService) UpdateWebhook(ctx context.Context, webhookID int64, name, url, secret string, events []string) (*model.Webhook, error) {
	webhook, err := s.webhookRepo.GetByID(ctx, webhookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	if name != "" {
		webhook.WebhookName = name
	}
	if url != "" {
		webhook.URL = url
	}
	if secret != "" {
		webhook.Secret = secret
	}
	if events != nil {
		eventsJSON, err := json.Marshal(events)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal events: %w", err)
		}
		webhook.Events = string(eventsJSON)
	}
	webhook.UpdatedAt = time.Now()

	if err := s.webhookRepo.Update(ctx, webhook); err != nil {
		return nil, fmt.Errorf("failed to update webhook: %w", err)
	}

	return webhook, nil
}

// DeleteWebhook 删除Webhook
func (s *SystemService) DeleteWebhook(ctx context.Context, webhookID int64) error {
	if err := s.webhookRepo.Delete(ctx, webhookID); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	return nil
}

// TriggerWebhook 触发Webhook
func (s *SystemService) TriggerWebhook(ctx context.Context, webhookID int64, eventType string, payload interface{}) error {
	webhook, err := s.webhookRepo.GetByID(ctx, webhookID)
	if err != nil {
		return fmt.Errorf("failed to get webhook: %w", err)
	}

	if !webhook.IsActive {
		return fmt.Errorf("webhook is disabled")
	}

	// 检查事件类型是否匹配
	var events []string
	if err := json.Unmarshal([]byte(webhook.Events), &events); err != nil {
		return fmt.Errorf("failed to unmarshal events: %w", err)
	}

	eventMatched := false
	for _, event := range events {
		if event == eventType || event == "*" {
			eventMatched = true
			break
		}
	}

	if !eventMatched {
		return fmt.Errorf("event type not matched: %s", eventType)
	}

	// 序列化payload
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 生成签名
	_ = s.generateSignature(webhook.Secret, string(payloadJSON))

	// TODO: 发送HTTP请求到webhook.URL
	// 这里简化处理，直接记录事件
	event := &model.WebhookEvent{
		WebhookID:   webhookID,
		EventType:   eventType,
		Payload:     string(payloadJSON),
		StatusCode:  200,
		ResponseBody: "OK",
		Success:     true,
		TriggeredAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := s.webhookRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to create webhook event: %w", err)
	}

	// 更新最后触发时间
	webhook.LastTrigger = &event.TriggeredAt
	webhook.FailCount = 0
	if err := s.webhookRepo.Update(ctx, webhook); err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	return nil
}

// TestWebhook 测试Webhook
func (s *SystemService) TestWebhook(ctx context.Context, webhookID int64) (*model.WebhookEvent, error) {
	testPayload := map[string]interface{}{
		"event": "test",
		"data": map[string]interface{}{
			"message": "This is a test webhook",
		},
		"timestamp": time.Now().Unix(),
	}

	if err := s.TriggerWebhook(ctx, webhookID, "test", testPayload); err != nil {
		return nil, fmt.Errorf("failed to test webhook: %w", err)
	}

	// 获取最新的事件
	events, _, err := s.webhookRepo.GetEvents(ctx, webhookID, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no events found")
	}

	return events[0], nil
}

// GetWebhookEvents 获取Webhook事件
func (s *SystemService) GetWebhookEvents(ctx context.Context, webhookID int64, page, pageSize int) ([]*model.WebhookEvent, int32, error) {
	events, total, err := s.webhookRepo.GetEvents(ctx, webhookID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get webhook events: %w", err)
	}
	return events, total, nil
}

// generateSignature 生成签名
func (s *SystemService) generateSignature(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// ==================== 租户管理 ====================

// CreateTenant 创建租户
func (s *SystemService) CreateTenant(ctx context.Context, name, code, contactName, contactEmail, quotaConfig string) (*model.Tenant, error) {
	tenant := &model.Tenant{
		TenantName:   name,
		TenantCode:   code,
		ContactName:  contactName,
		ContactEmail: contactEmail,
		Status:       1,
		QuotaConfig:  quotaConfig,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

// GetTenant 获取租户
func (s *SystemService) GetTenant(ctx context.Context, tenantID int64) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	return tenant, nil
}

// ListTenants 列出租户
func (s *SystemService) ListTenants(ctx context.Context, page, pageSize int) ([]*model.Tenant, int32, error) {
	tenants, total, err := s.tenantRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tenants: %w", err)
	}
	return tenants, total, nil
}

// UpdateTenant 更新租户
func (s *SystemService) UpdateTenant(ctx context.Context, tenantID int64, contactName, contactEmail, quotaConfig string) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	if contactName != "" {
		tenant.ContactName = contactName
	}
	if contactEmail != "" {
		tenant.ContactEmail = contactEmail
	}
	if quotaConfig != "" {
		tenant.QuotaConfig = quotaConfig
	}
	tenant.UpdatedAt = time.Now()

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return tenant, nil
}

// DeleteTenant 删除租户
func (s *SystemService) DeleteTenant(ctx context.Context, tenantID int64) error {
	if err := s.tenantRepo.Delete(ctx, tenantID); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}

// ==================== 国际化 ====================

// GetTranslation 获取翻译
func (s *SystemService) GetTranslation(ctx context.Context, locale, key string) (string, error) {
	value, err := s.translationRepo.Get(ctx, locale, key)
	if err != nil {
		return "", fmt.Errorf("failed to get translation: %w", err)
	}
	return value, nil
}

// SetTranslation 设置翻译
func (s *SystemService) SetTranslation(ctx context.Context, locale, key, value string) error {
	if err := s.translationRepo.Set(ctx, locale, key, value); err != nil {
		return fmt.Errorf("failed to set translation: %w", err)
	}
	return nil
}

// ListTranslations 列出翻译
func (s *SystemService) ListTranslations(ctx context.Context, locale string) (map[string]string, error) {
	translations, err := s.translationRepo.List(ctx, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to list translations: %w", err)
	}
	return translations, nil
}

// ==================== 审计日志 ====================

// CreateAuditLog 创建审计日志
func (s *SystemService) CreateAuditLog(ctx context.Context, userID int64, username, action, resourceType string, resourceID int64, details, ipAddress, userAgent string) error {
	log := &model.AuditLog{
		UserID:       userID,
		Username:     username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}

	if err := s.auditLogRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// ListAuditLogs 列出审计日志
func (s *SystemService) ListAuditLogs(ctx context.Context, userID int64, action string, page, pageSize int) ([]*model.AuditLog, int32, error) {
	logs, total, err := s.auditLogRepo.List(ctx, userID, action, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}
	return logs, total, nil
}