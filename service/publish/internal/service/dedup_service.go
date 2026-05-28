package service

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

// DeduplicationService 内容去重服务
type DeduplicationService struct {
	synonymMap map[string][]string
}

// NewDeduplicationService 创建去重服务
func NewDeduplicationService() *DeduplicationService {
	svc := &DeduplicationService{
		synonymMap: buildSynonymMap(),
	}
	return svc
}

// DeduplicateRequest 去重请求
type DeduplicateRequest struct {
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	Tags        []string `json:"tags"`
	MediaURLs   []string `json:"media_urls"`
	Strategy    string   `json:"strategy"`    // light, medium, heavy
	TargetSim   float32  `json:"target_sim"`  // 目标相似度 (0-1)
	Seed        int64    `json:"seed"`        // 随机种子
}

// DeduplicateResponse 去重响应
type DeduplicateResponse struct {
	Title         string   `json:"title"`
	Body          string   `json:"body"`
	Tags          []string `json:"tags"`
	MediaURLs     []string `json:"media_urls"`
	Similarity    float32  `json:"similarity"`
	Changes       []string `json:"changes"`
	OriginalLen   int      `json:"original_len"`
	DedupLen      int      `json:"dedup_len"`
}

// Deduplicate 内容去重
func (s *DeduplicationService) Deduplicate(req *DeduplicateRequest) *DeduplicateResponse {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if req.Seed != 0 {
		rng = rand.New(rand.NewSource(req.Seed))
	}

	changes := make([]string, 0)
	originalBody := req.Body

	// 空内容直接返回
	if req.Title == "" && req.Body == "" {
		return &DeduplicateResponse{
			Title:       req.Title,
			Body:        req.Body,
			Tags:        req.Tags,
			MediaURLs:   req.MediaURLs,
			Similarity:  1.0,
			Changes:     changes,
			OriginalLen: 0,
			DedupLen:    0,
		}
	}

	// 根据策略确定处理强度
	strategy := req.Strategy
	if strategy == "" {
		strategy = "medium"
	}

	// 处理标题
	newTitle := s.deduplicateTitle(req.Title, strategy, rng)
	if newTitle != req.Title {
		changes = append(changes, "标题已改写")
	}

	// 处理正文
	newBody := req.Body

	// 1. 段落重排
	if strategy == "medium" || strategy == "heavy" {
		paragraphs := s.splitParagraphs(newBody)
		if len(paragraphs) > 2 {
			reordered := s.reorderParagraphs(paragraphs, strategy, rng)
			if reordered != newBody {
				newBody = reordered
				changes = append(changes, "段落顺序已调整")
			}
		}
	}

	// 2. 同义词替换
	replaced, replaceCount := s.synonymReplace(newBody, strategy, rng)
	if replaceCount > 0 {
		newBody = replaced
		changes = append(changes, fmt.Sprintf("同义词替换 %d 处", replaceCount))
	}

	// 3. 句式变换
	if strategy == "heavy" {
		transformed, transformCount := s.sentenceTransform(newBody, rng)
		if transformCount > 0 {
			newBody = transformed
			changes = append(changes, fmt.Sprintf("句式变换 %d 处", transformCount))
		}
	}

	// 4. 添加差异化标记
	if strategy == "medium" || strategy == "heavy" {
		marked := s.addDifferentiationMarkers(newBody, rng)
		if marked != newBody {
			newBody = marked
			changes = append(changes, "添加差异化标记")
		}
	}

	// 处理标签
	newTags := s.deduplicateTags(req.Tags, strategy, rng)
	if !stringSliceEqual(newTags, req.Tags) {
		changes = append(changes, "标签已调整")
	}

	// 处理媒体URL
	newMediaURLs := s.deduplicateMediaURLs(req.MediaURLs, strategy, rng)
	if !stringSliceEqual(newMediaURLs, req.MediaURLs) {
		changes = append(changes, "媒体素材已差异化")
	}

	// 计算相似度
	similarity := s.calculateSimilarity(originalBody, newBody)

	return &DeduplicateResponse{
		Title:       newTitle,
		Body:        newBody,
		Tags:        newTags,
		MediaURLs:   newMediaURLs,
		Similarity:  similarity,
		Changes:     changes,
		OriginalLen: utf8.RuneCountInString(originalBody),
		DedupLen:    utf8.RuneCountInString(newBody),
	}
}

// ==================== 标题去重 ====================

func (s *DeduplicationService) deduplicateTitle(title, strategy string, rng *rand.Rand) string {
	if title == "" {
		return title
	}

	// 同义词替换
	replaced, count := s.synonymReplace(title, strategy, rng)
	if count > 0 {
		return replaced
	}

	return title
}

// ==================== 段落重排 ====================

