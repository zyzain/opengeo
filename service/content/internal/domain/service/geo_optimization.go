package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"opengeo/pkg/ai"
	"opengeo/service/content/internal/domain/model"
)

var (
	reMarkdownH2 = regexp.MustCompile(`(?m)^##\s+`)
	reMarkdownH3 = regexp.MustCompile(`(?m)^###\s+`)
	reMarkdownH4 = regexp.MustCompile(`(?m)^####\s+`)
	reHTMLTags   = regexp.MustCompile(`</?(h[1-6]|p|ul|ol|li|table|blockquote|pre|img)[^>]*>`)
	reURL        = regexp.MustCompile(`https?://[^\s\)]+`)
	reNumber     = regexp.MustCompile(`\b\d+(?:\.\d+)?%?\b`)
	reDate       = regexp.MustCompile(`\b\d{4}[-/年]\d{1,2}[-/月]\d{1,2}\b|\b\d{4}年\b`)
	reQuestion   = regexp.MustCompile(`(?m)^.*[？?].*$`)
	reListItem   = regexp.MustCompile(`(?m)^[\-\*\d]+[.)、]\s+`)
	reQuote      = regexp.MustCompile(`(?m)^>\s+`)
)

var chineseStopWords = map[string]struct{}{
	"的": {}, "了": {}, "在": {}, "是": {}, "和": {}, "与": {},
	"或": {}, "等": {}, "这": {}, "那": {}, "我": {}, "你": {},
	"他": {}, "她": {}, "它": {}, "们": {}, "个": {}, "之": {},
	"而": {}, "但": {}, "也": {}, "就": {}, "都": {}, "要": {},
	"会": {}, "能": {}, "可": {}, "可以": {}, "没有": {}, "有": {},
	"不": {}, "没": {}, "很": {}, "还": {}, "又": {}, "再": {},
	"才": {}, "已": {}, "已经": {}, "被": {}, "将": {}, "把": {},
	"从": {}, "到": {}, "对": {}, "为": {}, "所": {}, "其": {},
	"中": {}, "上": {}, "下": {}, "里": {}, "让": {}, "给": {},
	"着": {}, "过": {}, "地": {}, "得": {}, "来": {}, "去": {},
	"一": {}, "二": {}, "三": {}, "四": {}, "五": {}, "六": {},
	"七": {}, "八": {}, "九": {}, "十": {}, "百": {}, "千": {},
	"万": {}, "亿": {},
}

// chineseDictionary 常用中文复合词词典（按长度降序匹配）
var chineseDictionary = []string{
	// AI/技术 (4+ chars)
	"自然语言处理", "计算机视觉", "语音识别", "生成式人工智能",
	"大语言模型", "知识图谱", "数据挖掘", "人机交互",
	// AI/技术 (3 chars)
	"人工智能", "机器学习", "深度学习", "神经网络", "强化学习",
	"推荐系统", "搜索引擎", "区块链", "物联网", "虚拟现实",
	"增强现实", "云计算", "大数据", "数据库", "操作系统",
	// 营销/商业 (4+ chars)
	"内容营销", "数字营销", "品牌推广", "社交媒体", "用户增长",
	"搜索引擎优化", "搜索引擎营销", "投资回报率", "客户关系管理",
	// 营销/商业 (3 chars)
	"转化率", "点击率", "跳出率", "留存率", "活跃度",
	"电子商务", "用户体验", "数据分析", "数据驱动",
	"商业模式", "供应链", "产品经理", "市场营销",
	"品牌建设", "用户画像", "流量获取", "私域流量",
	"公域流量", "增长黑客", "内容创作", "用户需求",
	// SEO/内容 (3 chars)
	"关键词", "长尾词", "外链", "内链", "权重",
	"收录", "排名", "流量", "优化", "算法",
	"索引", "爬虫", "标签", "摘要", "元数据",
	"标题", "正文", "段落", "列表", "引用",
	// 通用高频 (2 chars)
	"技术", "产品", "服务", "方案", "系统", "平台",
	"数据", "信息", "内容", "用户", "市场", "行业",
	"企业", "公司", "品牌", "业务", "应用", "开发",
	"管理", "运营", "增长", "策略", "效果", "价值",
	"工具", "资源", "功能", "需求", "问题", "解决",
	"分析", "研究", "创新", "趋势", "未来", "智能",
	"全球", "中国", "行业", "领域", "方面", "方式",
	"方法", "技巧", "指南", "教程", "案例", "实践",
}

