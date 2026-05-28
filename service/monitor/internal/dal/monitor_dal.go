package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"opengeo/service/monitor/internal/domain/model"
)

// CitationRepository AI引用仓储
type CitationRepository struct {
	db *gorm.DB
}

// NewCitationRepository 创建AI引用仓储
func NewCitationRepository(db *gorm.DB) *CitationRepository {
	return &CitationRepository{db: db}
}

// Create 创建引用
func (r *CitationRepository) Create(ctx context.Context, citation *model.AICitation) error {
	if err := r.db.WithContext(ctx).Create(citation).Error; err != nil {
		return fmt.Errorf("failed to create citation: %w", err)
	}
	return nil
}

// List 列出引用
func (r *CitationRepository) List(ctx context.Context, contentID int64, aiModel string, page, pageSize int) ([]*model.AICitation, int32, error) {
	var citations []*model.AICitation
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AICitation{})

	if contentID > 0 {
		query = query.Where("content_id = ?", contentID)
	}
	if aiModel != "" {
		query = query.Where("ai_model = ?", aiModel)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count citations: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("tracked_at DESC").Find(&citations).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list citations: %w", err)
	}

	return citations, int32(total), nil
}

// GetCitationStats 获取引用统计
func (r *CitationRepository) GetCitationStats(ctx context.Context, contentID int64) (*CitationStats, error) {
	var stats CitationStats

	// 总引用数
	err := r.db.WithContext(ctx).
		Model(&model.AICitation{}).
		Where("content_id = ? AND is_cited = ?", contentID, true).
		Count(&stats.TotalCitations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count total citations: %w", err)
	}

	// 按模型统计
	type ModelCount struct {
		AIModel string
		Count   int64
	}
	var modelCounts []ModelCount
	err = r.db.WithContext(ctx).
		Model(&model.AICitation{}).
		Select("ai_model, COUNT(*) as count").
		Where("content_id = ? AND is_cited = ?", contentID, true).
		Group("ai_model").
		Find(&modelCounts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count by model: %w", err)
	}

	stats.ByModel = make(map[string]int64)
	for _, mc := range modelCounts {
		stats.ByModel[mc.AIModel] = mc.Count
	}

	// 平均位置
	err = r.db.WithContext(ctx).
		Model(&model.AICitation{}).
		Select("AVG(citation_position)").
		Where("content_id = ? AND is_cited = ? AND citation_position > 0", contentID, true).
		Scan(&stats.AvgPosition).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate avg position: %w", err)
	}

	return &stats, nil
}

// CitationStats 引用统计
type CitationStats struct {
	TotalCitations int64            `json:"total_citations"`
	ByModel        map[string]int64 `json:"by_model"`
	AvgPosition    float32          `json:"avg_position"`
}

// ScoreRepository 信源评分仓储
type ScoreRepository struct {
	db *gorm.DB
}

