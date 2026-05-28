package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/service/monitor/internal/query/repository"
)

// MonitorQueryService 监测查询服务（读侧）
type MonitorQueryService struct {
	queryRepo *repository.MonitorQueryRepository
}

// NewMonitorQueryService 创建监测查询服务
func NewMonitorQueryService(queryRepo *repository.MonitorQueryRepository) *MonitorQueryService {
	return &MonitorQueryService{queryRepo: queryRepo}
}

// GetAICitations 获取AI引用
func (s *MonitorQueryService) GetAICitations(ctx context.Context, contentID int64, aiModel string, page, pageSize int32) (*CitationQueryResult, error) {
	citations, total, err := s.queryRepo.GetAICitations(ctx, contentID, aiModel, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI citations: %w", err)
	}

	return &CitationQueryResult{
		Citations: citations,
		Total:     total,
	}, nil
}

// GetSourceScores 获取信源评分
func (s *MonitorQueryService) GetSourceScores(ctx context.Context, channelID, accountID int64) (*ScoreQueryResult, error) {
	scores, err := s.queryRepo.GetSourceScores(ctx, channelID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source scores: %w", err)
	}

	return &ScoreQueryResult{
		Scores: scores,
	}, nil
}

// GetCompetitorMonitors 获取竞品监测
func (s *MonitorQueryService) GetCompetitorMonitors(ctx context.Context, userID int64, page, pageSize int32) (*MonitorQueryResult, error) {
	monitors, total, err := s.queryRepo.GetCompetitorMonitors(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor monitors: %w", err)
	}

	return &MonitorQueryResult{
		Monitors: monitors,
		Total:    total,
	}, nil
}

// GetROIMetrics 获取ROI指标
func (s *MonitorQueryService) GetROIMetrics(ctx context.Context, contentID, channelID int64, startDate, endDate string) (*MetricQueryResult, error) {
	metrics, err := s.queryRepo.GetROIMetrics(ctx, contentID, channelID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get ROI metrics: %w", err)
	}

	return &MetricQueryResult{
		Metrics: metrics,
	}, nil
}

// GenerateSuggestions 生成优化建议
func (s *MonitorQueryService) GenerateSuggestions(ctx context.Context, contentID int64) (*SuggestionQueryResult, error) {
	suggestions, err := s.queryRepo.GenerateSuggestions(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate suggestions: %w", err)
	}

	return &SuggestionQueryResult{
		Suggestions: suggestions,
	}, nil
}

// 查询结果模型
type CitationQueryResult struct {
	Citations []CitationInfo `json:"citations"`
	Total     int32          `json:"total"`
}

type CitationInfo struct {
	ID              int64  `json:"id"`
	ContentID       int64  `json:"content_id"`
	AIModel         string `json:"ai_model"`
	QueryText       string `json:"query_text"`
	IsCited         bool   `json:"is_cited"`
	CitationPosition int32 `json:"citation_position"`
	CitationText    string `json:"citation_text"`
	Sentiment       string `json:"sentiment"`
	TrackedAt       string `json:"tracked_at"`
}

type ScoreQueryResult struct {
	Scores []ScoreInfo `json:"scores"`
}

type ScoreInfo struct {
	ID              int64   `json:"id"`
	ChannelID       int64   `json:"channel_id"`
	AccountID       int64   `json:"account_id"`
	Score           float32 `json:"score"`
	ScoreDimensions string  `json:"score_dimensions"`
	UpdatedAt       string  `json:"updated_at"`
}

type MonitorQueryResult struct {
	Monitors []MonitorInfo `json:"monitors"`
	Total    int32         `json:"total"`
}

type MonitorInfo struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	CompetitorName  string `json:"competitor_name"`
	CompetitorDomain string `json:"competitor_domain"`
	LastCheckTime   string `json:"last_check_time"`
	CreatedAt       string `json:"created_at"`
}

type MetricQueryResult struct {
	Metrics []MetricInfo `json:"metrics"`
}

type MetricInfo struct {
	ID          int64   `json:"id"`
	ContentID   int64   `json:"content_id"`
	ChannelID   int64   `json:"channel_id"`
	MetricType  string  `json:"metric_type"`
	MetricValue float32 `json:"metric_value"`
	RecordedAt  string  `json:"recorded_at"`
}

type SuggestionQueryResult struct {
	Suggestions []SuggestionInfo `json:"suggestions"`
}

type SuggestionInfo struct {
	ID             int64  `json:"id"`
	ContentID      int64  `json:"content_id"`
	SuggestionType string `json:"suggestion_type"`
	SuggestionData string `json:"suggestion_data"`
	Priority       int32  `json:"priority"`
	Status         int32  `json:"status"`
	CreatedAt      string `json:"created_at"`
}