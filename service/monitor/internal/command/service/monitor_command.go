package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/service/monitor/internal/domain/model"
	"opengeo/service/monitor/internal/adapter"
)

// MonitorCommandService 监测命令服务（写侧）
type MonitorCommandService struct {
	citationRepo *adapter.CitationRepository
	scoreRepo    *adapter.ScoreRepository
	metricsRepo  *adapter.MetricsRepository
}

// NewMonitorCommandService 创建监测命令服务
func NewMonitorCommandService(
	citationRepo *adapter.CitationRepository,
	scoreRepo *adapter.ScoreRepository,
	metricsRepo *adapter.MetricsRepository,
) *MonitorCommandService {
	return &MonitorCommandService{
		citationRepo: citationRepo,
		scoreRepo:    scoreRepo,
		metricsRepo:  metricsRepo,
	}
}

// RecordOptimization 记录优化结果
func (s *MonitorCommandService) RecordOptimization(ctx context.Context, event *model.ContentOptimizedEvent) error {
	// 创建优化记录
	record := &model.OptimizationRecord{
		ContentID:        event.ContentID,
		UserID:           event.UserID,
		OptimizationType: event.OptimizationType,
		Score:            event.Score,
		OptimizedAt:      event.OccurredAt,
		CreatedAt:        time.Now(),
	}

	if err := s.metricsRepo.CreateOptimizationRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to create optimization record: %w", err)
	}

	return nil
}

// RecordPublishSuccess 记录发布成功
func (s *MonitorCommandService) RecordPublishSuccess(ctx context.Context, event *model.PublishSuccessEvent) error {
	// 创建发布记录
	record := &model.PublishRecord{
		TaskID:      event.TaskID,
		UserID:      event.UserID,
		ContentID:   event.ContentID,
		ChannelID:   event.ChannelID,
		ExternalID:  event.ExternalID,
		Status:      "success",
		PublishedAt: event.PublishedAt,
		CreatedAt:   time.Now(),
	}

	if err := s.metricsRepo.CreatePublishRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to create publish record: %w", err)
	}

	return nil
}

// RecordAICitation 记录AI引用
func (s *MonitorCommandService) RecordAICitation(ctx context.Context, event *model.AICitationFoundEvent) error {
	// 创建AI引用记录
	citation := &model.AICitation{
		ContentID:        event.ContentID,
		AIModel:          event.AIModel,
		QueryText:        event.QueryText,
		IsCited:          event.IsCited,
		CitationPosition: event.CitationPosition,
		CitationText:     event.CitationText,
		Sentiment:        event.Sentiment,
		TrackedAt:        event.OccurredAt,
		CreatedAt:        time.Now(),
	}

	if err := s.citationRepo.Create(ctx, citation); err != nil {
		return fmt.Errorf("failed to create AI citation: %w", err)
	}

	return nil
}

// UpdateSourceScore 更新信源评分
func (s *MonitorCommandService) UpdateSourceScore(ctx context.Context, channelID, accountID int64, score float32, dimensions string) error {
	sourceScore := &model.SourceScore{
		ChannelID:       channelID,
		AccountID:       accountID,
		Score:           score,
		ScoreDimensions: dimensions,
		UpdatedAt:       time.Now(),
	}

	if err := s.scoreRepo.Upsert(ctx, sourceScore); err != nil {
		return fmt.Errorf("failed to update source score: %w", err)
	}

	return nil
}