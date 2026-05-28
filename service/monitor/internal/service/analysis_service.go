package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// CompetitorAnalysisService 竞品分析服务
type CompetitorAnalysisService struct {
	citationRepo   CitationRepo
	competitorRepo CompetitorRepo
	scoreRepo      ScoreRepo
}

type CitationRepo interface {
	GetCitationStats(ctx context.Context, contentID int64) (*CitationStats, error)
	List(ctx context.Context, contentID int64, aiModel string, page, pageSize int) ([]interface{}, int32, error)
}

type CompetitorRepo interface {
	Create(ctx context.Context, monitor interface{}) error
	GetByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context, userID int64, page, pageSize int) ([]interface{}, int32, error)
}

type ScoreRepo interface {
	List(ctx context.Context, channelID, accountID int64) ([]interface{}, error)
}

type CitationStats struct {
	TotalCitations int32
	AvgPosition    float32
	CitationRate   float32
	TopModels      []string
}

// CompetitorData 竞品数据
type CompetitorData struct {
	Name            string  `json:"name"`
	Domain          string  `json:"domain"`
	VisibilityScore float32 `json:"visibility_score"`
	CitationCount   int32   `json:"citation_count"`
	AvgPosition     float32 `json:"avg_position"`
	TopQueries      []string `json:"top_queries"`
	ContentCount    int32   `json:"content_count"`
}

// ComparisonReport 对比报告
type ComparisonReport struct {
	YourBrand       *CompetitorData   `json:"your_brand"`
	Competitors     []*CompetitorData `json:"competitors"`
	Rank            int               `json:"rank"`
	TotalCompetitors int              `json:"total_competitors"`
	Gaps            []*ContentGap     `json:"gaps"`
	Opportunities   []*Opportunity    `json:"opportunities"`
	GeneratedAt     time.Time         `json:"generated_at"`
}

// ContentGap 内容差距
type ContentGap struct {
	Topic       string  `json:"topic"`
	YourScore   float32 `json:"your_score"`
	BestScore   float32 `json:"best_score"`
	GapSize     float32 `json:"gap_size"`
	Priority    string  `json:"priority"`
	Suggestion  string  `json:"suggestion"`
}

// Opportunity 机会点
type Opportunity struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Effort      string  `json:"effort"`
}

// NewCompetitorAnalysisService 创建竞品分析服务
func NewCompetitorAnalysisService() *CompetitorAnalysisService {
	return &CompetitorAnalysisService{}
}

// AnalyzeCompetitor 分析单个竞品
func (s *CompetitorAnalysisService) AnalyzeCompetitor(ctx context.Context, competitor *CompetitorData) (*CompetitorAnalysisResult, error) {
	result := &CompetitorAnalysisResult{
		CompetitorName: competitor.Name,
		AnalyzedAt:     time.Now(),
	}

	// 计算可见性等级
	result.VisibilityLevel = getVisibilityLevel(competitor.VisibilityScore)

	// 分析引用表现
	result.CitationPerformance = analyzeCitationPerformance(competitor)

	// 分析内容策略
	result.ContentStrategy = analyzeContentStrategy(competitor)

	// 生成SWOT分析
	result.SWOT = generateSWOT(competitor)

	return result, nil
}

// CompetitorAnalysisResult 竞品分析结果
type CompetitorAnalysisResult struct {
	CompetitorName      string              `json:"competitor_name"`
	VisibilityLevel     string              `json:"visibility_level"`
	CitationPerformance *CitationPerformance `json:"citation_performance"`
	ContentStrategy     *ContentStrategy     `json:"content_strategy"`
	SWOT                *SWOTAnalysis        `json:"swot"`
	AnalyzedAt          time.Time            `json:"analyzed_at"`
}

// CitationPerformance 引用表现
type CitationPerformance struct {
	TotalCitations int32   `json:"total_citations"`
	AvgPosition    float32 `json:"avg_position"`
	TopModels      []string `json:"top_models"`
	Trend          string  `json:"trend"` // rising, stable, declining
}

