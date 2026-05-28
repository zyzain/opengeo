package service

import (
	"context"
	"testing"
)

func TestCompetitorAnalysisService_AnalyzeCompetitor(t *testing.T) {
	svc := NewCompetitorAnalysisService()
	ctx := context.Background()

	competitor := &CompetitorData{
		Name:            "竞品A",
		Domain:          "competitor-a.com",
		VisibilityScore: 85.0,
		CitationCount:   120,
		AvgPosition:     2.5,
		TopQueries:      []string{"GEO优化", "AI搜索", "内容营销"},
	}

	result, err := svc.AnalyzeCompetitor(ctx, competitor)
	if err != nil {
		t.Fatalf("analyze failed: %v", err)
	}

	if result.CompetitorName != "竞品A" {
		t.Errorf("expected 竞品A, got %s", result.CompetitorName)
	}
	if result.VisibilityLevel != "excellent" {
		t.Errorf("expected excellent, got %s", result.VisibilityLevel)
	}
	if result.CitationPerformance.Trend != "rising" {
		t.Errorf("expected rising, got %s", result.CitationPerformance.Trend)
	}
}

func TestGetVisibilityLevel(t *testing.T) {
	tests := []struct {
		score float32
		want  string
	}{
		{90, "excellent"},
		{80, "excellent"},
		{70, "good"},
		{60, "good"},
		{50, "moderate"},
		{40, "moderate"},
		{30, "low"},
	}

	for _, tt := range tests {
		if got := getVisibilityLevel(tt.score); got != tt.want {
			t.Errorf("getVisibilityLevel(%f) = %s, want %s", tt.score, got, tt.want)
		}
	}
}

func TestAnalyzeCitationPerformance(t *testing.T) {
	tests := []struct {
		data *CompetitorData
		want string
		desc string
	}{
		{
			data: &CompetitorData{CitationCount: 150},
			want: "rising",
			desc: "high citations",
		},
		{
			data: &CompetitorData{CitationCount: 75},
			want: "stable",
			desc: "medium citations",
		},
		{
			data: &CompetitorData{CitationCount: 30},
			want: "declining",
			desc: "low citations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := analyzeCitationPerformance(tt.data)
			if result.Trend != tt.want {
				t.Errorf("expected %s, got %s", tt.want, result.Trend)
			}
		})
	}
}

func TestGenerateSWOT(t *testing.T) {
	// 高可见性竞品
	highVis := &CompetitorData{
		VisibilityScore: 85,
		AvgPosition:     2,
		CitationCount:   100,
	}
	swot := generateSWOT(highVis)
	if len(swot.Strengths) == 0 {
		t.Error("expected strengths for high visibility")
	}
	if len(swot.Threats) == 0 {
		t.Error("expected threats for high citation count")
	}

	// 低可见性竞品
	lowVis := &CompetitorData{
		VisibilityScore: 30,
		AvgPosition:     8,
		CitationCount:   10,
	}
	swot = generateSWOT(lowVis)
	if len(swot.Weaknesses) == 0 {
		t.Error("expected weaknesses for low visibility")
	}
	if len(swot.Opportunities) == 0 {
		t.Error("expected opportunities for high avg position")
	}
}

func TestGenerateComparisonReport(t *testing.T) {
	svc := NewCompetitorAnalysisService()
	ctx := context.Background()

	yourBrand := &CompetitorData{
		Name:            "你的品牌",
		VisibilityScore: 70,
		CitationCount:   80,
		AvgPosition:     3.5,
		TopQueries:      []string{"GEO优化", "AI搜索"},
	}

	competitors := []*CompetitorData{
		{
			Name:            "竞品A",
			VisibilityScore: 85,
			CitationCount:   120,
			AvgPosition:     2.0,
			TopQueries:      []string{"GEO优化", "内容营销", "SEO"},
		},
		{
			Name:            "竞品B",
			VisibilityScore: 60,
			CitationCount:   50,
			AvgPosition:     5.0,
			TopQueries:      []string{"AI搜索", "搜索引擎优化"},
		},
	}

	report, err := svc.GenerateComparisonReport(ctx, yourBrand, competitors)
	if err != nil {
		t.Fatalf("generate report failed: %v", err)
	}

	if report.Rank != 2 {
		t.Errorf("expected rank 2, got %d", report.Rank)
	}
	if report.TotalCompetitors != 3 {
		t.Errorf("expected 3 competitors, got %d", report.TotalCompetitors)
	}
	if len(report.Gaps) == 0 {
		t.Error("expected content gaps")
	}
}

func TestSuggestionEngine_GenerateActionableSuggestions(t *testing.T) {
	engine := NewSuggestionEngine()
	ctx := context.Background()

	// 低优化内容
	sctx := &SuggestionContext{
		ContentID:       1,
		Title:           "测试标题",
		BodyLength:      500,
		HasHeadings:     false,
		HasSchemaMarkup: false,
		HasAIGCLabel:    false,
		CitationCount:   2,
		AvgPosition:     8,
		VisibilityScore: 40,
	}

	suggestions := engine.GenerateActionableSuggestions(ctx, sctx)

	if len(suggestions) < 5 {
		t.Errorf("expected at least 5 suggestions, got %d", len(suggestions))
	}

	// 检查高优先级建议
	hasHighPriority := false
	for _, s := range suggestions {
		if s.Priority == "high" {
			hasHighPriority = true
			break
		}
	}
	if !hasHighPriority {
		t.Error("expected high priority suggestions")
	}
}

