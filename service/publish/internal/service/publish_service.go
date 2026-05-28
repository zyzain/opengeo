package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/service/publish/internal/adapter/outbound/channel"
	"opengeo/service/publish/internal/dal"
	"opengeo/service/publish/internal/domain/model"
)

// PublishService 发布服务
type PublishService struct {
	publishRepo    *dal.PublishTaskRepository
	channelRepo    *dal.ChannelRepository
	channelFactory *channel.ChannelAdapterFactory
	antiBanEngine  *AntiBanEngine
	dedupSvc       *DeduplicationService
	retrySvc       *RetryService
}

// NewPublishService 创建发布服务
func NewPublishService(
	publishRepo *dal.PublishTaskRepository,
	channelRepo *dal.ChannelRepository,
	channelFactory *channel.ChannelAdapterFactory,
	dedupSvc *DeduplicationService,
	retrySvc *RetryService,
) *PublishService {
	return &PublishService{
		publishRepo:    publishRepo,
		channelRepo:    channelRepo,
		channelFactory: channelFactory,
		antiBanEngine:  NewAntiBanEngine(),
		dedupSvc:       dedupSvc,
		retrySvc:       retrySvc,
	}
}

// CreatePublishTask 创建发布任务
func (s *PublishService) CreatePublishTask(ctx context.Context, userID, contentID, channelID int64, scheduledTime *time.Time) (*model.PublishTask, error) {
	// 验证渠道是否存在
	channel, err := s.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	if !channel.IsEnabled {
		return nil, fmt.Errorf("channel is disabled")
	}

	// 创建发布任务
	task := &model.PublishTask{
		UserID:        userID,
		ContentID:     contentID,
		ChannelID:     channelID,
		Status:        model.PublishStatusPending,
		ScheduledTime: scheduledTime,
		MaxRetries:    3,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.publishRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create publish task: %w", err)
	}

	return task, nil
}

// GetPublishTask 获取发布任务
func (s *PublishService) GetPublishTask(ctx context.Context, taskID int64) (*model.PublishTask, error) {
	task, err := s.publishRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get publish task: %w", err)
	}

	return task, nil
}

// ListPublishTasks 列出发布任务
func (s *PublishService) ListPublishTasks(ctx context.Context, userID int64, status model.PublishStatus, page, pageSize int) ([]*model.PublishTask, int32, error) {
	tasks, total, err := s.publishRepo.List(ctx, userID, status, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list publish tasks: %w", err)
	}

	return tasks, total, nil
}

// CancelPublishTask 取消发布任务
func (s *PublishService) CancelPublishTask(ctx context.Context, taskID int64) error {
	task, err := s.publishRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get publish task: %w", err)
	}

	if task.Status != model.PublishStatusPending {
		return fmt.Errorf("cannot cancel task with status: %d", task.Status)
	}

	task.Status = model.PublishStatusCancelled
	task.UpdatedAt = time.Now()

	if err := s.publishRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to cancel publish task: %w", err)
	}

	return nil
}

// RetryPublishTask 重试发布任务
func (s *PublishService) RetryPublishTask(ctx context.Context, taskID int64) error {
	task, err := s.publishRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get publish task: %w", err)
	}

	if !task.CanRetry() {
		return fmt.Errorf("task cannot be retried")
	}

	task.IncrementRetry()

	if err := s.publishRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to retry publish task: %w", err)
	}

	return nil
}

// ==================== 发布预览 ====================

// PreviewPublish 发布预览
func (s *PublishService) PreviewPublish(ctx context.Context, channelID int64, title, body string, mediaURLs []string) (*channel.PreviewResponse, error) {
	// 获取渠道信息
	channelInfo, err := s.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	// 获取渠道适配器
	adapter, err := s.channelFactory.GetAdapter(channelInfo.ChannelType)
	if err != nil {
		return nil, fmt.Errorf("unsupported channel type: %w", err)
	}

	// 调用预览
	req := &channel.PreviewRequest{
		Title:     title,
		Body:      body,
		MediaURLs: mediaURLs,
	}

	resp, err := adapter.Preview(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to preview: %w", err)
	}

	return resp, nil
}