// ContentStrategy 内容策略
type ContentStrategy struct {
	ContentTypes    []string `json:"content_types"`
	UpdateFrequency string   `json:"update_frequency"`
	TopTopics       []string `json:"top_topics"`
	Strengths       []string `json:"strengths"`
	Weaknesses      []string `json:"weaknesses"`
}

// SWOTAnalysis SWOT分析
type SWOTAnalysis struct {
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	Opportunities []string `json:"opportunities"`
	Threats     []string `json:"threats"`
}

func getVisibilityLevel(score float32) string {
	if score >= 80 {
		return "excellent"
	} else if score >= 60 {
		return "good"
	} else if score >= 40 {
		return "moderate"
	} else {
		return "low"
	}
}

func analyzeCitationPerformance(data *CompetitorData) *CitationPerformance {
	perf := &CitationPerformance{
		TotalCitations: data.CitationCount,
		AvgPosition:    data.AvgPosition,
	}

	if data.CitationCount > 100 {
		perf.Trend = "rising"
	} else if data.CitationCount > 50 {
		perf.Trend = "stable"
	} else {
		perf.Trend = "declining"
	}

	return perf
}

func analyzeContentStrategy(data *CompetitorData) *ContentStrategy {
	return &ContentStrategy{
		ContentTypes:    []string{"article", "faq"},
		UpdateFrequency: "weekly",
		TopTopics:       data.TopQueries,
		Strengths:       []string{"内容更新频繁", "引用率高"},
		Weaknesses:      []string{"结构化不足"},
	}
}

func generateSWOT(data *CompetitorData) *SWOTAnalysis {
	swot := &SWOTAnalysis{
		Strengths:     []string{},
		Weaknesses:    []string{},
		Opportunities: []string{},
		Threats:       []string{},
	}

	if data.VisibilityScore >= 70 {
		swot.Strengths = append(swot.Strengths, "AI可见性高")
	} else {
		swot.Weaknesses = append(swot.Weaknesses, "AI可见性有待提升")
	}

	if data.AvgPosition <= 3 {
		swot.Strengths = append(swot.Strengths, "引用位置靠前")
	} else {
		swot.Opportunities = append(swot.Opportunities, "优化内容结构以提升引用位置")
	}

	if data.CitationCount > 50 {
		swot.Threats = append(swot.Threats, "竞品引用量大，需关注")
	}

	return swot
}

// GenerateComparisonReport 生成对比报告
func (s *CompetitorAnalysisService) GenerateComparisonReport(ctx context.Context, yourBrand *CompetitorData, competitors []*CompetitorData) (*ComparisonReport, error) {
	report := &ComparisonReport{
		YourBrand:        yourBrand,
		Competitors:      competitors,
		TotalCompetitors: len(competitors) + 1,
		GeneratedAt:      time.Now(),
	}

	// 计算排名
	rank := 1
	for _, c := range competitors {
		if c.VisibilityScore > yourBrand.VisibilityScore {
			rank++
		}
	}
	report.Rank = rank

	// 分析内容差距
	report.Gaps = s.analyzeContentGaps(yourBrand, competitors)

	// 生成机会点
	report.Opportunities = s.identifyOpportunities(yourBrand, competitors)

	return report, nil
}

func (s *CompetitorAnalysisService) analyzeContentGaps(yourBrand *CompetitorData, competitors []*CompetitorData) []*ContentGap {
	gaps := make([]*ContentGap, 0)

	// 找出竞品有但你没有的热门话题
	competitorTopics := make(map[string]float32)
	for _, c := range competitors {
		for _, topic := range c.TopQueries {
			if score, exists := competitorTopics[topic]; !exists || c.VisibilityScore > score {
				competitorTopics[topic] = c.VisibilityScore
			}
		}
	}

	yourTopics := make(map[string]bool)
	for _, topic := range yourBrand.TopQueries {
		yourTopics[topic] = true
	}

	for topic, bestScore := range competitorTopics {
		if !yourTopics[topic] {
			gap := &ContentGap{
				Topic:      topic,
				YourScore:  0,
				BestScore:  bestScore,
				GapSize:    bestScore,
				Priority:   "high",
				Suggestion: fmt.Sprintf("建议创建关于「%s」的高质量内容", topic),
			}
			gaps = append(gaps, gap)
		}
	}

	return gaps
}

