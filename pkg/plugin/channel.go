package plugin

import (
	"context"
	"sync"
)

// ChannelAdapter 渠道适配器接口
type ChannelAdapter interface {
	// Name 返回适配器名称
	Name() string

	// ChannelType 返回渠道类型
	ChannelType() string

	// Description 返回适配器描述
	Description() string

	// Version 返回版本
	Version() string

	// Publish 发布内容
	Publish(ctx context.Context, req *PublishRequest) (*PublishResponse, error)

	// Preview 预览内容
	Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error)

	// GetStatus 获取发布状态
	GetStatus(ctx context.Context, externalID string) (*PublishStatus, error)

	// Validate 验证内容是否符合渠道要求
	Validate(ctx context.Context, content *Content) ([]ValidationIssue, error)
}

// PublishRequest 发布请求
type PublishRequest struct {
	TenantID  int64    `json:"tenant_id"`
	BrandID   int64    `json:"brand_id"`
	Content   *Content `json:"content"`
	Channel   *Channel `json:"channel"`
	Schedule  string   `json:"schedule,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Content 内容
type Content struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Summary  string `json:"summary"`
	Tags     []string `json:"tags"`
	CoverURL string `json:"cover_url"`
}

// Channel 渠道配置
type Channel struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Config   map[string]string `json:"config"`
}

// PublishResponse 发布响应
type PublishResponse struct {
	ExternalID  string `json:"external_id"`
	ExternalURL string `json:"external_url"`
	PublishedAt string `json:"published_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// PreviewRequest 预览请求
type PreviewRequest struct {
	Content *Content `json:"content"`
	Channel *Channel `json:"channel"`
}

// PreviewResponse 预览响应
type PreviewResponse struct {
	HTML     string `json:"html"`
	Preview  string `json:"preview"`
	Warnings []string `json:"warnings,omitempty"`
}

// PublishStatus 发布状态
type PublishStatus struct {
	ExternalID string            `json:"external_id"`
	Status     string            `json:"status"` // pending/published/failed/deleted
	Metrics    map[string]float64 `json:"metrics,omitempty"`
	UpdatedAt  string            `json:"updated_at"`
}

// ValidationIssue 校验问题
type ValidationIssue struct {
	Field       string `json:"field"`
	IssueType   string `json:"issue_type"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
	Severity    string `json:"severity"` // error/warning/info
}

// ChannelAdapterRegistry 渠道适配器注册表
type ChannelAdapterRegistry struct {
	mu       sync.RWMutex
	adapters map[string]ChannelAdapter
}

// NewChannelAdapterRegistry 创建渠道适配器注册表
func NewChannelAdapterRegistry() *ChannelAdapterRegistry {
	return &ChannelAdapterRegistry{
		adapters: make(map[string]ChannelAdapter),
	}
}

// Register 注册适配器
func (r *ChannelAdapterRegistry) Register(adapter ChannelAdapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.ChannelType()] = adapter
}

// Get 获取适配器
func (r *ChannelAdapterRegistry) Get(channelType string) (ChannelAdapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.adapters[channelType]
	return a, ok
}

// List 列出所有适配器
func (r *ChannelAdapterRegistry) List() []ChannelAdapter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adapters := make([]ChannelAdapter, 0, len(r.adapters))
	for _, a := range r.adapters {
		adapters = append(adapters, a)
	}
	return adapters
}

// 全局渠道适配器注册表
var globalChannelAdapterRegistry = NewChannelAdapterRegistry()

// RegisterChannelAdapter 注册渠道适配器到全局注册表
func RegisterChannelAdapter(adapter ChannelAdapter) {
	globalChannelAdapterRegistry.Register(adapter)
}

// GetChannelAdapter 从全局注册表获取渠道适配器
func GetChannelAdapter(channelType string) (ChannelAdapter, bool) {
	return globalChannelAdapterRegistry.Get(channelType)
}

// ListChannelAdapters 列出全局注册表中的所有适配器
func ListChannelAdapters() []ChannelAdapter {
	return globalChannelAdapterRegistry.List()
}
