package model

import (
	"fmt"
	"time"

	"opengeo/pkg/config"
)

// PublishTask 发布任务聚合根
type PublishTask struct {
	ID                int64         `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID            int64         `json:"user_id" gorm:"index;not null"`
	ContentID         int64         `json:"content_id" gorm:"index;not null"`
	ChannelID         int64         `json:"channel_id" gorm:"index;not null"`
	FallbackChannelID *int64        `json:"fallback_channel_id"`
	Status            PublishStatus `json:"status" gorm:"index;default:0"`
	ScheduledTime     *time.Time    `json:"scheduled_time" gorm:"index"`
	PublishedTime     *time.Time    `json:"published_time"`
	RetryCount        int32         `json:"retry_count" gorm:"default:0"`
	MaxRetries        int32         `json:"max_retries" gorm:"default:3"`
	RetryDelay        int32         `json:"retry_delay" gorm:"default:30"`
	ErrorMessage      string        `json:"error_message" gorm:"size:512"`
	ErrorHistory      string        `json:"error_history" gorm:"type:text"`
	Priority          int32         `json:"priority" gorm:"index;default:0"`
	IsManuallyReview  bool          `json:"is_manually_review" gorm:"default:false"`
	ReviewNote        string        `json:"review_note" gorm:"size:512"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// PublishStatus 发布状态
type PublishStatus int32

const (
	PublishStatusPending    PublishStatus = 0 // 待发布
	PublishStatusPublishing PublishStatus = 1 // 发布中
	PublishStatusSuccess    PublishStatus = 2 // 已发布
	PublishStatusFailed     PublishStatus = 3 // 失败
	PublishStatusCancelled  PublishStatus = 4 // 已取消
	PublishStatusRetrying   PublishStatus = 5 // 重试中
	PublishStatusFallback   PublishStatus = 6 // 降级中
	PublishStatusReview     PublishStatus = 7 // 人工审核
)

// MarkAsPublishing 标记为发布中
func (t *PublishTask) MarkAsPublishing() {
	t.Status = PublishStatusPublishing
	t.UpdatedAt = time.Now()
}

// MarkAsSuccess 标记为成功
func (t *PublishTask) MarkAsSuccess() {
	t.Status = PublishStatusSuccess
	now := time.Now()
	t.PublishedTime = &now
	t.UpdatedAt = now
}

// MarkAsFailed 标记为失败
func (t *PublishTask) MarkAsFailed(errMsg string) {
	t.Status = PublishStatusFailed
	t.ErrorMessage = errMsg
	t.UpdatedAt = time.Now()
}

// MarkAsRetrying 标记为重试中
func (t *PublishTask) MarkAsRetrying() {
	t.Status = PublishStatusRetrying
	t.UpdatedAt = time.Now()
}

// MarkAsFallback 标记为降级中
func (t *PublishTask) MarkAsFallback() {
	t.Status = PublishStatusFallback
	t.UpdatedAt = time.Now()
}

// MarkAsReview 标记为人工审核
func (t *PublishTask) MarkAsReview(note string) {
	t.Status = PublishStatusReview
	t.IsManuallyReview = true
	t.ReviewNote = note
	t.UpdatedAt = time.Now()
}

// CanRetry 是否可重试
func (t *PublishTask) CanRetry() bool {
	return t.RetryCount < t.MaxRetries && t.Status == PublishStatusFailed
}

// CanFallback 是否可降级
func (t *PublishTask) CanFallback() bool {
	return t.FallbackChannelID != nil && t.Status == PublishStatusFailed
}

// IncrementRetry 增加重试次数
func (t *PublishTask) IncrementRetry() {
	t.RetryCount++
	t.Status = PublishStatusPending
	t.UpdatedAt = time.Now()
}

// GetNextRetryDelay 获取下次重试延迟（指数退避）
func (t *PublishTask) GetNextRetryDelay() time.Duration {
	cfg := config.GetConfig()
	baseDelay := t.RetryDelay
	if baseDelay <= 0 {
		baseDelay = cfg.Retry.InitialDelay
	}
	delay := float64(baseDelay) * float64(uint(1)<<uint(t.RetryCount))
	if delay > float64(cfg.Retry.MaxDelay) {
		delay = float64(cfg.Retry.MaxDelay)
	}
	return time.Duration(delay) * time.Second
}

// AddErrorHistory 添加错误历史
func (t *PublishTask) AddErrorHistory(errMsg string) {
	entry := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), errMsg)
	if t.ErrorHistory == "" {
		t.ErrorHistory = entry
	} else {
		t.ErrorHistory = t.ErrorHistory + "\n" + entry
	}
}

