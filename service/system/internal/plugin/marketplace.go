package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Marketplace 插件市场服务
type Marketplace struct {
	registry Registry
}

// NewMarketplace 创建插件市场服务
func NewMarketplace(registry Registry) *Marketplace {
	return &Marketplace{registry: registry}
}

// ==================== 插件查询 ====================

// SearchRequest 搜索插件请求
type SearchRequest struct {
	Keyword    string `json:"keyword"`
	Type       string `json:"type"`
	Tags       []string `json:"tags"`
	SortBy     string `json:"sort_by"` // downloads, rating, updated
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// PluginListing 插件列表项
type PluginListing struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Type        string   `json:"type"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Downloads   int64    `json:"downloads"`
	Rating      float32  `json:"rating"`
	RatingCount int32    `json:"rating_count"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsOfficial  bool     `json:"is_official"`
	IsVerified  bool     `json:"is_verified"`
}

// PluginDetail 插件详情
type PluginDetail struct {
	PluginListing
	License      string            `json:"license"`
	Homepage     string            `json:"homepage"`
	SourceURL    string            `json:"source_url"`
	BinaryURL    string            `json:"binary_url"`
	ConfigSchema string            `json:"config_schema"`
	MinVersion   string            `json:"min_version"`
	Changelog    string            `json:"changelog"`
	Examples     []PluginExample   `json:"examples"`
	Reviews      []PluginReview    `json:"reviews"`
}

// PluginExample 插件示例
type PluginExample struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Config      string `json:"config"`
}

// PluginReview 插件评价
type PluginReview struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Rating    int32     `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== 插件安装 ====================

// InstallRequest 安装插件请求
type InstallRequest struct {
	PluginID int64  `json:"plugin_id"`
	Version  string `json:"version"`
	Config   string `json:"config"`
}

// InstallResponse 安装插件响应
type InstallResponse struct {
	Success     bool   `json:"success"`
	PluginName  string `json:"plugin_name"`
	Version     string `json:"version"`
	Message     string `json:"message"`
}

// ==================== SDK 生成 ====================

// SDKTemplate SDK模板
type SDKTemplate struct {
	PluginType  string `json:"plugin_type"`
	Language    string `json:"language"`
	Template    string `json:"template"`
	Example     string `json:"example"`
	Docs        string `json:"docs"`
}

// GetSDKTemplate 获取SDK模板
func (m *Marketplace) GetSDKTemplate(ctx context.Context, pluginType, language string) (*SDKTemplate, error) {
	templates := map[string]map[string]*SDKTemplate{
		"channel": {
			"go": {
				PluginType: "channel",
				Language:   "go",
				Template: `package main

import (
	"context"
	"opengeo/plugin"
)

type MyChannelAdapter struct {
	platform string
	config   map[string]interface{}
}

func (a *MyChannelAdapter) Meta() plugin.Meta {
	return plugin.Meta{
		Name:        "my-channel",
		Version:     "1.0.0",
		Type:        "channel",
		Author:      "Your Name",
		Description: "My custom channel adapter",
	}
}

func (a *MyChannelAdapter) Init(ctx context.Context, config map[string]interface{}) error {
	a.config = config
	return nil
}

func (a *MyChannelAdapter) Start(ctx context.Context) error {
	return nil
}

func (a *MyChannelAdapter) Stop(ctx context.Context) error {
	return nil
}

func (a *MyChannelAdapter) HealthCheck(ctx context.Context) error {
	return nil
}

func (a *MyChannelAdapter) GetPlatform() string {
	return a.platform
}

func (a *MyChannelAdapter) Publish(ctx context.Context, req *plugin.PublishRequest) (*plugin.PublishResponse, error) {
	// Implement your publish logic here
	return &plugin.PublishResponse{
		Success:    true,
		ExternalID: "your-external-id",
	}, nil
}

func (a *MyChannelAdapter) Preview(ctx context.Context, req *plugin.PreviewRequest) (*plugin.PreviewResponse, error) {
	return &plugin.PreviewResponse{
		Success: true,
		HTML:    "<h1>" + req.Title + "</h1><p>" + req.Body + "</p>",
	}, nil
}