// chineseMaxWordLen 词典中最长词的长度
const chineseMaxWordLen = 6

// chineseWordSet 词典集合，用于快速查找
var chineseWordSet map[string]struct{}

func init() {
	chineseWordSet = make(map[string]struct{}, len(chineseDictionary))
	for _, w := range chineseDictionary {
		chineseWordSet[w] = struct{}{}
	}
}

// tokenizeChinese 对中文文本进行分词
// 优先匹配词典中最长匹配，未匹配部分回退到 bigram 切分
// 英文单词保持原样，过滤单字虚词
func tokenizeChinese(text string) []string {
	runes := []rune(text)
	length := len(runes)
	var tokens []string

	isChinese := func(r rune) bool {
		return r >= 0x4E00 && r <= 0x9FFF
	}
	isLetter := func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	}
	isDigit := func(r rune) bool {
		return r >= '0' && r <= '9'
	}

	i := 0
	for i < length {
		r := runes[i]

		// 英文+数字单词：连续读取
		if isLetter(r) || isDigit(r) {
			j := i + 1
			for j < length && (isLetter(runes[j]) || isDigit(runes[j])) {
				j++
			}
			word := string(runes[i:j])
			if _, stop := chineseStopWords[strings.ToLower(word)]; !stop {
				tokens = append(tokens, word)
			}
			i = j
			continue
		}

		// 中文字符：尝试最长匹配
		if isChinese(r) {
			matched := false
			for wordLen := chineseMaxWordLen; wordLen >= 2; wordLen-- {
				if i+wordLen > length {
					continue
				}
				candidate := string(runes[i : i+wordLen])
				// 检查候选词所有字符都是中文
				allChinese := true
				for _, c := range candidate {
					if !isChinese(c) {
						allChinese = false
						break
					}
				}
				if allChinese {
					if _, ok := chineseWordSet[candidate]; ok {
						tokens = append(tokens, candidate)
						i += wordLen
						matched = true
						break
					}
				}
			}
			if !matched {
				// 回退：跳过标点等非中文非字母数字字符
				// 如果是中文单字，尝试 bigram
				if i+1 < length && isChinese(runes[i+1]) {
					bigram := string(runes[i : i+2])
					if _, isStop := chineseStopWords[bigram]; !isStop {
						tokens = append(tokens, bigram)
					}
					i += 2
				} else {
					// 单个中文字符，跳过（过滤单字）
					i++
				}
			}
			continue
		}

		// 其他字符（标点、空格等）跳过
		i++
	}

	return tokens
}

var authorityDomains = []string{
	"wikipedia.org", "gov.cn", "gov", "edu.cn", "edu",
	"github.com", "stackoverflow.com", "arxiv.org",
	"nature.com", "science.org", "ieee.org",
	"who.int", "un.org",
	"zhihu.com", "csdn.net", "segmentfault.com",
	"baike.baidu.com", "wikipedia",
}

// GEOOptimizationService GEO优化领域服务
type GEOOptimizationService struct {
	authorityDomainSet map[string]struct{}
}

// NewGEOOptimizationService 创建GEO优化领域服务
func NewGEOOptimizationService() *GEOOptimizationService {
	s := &GEOOptimizationService{
		authorityDomainSet: make(map[string]struct{}, len(authorityDomains)),
	}
	for _, d := range authorityDomains {
		s.authorityDomainSet[d] = struct{}{}
	}
	return s
}

// OptimizeForAI 为AI优化内容
func (s *GEOOptimizationService) OptimizeForAI(ctx context.Context, content *model.Content) (*OptimizationResult, error) {
	if !content.IsValid() {
		return nil, fmt.Errorf("content is not valid")
	}

	structScore, structDetails := s.analyzeStructure(content)
	readScore, readDetails := s.analyzeReadability(content)
	intentScore := s.analyzeAIIntent(content)
	schemaScore, generatedSchema := s.analyzeAndGenerateSchema(content)

	totalScore := structScore*0.3 + readScore*0.3 + intentScore*0.2 + schemaScore*0.2
	totalScore = float32(math.Min(float64(totalScore), 100))

	suggestions := s.generateSuggestions(content, structDetails, readDetails, schemaScore)

	return &OptimizationResult{
		ContentID:          content.ID,
		StructureScore:     structScore,
		ReadabilityScore:   readScore,
		TotalScore:         totalScore,
		Suggestions:        suggestions,
		SchemaMarkup:       generatedSchema,
		StructuralChanges:  structDetails,
		ReadabilityDetails: readDetails,
	}, nil
}