func (s *CompetitorAnalysisService) identifyOpportunities(yourBrand *CompetitorData, competitors []*CompetitorData) []*Opportunity {
	opportunities := make([]*Opportunity, 0)

	// 分析位置机会
	if yourBrand.AvgPosition > 3 {
		opportunities = append(opportunities, &Opportunity{
			Type:        "position",
			Description: "当前引用位置偏后，可通过优化内容结构提升位置",
			Impact:      "high",
			Effort:      "medium",
		})
	}

	// 分析引用量机会
	maxCitations := yourBrand.CitationCount
	for _, c := range competitors {
		if c.CitationCount > maxCitations {
			maxCitations = c.CitationCount
		}
	}
	if maxCitations > yourBrand.CitationCount*2 {
		opportunities = append(opportunities, &Opportunity{
			Type:        "citation",
			Description: "竞品引用量远超你，需增加内容覆盖面和权威性",
			Impact:      "high",
			Effort:      "high",
		})
	}

	return opportunities
}

// ==================== 自动化优化建议 ====================

// SuggestionEngine 优化建议引擎
type SuggestionEngine struct{}

// NewSuggestionEngine 创建建议引擎
func NewSuggestionEngine() *SuggestionEngine {
	return &SuggestionEngine{}
}

// SuggestionContext 建议上下文
type SuggestionContext struct {
	ContentID       int64   `json:"content_id"`
	Title           string  `json:"title"`
	BodyLength      int     `json:"body_length"`
	HasHeadings     bool    `json:"has_headings"`
	HasSchemaMarkup bool    `json:"has_schema_markup"`
	HasAIGCLabel    bool    `json:"has_aigc_label"`
	CitationCount   int32   `json:"citation_count"`
	AvgPosition     float32 `json:"avg_position"`
	CitationRate    float32 `json:"citation_rate"`
	VisibilityScore float32 `json:"visibility_score"`
	Keywords        []string `json:"keywords"`
}

