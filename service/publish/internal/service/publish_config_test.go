package service

import (
	"context"
	"encoding/json"
	"testing"

	"opengeo/service/publish/internal/domain/model"
)

func newTestChannel() *model.Channel {
	geoConfig := &model.ChannelGEOConfig{
		InjectHiddenMarkup: true,
		HiddenMarkupFormat: "html_comment",
		InjectSchemaMarkup: true,
		SchemaType:         "Article",
		MaxTitleLength:     64,
		MaxBodyLength:      5000,
		InjectAIGCLabel:    true,
		AIGCLabelTemplate:  "本文由AI辅助生成（AIGC）。",
		Variables:          map[string]string{"brand": "OpenGEO"},
	}
	configJSON, _ := json.Marshal(geoConfig)

	return &model.Channel{
		ID:            1,
		UserID:        1,
		ChannelType:   "wechat",
		ChannelName:   "微信公众号",
		TitleTemplate: "{{title}} - {{brand}}",
		BodyTemplate:  "{{body}}\n\n来源：{{channel_name}}",
		TagsTemplate:  "{{tags}},GEO",
		GEOConfig:     string(configJSON),
		IsEnabled:     true,
	}
}

func TestPreparePublishConfig_Basic(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	config, err := svc.PreparePublishConfig(ctx, channel,
		"测试标题", "测试正文内容", []string{"tag1", "tag2"}, "https://example.com/cover.jpg", nil)

	if err != nil {
		t.Fatalf("PreparePublishConfig failed: %v", err)
	}

	if config.Title == "" {
		t.Error("expected non-empty title")
	}
	if config.Body == "" {
		t.Error("expected non-empty body")
	}
	if len(config.Tags) == 0 {
		t.Error("expected tags")
	}
}

func TestPreparePublishConfig_VariableReplacement(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	config, err := svc.PreparePublishConfig(ctx, channel,
		"测试标题", "测试正文", []string{"tag1"}, "", nil)

	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// 标题模板是 "{{title}} - {{brand}}"
	if config.Title != "测试标题 - OpenGEO" {
		t.Errorf("title variable replacement failed: %s", config.Title)
	}

	// 正文模板是 "{{body}}\n\n来源：{{channel_name}}"
	if config.Body != "测试正文\n\n来源：微信公众号" {
		t.Errorf("body variable replacement failed: %s", config.Body)
	}
}

func TestPreparePublishConfig_CustomVariables(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	customVars := map[string]string{
		"brand": "CustomBrand",
	}

	config, err := svc.PreparePublishConfig(ctx, channel,
		"标题", "正文", []string{}, "", customVars)

	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// 自定义变量应覆盖渠道配置变量
	if config.Title != "标题 - CustomBrand" {
		t.Errorf("custom variable not applied: %s", config.Title)
	}
}

func TestPreparePublishConfig_HiddenMarkup(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	config, err := svc.PreparePublishConfig(ctx, channel,
		"测试标题", "测试正文", []string{}, "", nil)

	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if config.HiddenMarkup == "" {
		t.Error("expected hidden markup")
	}
	if !containsStr(config.HiddenMarkup, "<!-- GEO:") {
		t.Error("expected HTML comment format")
	}
}

func TestPreparePublishConfig_SchemaMarkup(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	config, err := svc.PreparePublishConfig(ctx, channel,
		"测试标题", "测试正文", []string{}, "", nil)

	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if config.SchemaMarkup == "" {
		t.Error("expected schema markup")
	}
	if !containsStr(config.SchemaMarkup, "schema.org") {
		t.Error("expected schema.org in markup")
	}
	if !containsStr(config.SchemaMarkup, "Article") {
		t.Error("expected Article type")
	}
}

func TestPreparePublishConfig_AIGCLabel(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	config, err := svc.PreparePublishConfig(ctx, channel,
		"标题", "正文", []string{}, "", nil)

	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if config.AIGCLabel == "" {
		t.Error("expected AIGC label")
	}
	if !containsStr(config.AIGCLabel, "AIGC") {
		t.Error("expected AIGC in label")
	}
}

func TestPreparePublishConfig_TitleLengthLimit(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	// 创建超长标题
	longTitle := ""
	for i := 0; i < 100; i++ {
		longTitle += "测试"
	}

	config, err := svc.PreparePublishConfig(ctx, channel,
		longTitle, "正文", []string{}, "", nil)

	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// 标题应该被截断
	if len(config.Title) > 100 { // 64 + 模板后缀
		t.Errorf("title not truncated: len=%d", len(config.Title))
	}
}

