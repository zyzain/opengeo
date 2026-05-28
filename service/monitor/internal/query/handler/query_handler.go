package handler

import (
	"context"

	"opengeo/service/monitor/internal/query/service"
)

// MonitorQueryHandler 监测查询处理器（读侧）
type MonitorQueryHandler struct {
	queryService *service.MonitorQueryService
}

// NewMonitorQueryHandler 创建监测查询处理器
func NewMonitorQueryHandler(queryService *service.MonitorQueryService) *MonitorQueryHandler {
	return &MonitorQueryHandler{queryService: queryService}
}

// GetAICitations 获取AI引用
func (h *MonitorQueryHandler) GetAICitations(ctx context.Context, req *GetAICitationsRequest) (*GetAICitationsResponse, error) {
	result, err := h.queryService.GetAICitations(ctx, req.ContentID, req.AIModel, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &GetAICitationsResponse{
		Citations: result.Citations,
		Total:     result.Total,
	}, nil
}

// GetSourceScores 获取信源评分
func (h *MonitorQueryHandler) GetSourceScores(ctx context.Context, req *GetSourceScoresRequest) (*GetSourceScoresResponse, error) {
	result, err := h.queryService.GetSourceScores(ctx, req.ChannelID, req.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetSourceScoresResponse{
		Scores: result.Scores,
	}, nil
}

// GetCompetitorMonitors 获取竞品监测
func (h *MonitorQueryHandler) GetCompetitorMonitors(ctx context.Context, req *GetCompetitorMonitorsRequest) (*GetCompetitorMonitorsResponse, error) {
	result, err := h.queryService.GetCompetitorMonitors(ctx, req.UserID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &GetCompetitorMonitorsResponse{
		Monitors: result.Monitors,
		Total:    result.Total,
	}, nil
}

// GetROIMetrics 获取ROI指标
func (h *MonitorQueryHandler) GetROIMetrics(ctx context.Context, req *GetROIMetricsRequest) (*GetROIMetricsResponse, error) {
	result, err := h.queryService.GetROIMetrics(ctx, req.ContentID, req.ChannelID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return &GetROIMetricsResponse{
		Metrics: result.Metrics,
	}, nil
}

// GenerateSuggestions 生成优化建议
func (h *MonitorQueryHandler) GenerateSuggestions(ctx context.Context, req *GenerateSuggestionsRequest) (*GenerateSuggestionsResponse, error) {
	result, err := h.queryService.GenerateSuggestions(ctx, req.ContentID)
	if err != nil {
		return nil, err
	}

	return &GenerateSuggestionsResponse{
		Suggestions: result.Suggestions,
	}, nil
}

// 请求/响应模型
type GetAICitationsRequest struct {
	ContentID int64  `json:"content_id"`
	AIModel   string `json:"ai_model"`
	Page      int32  `json:"page"`
	PageSize  int32  `json:"page_size"`
}

type GetAICitationsResponse struct {
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

type GetSourceScoresRequest struct {
	ChannelID int64 `json:"channel_id"`
	AccountID int64 `json:"account_id"`
}

type GetSourceScoresResponse struct {
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

type GetCompetitorMonitorsRequest struct {
	UserID   int64 `json:"user_id"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type GetCompetitorMonitorsResponse struct {
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

type GetROIMetricsRequest struct {
	ContentID int64  `json:"content_id"`
	ChannelID int64  `json:"channel_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type GetROIMetricsResponse struct {
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

type GenerateSuggestionsRequest struct {
	ContentID int64 `json:"content_id"`
}

type GenerateSuggestionsResponse struct {
	Suggestions []SuggestionInfo `json:"suggestions"`
}

type SuggestionInfo struct {
	ID            int64  `json:"id"`
	ContentID     int64  `json:"content_id"`
	SuggestionType string `json:"suggestion_type"`
	SuggestionData string `json:"suggestion_data"`
	Priority      int32  `json:"priority"`
	Status        int32  `json:"status"`
	CreatedAt     string `json:"created_at"`
}