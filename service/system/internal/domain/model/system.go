package model

import "time"

// SystemConfig 系统配置
type SystemConfig struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ConfigKey   string    `json:"config_key" gorm:"uniqueIndex;size:128;not null"`
	ConfigValue string    `json:"config_value" gorm:"type:text"`
	ConfigType  string    `json:"config_type" gorm:"size:32"` // string, number, json, boolean
	Description string    `json:"description" gorm:"size:256"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Plugin 插件
type Plugin struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	PluginName    string    `json:"plugin_name" gorm:"uniqueIndex;size:64;not null"`
	PluginType    string    `json:"plugin_type" gorm:"size:32"` // channel, ai, analyzer
	Description   string    `json:"description" gorm:"size:256"`
	Version       string    `json:"version" gorm:"size:32"`
	Author        string    `json:"author" gorm:"size:64"`
	ConfigSchema  string    `json:"config_schema" gorm:"type:text"` // JSON格式的配置模式
	IsEnabled     bool      `json:"is_enabled" gorm:"default:true"`
	InstalledAt   time.Time `json:"installed_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Webhook Webhook配置
type Webhook struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int64     `json:"user_id" gorm:"index;not null"`
	WebhookName string    `json:"webhook_name" gorm:"size:128;not null"`
	URL         string    `json:"url" gorm:"size:512;not null"`
	Secret      string    `json:"secret" gorm:"size:128"`
	Events      string    `json:"events" gorm:"type:text"` // JSON格式的事件列表
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	LastTrigger *time.Time `json:"last_trigger"`
	FailCount   int32     `json:"fail_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WebhookEvent Webhook事件
type WebhookEvent struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	WebhookID   int64     `json:"webhook_id" gorm:"index;not null"`
	EventType   string    `json:"event_type" gorm:"size:64;not null"`
	Payload     string    `json:"payload" gorm:"type:text"`
	StatusCode  int32     `json:"status_code"`
	ResponseBody string   `json:"response_body" gorm:"type:text"`
	Success     bool      `json:"success"`
	TriggeredAt time.Time `json:"triggered_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Tenant 租户
type Tenant struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantName   string    `json:"tenant_name" gorm:"uniqueIndex;size:128;not null"`
	TenantCode   string    `json:"tenant_code" gorm:"uniqueIndex;size:64;not null"`
	ContactName  string    `json:"contact_name" gorm:"size:64"`
	ContactEmail string    `json:"contact_email" gorm:"size:128"`
	Status       int32     `json:"status" gorm:"default:1"` // 1:正常 0:禁用
	QuotaConfig  string    `json:"quota_config" gorm:"type:text"` // JSON格式的配额配置
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Translation 国际化翻译
type Translation struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Locale    string    `json:"locale" gorm:"size:32;not null;index:idx_translation_lookup"`
	Key       string    `json:"key" gorm:"size:256;not null;index:idx_translation_lookup"`
	Value     string    `json:"value" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuditLog 审计日志
type AuditLog struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       int64     `json:"user_id" gorm:"index"`
	Username     string    `json:"username" gorm:"size:64"`
	Action       string    `json:"action" gorm:"size:64;not null;index"`
	ResourceType string    `json:"resource_type" gorm:"size:64"`
	ResourceID   int64     `json:"resource_id"`
	Details      string    `json:"details" gorm:"type:text"`
	IPAddress    string    `json:"ip_address" gorm:"size:64"`
	UserAgent    string    `json:"user_agent" gorm:"size:256"`
	CreatedAt    time.Time `json:"created_at"`
}