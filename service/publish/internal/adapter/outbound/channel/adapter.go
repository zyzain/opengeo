package channel

import (
	"context"
	"fmt"
	"time"
)

// ChannelAdapter 渠道适配器接口
type ChannelAdapter interface {
	// GetName 获取适配器名称
	GetName() string
	// GetPlatform 获取平台类型
	GetPlatform() string
	// Publish 发布内容
	Publish(ctx context.Context, req *PublishRequest, publishCtx *AdapterPublishContext) (*PublishResponse, error)
	// Preview 预览内容
	Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error)
	// Validate 校验内容
	Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error)
	// GetMetrics 获取发布指标
	GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error)
}

// AdapterPublishContext 发布上下文（来自防封引擎）
type AdapterPublishContext struct {
	ProxyURL      string            `json:"proxy_url"`
	Headers       map[string]string `json:"headers"`
	UserAgent     string            `json:"user_agent"`
	FingerprintID int64             `json:"fingerprint_id"`
}

// PublishRequest 发布请求
type PublishRequest struct {
	ContentID   int64             `json:"content_id"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	MediaURLs   []string          `json:"media_urls"`
	Tags        []string          `json:"tags"`
	Config      map[string]string `json:"config"`
	ScheduledAt *time.Time        `json:"scheduled_at"`
}

// PublishResponse 发布响应
type PublishResponse struct {
	Success     bool      `json:"success"`
	ExternalID  string    `json:"external_id"`
	ExternalURL string    `json:"external_url"`
	PublishedAt time.Time `json:"published_at"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
}

// PreviewRequest 预览请求
type PreviewRequest struct {
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	MediaURLs []string          `json:"media_urls"`
	Config    map[string]string `json:"config"`
}

// PreviewResponse 预览响应
type PreviewResponse struct {
	Success  bool   `json:"success"`
	HTML     string `json:"html"`
	Markdown string `json:"markdown"`
	Warnings []string `json:"warnings,omitempty"`
}

// ValidateRequest 校验请求
type ValidateRequest struct {
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	MediaURLs []string          `json:"media_urls"`
	Config    map[string]string `json:"config"`
}

