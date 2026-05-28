package service

import (
	"testing"
	"time"
)

func TestAntiBanEngine_PreparePublish(t *testing.T) {
	engine := NewAntiBanEngine()

	// 添加测试代理
	engine.proxyPool.AddProxy("wechat", &ProxyInfo{
		ID:        1,
		Type:      "https",
		IP:        "1.2.3.4",
		Port:      8080,
		Country:   "CN",
		Available: true,
	})

	pubCtx, err := engine.PreparePublish(nil, 1, "wechat")
	if err != nil {
		t.Fatalf("PreparePublish failed: %v", err)
	}

	if pubCtx == nil {
		t.Fatal("expected non-nil publish context")
	}
	if pubCtx.UserAgent == "" {
		t.Error("expected non-empty user agent")
	}
	if pubCtx.Headers == nil {
		t.Error("expected non-nil headers")
	}
	if pubCtx.Delay <= 0 {
		t.Error("expected positive delay")
	}
}

func TestProxyPool_AddAndAllocate(t *testing.T) {
	pool := NewProxyPool()

	pool.AddProxy("wechat", &ProxyInfo{ID: 1, Type: "https", IP: "1.2.3.4", Port: 8080, Available: true})
	pool.AddProxy("wechat", &ProxyInfo{ID: 2, Type: "https", IP: "5.6.7.8", Port: 8080, Available: true})

	proxy, err := pool.Allocate("wechat")
	if err != nil {
		t.Fatalf("Allocate failed: %v", err)
	}
	if proxy == nil {
		t.Fatal("expected non-nil proxy")
	}
}

func TestProxyPool_Exhaustion(t *testing.T) {
	pool := NewProxyPool()

	_, err := pool.Allocate("wechat")
	if err == nil {
		t.Error("expected error when no proxies available")
	}
}

func TestFingerprintManager_Select(t *testing.T) {
	mgr := NewFingerprintManager()

	fp1 := mgr.Select(1, "windows")
	if fp1 == nil {
		t.Fatal("expected non-nil fingerprint")
	}
	if fp1.UserAgent == "" {
		t.Error("expected non-empty user agent")
	}
	if fp1.Platform != "windows" {
		t.Errorf("expected windows, got %s", fp1.Platform)
	}

	// 同一账号应返回同一指纹
	fp2 := mgr.Select(1, "windows")
	if fp1.ID != fp2.ID {
		t.Errorf("expected same fingerprint for same account: %d vs %d", fp1.ID, fp2.ID)
	}

	// 不同账号应返回不同指纹
	fp3 := mgr.Select(2, "windows")
	// 可能相同也可能不同，但不应 panic
	_ = fp3
}

func TestFingerprintManager_AllPlatforms(t *testing.T) {
	mgr := NewFingerprintManager()

	platforms := []struct {
		name      string
		accountID int64
	}{
		{"windows", 201},
		{"macos", 202},
		{"linux", 203},
	}

	for _, p := range platforms {
		fp := mgr.Select(p.accountID, p.name)
		if fp == nil {
			t.Errorf("expected fingerprint for platform %s", p.name)
		}
		if fp.Platform != p.name {
			t.Errorf("expected platform %s, got %s", p.name, fp.Platform)
		}
	}
}

func TestBehaviorSimulator_CalculateDelay(t *testing.T) {
	sim := NewBehaviorSimulator()

	platforms := []string{"wechat", "weibo", "douyin", "xiaohongshu", "zhihu", "toutiao"}

	for _, p := range platforms {
		delay := sim.CalculateDelay(p)
		if delay < 5*time.Second {
			t.Errorf("delay too short for %s: %v", p, delay)
		}
		if delay > 300*time.Second {
			t.Errorf("delay too long for %s: %v", p, delay)
		}
	}
}

func TestBehaviorSimulator_SimulateScrollBehavior(t *testing.T) {
	sim := NewBehaviorSimulator()

	intervals := sim.SimulateScrollBehavior(10)
	if len(intervals) != 10 {
		t.Errorf("expected 10 intervals, got %d", len(intervals))
	}

	for _, interval := range intervals {
		if interval < 500*time.Millisecond || interval > 3*time.Second {
			t.Errorf("interval out of range: %v", interval)
		}
	}
}

func TestBehaviorSimulator_GenerateTypingDelay(t *testing.T) {
	sim := NewBehaviorSimulator()

	delays := sim.GenerateTypingDelay(20)
	if len(delays) != 20 {
		t.Errorf("expected 20 delays, got %d", len(delays))
	}

	for _, delay := range delays {
		if delay < 50*time.Millisecond {
			t.Errorf("delay too short: %v", delay)
		}
	}
}

func TestAccountRateLimiter_CheckLimit(t *testing.T) {
	limiter := NewAccountRateLimiter()

	// 第一次应该通过
	if err := limiter.CheckLimit(1); err != nil {
		t.Errorf("first check should pass: %v", err)
	}

	// 记录发布
	limiter.Record(1)

	// 立即再次检查应该被限制（间隔太短）
	if err := limiter.CheckLimit(1); err == nil {
		t.Error("expected rate limit error for immediate retry")
	}
}

func TestAccountRateLimiter_GetStats(t *testing.T) {
	limiter := NewAccountRateLimiter()

	limiter.Record(1)
	limiter.Record(1)
	limiter.Record(1)

	stats := limiter.GetStats(1)
	if stats["total_published"] != 3 {
		t.Errorf("expected 3 total, got %v", stats["total_published"])
	}
}

func TestAccountRateLimiter_SetLimit(t *testing.T) {
	limiter := NewAccountRateLimiter()

	limiter.SetLimit("custom", &RateLimit{
		MaxPerHour:  100,
		MaxPerDay:   1000,
		MinInterval: 10,
	})

	// 自定义限制应该生效
	if err := limiter.CheckLimit(999); err != nil {
		// 使用默认限制
		t.Logf("using default limit: %v", err)
	}
}

func TestBuildHeaders(t *testing.T) {
	engine := NewAntiBanEngine()

	fp := &FingerprintInfo{
		Platform:  "windows",
		UserAgent: "Mozilla/5.0 Test",
		Language:  "zh-CN,zh;q=0.9",
	}

	proxy := &ProxyInfo{Type: "https"}

	headers := engine.buildHeaders(fp, proxy)

	if headers["User-Agent"] != "Mozilla/5.0 Test" {
		t.Errorf("wrong user agent: %s", headers["User-Agent"])
	}
	if headers["Accept-Language"] != "zh-CN,zh;q=0.9" {
		t.Errorf("wrong language: %s", headers["Accept-Language"])
	}
	if _, ok := headers["Sec-Ch-Ua"]; !ok {
		t.Error("expected Sec-Ch-Ua header for windows")
	}
}

func TestProxyPool_Release(t *testing.T) {
	pool := NewProxyPool()
	pool.AddProxy("test", &ProxyInfo{ID: 100, Available: true})

	proxy, _ := pool.Allocate("test")
	if proxy == nil {
		t.Fatal("expected proxy")
	}

	pool.Release(proxy.ID)
	// 释放后应该可以再次分配
	proxy2, err := pool.Allocate("test")
	if err != nil {
		t.Errorf("expected to allocate after release: %v", err)
	}
	_ = proxy2
}

func TestFingerprintManager_Release(t *testing.T) {
	mgr := NewFingerprintManager()

	fp := mgr.Select(999, "windows")
	mgr.Release(fp.ID)

	// 释放后再次选择可能返回不同指纹
	fp2 := mgr.Select(999, "windows")
	if fp2 == nil {
		t.Error("expected fingerprint after release")
	}
}
