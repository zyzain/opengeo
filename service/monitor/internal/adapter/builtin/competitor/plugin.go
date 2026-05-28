package competitor

import (
	"context"

	"opengeo/pkg/plugin"
)

// CompetitorMonitorPlugin 竞品监测插件
type CompetitorMonitorPlugin struct{}

func NewCompetitorMonitorPlugin() *CompetitorMonitorPlugin {
	return &CompetitorMonitorPlugin{}
}

func (p *CompetitorMonitorPlugin) Name() string        { return "competitor_monitor" }
func (p *CompetitorMonitorPlugin) Version() string     { return "1.0.0" }
func (p *CompetitorMonitorPlugin) Description() string { return "监测竞品动态和差距" }
func (p *CompetitorMonitorPlugin) SupportedMetrics() []string {
	return []string{"visibility", "content_gap", "market_share"}
}

func (p *CompetitorMonitorPlugin) Diagnose(ctx context.Context, input *plugin.DiagnosticInput) (*plugin.DiagnosticOutput, error) {
	return &plugin.DiagnosticOutput{
		PluginName: p.Name(),
		Score:      75.0,
		Dimensions: []plugin.DimensionScore{
			{Name: "visibility", Score: 70.0, Weight: 0.4},
			{Name: "content_gap", Score: 80.0, Weight: 0.3},
			{Name: "market_share", Score: 75.0, Weight: 0.3},
		},
		Suggestions: []plugin.Suggestion{
			{Type: "content", Priority: "high", Title: "填补内容差距", Content: "竞品在以下领域有更多内容覆盖"},
		},
	}, nil
}

func init() {
	plugin.RegisterDiagnosticPlugin(NewCompetitorMonitorPlugin())
}
