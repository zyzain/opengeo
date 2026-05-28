package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"opengeo/service/monitor/internal/dal"
	"opengeo/service/monitor/internal/domain/model"
)

// MonitorService 监测服务
type MonitorService struct {
	citationRepo    *dal.CitationRepository
	scoreRepo       *dal.ScoreRepository
	competitorRepo  *dal.CompetitorRepository
	roiRepo         *dal.ROIRepository
	suggestionRepo  *dal.SuggestionRepository
}

// NewMonitorService 创建监测服务
func NewMonitorService(
	citationRepo *dal.CitationRepository,
	scoreRepo *dal.ScoreRepository,
	competitorRepo *dal.CompetitorRepository,
	roiRepo *dal.ROIRepository,
	suggestionRepo *dal.SuggestionRepository,
) *MonitorService {
	return &MonitorService{
		citationRepo:   citationRepo,
		scoreRepo:      scoreRepo,
		competitorRepo: competitorRepo,
		roiRepo:        roiRepo,
		suggestionRepo: suggestionRepo,
	}
}

// ==================== AI引用追踪 ====================

// TrackAICitation 追踪AI引用
func (s *MonitorService) TrackAICitation(ctx context.Context, contentID int64, aiModel, queryText string, isCited bool, position int32, citationText, sentiment string) (*model.AICitation, error) {
	citation := &model.AICitation{
		ContentID:        contentID,
		AIModel:          aiModel,
		QueryText:        queryText,
		IsCited:          isCited,
		CitationPosition: position,
		CitationText:     citationText,
		Sentiment:        sentiment,
		TrackedAt:        time.Now(),
		CreatedAt:        time.Now(),
	}

	if err := s.citationRepo.Create(ctx, citation); err != nil {
		return nil, fmt.Errorf("failed to track citation: %w", err)
	}

	return citation, nil
}

// GetAICitations 获取AI引用列表
func (s *MonitorService) GetAICitations(ctx context.Context, contentID int64, aiModel string, page, pageSize int) ([]*model.AICitation, int32, error) {
	citations, total, err := s.citationRepo.List(ctx, contentID, aiModel, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get citations: %w", err)
	}

	return citations, total, nil
}

// GetCitationStats 获取引用统计
func (s *MonitorService) GetCitationStats(ctx context.Context, contentID int64) (*dal.CitationStats, error) {
	stats, err := s.citationRepo.GetCitationStats(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get citation stats: %w", err)
	}

	return stats, nil
}

// ==================== 信源评分 ====================

// UpdateSourceScore 更新信源评分
func (s *MonitorService) UpdateSourceScore(ctx context.Context, channelID, accountID int64, score float32, dimensions *model.ScoreDimension) error {
	dimensionsJSON, err := json.Marshal(dimensions)
	if err != nil {
		return fmt.Errorf("failed to marshal dimensions: %w", err)
	}

	sourceScore := &model.SourceScore{
		ChannelID:       channelID,
		AccountID:       accountID,
		Score:           score,
		ScoreDimensions: string(dimensionsJSON),
		UpdatedAt:       time.Now(),
	}

	if err := s.scoreRepo.Upsert(ctx, sourceScore); err != nil {
		return fmt.Errorf("failed to update source score: %w", err)
	}

	return nil
}

// GetSourceScores 获取信源评分
func (s *MonitorService) GetSourceScores(ctx context.Context, channelID, accountID int64) ([]*model.SourceScore, error) {
	scores, err := s.scoreRepo.List(ctx, channelID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source scores: %w", err)
	}

	return scores, nil
}

// ==================== 竞品监测 ====================

// CreateCompetitorMonitor 创建竞品监测
func (s *MonitorService) CreateCompetitorMonitor(ctx context.Context, userID int64, competitorName, competitorDomain, monitorConfig string) (*model.CompetitorMonitor, error) {
	monitor := &model.CompetitorMonitor{
		UserID:           userID,
		CompetitorName:   competitorName,
		CompetitorDomain: competitorDomain,
		MonitorConfig:    monitorConfig,
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.competitorRepo.Create(ctx, monitor); err != nil {
		return nil, fmt.Errorf("failed to create competitor monitor: %w", err)
	}

	return monitor, nil
}

// GetCompetitorMonitors 获取竞品监测列表
func (s *MonitorService) GetCompetitorMonitors(ctx context.Context, userID int64, page, pageSize int) ([]*model.CompetitorMonitor, int32, error) {
	monitors, total, err := s.competitorRepo.List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get competitor monitors: %w", err)
	}

	return monitors, total, nil
}

