package weibo

import (
	"context"

	"opengeo/pkg/plugin"
)

type WeiboAdapter struct{}

func NewWeiboAdapter() *WeiboAdapter { return &WeiboAdapter{} }

func (a *WeiboAdapter) Name() string        { return "weibo_adapter" }
func (a *WeiboAdapter) ChannelType() string  { return "weibo" }
func (a *WeiboAdapter) Description() string { return "微博内容发布" }
func (a *WeiboAdapter) Version() string     { return "1.0.0" }

func (a *WeiboAdapter) Publish(ctx context.Context, req *plugin.PublishRequest) (*plugin.PublishResponse, error) {
	return &plugin.PublishResponse{
		ExternalID:  "wb_123456",
		ExternalURL: "https://weibo.com/xxx",
		PublishedAt: "2026-05-28T10:00:00Z",
	}, nil
}

func (a *WeiboAdapter) Preview(ctx context.Context, req *plugin.PreviewRequest) (*plugin.PreviewResponse, error) {
	return &plugin.PreviewResponse{HTML: "<div>预览</div>", Preview: "预览"}, nil
}

func (a *WeiboAdapter) GetStatus(ctx context.Context, externalID string) (*plugin.PublishStatus, error) {
	return &plugin.PublishStatus{ExternalID: externalID, Status: "published", UpdatedAt: "2026-05-28T10:00:00Z"}, nil
}

func (a *WeiboAdapter) Validate(ctx context.Context, content *plugin.Content) ([]plugin.ValidationIssue, error) {
	return nil, nil
}

func init() { plugin.RegisterChannelAdapter(NewWeiboAdapter()) }