// GenerateActionableSuggestions 生成可执行的优化建议
func (e *SuggestionEngine) GenerateActionableSuggestions(ctx context.Context, sctx *SuggestionContext) []*ActionableSuggestion {
	suggestions := make([]*ActionableSuggestion, 0)

	// 1. 内容结构建议
	if !sctx.HasHeadings {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "struct_headings",
			Category:    "content_structure",
			Title:       "添加标题层级结构",
			Description: "使用 ## / ### 标记创建清晰的标题层级，帮助AI理解内容结构",
			Priority:    "high",
			Impact:      "high",
			Effort:      "low",
			Steps: []string{
				"分析内容逻辑，确定主要章节",
				"使用 ## 标记主要章节标题",
				"使用 ### 标记子章节标题",
				"确保标题包含核心关键词",
			},
			Metric: "AI引用位置预计提升20-30%",
		})
	}

	// 2. Schema Markup建议
	if !sctx.HasSchemaMarkup {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "struct_schema",
			Category:    "content_structure",
			Title:       "添加Schema Markup结构化数据",
			Description: "JSON-LD格式的结构化数据帮助AI快速理解内容语义",
			Priority:    "high",
			Impact:      "high",
			Effort:      "medium",
			Steps: []string{
				"选择合适的Schema类型（Article/FAQ/Product）",
				"生成JSON-LD代码",
				"添加到页面<head>标签中",
				"使用Google结构化数据测试工具验证",
			},
			Metric: "AI收录速度预计提升50%",
		})
	}

	// 3. 引用优化建议
	if sctx.CitationCount < 5 {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "citation_authority",
			Category:    "authority",
			Title:       "增加权威引用",
			Description: "添加权威信源引用，提升内容可信度和AI引用意愿",
			Priority:    "high",
			Impact:      "high",
			Effort:      "medium",
			Steps: []string{
				"识别内容中的关键声明和数据",
				"查找权威来源（学术论文、官方报告、百科）",
				"添加引用链接和出处说明",
				"确保引用格式规范",
			},
			Metric: "AI引用率预计提升30-50%",
		})
	}

	// 4. 位置优化建议
	if sctx.AvgPosition > 5 {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "citation_position",
			Category:    "citation",
			Title:       "优化AI引用位置",
			Description: "当前内容在AI回答中位置偏后，需优化以获得更靠前的位置",
			Priority:    "medium",
			Impact:      "high",
			Effort:      "medium",
			Steps: []string{
				"分析当前排名靠前的内容特征",
				"优化内容开头，直接回答核心问题",
				"增加数据支撑和具体案例",
				"使用列表和要点提炼关键信息",
			},
			Metric: "引用位置预计从第5位提升到前3位",
		})
	}

	// 5. AIGC标识建议
	if !sctx.HasAIGCLabel {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "compliance_aigc",
			Category:    "compliance",
			Title:       "添加AIGC标识",
			Description: "添加AI生成内容标识，符合监管要求并避免被AI模型降权",
			Priority:    "medium",
			Impact:      "medium",
			Effort:      "low",
			Steps: []string{
				"在内容末尾添加AIGC声明",
				"说明内容的AI辅助程度",
				"添加人工审核说明",
			},
			Metric: "避免被AI模型降权风险",
		})
	}

	// 6. 内容长度建议
	if sctx.BodyLength < 1000 {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "content_length",
			Category:    "content",
			Title:       "增加内容深度",
			Description: "当前内容较短，建议扩展到1000字以上以提升权威性",
			Priority:    "medium",
			Impact:      "medium",
			Effort:      "medium",
			Steps: []string{
				"添加更多细节和案例",
				"增加数据支撑",
				"添加常见问题解答",
				"补充相关背景信息",
			},
			Metric: "内容深度提升后AI引用概率增加40%",
		})
	}

	// 7. 关键词优化建议
	if len(sctx.Keywords) < 3 {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "content_keywords",
			Category:    "seo",
			Title:       "优化关键词覆盖",
			Description: "增加核心关键词在内容中的出现频率",
			Priority:    "low",
			Impact:      "medium",
			Effort:      "low",
			Steps: []string{
				"确定3-5个核心关键词",
				"在标题、开头、结尾自然使用关键词",
				"在子标题中包含关键词变体",
				"确保关键词密度在1-3%之间",
			},
			Metric: "关键词匹配度提升后引用率增加20%",
		})
	}

	// 8. 可见性提升建议
	if sctx.VisibilityScore < 60 {
		suggestions = append(suggestions, &ActionableSuggestion{
			ID:          "visibility_boost",
			Category:    "strategy",
			Title:       "提升AI可见性综合策略",
			Description: "当前AI可见性较低，建议采用综合优化策略",
			Priority:    "high",
			Impact:      "high",
			Effort:      "high",
			Steps: []string{
				"分析目标AI模型的内容偏好",
				"优化内容结构和格式",
				"增加权威引用和数据支撑",
				"定期更新内容保持时效性",
				"在多个渠道发布建立信源矩阵",
			},
			Metric: "综合优化后可见性预计提升50-80%",
		})
	}

	return suggestions
}

// ActionableSuggestion 可执行建议
type ActionableSuggestion struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Impact      string   `json:"impact"`
	Effort      string   `json:"effort"`
	Steps       []string `json:"steps"`
	Metric      string   `json:"metric"`
}

// FormatSuggestionReport 格式化建议报告
func FormatSuggestionReport(suggestions []*ActionableSuggestion) string {
	if len(suggestions) == 0 {
		return "当前内容优化良好，暂无具体建议。"
	}

	report := fmt.Sprintf("## 优化建议报告\n\n共 %d 条建议：\n\n", len(suggestions))

	for i, s := range suggestions {
		report += fmt.Sprintf("### %d. %s\n", i+1, s.Title)
		report += fmt.Sprintf("- **分类**: %s\n", s.Category)
		report += fmt.Sprintf("- **优先级**: %s\n", s.Priority)
		report += fmt.Sprintf("- **影响**: %s | **工作量**: %s\n", s.Impact, s.Effort)
		report += fmt.Sprintf("- **预期效果**: %s\n", s.Metric)
		report += fmt.Sprintf("- **描述**: %s\n", s.Description)
		report += "- **执行步骤**:\n"
		for j, step := range s.Steps {
			report += fmt.Sprintf("  %d. %s\n", j+1, step)
		}
		report += "\n"
	}

	return report
}