// GetCompetitorAnalysis 获取竞品分析
func (s *MonitorService) GetCompetitorAnalysis(ctx context.Context, monitorID int64) (*model.CompetitorAnalysis, error) {
	// TODO: 实现竞品分析逻辑
	analysis := &model.CompetitorAnalysis{
		MonitorID:       monitorID,
		VisibilityScore: 75.5,
		TopQueries:      `["query1", "query2", "query3"]`,
		ContentGaps:     `["gap1", "gap2"]`,
		Recommendations: `["recommendation1", "recommendation2"]`,
		AnalyzedAt:      time.Now(),
		CreatedAt:       time.Now(),
	}

	return analysis, nil
}

// ==================== ROI分析 ====================

// TrackROIMetric 追踪ROI指标
func (s *MonitorService) TrackROIMetric(ctx context.Context, contentID, channelID int64, metricType string, metricValue float64, utmSource, utmMedium, utmCampaign string) (*model.ROIMetric, error) {
	metric := &model.ROIMetric{
		ContentID:   contentID,
		ChannelID:   channelID,
		MetricType:  metricType,
		MetricValue: metricValue,
		UTMSource:   utmSource,
		UTMMedium:   utmMedium,
		UTMCampaign: utmCampaign,
		RecordedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := s.roiRepo.Create(ctx, metric); err != nil {
		return nil, fmt.Errorf("failed to track ROI metric: %w", err)
	}

	return metric, nil
}

// GetROIMetrics 获取ROI指标
func (s *MonitorService) GetROIMetrics(ctx context.Context, contentID, channelID int64, startDate, endDate string) ([]*model.ROIMetric, error) {
	metrics, err := s.roiRepo.List(ctx, contentID, channelID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get ROI metrics: %w", err)
	}

	return metrics, nil
}

// GetROIStats 获取ROI统计
func (s *MonitorService) GetROIStats(ctx context.Context, contentID int64) (*dal.ROIStats, error) {
	stats, err := s.roiRepo.GetROIStats(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ROI stats: %w", err)
	}

	return stats, nil
}

// ==================== 优化建议 ====================

// GenerateSuggestions 生成优化建议
func (s *MonitorService) GenerateSuggestions(ctx context.Context, contentID int64) ([]*model.OptimizationSuggestion, error) {
	// 获取引用统计
	citationStats, err := s.citationRepo.GetCitationStats(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get citation stats: %w", err)
	}

	suggestions := make([]*model.OptimizationSuggestion, 0)

	// 基于引用统计生成建议
	if citationStats.TotalCitations < 5 {
		suggestions = append(suggestions, &model.OptimizationSuggestion{
			ContentID:      contentID,
			SuggestionType: "content",
			SuggestionData: `{"suggestion": "增加权威引用和数据支撑，提升AI引用率"}`,
			Priority:       2,
			Status:         0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}

	if citationStats.AvgPosition > 5 {
		suggestions = append(suggestions, &model.OptimizationSuggestion{
			ContentID:      contentID,
			SuggestionType: "structure",
			SuggestionData: `{"suggestion": "优化内容结构，提升在AI回答中的位置"}`,
			Priority:       1,
			Status:         0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}

	// 保存建议
	for _, suggestion := range suggestions {
		if err := s.suggestionRepo.Create(ctx, suggestion); err != nil {
			return nil, fmt.Errorf("failed to create suggestion: %w", err)
		}
	}

	return suggestions, nil
}

// GetSuggestions 获取优化建议
func (s *MonitorService) GetSuggestions(ctx context.Context, contentID int64, status *int32) ([]*model.OptimizationSuggestion, error) {
	suggestions, err := s.suggestionRepo.List(ctx, contentID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}

	return suggestions, nil
}

// ApplySuggestion 应用优化建议
func (s *MonitorService) ApplySuggestion(ctx context.Context, suggestionID int64) error {
	if err := s.suggestionRepo.UpdateStatus(ctx, suggestionID, 1); err != nil {
		return fmt.Errorf("failed to apply suggestion: %w", err)
	}

	return nil
}

// IgnoreSuggestion 忽略优化建议
func (s *MonitorService) IgnoreSuggestion(ctx context.Context, suggestionID int64) error {
	if err := s.suggestionRepo.UpdateStatus(ctx, suggestionID, 2); err != nil {
		return fmt.Errorf("failed to ignore suggestion: %w", err)
	}

	return nil
}