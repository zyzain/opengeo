package service

import (
	"fmt"
	"time"

	"opengeo/service/publish/internal/domain/model"
)

// RetryService 重试与降级服务
type RetryService struct {
	// 依赖注入
}

// NewRetryService 创建重试与降级服务
func NewRetryService() *RetryService {
	return &RetryService{}
}

// ProcessRetryResult 处理重试结果
func (s *RetryService) ProcessRetryResult(
	task *model.PublishTask,
	config *model.RetryConfig,
	publishErr error,
) *RetryDecision {

	if publishErr == nil {
		return &RetryDecision{Action: "success"}
	}

	// 记录错误历史
	task.AddErrorHistory(publishErr.Error())
	task.MarkAsFailed(publishErr.Error())

	// 判断是否可以自动重试
	if config.EnableAutoRetry && task.CanRetry() {
		delay := task.GetNextRetryDelay()
		return &RetryDecision{
			Action:    "retry",
			Delay:     delay,
			Reason:    fmt.Sprintf("自动重试 (%d/%d)，延迟 %v", task.RetryCount+1, task.MaxRetries, delay),
		}
	}

	// 判断是否可以降级
	if config.EnableFallback && task.CanFallback() {
		return &RetryDecision{
			Action:         "fallback",
			FallbackChannelID: *task.FallbackChannelID,
			Reason:         "主渠道失败，切换到备用渠道",
		}
	}

	// 进入人工审核队列
	return &RetryDecision{
		Action: "review",
		Reason: "重试次数已用完且无备用渠道，进入人工审核队列",
	}
}

// RetryDecision 重试决策
type RetryDecision struct {
	Action            string        `json:"action"` // success, retry, fallback, review
	Delay             time.Duration `json:"delay"`
	FallbackChannelID int64         `json:"fallback_channel_id"`
	Reason            string        `json:"reason"`
}

// PrepareRetry 准备重试
func (s *RetryService) PrepareRetry(task *model.PublishTask) {
	task.IncrementRetry()
	task.MarkAsRetrying()
}

// PrepareFallback 准备降级
func (s *RetryService) PrepareFallback(task *model.PublishTask, fallbackChannelID int64) {
	task.FallbackChannelID = &fallbackChannelID
	task.MarkAsFallback()
}

// PrepareReview 准备人工审核
func (s *RetryService) PrepareReview(task *model.PublishTask, reason string) {
	task.MarkAsReview(reason)
}

// CreateFallbackQueueEntry 创建降级队列条目
func (s *RetryService) CreateFallbackQueueEntry(
	task *model.PublishTask,
	originalChannelID, fallbackChannelID int64,
	reason string,
) *model.FallbackQueue {
	return &model.FallbackQueue{
		TaskID:            task.ID,
		UserID:            task.UserID,
		ContentID:         task.ContentID,
		OriginalChannelID: originalChannelID,
		FallbackChannelID: fallbackChannelID,
		Reason:            reason,
		Status:            "pending",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// ApproveFallback 批准降级
func (s *RetryService) ApproveFallback(entry *model.FallbackQueue, reviewedBy, note string) {
	now := time.Now()
	entry.Status = "approved"
	entry.ReviewedBy = reviewedBy
	entry.ReviewNote = note
	entry.ReviewedAt = &now
	entry.UpdatedAt = now
}

// RejectFallback 拒绝降级
func (s *RetryService) RejectFallback(entry *model.FallbackQueue, reviewedBy, note string) {
	now := time.Now()
	entry.Status = "rejected"
	entry.ReviewedBy = reviewedBy
	entry.ReviewNote = note
	entry.ReviewedAt = &now
	entry.UpdatedAt = now
}
