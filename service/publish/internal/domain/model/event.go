package model

import "time"

// ==================== 领域事件 ====================

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventType() string
	GetOccurredAt() time.Time
}

// ContentOptimizedEvent 内容优化完成事件
type ContentOptimizedEvent struct {
	ContentID        int64     `json:"content_id"`
	UserID           int64     `json:"user_id"`
	OptimizationType string    `json:"optimization_type"`
	Score            float32   `json:"score"`
	EventTime        time.Time `json:"event_time"`
}

func (e *ContentOptimizedEvent) EventType() string    { return "content.optimized" }
func (e *ContentOptimizedEvent) GetOccurredAt() time.Time { return e.EventTime }

// PublishRequestedEvent 发布请求事件
type PublishRequestedEvent struct {
	TaskID        int64     `json:"task_id"`
	UserID        int64     `json:"user_id"`
	ContentID     int64     `json:"content_id"`
	ChannelID     int64     `json:"channel_id"`
	ScheduledTime *time.Time `json:"scheduled_time"`
	EventTime     time.Time `json:"event_time"`
}

func (e *PublishRequestedEvent) EventType() string    { return "publish.requested" }
func (e *PublishRequestedEvent) GetOccurredAt() time.Time { return e.EventTime }

// PublishSuccessEvent 发布成功事件
type PublishSuccessEvent struct {
	TaskID       int64     `json:"task_id"`
	UserID       int64     `json:"user_id"`
	ContentID    int64     `json:"content_id"`
	ChannelID    int64     `json:"channel_id"`
	ExternalID   string    `json:"external_id"`
	PublishedAt  time.Time `json:"published_at"`
	EventTime    time.Time `json:"event_time"`
}

func (e *PublishSuccessEvent) EventType() string    { return "publish.success" }
func (e *PublishSuccessEvent) GetOccurredAt() time.Time { return e.EventTime }

// PublishFailedEvent 发布失败事件
type PublishFailedEvent struct {
	TaskID      int64     `json:"task_id"`
	UserID      int64     `json:"user_id"`
	ContentID   int64     `json:"content_id"`
	ChannelID   int64     `json:"channel_id"`
	ErrorMsg    string    `json:"error_msg"`
	EventTime   time.Time `json:"event_time"`
}

func (e *PublishFailedEvent) EventType() string    { return "publish.failed" }
func (e *PublishFailedEvent) GetOccurredAt() time.Time { return e.EventTime }

// ==================== 集成事件（Kafka消息） ====================

// KafkaMessage Kafka消息封装
type KafkaMessage struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key"`
	Value     []byte            `json:"value"`
	Headers   map[string]string `json:"headers"`
	Timestamp time.Time         `json:"timestamp"`
}

// EventMetadata 事件元数据
type EventMetadata struct {
	EventID       string `json:"event_id"`
	EventType     string `json:"event_type"`
	Source        string `json:"source"`
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"`
	Timestamp     time.Time `json:"timestamp"`
}
