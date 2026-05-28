package service

import (
	"context"
	"encoding/json"
	"testing"

	"opengeo/service/publish/internal/domain/model"
)

func newTestValidationChannel() *model.Channel {
	geoConfig := &model.ChannelGEOConfig{
		InjectSchemaMarkup: true,
		InjectAIGCLabel:    true,
		MaxTitleLength:     64,
		MaxBodyLength:      5000,
	}
	configJSON, _ := json.Marshal(geoConfig)

	return &model.Channel{
		ID:         1,
		ChannelType: "wechat",
		GEOConfig:  string(configJSON),
	}
}

func TestValidateForPublish_ValidConfig(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title:        "这是一个有效的测试标题",
		Body:         "## 测试标题\n\n这是一段有效的测试正文内容，长度足够通过校验。包含足够的文字来满足最低长度要求。",
		Tags:         []string{"测试", "GEO"},
		SchemaMarkup: `{"@context":"https://schema.org","@type":"Article"}`,
		AIGCLabel:    "本文由AI辅助生成。",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
	if result.Score < 80 {
		t.Errorf("expected score >= 80, got %f", result.Score)
	}
}

func TestValidateForPublish_EmptyTitle(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "",
		Body:  "测试正文",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	if result.Valid {
		t.Error("expected invalid for empty title")
	}

	found := false
	for _, err := range result.Errors {
		if err.Field == "title" && err.Rule == "required" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected title required error")
	}
}

func TestValidateForPublish_EmptyBody(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "测试标题",
		Body:  "",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	if result.Valid {
		t.Error("expected invalid for empty body")
	}
}

func TestValidateForPublish_ShortTitle(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "短",
		Body:  "测试正文内容",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, warn := range result.Warnings {
		if warn.Field == "title" && warn.Rule == "min_length" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected short title warning")
	}
}

func TestValidateForPublish_LongTitle(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	longTitle := ""
	for i := 0; i < 100; i++ {
		longTitle += "测试"
	}

	config := &model.PublishConfig{
		Title: longTitle,
		Body:  "测试正文",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, warn := range result.Warnings {
		if warn.Field == "title" && warn.Rule == "max_length" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected long title warning")
	}
}

func TestValidateForPublish_ShortBody(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "测试标题",
		Body:  "短内容",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, warn := range result.Warnings {
		if warn.Field == "body" && warn.Rule == "min_length" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected short body warning")
	}
}

func TestValidateForPublish_InvalidSchemaMarkup(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title:        "测试标题",
		Body:         "测试正文内容",
		SchemaMarkup: "invalid json",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, err := range result.Errors {
		if err.Field == "schema_markup" && err.Rule == "json_format" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected invalid schema error")
	}
}

func TestValidateForPublish_MissingSchema(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "测试标题",
		Body:  "测试正文",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, warn := range result.Warnings {
		if warn.Field == "schema_markup" && warn.Rule == "missing" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected missing schema warning")
	}
}

func TestValidateForPublish_NoTags(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "测试标题",
		Body:  "测试正文内容",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, info := range result.Infos {
		if info.Field == "tags" && info.Rule == "empty" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected empty tags info")
	}
}

func TestValidateForPublish_TooManyTags(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	tags := make([]string, 15)
	for i := range tags {
		tags[i] = "tag"
	}

	config := &model.PublishConfig{
		Title: "测试标题",
		Body:  "测试正文内容",
		Tags:  tags,
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, warn := range result.Warnings {
		if warn.Field == "tags" && warn.Rule == "too_many" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected too many tags warning")
	}
}

func TestValidateForPublish_InvalidCoverURL(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title:    "测试标题",
		Body:     "测试正文内容",
		CoverURL: "invalid-url",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, err := range result.Errors {
		if err.Field == "cover_url" && err.Rule == "url_format" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected invalid cover URL error")
	}
}

func TestValidateForPublish_ImageWithoutALT(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "测试标题",
		Body:  `## 测试\n\n<img src="test.jpg">\n\n正文内容`,
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, warn := range result.Warnings {
		if warn.Rule == "img_alt" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected image ALT warning")
	}
}

func TestValidateForPublish_NoHeadings(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	longBody := ""
	for i := 0; i < 100; i++ {
		longBody += "这是一段没有标题的长文本内容。"
	}

	config := &model.PublishConfig{
		Title: "测试标题",
		Body:  longBody,
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	found := false
	for _, info := range result.Infos {
		if info.Rule == "structure" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected structure info")
	}
}

func TestValidateForPublish_ScoreCalculation(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title: "",
		Body:  "",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	if result.Score >= 100 {
		t.Error("expected score < 100 for invalid config")
	}
	if result.Valid {
		t.Error("expected invalid")
	}
}

func TestValidateForPublish_CountCheckItems(t *testing.T) {
	svc := NewValidationService()
	ctx := context.Background()
	channel := newTestValidationChannel()

	config := &model.PublishConfig{
		Title:        "这是一个有效的测试标题",
		Body:         "## 测试\n\n有效正文内容",
		Tags:         []string{"测试"},
		SchemaMarkup: `{"@context":"https://schema.org"}`,
		AIGCLabel:    "AIGC",
		CoverURL:     "https://example.com/cover.jpg",
	}

	result := svc.ValidateForPublish(ctx, channel, config)

	totalItems := len(result.Errors) + len(result.Warnings) + len(result.Infos)
	t.Logf("Total check items: %d (errors: %d, warnings: %d, infos: %d)",
		totalItems, len(result.Errors), len(result.Warnings), len(result.Infos))
}

func TestFormatValidationResult(t *testing.T) {
	svc := NewValidationService()

	result := &ValidationResult{
		Valid: false,
		Errors: []ValidationItem{
			{Field: "title", Rule: "required", Message: "标题不能为空"},
		},
		Warnings: []ValidationItem{
			{Field: "body", Rule: "min_length", Message: "正文过短"},
		},
		Infos: []ValidationItem{
			{Field: "tags", Rule: "empty", Message: "未设置标签"},
		},
		Score: 80,
	}

	output := svc.FormatValidationResult(result)

	if output == "" {
		t.Error("expected non-empty output")
	}
	if !containsStr(output, "未通过") {
		t.Error("expected status in output")
	}
	if !containsStr(output, "标题不能为空") {
		t.Error("expected error in output")
	}
}
