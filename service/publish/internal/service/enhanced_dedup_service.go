package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"opengeo/pkg/similarity"
	"opengeo/service/publish/internal/dal"
	"opengeo/service/publish/internal/domain/model"
)

// EnhancedDeduplicationService 增强版内容去重服务
type EnhancedDeduplicationService struct {
	fpRepo       *dal.ContentFingerprintRepository
	synonymRepo  *dal.SynonymDictRepository
	historyRepo  *dal.DedupHistoryRepository
	similarity   *similarity.CombinedSimilarity
	aiRewriter   *AIRewriter
	synonymMap   map[string][]string
}

// NewEnhancedDeduplicationService 创建增强版去重服务
func NewEnhancedDeduplicationService(
	fpRepo *dal.ContentFingerprintRepository,
	synonymRepo *dal.SynonymDictRepository,
	historyRepo *dal.DedupHistoryRepository,
	aiRewriter *AIRewriter,
) *EnhancedDeduplicationService {
	return &EnhancedDeduplicationService{
		fpRepo:      fpRepo,
		synonymRepo: synonymRepo,
		historyRepo: historyRepo,
		similarity:  similarity.NewCombinedSimilarity(),
		aiRewriter:  aiRewriter,
	}
}

// InitSynonymMap 初始化同义词映射
func (s *EnhancedDeduplicationService) InitSynonymMap(ctx context.Context) error {
	synonymMap, err := s.synonymRepo.BuildSynonymMap(ctx)
	if err != nil {
		// 如果数据库加载失败，使用默认映射
		synonymMap = buildDefaultSynonymMap()
	}
	s.synonymMap = synonymMap
	return nil
}

// EnhancedDeduplicateRequest 增强版去重请求
type EnhancedDeduplicateRequest struct {
	UserID       int64    `json:"user_id"`
	ContentID    int64    `json:"content_id"`
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Tags         []string `json:"tags"`
	ContentType  string   `json:"content_type"`
	Strategy     string   `json:"strategy"`      // light, medium, heavy, ai
	TargetSim    float32  `json:"target_sim"`    // 目标相似度 (0-1)
	EnableAI     bool     `json:"enable_ai"`     // 是否启用AI改写
	EnableDBCheck bool   `json:"enable_db_check"` // 是否检查数据库历史内容
}

// EnhancedDeduplicateResponse 增强版去重响应
type EnhancedDeduplicateResponse struct {
	Title            string                    `json:"title"`
	Body             string                    `json:"body"`
	Tags             []string                  `json:"tags"`
	Similarity       float32                   `json:"similarity"`
	Changes          []string                  `json:"changes"`
	OriginalLen      int                       `json:"original_len"`
	DedupLen         int                       `json:"dedup_len"`
	DuplicatesFound  int                       `json:"duplicates_found"`
	DuplicateIDs     []int64                   `json:"duplicate_ids"`
	FingerprintData  *similarity.ContentFingerprintData `json:"fingerprint_data"`
	AITransformed    bool                      `json:"ai_transformed"`
	TransformDetails string                    `json:"transform_details"`
}