// FallbackQueue 降级队列
type FallbackQueue struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID        int64     `json:"task_id" gorm:"index;not null"`
	UserID        int64     `json:"user_id" gorm:"index;not null"`
	ContentID     int64     `json:"content_id"`
	OriginalChannelID int64 `json:"original_channel_id"`
	FallbackChannelID int64 `json:"fallback_channel_id"`
	Reason        string    `json:"reason" gorm:"type:text"`
	Status        string    `json:"status" gorm:"size:32;default:pending"` // pending, approved, rejected, processed
	ReviewNote    string    `json:"review_note"`
	ReviewedBy    string    `json:"reviewed_by"`
	ReviewedAt    *time.Time `json:"reviewed_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int32 `json:"max_retries"`
	InitialDelay    int32 `json:"initial_delay"`    // 初始延迟(秒)
	MaxDelay        int32 `json:"max_delay"`        // 最大延迟(秒)
	BackoffFactor   float32 `json:"backoff_factor"` // 退避因子
	EnableFallback  bool  `json:"enable_fallback"`
	EnableAutoRetry bool  `json:"enable_auto_retry"`
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	cfg := config.GetConfig()
	return &RetryConfig{
		MaxRetries:      cfg.Retry.MaxRetries,
		InitialDelay:    cfg.Retry.InitialDelay,
		MaxDelay:        cfg.Retry.MaxDelay,
		BackoffFactor:   cfg.Retry.BackoffFactor,
		EnableFallback:  cfg.Retry.EnableFallback,
		EnableAutoRetry: cfg.Retry.EnableAutoRetry,
	}
}

// Channel 渠道
type Channel struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int64     `json:"user_id" gorm:"index;not null"`
	ChannelType    string    `json:"channel_type" gorm:"size:32;not null"`
	ChannelName    string    `json:"channel_name" gorm:"size:128;not null"`
	ChannelConfig  string    `json:"channel_config" gorm:"type:text"`
	IsEnabled      bool      `json:"is_enabled" gorm:"default:true;index"`
	TitleTemplate  string    `json:"title_template" gorm:"type:text"`
	BodyTemplate   string    `json:"body_template" gorm:"type:text"`
	TagsTemplate   string    `json:"tags_template" gorm:"size:512"`
	CoverTemplate  string    `json:"cover_template" gorm:"size:512"`
	GEOConfig      string    `json:"geo_config" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ChannelGEOConfig 渠道GEO配置
type ChannelGEOConfig struct {
	// 隐藏语义标记
	InjectHiddenMarkup  bool   `json:"inject_hidden_markup"`
	HiddenMarkupFormat  string `json:"hidden_markup_format"`  // html_comment, microdata, json_ld
	HiddenMarkupPrefix  string `json:"hidden_markup_prefix"`

	// Schema Markup
	InjectSchemaMarkup  bool   `json:"inject_schema_markup"`
	SchemaType          string `json:"schema_type"`          // article, product, faq

	// 内容优化
	AutoOptimizeTitle   bool   `json:"auto_optimize_title"`
	AutoOptimizeBody    bool   `json:"auto_optimize_body"`
	MaxTitleLength      int    `json:"max_title_length"`
	MaxBodyLength       int    `json:"max_body_length"`

	// AIGC标识
	InjectAIGCLabel     bool   `json:"inject_aigc_label"`
	AIGCLabelTemplate   string `json:"aigc_label_template"`

	// 变量替换
	Variables           map[string]string `json:"variables"`
}

// PublishConfig 发布配置（用于预览和实际发布）
type PublishConfig struct {
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Tags        []string          `json:"tags"`
	CoverURL    string            `json:"cover_url"`
	SchemaMarkup string           `json:"schema_markup"`
	HiddenMarkup string           `json:"hidden_markup"`
	AIGCLabel    string           `json:"aigc_label"`
	Variables    map[string]string `json:"variables"`
	ChannelConfig *ChannelGEOConfig `json:"channel_config"`
}

// ChannelType 渠道类型
const (
	ChannelTypeWechat  = "wechat"
	ChannelTypeWeibo   = "weibo"
	ChannelTypeZhihu   = "zhihu"
	ChannelTypeToutiao = "toutiao"
	ChannelTypeDouyin  = "douyin"
	ChannelTypeXiaohongshu = "xiaohongshu"
)

// PublishResult 发布结果
type PublishResult struct {
	TaskID      int64  `json:"task_id"`
	ChannelID   int64  `json:"channel_id"`
	Success     bool   `json:"success"`
	ExternalID  string `json:"external_id"`
	ErrorMsg    string `json:"error_msg"`
	PublishedAt time.Time `json:"published_at"`
}