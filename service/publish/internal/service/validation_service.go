package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"opengeo/service/publish/internal/domain/model"
)

var (
	reURL       = regexp.MustCompile(`https?://[^\s\)]+`)
	reImgTag    = regexp.MustCompile(`<img[^>]*>`)
	reAltAttr   = regexp.MustCompile(`alt=["'][^"']*["']`)
	reLink      = regexp.MustCompile(`<a[^>]*href=["']([^"']*)["'][^>]*>`)
	reHeading   = regexp.MustCompile(`(?m)^#{1,6}\s+`)
)

// ValidationService 发布校验服务
type ValidationService struct{}

// NewValidationService 创建校验服务
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// ValidationResult 校验结果
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationItem `json:"errors"`
	Warnings []ValidationItem `json:"warnings"`
	Infos    []ValidationItem `json:"infos"`
	Score    float32           `json:"score"`
}

// ValidationItem 校验项
type ValidationItem struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
}

// ValidateForPublish 发布前完整校验
func (s *ValidationService) ValidateForPublish(
	ctx context.Context,
	channel *model.Channel,
	config *model.PublishConfig,
) *ValidationResult {

	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]ValidationItem, 0),
		Warnings: make([]ValidationItem, 0),
		Infos:    make([]ValidationItem, 0),
		Score:    100,
	}

	geoConfig := s.parseGEOConfig(channel.GEOConfig)

	// 1-4: 标题校验
	s.validateTitle(config.Title, geoConfig, result)

	// 5-9: 正文校验
	s.validateBody(config.Body, geoConfig, result)

	// 10-11: Schema Markup校验
	s.validateSchemaMarkup(config.SchemaMarkup, geoConfig, result)

	// 12-13: 标签校验
	s.validateTags(config.Tags, result)

	// 14: AIGC标识校验
	s.validateAIGCLabel(config.AIGCLabel, geoConfig, result)

	// 15: 封面URL校验
	s.validateCoverURL(config.CoverURL, result)

	// 16: 链接有效性校验
	s.validateLinks(config.Body, result)

	// 17: 图片ALT标签校验
	s.validateImageALT(config.Body, result)

	// 计算最终分数
	s.calculateScore(result)

	return result
}

// ==================== 标题校验 (1-4) ====================

func (s *ValidationService) validateTitle(title string, geoConfig *model.ChannelGEOConfig, result *ValidationResult) {
	// 1. 标题非空
	if title == "" {
		result.Errors = append(result.Errors, ValidationItem{
			Field:   "title",
			Rule:    "required",
			Message: "标题不能为空",
		})
		return
	}

	// 2. 标题长度
	titleLen := utf8.RuneCountInString(title)
	if titleLen < 5 {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "title",
			Rule:    "min_length",
			Message: fmt.Sprintf("标题过短（%d字），建议≥5字", titleLen),
		})
	}
	if geoConfig.MaxTitleLength > 0 && titleLen > geoConfig.MaxTitleLength {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "title",
			Rule:    "max_length",
			Message: fmt.Sprintf("标题超过%d字符（当前%d字），将被截断", geoConfig.MaxTitleLength, titleLen),
		})
	}

	// 3. 标题特殊字符
	specialChars := []string{"<", ">", "{", "}", "|", "\\"}
	for _, ch := range specialChars {
		if strings.Contains(title, ch) {
			result.Warnings = append(result.Warnings, ValidationItem{
				Field:   "title",
				Rule:    "special_chars",
				Message: fmt.Sprintf("标题包含特殊字符「%s」，可能影响显示", ch),
			})
			break
		}
	}

	// 4. 标题格式
	if strings.HasPrefix(title, " ") || strings.HasSuffix(title, " ") {
		result.Infos = append(result.Infos, ValidationItem{
			Field:   "title",
			Rule:    "format",
			Message: "标题首尾包含空格，建议去除",
		})
	}
}

// ==================== 正文校验 (5-9) ====================