// AdaptForModel 为特定AI模型适配内容
func (s *GEOOptimizationService) AdaptForModel(ctx context.Context, content *model.Content, targetModel string) (*AdaptationResult, error) {
	if !content.IsValid() {
		return nil, fmt.Errorf("content is not valid")
	}

	var adaptedBody string
	var formatChanges []string

	switch targetModel {
	case "deepseek":
		adaptedBody, formatChanges = s.adaptForDeepSeek(content)
	case "kimi":
		adaptedBody, formatChanges = s.adaptForKimi(content)
	case "doubao":
		adaptedBody, formatChanges = s.adaptForDoubao(content)
	case "chatgpt":
		adaptedBody, formatChanges = s.adaptForChatGPT(content)
	default:
		adaptedBody = content.Body
		formatChanges = []string{"未识别的目标模型，保持原始内容"}
	}

	return &AdaptationResult{
		ContentID:      content.ID,
		TargetModel:    targetModel,
		AdaptedContent: adaptedBody,
		FormatChanges:  formatChanges,
	}, nil
}

// CheckCompliance 检查内容合规性
func (s *GEOOptimizationService) CheckCompliance(ctx context.Context, content *model.Content) (*ComplianceResult, error) {
	if !content.IsValid() {
		return nil, fmt.Errorf("content is not valid")
	}

	issues := make([]ai.ComplianceIssue, 0)
	issues = append(issues, s.checkSensitiveWords(content)...)
	issues = append(issues, s.checkAdLawCompliance(content)...)
	issues = append(issues, s.checkAIGCLabeling(content)...)

	return &ComplianceResult{
		ContentID: content.ID,
		Compliant: len(issues) == 0,
		Issues:    issues,
		Report:    s.generateReport(content, issues),
	}, nil
}

// ==================== 结构分析 ====================

func (s *GEOOptimizationService) analyzeStructure(content *model.Content) (float32, []string) {
	score := float32(0)
	details := make([]string, 0)
	body := content.Body

	// 标题层级分析 (30分)
	h2Count := len(reMarkdownH2.FindAllString(body, -1))
	h3Count := len(reMarkdownH3.FindAllString(body, -1))
	htmlHeadingCount := len(reHTMLTags.FindAllString(body, -1))

	if h2Count > 0 || h3Count > 0 || htmlHeadingCount > 0 {
		score += 15
		details = append(details, "检测到标题层级结构")
		if h2Count > 0 && h3Count > 0 {
			score += 15
			details = append(details, fmt.Sprintf("标题层级完整（H2: %d, H3: %d）", h2Count, h3Count))
		}
	} else {
		details = append(details, "缺少标题层级结构，建议使用 ## / ### 标记")
	}

	// 段落结构分析 (25分)
	paragraphs := s.splitParagraphs(body)
	if len(paragraphs) >= 3 {
		score += 10
		details = append(details, fmt.Sprintf("段落数量充足（%d段）", len(paragraphs)))

		idealCount := 0
		for _, p := range paragraphs {
			wordCount := utf8.RuneCountInString(p)
			if wordCount >= 50 && wordCount <= 300 {
				idealCount++
			}
		}
		idealRatio := float32(idealCount) / float32(len(paragraphs))
		if idealRatio >= 0.6 {
			score += 15
			details = append(details, "段落长度分布合理（50-300字/段）")
		} else {
			score += 5
			details = append(details, "部分段落过长或过短，建议控制在50-300字")
		}
	} else {
		details = append(details, "段落数量不足，建议至少分为3段")
	}

	// 列表/结构化元素检测 (20分)
	listItems := reListItem.FindAllString(body, -1)
	if len(listItems) > 0 {
		score += 10
		details = append(details, fmt.Sprintf("包含列表结构（%d项）", len(listItems)))
	}
	blockquotes := reQuote.FindAllString(body, -1)
	if len(blockquotes) > 0 {
		score += 10
		details = append(details, fmt.Sprintf("包含引用块（%d处）", len(blockquotes)))
	}

	// 正文长度评估 (25分)
	bodyLen := utf8.RuneCountInString(body)
	switch {
	case bodyLen >= 2000:
		score += 25
		details = append(details, fmt.Sprintf("正文长度充足（%d字）", bodyLen))
	case bodyLen >= 1000:
		score += 20
		details = append(details, fmt.Sprintf("正文长度适中（%d字），建议≥2000字", bodyLen))
	case bodyLen >= 500:
		score += 10
		details = append(details, fmt.Sprintf("正文较短（%d字），建议≥1000字", bodyLen))
	default:
		details = append(details, fmt.Sprintf("正文过短（%d字），建议≥500字", bodyLen))
	}

	return float32(math.Min(float64(score), 100)), details
}

