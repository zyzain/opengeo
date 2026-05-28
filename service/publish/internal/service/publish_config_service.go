package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"opengeo/service/publish/internal/domain/model"
)

var (
	reVariable     = regexp.MustCompile(`\{\{(\w+)\}\}`)
	reHTMLTag      = regexp.MustCompile(`<[^>]+>`)
)

// PublishConfigService GEO发布配置服务
type PublishConfigService struct{}

// NewPublishConfigService 创建GEO发布配置服务
func NewPublishConfigService() *PublishConfigService {
	return &PublishConfigService{}
}

// PreparePublishConfig 准备发布配置
func (s *PublishConfigService) PreparePublishConfig(
	ctx context.Context,
	channel *model.Channel,
	title, body string,
	tags []string,
	coverURL string,
	customVars map[string]string,
) (*model.PublishConfig, error) {

	// 解析渠道GEO配置
	geoConfig := s.parseGEOConfig(channel.GEOConfig)

	// 合并变量
	variables := s.buildVariables(channel, title, body, tags, coverURL, customVars)

	// 应用模板和变量替换
	finalTitle := s.applyTemplate(channel.TitleTemplate, title, variables)
	finalBody := s.applyTemplate(channel.BodyTemplate, body, variables)
	finalTags := s.applyTagsTemplate(channel.TagsTemplate, tags, variables)
	finalCover := s.applyTemplate(channel.CoverTemplate, coverURL, variables)

	// 限制长度
	if geoConfig.MaxTitleLength > 0 && len(finalTitle) > geoConfig.MaxTitleLength {
		finalTitle = finalTitle[:geoConfig.MaxTitleLength] + "..."
	}
	if geoConfig.MaxBodyLength > 0 && len([]rune(finalBody)) > geoConfig.MaxBodyLength {
		runes := []rune(finalBody)
		finalBody = string(runes[:geoConfig.MaxBodyLength]) + "..."
	}

	// 生成隐藏语义标记
	hiddenMarkup := ""
	if geoConfig.InjectHiddenMarkup {
		hiddenMarkup = s.generateHiddenMarkup(geoConfig, title, body)
	}

	// 生成Schema Markup
	schemaMarkup := ""
	if geoConfig.InjectSchemaMarkup {
		schemaMarkup = s.generateSchemaMarkup(geoConfig, title, body)
	}

	// 生成AIGC标识
	aigcLabel := ""
	if geoConfig.InjectAIGCLabel {
		aigcLabel = s.generateAIGCLabel(geoConfig)
	}

	return &model.PublishConfig{
		Title:        finalTitle,
		Body:         finalBody,
		Tags:         finalTags,
		CoverURL:     finalCover,
		SchemaMarkup: schemaMarkup,
		HiddenMarkup: hiddenMarkup,
		AIGCLabel:    aigcLabel,
		Variables:    variables,
		ChannelConfig: geoConfig,
	}, nil
}

// PreviewPublish 发布预览
func (s *PublishConfigService) PreviewPublish(
	ctx context.Context,
	channel *model.Channel,
	config *model.PublishConfig,
) (string, error) {

	geoConfig := s.parseGEOConfig(channel.GEOConfig)

	var sb strings.Builder

	// 标题
	sb.WriteString(fmt.Sprintf("# %s\n\n", config.Title))

	// Schema Markup（如果启用）
	if config.SchemaMarkup != "" && geoConfig.InjectSchemaMarkup {
		sb.WriteString("<!-- Schema Markup -->\n")
		sb.WriteString(fmt.Sprintf("<script type=\"application/ld+json\">\n%s\n</script>\n\n", config.SchemaMarkup))
	}

	// 隐藏语义标记（如果启用）
	if config.HiddenMarkup != "" {
		sb.WriteString(config.HiddenMarkup)
		sb.WriteString("\n\n")
	}

	// 正文
	sb.WriteString(config.Body)
	sb.WriteString("\n\n")

	// AIGC标识
	if config.AIGCLabel != "" {
		sb.WriteString("---\n")
		sb.WriteString(config.AIGCLabel)
		sb.WriteString("\n")
	}

	// 标签
	if len(config.Tags) > 0 {
		sb.WriteString("\n标签: ")
		sb.WriteString(strings.Join(config.Tags, ", "))
	}

	return sb.String(), nil
}

