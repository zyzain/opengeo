package douyin

import (
	"context"

	"opengeo/pkg/plugin"
)

type DouyinAdapter struct{}

func NewDouyinAdapter() *DouyinAdapter { return &DouyinAdapter{} }

func (a *DouyinAdapter) Name() string        { return "douyin_adapter" }
func (a *DouyinAdapter) ChannelType() string  { return "douyin" }
func (a *DouyinAdapter) Description() string { return "抖音短视频/图文发布" }
func (a *DouyinAdapter) Version() string     { return "1.0.0" }

func (a *DouyinAdapter) Publish(ctx context.Context, req *plugin.PublishRequest) (*plugin.PublishResponse, error) {
	return &plugin.PublishResponse{ExternalID: "dy_123456", ExternalURL: "https://douyin.com/xxx", PublishedAt: "2026-05-28T10:00:00Z"}, nil
}

func (a *DouyinAdapter) Preview(ctx context.Context, req *plugin.PreviewRequest) (*plugin.PreviewResponse, error) {
	return &plugin.PreviewResponse{HTML: "<div>预览</div>", Preview: "预览"}, nil
}

func (a *DouyinAdapter) GetStatus(ctx context.Context, externalID string) (*plugin.PublishStatus, error) {
	return &plugin.PublishStatus{ExternalID: externalID, Status: "published", UpdatedAt: "2026-05-28T10:00:00Z"}, nil
}

func (a *DouyinAdapter) Validate(ctx context.Context, content *plugin.Content) ([]plugin.ValidationIssue, error) {
	return nil, nil
}

func init() { plugin.RegisterChannelAdapter(NewDouyinAdapter()) }
