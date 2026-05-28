package similarity

import (
	"hash/fnv"
	"math"
	"sort"
	"strings"
	"unicode/utf8"
)

// SimHash SimHash指纹计算
type SimHash struct {
	hashBits int
}

// NewSimHash 创建SimHash计算器
func NewSimHash() *SimHash {
	return &SimHash{
		hashBits: 64,
	}
}

// Compute 计算文本的SimHash值
func (s *SimHash) Compute(text string) uint64 {
	// 分词
	tokens := tokenize(text)
	if len(tokens) == 0 {
		return 0
	}

	// 初始化向量
	vectors := make([]int, s.hashBits)

	// 对每个token计算hash并累加
	for _, token := range tokens {
		hash := s.hashToken(token)
		for i := 0; i < s.hashBits; i++ {
			if (hash >> uint(i)) & 1 == 1 {
				vectors[i]++
			} else {
				vectors[i]--
			}
		}
	}

	// 生成最终hash
	var result uint64
	for i := 0; i < s.hashBits; i++ {
		if vectors[i] > 0 {
			result |= 1 << uint(i)
		}
	}

	return result
}

// HammingDistance 计算两个SimHash的汉明距离
func (s *SimHash) HammingDistance(hash1, hash2 uint64) int {
	xor := hash1 ^ hash2
	count := 0
	for xor != 0 {
		count++
		xor &= xor - 1
	}
	return count
}

// Similarity 计算基于SimHash的相似度 (0-1)
func (s *SimHash) Similarity(hash1, hash2 uint64) float64 {
	distance := s.HammingDistance(hash1, hash2)
	return 1.0 - float64(distance)/float64(s.hashBits)
}

func (s *SimHash) hashToken(token string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(token))
	return h.Sum64()
}

// CosineSimilarity 余弦相似度计算器
type CosineSimilarity struct{}

// NewCosineSimilarity 创建余弦相似度计算器
func NewCosineSimilarity() *CosineSimilarity {
	return &CosineSimilarity{}
}

// Compute 计算两个文本的余弦相似度
func (c *CosineSimilarity) Compute(text1, text2 string) float64 {
	// 分词
	tokens1 := tokenize(text1)
	tokens2 := tokenize(text2)

	// 构建词频向量
	vec1 := buildVector(tokens1)
	vec2 := buildVector(tokens2)

	// 计算余弦相似度
	return cosineSimilarity(vec1, vec2)
}

// ComputeWithKeywords 计算带关键词权重的余弦相似度
func (c *CosineSimilarity) ComputeWithKeywords(text1, text2 string, keywords map[string]float64) float64 {
	tokens1 := tokenize(text1)
	tokens2 := tokenize(text2)

	vec1 := buildWeightedVector(tokens1, keywords)
	vec2 := buildWeightedVector(tokens2, keywords)

	return cosineSimilarity(vec1, vec2)
}

// JaccardSimilarity Jaccard相似度
type JaccardSimilarity struct{}

// NewJaccardSimilarity 创建Jaccard相似度计算器
func NewJaccardSimilarity() *JaccardSimilarity {
	return &JaccardSimilarity{}
}

// Compute 计算两个文本的Jaccard相似度
func (j *JaccardSimilarity) Compute(text1, text2 string) float64 {
	tokens1 := tokenize(text1)
	tokens2 := tokenize(text2)

	set1 := make(map[string]bool)
	for _, t := range tokens1 {
		set1[t] = true
	}

	set2 := make(map[string]bool)
	for _, t := range tokens2 {
		set2[t] = true
	}

	// 计算交集
	intersection := 0
	for t := range set1 {
		if set2[t] {
			intersection++
		}
	}

	// 计算并集
	union := len(set1) + len(set2) - intersection

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// TextExtractor 文本特征提取器
type TextExtractor struct{}

// NewTextExtractor 创建文本特征提取器
func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

// ExtractKeywords 提取关键词
func (e *TextExtractor) ExtractKeywords(text string, topK int) []Keyword {
	tokens := tokenize(text)
	
	// 计算词频
	freq := make(map[string]int)
	for _, t := range tokens {
		freq[t]++
	}

	// 转换为排序列表
	keywords := make([]Keyword, 0, len(freq))
	for word, count := range freq {
		keywords = append(keywords, Keyword{
			Word:  word,
			Count: count,
			TF:    float64(count) / float64(len(tokens)),
		})
	}

	// 按频率排序
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Count > keywords[j].Count
	})

	// 返回topK
	if topK > 0 && topK < len(keywords) {
		keywords = keywords[:topK]
	}

	return keywords
}

// ExtractNgrams 提取N-gram
func (e *TextExtractor) ExtractNgrams(text string, n int) []string {
	tokens := tokenize(text)
	if len(tokens) < n {
		return tokens
	}

	ngrams := make([]string, 0, len(tokens)-n+1)
	for i := 0; i <= len(tokens)-n; i++ {
		ngrams = append(ngrams, strings.Join(tokens[i:i+n], ""))
	}

	return ngrams
}

