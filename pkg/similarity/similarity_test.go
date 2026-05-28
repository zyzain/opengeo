package similarity

import (
	"testing"
)

func TestSimHash_Compute(t *testing.T) {
	sh := NewSimHash()

	hash1 := sh.Compute("GEO优化是提升AI搜索可见性的关键方法")
	hash2 := sh.Compute("GEO优化是提升AI搜索可见性的重要手段")
	hash3 := sh.Compute("今天天气真好")

	if hash1 == 0 {
		t.Error("expected non-zero hash")
	}

	// 相似文本应该有较小的汉明距离
	dist12 := sh.HammingDistance(hash1, hash2)
	dist13 := sh.HammingDistance(hash1, hash3)

	if dist12 >= dist13 {
		t.Errorf("similar texts should have smaller distance: dist12=%d, dist13=%d", dist12, dist13)
	}
}

func TestSimHash_Similarity(t *testing.T) {
	sh := NewSimHash()

	text1 := "GEO优化完全指南"
	text2 := "GEO优化详细指南"
	text3 := "完全不同的内容"

	hash1 := sh.Compute(text1)
	hash2 := sh.Compute(text2)
	hash3 := sh.Compute(text3)

	sim12 := sh.Similarity(hash1, hash2)
	sim13 := sh.Similarity(hash1, hash3)

	if sim12 <= sim13 {
		t.Errorf("similar texts should have higher similarity: sim12=%.4f, sim13=%.4f", sim12, sim13)
	}
	if sim12 < 0 || sim12 > 1 {
		t.Errorf("similarity out of range: %f", sim12)
	}
}

func TestSimHash_EmptyText(t *testing.T) {
	sh := NewSimHash()
	hash := sh.Compute("")
	if hash != 0 {
		t.Errorf("expected 0 for empty text, got %d", hash)
	}
}

func TestCosineSimilarity(t *testing.T) {
	cs := NewCosineSimilarity()

	tests := []struct {
		text1   string
		text2   string
		minSim  float64
		maxSim  float64
		desc    string
	}{
		{"hello world", "hello world", 0.99, 1.01, "identical"},
		{"hello world", "hello there", 0.0, 0.8, "partial match"},
		{"abc", "xyz", 0.0, 0.3, "different"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			sim := cs.Compute(tt.text1, tt.text2)
			if sim < tt.minSim || sim > tt.maxSim {
				t.Errorf("cosine similarity %.4f not in range [%.2f, %.2f]", sim, tt.minSim, tt.maxSim)
			}
		})
	}
}

func TestJaccardSimilarity(t *testing.T) {
	js := NewJaccardSimilarity()

	sim := js.Compute("hello world", "hello world")
	if sim != 1.0 {
		t.Errorf("expected 1.0 for identical, got %f", sim)
	}

	sim = js.Compute("hello", "world")
	if sim < 0 || sim > 1 {
		t.Errorf("similarity out of range: %f", sim)
	}
}

func TestCombinedSimilarity(t *testing.T) {
	cs := NewCombinedSimilarity()

	text1 := "GEO优化是提升搜索引擎可见性的方法"
	text2 := "GEO优化是提升AI搜索可见性的手段"
	text3 := "完全不相关的文本内容"

	result := cs.Compute(text1, text2, 0.7)
	if result.CombinedSimilarity <= 0 {
		t.Error("expected positive combined similarity")
	}

	result3 := cs.Compute(text1, text3, 0.7)
	if result3.CombinedSimilarity >= result.CombinedSimilarity {
		t.Error("different texts should have lower similarity")
	}
}

func TestCombinedSimilarity_Fingerprint(t *testing.T) {
	cs := NewCombinedSimilarity()

	fp := cs.ComputeFingerprint("这是一段测试文本，用于验证内容指纹功能")
	if fp.SimHash == 0 {
		t.Error("expected non-zero simhash")
	}
	if fp.WordCount == 0 {
		t.Error("expected non-zero word count")
	}
	if len(fp.Keywords) == 0 {
		t.Error("expected keywords")
	}
}

func TestTextExtractor_ExtractKeywords(t *testing.T) {
	e := NewTextExtractor()

	text := "GEO优化非常重要，优化可以提升搜索排名，优化是现代SEO的核心"
	keywords := e.ExtractKeywords(text, 5)

	if len(keywords) == 0 {
		t.Error("expected keywords")
	}
	if len(keywords) > 5 {
		t.Errorf("expected at most 5 keywords, got %d", len(keywords))
	}

	// "优化" 应该是高频词
	found := false
	for _, kw := range keywords {
		if kw.Word == "优化" {
			found = true
			if kw.Count < 2 {
				t.Errorf("expected '优化' count >= 2, got %d", kw.Count)
			}
			break
		}
	}
	if !found {
		t.Error("expected '优化' in keywords")
	}
}

func TestTextExtractor_ExtractNgrams(t *testing.T) {
	e := NewTextExtractor()

	ngrams := e.ExtractNgrams("hello world test", 2)
	if len(ngrams) == 0 {
		t.Error("expected ngrams")
	}
}

func BenchmarkSimHash_Compute(b *testing.B) {
	sh := NewSimHash()
	text := "GEO优化是提升AI搜索引擎可见性的关键方法，通过优化内容结构和格式可以提高被引用的概率"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sh.Compute(text)
	}
}

func BenchmarkCosineSimilarity_Compute(b *testing.B) {
	cs := NewCosineSimilarity()
	text1 := "GEO优化是提升搜索引擎可见性的方法"
	text2 := "GEO优化是提升AI搜索可见性的手段"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs.Compute(text1, text2)
	}
}

func BenchmarkCombinedSimilarity_Compute(b *testing.B) {
	cs := NewCombinedSimilarity()
	text1 := "GEO优化是提升搜索引擎可见性的方法，通过结构化数据和权威引用来实现"
	text2 := "GEO优化是提升AI搜索可见性的手段，需要添加结构化数据和引用来源"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs.Compute(text1, text2, 0.7)
	}
}

func BenchmarkTokenize(b *testing.B) {
	text := "这是一段用于测试分词性能的中文文本，包含各种标点符号和空格。"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenize(text)
	}
}
