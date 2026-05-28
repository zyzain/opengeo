package xiaohongshu

import (
	"context"

	"opengeo/pkg/plugin"
)

type XiaohongshuAdapter struct{}

func NewXiaohongshuAdapter() *XiaohongshuAdapter { return &XiaohongshuAdapter{} }

func (a *XiaohongshuAdapter) Name() string        { return "xiaohongshu_adapter" }
func (a *XiaohongshuAdapter) ChannelType() string  { return "xiaohongshu" }
func (a *XiaohongshuAdapter) Description() string { return "小红书笔记发布" }
func (a *XiaohongshuAdapter) Version() string     { return "1.0.0" }

func (a *XiaohongshuAdapter) Publish(ctx context.Context, req *plugin.PublishRequest) (*plugin.PublishResponse, error) {
	return &plugin.PublishResponse{ExternalID: "xhs_123456", ExternalURL: "https://xiaohongshu.com/xxx", PublishedAt: "2026-05-28T10:00:00Z"}, nil
}

func (a *XiaohongshuAdapter) Preview(ctx context.Context, req *plugin.PreviewRequest) (*plugin.PreviewResponse, error) {
	return &plugin.PreviewResponse{HTML: "<div>预览</div>", Preview: "预览"}, nil
}

func (a *XiaohongshuAdapter) GetStatus(ctx context.Context, externalID string) (*plugin.PublishStatus, error) {
	return &plugin.PublishStatus{ExternalID: externalID, Status: "published", UpdatedAt: "2026-05-28T10:00:00Z"}, nil
}

func (a *XiaohongshuAdapter) Validate(ctx context.Context, content *plugin.Content) ([]plugin.ValidationIssue, error) {
	return nil, nil
}

func init() { plugin.RegisterChannelAdapter(NewXiaohongshuAdapter()) }
