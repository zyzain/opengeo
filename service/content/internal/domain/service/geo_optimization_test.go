package service

import (
	"context"
	"strings"
	"testing"

	"opengeo/service/content/internal/domain/model"
)

func newTestContent(title, body, contentType, schemaMarkup string) *model.Content {
	return &model.Content{
		ID:           1,
		UserID:       1,
		Title:        title,
		Body:         body,
		ContentType:  contentType,
		Status:       0,
		SchemaMarkup: schemaMarkup,
	}
}

func TestOptimizeForAI_BasicStructure(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	body := `## 什么是GEO优化？

GEO（Generative Engine Optimization）是一种针对AI搜索引擎的内容优化方法。

### 核心要点

- 结构化数据：使用Schema Markup标注内容
- 权威引用：链接到权威信源
- 事实密度：提供具体数字和数据

## 如何实施GEO优化？

根据2024年研究数据显示，采用GEO优化的内容被AI引用率提升约35%。

1. 分析目标AI模型的偏好
2. 优化内容结构和格式
3. 添加权威引用和数据支撑

> 以上方法已被验证有效。

更多信息请参考 https://schema.org 和 https://arxiv.org/abs/2024.example。

## 总结

综上所述，GEO优化是提升品牌在AI搜索中可见性的关键策略。`

	content := newTestContent("GEO优化完全指南：提升AI搜索可见性", body, "article", "")
	result, err := svc.OptimizeForAI(ctx, content)
	if err != nil {
		t.Fatalf("OptimizeForAI failed: %v", err)
	}

	if result.TotalScore <= 0 {
		t.Error("expected positive total score")
	}
	if result.StructureScore <= 0 {
		t.Error("expected positive structure score")
	}
	if result.ReadabilityScore <= 0 {
		t.Error("expected positive readability score")
	}
	if len(result.Suggestions) == 0 {
		t.Error("expected at least one suggestion")
	}
	if result.SchemaMarkup == "" {
		t.Error("expected generated schema markup")
	}

	t.Logf("Total Score: %.1f", result.TotalScore)
	t.Logf("Structure Score: %.1f", result.StructureScore)
	t.Logf("Readability Score: %.1f", result.ReadabilityScore)
	t.Logf("Schema Markup generated: %d bytes", len(result.SchemaMarkup))
	for i, s := range result.Suggestions {
		t.Logf("Suggestion %d: %s", i+1, s)
	}
}

func TestOptimizeForAI_WithExistingSchema(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	existingSchema := `{
		"@context": "https://schema.org",
		"@type": "Article",
		"headline": "Test Article",
		"description": "A test article",
		"author": {"@type": "Organization", "name": "Test"},
		"datePublished": "2024-01-01",
		"articleBody": "Test body"
	}`

	body := `## 标题

这是一篇测试文章，包含一些内容。2024年数据显示效果提升了50%。

参考 https://example.com 了解更多。

## 总结

总结一下要点。`

	content := newTestContent("测试文章", body, "article", existingSchema)
	result, err := svc.OptimizeForAI(ctx, content)
	if err != nil {
		t.Fatalf("OptimizeForAI failed: %v", err)
	}

	if result.SchemaMarkup != existingSchema {
		t.Error("expected existing schema markup to be preserved")
	}
}

func TestOptimizeForAI_InvalidContent(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	content := newTestContent("", "", "article", "")
	_, err := svc.OptimizeForAI(ctx, content)
	if err == nil {
		t.Error("expected error for invalid content")
	}
}

func TestAdaptForModel(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	body := `## 测试内容

这是一段测试内容，用于验证模型适配功能。`

	content := newTestContent("测试", body, "article", "")

	tests := []struct {
		model    string
		wantFunc func(string) bool
		desc     string
	}{
		{
			model: "deepseek",
			wantFunc: func(s string) bool {
				return strings.Contains(s, "总结")
			},
			desc: "DeepSeek should add conclusion",
		},
		{
			model: "chatgpt",
			wantFunc: func(s string) bool {
				return strings.Contains(s, "FAQ") || strings.Contains(s, "常见问题")
			},
			desc: "ChatGPT should add FAQ section",
		},
		{
			model: "kimi",
			wantFunc: func(s string) bool {
				return len(s) >= len(body)
			},
			desc: "Kimi should keep or expand content",
		},
		{
			model: "doubao",
			wantFunc: func(s string) bool {
				return strings.Contains(s, "要点提炼")
			},
			desc: "Doubao should add summary header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			result, err := svc.AdaptForModel(ctx, content, tt.model)
			if err != nil {
				t.Fatalf("AdaptForModel failed: %v", err)
			}
			if !tt.wantFunc(result.AdaptedContent) {
				t.Errorf("adaptation for %s failed: %s", tt.model, tt.desc)
			}
			if len(result.FormatChanges) == 0 {
				t.Errorf("expected format changes for model %s", tt.model)
			}
		})
	}
}

func TestCheckCompliance_SensitiveWords(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	body := `## 测试

这段内容包含赌博相关信息，需要被检测出来。`

	content := newTestContent("测试", body, "article", "")
	result, err := svc.CheckCompliance(ctx, content)
	if err != nil {
		t.Fatalf("CheckCompliance failed: %v", err)
	}

	if result.Compliant {
		t.Error("expected non-compliant content with sensitive word")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.IssueType == "sensitive" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected sensitive word issue")
	}
}

