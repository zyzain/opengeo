# 诊断插件开发指南

## 概述

诊断插件用于扩展品牌诊断能力，支持自定义指标和评分逻辑。

## 接口定义

```go
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
```

## 开发步骤

### 1. 创建插件目录

```
service/monitor/internal/adapter/builtin/your_plugin/
├── plugin.go
└── plugin_test.go
```

### 2. 实现插件接口

```go
package your_plugin

import (
    "context"
    "opengeo/pkg/plugin"
)

type YourPlugin struct{}

func NewYourPlugin() *YourPlugin {
    return &YourPlugin{}
}

func (p *YourPlugin) Name() string        { return "your_plugin" }
func (p *YourPlugin) Version() string     { return "1.0.0" }
func (p *YourPlugin) Description() string { return "你的插件描述" }
func (p *YourPlugin) SupportedMetrics() []string {
    return []string{"metric1", "metric2"}
}

func (p *YourPlugin) Diagnose(ctx context.Context, input *plugin.DiagnosticInput) (*plugin.DiagnosticOutput, error) {
    // 实现诊断逻辑
    return &plugin.DiagnosticOutput{
        PluginName: p.Name(),
        Score:      85.0,
        Dimensions: []plugin.DimensionScore{
            {Name: "metric1", Score: 90.0, Weight: 0.5},
            {Name: "metric2", Score: 80.0, Weight: 0.5},
        },
    }, nil
}
```

### 3. 注册插件

在 `init()` 函数中注册插件：

```go
func init() {
    plugin.RegisterDiagnosticPlugin(NewYourPlugin())
}
```

### 4. 导入插件

在 `service/monitor/cmd/main.go` 中导入插件包：

```go
import (
    _ "opengeo/service/monitor/internal/adapter/builtin/your_plugin"
)
```

## 输入/输出结构

### DiagnosticInput

```go
type DiagnosticInput struct {
    TenantID    int64      `json:"tenant_id"`
    BrandID     int64      `json:"brand_id"`
    ContentID   int64      `json:"content_id"`
    ContentBody string     `json:"content_body"`
    Metrics     []string   `json:"metrics"`
    TimeRange   *TimeRange `json:"time_range,omitempty"`
}
```

### DiagnosticOutput

```go
type DiagnosticOutput struct {
    PluginName  string           `json:"plugin_name"`
    Score       float64          `json:"score"`
    Dimensions  []DimensionScore `json:"dimensions"`
    Attribution *Attribution     `json:"attribution,omitempty"`
    Suggestions []Suggestion     `json:"suggestions"`
}
```

## 测试

```go
func TestYourPlugin_Diagnose(t *testing.T) {
    plugin := NewYourPlugin()
    
    input := &plugin.DiagnosticInput{
        TenantID:  1,
        BrandID:   1,
        ContentID: 1,
        Metrics:   []string{"metric1", "metric2"},
    }
    
    output, err := plugin.Diagnose(context.Background(), input)
    assert.NoError(t, err)
    assert.GreaterOrEqual(t, output.Score, 0.0)
    assert.LessOrEqual(t, output.Score, 100.0)
}
```