func (a *MyChannelAdapter) Validate(ctx context.Context, req *plugin.ValidateRequest) (*plugin.ValidateResponse, error) {
	return &plugin.ValidateResponse{Valid: true}, nil
}

func (a *MyChannelAdapter) GetMetrics(ctx context.Context, req *plugin.MetricsRequest) (*plugin.MetricsResponse, error) {
	return &plugin.MetricsResponse{}, nil
}

func main() {
	adapter := &MyChannelAdapter{platform: "my-platform"}
	plugin.Register(adapter)
}`,
				Example: `// Usage example
adapter := &MyChannelAdapter{platform: "my-platform"}
registry.Register(adapter)

// Get and use the plugin
p, _ := registry.Get("my-channel")
channelPlugin := p.(plugin.ChannelAdapterPlugin)
result, _ := channelPlugin.Publish(ctx, &plugin.PublishRequest{
    Title: "Hello",
    Body:  "World",
})`,
				Docs: `## Channel Adapter Plugin SDK

### Interface Methods

- Meta() - Return plugin metadata
- Init() - Initialize with config
- Start()/Stop() - Lifecycle management
- GetPlatform() - Return platform identifier
- Publish() - Publish content to platform
- Preview() - Preview content rendering
- Validate() - Validate content format
- GetMetrics() - Get publishing metrics

### Best Practices

1. Handle errors gracefully
2. Implement retry logic for API calls
3. Cache authentication tokens
4. Validate content before publishing`,
			},
		},
		"ai_model": {
			"go": {
				PluginType: "ai_model",
				Language:   "go",
				Template: `package main

import (
	"context"
	"opengeo/plugin"
)

type MyAIModel struct {
	modelName string
	apiKey    string
	config    map[string]interface{}
}

func (m *MyAIModel) Meta() plugin.Meta {
	return plugin.Meta{
		Name:        "my-ai-model",
		Version:     "1.0.0",
		Type:        "ai_model",
		Author:      "Your Name",
		Description: "My custom AI model connector",
	}
}

func (m *MyAIModel) Init(ctx context.Context, config map[string]interface{}) error {
	m.config = config
	if key, ok := config["api_key"].(string); ok {
		m.apiKey = key
	}
	return nil
}

func (m *MyAIModel) Start(ctx context.Context) error {
	return nil
}

func (m *MyAIModel) Stop(ctx context.Context) error {
	return nil
}

func (m *MyAIModel) HealthCheck(ctx context.Context) error {
	return nil
}

func (m *MyAIModel) GetModelName() string {
	return m.modelName
}

func (m *MyAIModel) Chat(ctx context.Context, req *plugin.ChatRequest) (*plugin.ChatResponse, error) {
	// Implement your chat logic here
	return &plugin.ChatResponse{
		Content:      "Response from my AI model",
		FinishReason: "stop",
	}, nil
}

func (m *MyAIModel) Embedding(ctx context.Context, req *plugin.EmbeddingRequest) (*plugin.EmbeddingResponse, error) {
	return &plugin.EmbeddingResponse{
		Embeddings: make([][]float32, len(req.Texts)),
		Dimensions: 1536,
	}, nil
}

func (m *MyAIModel) GetUsage(ctx context.Context) (*plugin.UsageStats, error) {
	return &plugin.UsageStats{}, nil
}`,
			},
		},
		"analyzer": {
			"go": {
				PluginType: "analyzer",
				Language:   "go",
				Template: `package main

import (
	"context"
	"opengeo/plugin"
)

type MyAnalyzer struct {
	analyzerType string
	config       map[string]interface{}
}

func (a *MyAnalyzer) Meta() plugin.Meta {
	return plugin.Meta{
		Name:        "my-analyzer",
		Version:     "1.0.0",
		Type:        "analyzer",
		Author:      "Your Name",
		Description: "My custom data analyzer",
	}
}

func (a *MyAnalyzer) Init(ctx context.Context, config map[string]interface{}) error {
	a.config = config
	return nil
}

func (a *MyAnalyzer) Start(ctx context.Context) error {
	return nil
}

func (a *MyAnalyzer) Stop(ctx context.Context) error {
	return nil
}