// ==================== 可读性分析 ====================

func (s *GEOOptimizationService) analyzeReadability(content *model.Content) (float32, []string) {
	score := float32(0)
	details := make([]string, 0)
	body := content.Body

	// 标题关键词在正文中的出现密度 (30分)
	keywordScore := s.analyzeKeywordDensity(content.Title, body)
	score += keywordScore * 0.3
	if keywordScore >= 70 {
		details = append(details, fmt.Sprintf("标题关键词密度良好（%.0f分）", keywordScore))
	} else {
		details = append(details, fmt.Sprintf("标题关键词密度偏低（%.0f分），建议在正文中多次提及核心关键词", keywordScore))
	}

	// 事实性数据密度 (30分)
	factualScore := s.analyzeFactualDensity(body)
	score += factualScore * 0.3
	if factualScore >= 60 {
		details = append(details, "事实性数据密度良好")
	} else {
		details = append(details, "事实性数据不足，建议添加具体数字、百分比、日期等")
	}

	// 权威引用检测 (20分)
	authScore, authCount := s.analyzeAuthorityReferences(body)
	score += authScore * 0.2
	if authCount > 0 {
		details = append(details, fmt.Sprintf("包含 %d 处权威引用", authCount))
	} else {
		details = append(details, "缺少权威引用，建议链接到权威信源")
	}

	// 句子长度分析 (20分)
	sentScore := s.analyzeSentenceStructure(body)
	score += sentScore * 0.2
	if sentScore >= 70 {
		details = append(details, "句子结构清晰")
	} else {
		details = append(details, "部分句子过长，建议拆分为短句以提升AI可读性")
	}

	return float32(math.Min(float64(score), 100)), details
}

func (s *GEOOptimizationService) analyzeKeywordDensity(title, body string) float32 {
	if title == "" || body == "" {
		return 0
	}

	// 提取标题中的关键词（按空格/标点分词，过滤停用词）
	keywords := s.extractKeywords(title)
	if len(keywords) == 0 {
		return 30
	}

	bodyLower := strings.ToLower(body)
	totalHits := 0
	for _, kw := range keywords {
		kwLower := strings.ToLower(kw)
		totalHits += strings.Count(bodyLower, kwLower)
	}

	// 每个关键词至少出现2次为理想
	expectedMin := len(keywords) * 2
	ratio := float32(totalHits) / float32(expectedMin)
	return float32(math.Min(float64(ratio*100), 100))
}

func (s *GEOOptimizationService) extractKeywords(text string) []string {
	words := tokenizeChinese(text)

	var keywords []string
	for _, w := range words {
		w = strings.TrimSpace(w)
		if utf8.RuneCountInString(w) < 2 {
			continue
		}
		if _, isStop := chineseStopWords[w]; isStop {
			continue
		}
		if _, isStop := chineseStopWords[strings.ToLower(w)]; isStop {
			continue
		}
		keywords = append(keywords, w)
	}

	// 去重
	seen := make(map[string]struct{})
	var unique []string
	for _, kw := range keywords {
		lower := strings.ToLower(kw)
		if _, exists := seen[lower]; !exists {
			seen[lower] = struct{}{}
			unique = append(unique, kw)
		}
	}

	return unique
}

func (s *GEOOptimizationService) analyzeFactualDensity(body string) float32 {
	if body == "" {
		return 0
	}

	totalWords := float32(utf8.RuneCountInString(body))
	if totalWords == 0 {
		return 0
	}

	// 数字/百分比检测
	numberMatches := reNumber.FindAllString(body, -1)
	// 日期检测
	dateMatches := reDate.FindAllString(body, -1)

	factualCount := float32(len(numberMatches) + len(dateMatches))
	// 每200字至少1个事实性数据为理想
	idealRatio := totalWords / 200
	if idealRatio < 1 {
		idealRatio = 1
	}

	ratio := factualCount / idealRatio
	return float32(math.Min(float64(ratio*100), 100))
}