func (s *DeduplicationService) splitParagraphs(body string) []string {
	lines := strings.Split(body, "\n")
	var paragraphs []string
	current := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if current != "" {
				paragraphs = append(paragraphs, current)
				current = ""
			}
		} else {
			if current != "" {
				current += "\n"
			}
			current += trimmed
		}
	}
	if current != "" {
		paragraphs = append(paragraphs, current)
	}

	return paragraphs
}

func (s *DeduplicationService) reorderParagraphs(paragraphs []string, strategy string, rng *rand.Rand) string {
	if len(paragraphs) <= 2 {
		return strings.Join(paragraphs, "\n\n")
	}

	// 保留第一个和最后一个段落，重排中间部分
	result := make([]string, len(paragraphs))
	copy(result, paragraphs)

	if strategy == "heavy" {
		// 重排除第一个外的所有段落
		middle := make([]string, len(result)-1)
		copy(middle, result[1:])
		rng.Shuffle(len(middle), func(i, j int) {
			middle[i], middle[j] = middle[j], middle[i]
		})
		result = append(result[:1], middle...)
	} else {
		// 只重排中间段落（保留首尾）
		if len(result) > 3 {
			middle := make([]string, len(result)-2)
			copy(middle, result[1:len(result)-1])
			rng.Shuffle(len(middle), func(i, j int) {
				middle[i], middle[j] = middle[j], middle[i]
			})
			result = append(append([]string{result[0]}, middle...), result[len(result)-1])
		}
	}

	return strings.Join(result, "\n\n")
}

// ==================== 同义词替换 ====================

func (s *DeduplicationService) synonymReplace(text, strategy string, rng *rand.Rand) (string, int) {
	count := 0
	result := text

	// 确定替换概率
	probability := 0.3
	switch strategy {
	case "light":
		probability = 0.2
	case "medium":
		probability = 0.4
	case "heavy":
		probability = 0.6
	}

	for word, synonyms := range s.synonymMap {
		if strings.Contains(result, word) && rng.Float64() < probability {
			synonym := synonyms[rng.Intn(len(synonyms))]
			result = strings.ReplaceAll(result, word, synonym)
			count++
		}
	}

	return result, count
}

// ==================== 句式变换 ====================

func (s *DeduplicationService) sentenceTransform(text string, rng *rand.Rand) (string, int) {
	count := 0
	sentences := splitSentences(text)
	result := make([]string, len(sentences))

	for i, sent := range sentences {
		transformed := sent

		// 变换1: 添加过渡词
		if rng.Float64() < 0.2 {
			transitions := []string{"此外，", "同时，", "另外，", "值得注意的是，", "具体来说，"}
			if !hasTransition(sent) {
				transformed = transitions[rng.Intn(len(transitions))] + transformed
				count++
			}
		}

		// 变换2: 调整语序（简单实现）
		if rng.Float64() < 0.1 && len([]rune(sent)) > 20 {
			// 在句首添加限定词
			qualifiers := []string{"一般来说，", "通常情况下，", "从实践来看，"}
			transformed = qualifiers[rng.Intn(len(qualifiers))] + transformed
			count++
		}

		result[i] = transformed
	}

	return strings.Join(result, ""), count
}

func splitSentences(text string) []string {
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '。' || r == '.' || r == '！' || r == '!' || r == '？' || r == '?'
	})

	var result []string
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s+"。")
		}
	}

	return result
}

func hasTransition(sent string) bool {
	transitions := []string{"此外", "同时", "另外", "因此", "所以", "然而", "但是", "不过", "首先", "其次", "最后"}
	for _, t := range transitions {
		if strings.Contains(sent, t) {
			return true
		}
	}
	return false
}

// ==================== 差异化标记 ====================

func (s *DeduplicationService) addDifferentiationMarkers(body string, rng *rand.Rand) string {
	// 添加不同的引用格式
	markers := []string{
		"\n\n---\n*以上内容仅供参考*",
		"\n\n> 本文观点仅代表作者立场",
		"\n\n【声明】本文内容基于公开信息整理",
	}

	if rng.Float64() < 0.3 {
		return body + markers[rng.Intn(len(markers))]
	}

	return body
}

// ==================== 标签去重 ====================

func (s *DeduplicationService) deduplicateTags(tags []string, strategy string, rng *rand.Rand) []string {
	if len(tags) == 0 {
		return tags
	}

	result := make([]string, len(tags))
	copy(result, tags)

	// 替换部分标签的同义词
	for i, tag := range result {
		if synonyms, ok := s.synonymMap[tag]; ok && rng.Float64() < 0.3 {
			result[i] = synonyms[rng.Intn(len(synonyms))]
		}
	}

	// 打乱标签顺序
	if strategy == "heavy" {
		rng.Shuffle(len(result), func(i, j int) {
			result[i], result[j] = result[j], result[i]
		})
	}

	return result
}

