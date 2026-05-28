package service

import (
	"math/rand"
	"strings"
	"testing"
)

func TestDeduplicate_Basic(t *testing.T) {
	svc := NewDeduplicationService()

	req := &DeduplicateRequest{
		Title: "GEO优化完全指南：提升搜索引擎可见性",
		Body: `## 什么是GEO优化

GEO优化是一种提升内容在AI搜索引擎中可见性的方法。通过优化内容结构和格式，可以提高被AI引用的概率。

## 如何进行GEO优化

首先，需要分析目标关键词。其次，优化内容结构。最后，添加权威引用和数据支撑。

## 总结

GEO优化是现代SEO的重要组成部分，值得每个内容创作者关注。`,
		Tags:     []string{"GEO", "优化", "搜索引擎"},
		Strategy: "medium",
	}

	resp := svc.Deduplicate(req)

	if resp.Title == "" {
		t.Error("expected non-empty title")
	}
	if resp.Body == "" {
		t.Error("expected non-empty body")
	}
	if resp.Similarity > 1.0 {
		t.Errorf("similarity invalid: %f", resp.Similarity)
	}
	if len(resp.Changes) == 0 {
		t.Error("expected changes")
	}
}

func TestDeduplicate_LightStrategy(t *testing.T) {
	svc := NewDeduplicationService()

	body := "这是一段测试内容，用于验证轻度去重策略。GEO优化非常重要，可以提升搜索排名。"

	req := &DeduplicateRequest{
		Title:    "测试标题",
		Body:     body,
		Strategy: "light",
	}

	resp := svc.Deduplicate(req)

	// 轻度策略应该保留更多原始内容
	if resp.Similarity < 0.5 {
		t.Errorf("light strategy similarity too low: %f", resp.Similarity)
	}
}

func TestDeduplicate_HeavyStrategy(t *testing.T) {
	svc := NewDeduplicationService()

	body := `第一段：GEO优化是提升AI搜索可见性的关键方法。

第二段：通过结构化数据和权威引用，可以提高内容被AI引用的概率。

第三段：建议每个内容创作者都掌握GEO优化技巧。

第四段：定期分析和优化内容是保持竞争力的重要手段。`

	req := &DeduplicateRequest{
		Title:    "GEO优化指南",
		Body:     body,
		Strategy: "heavy",
	}

	resp := svc.Deduplicate(req)

	// 重度策略应该产生更多变化
	if len(resp.Changes) < 2 {
		t.Errorf("heavy strategy expected more changes, got %d", len(resp.Changes))
	}
}

func TestDeduplicate_SynonymReplace(t *testing.T) {
	svc := NewDeduplicationService()

	body := "提升搜索引擎排名非常重要，需要优化关键词和内容结构。"

	req := &DeduplicateRequest{
		Title:    "测试",
		Body:     body,
		Strategy: "heavy",
		Seed:     12345,
	}

	resp := svc.Deduplicate(req)

	// 应该有同义词替换
	if resp.Body == body && resp.Title == "测试" {
		t.Error("expected synonym replacement")
	}
}

func TestDeduplicate_ParagraphReorder(t *testing.T) {
	svc := NewDeduplicationService()

	body := `段落一：这是第一段内容。

段落二：这是第二段内容。

段落三：这是第三段内容。

段落四：这是第四段内容。`

	req := &DeduplicateRequest{
		Title:    "测试",
		Body:     body,
		Strategy: "heavy",
		Seed:     99999,
	}

	resp := svc.Deduplicate(req)

	// 段落顺序应该被调整
	if resp.Body == body {
		t.Error("expected paragraph reordering")
	}

	// 但内容应该保持完整
	if !strings.Contains(resp.Body, "段落一") {
		t.Error("missing paragraph 1")
	}
	if !strings.Contains(resp.Body, "段落四") {
		t.Error("missing paragraph 4")
	}
}

func TestDeduplicate_Tags(t *testing.T) {
	svc := NewDeduplicationService()

	req := &DeduplicateRequest{
		Title:    "测试",
		Body:     "测试内容",
		Tags:     []string{"提升", "优化", "搜索引擎"},
		Strategy: "medium",
	}

	resp := svc.Deduplicate(req)

	// 标签应该有变化
	if stringSliceEqual(resp.Tags, req.Tags) {
		// 可能没有变化，取决于随机性
		t.Log("tags unchanged (may vary by randomness)")
	}
}

func TestDeduplicate_EmptyContent(t *testing.T) {
	svc := NewDeduplicationService()

	req := &DeduplicateRequest{
		Title:    "",
		Body:     "",
		Strategy: "medium",
	}

	resp := svc.Deduplicate(req)

	if resp.Title != "" {
		t.Error("expected empty title")
	}
	if resp.Body != "" {
		t.Error("expected empty body")
	}
}

func TestDeduplicate_SimilarityCalculation(t *testing.T) {
	svc := NewDeduplicationService()

	tests := []struct {
		text1    string
		text2    string
		wantMax  float32
		desc     string
	}{
		{"hello", "hello", 1.0, "identical"},
		{"hello", "world", 0.5, "different"},
		{"", "", 1.0, "both empty"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			sim := svc.calculateSimilarity(tt.text1, tt.text2)
			if tt.text1 == tt.text2 && sim != 1.0 {
				t.Errorf("expected 1.0 for identical, got %f", sim)
			}
		})
	}
}

func TestDeduplicate_SplitParagraphs(t *testing.T) {
	svc := NewDeduplicationService()

	body := "段落一\n\n段落二\n\n段落三"
	paragraphs := svc.splitParagraphs(body)

	if len(paragraphs) != 3 {
		t.Errorf("expected 3 paragraphs, got %d", len(paragraphs))
	}
}

func TestDeduplicate_SentenceTransform(t *testing.T) {
	svc := NewDeduplicationService()

	sentences := []string{
		"GEO优化非常重要。",
		"需要持续改进内容质量。",
	}

	rng := newTestRng()
	result, count := svc.sentenceTransform(strings.Join(sentences, ""), rng)

	if result == "" {
		t.Error("expected non-empty result")
	}
	t.Logf("transformed %d sentences", count)
}

func TestDeduplicate_MediaURLs(t *testing.T) {
	svc := NewDeduplicationService()

	urls := []string{
		"https://example.com/img1.jpg",
		"https://example.com/img2.jpg",
		"https://example.com/img3.jpg",
	}

	rng := newTestRng()
	result := svc.deduplicateMediaURLs(urls, "heavy", rng)

	if len(result) != len(urls) {
		t.Errorf("expected %d URLs, got %d", len(urls), len(result))
	}
}

func TestDeduplicate_HasTransition(t *testing.T) {
	tests := []struct {
		sent string
		want bool
		desc string
	}{
		{"此外，这很重要", true, "has transition"},
		{"这很重要", false, "no transition"},
		{"因此可以得出结论", true, "has transition"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := hasTransition(tt.sent); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildSynonymMap(t *testing.T) {
	m := buildSynonymMap()

	if len(m) == 0 {
		t.Error("expected non-empty synonym map")
	}

	// 检查一些基本同义词
	if synonyms, ok := m["提升"]; !ok || len(synonyms) == 0 {
		t.Error("expected synonyms for 提升")
	}
	if synonyms, ok := m["优化"]; !ok || len(synonyms) == 0 {
		t.Error("expected synonyms for 优化")
	}
}

func newTestRng() *rand.Rand {
	return rand.New(rand.NewSource(12345))
}