// EnhancedDeduplicate 增强版内容去重
func (s *EnhancedDeduplicationService) EnhancedDeduplicate(ctx context.Context, req *EnhancedDeduplicateRequest) (*EnhancedDeduplicateResponse, error) {
	// 初始化同义词映射
	if s.synonymMap == nil {
		s.InitSynonymMap(ctx)
	}
	
	// 设置默认值
	strategy := req.Strategy
	if strategy == "" {
		strategy = "medium"
	}

	targetSim := req.TargetSim
	if targetSim == 0 {
		targetSim = 0.7
	}

	changes := make([]string, 0)
	originalBody := req.Body
	duplicateIDs := make([]int64, 0)
	aiTransformed := false
	transformDetails := ""

	// 空内容直接返回
	if req.Title == "" && req.Body == "" {
		return &EnhancedDeduplicateResponse{
			Title:       req.Title,
			Body:        req.Body,
			Tags:        req.Tags,
			Similarity:  1.0,
			Changes:     changes,
			OriginalLen: 0,
			DedupLen:    0,
		}, nil
	}

	// 1. 检查数据库历史内容
	if req.EnableDBCheck && req.UserID > 0 {
		duplicates, err := s.findDuplicatesInDB(ctx, req.UserID, req.Title, req.Body, req.ContentType)
		if err == nil && len(duplicates) > 0 {
			for _, dup := range duplicates {
				duplicateIDs = append(duplicateIDs, dup.ContentID)
			}
			changes = append(changes, fmt.Sprintf("发现 %d 条相似历史内容", len(duplicates)))
		}
	}

	// 2. 计算原始内容指纹
	fingerprintData := s.similarity.ComputeFingerprint(originalBody)

	// 3. 根据策略处理
	newTitle := req.Title
	newBody := originalBody
	newTags := req.Tags

	// 检查是否需要AI改写
	if req.EnableAI && (strategy == "ai" || strategy == "heavy") {
		// 使用AI改写
		aiResult, err := s.aiRewriter.Rewrite(ctx, &AIRewriteRequest{
			Title:   newTitle,
			Body:    newBody,
			Style:   strategy,
			Purpose: "dedup",
		})
		if err == nil && aiResult.Success {
			newTitle = aiResult.Title
			newBody = aiResult.Body
			aiTransformed = true
			transformDetails = aiResult.Details
			changes = append(changes, "AI智能改写完成")
		}
	} else {
		// 使用传统方法
		// 处理标题
		newTitle = s.deduplicateTitle(newTitle, strategy)

		// 处理正文
		// 段落重排
		if strategy == "medium" || strategy == "heavy" {
			paragraphs := s.splitParagraphs(newBody)
			if len(paragraphs) > 2 {
				reordered := s.reorderParagraphs(paragraphs, strategy)
				if reordered != newBody {
					newBody = reordered
					changes = append(changes, "段落顺序已调整")
				}
			}
		}

		// 同义词替换
		replaced, replaceCount := s.synonymReplace(newBody, strategy)
		if replaceCount > 0 {
			newBody = replaced
			changes = append(changes, fmt.Sprintf("同义词替换 %d 处", replaceCount))
		}

		// 句式变换
		if strategy == "heavy" {
			transformed, transformCount := s.sentenceTransform(newBody)
			if transformCount > 0 {
				newBody = transformed
				changes = append(changes, fmt.Sprintf("句式变换 %d 处", transformCount))
			}
		}
	}

	// 处理标签
	newTags = s.deduplicateTags(newTags, strategy)

	// 计算相似度
	simResult := s.similarity.Compute(originalBody, newBody, float64(targetSim))
	similarityScore := float32(simResult.CombinedSimilarity)

	// 保存内容指纹
	if req.UserID > 0 {
		s.saveFingerprint(ctx, req.UserID, req.ContentID, req.Title, originalBody, req.ContentType)
	}

	// 保存去重历史
	if req.UserID > 0 {
		s.saveHistory(ctx, req.UserID, req.ContentID, originalBody, newBody, similarityScore, len(duplicateIDs), duplicateIDs, strategy, aiTransformed)
	}

	return &EnhancedDeduplicateResponse{
		Title:            newTitle,
		Body:             newBody,
		Tags:             newTags,
		Similarity:       similarityScore,
		Changes:          changes,
		OriginalLen:      utf8.RuneCountInString(originalBody),
		DedupLen:         utf8.RuneCountInString(newBody),
		DuplicatesFound:  len(duplicateIDs),
		DuplicateIDs:     duplicateIDs,
		FingerprintData:  fingerprintData,
		AITransformed:    aiTransformed,
		TransformDetails: transformDetails,
	}, nil
}

// findDuplicatesInDB 在数据库中查找重复内容
func (s *EnhancedDeduplicationService) findDuplicatesInDB(ctx context.Context, userID int64, title, body, contentType string) ([]*model.ContentFingerprint, error) {
	// 计算内容指纹
	titleHash := fmt.Sprintf("%d", s.similarity.ComputeFingerprint(title).SimHash)
	bodyHash := fmt.Sprintf("%d", s.similarity.ComputeFingerprint(body).SimHash)

	// 查找相似内容
	fps, err := s.fpRepo.FindSimilarByHash(ctx, userID, titleHash, bodyHash, 10)
	if err != nil {
		return nil, err
	}

	// 进一步验证相似度
	var duplicates []*model.ContentFingerprint
	for _, fp := range fps {
		// 使用余弦相似度进行精确验证
		sim := s.similarity.Compute(body, fp.BodyFingerprint, 0.7)
		if sim.CombinedSimilarity >= 0.7 {
			duplicates = append(duplicates, fp)
		}
	}

	return duplicates, nil
}

// saveFingerprint 保存内容指纹
func (s *EnhancedDeduplicationService) saveFingerprint(ctx context.Context, userID, contentID int64, title, body, contentType string) {
	fpData := s.similarity.ComputeFingerprint(body)
	titleFpData := s.similarity.ComputeFingerprint(title)

	keywordsJSON, _ := json.Marshal(fpData.Keywords)

	fp := &model.ContentFingerprint{
		UserID:           userID,
		ContentID:        contentID,
		TitleHash:        fmt.Sprintf("%d", titleFpData.SimHash),
		BodyHash:         fmt.Sprintf("%d", fpData.SimHash),
		TitleFingerprint: fmt.Sprintf("%d", titleFpData.SimHash),
		BodyFingerprint:  fmt.Sprintf("%d", fpData.SimHash),
		Keywords:         string(keywordsJSON),
		WordCount:        fpData.WordCount,
		ContentType:      contentType,
		CreatedAt:        time.Now(),
	}

	s.fpRepo.Create(ctx, fp)
}

