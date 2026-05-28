package zhihu

import (
	"context"

	"opengeo/pkg/plugin"
)

type ZhihuAdapter struct{}

func NewZhihuAdapter() *ZhihuAdapter { return &ZhihuAdapter{} }

func (a *ZhihuAdapter) Name() string        { return "zhihu_adapter" }
func (a *ZhihuAdapter) ChannelType() string  { return "zhihu" }
func (a *ZhihuAdapter) Description() string { return "知乎文章/回答发布" }
func (a *ZhihuAdapter) Version() string     { return "1.0.0" }

func (a *ZhihuAdapter) Publish(ctx context.Context, req *plugin.PublishRequest) (*plugin.PublishResponse, error) {
	return &plugin.PublishResponse{ExternalID: "zh_123456", ExternalURL: "https://zhihu.com/xxx", PublishedAt: "2026-05-28T10:00:00Z"}, nil
}

func (a *ZhihuAdapter) Preview(ctx context.Context, req *plugin.PreviewRequest) (*plugin.PreviewResponse, error) {
	return &plugin.PreviewResponse{HTML: "<div>预览</div>", Preview: "预览"}, nil
}

func (a *ZhihuAdapter) GetStatus(ctx context.Context, externalID string) (*plugin.PublishStatus, error) {
	return &plugin.PublishStatus{ExternalID: externalID, Status: "published", UpdatedAt: "2026-05-28T10:00:00Z"}, nil
}

func (a *ZhihuAdapter) Validate(ctx context.Context, content *plugin.Content) ([]plugin.ValidationIssue, error) {
	return nil, nil
}

func init() { plugin.RegisterChannelAdapter(NewZhihuAdapter()) }
