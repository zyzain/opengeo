package plugin

import (
	"context"
	"sync"
)

// DiagnosticPlugin 诊断插件接口
type DiagnosticPlugin interface {
	// Name 返回插件唯一标识
	Name() string

	// Version 返回插件版本
	Version() string

	// Description 返回插件描述
	Description() string

	// SupportedMetrics 返回支持的指标列表
	SupportedMetrics() []string

	// Diagnose 执行诊断
	Diagnose(ctx context.Context, input *DiagnosticInput) (*DiagnosticOutput, error)
}

// DiagnosticInput 诊断输入
type DiagnosticInput struct {
	TenantID     int64    `json:"tenant_id"`
	BrandID      int64    `json:"brand_id"`
	ContentID    int64    `json:"content_id"`
	ContentBody  string   `json:"content_body"`
	Metrics      []string `json:"metrics"`
	TimeRange    *TimeRange `json:"time_range,omitempty"`
}

// TimeRange 时间范围
type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// DiagnosticOutput 诊断输出
type DiagnosticOutput struct {
	PluginName   string           `json:"plugin_name"`
	Score        float64          `json:"score"`
	Dimensions   []DimensionScore `json:"dimensions"`
	Attribution  *Attribution     `json:"attribution,omitempty"`
	Suggestions  []Suggestion     `json:"suggestions"`
}

// DimensionScore 维度评分
type DimensionScore struct {
	Name    string  `json:"name"`
	Score   float64 `json:"score"`
	Weight  float64 `json:"weight"`
	Details string  `json:"details,omitempty"`
}

// Attribution 归因链路
type Attribution struct {
	RootCause    string            `json:"root_cause"`
	TraceID      string            `json:"trace_id"`
	SpanID       string            `json:"span_id"`
	RelatedItems []AttributionItem `json:"related_items"`
}

// AttributionItem 归因项
type AttributionItem struct {
	Type    string `json:"type"`
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Impact  string `json:"impact"`
}

// Suggestion 优化建议
type Suggestion struct {
	Type     string `json:"type"`
	Priority string `json:"priority"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// DiagnosticRegistry 诊断插件注册表
type DiagnosticRegistry struct {
	mu      sync.RWMutex
	plugins map[string]DiagnosticPlugin
}

// NewDiagnosticRegistry 创建诊断插件注册表
func NewDiagnosticRegistry() *DiagnosticRegistry {
	return &DiagnosticRegistry{
		plugins: make(map[string]DiagnosticPlugin),
	}
}

// Register 注册插件
func (r *DiagnosticRegistry) Register(plugin DiagnosticPlugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[plugin.Name()] = plugin
}

// Get 获取插件
func (r *DiagnosticRegistry) Get(name string) (DiagnosticPlugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[name]
	return p, ok
}

// List 列出所有插件
func (r *DiagnosticRegistry) List() []DiagnosticPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	plugins := make([]DiagnosticPlugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// ListByMetric 根据指标列出支持的插件
func (r *DiagnosticRegistry) ListByMetric(metric string) []DiagnosticPlugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var plugins []DiagnosticPlugin
	for _, p := range r.plugins {
		for _, m := range p.SupportedMetrics() {
			if m == metric {
				plugins = append(plugins, p)
				break
			}
		}
	}
	return plugins
}

// 全局诊断插件注册表
var globalDiagnosticRegistry = NewDiagnosticRegistry()

// RegisterDiagnosticPlugin 注册诊断插件到全局注册表
func RegisterDiagnosticPlugin(plugin DiagnosticPlugin) {
	globalDiagnosticRegistry.Register(plugin)
}

// GetDiagnosticPlugin 从全局注册表获取诊断插件
func GetDiagnosticPlugin(name string) (DiagnosticPlugin, bool) {
	return globalDiagnosticRegistry.Get(name)
}

// ListDiagnosticPlugins 列出全局注册表中的所有插件
func ListDiagnosticPlugins() []DiagnosticPlugin {
	return globalDiagnosticRegistry.List()
}
