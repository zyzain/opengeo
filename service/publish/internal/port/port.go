package port

import (
	"context"

	"opengeo/service/publish/internal/domain/model"
)

// ==================== 入站端口 ====================

// PublishService 发布服务接口（RPC入站端口）
type PublishService interface {
	CreatePublishTask(ctx context.Context, req *CreatePublishTaskRequest) (*PublishTask, error)
	GetPublishTask(ctx context.Context, req *GetPublishTaskRequest) (*PublishTask, error)
	ListPublishTasks(ctx context.Context, req *ListPublishTasksRequest) (*ListPublishTasksResponse, error)
	CancelPublishTask(ctx context.Context, req *CancelPublishTaskRequest) error
	RetryPublishTask(ctx context.Context, req *RetryPublishTaskRequest) error
}

// EventConsumer 事件消费者接口（事件入站端口）
type EventConsumer interface {
	// SubscribeContentOptimized 订阅内容优化完成事件
	SubscribeContentOptimized(ctx context.Context, handler ContentOptimizedHandler) error
	// SubscribePublishRequested 订阅发布请求事件
	SubscribePublishRequested(ctx context.Context, handler PublishRequestedHandler) error
	// Start 启动消费者
	Start(ctx context.Context) error
	// Stop 停止消费者
	Stop() error
}

// ContentOptimizedHandler 内容优化事件处理器
type ContentOptimizedHandler func(ctx context.Context, event *model.ContentOptimizedEvent) error

// PublishRequestedHandler 发布请求事件处理器
type PublishRequestedHandler func(ctx context.Context, event *model.PublishRequestedEvent) error

// ==================== 出站端口 ====================

// PublishTaskRepository 发布任务仓储接口
type PublishTaskRepository interface {
	Create(ctx context.Context, task *model.PublishTask) error
	GetByID(ctx context.Context, id int64) (*model.PublishTask, error)
	Update(ctx context.Context, task *model.PublishTask) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter *PublishTaskFilter) ([]*model.PublishTask, int32, error)
	GetPendingTasks(ctx context.Context, limit int) ([]*model.PublishTask, error)
}

// PublishTaskFilter 发布任务过滤器
type PublishTaskFilter struct {
	UserID   int64
	Status   model.PublishStatus
	Page     int32
	PageSize int32
}

// ChannelRepository 渠道仓储接口
type ChannelRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Channel, error)
	ListByUser(ctx context.Context, userID int64) ([]*model.Channel, error)
}

// ChannelSDK 渠道SDK接口（渠道适配器）
type ChannelSDK interface {
	// Publish 发布内容
	Publish(ctx context.Context, content *ContentData, config *ChannelConfig, publishCtx *PublishContext) (*model.PublishResult, error)
	// Preview 预览
	Preview(ctx context.Context, content *ContentData, config *ChannelConfig) (*PreviewResult, error)
	// Validate 验证
	Validate(ctx context.Context, content *ContentData, config *ChannelConfig) error
}

// PublishContext 发布上下文（防封引擎生成的环境信息）
type PublishContext struct {
	ProxyURL      string            `json:"proxy_url"`
	Headers       map[string]string `json:"headers"`
	UserAgent     string            `json:"user_agent"`
	FingerprintID int64             `json:"fingerprint_id"`
}

// ContentData 内容数据
type ContentData struct {
	Title        string `json:"title"`
	Body         string `json:"body"`
	ContentType  string `json:"content_type"`
	SchemaMarkup string `json:"schema_markup"`
	MediaURLs    []string `json:"media_urls"`
}

// ChannelConfig 渠道配置
type ChannelConfig struct {
	Platform     string            `json:"platform"`
	Credentials  map[string]string `json:"credentials"`
	Settings     map[string]string `json:"settings"`
}

// PreviewResult 预览结果
type PreviewResult struct {
	HTML     string `json:"html"`
	Metadata map[string]interface{} `json:"metadata"`
}

// EventProducer 事件生产者接口（事件出站端口）
type EventProducer interface {
	// PublishContentOptimized 发布内容优化事件
	PublishContentOptimized(ctx context.Context, event *model.ContentOptimizedEvent) error
	// PublishPublishSuccess 发布成功事件
	PublishPublishSuccess(ctx context.Context, event *model.PublishSuccessEvent) error
	// PublishPublishFailed 发布失败事件
	PublishPublishFailed(ctx context.Context, event *model.PublishFailedEvent) error
	// Close 关闭生产者
	Close() error
}

// ==================== 请求/响应模型 ====================

type CreatePublishTaskRequest struct {
	UserID        int64  `json:"user_id"`
	ContentID     int64  `json:"content_id"`
	ChannelID     int64  `json:"channel_id"`
	ScheduledTime string `json:"scheduled_time"`
}

type PublishTask struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id"`
	ContentID     int64  `json:"content_id"`
	ChannelID     int64  `json:"channel_id"`
	Status        int32  `json:"status"`
	ScheduledTime string `json:"scheduled_time"`
	PublishedTime string `json:"published_time"`
	RetryCount    int32  `json:"retry_count"`
	ErrorMessage  string `json:"error_message"`
	CreatedAt     string `json:"created_at"`
}

type GetPublishTaskRequest struct {
	TaskID int64 `json:"task_id"`
}

type ListPublishTasksRequest struct {
	UserID   int64 `json:"user_id"`
	Status   int32 `json:"status"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type ListPublishTasksResponse struct {
	Tasks []*PublishTask `json:"tasks"`
	Total int32          `json:"total"`
}

type CancelPublishTaskRequest struct {
	TaskID int64 `json:"task_id"`
}

type RetryPublishTaskRequest struct {
	TaskID int64 `json:"task_id"`
}