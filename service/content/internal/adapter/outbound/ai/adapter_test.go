package ai

import (
	"testing"

	"opengeo/pkg/ai"
)

func TestParseJSONFromResponse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		desc     string
	}{
		{
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
			desc:     "plain JSON",
		},
		{
			input:    "```json\n{\"key\": \"value\"}\n```",
			expected: `{"key": "value"}`,
			desc:     "markdown code block",
		},
		{
			input:    "```\n{\"key\": \"value\"}\n```",
			expected: `{"key": "value"}`,
			desc:     "code block without language",
		},
		{
			input:    "Here is the result:\n{\"key\": \"value\"}\nDone.",
			expected: `{"key": "value"}`,
			desc:     "JSON embedded in text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := parseJSONFromResponse(tt.input)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestParseOptimizeResponse(t *testing.T) {
	raw := `{
		"optimized_title": "优化后的标题",
		"optimized_body": "优化后的正文",
		"schema_markup": "{\"@context\":\"https://schema.org\"}",
		"score": 85.5,
		"suggestions": ["建议1", "建议2"],
		"structural_changes": ["变更1"]
	}`

	result, err := parseOptimizeResponse(raw)
	if err != nil {
		t.Fatalf("parseOptimizeResponse failed: %v", err)
	}

	if !result.Success {
		t.Error("expected success=true")
	}
	if result.OptimizedTitle != "优化后的标题" {
		t.Errorf("wrong title: %s", result.OptimizedTitle)
	}
	if result.Score != 85.5 {
		t.Errorf("wrong score: %f", result.Score)
	}
	if len(result.Suggestions) != 2 {
		t.Errorf("wrong suggestions count: %d", len(result.Suggestions))
	}
}

func TestParseOptimizeResponse_WithMarkdown(t *testing.T) {
	raw := "Here is the result:\n```json\n" + `{
		"optimized_title": "Test",
		"optimized_body": "Body",
		"schema_markup": "",
		"score": 80,
		"suggestions": [],
		"structural_changes": []
	}` + "\n```"

	result, err := parseOptimizeResponse(raw)
	if err != nil {
		t.Fatalf("parseOptimizeResponse failed: %v", err)
	}
	if result.OptimizedTitle != "Test" {
		t.Errorf("wrong title: %s", result.OptimizedTitle)
	}
}

func TestParseAdaptResponse(t *testing.T) {
	raw := `{
		"adapted_title": "适配标题",
		"adapted_body": "适配正文",
		"format_changes": ["变更1", "变更2"]
	}`

	result, err := parseAdaptResponse(raw)
	if err != nil {
		t.Fatalf("parseAdaptResponse failed: %v", err)
	}

	if !result.Success {
		t.Error("expected success=true")
	}
	if result.AdaptedTitle != "适配标题" {
		t.Errorf("wrong title: %s", result.AdaptedTitle)
	}
	if len(result.FormatChanges) != 2 {
		t.Errorf("wrong changes count: %d", len(result.FormatChanges))
	}
}

func TestParseComplianceResponse(t *testing.T) {
	raw := `{
		"compliant": false,
		"issues": [
			{
				"issue_type": "sensitive",
				"description": "检测到敏感词",
				"severity": "high",
				"suggestion": "建议删除",
				"location": "正文第1段"
			},
			{
				"issue_type": "aigc_label",
				"description": "缺少AIGC标识",
				"severity": "low",
				"suggestion": "添加标识",
				"location": "全文"
			}
		],
		"sensitive_words": ["敏感词"],
		"aigc_label_required": true,
		"report": "合规检测报告",
		"score": 75.0
	}`

	result, err := parseComplianceResponse(raw)
	if err != nil {
		t.Fatalf("parseComplianceResponse failed: %v", err)
	}

	if result.Compliant {
		t.Error("expected compliant=false")
	}
	if len(result.Issues) != 2 {
		t.Errorf("wrong issues count: %d", len(result.Issues))
	}
	if result.Issues[0].IssueType != "sensitive" {
		t.Errorf("wrong issue type: %s", result.Issues[0].IssueType)
	}
	if result.Score != 75.0 {
		t.Errorf("wrong score: %f", result.Score)
	}
	if !result.AIGCLabelRequired {
		t.Error("expected AIGC label required")
	}
}

func TestParseComplianceResponse_Compliant(t *testing.T) {
	raw := `{
		"compliant": true,
		"issues": [],
		"sensitive_words": [],
		"aigc_label_required": false,
		"report": "内容合规",
		"score": 95.0
	}`

	result, err := parseComplianceResponse(raw)
	if err != nil {
		t.Fatalf("parseComplianceResponse failed: %v", err)
	}

	if !result.Compliant {
		t.Error("expected compliant=true")
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(result.Issues))
	}
}