func TestPreviewPublish(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	config := &model.PublishConfig{
		Title:        "测试标题",
		Body:         "测试正文",
		Tags:         []string{"tag1", "tag2"},
		SchemaMarkup: `{"@context":"https://schema.org"}`,
		HiddenMarkup: "<!-- GEO: 测试 -->",
		AIGCLabel:    "本文由AI辅助生成。",
	}

	preview, err := svc.PreviewPublish(ctx, channel, config)
	if err != nil {
		t.Fatalf("PreviewPublish failed: %v", err)
	}

	if preview == "" {
		t.Error("expected non-empty preview")
	}
	if !containsStr(preview, "测试标题") {
		t.Error("expected title in preview")
	}
	if !containsStr(preview, "测试正文") {
		t.Error("expected body in preview")
	}
	if !containsStr(preview, "tag1") {
		t.Error("expected tags in preview")
	}
}

func TestValidatePublishConfig(t *testing.T) {
	svc := NewPublishConfigService()
	ctx := context.Background()
	channel := newTestChannel()

	tests := []struct {
		config    *model.PublishConfig
		wantErrs  int
		wantWarns int
		desc      string
	}{
		{
			config:    &model.PublishConfig{Title: "标题", Body: "正文", SchemaMarkup: "{}", AIGCLabel: "label"},
			wantErrs:  0,
			wantWarns: 0,
			desc:      "valid config",
		},
		{
			config:    &model.PublishConfig{Title: "", Body: "正文"},
			wantErrs:  1,
			wantWarns: 2, // missing schema + aigc
			desc:      "missing title",
		},
		{
			config:    &model.PublishConfig{Title: "标题", Body: ""},
			wantErrs:  1,
			wantWarns: 2,
			desc:      "missing body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			errs, warns := svc.ValidatePublishConfig(ctx, channel, tt.config)
			if len(errs) != tt.wantErrs {
				t.Errorf("errors: got %d, want %d: %v", len(errs), tt.wantErrs, errs)
			}
			if len(warns) != tt.wantWarns {
				t.Errorf("warnings: got %d, want %d: %v", len(warns), tt.wantWarns, warns)
			}
		})
	}
}

func TestUpdateChannelGEOConfig(t *testing.T) {
	svc := NewPublishConfigService()
	channel := &model.Channel{}

	config := &model.ChannelGEOConfig{
		InjectHiddenMarkup: true,
		InjectSchemaMarkup: true,
		SchemaType:         "FAQ",
	}

	err := svc.UpdateChannelGEOConfig(channel, config)
	if err != nil {
		t.Fatalf("UpdateChannelGEOConfig failed: %v", err)
	}

	if channel.GEOConfig == "" {
		t.Error("expected GEO config to be set")
	}

	// 解析并验证
	var parsed model.ChannelGEOConfig
	json.Unmarshal([]byte(channel.GEOConfig), &parsed)

	if !parsed.InjectHiddenMarkup {
		t.Error("expected inject_hidden_markup=true")
	}
	if parsed.SchemaType != "FAQ" {
		t.Errorf("expected schema_type=FAQ, got %s", parsed.SchemaType)
	}
}

func TestUpdateChannelTemplates(t *testing.T) {
	svc := NewPublishConfigService()
	channel := &model.Channel{}

	svc.UpdateChannelTemplates(channel, "new title", "new body", "new tags", "new cover")

	if channel.TitleTemplate != "new title" {
		t.Errorf("title template not updated: %s", channel.TitleTemplate)
	}
	if channel.BodyTemplate != "new body" {
		t.Errorf("body template not updated: %s", channel.BodyTemplate)
	}
}

func TestGetChannelGEOConfig(t *testing.T) {
	svc := NewPublishConfigService()

	// 空配置返回默认值
	channel := &model.Channel{}
	config := svc.GetChannelGEOConfig(channel)

	if config == nil {
		t.Fatal("expected non-nil config")
	}
	if config.MaxTitleLength != 64 {
		t.Errorf("expected default max title length 64, got %d", config.MaxTitleLength)
	}
}

func TestExtractKeywords(t *testing.T) {
	svc := NewPublishConfigService()

	tests := []struct {
		input string
		want  int
		desc  string
	}{
		{"GEO优化完全指南", 4, "chinese keywords"},
		{"Hello World Test", 3, "english keywords"},
		{"a b c", 0, "too short"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := svc.extractKeywords(tt.input)
			if result == "" && tt.want > 0 {
				t.Error("expected non-empty result")
			}
		})
	}
}

func TestTruncateText(t *testing.T) {
	svc := NewPublishConfigService()

	tests := []struct {
		text   string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello..."},
		{"你好世界测试", 3, "你好世..."},
	}

	for _, tt := range tests {
		result := svc.truncateText(tt.text, tt.maxLen)
		runes := []rune(result)
		if len(runes) > tt.maxLen+3 { // +3 for "..."
			t.Errorf("truncate(%q, %d) = %q, too long", tt.text, tt.maxLen, result)
		}
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