// ValidatePublishConfig 校验发布配置
func (s *PublishConfigService) ValidatePublishConfig(
	ctx context.Context,
	channel *model.Channel,
	config *model.PublishConfig,
) ([]string, []string) {

	errors := make([]string, 0)
	warnings := make([]string, 0)

	geoConfig := s.parseGEOConfig(channel.GEOConfig)

	// 标题校验
	if config.Title == "" {
		errors = append(errors, "标题不能为空")
	}
	if geoConfig.MaxTitleLength > 0 && len(config.Title) > geoConfig.MaxTitleLength {
		warnings = append(warnings, fmt.Sprintf("标题超过%d字符，将被截断", geoConfig.MaxTitleLength))
	}

	// 正文校验
	if config.Body == "" {
		errors = append(errors, "正文不能为空")
	}
	if geoConfig.MaxBodyLength > 0 && len([]rune(config.Body)) > geoConfig.MaxBodyLength {
		warnings = append(warnings, fmt.Sprintf("正文超过%d字符，将被截断", geoConfig.MaxBodyLength))
	}

	// Schema Markup校验
	if geoConfig.InjectSchemaMarkup && config.SchemaMarkup == "" {
		warnings = append(warnings, "已启用Schema Markup注入但未生成Schema数据")
	}

	// AIGC标识校验
	if geoConfig.InjectAIGCLabel && config.AIGCLabel == "" {
		warnings = append(warnings, "已启用AIGC标识但未生成标识内容")
	}

	return errors, warnings
}

// ==================== 内部方法 ====================

// parseGEOConfig 解析GEO配置
func (s *PublishConfigService) parseGEOConfig(configJSON string) *model.ChannelGEOConfig {
	config := &model.ChannelGEOConfig{
		InjectHiddenMarkup: false,
		HiddenMarkupFormat: "html_comment",
		InjectSchemaMarkup: false,
		SchemaType:         "article",
		MaxTitleLength:     64,
		MaxBodyLength:      5000,
		InjectAIGCLabel:    true,
		AIGCLabelTemplate:  "本文由AI辅助生成（AIGC），仅供参考。",
		Variables:          make(map[string]string),
	}

	if configJSON == "" {
		return config
	}

	json.Unmarshal([]byte(configJSON), config)
	return config
}

// buildVariables 构建变量映射
func (s *PublishConfigService) buildVariables(
	channel *model.Channel,
	title, body string,
	tags []string,
	coverURL string,
	customVars map[string]string,
) map[string]string {

	variables := map[string]string{
		"title":        title,
		"body":         body,
		"tags":         strings.Join(tags, ","),
		"cover_url":    coverURL,
		"channel_name": channel.ChannelName,
		"channel_type": channel.ChannelType,
		"date":         time.Now().Format("2006-01-02"),
		"time":         time.Now().Format("15:04:05"),
		"datetime":     time.Now().Format("2006-01-02 15:04:05"),
		"timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
	}

	// 合并渠道配置中的变量
	geoConfig := s.parseGEOConfig(channel.GEOConfig)
	for k, v := range geoConfig.Variables {
		variables[k] = v
	}

	// 合并自定义变量（最高优先级）
	for k, v := range customVars {
		variables[k] = v
	}

	return variables
}

// applyTemplate 应用模板变量替换
func (s *PublishConfigService) applyTemplate(templateStr, defaultValue string, variables map[string]string) string {
	if templateStr == "" {
		return defaultValue
	}

	result := reVariable.ReplaceAllStringFunc(templateStr, func(match string) string {
		varName := match[2 : len(match)-2] // 去掉 {{ 和 }}
		if val, ok := variables[varName]; ok {
			return val
		}
		return match
	})

	return result
}

// applyTagsTemplate 应用标签模板
func (s *PublishConfigService) applyTagsTemplate(templateStr string, tags []string, variables map[string]string) []string {
	if templateStr == "" {
		return tags
	}

	// 替换变量
	result := reVariable.ReplaceAllStringFunc(templateStr, func(match string) string {
		varName := match[2 : len(match)-2]
		if varName == "tags" {
			return strings.Join(tags, ",")
		}
		if val, ok := variables[varName]; ok {
			return val
		}
		return match
	})

	// 分割标签
	return strings.Split(result, ",")
}