func TestCheckCompliance_AdLawWords(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	body := `## 测试

我们的产品是行业第一，效果最好，绝对领先。`

	content := newTestContent("测试", body, "article", "")
	result, err := svc.CheckCompliance(ctx, content)
	if err != nil {
		t.Fatalf("CheckCompliance failed: %v", err)
	}

	if result.Compliant {
		t.Error("expected non-compliant content with ad law violation")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.IssueType == "ad_law" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected ad law issue")
	}
}

func TestCheckCompliance_AIGCLabel(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	body := `## 测试

这是一段没有标识的普通内容，需要检测是否缺少生成标识。`

	content := newTestContent("测试", body, "article", "")
	result, err := svc.CheckCompliance(ctx, content)
	if err != nil {
		t.Fatalf("CheckCompliance failed: %v", err)
	}

	if result.Compliant {
		t.Error("expected non-compliant content without AIGC label")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.IssueType == "aigc_label" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected AIGC label issue")
	}
}

func TestCheckCompliance_CleanContent(t *testing.T) {
	svc := NewGEOOptimizationService()
	ctx := context.Background()

	body := `## 测试

这是一段合规的内容。本文由AI辅助生成（AIGC）。

参考 https://example.com`

	content := newTestContent("测试", body, "article", "")
	result, err := svc.CheckCompliance(ctx, content)
	if err != nil {
		t.Fatalf("CheckCompliance failed: %v", err)
	}

	if !result.Compliant {
		t.Errorf("expected compliant, but got issues: %v", result.Issues)
	}
	if result.Report == "" {
		t.Error("expected non-empty report")
	}
}

func TestAnalyzeFactualDensity(t *testing.T) {
	svc := NewGEOOptimizationService()

	tests := []struct {
		body     string
		expected bool
		desc     string
	}{
		{
			body:     "2024年数据显示，采用率提升了35%，覆盖了1000万用户。",
			expected: true,
			desc:     "high factual density",
		},
		{
			body:     "这是一段没有任何数据支撑的普通文字内容描述。",
			expected: false,
			desc:     "low factual density",
		},
	}

	for _, tt := range tests {
		score := svc.analyzeFactualDensity(tt.body)
		t.Logf("%s: score=%.1f", tt.desc, score)
		if tt.expected && score < 30 {
			t.Errorf("%s: expected score >= 30, got %.1f", tt.desc, score)
		}
		if !tt.expected && score > 30 {
			t.Errorf("%s: expected score <= 30, got %.1f", tt.desc, score)
		}
	}
}

func TestAnalyzeAuthorityReferences(t *testing.T) {
	svc := NewGEOOptimizationService()

	tests := []struct {
		body      string
		wantCount int
		desc      string
	}{
		{
			body:      "参考 https://arxiv.org/abs/2024.example 和 https://en.wikipedia.org/wiki/GEO",
			wantCount: 2,
			desc:      "two authority references",
		},
		{
			body:      "参考 https://example.com/page",
			wantCount: 0,
			desc:      "no authority references",
		},
	}

	for _, tt := range tests {
		_, count := svc.analyzeAuthorityReferences(tt.body)
		if count != tt.wantCount {
			t.Errorf("%s: expected %d authority refs, got %d", tt.desc, tt.wantCount, count)
		}
	}
}

func TestGenerateSchemaMarkup(t *testing.T) {
	svc := NewGEOOptimizationService()

	content := newTestContent("测试标题", "测试正文内容", "article", "")

	t.Run("default author", func(t *testing.T) {
		schema := svc.generateSchemaMarkup(content, "")
		if schema == "" {
			t.Error("expected non-empty schema markup")
		}
		if !strings.Contains(schema, "schema.org") {
			t.Error("expected schema.org in markup")
		}
		if !strings.Contains(schema, "Article") {
			t.Error("expected Article type in markup")
		}
		if !strings.Contains(schema, "测试标题") {
			t.Error("expected title in markup")
		}
		if !strings.Contains(schema, "OpenGEO") {
			t.Error("expected default author OpenGEO in markup")
		}
	})

	t.Run("custom author", func(t *testing.T) {
		schema := svc.generateSchemaMarkup(content, "CustomAuthor")
		if schema == "" {
			t.Error("expected non-empty schema markup")
		}
		if !strings.Contains(schema, "CustomAuthor") {
			t.Error("expected custom author in markup")
		}
		if strings.Contains(schema, `"name":  "OpenGEO"`) {
			// publisher should still be OpenGEO, but author should be custom
		}
	})
}

func TestSplitParagraphs(t *testing.T) {
	svc := NewGEOOptimizationService()

	body := "段落一\n\n段落二\n\n段落三"
	paragraphs := svc.splitParagraphs(body)

	if len(paragraphs) != 3 {
		t.Errorf("expected 3 paragraphs, got %d", len(paragraphs))
	}
}

func TestKeywordDensity(t *testing.T) {
	svc := NewGEOOptimizationService()

	body := "GEO优化是一种新的方法。GEO优化帮助提升AI搜索可见性。GEO优化非常重要。"
	score := svc.analyzeKeywordDensity("GEO优化", body)

	if score < 50 {
		t.Errorf("expected high keyword density score, got %.1f", score)
	}
}
