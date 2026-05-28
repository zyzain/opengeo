# 渠道适配器开发指南

## 概述

渠道适配器用于扩展内容发布渠道，支持自定义发布逻辑。

## 接口定义

```go
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
```

## 开发步骤

### 1. 创建适配器目录

```
service/publish/internal/adapter/your_channel/
├── adapter.go
└── adapter_test.go
```

### 2. 实现适配器接口

```go
package your_channel

import (
    "context"
    "opengeo/pkg/plugin"
)

type YourChannelAdapter struct{}

func NewYourChannelAdapter() *YourChannelAdapter {
    return &YourChannelAdapter{}
}

func (a *YourChannelAdapter) Name() string        { return "your_channel_adapter" }
func (a *YourChannelAdapter) ChannelType() string  { return "your_channel" }
func (a *YourChannelAdapter) Description() string { return "你的渠道适配器" }
func (a *YourChannelAdapter) Version() string     { return "1.0.0" }

func (a *YourChannelAdapter) Publish(ctx context.Context, req *plugin.PublishRequest) (*plugin.PublishResponse, error) {
    // 实现发布逻辑
    return &plugin.PublishResponse{
        ExternalID:  "ext_123",
        ExternalURL: "https://your-channel.com/xxx",
        PublishedAt: "2026-05-28T10:00:00Z",
    }, nil
}

func (a *YourChannelAdapter) Preview(ctx context.Context, req *plugin.PreviewRequest) (*plugin.PreviewResponse, error) {
    return &plugin.PreviewResponse{
        HTML:    "<div>预览</div>",
        Preview: "预览文本",
    }, nil
}

func (a *YourChannelAdapter) GetStatus(ctx context.Context, externalID string) (*plugin.PublishStatus, error) {
    return &plugin.PublishStatus{
        ExternalID: externalID,
        Status:     "published",
        UpdatedAt:  "2026-05-28T10:00:00Z",
    }, nil
}

func (a *YourChannelAdapter) Validate(ctx context.Context, content *plugin.Content) ([]plugin.ValidationIssue, error) {
    // 实现校验逻辑
    return nil, nil
}
```

### 3. 注册适配器

```go
func init() {
    plugin.RegisterChannelAdapter(NewYourChannelAdapter())
}
```

### 4. 导入适配器

在 `service/publish/cmd/main.go` 中导入：

```go
import (
    _ "opengeo/service/publish/internal/adapter/your_channel"
)
```

## 输入/输出结构

### PublishRequest

```go
type PublishRequest struct {
    TenantID int64              `json:"tenant_id"`
    BrandID  int64              `json:"brand_id"`
    Content  *Content           `json:"content"`
    Channel  *Channel           `json:"channel"`
    Schedule string             `json:"schedule,omitempty"`
    Metadata map[string]string  `json:"metadata,omitempty"`
}
```

### PublishResponse

```go
type PublishResponse struct {
    ExternalID  string            `json:"external_id"`
    ExternalURL string            `json:"external_url"`
    PublishedAt string            `json:"published_at"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}
```

### Content

```go
type Content struct {
    ID       int64    `json:"id"`
    Title    string   `json:"title"`
    Body     string   `json:"body"`
    Summary  string   `json:"summary"`
    Tags     []string `json:"tags"`
    CoverURL string   `json:"cover_url"`
}
```

## 测试

```go
func TestYourChannelAdapter_Publish(t *testing.T) {
    adapter := NewYourChannelAdapter()
    
    req := &plugin.PublishRequest{
        TenantID: 1,
        BrandID:  1,
        Content: &plugin.Content{
            Title: "测试标题",
            Body:  "测试内容",
        },
        Channel: &plugin.Channel{
            ID:   1,
            Name: "测试渠道",
        },
    }
    
    resp, err := adapter.Publish(context.Background(), req)
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.ExternalID)
}
```
