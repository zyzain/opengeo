package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"opengeo/pkg/config"
	"opengeo/service/system/internal/dal"
	"opengeo/service/system/internal/domain/model"
)

// WebhookDeliveryService Webhook投递服务
type WebhookDeliveryService struct {
	webhookRepo *dal.WebhookRepository
	httpClient  *http.Client
	maxRetries  int
}

// NewWebhookDeliveryService 创建Webhook投递服务
func NewWebhookDeliveryService(webhookRepo *dal.WebhookRepository) *WebhookDeliveryService {
	cfg := config.GetConfig()
	return &WebhookDeliveryService{
		webhookRepo: webhookRepo,
		httpClient: &http.Client{
			Timeout: cfg.Webhook.Timeout,
		},
		maxRetries: cfg.Webhook.MaxRetries,
	}
}

// DeliveryResult 投递结果
type DeliveryResult struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"status_code"`
	ResponseBody string `json:"response_body"`
	Error        string `json:"error,omitempty"`
	Duration     int64  `json:"duration_ms"`
	Attempt      int    `json:"attempt"`
}

// DeliverPayload 投递Webhook负载
func (s *WebhookDeliveryService) DeliverPayload(ctx context.Context, webhook *model.Webhook, eventType string, payload interface{}) (*DeliveryResult, error) {
	// 序列化负载
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	// 生成签名
	signature := GenerateHMACSignature(webhook.Secret, string(payloadJSON))

	// 构建请求头
	headers := map[string]string{
		"Content-Type":          "application/json",
		"X-Webhook-Event":       eventType,
		"X-Webhook-Signature":   signature,
		"X-Webhook-Timestamp":   fmt.Sprintf("%d", time.Now().Unix()),
		"User-Agent":            "OpenGEO-Webhook/1.0",
	}

	// 带重试的投递
	var result *DeliveryResult
	for attempt := 1; attempt <= s.maxRetries; attempt++ {
		result, err = s.doHTTPRequest(ctx, webhook.URL, payloadJSON, headers)
		if err == nil && result.StatusCode >= 200 && result.StatusCode < 300 {
			result.Attempt = attempt
			break
		}

		if attempt < s.maxRetries {
			// 指数退避
			backoff := time.Duration(attempt*attempt) * time.Second
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(backoff):
			}
		}
		result.Attempt = attempt
	}

	// 记录事件
	webhookEvent := &model.WebhookEvent{
		WebhookID:    webhook.ID,
		EventType:    eventType,
		Payload:      string(payloadJSON),
		StatusCode:   int32(result.StatusCode),
		ResponseBody: truncateString(result.ResponseBody, 2000),
		Success:      result.Success,
		TriggeredAt:  time.Now(),
		CreatedAt:    time.Now(),
	}

	if err := s.webhookRepo.CreateEvent(ctx, webhookEvent); err != nil {
		// 记录事件失败不影响投递结果
		fmt.Printf("failed to create webhook event: %v\n", err)
	}

	// 更新Webhook状态
	if result.Success {
		webhook.LastTrigger = &webhookEvent.TriggeredAt
		webhook.FailCount = 0
	} else {
		webhook.FailCount++
	}
	s.webhookRepo.Update(ctx, webhook)

	return result, nil
}

// doHTTPRequest 执行HTTP请求
func (s *WebhookDeliveryService) doHTTPRequest(ctx context.Context, url string, body []byte, headers map[string]string) (*DeliveryResult, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return &DeliveryResult{
			Success: false,
			Error:   fmt.Sprintf("create request: %v", err),
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &DeliveryResult{
			Success: false,
			Error:   fmt.Sprintf("http request: %v", err),
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &DeliveryResult{
			Success:    false,
			StatusCode: resp.StatusCode,
			Error:      fmt.Sprintf("read response: %v", err),
			Duration:   time.Since(start).Milliseconds(),
		}, nil
	}

	return &DeliveryResult{
		Success:      resp.StatusCode >= 200 && resp.StatusCode < 300,
		StatusCode:   resp.StatusCode,
		ResponseBody: string(respBody),
		Duration:     time.Since(start).Milliseconds(),
	}, nil
}

// ==================== 签名验证 ====================

// GenerateHMACSignature 生成HMAC-SHA256签名
func GenerateHMACSignature(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// VerifyHMACSignature 验证HMAC-SHA256签名
func VerifyHMACSignature(secret, payload, signature string) bool {
	expected := GenerateHMACSignature(secret, payload)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// ==================== 批量投递 ====================

// BatchDeliveryRequest 批量投递请求
type BatchDeliveryRequest struct {
	WebhookIDs []int64     `json:"webhook_ids"`
	EventType  string      `json:"event_type"`
	Payload    interface{} `json:"payload"`
}

// BatchDeliveryResult 批量投递结果
type BatchDeliveryResult struct {
	Total   int              `json:"total"`
	Success  int             `json:"success"`
	Failed   int             `json:"failed"`
	Results  []*DeliveryResult `json:"results"`
}

// BatchDeliver 批量投递Webhook
func (s *WebhookDeliveryService) BatchDeliver(ctx context.Context, req *BatchDeliveryRequest) (*BatchDeliveryResult, error) {
	result := &BatchDeliveryResult{
		Total:   len(req.WebhookIDs),
		Results: make([]*DeliveryResult, 0, len(req.WebhookIDs)),
	}

	for _, webhookID := range req.WebhookIDs {
		webhook, err := s.webhookRepo.GetByID(ctx, webhookID)
		if err != nil {
			result.Failed++
			result.Results = append(result.Results, &DeliveryResult{
				Success: false,
				Error:   fmt.Sprintf("webhook not found: %d", webhookID),
			})
			continue
		}

		if !webhook.IsActive {
			result.Failed++
			result.Results = append(result.Results, &DeliveryResult{
				Success: false,
				Error:   "webhook is disabled",
			})
			continue
		}

		deliveryResult, err := s.DeliverPayload(ctx, webhook, req.EventType, req.Payload)
		if err != nil {
			result.Failed++
			result.Results = append(result.Results, &DeliveryResult{
				Success: false,
				Error:   err.Error(),
			})
			continue
		}

		if deliveryResult.Success {
			result.Success++
		} else {
			result.Failed++
		}
		result.Results = append(result.Results, deliveryResult)
	}

	return result, nil
}

// ==================== Webhook管理增强 ====================

// TriggerWebhookByEvent 通过事件类型触发所有匹配的Webhook
func (s *WebhookDeliveryService) TriggerWebhookByEvent(ctx context.Context, userID int64, eventType string, payload interface{}) (*BatchDeliveryResult, error) {
	// 获取用户的所有Webhook
	webhooks, _, err := s.webhookRepo.List(ctx, userID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("list webhooks: %w", err)
	}

	// 过滤匹配事件类型的Webhook
	matchedIDs := make([]int64, 0)
	for _, wh := range webhooks {
		if !wh.IsActive {
			continue
		}

		var events []string
		if err := json.Unmarshal([]byte(wh.Events), &events); err != nil {
			continue
		}

		for _, event := range events {
			if event == eventType || event == "*" {
				matchedIDs = append(matchedIDs, wh.ID)
				break
			}
		}
	}

	if len(matchedIDs) == 0 {
		return &BatchDeliveryResult{Total: 0}, nil
	}

	return s.BatchDeliver(ctx, &BatchDeliveryRequest{
		WebhookIDs: matchedIDs,
		EventType:  eventType,
		Payload:    payload,
	})
}

// GetDeliveryStats 获取投递统计
func (s *WebhookDeliveryService) GetDeliveryStats(ctx context.Context, webhookID int64) (*DeliveryStats, error) {
	events, total, err := s.webhookRepo.GetEvents(ctx, webhookID, 1, 1000)
	if err != nil {
		return nil, err
	}

	stats := &DeliveryStats{
		Total:   total,
		Success: 0,
		Failed:  0,
	}

	for _, event := range events {
		if event.Success {
			stats.Success++
		} else {
			stats.Failed++
		}
	}

	if stats.Total > 0 {
		stats.SuccessRate = float32(stats.Success) / float32(stats.Total) * 100
	}

	return stats, nil
}

// DeliveryStats 投递统计
type DeliveryStats struct {
	Total       int32   `json:"total"`
	Success     int32   `json:"success"`
	Failed      int32   `json:"failed"`
	SuccessRate float32 `json:"success_rate"`
}

// ==================== 辅助函数 ====================

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// ==================== 内置事件类型 ====================

// EventTypes 内置事件类型常量
const (
	EventContentCreated    = "content.created"
	EventContentUpdated    = "content.updated"
	EventContentDeleted    = "content.deleted"
	EventContentPublished  = "content.published"
	EventPublishSuccess    = "publish.success"
	EventPublishFailed     = "publish.failed"
	EventAccountHealthAlert = "account.health_alert"
	EventScheduleTriggered = "schedule.triggered"
)

// ValidateEventType 验证事件类型
func ValidateEventType(eventType string) bool {
	validTypes := map[string]bool{
		EventContentCreated:    true,
		EventContentUpdated:    true,
		EventContentDeleted:    true,
		EventContentPublished:  true,
		EventPublishSuccess:    true,
		EventPublishFailed:     true,
		EventAccountHealthAlert: true,
		EventScheduleTriggered: true,
	}
	return validTypes[eventType] || strings.HasPrefix(eventType, "custom.")
}