// saveHistory 保存去重历史
func (s *EnhancedDeduplicationService) saveHistory(ctx context.Context, userID, contentID int64, original, deduped string, similarity float32, dupCount int, dupIDs []int64, strategy string, aiTransformed bool) {
	originalHash := fmt.Sprintf("%d", s.similarity.ComputeFingerprint(original).SimHash)
	dedupedHash := fmt.Sprintf("%d", s.similarity.ComputeFingerprint(deduped).SimHash)

	dupIDsJSON, _ := json.Marshal(dupIDs)

	history := &model.DedupHistory{
		UserID:         userID,
		ContentID:      contentID,
		OriginalHash:   originalHash,
		DedupedHash:    dedupedHash,
		Similarity:     similarity,
		DuplicateCount: dupCount,
		DuplicateIDs:   string(dupIDsJSON),
		Strategy:       strategy,
		AITransformed:  aiTransformed,
		CreatedAt:      time.Now(),
	}

	s.historyRepo.Create(ctx, history)
}

// ==================== 传统去重方法 ====================

func (s *EnhancedDeduplicationService) deduplicateTitle(title, strategy string) string {
	if title == "" {
		return title
	}

	replaced, count := s.synonymReplace(title, strategy)
	if count > 0 {
		return replaced
	}

	return title
}

func (s *EnhancedDeduplicationService) splitParagraphs(body string) []string {
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

func (s *EnhancedDeduplicationService) reorderParagraphs(paragraphs []string, strategy string) string {
	if len(paragraphs) <= 2 {
		return strings.Join(paragraphs, "\n\n")
	}

	result := make([]string, len(paragraphs))
	copy(result, paragraphs)

	// 简单的段落重排逻辑
	if len(result) > 3 {
		// 保留首尾，重排中间
		middle := result[1 : len(result)-1]
		// 简单反转
		for i, j := 0, len(middle)-1; i < j; i, j = i+1, j-1 {
			middle[i], middle[j] = middle[j], middle[i]
		}
	}

	return strings.Join(result, "\n\n")
}

func (s *EnhancedDeduplicationService) synonymReplace(text, strategy string) (string, int) {
	count := 0
	result := text

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
		if strings.Contains(result, word) && len(synonyms) > 0 {
			// 随机决定是否替换
			if float64(count)/float64(len(s.synonymMap)) < probability {
				synonym := synonyms[0] // 使用第一个同义词
				result = strings.Replace(result, word, synonym, 1)
				count++
			}
		}
	}

	return result, count
}

func (s *EnhancedDeduplicationService) sentenceTransform(text string) (string, int) {
	count := 0
	sentences := splitSentences(text)
	result := make([]string, len(sentences))

	transitions := []string{"此外，", "同时，", "另外，", "值得注意的是，", "具体来说，"}
	qualifiers := []string{"一般来说，", "通常情况下，", "从实践来看，"}

	for i, sent := range sentences {
		transformed := sent

		// 添加过渡词
		if i > 0 && i < len(sentences)-1 && count < 3 {
			if !hasTransition(sent) {
				transformed = transitions[count%len(transitions)] + transformed
				count++
			}
		}

		// 添加限定词
		if i == 0 && count < 5 {
			transformed = qualifiers[0] + transformed
			count++
		}

		result[i] = transformed
	}

	return strings.Join(result, ""), count
}

func (s *EnhancedDeduplicationService) deduplicateTags(tags []string, strategy string) []string {
	if len(tags) == 0 {
		return tags
	}

	result := make([]string, len(tags))
	copy(result, tags)

	// 替换部分标签的同义词
	for i, tag := range result {
		if synonyms, ok := s.synonymMap[tag]; ok && len(synonyms) > 0 {
			result[i] = synonyms[0]
		}
	}

	return result
}

// buildDefaultSynonymMap 构建默认同义词映射
func buildDefaultSynonymMap() map[string][]string {
	return map[string][]string{
		"提升": {"提高", "增强", "改善", "优化"},
		"提高": {"提升", "增强", "改善", "优化"},
		"增强": {"提升", "提高", "强化", "加强"},
		"优化": {"改进", "完善", "提升", "改善"},
		"实现": {"达成", "完成", "达到", "做到"},
		"提供": {"供给", "供应", "给予", "带来"},
		"使用": {"采用", "运用", "利用", "应用"},
		"重要": {"关键", "核心", "主要", "关键性"},
		"有效": {"高效", "有力", "显著", "明显"},
		"快速": {"迅速", "高效", "敏捷", "即时"},
		"简单": {"简便", "便捷", "易用", "轻松"},
		"功能": {"特性", "能力", "模块", "组件"},
		"性能": {"效率", "表现", "速度", "指标"},
		"系统": {"平台", "方案", "工具", "产品"},
		"数据": {"信息", "资料", "内容", "素材"},
	}
}