// SuggestionToModel 将建议转换为数据模型
func SuggestionToModel(suggestion *ActionableSuggestion, contentID int64) map[string]interface{} {
	data, _ := json.Marshal(suggestion)
	return map[string]interface{}{
		"content_id":      contentID,
		"suggestion_type": suggestion.Category,
		"suggestion_data": string(data),
		"priority":        priorityToInt(suggestion.Priority),
		"status":          0,
	}
}

func priorityToInt(priority string) int32 {
	switch priority {
	case "high":
		return 2
	case "medium":
		return 1
	default:
		return 0
	}
}

// BatchGenerateSuggestions 批量生成建议
func (e *SuggestionEngine) BatchGenerateSuggestions(ctx context.Context, contexts []*SuggestionContext) map[int64][]*ActionableSuggestion {
	result := make(map[int64][]*ActionableSuggestion)

	for _, sctx := range contexts {
		suggestions := e.GenerateActionableSuggestions(ctx, sctx)
		if len(suggestions) > 0 {
			result[sctx.ContentID] = suggestions
		}
	}

	return result
}

// GetTopSuggestions 获取优先级最高的建议
func (e *SuggestionEngine) GetTopSuggestions(ctx context.Context, sctx *SuggestionContext, limit int) []*ActionableSuggestion {
	all := e.GenerateActionableSuggestions(ctx, sctx)

	// 按优先级排序
	high := make([]*ActionableSuggestion, 0)
	medium := make([]*ActionableSuggestion, 0)
	low := make([]*ActionableSuggestion, 0)

	for _, s := range all {
		switch s.Priority {
		case "high":
			high = append(high, s)
		case "medium":
			medium = append(medium, s)
		default:
			low = append(low, s)
		}
	}

	result := make([]*ActionableSuggestion, 0, limit)
	result = append(result, high...)
	result = append(result, medium...)
	result = append(result, low...)

	if len(result) > limit {
		result = result[:limit]
	}

	return result
}

// CalculateOptimizationScore 计算优化分数
func (e *SuggestionEngine) CalculateOptimizationScore(sctx *SuggestionContext) *OptimizationScore {
	score := &OptimizationScore{
		ContentID: sctx.ContentID,
	}

	// 结构分数
	structScore := float32(50)
	if sctx.HasHeadings {
		structScore += 25
	}
	if sctx.HasSchemaMarkup {
		structScore += 25
	}
	score.StructureScore = math.Min(float64(structScore), 100)

	// 权威分数
	authScore := float32(30)
	if sctx.CitationCount >= 5 {
		authScore += 35
	} else if sctx.CitationCount >= 2 {
		authScore += 20
	}
	if sctx.AvgPosition <= 3 {
		authScore += 35
	} else if sctx.AvgPosition <= 5 {
		authScore += 20
	}
	score.AuthorityScore = math.Min(float64(authScore), 100)

	// 合规分数
	complianceScore := float32(60)
	if sctx.HasAIGCLabel {
		complianceScore += 40
	}
	score.ComplianceScore = math.Min(float64(complianceScore), 100)

	// 综合分数
	score.OverallScore = (score.StructureScore + score.AuthorityScore + score.ComplianceScore) / 3

	return score
}

// OptimizationScore 优化分数
type OptimizationScore struct {
	ContentID       int64   `json:"content_id"`
	StructureScore  float64 `json:"structure_score"`
	AuthorityScore  float64 `json:"authority_score"`
	ComplianceScore float64 `json:"compliance_score"`
	OverallScore    float64 `json:"overall_score"`
}