// Keyword 关键词
type Keyword struct {
	Word  string
	Count int
	TF    float64
}

// CombinedSimilarity 组合相似度计算器
type CombinedSimilarity struct {
	simHash     *SimHash
	cosine      *CosineSimilarity
	jaccard     *JaccardSimilarity
	extractor   *TextExtractor
}

// NewCombinedSimilarity 创建组合相似度计算器
func NewCombinedSimilarity() *CombinedSimilarity {
	return &CombinedSimilarity{
		simHash:   NewSimHash(),
		cosine:    NewCosineSimilarity(),
		jaccard:   NewJaccardSimilarity(),
		extractor: NewTextExtractor(),
	}
}

// SimilarityResult 相似度结果
type SimilarityResult struct {
	SimHashSimilarity    float64 `json:"simhash_similarity"`
	CosineSimilarity     float64 `json:"cosine_similarity"`
	JaccardSimilarity    float64 `json:"jaccard_similarity"`
	CombinedSimilarity   float64 `json:"combined_similarity"`
	IsDuplicate          bool    `json:"is_duplicate"`
	DuplicateThreshold   float64 `json:"duplicate_threshold"`
}

// Compute 计算组合相似度
func (c *CombinedSimilarity) Compute(text1, text2 string, threshold float64) *SimilarityResult {
	// 计算各种相似度
	simHash1 := c.simHash.Compute(text1)
	simHash2 := c.simHash.Compute(text2)
	simHashSim := c.simHash.Similarity(simHash1, simHash2)

	cosineSim := c.cosine.Compute(text1, text2)
	jaccardSim := c.jaccard.Compute(text1, text2)

	// 组合相似度 (加权平均)
	combinedSim := simHashSim*0.3 + cosineSim*0.5 + jaccardSim*0.2

	return &SimilarityResult{
		SimHashSimilarity:  simHashSim,
		CosineSimilarity:   cosineSim,
		JaccardSimilarity:  jaccardSim,
		CombinedSimilarity: combinedSim,
		IsDuplicate:        combinedSim >= threshold,
		DuplicateThreshold: threshold,
	}
}

// ComputeFingerprint 计算内容指纹
func (c *CombinedSimilarity) ComputeFingerprint(text string) *ContentFingerprintData {
	simHash := c.simHash.Compute(text)
	keywords := c.extractor.ExtractKeywords(text, 20)
	
	return &ContentFingerprintData{
		SimHash:   simHash,
		Keywords:  keywords,
		WordCount: utf8.RuneCountInString(text),
	}
}

// ContentFingerprintData 内容指纹数据
type ContentFingerprintData struct {
	SimHash   uint64    `json:"simhash"`
	Keywords  []Keyword `json:"keywords"`
	WordCount int       `json:"word_count"`
}

// 辅助函数

// tokenize 中文分词（简单实现，按字符和标点分割）
func tokenize(text string) []string {
	// 转小写
	text = strings.ToLower(text)

	// 按标点和空格分割
	var tokens []string
	var current strings.Builder

	for _, r := range text {
		if isDelimiter(r) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	// 对中文进行字符级分词（2-gram）
	var result []string
	for _, token := range tokens {
		runes := []rune(token)
		if len(runes) <= 2 {
			result = append(result, token)
		} else {
			// 2-gram
			for i := 0; i < len(runes)-1; i++ {
				result = append(result, string(runes[i:i+2]))
			}
		}
	}

	return result
}

func isDelimiter(r rune) bool {
	return r == ' ' || r == '，' || r == '。' || r == '！' || r == '？' || 
		r == '；' || r == '：' || r == '\u201c' || r == '\u201d' || r == '\u2018' || r == '\u2019' ||
		r == '（' || r == '）' || r == '【' || r == '】' || r == '《' || r == '》' ||
		r == ',' || r == '.' || r == '!' || r == '?' || r == ';' || r == ':' ||
		r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' ||
		r == '\n' || r == '\t' || r == '\r'
}

func buildVector(tokens []string) map[string]float64 {
	vec := make(map[string]float64)
	for _, t := range tokens {
		vec[t]++
	}
	return vec
}

func buildWeightedVector(tokens []string, keywords map[string]float64) map[string]float64 {
	vec := make(map[string]float64)
	for _, t := range tokens {
		weight := 1.0
		if w, ok := keywords[t]; ok {
			weight = w
		}
		vec[t] += weight
	}
	return vec
}

func cosineSimilarity(vec1, vec2 map[string]float64) float64 {
	// 计算点积
	dotProduct := 0.0
	for k, v1 := range vec1 {
		if v2, ok := vec2[k]; ok {
			dotProduct += v1 * v2
		}
	}

	// 计算向量长度
	norm1 := 0.0
	for _, v := range vec1 {
		norm1 += v * v
	}
	norm1 = math.Sqrt(norm1)

	norm2 := 0.0
	for _, v := range vec2 {
		norm2 += v * v
	}
	norm2 = math.Sqrt(norm2)

	// 避免除零
	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (norm1 * norm2)
}
