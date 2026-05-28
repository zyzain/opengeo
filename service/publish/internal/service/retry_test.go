package service

import (
	"fmt"
	"testing"
	"time"

	"opengeo/service/publish/internal/domain/model"
)

func newTestTask() *model.PublishTask {
	return &model.PublishTask{
		ID:           1,
		UserID:       1,
		ContentID:    1,
		ChannelID:    1,
		Status:       model.PublishStatusFailed,
		RetryCount:   0,
		MaxRetries:   3,
		RetryDelay:   30,
		ErrorMessage: "test error",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func newTestRetryConfig() *model.RetryConfig {
	return &model.RetryConfig{
		MaxRetries:      3,
		InitialDelay:    30,
		MaxDelay:        3600,
		BackoffFactor:   2.0,
		EnableFallback:  true,
		EnableAutoRetry: true,
	}
}

func TestProcessRetryResult_Success(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()
	config := newTestRetryConfig()

	decision := svc.ProcessRetryResult(task, config, nil)

	if decision.Action != "success" {
		t.Errorf("expected success, got %s", decision.Action)
	}
}

func TestProcessRetryResult_Retry(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()
	config := newTestRetryConfig()

	decision := svc.ProcessRetryResult(task, config, fmt.Errorf("publish failed"))

	if decision.Action != "retry" {
		t.Errorf("expected retry, got %s", decision.Action)
	}
	if decision.Delay == 0 {
		t.Error("expected non-zero delay")
	}
}

func TestProcessRetryResult_RetryExhausted(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()
	task.RetryCount = 3 // 已用完重试次数
	config := newTestRetryConfig()

	decision := svc.ProcessRetryResult(task, config, fmt.Errorf("publish failed"))

	if decision.Action != "review" {
		t.Errorf("expected review, got %s", decision.Action)
	}
}

func TestProcessRetryResult_Fallback(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()
	task.RetryCount = 3 // 重试次数用完
	fallbackID := int64(2)
	task.FallbackChannelID = &fallbackID
	config := newTestRetryConfig()

	decision := svc.ProcessRetryResult(task, config, fmt.Errorf("publish failed"))

	if decision.Action != "fallback" {
		t.Errorf("expected fallback, got %s", decision.Action)
	}
	if decision.FallbackChannelID != 2 {
		t.Errorf("expected fallback channel 2, got %d", decision.FallbackChannelID)
	}
}

func TestProcessRetryResult_DisabledAutoRetry(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()
	config := newTestRetryConfig()
	config.EnableAutoRetry = false

	decision := svc.ProcessRetryResult(task, config, fmt.Errorf("publish failed"))

	if decision.Action != "review" {
		t.Errorf("expected review when auto retry disabled, got %s", decision.Action)
	}
}

func TestPrepareRetry(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()

	svc.PrepareRetry(task)

	if task.RetryCount != 1 {
		t.Errorf("expected retry count 1, got %d", task.RetryCount)
	}
	if task.Status != model.PublishStatusRetrying {
		t.Errorf("expected retrying status, got %d", task.Status)
	}
}

func TestPrepareFallback(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()

	svc.PrepareFallback(task, 2)

	if task.FallbackChannelID == nil || *task.FallbackChannelID != 2 {
		t.Error("expected fallback channel 2")
	}
	if task.Status != model.PublishStatusFallback {
		t.Errorf("expected fallback status, got %d", task.Status)
	}
}

func TestPrepareReview(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()

	svc.PrepareReview(task, "需要人工审核")

	if !task.IsManuallyReview {
		t.Error("expected manually review flag")
	}
	if task.ReviewNote != "需要人工审核" {
		t.Errorf("expected review note, got %s", task.ReviewNote)
	}
	if task.Status != model.PublishStatusReview {
		t.Errorf("expected review status, got %d", task.Status)
	}
}

func TestCreateFallbackQueueEntry(t *testing.T) {
	svc := NewRetryService()
	task := newTestTask()

	entry := svc.CreateFallbackQueueEntry(task, 1, 2, "主渠道失败")

	if entry.TaskID != 1 {
		t.Errorf("expected task ID 1, got %d", entry.TaskID)
	}
	if entry.OriginalChannelID != 1 {
		t.Errorf("expected original channel 1, got %d", entry.OriginalChannelID)
	}
	if entry.FallbackChannelID != 2 {
		t.Errorf("expected fallback channel 2, got %d", entry.FallbackChannelID)
	}
	if entry.Status != "pending" {
		t.Errorf("expected pending status, got %s", entry.Status)
	}
}

func TestApproveFallback(t *testing.T) {
	svc := NewRetryService()
	entry := &model.FallbackQueue{Status: "pending"}

	svc.ApproveFallback(entry, "admin", "批准降级")

	if entry.Status != "approved" {
		t.Errorf("expected approved, got %s", entry.Status)
	}
	if entry.ReviewedBy != "admin" {
		t.Errorf("expected admin, got %s", entry.ReviewedBy)
	}
	if entry.ReviewedAt == nil {
		t.Error("expected reviewed at time")
	}
}

func TestRejectFallback(t *testing.T) {
	svc := NewRetryService()
	entry := &model.FallbackQueue{Status: "pending"}

	svc.RejectFallback(entry, "admin", "拒绝降级")

	if entry.Status != "rejected" {
		t.Errorf("expected rejected, got %s", entry.Status)
	}
}

func TestGetNextRetryDelay(t *testing.T) {
	task := newTestTask()
	task.RetryDelay = 30

	// 第1次重试: 30 * 2^0 = 30秒
	task.RetryCount = 0
	delay1 := task.GetNextRetryDelay()
	if delay1 != 30*time.Second {
		t.Errorf("expected 30s, got %v", delay1)
	}

	// 第2次重试: 30 * 2^1 = 60秒
	task.RetryCount = 1
	delay2 := task.GetNextRetryDelay()
	if delay2 != 60*time.Second {
		t.Errorf("expected 60s, got %v", delay2)
	}

	// 第3次重试: 30 * 2^2 = 120秒
	task.RetryCount = 2
	delay3 := task.GetNextRetryDelay()
	if delay3 != 120*time.Second {
		t.Errorf("expected 120s, got %v", delay3)
	}
}

func TestGetNextRetryDelay_MaxDelay(t *testing.T) {
	task := newTestTask()
	task.RetryDelay = 30
	task.RetryCount = 10 // 大次数

	delay := task.GetNextRetryDelay()
	if delay > 3600*time.Second {
		t.Errorf("delay exceeds max: %v", delay)
	}
}

func TestAddErrorHistory(t *testing.T) {
	task := newTestTask()

	task.AddErrorHistory("error 1")
	task.AddErrorHistory("error 2")

	if task.ErrorHistory == "" {
		t.Error("expected error history")
	}
	if !containsStr(task.ErrorHistory, "error 1") {
		t.Error("expected error 1 in history")
	}
	if !containsStr(task.ErrorHistory, "error 2") {
		t.Error("expected error 2 in history")
	}
}

func TestCanRetry(t *testing.T) {
	task := newTestTask()

	// 可以重试
	task.RetryCount = 0
	task.MaxRetries = 3
	task.Status = model.PublishStatusFailed
	if !task.CanRetry() {
		t.Error("expected can retry")
	}

	// 重试次数用完
	task.RetryCount = 3
	if task.CanRetry() {
		t.Error("expected cannot retry")
	}
}

func TestCanFallback(t *testing.T) {
	task := newTestTask()

	// 无备用渠道
	task.Status = model.PublishStatusFailed
	if task.CanFallback() {
		t.Error("expected cannot fallback")
	}

	// 有备用渠道
	fallbackID := int64(2)
	task.FallbackChannelID = &fallbackID
	if !task.CanFallback() {
		t.Error("expected can fallback")
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	config := model.DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("expected max retries 3, got %d", config.MaxRetries)
	}
	if config.InitialDelay != 30 {
		t.Errorf("expected initial delay 30, got %d", config.InitialDelay)
	}
	if !config.EnableAutoRetry {
		t.Error("expected auto retry enabled")
	}
}

func TestPublishStatusConstants(t *testing.T) {
	if model.PublishStatusPending != 0 {
		t.Error("expected pending = 0")
	}
	if model.PublishStatusSuccess != 2 {
		t.Error("expected success = 2")
	}
	if model.PublishStatusReview != 7 {
		t.Error("expected review = 7")
	}
}