func (s *ValidationService) validateBody(body string, geoConfig *model.ChannelGEOConfig, result *ValidationResult) {
	// 5. 正文非空
	if body == "" {
		result.Errors = append(result.Errors, ValidationItem{
			Field:   "body",
			Rule:    "required",
			Message: "正文不能为空",
		})
		return
	}

	// 6. 正文长度
	bodyLen := utf8.RuneCountInString(body)
	if bodyLen < 100 {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "body",
			Rule:    "min_length",
			Message: fmt.Sprintf("正文过短（%d字），建议≥100字", bodyLen),
		})
	}
	if geoConfig.MaxBodyLength > 0 && bodyLen > geoConfig.MaxBodyLength {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "body",
			Rule:    "max_length",
			Message: fmt.Sprintf("正文超过%d字符（当前%d字），将被截断", geoConfig.MaxBodyLength, bodyLen),
		})
	}

	// 7. 正文结构 - 标题层级
	headings := reHeading.FindAllString(body, -1)
	if len(headings) == 0 && bodyLen > 500 {
		result.Infos = append(result.Infos, ValidationItem{
			Field:   "body",
			Rule:    "structure",
			Message: "正文缺少标题层级结构，建议使用 ## / ### 标记",
		})
	}

	// 8. 正文结构 - 段落
	paragraphs := strings.Split(body, "\n\n")
	if len(paragraphs) < 2 && bodyLen > 300 {
		result.Infos = append(result.Infos, ValidationItem{
			Field:   "body",
			Rule:    "paragraphs",
			Message: "正文段落较少，建议分段以提升可读性",
		})
	}

	// 9. 正文空白内容
	trimmed := strings.TrimSpace(body)
	if len(trimmed) == 0 {
		result.Errors = append(result.Errors, ValidationItem{
			Field:   "body",
			Rule:    "empty_content",
			Message: "正文仅包含空白字符",
		})
	}
}

// ==================== Schema Markup校验 (10-11) ====================

func (s *ValidationService) validateSchemaMarkup(schemaMarkup string, geoConfig *model.ChannelGEOConfig, result *ValidationResult) {
	// 10. Schema Markup存在性
	if geoConfig.InjectSchemaMarkup && schemaMarkup == "" {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "schema_markup",
			Rule:    "missing",
			Message: "已启用Schema Markup注入但未生成Schema数据",
		})
		return
	}

	if schemaMarkup == "" {
		return
	}

	// 11. Schema Markup JSON格式
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaMarkup), &schema); err != nil {
		result.Errors = append(result.Errors, ValidationItem{
			Field:   "schema_markup",
			Rule:    "json_format",
			Message: fmt.Sprintf("Schema Markup JSON格式错误: %v", err),
		})
		return
	}

	// 检查必要字段
	if _, ok := schema["@context"]; !ok {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "schema_markup",
			Rule:    "required_field",
			Message: "Schema Markup缺少 @context 字段",
		})
	}
	if _, ok := schema["@type"]; !ok {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "schema_markup",
			Rule:    "required_field",
			Message: "Schema Markup缺少 @type 字段",
		})
	}
}

// ==================== 标签校验 (12-13) ====================

func (s *ValidationService) validateTags(tags []string, result *ValidationResult) {
	// 12. 标签数量
	if len(tags) == 0 {
		result.Infos = append(result.Infos, ValidationItem{
			Field:   "tags",
			Rule:    "empty",
			Message: "未设置标签，建议添加以提升发现性",
		})
	} else if len(tags) > 10 {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "tags",
			Rule:    "too_many",
			Message: fmt.Sprintf("标签过多（%d个），建议≤10个", len(tags)),
		})
	}

	// 13. 标签长度
	for _, tag := range tags {
		if utf8.RuneCountInString(tag) > 20 {
			result.Warnings = append(result.Warnings, ValidationItem{
				Field:   "tags",
				Rule:    "tag_length",
				Message: fmt.Sprintf("标签「%s」过长（%d字），建议≤20字", tag, utf8.RuneCountInString(tag)),
			})
		}
	}
}