// ==================== 媒体URL去重 ====================

func (s *DeduplicationService) deduplicateMediaURLs(urls []string, strategy string, rng *rand.Rand) []string {
	if len(urls) <= 1 {
		return urls
	}

	result := make([]string, len(urls))
	copy(result, urls)

	// 打乱媒体顺序
	if strategy == "medium" || strategy == "heavy" {
		rng.Shuffle(len(result), func(i, j int) {
			result[i], result[j] = result[j], result[i]
		})
	}

	return result
}

// ==================== 相似度计算 ====================

func (s *DeduplicationService) calculateSimilarity(text1, text2 string) float32 {
	if text1 == text2 {
		return 1.0
	}

	// 使用编辑距离计算相似度
	len1 := utf8.RuneCountInString(text1)
	len2 := utf8.RuneCountInString(text2)

	if len1 == 0 || len2 == 0 {
		return 0
	}

	// 计算共有字符比例
	commonChars := 0
	runes1 := []rune(text1)
	runes2 := []rune(text2)

	charMap := make(map[rune]int)
	for _, r := range runes1 {
		charMap[r]++
	}
	for _, r := range runes2 {
		if charMap[r] > 0 {
			commonChars++
			charMap[r]--
		}
	}

	maxLen := len1
	if len2 > maxLen {
		maxLen = len2
	}

	return float32(commonChars) / float32(maxLen)
}

// ==================== 辅助函数 ====================

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ==================== 同义词库 ====================

func buildSynonymMap() map[string][]string {
	return map[string][]string{
		// 动词
		"提升": {"提高", "增强", "改善", "优化"},
		"提高": {"提升", "增强", "改善", "优化"},
		"增强": {"提升", "提高", "强化", "加强"},
		"优化": {"改进", "完善", "提升", "改善"},
		"实现": {"达成", "完成", "达到", "做到"},
		"提供": {"供给", "供应", "给予", "带来"},
		"使用": {"采用", "运用", "利用", "应用"},
		"采用": {"使用", "运用", "利用", "应用"},
		"支持": {"支撑", "兼容", "承载", "适配"},
		"包含": {"涵盖", "包括", "含有", "囊括"},
		"需要": {"要求", "需", "必备", "必须"},
		"建议": {"推荐", "提倡", "倡导", "提议"},
		"分析": {"解析", "研究", "剖析", "解读"},
		"创建": {"建立", "构建", "搭建", "生成"},
		"管理": {"管控", "治理", "运维", "运营"},

		// 形容词
		"重要": {"关键", "核心", "主要", "关键性"},
		"有效": {"高效", "有力", "显著", "明显"},
		"快速": {"迅速", "高效", "敏捷", "即时"},
		"简单": {"简便", "便捷", "易用", "轻松"},
		"完整": {"全面", "完善", "齐全", "完备"},
		"强大": {"强劲", "出色", "卓越", "优秀"},
		"灵活": {"弹性", "可扩展", "可定制", "柔性"},
		"稳定": {"可靠", "稳固", "牢固", "持续"},
		"安全": {"可靠", "受保护", "无风险", "可信"},
		"智能": {"智慧", "自动化", "AI驱动", "智能化"},

		// 名词
		"功能": {"特性", "能力", "模块", "组件"},
		"性能": {"效率", "表现", "速度", "指标"},
		"用户": {"使用者", "客户", "操作者", "终端用户"},
		"系统": {"平台", "方案", "工具", "产品"},
		"数据": {"信息", "资料", "内容", "素材"},
		"结果": {"成果", "效果", "产出", "输出"},
		"方法": {"方式", "手段", "途径", "策略"},
		"问题": {"挑战", "难点", "痛点", "困难"},
		"优势": {"优点", "长处", "亮点", "特色"},
		"目标": {"目的", "方向", "愿景", "指标"},

		// 副词/连接词
		"非常": {"十分", "极其", "相当", "特别"},
		"因此": {"所以", "故而", "由此", "于是"},
		"但是": {"然而", "不过", "可是", "只是"},
		"而且": {"并且", "同时", "此外", "另外"},
		"如果": {"假如", "倘若", "若是", "假设"},
		"可以": {"能够", "能够", "得以", "可"},
		"已经": {"已", "已然", "早已", "业已"},

		// GEO相关
		"搜索引擎": {"搜索平台", "检索引擎", "搜索工具"},
		"内容优化": {"内容改进", "内容提升", "内容完善"},
		"关键词": {"搜索词", "核心词", "目标词"},
		"流量": {"访问量", "曝光量", "用户量"},
		"转化": {"转化率", "成交", "转化效果"},
		"排名": {"排序", "位置", "名次"},
		"收录": {"索引", "抓取", "入库"},
		"权重": {"权威度", "信任度", "评分"},
	}
}