func TestSuggestionEngine_HighOptimizationContent(t *testing.T) {
	engine := NewSuggestionEngine()
	ctx := context.Background()

	// 高优化内容
	sctx := &SuggestionContext{
		ContentID:       1,
		Title:           "测试标题",
		BodyLength:      2000,
		HasHeadings:     true,
		HasSchemaMarkup: true,
		HasAIGCLabel:    true,
		CitationCount:   15,
		AvgPosition:     2,
		VisibilityScore: 85,
		Keywords:        []string{"GEO", "AI", "优化"},
	}

	suggestions := engine.GenerateActionableSuggestions(ctx, sctx)

	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions for well-optimized content, got %d", len(suggestions))
	}
}

func TestSuggestionEngine_GetTopSuggestions(t *testing.T) {
	engine := NewSuggestionEngine()
	ctx := context.Background()

	sctx := &SuggestionContext{
		ContentID:       1,
		BodyLength:      500,
		HasHeadings:     false,
		HasSchemaMarkup: false,
		HasAIGCLabel:    false,
		CitationCount:   1,
		AvgPosition:     10,
		VisibilityScore: 30,
	}

	top := engine.GetTopSuggestions(ctx, sctx, 3)

	if len(top) > 3 {
		t.Errorf("expected max 3 suggestions, got %d", len(top))
	}

	// 第一个应该是高优先级
	if len(top) > 0 && top[0].Priority != "high" {
		t.Errorf("expected first suggestion to be high priority, got %s", top[0].Priority)
	}
}

func TestSuggestionEngine_CalculateOptimizationScore(t *testing.T) {
	engine := NewSuggestionEngine()

	sctx := &SuggestionContext{
		ContentID:       1,
		HasHeadings:     true,
		HasSchemaMarkup: true,
		HasAIGCLabel:    true,
		CitationCount:   10,
		AvgPosition:     2,
	}

	score := engine.CalculateOptimizationScore(sctx)

	if score.StructureScore < 80 {
		t.Errorf("expected high structure score, got %f", score.StructureScore)
	}
	if score.AuthorityScore < 80 {
		t.Errorf("expected high authority score, got %f", score.AuthorityScore)
	}
	if score.ComplianceScore < 80 {
		t.Errorf("expected high compliance score, got %f", score.ComplianceScore)
	}
	if score.OverallScore < 80 {
		t.Errorf("expected high overall score, got %f", score.OverallScore)
	}
}

func TestSuggestionEngine_BatchGenerateSuggestions(t *testing.T) {
	engine := NewSuggestionEngine()
	ctx := context.Background()

	contexts := []*SuggestionContext{
		{ContentID: 1, BodyLength: 500, HasHeadings: false},
		{ContentID: 2, BodyLength: 2000, HasHeadings: true, HasSchemaMarkup: true, HasAIGCLabel: true, CitationCount: 10, AvgPosition: 2, VisibilityScore: 80, Keywords: []string{"GEO", "AI", "优化", "搜索引擎", "内容营销"}},
		{ContentID: 3, BodyLength: 300, HasSchemaMarkup: false},
	}

	result := engine.BatchGenerateSuggestions(ctx, contexts)

	if len(result) != 2 {
		t.Errorf("expected 2 contents with suggestions, got %d", len(result))
	}
	if _, exists := result[2]; exists {
		t.Error("content 2 should not have suggestions (well optimized)")
	}
}

func TestFormatSuggestionReport(t *testing.T) {
	suggestions := []*ActionableSuggestion{
		{
			ID:          "test1",
			Category:    "structure",
			Title:       "添加标题",
			Description: "使用标题标记",
			Priority:    "high",
			Impact:      "high",
			Effort:      "low",
			Steps:       []string{"步骤1", "步骤2"},
			Metric:      "提升20%",
		},
	}

	report := FormatSuggestionReport(suggestions)
	if report == "" {
		t.Error("expected non-empty report")
	}
	if !contains(report, "添加标题") {
		t.Error("expected report to contain suggestion title")
	}

	// 空建议
	emptyReport := FormatSuggestionReport([]*ActionableSuggestion{})
	if !contains(emptyReport, "暂无") {
		t.Error("expected 'no suggestions' message")
	}
}

func TestSuggestionToModel(t *testing.T) {
	suggestion := &ActionableSuggestion{
		ID:          "test",
		Category:    "structure",
		Title:       "测试建议",
		Description: "测试描述",
		Priority:    "high",
		Steps:       []string{"步骤1"},
		Metric:      "提升10%",
	}

	model := SuggestionToModel(suggestion, 123)

	if model["content_id"] != int64(123) {
		t.Errorf("expected 123, got %v", model["content_id"])
	}
	if model["priority"] != int32(2) {
		t.Errorf("expected 2, got %v", model["priority"])
	}
}

func TestPriorityToInt(t *testing.T) {
	tests := []struct {
		priority string
		want     int32
	}{
		{"high", 2},
		{"medium", 1},
		{"low", 0},
		{"unknown", 0},
	}

	for _, tt := range tests {
		if got := priorityToInt(tt.priority); got != tt.want {
			t.Errorf("priorityToInt(%s) = %d, want %d", tt.priority, got, tt.want)
		}
	}
}

func TestActionableSuggestion_Fields(t *testing.T) {
	suggestion := &ActionableSuggestion{
		ID:          "test",
		Category:    "structure",
		Title:       "测试标题",
		Description: "测试描述",
		Priority:    "high",
		Impact:      "high",
		Effort:      "low",
		Steps:       []string{"步骤1", "步骤2"},
		Metric:      "提升20%",
	}

	if suggestion.ID != "test" {
		t.Errorf("expected test, got %s", suggestion.ID)
	}
	if len(suggestion.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(suggestion.Steps))
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