func (s *GEOOptimizationService) analyzeAuthorityReferences(body string) (float32, int) {
	urls := reURL.FindAllString(body, -1)
	if len(urls) == 0 {
		return 0, 0
	}

	authCount := 0
	for _, u := range urls {
		uLower := strings.ToLower(u)
		for domain := range s.authorityDomainSet {
			if strings.Contains(uLower, domain) {
				authCount++
				break
			}
		}
	}

	// 至少2个权威引用为理想
	score := float32(authCount) / 2 * 100
	return float32(math.Min(float64(score), 100)), authCount
}

func (s *GEOOptimizationService) analyzeSentenceStructure(body string) float32 {
	sentences := strings.FieldsFunc(body, func(r rune) bool {
		return r == '。' || r == '.' || r == '！' || r == '!' || r == '？' || r == '?'
	})

	if len(sentences) == 0 {
		return 50
	}

	totalLen := 0
	longSentences := 0
	for _, sent := range sentences {
		sentLen := utf8.RuneCountInString(strings.TrimSpace(sent))
		totalLen += sentLen
		if sentLen > 80 {
			longSentences++
		}
	}

	avgLen := totalLen / len(sentences)
	score := float32(80)

	// 平均句长30-60字为理想
	if avgLen > 60 {
		penalty := float32(avgLen-60) * 1.5
		score -= penalty
	}
	if avgLen < 10 {
		score -= 20
	}

	// 过长句子比例惩罚
	longRatio := float32(longSentences) / float32(len(sentences))
	if longRatio > 0.3 {
		score -= longRatio * 30
	}

	return float32(math.Max(math.Min(float64(score), 100), 0))
}

// ==================== AI意图匹配分析 ====================

func (s *GEOOptimizationService) analyzeAIIntent(content *model.Content) float32 {
	score := float32(0)
	body := content.Body

	// 问答模式检测 - AI偏好结构化问答
	questions := reQuestion.FindAllString(body, -1)
	if len(questions) > 0 {
		score += 25
	}

	// 总结/结论检测 - AI偏好内容有明确结论
	conclusionKeywords := []string{"总结", "结论", "综上", "总的来说", "总而言之", "summary", "conclusion", "in summary", "to sum up"}
	bodyLower := strings.ToLower(body)
	for _, kw := range conclusionKeywords {
		if strings.Contains(bodyLower, strings.ToLower(kw)) {
			score += 25
			break
		}
	}

	// 列表/要点检测 - AI偏好提取要点
	if reListItem.MatchString(body) {
		score += 25
	}

	// 数据支撑检测 - AI偏好有数据的内容
	if reNumber.MatchString(body) {
		score += 25
	}

	return float32(math.Min(float64(score), 100))
}

// ==================== Schema Markup 分析与生成 ====================

func (s *GEOOptimizationService) analyzeAndGenerateSchema(content *model.Content) (float32, string) {
	score := float32(0)

	// 分析现有 Schema Markup 完整度
	if content.SchemaMarkup != "" {
		completeness := s.evaluateSchemaCompleteness(content.SchemaMarkup)
		score = completeness
		return score, content.SchemaMarkup
	}

	// 自动生成 Schema Markup
	generated := s.generateSchemaMarkup(content, "")
	score = 60
	return score, generated
}

func (s *GEOOptimizationService) evaluateSchemaCompleteness(schemaMarkup string) float32 {
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaMarkup), &schema); err != nil {
		return 20
	}

	score := float32(0)

	// 基本结构完整性
	if _, ok := schema["@context"]; ok {
		score += 15
	}
	if _, ok := schema["@type"]; ok {
		score += 15
	}

	// 关键字段完整性
	requiredFields := []string{"headline", "name", "description", "author", "datePublished"}
	for _, field := range requiredFields {
		if _, ok := schema[field]; ok {
			score += 10
		}
	}

	// articleBody 或 text
	if _, ok := schema["articleBody"]; ok {
		score += 10
	} else if _, ok := schema["text"]; ok {
		score += 10
	}

	return float32(math.Min(float64(score), 100))
}

