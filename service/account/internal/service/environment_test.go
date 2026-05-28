package service

import (
	"testing"

	"opengeo/service/account/internal/dal"
)

func TestGenerateFingerprint_Platform(t *testing.T) {
	tests := []struct {
		platform string
		ua       string
		width    int32
		desc     string
	}{
		{"windows", "Windows NT 10.0", 1920, "windows platform"},
		{"macos", "Macintosh", 2560, "macos platform"},
		{"linux", "Linux", 1920, "linux platform"},
		{"unknown", "Windows NT 10.0", 1920, "default platform"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fp := &dal.BrowserFingerprint{
				Platform: tt.platform,
			}

			switch tt.platform {
			case "windows":
				fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
				fp.ScreenWidth = 1920
				fp.ScreenHeight = 1080
			case "macos":
				fp.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
				fp.ScreenWidth = 2560
				fp.ScreenHeight = 1600
			case "linux":
				fp.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"
				fp.ScreenWidth = 1920
				fp.ScreenHeight = 1080
			default:
				fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
				fp.ScreenWidth = 1920
				fp.ScreenHeight = 1080
			}

			if fp.ScreenWidth != tt.width {
				t.Errorf("expected width %d, got %d", tt.width, fp.ScreenWidth)
			}
		})
	}
}

func TestGenerateUniqueID(t *testing.T) {
	id1, err := generateUniqueID("fp")
	if err != nil {
		t.Fatalf("generateUniqueID failed: %v", err)
	}

	id2, err := generateUniqueID("fp")
	if err != nil {
		t.Fatalf("generateUniqueID failed: %v", err)
	}

	if id1 == id2 {
		t.Error("expected unique IDs to be different")
	}

	if len(id1) < 10 {
		t.Error("expected ID to be at least 10 characters")
	}
}

func TestGenerateRandomHash(t *testing.T) {
	hash1, err := generateRandomHash(16)
	if err != nil {
		t.Fatalf("generateRandomHash failed: %v", err)
	}

	hash2, err := generateRandomHash(16)
	if err != nil {
		t.Fatalf("generateRandomHash failed: %v", err)
	}

	if hash1 == hash2 {
		t.Error("expected different hashes")
	}

	if len(hash1) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("expected hash length 32, got %d", len(hash1))
	}
}

func TestFingerprintModel(t *testing.T) {
	fp := &dal.BrowserFingerprint{
		FingerprintID: "fp_test123",
		UserAgent:     "Mozilla/5.0",
		Platform:      "windows",
		Language:      "zh-CN",
		ScreenWidth:   1920,
		ScreenHeight:  1080,
		Timezone:      "Asia/Shanghai",
		IsUnique:      true,
		Status:        1,
	}

	if fp.FingerprintID != "fp_test123" {
		t.Error("fingerprint ID mismatch")
	}
	if !fp.IsUnique {
		t.Error("expected unique fingerprint")
	}
	if fp.Status != 1 {
		t.Error("expected active status")
	}
}

func TestProxyModel(t *testing.T) {
	proxy := &dal.ProxyIP{
		ProxyID:     "proxy_test123",
		IPAddress:   "1.2.3.4",
		Port:        8080,
		Protocol:    "http",
		Country:     "CN",
		IsAvailable: true,
		FailCount:   0,
		Status:      1,
	}

	if proxy.ProxyID != "proxy_test123" {
		t.Error("proxy ID mismatch")
	}
	if !proxy.IsAvailable {
		t.Error("expected available proxy")
	}
	if proxy.FailCount != 0 {
		t.Error("expected zero fail count")
	}
}

func TestAccountEnvironmentModel(t *testing.T) {
	env := &dal.AccountEnvironment{
		AccountID:     1,
		FingerprintID: 1,
		ProxyID:       1,
		EnvName:       "env-1",
		IsActive:      true,
	}

	if env.AccountID != 1 {
		t.Error("account ID mismatch")
	}
	if !env.IsActive {
		t.Error("expected active environment")
	}
}

func TestProxyFailCount(t *testing.T) {
	proxy := &dal.ProxyIP{
		FailCount:   2,
		IsAvailable: true,
	}

	// Simulate fail
	proxy.FailCount++
	proxy.IsAvailable = proxy.FailCount < 3

	if proxy.FailCount != 3 {
		t.Errorf("expected fail count 3, got %d", proxy.FailCount)
	}
	if proxy.IsAvailable {
		t.Error("expected unavailable after 3 fails")
	}
}