func TestBuildOptimizePrompt(t *testing.T) {
		req := &ai.OptimizeRequest{
			Title:       "测试标题",
			Body:        "测试正文",
			ContentType: "article",
		}

	prompt := buildOptimizePrompt(req)

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}
	if !contains(prompt, "测试标题") {
		t.Error("expected prompt to contain title")
	}
	if !contains(prompt, "测试正文") {
		t.Error("expected prompt to contain body")
	}
	if !contains(prompt, "JSON") {
		t.Error("expected prompt to mention JSON format")
	}
}

func TestBuildAdaptPrompt(t *testing.T) {
	tests := []struct {
		model string
		hint  string
	}{
		{"deepseek", "结构化"},
		{"kimi", "长文本"},
		{"doubao", "简洁"},
		{"chatgpt", "对话式"},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			req := &ai.AdaptRequest{
				Title:       "标题",
				Body:        "正文",
				TargetModel: tt.model,
			}
			prompt := buildAdaptPrompt(req)
			if !contains(prompt, tt.hint) {
				t.Errorf("expected prompt to contain hint '%s' for model %s", tt.hint, tt.model)
			}
		})
	}
}

func TestBuildCompliancePrompt(t *testing.T) {
	req := &ai.ComplianceRequest{
		Title: "标题",
		Body:  "正文",
	}

	prompt := buildCompliancePrompt(req)

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}
	if !contains(prompt, "敏感词") {
		t.Error("expected prompt to mention sensitive words")
	}
	if !contains(prompt, "广告法") {
		t.Error("expected prompt to mention ad law")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s      string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello..."},
		{"你好世界", 2, "你好..."},
	}

	for _, tt := range tests {
		result := truncate(tt.s, tt.maxLen)
		if result != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.s, tt.maxLen, result, tt.want)
		}
	}
}

func TestAIServiceFactory(t *testing.T) {
	factory := NewAIServiceFactory()

	adapter := NewDeepSeekAdapter("test-key")
	factory.RegisterAdapter("deepseek", adapter)

	got, err := factory.GetAdapter("deepseek")
	if err != nil {
		t.Fatalf("GetAdapter failed: %v", err)
	}
	if got == nil {
		t.Error("expected non-nil adapter")
	}

	defaultAdapter, err := factory.GetDefaultAdapter()
	if err != nil {
		t.Fatalf("GetDefaultAdapter failed: %v", err)
	}
	if defaultAdapter == nil {
		t.Error("expected non-nil default adapter")
	}

	_, err = factory.GetAdapter("unknown")
	if err == nil {
		t.Error("expected error for unknown model")
	}
}

func TestLLMConfig_DefaultConfig(t *testing.T) {
	config := LLMConfig{
		BaseURL: "https://api.example.com",
		APIKey:  "test-key",
		Model:   "test-model",
	}
	config.DefaultConfig()

	if config.Timeout == 0 {
		t.Error("expected default timeout")
	}
	if config.MaxTokens == 0 {
		t.Error("expected default max tokens")
	}
	if config.Temperature == 0 {
		t.Error("expected default temperature")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