func (s *GEOOptimizationService) generateSchemaMarkup(content *model.Content, authorName string) string {
	title := content.Title
	body := content.Body

	if authorName == "" {
		authorName = "OpenGEO"
	}

	// 截取摘要
	summary := body
	runes := []rune(body)
	if len(runes) > 200 {
		summary = string(runes[:200]) + "..."
	}

	now := time.Now()
	publishDate := now.Format("2006-01-02")

	schemaType := "Article"
	if content.ContentType == model.ContentTypeVideo {
		schemaType = "VideoObject"
	} else if content.ContentType == model.ContentTypeImage {
		schemaType = "ImageObject"
	}

	schema := map[string]interface{}{
		"@context":       "https://schema.org",
		"@type":          schemaType,
		"headline":       title,
		"description":    summary,
		"articleBody":    body,
		"datePublished":  publishDate,
		"dateModified":   publishDate,
		"author": map[string]interface{}{
			"@type": "Organization",
			"name":  authorName,
		},
		"publisher": map[string]interface{}{
			"@type": "Organization",
			"name":  "OpenGEO",
		},
		"mainEntityOfPage": map[string]interface{}{
			"@type": "WebPage",
			"@id":   fmt.Sprintf("https://opengeo.ai/content/%d", content.ID),
		},
	}

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}

// ==================== 优化建议生成 ====================

func (s *GEOOptimizationService) generateSuggestions(content *model.Content, structDetails, readDetails []string, schemaScore float32) []string {
	suggestions := make([]string, 0)

	if content.SchemaMarkup == "" {
		suggestions = append(suggestions, "【P0】添加Schema Markup结构化数据，已自动生成建议版本，请检查并补充")
	} else if schemaScore < 80 {
		suggestions = append(suggestions, "【P1】现有Schema Markup字段不完整，建议补充 author/datePublished/description 字段")
	}

	bodyLen := utf8.RuneCountInString(content.Body)
	if bodyLen < 500 {
		suggestions = append(suggestions, "【P0】正文过短，建议至少500字以上以提升AI引用率")
	} else if bodyLen < 1000 {
		suggestions = append(suggestions, "【P1】正文长度适中，建议≥1000字以提升权威性")
	}

	// 基于结构分析的建议
	hasHeading := false
	for _, d := range structDetails {
		if strings.Contains(d, "标题层级完整") {
			hasHeading = true
			break
		}
	}
	if !hasHeading {
		suggestions = append(suggestions, "【P0】添加 ## / ### 标题层级结构，帮助AI理解内容层次")
	}

	hasList := false
	for _, d := range structDetails {
		if strings.Contains(d, "列表结构") {
			hasList = true
			break
		}
	}
	if !hasList {
		suggestions = append(suggestions, "【P1】使用列表（有序/无序）呈现要点，便于AI提取关键信息")
	}

	// 基于可读性分析的建议
	hasFactual := false
	for _, d := range readDetails {
		if strings.Contains(d, "事实性数据密度良好") {
			hasFactual = true
			break
		}
	}
	if !hasFactual {
		suggestions = append(suggestions, "【P1】添加具体数据（数字/百分比/日期），提升内容的可信度")
	}

	hasAuth := false
	for _, d := range readDetails {
		if strings.Contains(d, "权威引用") && !strings.Contains(d, "缺少") {
			hasAuth = true
			break
		}
	}
	if !hasAuth {
		suggestions = append(suggestions, "【P1】添加权威信源引用链接（如百科/官网/学术文献），提升AI信任度")
	}

	// 检查问答模式
	bodyLower := strings.ToLower(content.Body)
	qaKeywords := []string{"？", "?", "什么", "如何", "为什么", "what", "how", "why"}
	hasQA := false
	for _, kw := range qaKeywords {
		if strings.Contains(bodyLower, kw) {
			hasQA = true
			break
		}
	}
	if !hasQA {
		suggestions = append(suggestions, "【P2】添加问答式段落（如「什么是XX？」），匹配AI搜索意图模式")
	}

	conclusionKeywords := []string{"总结", "结论", "综上", "总而言之", "summary", "conclusion"}
	hasConclusion := false
	for _, kw := range conclusionKeywords {
		if strings.Contains(bodyLower, kw) {
			hasConclusion = true
			break
		}
	}
	if !hasConclusion {
		suggestions = append(suggestions, "【P2】添加总结/结论段落，便于AI直接引用核心观点")
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "内容GEO优化良好，继续保持")
	}

	return suggestions
}

