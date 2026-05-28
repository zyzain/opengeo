package wechat

import (
	"context"

	"opengeo/pkg/plugin"
)

// WechatAdapter 微信公众号适配器
type WechatAdapter struct{}

func NewWechatAdapter() *WechatAdapter {
	return &WechatAdapter{}
}

func (a *WechatAdapter) Name() string        { return "wechat_adapter" }
func (a *WechatAdapter) ChannelType() string  { return "wechat" }
func (a *WechatAdapter) Description() string { return "微信公众号内容发布" }
func (a *WechatAdapter) Version() string     { return "1.0.0" }

func (a *WechatAdapter) Publish(ctx context.Context, req *plugin.PublishRequest) (*plugin.PublishResponse, error) {
	// 实现微信公众号发布逻辑
	return &plugin.PublishResponse{
		ExternalID:  "wx_123456",
		ExternalURL: "https://mp.weixin.qq.com/s/xxx",
		PublishedAt: "2026-05-28T10:00:00Z",
	}, nil
}

func (a *WechatAdapter) Preview(ctx context.Context, req *plugin.PreviewRequest) (*plugin.PreviewResponse, error) {
	return &plugin.PreviewResponse{
		HTML:    "<div>预览内容</div>",
		Preview: "预览文本",
	}, nil
}

func (a *WechatAdapter) GetStatus(ctx context.Context, externalID string) (*plugin.PublishStatus, error) {
	return &plugin.PublishStatus{
		ExternalID: externalID,
		Status:     "published",
		UpdatedAt:  "2026-05-28T10:00:00Z",
	}, nil
}

func (a *WechatAdapter) Validate(ctx context.Context, content *plugin.Content) ([]plugin.ValidationIssue, error) {
	var issues []plugin.ValidationIssue
	if len(content.Title) > 64 {
		issues = append(issues, plugin.ValidationIssue{
			Field:       "title",
			IssueType:   "length",
			Description: "标题长度超过限制",
			Suggestion:  "请将标题缩短到64个字符以内",
			Severity:    "error",
		})
	}
	return issues, nil
}

func init() {
	plugin.RegisterChannelAdapter(NewWechatAdapter())
}
