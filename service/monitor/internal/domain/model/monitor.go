package model

import "time"

// AICitation AI引用
type AICitation struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ContentID        int64     `json:"content_id" gorm:"index;not null"`
	AIModel          string    `json:"ai_model" gorm:"size:64;not null;index:idx_citation_stats"`
	QueryText        string    `json:"query_text" gorm:"size:512;not null"`
	IsCited          bool      `json:"is_cited" gorm:"index:idx_citation_stats"`
	CitationPosition int32     `json:"citation_position"`
	CitationText     string    `json:"citation_text" gorm:"type:text"`
	Sentiment        string    `json:"sentiment" gorm:"size:32"` // positive, neutral, negative
	TrackedAt        time.Time `json:"tracked_at"`
	CreatedAt        time.Time `json:"created_at"`
}

// SourceScore 信源评分
type SourceScore struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ChannelID       int64     `json:"channel_id" gorm:"not null;index:idx_source_score_lookup"`
	AccountID       int64     `json:"account_id" gorm:"not null;index:idx_source_score_lookup"`
	Score           float32   `json:"score"`
	ScoreDimensions string    `json:"score_dimensions" gorm:"type:text"` // JSON格式的各维度评分
	UpdatedAt       time.Time `json:"updated_at"`
}

// ScoreDimension 评分维度
type ScoreDimension struct {
	RecencySpeed    float32 `json:"recency_speed"`    // 收录速度
	RankingStability float32 `json:"ranking_stability"` // 排名稳定性
	CitationFrequency float32 `json:"citation_frequency"` // 引用频次
	AuthorityScore  float32 `json:"authority_score"`   // 权威性评分
	ContentQuality  float32 `json:"content_quality"`   // 内容质量
}

// CompetitorMonitor 竞品监测
type CompetitorMonitor struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int64      `json:"user_id" gorm:"index;not null"`
	CompetitorName  string     `json:"competitor_name" gorm:"size:128;not null"`
	CompetitorDomain string    `json:"competitor_domain" gorm:"size:256"`
	MonitorConfig   string     `json:"monitor_config" gorm:"type:text"` // JSON格式的监测配置
	LastCheckTime   *time.Time `json:"last_check_time"`
	IsActive        bool       `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CompetitorAnalysis 竞品分析结果
type CompetitorAnalysis struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MonitorID       int64     `json:"monitor_id" gorm:"index;not null"`
	VisibilityScore float32   `json:"visibility_score"` // 可见性评分
	TopQueries      string    `json:"top_queries" gorm:"type:text"` // JSON格式的热门查询
	ContentGaps     string    `json:"content_gaps" gorm:"type:text"` // JSON格式的内容差距
	Recommendations string    `json:"recommendations" gorm:"type:text"` // 建议
	AnalyzedAt      time.Time `json:"analyzed_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// ROIMetric ROI指标
type ROIMetric struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ContentID   int64     `json:"content_id" gorm:"not null;index:idx_roi_lookup"`
	ChannelID   int64     `json:"channel_id" gorm:"not null;index:idx_roi_lookup"`
	MetricType  string    `json:"metric_type" gorm:"size:32;not null"` // inquiry, visit, consult
	MetricValue float64   `json:"metric_value"`
	UTMSource   string    `json:"utm_source" gorm:"size:64"`
	UTMMedium   string    `json:"utm_medium" gorm:"size:64"`
	UTMCampaign string    `json:"utm_campaign" gorm:"size:64"`
	RecordedAt  time.Time `json:"recorded_at" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
}

// OptimizationSuggestion 优化建议
type OptimizationSuggestion struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ContentID      int64     `json:"content_id" gorm:"not null;index:idx_suggestion_lookup"`
	SuggestionType string    `json:"suggestion_type" gorm:"size:32"` // content, structure, authority
	SuggestionData string    `json:"suggestion_data" gorm:"type:text"` // JSON格式的建议内容
	Priority       int32     `json:"priority" gorm:"default:0"` // 0:低 1:中 2:高
	Status         int32     `json:"status" gorm:"default:0;index:idx_suggestion_lookup"` // 0:待处理 1:已应用 2:已忽略
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CitationTrend 引用趋势
type CitationTrend struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ContentID     int64     `json:"content_id" gorm:"index;not null"`
	AIModel       string    `json:"ai_model" gorm:"size:64;not null"`
	CitationCount int32     `json:"citation_count"`
	CitationRate  float32   `json:"citation_rate"` // 引用率
	TrendDate     time.Time `json:"trend_date" gorm:"index;not null"`
	CreatedAt     time.Time `json:"created_at"`
}