// ==================== 模型适配 ====================

func (s *GEOOptimizationService) adaptForDeepSeek(content *model.Content) (string, []string) {
	changes := make([]string, 0)
	body := content.Body

	// DeepSeek偏好：结构化、有明确结论、分点论述
	// 1. 确保有结论段
	conclusionKeywords := []string{"总结", "结论", "综上"}
	hasConclusion := false
	bodyLower := strings.ToLower(body)
	for _, kw := range conclusionKeywords {
		if strings.Contains(bodyLower, kw) {
			hasConclusion = true
			break
		}
	}

	adaptedBody := body
	if !hasConclusion {
		adaptedBody += "\n\n## 总结\n\n以上内容从多个维度进行了分析，具体结论请参考各章节详细说明。"
		changes = append(changes, "添加了总结段落（DeepSeek偏好明确结论）")
	}

	// 2. 检查是否有分点
	if !reListItem.MatchString(body) {
		changes = append(changes, "建议将关键信息以列表形式呈现（DeepSeek偏好分点论述）")
	}

	changes = append(changes, "适配DeepSeek模型格式：强调结构化与结论导向")
	return adaptedBody, changes
}

func (s *GEOOptimizationService) adaptForKimi(content *model.Content) (string, []string) {
	changes := make([]string, 0)
	body := content.Body

	// Kimi偏好：长文本、详细分析、数据支撑
	bodyLen := utf8.RuneCountInString(body)
	adaptedBody := body

	if bodyLen < 1000 {
		adaptedBody += "\n\n---\n\n> 本文内容经过AI优化，建议补充更多细节数据与案例分析以提升引用率。"
		changes = append(changes, "添加了补充提示（Kimi偏好长文本详细分析）")
	}

	changes = append(changes, "适配Kimi模型格式：强调详细分析与数据密度")
	return adaptedBody, changes
}

func (s *GEOOptimizationService) adaptForDoubao(content *model.Content) (string, []string) {
	changes := make([]string, 0)
	body := content.Body

	// 豆包偏好：简洁明了、有数据支撑
	adaptedBody := body

	// 检查是否有要点提炼
	if !reListItem.MatchString(body) {
		adaptedBody = "## 要点提炼\n\n" + adaptedBody
		changes = append(changes, "添加了要点提炼标题（豆包偏好简洁结构）")
	}

	changes = append(changes, "适配豆包模型格式：强调简洁与数据支撑")
	return adaptedBody, changes
}

func (s *GEOOptimizationService) adaptForChatGPT(content *model.Content) (string, []string) {
	changes := make([]string, 0)
	body := content.Body

	// ChatGPT偏好：对话式、有上下文、FAQ格式
	adaptedBody := body

	// 检查是否有FAQ风格
	questions := reQuestion.FindAllString(body, -1)
	if len(questions) == 0 {
		adaptedBody += "\n\n## 常见问题\n\n**Q: 以上内容的核心要点是什么？**\n\nA: 请参考文中各章节详细说明。"
		changes = append(changes, "添加了FAQ段落（ChatGPT偏好对话式结构）")
	}

	changes = append(changes, "适配ChatGPT模型格式：强调对话式与上下文连贯性")
	return adaptedBody, changes
}

// ==================== 合规检测 ====================

var sensitiveWords = []string{
	"赌博", "博彩", "色情", "暴力", "恐怖", "毒品",
	"枪支", "弹药", "管制刀具", "传销", "诈骗",
	"非法集资", "洗钱", "走私", "贩卖",
}

var adLawWords = []string{
	"最好", "第一", "顶级", "绝对", "100%", "永远",
	"最强", "最佳", "最优秀", "最先进", "独一无二",
	"史上最", "全球第一", "行业领先", "遥遥领先",
	"best", "No.1", "the best", "top 1", "number one",
}

var aigcLabelKeywords = []string{
	"AI生成", "AIGC", "人工智能辅助", "AI辅助",
	"generated by AI", "AI-generated", "AIGC content",
}

