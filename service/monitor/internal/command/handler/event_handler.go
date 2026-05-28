package handler

import (
	"context"
	"log"

	"opengeo/service/monitor/internal/command/service"
	"opengeo/service/monitor/internal/domain/model"
)

// ContentOptimizedEventHandler 内容优化事件处理器（写侧）
type ContentOptimizedEventHandler struct {
	monitorService *service.MonitorCommandService
}

// NewContentOptimizedEventHandler 创建内容优化事件处理器
func NewContentOptimizedEventHandler(monitorService *service.MonitorCommandService) *ContentOptimizedEventHandler {
	return &ContentOptimizedEventHandler{monitorService: monitorService}
}

// Handle 处理内容优化事件
func (h *ContentOptimizedEventHandler) Handle(ctx context.Context, event *model.ContentOptimizedEvent) error {
	log.Printf("Handling content optimized event: content_id=%d, score=%f", event.ContentID, event.Score)

	// 记录优化结果
	if err := h.monitorService.RecordOptimization(ctx, event); err != nil {
		log.Printf("Failed to record optimization: %v", err)
		return err
	}

	return nil
}

// PublishSuccessEventHandler 发布成功事件处理器（写侧）
type PublishSuccessEventHandler struct {
	monitorService *service.MonitorCommandService
}

// NewPublishSuccessEventHandler 创建发布成功事件处理器
func NewPublishSuccessEventHandler(monitorService *service.MonitorCommandService) *PublishSuccessEventHandler {
	return &PublishSuccessEventHandler{monitorService: monitorService}
}

// Handle 处理发布成功事件
func (h *PublishSuccessEventHandler) Handle(ctx context.Context, event *model.PublishSuccessEvent) error {
	log.Printf("Handling publish success event: task_id=%d, channel_id=%d", event.TaskID, event.ChannelID)

	// 记录发布成功
	if err := h.monitorService.RecordPublishSuccess(ctx, event); err != nil {
		log.Printf("Failed to record publish success: %v", err)
		return err
	}

	return nil
}

// AICitationFoundEventHandler AI引用发现事件处理器（写侧）
type AICitationFoundEventHandler struct {
	monitorService *service.MonitorCommandService
}

// NewAICitationFoundEventHandler 创建AI引用发现事件处理器
func NewAICitationFoundEventHandler(monitorService *service.MonitorCommandService) *AICitationFoundEventHandler {
	return &AICitationFoundEventHandler{monitorService: monitorService}
}

// Handle 处理AI引用发现事件
func (h *AICitationFoundEventHandler) Handle(ctx context.Context, event *model.AICitationFoundEvent) error {
	log.Printf("Handling AI citation found event: content_id=%d, ai_model=%s", event.ContentID, event.AIModel)

	// 记录AI引用
	if err := h.monitorService.RecordAICitation(ctx, event); err != nil {
		log.Printf("Failed to record AI citation: %v", err)
		return err
	}

	return nil
}