// ==================== AIGC标识校验 (14) ====================

func (s *ValidationService) validateAIGCLabel(aigcLabel string, geoConfig *model.ChannelGEOConfig, result *ValidationResult) {
	if geoConfig.InjectAIGCLabel && aigcLabel == "" {
		result.Warnings = append(result.Warnings, ValidationItem{
			Field:   "aigc_label",
			Rule:    "missing",
			Message: "已启用AIGC标识但未生成标识内容",
		})
	}
}

// ==================== 封面URL校验 (15) ====================

func (s *ValidationService) validateCoverURL(coverURL string, result *ValidationResult) {
	if coverURL == "" {
		return
	}

	// 15. URL格式
	u, err := url.Parse(coverURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		result.Errors = append(result.Errors, ValidationItem{
			Field:   "cover_url",
			Rule:    "url_format",
			Message: fmt.Sprintf("封面URL格式错误: %s", coverURL),
		})
	}
}

// ==================== 链接校验 (16) ====================

func (s *ValidationService) validateLinks(body string, result *ValidationResult) {
	links := reLink.FindAllStringSubmatch(body, -1)
	for _, link := range links {
		if len(link) < 2 {
			continue
		}
		href := link[1]

		// 16. 链接格式
		if !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "/") {
			result.Warnings = append(result.Warnings, ValidationItem{
				Field:   "body",
				Rule:    "link_format",
				Message: fmt.Sprintf("链接格式可能无效: %s", href),
			})
		}
	}
}

// ==================== 图片ALT校验 (17) ====================

func (s *ValidationService) validateImageALT(body string, result *ValidationResult) {
	imgTags := reImgTag.FindAllString(body, -1)
	for _, img := range imgTags {
		// 17. 图片ALT属性
		if !reAltAttr.MatchString(img) {
			result.Warnings = append(result.Warnings, ValidationItem{
				Field:   "body",
				Rule:    "img_alt",
				Message: "图片缺少ALT属性，影响可访问性和SEO",
			})
		}
	}
}

// ==================== 分数计算 ====================

func (s *ValidationService) calculateScore(result *ValidationResult) {
	score := float32(100)

	// 错误扣分
	for range result.Errors {
		score -= 15
	}

	// 警告扣分
	for range result.Warnings {
		score -= 5
	}

	// 信息不扣分

	if score < 0 {
		score = 0
	}

	result.Score = score
	result.Valid = len(result.Errors) == 0
}

// ==================== 辅助方法 ====================

func (s *ValidationService) parseGEOConfig(configJSON string) *model.ChannelGEOConfig {
	config := &model.ChannelGEOConfig{
		MaxTitleLength: 64,
		MaxBodyLength:  5000,
	}

	if configJSON == "" {
		return config
	}

	json.Unmarshal([]byte(configJSON), config)
	return config
}

// FormatValidationResult 格式化校验结果
func (s *ValidationService) FormatValidationResult(result *ValidationResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("校验结果: %s (分数: %.0f/100)\n", boolToStatus(result.Valid), result.Score))

	if len(result.Errors) > 0 {
		sb.WriteString(fmt.Sprintf("\n❌ 错误 (%d项):\n", len(result.Errors)))
		for i, err := range result.Errors {
			sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, err.Field, err.Message))
		}
	}

	if len(result.Warnings) > 0 {
		sb.WriteString(fmt.Sprintf("\n⚠️ 警告 (%d项):\n", len(result.Warnings)))
		for i, warn := range result.Warnings {
			sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, warn.Field, warn.Message))
		}
	}

	if len(result.Infos) > 0 {
		sb.WriteString(fmt.Sprintf("\n💡 建议 (%d项):\n", len(result.Infos)))
		for i, info := range result.Infos {
			sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, info.Field, info.Message))
		}
	}

	return sb.String()
}

func boolToStatus(b bool) string {
	if b {
		return "✅ 通过"
	}
	return "❌ 未通过"
}
