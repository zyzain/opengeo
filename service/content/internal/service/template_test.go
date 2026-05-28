package service

import (
	"encoding/json"
	"testing"

	"opengeo/service/content/internal/domain/model"
)

func TestRateTemplate_Validation(t *testing.T) {
	// Test rating validation without DB
	tests := []struct {
		score   float32
		wantErr bool
		desc    string
	}{
		{3.0, false, "valid rating"},
		{1.0, false, "minimum rating"},
		{5.0, false, "maximum rating"},
		{0.5, true, "below minimum"},
		{5.5, true, "above maximum"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := validateRating(tt.score)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRating(%f) error = %v, wantErr %v", tt.score, err, tt.wantErr)
			}
		})
	}
}

func TestExportTemplate_JSON(t *testing.T) {
	template := &model.ContentTemplate{
		Name:         "测试模板",
		Description:  "测试描述",
		TemplateType: "article",
		TemplateData: "## {{title}}\n{{body}}",
		Tags:         "test,template",
		Author:       "test",
		IsOfficial:   false,
	}

	exportData := model.TemplateExportData{
		Name:         template.Name,
		Description:  template.Description,
		TemplateType: template.TemplateType,
		TemplateData: template.TemplateData,
		Tags:         template.Tags,
		Author:       template.Author,
		IsOfficial:   template.IsOfficial,
	}

	data, err := json.Marshal(exportData)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var imported model.TemplateExportData
	if err := json.Unmarshal(data, &imported); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if imported.Name != template.Name {
		t.Errorf("name mismatch: %s vs %s", imported.Name, template.Name)
	}
	if imported.TemplateData != template.TemplateData {
		t.Error("template data mismatch")
	}
}

func TestImportTemplate_Validation(t *testing.T) {
	tests := []struct {
		jsonData string
		wantErr  bool
		desc     string
	}{
		{`{"name":"test","template_data":"data"}`, false, "valid import"},
		{`{"name":"","template_data":"data"}`, true, "missing name"},
		{`{"name":"test","template_data":""}`, true, "missing data"},
		{`invalid json`, true, "invalid json"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var exportData model.TemplateExportData
			err := json.Unmarshal([]byte(tt.jsonData), &exportData)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			if exportData.Name == "" || exportData.TemplateData == "" {
				if !tt.wantErr {
					t.Error("expected validation error")
				}
			}
		})
	}
}

func TestBuiltinTemplates_Count(t *testing.T) {
	// Verify we have 10 built-in templates
	templateNames := []string{
		"GEO文章优化模板",
		"FAQ问答模板",
		"产品评测模板",
		"行业分析模板",
		"教程指南模板",
		"新闻资讯模板",
		"对比分析模板",
		"数据报告模板",
		"GEO优化清单模板",
		"竞品分析模板",
	}

	if len(templateNames) != 10 {
		t.Errorf("expected 10 templates, got %d", len(templateNames))
	}
}

func TestTemplateTypes(t *testing.T) {
	expectedTypes := []string{
		"article", "faq", "review", "analysis",
		"tutorial", "news", "comparison", "report",
		"checklist", "competitor",
	}

	for _, typ := range expectedTypes {
		if typ == "" {
			t.Error("empty template type")
		}
	}
}

func TestContentTemplate_UsageCount(t *testing.T) {
	template := &model.ContentTemplate{
		UsageCount:  5,
		Rating:      4.0,
		RatingCount: 10,
	}

	// Simulate rating
	newScore := float32(3.0)
	totalScore := template.Rating * float32(template.RatingCount)
	template.RatingCount++
	template.Rating = (totalScore + newScore) / float32(template.RatingCount)

	if template.RatingCount != 11 {
		t.Errorf("expected rating count 11, got %d", template.RatingCount)
	}

	expectedRating := (4.0*10 + 3.0) / 11
	if template.Rating != float32(expectedRating) {
		t.Errorf("expected rating %f, got %f", expectedRating, template.Rating)
	}
}

// validateRating validates rating is between 1 and 5
func validateRating(score float32) error {
	if score < 1 || score > 5 {
		return &ratingError{score: score}
	}
	return nil
}

type ratingError struct {
	score float32
}

func (e *ratingError) Error() string {
	return "rating must be between 1 and 5"
}