// generateHiddenMarkup 生成隐藏语义标记
func (s *PublishConfigService) generateHiddenMarkup(config *model.ChannelGEOConfig, title, body string) string {
	prefix := config.HiddenMarkupPrefix
	if prefix == "" {
		prefix = "GEO"
	}

	switch config.HiddenMarkupFormat {
	case "html_comment":
		return fmt.Sprintf("<!-- %s: %s --><!-- %s-Keywords: %s -->",
			prefix, title, prefix, s.extractKeywords(title))
	case "microdata":
		return fmt.Sprintf(`<span itemscope itemtype="https://schema.org/Article">
<meta itemprop="headline" content="%s"/>
<meta itemprop="description" content="%s"/>
</span>`, title, s.truncateText(body, 200))
	case "json_ld":
		return fmt.Sprintf(`<script type="application/ld+json">
{"@context":"https://schema.org","@type":"Article","headline":"%s","description":"%s"}
</script>`, title, s.truncateText(body, 200))
	default:
		return fmt.Sprintf("<!-- %s: %s -->", prefix, title)
	}
}

// generateSchemaMarkup 生成Schema Markup
func (s *PublishConfigService) generateSchemaMarkup(config *model.ChannelGEOConfig, title, body string) string {
	schemaType := config.SchemaType
	if schemaType == "" {
		schemaType = "Article"
	}

	description := s.truncateText(body, 200)
	now := time.Now().Format("2006-01-02")

	return fmt.Sprintf(`{
  "@context": "https://schema.org",
  "@type": "%s",
  "headline": "%s",
  "description": "%s",
  "datePublished": "%s",
  "dateModified": "%s"
}`, schemaType, title, description, now, now)
}

// generateAIGCLabel 生成AIGC标识
func (s *PublishConfigService) generateAIGCLabel(config *model.ChannelGEOConfig) string {
	if config.AIGCLabelTemplate != "" {
		return config.AIGCLabelTemplate
	}
	return "本文由AI辅助生成（AIGC），仅供参考。"
}

// extractKeywords 提取关键词
func (s *PublishConfigService) extractKeywords(text string) string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ',' || r == '，' || r == '。' || r == '、'
	})

	keywords := make([]string, 0)
	for _, w := range words {
		w = strings.TrimSpace(w)
		if len([]rune(w)) >= 2 {
			keywords = append(keywords, w)
		}
	}

	if len(keywords) > 5 {
		keywords = keywords[:5]
	}

	return strings.Join(keywords, ", ")
}

// truncateText 截断文本
func (s *PublishConfigService) truncateText(text string, maxLen int) string {
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	return string(runes[:maxLen]) + "..."
}

// ==================== 渠道配置管理 ====================

// UpdateChannelGEOConfig 更新渠道GEO配置
func (s *PublishConfigService) UpdateChannelGEOConfig(channel *model.Channel, config *model.ChannelGEOConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal geo config: %w", err)
	}
	channel.GEOConfig = string(configJSON)
	channel.UpdatedAt = time.Now()
	return nil
}

// UpdateChannelTemplates 更新渠道模板
func (s *PublishConfigService) UpdateChannelTemplates(
	channel *model.Channel,
	titleTemplate, bodyTemplate, tagsTemplate, coverTemplate string,
) {
	if titleTemplate != "" {
		channel.TitleTemplate = titleTemplate
	}
	if bodyTemplate != "" {
		channel.BodyTemplate = bodyTemplate
	}
	if tagsTemplate != "" {
		channel.TagsTemplate = tagsTemplate
	}
	if coverTemplate != "" {
		channel.CoverTemplate = coverTemplate
	}
	channel.UpdatedAt = time.Now()
}

// GetChannelGEOConfig 获取渠道GEO配置
func (s *PublishConfigService) GetChannelGEOConfig(channel *model.Channel) *model.ChannelGEOConfig {
	return s.parseGEOConfig(channel.GEOConfig)
}
