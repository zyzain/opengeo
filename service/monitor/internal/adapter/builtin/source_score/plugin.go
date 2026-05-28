package source_score

import (
	"context"

	"opengeo/pkg/plugin"
)

// SourceScorePlugin 信源评分插件
type SourceScorePlugin struct{}

// NewSourceScorePlugin 创建信源评分插件
func NewSourceScorePlugin() *SourceScorePlugin {
	return &SourceScorePlugin{}
}

// Name 返回插件名称
func (p *SourceScorePlugin) Name() string {
	return "source_score"
}

// Version 返回插件版本
func (p *SourceScorePlugin) Version() string {
	return "1.0.0"
}

// Description 返回插件描述
func (p *SourceScorePlugin) Description() string {
	return "评估渠道和账号的权威度"
}

// SupportedMetrics 返回支持的指标列表
func (p *SourceScorePlugin) SupportedMetrics() []string {
	return []string{"authority", "freshness", "relevance"}
}

// Diagnose 执行诊断
func (p *SourceScorePlugin) Diagnose(ctx context.Context, input *plugin.DiagnosticInput) (*plugin.DiagnosticOutput, error) {
	// 获取租户 ID
	tenantID := input.TenantID

	// 这里应该实现实际的信源评分逻辑
	// 简化实现，返回示例数据
	output := &plugin.DiagnosticOutput{
		PluginName: p.Name(),
		Score:      85.5,
		Dimensions: []plugin.DimensionScore{
			{
				Name:   "authority",
				Score:  90.0,
				Weight: 0.4,
			},
			{
				Name:   "freshness",
				Score:  80.0,
				Weight: 0.3,
			},
			{
				Name:   "relevance",
				Score:  85.0,
				Weight: 0.3,
			},
		},
		Suggestions: []plugin.Suggestion{
			{
				Type:     "content",
				Priority: "medium",
				Title:    "提升内容时效性",
				Content:  "建议定期更新内容，保持内容的新鲜度",
			},
		},
	}

	return output, nil
}

func init() {
	// 注册插件到全局注册表
	plugin.RegisterDiagnosticPlugin(NewSourceScorePlugin())
}