func (s *GEOOptimizationService) checkSensitiveWords(content *model.Content) []ai.ComplianceIssue {
	issues := make([]ai.ComplianceIssue, 0)
	bodyLower := strings.ToLower(content.Body)
	titleLower := strings.ToLower(content.Title)
	fullText := titleLower + " " + bodyLower

	for _, word := range sensitiveWords {
		if strings.Contains(fullText, strings.ToLower(word)) {
			issues = append(issues, ai.ComplianceIssue{
				IssueType:   "sensitive",
				Description: fmt.Sprintf("检测到敏感词「%s」", word),
				Severity:    "high",
				Suggestion:  fmt.Sprintf("请删除或替换敏感词「%s」", word),
			})
		}
	}

	return issues
}

func (s *GEOOptimizationService) checkAdLawCompliance(content *model.Content) []ai.ComplianceIssue {
	issues := make([]ai.ComplianceIssue, 0)
	bodyLower := strings.ToLower(content.Body)
	titleLower := strings.ToLower(content.Title)
	fullText := titleLower + " " + bodyLower

	for _, word := range adLawWords {
		if strings.Contains(fullText, strings.ToLower(word)) {
			issues = append(issues, ai.ComplianceIssue{
				IssueType:   "ad_law",
				Description: fmt.Sprintf("检测到广告法敏感词「%s」，可能违反《广告法》规定", word),
				Severity:    "medium",
				Suggestion:  fmt.Sprintf("建议替换「%s」为更客观的表述", word),
			})
		}
	}

	return issues
}

func (s *GEOOptimizationService) checkAIGCLabeling(content *model.Content) []ai.ComplianceIssue {
	issues := make([]ai.ComplianceIssue, 0)
	bodyLower := strings.ToLower(content.Body)

	hasLabel := false
	for _, kw := range aigcLabelKeywords {
		if strings.Contains(bodyLower, strings.ToLower(kw)) {
			hasLabel = true
			break
		}
	}

	if !hasLabel {
		issues = append(issues, ai.ComplianceIssue{
			IssueType:   "aigc_label",
			Description: "内容缺少AIGC标识，可能被AI模型降权或不符合监管要求",
			Severity:    "low",
			Suggestion:  "建议在内容末尾添加AIGC标识，如「本文由AI辅助生成」",
		})
	}

	return issues
}

func (s *GEOOptimizationService) generateReport(content *model.Content, issues []ai.ComplianceIssue) string {
	if len(issues) == 0 {
		return fmt.Sprintf("合规检测通过\n内容ID: %d\n标题: %s\n检测时间: %s",
			content.ID, content.Title, time.Now().Format("2006-01-02 15:04:05"))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("合规检测报告\n内容ID: %d\n标题: %s\n检测时间: %s\n\n",
		content.ID, content.Title, time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("发现问题: %d 个\n\n", len(issues)))

	highCount, mediumCount, lowCount := 0, 0, 0
	for _, issue := range issues {
		switch issue.Severity {
		case "high":
			highCount++
		case "medium":
			mediumCount++
		case "low":
			lowCount++
		}
	}

	sb.WriteString(fmt.Sprintf("高危: %d | 中危: %d | 低危: %d\n\n", highCount, mediumCount, lowCount))

	for i, issue := range issues {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n   建议: %s\n", i+1, issue.Severity, issue.Description, issue.Suggestion))
	}

	return sb.String()
}

// ==================== 辅助方法 ====================

func (s *GEOOptimizationService) splitParagraphs(body string) []string {
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

// ==================== 结果类型 ====================

// OptimizationResult 优化结果
type OptimizationResult struct {
	ContentID          int64    `json:"content_id"`
	StructureScore     float32  `json:"structure_score"`
	ReadabilityScore   float32  `json:"readability_score"`
	TotalScore         float32  `json:"total_score"`
	Suggestions        []string `json:"suggestions"`
	SchemaMarkup       string   `json:"schema_markup"`
	StructuralChanges  []string `json:"structural_changes"`
	ReadabilityDetails []string `json:"readability_details"`
}

// AdaptationResult 适配结果
type AdaptationResult struct {
	ContentID      int64    `json:"content_id"`
	TargetModel    string   `json:"target_model"`
	AdaptedContent string   `json:"adapted_content"`
	FormatChanges  []string `json:"format_changes"`
}

// ComplianceResult 合规结果
type ComplianceResult struct {
	ContentID int64                `json:"content_id"`
	Compliant bool                 `json:"compliant"`
	Issues    []ai.ComplianceIssue `json:"issues"`
	Report    string               `json:"report"`
}