func (a *MyAnalyzer) HealthCheck(ctx context.Context) error {
	return nil
}

func (a *MyAnalyzer) GetAnalyzerType() string {
	return a.analyzerType
}

func (a *MyAnalyzer) Analyze(ctx context.Context, req *plugin.AnalyzeRequest) (*plugin.AnalyzeResponse, error) {
	return &plugin.AnalyzeResponse{
		Success: true,
		Results: map[string]interface{}{},
		Score:   85.0,
		Summary: "Analysis complete",
	}, nil
}

func (a *MyAnalyzer) GetCapabilities() []string {
	return []string{"seo_analysis", "content_scoring", "keyword_extraction"}
}`,
			},
		},
	}

	if typeTemplates, ok := templates[pluginType]; ok {
		if tmpl, ok := typeTemplates[language]; ok {
			return tmpl, nil
		}
	}

	return nil, fmt.Errorf("unsupported plugin type %s or language %s", pluginType, language)
}

// ==================== 内置插件注册 ====================

// RegisterBuiltinPlugins 注册内置插件
func RegisterBuiltinPlugins(registry Registry) error {
	builtinPlugins := []struct {
		name        string
		pluginType  string
		description string
		author      string
		version     string
	}{
		{"wechat", "channel", "微信公众号适配器", "OpenGEO", "1.0.0"},
		{"weibo", "channel", "微博适配器", "OpenGEO", "1.0.0"},
		{"douyin", "channel", "抖音适配器", "OpenGEO", "1.0.0"},
		{"xiaohongshu", "channel", "小红书适配器", "OpenGEO", "1.0.0"},
		{"zhihu", "channel", "知乎适配器", "OpenGEO", "1.0.0"},
		{"deepseek", "ai_model", "DeepSeek AI模型连接器", "OpenGEO", "1.0.0"},
		{"kimi", "ai_model", "Kimi AI模型连接器", "OpenGEO", "1.0.0"},
		{"doubao", "ai_model", "豆包AI模型连接器", "OpenGEO", "1.0.0"},
		{"chatgpt", "ai_model", "ChatGPT模型连接器", "OpenGEO", "1.0.0"},
		{"seo_analyzer", "analyzer", "SEO分析器", "OpenGEO", "1.0.0"},
		{"content_scorer", "analyzer", "内容评分器", "OpenGEO", "1.0.0"},
		{"geo_optimizer", "analyzer", "GEO优化分析器", "OpenGEO", "1.0.0"},
	}

	for _, bp := range builtinPlugins {
		p := &BuiltinPlugin{
			meta: Meta{
				Name:        bp.name,
				Version:     bp.version,
				Type:        bp.pluginType,
				Author:      bp.author,
				Description: bp.description,
				License:     "MIT",
			},
		}
		if err := registry.Register(p); err != nil {
			// 忽略重复注册
			continue
		}
		// 内置插件默认启用
		registry.Enable(bp.name)
	}

	return nil
}

// BuiltinPlugin 内置插件基类
type BuiltinPlugin struct {
	meta   Meta
	config map[string]interface{}
}

func (p *BuiltinPlugin) Meta() Meta { return p.meta }

func (p *BuiltinPlugin) Init(ctx context.Context, config map[string]interface{}) error {
	p.config = config
	return nil
}

func (p *BuiltinPlugin) Start(ctx context.Context) error  { return nil }
func (p *BuiltinPlugin) Stop(ctx context.Context) error   { return nil }
func (p *BuiltinPlugin) HealthCheck(ctx context.Context) error { return nil }

// ==================== 插件配置验证 ====================

// ValidatePluginConfig 验证插件配置
func ValidatePluginConfig(configSchema string, config map[string]interface{}) error {
	if configSchema == "" {
		return nil
	}

	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(configSchema), &schema); err != nil {
		return fmt.Errorf("invalid config schema: %w", err)
	}

	// 检查必填字段
	if required, ok := schema["required"].([]interface{}); ok {
		for _, field := range required {
			if fieldName, ok := field.(string); ok {
				if _, exists := config[fieldName]; !exists {
					return fmt.Errorf("required field %s is missing", fieldName)
				}
			}
		}
	}

	return nil
}
