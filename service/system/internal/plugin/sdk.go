package plugin

import "context"

// ==================== 插件 SDK 接口定义 ====================

// Plugin 插件基础接口
type Plugin interface {
	// Meta 返回插件元信息
	Meta() Meta
	// Init 初始化插件
	Init(ctx context.Context, config map[string]interface{}) error
	// Start 启动插件
	Start(ctx context.Context) error
	// Stop 停止插件
	Stop(ctx context.Context) error
	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error
}

// Meta 插件元信息
type Meta struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Type        string   `json:"type"` // channel, ai_model, analyzer
	Author      string   `json:"author"`
	Description string   `json:"description"`
	License     string   `json:"license"`
	Homepage    string   `json:"homepage"`
	Tags        []string `json:"tags"`
	MinVersion  string   `json:"min_version"` // 最低平台版本
}

// ==================== 渠道适配器插件 ====================

// ChannelAdapterPlugin 渠道适配器插件接口
type ChannelAdapterPlugin interface {
	Plugin

	// GetPlatform 获取平台标识
	GetPlatform() string
	// Publish 发布内容
	Publish(ctx context.Context, req *PublishRequest) (*PublishResponse, error)
	// Preview 预览内容
	Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error)
	// Validate 校验内容
	Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error)
	// GetMetrics 获取发布指标
	GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error)
}

// PublishRequest 发布请求
type PublishRequest struct {
	ContentID   int64             `json:"content_id"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	MediaURLs   []string          `json:"media_urls"`
	Tags        []string          `json:"tags"`
	Config      map[string]string `json:"config"`
}

// PublishResponse 发布响应
type PublishResponse struct {
	Success     bool   `json:"success"`
	ExternalID  string `json:"external_id"`
	ExternalURL string `json:"external_url"`
	ErrorMsg    string `json:"error_msg,omitempty"`
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
	Views    int64 `json:"views"`
	Likes    int64 `json:"likes"`
	Shares   int64 `json:"shares"`
	Comments int64 `json:"comments"`
}

// ==================== AI 模型连接器插件 ====================

// AIModelPlugin AI模型连接器插件接口
type AIModelPlugin interface {
	Plugin

	// GetModelName 获取模型名称
	GetModelName() string
	// Chat 对话补全
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	// Embedding 向量化
	Embedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
	// GetUsage 获取用量统计
	GetUsage(ctx context.Context) (*UsageStats, error)
}

// ChatRequest 对话请求
type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float32       `json:"temperature"`
	Stream      bool          `json:"stream"`
}

// ChatMessage 对话消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 对话响应
type ChatResponse struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
	PromptTokens int    `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

// EmbeddingRequest 向量化请求
type EmbeddingRequest struct {
	Texts []string `json:"texts"`
	Model string   `json:"model"`
}

// EmbeddingResponse 向量化响应
type EmbeddingResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
	Dimensions int         `json:"dimensions"`
}

// UsageStats 用量统计
type UsageStats struct {
	TotalRequests   int64   `json:"total_requests"`
	TotalTokens     int64   `json:"total_tokens"`
	TotalCost       float64 `json:"total_cost"`
	AvgLatency      float64 `json:"avg_latency"`
}

// ==================== 数据分析插件 ====================

// AnalyzerPlugin 数据分析插件接口
type AnalyzerPlugin interface {
	Plugin

	// GetAnalyzerType 获取分析器类型
	GetAnalyzerType() string
	// Analyze 执行分析
	Analyze(ctx context.Context, req *AnalyzeRequest) (*AnalyzeResponse, error)
	// GetCapabilities 获取分析能力
	GetCapabilities() []string
}

// AnalyzeRequest 分析请求
type AnalyzeRequest struct {
	ContentType string            `json:"content_type"`
	Data        interface{}       `json:"data"`
	Options     map[string]string `json:"options"`
}

// AnalyzeResponse 分析响应
type AnalyzeResponse struct {
	Success  bool                   `json:"success"`
	Results  map[string]interface{} `json:"results"`
	Score    float32                `json:"score"`
	Summary  string                 `json:"summary"`
	Details  []AnalyzeDetail        `json:"details"`
}

// AnalyzeDetail 分析详情
type AnalyzeDetail struct {
	Name    string      `json:"name"`
	Value   interface{} `json:"value"`
	Score   float32     `json:"score"`
	Suggest string      `json:"suggest"`
}

// ==================== 插件注册表 ====================

// Registry 插件注册表接口
type Registry interface {
	// Register 注册插件
	Register(plugin Plugin) error
	// Unregister 注销插件
	Unregister(name string) error
	// Get 获取插件
	Get(name string) (Plugin, error)
	// List 列出所有插件
	List() []Plugin
	// GetByType 按类型获取插件
	GetByType(pluginType string) []Plugin
	// Enable 启用插件
	Enable(name string) error
	// Disable 禁用插件
	Disable(name string) error
}

// PluginStatus 插件状态
type PluginStatus string

const (
	StatusRegistered PluginStatus = "registered"
	StatusEnabled    PluginStatus = "enabled"
	StatusDisabled   PluginStatus = "disabled"
	StatusError      PluginStatus = "error"
)

// PluginInfo 插件运行时信息
type PluginInfo struct {
	Meta      Meta         `json:"meta"`
	Status    PluginStatus `json:"status"`
	Error     string       `json:"error,omitempty"`
	LoadTime  int64        `json:"load_time"`
	CallCount int64        `json:"call_count"`
	AvgLatency float64     `json:"avg_latency"`
}