// ValidateResponse 校验响应
type ValidateResponse struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// MetricsRequest 指标请求
type MetricsRequest struct {
	ExternalID string `json:"external_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

// MetricsResponse 指标响应
type MetricsResponse struct {
	Views    int64   `json:"views"`
	Likes    int64   `json:"likes"`
	Shares   int64   `json:"shares"`
	Comments int64   `json:"comments"`
	Reach    int64   `json:"reach"`
}

// ==================== 微信公众号适配器 ====================

// WechatAdapter 微信公众号适配器
type WechatAdapter struct {
	appID     string
	appSecret string
}

// NewWechatAdapter 创建微信公众号适配器
func NewWechatAdapter(appID, appSecret string) *WechatAdapter {
	return &WechatAdapter{
		appID:     appID,
		appSecret: appSecret,
	}
}

func (a *WechatAdapter) GetName() string {
	return "wechat"
}

func (a *WechatAdapter) GetPlatform() string {
	return "wechat"
}

func (a *WechatAdapter) Publish(ctx context.Context, req *PublishRequest, publishCtx *AdapterPublishContext) (*PublishResponse, error) {
	// TODO: 调用微信公众号API
	return &PublishResponse{
		Success:     true,
		ExternalID:  fmt.Sprintf("wechat_%d", time.Now().Unix()),
		ExternalURL: "https://mp.weixin.qq.com/...",
		PublishedAt: time.Now(),
	}, nil
}

func (a *WechatAdapter) Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error) {
	// 生成微信公众号预览HTML
	html := fmt.Sprintf(`
		<div class="wechat-article">
			<h1>%s</h1>
			<div class="content">%s</div>
		</div>
	`, req.Title, req.Body)

	return &PreviewResponse{
		Success: true,
		HTML:    html,
	}, nil
}

func (a *WechatAdapter) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	errors := make([]string, 0)
	warnings := make([]string, 0)

	if req.Title == "" {
		errors = append(errors, "标题不能为空")
	}
	if len(req.Title) > 64 {
		warnings = append(warnings, "标题超过64字可能被截断")
	}
	if req.Body == "" {
		errors = append(errors, "正文不能为空")
	}

	return &ValidateResponse{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

func (a *WechatAdapter) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	// TODO: 调用微信公众号API获取指标
	return &MetricsResponse{}, nil
}

// ==================== 微博适配器 ====================

// WeiboAdapter 微博适配器
type WeiboAdapter struct {
	accessToken string
}

// NewWeiboAdapter 创建微博适配器
func NewWeiboAdapter(accessToken string) *WeiboAdapter {
	return &WeiboAdapter{accessToken: accessToken}
}

func (a *WeiboAdapter) GetName() string {
	return "weibo"
}

func (a *WeiboAdapter) GetPlatform() string {
	return "weibo"
}

func (a *WeiboAdapter) Publish(ctx context.Context, req *PublishRequest, publishCtx *AdapterPublishContext) (*PublishResponse, error) {
	// TODO: 调用微博API
	return &PublishResponse{
		Success:     true,
		ExternalID:  fmt.Sprintf("weibo_%d", time.Now().Unix()),
		PublishedAt: time.Now(),
	}, nil
}

func (a *WeiboAdapter) Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error) {
	// 微博限制140字
	body := req.Body
	if len(body) > 140 {
		body = body[:140] + "..."
	}

	return &PreviewResponse{
		Success:  true,
		Markdown: fmt.Sprintf("%s\n\n%s", req.Title, body),
		Warnings: []string{"微博内容限制140字"},
	}, nil
}

func (a *WeiboAdapter) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	errors := make([]string, 0)
	warnings := make([]string, 0)

	if req.Title == "" && req.Body == "" {
		errors = append(errors, "内容不能为空")
	}
	if len(req.Body) > 140 {
		warnings = append(warnings, "内容超过140字将被截断")
	}

	return &ValidateResponse{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

func (a *WeiboAdapter) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	return &MetricsResponse{}, nil
}

// ==================== 抖音适配器 ====================

// DouyinAdapter 抖音适配器
type DouyinAdapter struct {
	accessToken string
}

// NewDouyinAdapter 创建抖音适配器
func NewDouyinAdapter(accessToken string) *DouyinAdapter {
	return &DouyinAdapter{accessToken: accessToken}
}

func (a *DouyinAdapter) GetName() string {
	return "douyin"
}

func (a *DouyinAdapter) GetPlatform() string {
	return "douyin"
}

func (a *DouyinAdapter) Publish(ctx context.Context, req *PublishRequest, publishCtx *AdapterPublishContext) (*PublishResponse, error) {
	// TODO: 调用抖音API
	return &PublishResponse{
		Success:     true,
		ExternalID:  fmt.Sprintf("douyin_%d", time.Now().Unix()),
		PublishedAt: time.Now(),
	}, nil
}

func (a *DouyinAdapter) Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error) {
	return &PreviewResponse{
		Success:  true,
		Markdown: fmt.Sprintf("%s\n\n%s", req.Title, req.Body),
	}, nil
}

func (a *DouyinAdapter) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	errors := make([]string, 0)

	if req.Title == "" {
		errors = append(errors, "标题不能为空")
	}
	if len(req.MediaURLs) == 0 {
		errors = append(errors, "抖音需要至少一个视频或图片")
	}

	return &ValidateResponse{
		Valid:  len(errors) == 0,
		Errors: errors,
	}, nil
}

func (a *DouyinAdapter) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	return &MetricsResponse{}, nil
}

// ==================== 小红书适配器 ====================

// XiaohongshuAdapter 小红书适配器
type XiaohongshuAdapter struct {
	accessToken string
}

// NewXiaohongshuAdapter 创建小红书适配器
func NewXiaohongshuAdapter(accessToken string) *XiaohongshuAdapter {
	return &XiaohongshuAdapter{accessToken: accessToken}
}

func (a *XiaohongshuAdapter) GetName() string {
	return "xiaohongshu"
}

func (a *XiaohongshuAdapter) GetPlatform() string {
	return "xiaohongshu"
}

func (a *XiaohongshuAdapter) Publish(ctx context.Context, req *PublishRequest, publishCtx *AdapterPublishContext) (*PublishResponse, error) {
	// TODO: 调用小红书API
	return &PublishResponse{
		Success:     true,
		ExternalID:  fmt.Sprintf("xhs_%d", time.Now().Unix()),
		PublishedAt: time.Now(),
	}, nil
}

func (a *XiaohongshuAdapter) Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error) {
	return &PreviewResponse{
		Success:  true,
		Markdown: fmt.Sprintf("%s\n\n%s", req.Title, req.Body),
	}, nil
}

func (a *XiaohongshuAdapter) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	errors := make([]string, 0)

	if req.Title == "" {
		errors = append(errors, "标题不能为空")
	}
	if len(req.Title) > 20 {
		errors = append(errors, "小红书标题限制20字")
	}

	return &ValidateResponse{
		Valid:  len(errors) == 0,
		Errors: errors,
	}, nil
}

func (a *XiaohongshuAdapter) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	return &MetricsResponse{}, nil
}

// ==================== 适配器工厂 ====================

// ChannelAdapterFactory 渠道适配器工厂
type ChannelAdapterFactory struct {
	adapters map[string]ChannelAdapter
}

// NewChannelAdapterFactory 创建渠道适配器工厂
func NewChannelAdapterFactory() *ChannelAdapterFactory {
	return &ChannelAdapterFactory{
		adapters: make(map[string]ChannelAdapter),
	}
}

// RegisterAdapter 注册适配器
func (f *ChannelAdapterFactory) RegisterAdapter(adapter ChannelAdapter) {
	f.adapters[adapter.GetPlatform()] = adapter
}

// GetAdapter 获取适配器
func (f *ChannelAdapterFactory) GetAdapter(platform string) (ChannelAdapter, error) {
	adapter, ok := f.adapters[platform]
	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
	return adapter, nil
}

// GetAllAdapters 获取所有适配器
func (f *ChannelAdapterFactory) GetAllAdapters() map[string]ChannelAdapter {
	return f.adapters
}

// GetSupportedPlatforms 获取支持的平台
func (f *ChannelAdapterFactory) GetSupportedPlatforms() []string {
	platforms := make([]string, 0, len(f.adapters))
	for platform := range f.adapters {
		platforms = append(platforms, platform)
	}
	return platforms
}