// NewScoreRepository 创建信源评分仓储
func NewScoreRepository(db *gorm.DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

// Upsert 更新或插入评分
func (r *ScoreRepository) Upsert(ctx context.Context, score *model.SourceScore) error {
	if err := r.db.WithContext(ctx).
		Where("channel_id = ? AND account_id = ?", score.ChannelID, score.AccountID).
		Assign(map[string]interface{}{
			"score":            score.Score,
			"score_dimensions": score.ScoreDimensions,
			"updated_at":       time.Now(),
		}).
		FirstOrCreate(score).Error; err != nil {
		return fmt.Errorf("failed to upsert score: %w", err)
	}
	return nil
}

// Get 获取评分
func (r *ScoreRepository) Get(ctx context.Context, channelID, accountID int64) (*model.SourceScore, error) {
	var score model.SourceScore

	err := r.db.WithContext(ctx).
		Where("channel_id = ? AND account_id = ?", channelID, accountID).
		First(&score).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &model.SourceScore{
				ChannelID: channelID,
				AccountID: accountID,
				Score:     0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get score: %w", err)
	}

	return &score, nil
}

// List 列出评分
func (r *ScoreRepository) List(ctx context.Context, channelID, accountID int64) ([]*model.SourceScore, error) {
	var scores []*model.SourceScore

	query := r.db.WithContext(ctx).Model(&model.SourceScore{})

	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if accountID > 0 {
		query = query.Where("account_id = ?", accountID)
	}

	if err := query.Order("score DESC").Limit(100).Find(&scores).Error; err != nil {
		return nil, fmt.Errorf("failed to list scores: %w", err)
	}

	return scores, nil
}

// CompetitorRepository 竞品监测仓储
type CompetitorRepository struct {
	db *gorm.DB
}

// NewCompetitorRepository 创建竞品监测仓储
func NewCompetitorRepository(db *gorm.DB) *CompetitorRepository {
	return &CompetitorRepository{db: db}
}

// Create 创建监测
func (r *CompetitorRepository) Create(ctx context.Context, monitor *model.CompetitorMonitor) error {
	if err := r.db.WithContext(ctx).Create(monitor).Error; err != nil {
		return fmt.Errorf("failed to create competitor monitor: %w", err)
	}
	return nil
}

// GetByID 根据ID获取监测
func (r *CompetitorRepository) GetByID(ctx context.Context, id int64) (*model.CompetitorMonitor, error) {
	var monitor model.CompetitorMonitor
	if err := r.db.WithContext(ctx).First(&monitor, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("competitor monitor not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get competitor monitor: %w", err)
	}
	return &monitor, nil
}

// List 列出监测
func (r *CompetitorRepository) List(ctx context.Context, userID int64, page, pageSize int) ([]*model.CompetitorMonitor, int32, error) {
	var monitors []*model.CompetitorMonitor
	var total int64

	query := r.db.WithContext(ctx).Model(&model.CompetitorMonitor{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count monitors: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&monitors).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list monitors: %w", err)
	}

	return monitors, int32(total), nil
}

// UpdateLastCheckTime 更新最后检查时间
func (r *CompetitorRepository) UpdateLastCheckTime(ctx context.Context, monitorID int64) error {
	if err := r.db.WithContext(ctx).
		Model(&model.CompetitorMonitor{}).
		Where("id = ?", monitorID).
		Update("last_check_time", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to update last check time: %w", err)
	}
	return nil
}

// ROIRepository ROI指标仓储
type ROIRepository struct {
	db *gorm.DB
}

// NewROIRepository 创建ROI指标仓储
func NewROIRepository(db *gorm.DB) *ROIRepository {
	return &ROIRepository{db: db}
}

// Create 创建指标
func (r *ROIRepository) Create(ctx context.Context, metric *model.ROIMetric) error {
	if err := r.db.WithContext(ctx).Create(metric).Error; err != nil {
		return fmt.Errorf("failed to create ROI metric: %w", err)
	}
	return nil
}

// List 列出指标
func (r *ROIRepository) List(ctx context.Context, contentID, channelID int64, startDate, endDate string) ([]*model.ROIMetric, error) {
	var metrics []*model.ROIMetric

	query := r.db.WithContext(ctx).Model(&model.ROIMetric{})

	if contentID > 0 {
		query = query.Where("content_id = ?", contentID)
	}
	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if startDate != "" {
		query = query.Where("recorded_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("recorded_at <= ?", endDate)
	}

	if err := query.Order("recorded_at DESC").Limit(500).Find(&metrics).Error; err != nil {
		return nil, fmt.Errorf("failed to list ROI metrics: %w", err)
	}

	return metrics, nil
}

// GetROIStats 获取ROI统计
func (r *ROIRepository) GetROIStats(ctx context.Context, contentID int64) (*ROIStats, error) {
	var stats ROIStats

	// 按类型统计
	type TypeSum struct {
		MetricType string
		Total      float64
	}
	var typeSums []TypeSum
	err := r.db.WithContext(ctx).
		Model(&model.ROIMetric{}).
		Select("metric_type, SUM(metric_value) as total").
		Where("content_id = ?", contentID).
		Group("metric_type").
		Find(&typeSums).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate ROI stats: %w", err)
	}

	stats.ByType = make(map[string]float64)
	for _, ts := range typeSums {
		stats.ByType[ts.MetricType] = ts.Total
		stats.Total += ts.Total
	}

	return &stats, nil
}

// ROIStats ROI统计
type ROIStats struct {
	Total  float64            `json:"total"`
	ByType map[string]float64 `json:"by_type"`
}

// SuggestionRepository 优化建议仓储
type SuggestionRepository struct {
	db *gorm.DB
}

// NewSuggestionRepository 创建优化建议仓储
func NewSuggestionRepository(db *gorm.DB) *SuggestionRepository {
	return &SuggestionRepository{db: db}
}

// Create 创建建议
func (r *SuggestionRepository) Create(ctx context.Context, suggestion *model.OptimizationSuggestion) error {
	if err := r.db.WithContext(ctx).Create(suggestion).Error; err != nil {
		return fmt.Errorf("failed to create suggestion: %w", err)
	}
	return nil
}

// List 列出建议
func (r *SuggestionRepository) List(ctx context.Context, contentID int64, status *int32) ([]*model.OptimizationSuggestion, error) {
	var suggestions []*model.OptimizationSuggestion

	query := r.db.WithContext(ctx).Where("content_id = ?", contentID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("priority DESC, created_at DESC").Limit(100).Find(&suggestions).Error; err != nil {
		return nil, fmt.Errorf("failed to list suggestions: %w", err)
	}

	return suggestions, nil
}

// UpdateStatus 更新建议状态
func (r *SuggestionRepository) UpdateStatus(ctx context.Context, suggestionID int64, status int32) error {
	if err := r.db.WithContext(ctx).
		Model(&model.OptimizationSuggestion{}).
		Where("id = ?", suggestionID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update suggestion status: %w", err)
	}
	return nil
}