// ==================== 发布校验 ====================

// ValidatePublish 发布校验
func (s *PublishService) ValidatePublish(ctx context.Context, channelID int64, title, body string, mediaURLs []string) (*channel.ValidateResponse, error) {
	// 获取渠道信息
	channelInfo, err := s.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	// 获取渠道适配器
	adapter, err := s.channelFactory.GetAdapter(channelInfo.ChannelType)
	if err != nil {
		return nil, fmt.Errorf("unsupported channel type: %w", err)
	}

	// 调用校验
	req := &channel.ValidateRequest{
		Title:     title,
		Body:      body,
		MediaURLs: mediaURLs,
	}

	resp, err := adapter.Validate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate: %w", err)
	}

	return resp, nil
}

// ==================== 执行发布 ====================

// ExecutePublishTask 执行发布任务（带防封保护）
func (s *PublishService) ExecutePublishTask(ctx context.Context, taskID int64, content *ContentData) (*model.PublishResult, error) {
	task, err := s.publishRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get publish task: %w", err)
	}

	if task.Status != model.PublishStatusPending {
		return nil, fmt.Errorf("task is not in pending status")
	}

	// 标记为发布中
	task.MarkAsPublishing()
	if err := s.publishRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	// 获取渠道信息
	channelInfo, err := s.channelRepo.GetByID(ctx, task.ChannelID)
	if err != nil {
		task.MarkAsFailed("channel not found")
		s.publishRepo.Update(ctx, task)
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	// ===== 防封引擎：准备发布环境 =====
	pubCtx, err := s.antiBanEngine.PreparePublish(ctx, task.UserID, channelInfo.ChannelType)
	if err != nil {
		task.MarkAsFailed(fmt.Sprintf("anti-ban prepare failed: %v", err))
		s.publishRepo.Update(ctx, task)
		return nil, fmt.Errorf("anti-ban prepare failed: %w", err)
	}

	// 记录防封信息到任务
	task.ReviewNote = fmt.Sprintf("proxy=%s, fingerprint=%d, delay=%v",
		pubCtx.Proxy.Type, pubCtx.Fingerprint.ID, pubCtx.Delay)

	// 等待随机延迟（模拟真人操作间隔）
	select {
	case <-ctx.Done():
		task.MarkAsFailed("context cancelled")
		s.publishRepo.Update(ctx, task)
		return nil, ctx.Err()
	case <-time.After(pubCtx.Delay):
		// 延迟结束，继续执行
	}

	// 获取渠道适配器
	adapter, err := s.channelFactory.GetAdapter(channelInfo.ChannelType)
	if err != nil {
		task.MarkAsFailed("unsupported channel type")
		s.publishRepo.Update(ctx, task)
		return nil, fmt.Errorf("unsupported channel type: %w", err)
	}

	// 先进行校验
	validateReq := &channel.ValidateRequest{
		Title:     content.Title,
		Body:      content.Body,
		MediaURLs: content.MediaURLs,
	}
	validateResp, err := adapter.Validate(ctx, validateReq)
	if err != nil {
		task.MarkAsFailed("validation failed")
		s.publishRepo.Update(ctx, task)
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if !validateResp.Valid {
		task.MarkAsFailed(fmt.Sprintf("validation failed: %v", validateResp.Errors))
		s.publishRepo.Update(ctx, task)
		return nil, fmt.Errorf("validation failed: %v", validateResp.Errors)
	}

	// ===== 去重处理 =====
	if s.dedupSvc != nil {
		dedupReq := &DeduplicateRequest{
			Title:    content.Title,
			Body:     content.Body,
			Tags:     content.Tags,
			Strategy: "medium",
		}
		dedupResult := s.dedupSvc.Deduplicate(dedupReq)
		content.Title = dedupResult.Title
		content.Body = dedupResult.Body
		content.Tags = dedupResult.Tags
	}

	// 执行发布
	publishReq := &channel.PublishRequest{
		ContentID: content.ContentID,
		Title:     content.Title,
		Body:      content.Body,
		MediaURLs: content.MediaURLs,
		Tags:      content.Tags,
	}

	adapterPubCtx := &channel.AdapterPublishContext{
		ProxyURL:      fmt.Sprintf("%s://%s:%d", pubCtx.Proxy.Type, pubCtx.Proxy.IP, pubCtx.Proxy.Port),
		Headers:       pubCtx.Headers,
		UserAgent:     pubCtx.UserAgent,
		FingerprintID: pubCtx.Fingerprint.ID,
	}

	publishResp, err := adapter.Publish(ctx, publishReq, adapterPubCtx)
	if err != nil {
		if s.retrySvc != nil {
			retryConfig := model.DefaultRetryConfig()
			retryDecision := s.retrySvc.ProcessRetryResult(task, retryConfig, err)
			switch retryDecision.Action {
			case "retry":
				task.ReviewNote = retryDecision.Reason
			case "fallback":
				task.ReviewNote = retryDecision.Reason
			case "review":
				task.ReviewNote = retryDecision.Reason
			}
		} else {
			task.MarkAsFailed(err.Error())
		}
		s.publishRepo.Update(ctx, task)
		return nil, fmt.Errorf("failed to publish: %w", err)
	}

	if publishResp.Success {
		task.MarkAsSuccess()
		s.antiBanEngine.AfterPublish(task.UserID, channelInfo.ChannelType, pubCtx.Proxy.ID, true)
	} else {
		if s.retrySvc != nil {
			retryConfig := model.DefaultRetryConfig()
			retryDecision := s.retrySvc.ProcessRetryResult(task, retryConfig, fmt.Errorf("%s", publishResp.ErrorMsg))
			switch retryDecision.Action {
			case "retry":
				task.ReviewNote = retryDecision.Reason
			case "fallback":
				task.ReviewNote = retryDecision.Reason
			case "review":
				task.ReviewNote = retryDecision.Reason
			}
		} else {
			task.MarkAsFailed(publishResp.ErrorMsg)
		}
		s.antiBanEngine.AfterPublish(task.UserID, channelInfo.ChannelType, pubCtx.Proxy.ID, false)
	}

	if err := s.publishRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return &model.PublishResult{
		TaskID:      task.ID,
		ChannelID:   task.ChannelID,
		Success:     publishResp.Success,
		ExternalID:  publishResp.ExternalID,
		ErrorMsg:    publishResp.ErrorMsg,
		PublishedAt: publishResp.PublishedAt,
	}, nil
}

// ContentData 内容数据
type ContentData struct {
	ContentID int64    `json:"content_id"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	MediaURLs []string `json:"media_urls"`
	Tags      []string `json:"tags"`
}

// ==================== 渠道管理 ====================

// CreateChannel 创建渠道
func (s *PublishService) CreateChannel(ctx context.Context, userID int64, channelType, channelName, channelConfig string) (*model.Channel, error) {
	// 验证渠道类型是否支持
	_, err := s.channelFactory.GetAdapter(channelType)
	if err != nil {
		return nil, fmt.Errorf("unsupported channel type: %s", channelType)
	}

	channel := &model.Channel{
		UserID:        userID,
		ChannelType:   channelType,
		ChannelName:   channelName,
		ChannelConfig: channelConfig,
		IsEnabled:     true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.channelRepo.Create(ctx, channel); err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return channel, nil
}

// GetChannel 获取渠道
func (s *PublishService) GetChannel(ctx context.Context, channelID int64) (*model.Channel, error) {
	channel, err := s.channelRepo.GetByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	return channel, nil
}

// ListChannels 列出渠道
func (s *PublishService) ListChannels(ctx context.Context, userID int64) ([]*model.Channel, error) {
	channels, err := s.channelRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	return channels, nil
}

// GetSupportedPlatforms 获取支持的平台
func (s *PublishService) GetSupportedPlatforms() []string {
	return s.channelFactory.GetSupportedPlatforms()
}

// GetPendingTasks 获取待处理任务
func (s *PublishService) GetPendingTasks(ctx context.Context, limit int) ([]*model.PublishTask, error) {
	tasks, err := s.publishRepo.GetPendingTasks(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending tasks: %w", err)
	}

	return tasks, nil
}