package service

import (
	"testing"
)

func TestGenerateHMACSignature(t *testing.T) {
	secret := "test-secret"
	payload := `{"event":"test","data":"hello"}`

	sig1 := GenerateHMACSignature(secret, payload)
	sig2 := GenerateHMACSignature(secret, payload)

	if sig1 != sig2 {
		t.Error("expected consistent signatures")
	}
	if sig1 == "" {
		t.Error("expected non-empty signature")
	}
	// HMAC-SHA256 格式: sha256=hex
	if len(sig1) != 71 { // "sha256=" (7) + 64 hex chars
		t.Errorf("expected signature length 71, got %d", len(sig1))
	}
}

func TestVerifyHMACSignature(t *testing.T) {
	secret := "test-secret"
	payload := `{"event":"test","data":"hello"}`

	signature := GenerateHMACSignature(secret, payload)

	// 正确签名应该验证通过
	if !VerifyHMACSignature(secret, payload, signature) {
		t.Error("expected valid signature")
	}

	// 错误签名应该验证失败
	if VerifyHMACSignature(secret, payload, "sha256=invalid") {
		t.Error("expected invalid signature")
	}

	// 不同secret应该验证失败
	if VerifyHMACSignature("wrong-secret", payload, signature) {
		t.Error("expected invalid signature with wrong secret")
	}

	// 不同payload应该验证失败
	if VerifyHMACSignature(secret, `{"different":"payload"}`, signature) {
		t.Error("expected invalid signature with different payload")
	}
}

func TestValidateEventType(t *testing.T) {
	tests := []struct {
		eventType string
		want      bool
	}{
		{EventContentCreated, true},
		{EventContentUpdated, true},
		{EventContentDeleted, true},
		{EventContentPublished, true},
		{EventPublishSuccess, true},
		{EventPublishFailed, true},
		{EventAccountHealthAlert, true},
		{EventScheduleTriggered, true},
		{"custom.my_event", true},
		{"invalid_event", false},
		{"", false},
		{"random.string", false},
	}

	for _, tt := range tests {
		t.Run(tt.eventType, func(t *testing.T) {
			if got := ValidateEventType(tt.eventType); got != tt.want {
				t.Errorf("ValidateEventType(%q) = %v, want %v", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		s      string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hello..."},
		{"", 5, ""},
		{"abc", 3, "abc"},
	}

	for _, tt := range tests {
		result := truncateString(tt.s, tt.maxLen)
		if result != tt.want {
			t.Errorf("truncateString(%q, %d) = %q, want %q", tt.s, tt.maxLen, result, tt.want)
		}
	}
}

func TestDeliveryResult(t *testing.T) {
	result := &DeliveryResult{
		Success:    true,
		StatusCode: 200,
		Duration:   150,
		Attempt:    1,
	}

	if !result.Success {
		t.Error("expected success")
	}
	if result.StatusCode != 200 {
		t.Errorf("expected 200, got %d", result.StatusCode)
	}
}

func TestDeliveryStats(t *testing.T) {
	stats := &DeliveryStats{
		Total:   100,
		Success: 95,
		Failed:  5,
	}

	if stats.Total != 100 {
		t.Errorf("expected 100, got %d", stats.Total)
	}

	// 计算成功率
	rate := float32(stats.Success) / float32(stats.Total) * 100
	if rate != 95.0 {
		t.Errorf("expected 95%%, got %f%%", rate)
	}
}

func TestBatchDeliveryRequest(t *testing.T) {
	req := &BatchDeliveryRequest{
		WebhookIDs: []int64{1, 2, 3},
		EventType:  "content.created",
		Payload:    map[string]string{"id": "123"},
	}

	if len(req.WebhookIDs) != 3 {
		t.Errorf("expected 3 webhook IDs, got %d", len(req.WebhookIDs))
	}
	if req.EventType != "content.created" {
		t.Errorf("expected content.created, got %s", req.EventType)
	}
}

func TestEventTypes(t *testing.T) {
	// 验证所有事件类型常量
	events := []string{
		EventContentCreated,
		EventContentUpdated,
		EventContentDeleted,
		EventContentPublished,
		EventPublishSuccess,
		EventPublishFailed,
		EventAccountHealthAlert,
		EventScheduleTriggered,
	}

	for _, event := range events {
		if event == "" {
			t.Error("empty event type")
		}
		if !ValidateEventType(event) {
			t.Errorf("event type %s should be valid", event)
		}
	}
}

func TestHMACSignatureFormat(t *testing.T) {
	secret := "my-secret"
	payload := "test-payload"

	sig := GenerateHMACSignature(secret, payload)

	// 检查格式
	if sig[:7] != "sha256=" {
		t.Errorf("expected prefix 'sha256=', got %s", sig[:7])
	}

	// 检查十六进制部分
	hexPart := sig[7:]
	if len(hexPart) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(hexPart))
	}

	// 验证是有效的十六进制
	for _, c := range hexPart {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("invalid hex char: %c", c)
		}
	}
}
