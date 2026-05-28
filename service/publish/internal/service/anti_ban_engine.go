package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"opengeo/pkg/config"
)

// ==================== 防封发布引擎 ====================

// AntiBanEngine 防封引擎：模拟真人发布行为，防止平台封号
type AntiBanEngine struct {
	proxyPool     *ProxyPool
	fingerprintMgr *FingerprintManager
	behaviorSim   *BehaviorSimulator
	rateLimiter   *AccountRateLimiter
}

// NewAntiBanEngine 创建防封引擎
func NewAntiBanEngine() *AntiBanEngine {
	return &AntiBanEngine{
		proxyPool:      NewProxyPool(),
		fingerprintMgr: NewFingerprintManager(),
		behaviorSim:    NewBehaviorSimulator(),
		rateLimiter:    NewAccountRateLimiter(),
	}
}

// PublishContext 发布上下文：每次发布携带的环境信息
type PublishContext struct {
	AccountID   int64              `json:"account_id"`
	Platform    string             `json:"platform"`
	Proxy       *ProxyInfo         `json:"proxy"`
	Fingerprint *FingerprintInfo   `json:"fingerprint"`
	Delay       time.Duration      `json:"delay"`
	Headers     map[string]string  `json:"headers"`
	UserAgent   string             `json:"user_agent"`
}

// PreparePublish 准备发布环境（IP + 指纹 + 延迟）
func (e *AntiBanEngine) PreparePublish(ctx context.Context, accountID int64, platform string) (*PublishContext, error) {
	// 1. 检查账号发布频率限制
	if err := e.rateLimiter.CheckLimit(accountID); err != nil {
		return nil, err
	}

	// 2. 分配代理IP
	proxy, err := e.proxyPool.Allocate(platform)
	if err != nil {
		// 无可用代理时使用直连（降级）
		proxy = &ProxyInfo{Type: "direct"}
	}

	// 3. 选择浏览器指纹
	fingerprint := e.fingerprintMgr.Select(accountID, platform)

	// 4. 计算随机延迟（模拟真人操作间隔）
	delay := e.behaviorSim.CalculateDelay(platform)

	// 5. 构建请求头
	headers := e.buildHeaders(fingerprint, proxy)

	return &PublishContext{
		AccountID:   accountID,
		Platform:    platform,
		Proxy:       proxy,
		Fingerprint: fingerprint,
		Delay:       delay,
		Headers:     headers,
		UserAgent:   fingerprint.UserAgent,
	}, nil
}

// AfterPublish 发布后处理（更新统计、释放资源）
func (e *AntiBanEngine) AfterPublish(accountID int64, platform string, proxyID int64, success bool) {
	e.rateLimiter.Record(accountID)

	if !success {
		e.proxyPool.MarkFailed(proxyID)
	}
}

// buildHeaders 构建模拟真人请求头
func (e *AntiBanEngine) buildHeaders(fp *FingerprintInfo, proxy *ProxyInfo) map[string]string {
	headers := map[string]string{
		"User-Agent":                fp.UserAgent,
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language":           fp.Language,
		"Accept-Encoding":           "gzip, deflate, br",
		"Connection":                "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Cache-Control":             "max-age=0",
	}

	// 添加平台特定头
	switch fp.Platform {
	case "windows":
		headers["Sec-Ch-Ua"] = `"Chromium";v="120", "Google Chrome";v="120", "Not-A.Brand";v="99"`
		headers["Sec-Ch-Ua-Mobile"] = "?0"
		headers["Sec-Ch-Ua-Platform"] = `"Windows"`
	case "macos":
		headers["Sec-Ch-Ua"] = `"Chromium";v="120", "Google Chrome";v="120", "Not-A.Brand";v="99"`
		headers["Sec-Ch-Ua-Mobile"] = "?0"
		headers["Sec-Ch-Ua-Platform"] = `"macOS"`
	}

	return headers
}

// ==================== 代理IP池 ====================

// ProxyInfo 代理信息
type ProxyInfo struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"` // http, https, socks5
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Country   string `json:"country"`
	Latency   int    `json:"latency"`
	Available bool   `json:"available"`
}

// ProxyPool 代理IP池
type ProxyPool struct {
	mu       sync.RWMutex
	proxies  map[string][]*ProxyInfo // platform -> proxies
	usage    map[int64]time.Time     // proxy_id -> last_used
	failures map[int64]int           // proxy_id -> consecutive_failures
}

// NewProxyPool 创建代理池
func NewProxyPool() *ProxyPool {
	return &ProxyPool{
		proxies:  make(map[string][]*ProxyInfo),
		usage:    make(map[int64]time.Time),
		failures: make(map[int64]int),
	}
}

// AddProxy 添加代理到池
func (p *ProxyPool) AddProxy(platform string, proxy *ProxyInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.proxies[platform] = append(p.proxies[platform], proxy)
}

// Allocate 分配代理（轮询 + 延迟最小优先）
func (p *ProxyPool) Allocate(platform string) (*ProxyInfo, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg := config.GetConfig()
	proxies, ok := p.proxies[platform]
	if !ok || len(proxies) == 0 {
		// 尝试通用代理
		proxies = p.proxies["*"]
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxy available for platform: %s", platform)
	}

	// 选择最近未使用且可用的代理
	var best *ProxyInfo
	var bestTime time.Time

	for _, proxy := range proxies {
		if !proxy.Available {
			continue
		}
		// 检查连续失败次数
		if p.failures[proxy.ID] >= cfg.AntiBan.ProxyFailureThreshold {
			continue
		}

		lastUsed := p.usage[proxy.ID]
		if best == nil || lastUsed.Before(bestTime) {
			best = proxy
			bestTime = lastUsed
		}
	}

	if best == nil {
		return nil, fmt.Errorf("all proxies exhausted for platform: %s", platform)
	}

	p.usage[best.ID] = time.Now()
	return best, nil
}

// MarkFailed 标记代理失败
func (p *ProxyPool) MarkFailed(proxyID int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.failures[proxyID]++

	cfg := config.GetConfig()
	threshold := cfg.AntiBan.ProxyFailureThreshold
	if threshold <= 0 {
		threshold = 3
	}

	if p.failures[proxyID] >= threshold {
		for _, proxies := range p.proxies {
			for _, proxy := range proxies {
				if proxy.ID == proxyID {
					proxy.Available = false
					break
				}
			}
		}
	}
}

// LoadProxies 批量加载代理配置到池中
func (p *ProxyPool) LoadProxies(platform string, proxyConfigs []*ProxyInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, proxy := range proxyConfigs {
		proxy.Available = true
		p.proxies[platform] = append(p.proxies[platform], proxy)
	}
}

// Release 释放代理
func (p *ProxyPool) Release(proxyID int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.usage, proxyID)
}

// ResetFailures 重置失败计数
func (p *ProxyPool) ResetFailures(proxyID int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.failures, proxyID)
}

// ==================== 浏览器指纹管理 ====================

// FingerprintInfo 浏览器指纹信息
type FingerprintInfo struct {
	ID           int64  `json:"id"`
	Platform     string `json:"platform"`
	UserAgent    string `json:"user_agent"`
	Language     string `json:"language"`
	ScreenWidth  int    `json:"screen_width"`
	ScreenHeight int    `json:"screen_height"`
	Timezone     string `json:"timezone"`
	WebGLVendor  string `json:"webgl_vendor"`
	WebGLRenderer string `json:"webgl_renderer"`
	CanvasHash   string `json:"canvas_hash"`
	AudioHash    string `json:"audio_hash"`
}

// FingerprintManager 指纹管理器
type FingerprintManager struct {
	mu          sync.RWMutex
	fingerprints map[string][]*FingerprintInfo // platform -> fingerprints
	usage        map[int64]int64               // fingerprint_id -> account_id
}

// NewFingerprintManager 创建指纹管理器
func NewFingerprintManager() *FingerprintManager {
	mgr := &FingerprintManager{
		fingerprints: make(map[string][]*FingerprintInfo),
		usage:        make(map[int64]int64),
	}

	// 预置常用指纹
	mgr.initBuiltinFingerprints()
	return mgr
}

// initBuiltinFingerprints 初始化内置指纹
func (m *FingerprintManager) initBuiltinFingerprints() {
	// Windows Chrome 指纹组
	windowsFPS := []*FingerprintInfo{
		{ID: 1, Platform: "windows", UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", Language: "zh-CN,zh;q=0.9,en;q=0.8", ScreenWidth: 1920, ScreenHeight: 1080, Timezone: "Asia/Shanghai", WebGLVendor: "Google Inc. (NVIDIA)", WebGLRenderer: "ANGLE (NVIDIA, NVIDIA GeForce GTX 1080)", CanvasHash: "win_chrome_001", AudioHash: "win_audio_001"},
		{ID: 2, Platform: "windows", UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36", Language: "zh-CN,zh;q=0.9", ScreenWidth: 2560, ScreenHeight: 1440, Timezone: "Asia/Shanghai", WebGLVendor: "Google Inc. (Intel)", WebGLRenderer: "ANGLE (Intel, Intel(R) UHD Graphics 630)", CanvasHash: "win_chrome_002", AudioHash: "win_audio_002"},
		{ID: 3, Platform: "windows", UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36", Language: "zh-TW,zh;q=0.9,en;q=0.8", ScreenWidth: 1920, ScreenHeight: 1080, Timezone: "Asia/Taipei", WebGLVendor: "Google Inc. (AMD)", WebGLRenderer: "ANGLE (AMD, AMD Radeon RX 580)", CanvasHash: "win_chrome_003", AudioHash: "win_audio_003"},
	}
	m.fingerprints["windows"] = windowsFPS

	// macOS Chrome 指纹组
	macFPS := []*FingerprintInfo{
		{ID: 10, Platform: "macos", UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", Language: "zh-CN,zh;q=0.9,en;q=0.8", ScreenWidth: 2560, ScreenHeight: 1600, Timezone: "Asia/Shanghai", WebGLVendor: "Apple", WebGLRenderer: "Apple M1 Pro", CanvasHash: "mac_chrome_001", AudioHash: "mac_audio_001"},
		{ID: 11, Platform: "macos", UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36", Language: "en-US,en;q=0.9", ScreenWidth: 1440, ScreenHeight: 900, Timezone: "America/New_York", WebGLVendor: "Apple", WebGLRenderer: "Apple M2", CanvasHash: "mac_chrome_002", AudioHash: "mac_audio_002"},
	}
	m.fingerprints["macos"] = macFPS

	// Linux Chrome 指纹组
	linuxFPS := []*FingerprintInfo{
		{ID: 20, Platform: "linux", UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", Language: "zh-CN,zh;q=0.9", ScreenWidth: 1920, ScreenHeight: 1080, Timezone: "Asia/Shanghai", WebGLVendor: "Mesa", WebGLRenderer: "Mesa Intel(R) UHD Graphics 630", CanvasHash: "linux_chrome_001", AudioHash: "linux_audio_001"},
	}
	m.fingerprints["linux"] = linuxFPS
}

// Select 为账号选择指纹（同一平台同一账号固定使用同一指纹）
func (m *FingerprintManager) Select(accountID int64, platform string) *FingerprintInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已有绑定
	for fpID, boundAccountID := range m.usage {
		if boundAccountID == accountID {
			// 查找对应指纹
			for _, fps := range m.fingerprints {
				for _, fp := range fps {
					if fp.ID == fpID {
						return fp
					}
				}
			}
		}
	}

	// 新分配：轮询选择
	fps := m.fingerprints[platform]
	if len(fps) == 0 {
		fps = m.fingerprints["windows"] // 默认
	}

	// 选择使用次数最少的指纹
	selected := fps[0]
	minUsage := int(^uint(0) >> 1)
	for _, fp := range fps {
		usageCount := 0
		for fpID := range m.usage {
			if fpID == fp.ID {
				usageCount++
			}
		}
		if usageCount < minUsage {
			minUsage = usageCount
			selected = fp
		}
	}

	m.usage[selected.ID] = accountID
	return selected
}

// Release 释放指纹
func (m *FingerprintManager) Release(fingerprintID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.usage, fingerprintID)
}

// ==================== 行为模拟器 ====================

// BehaviorSimulator 行为模拟器：模拟真人操作模式
type BehaviorSimulator struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// NewBehaviorSimulator 创建行为模拟器
func NewBehaviorSimulator() *BehaviorSimulator {
	return &BehaviorSimulator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CalculateDelay 计算随机延迟（模拟真人操作间隔）
func (s *BehaviorSimulator) CalculateDelay(platform string) time.Duration {
	cfg := config.GetConfig()

	bounds, ok := cfg.AntiBan.PlatformDelays[platform]
	if !ok {
		bounds = cfg.AntiBan.DefaultDelay
	}

	s.mu.Lock()
	base := bounds[0] + s.rng.Intn(bounds[1]-bounds[0])
	offset := int(float64(base) * cfg.AntiBan.OffsetRatio * (s.rng.Float64()*2 - 1))
	s.mu.Unlock()

	delay := base + offset
	if delay < cfg.AntiBan.MinDelay {
		delay = cfg.AntiBan.MinDelay
	}

	return time.Duration(delay) * time.Second
}

// SimulateScrollBehavior 模拟滚动行为（返回滚动间隔序列）
func (s *BehaviorSimulator) SimulateScrollBehavior(scrollCount int) []time.Duration {
	cfg := config.GetConfig()
	intervals := make([]time.Duration, scrollCount)
	min := cfg.AntiBan.ScrollInterval[0]
	max := cfg.AntiBan.ScrollInterval[1]
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < scrollCount; i++ {
		intervals[i] = time.Duration(min+s.rng.Intn(max-min)) * time.Millisecond
	}
	return intervals
}

// GenerateTypingDelay 模拟打字延迟（每个字符的间隔）
func (s *BehaviorSimulator) GenerateTypingDelay(textLen int) []time.Duration {
	cfg := config.GetConfig()
	delays := make([]time.Duration, textLen)
	min := cfg.AntiBan.TypingDelay[0]
	max := cfg.AntiBan.TypingDelay[1]
	pauseMin := cfg.AntiBan.TypingPauseDelay[0]
	pauseMax := cfg.AntiBan.TypingPauseDelay[1]
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < textLen; i++ {
		delays[i] = time.Duration(min+s.rng.Intn(max-min)) * time.Millisecond
		if s.rng.Float64() < cfg.AntiBan.TypingPauseProbability {
			delays[i] += time.Duration(pauseMin+s.rng.Intn(pauseMax-pauseMin)) * time.Millisecond
		}
	}
	return delays
}

// ==================== 账号频率限制器 ====================

// AccountRateLimiter 账号频率限制器
type AccountRateLimiter struct {
	mu       sync.RWMutex
	records  map[int64][]time.Time // account_id -> publish_times
	limits   map[string]*RateLimit // platform -> limit config
}

// RateLimit 频率限制配置
type RateLimit struct {
	MaxPerHour  int `json:"max_per_hour"`
	MaxPerDay   int `json:"max_per_day"`
	MinInterval int `json:"min_interval"` // 最小间隔（秒）
}

// NewAccountRateLimiter 创建频率限制器
func NewAccountRateLimiter() *AccountRateLimiter {
	cfg := config.GetConfig()
	limiter := &AccountRateLimiter{
		records: make(map[int64][]time.Time),
		limits:  make(map[string]*RateLimit),
	}

	// 从配置加载各平台频率限制
	for platform, limitCfg := range cfg.AntiBan.RateLimits {
		limiter.limits[platform] = &RateLimit{
			MaxPerHour:  limitCfg.MaxPerHour,
			MaxPerDay:   limitCfg.MaxPerDay,
			MinInterval: limitCfg.MinInterval,
		}
	}

	return limiter
}

// CheckLimit 检查是否超过频率限制
func (r *AccountRateLimiter) CheckLimit(accountID int64) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cfg := config.GetConfig()
	records := r.records[accountID]
	now := time.Now()

	// 清理过期记录
	var validRecords []time.Time
	for _, t := range records {
		if now.Sub(t) < 24*time.Hour {
			validRecords = append(validRecords, t)
		}
	}

	// 检查每小时限制
	hourAgo := now.Add(-1 * time.Hour)
	hourCount := 0
	for _, t := range validRecords {
		if t.After(hourAgo) {
			hourCount++
		}
	}

	// 使用配置中的默认限制
	maxPerHour := cfg.AntiBan.DefaultRateLimit.MaxPerHour
	if hourCount >= maxPerHour {
		return fmt.Errorf("rate limit exceeded: %d requests in last hour (max: %d)", hourCount, maxPerHour)
	}

	// 检查最小间隔
	if len(validRecords) > 0 {
		lastPublish := validRecords[len(validRecords)-1]
		minInterval := time.Duration(cfg.AntiBan.DefaultRateLimit.MinInterval) * time.Second
		if now.Sub(lastPublish) < minInterval {
			return fmt.Errorf("too frequent: minimum interval %v required", minInterval)
		}
	}

	return nil
}

// Record 记录发布
func (r *AccountRateLimiter) Record(accountID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[accountID] = append(r.records[accountID], time.Now())
}

// SetLimit 设置平台频率限制
func (r *AccountRateLimiter) SetLimit(platform string, limit *RateLimit) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.limits[platform] = limit
}

// GetStats 获取账号发布统计
func (r *AccountRateLimiter) GetStats(accountID int64) map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	records := r.records[accountID]
	now := time.Now()

	hourAgo := now.Add(-1 * time.Hour)
	dayAgo := now.Add(-24 * time.Hour)

	hourCount := 0
	dayCount := 0
	for _, t := range records {
		if t.After(hourAgo) {
			hourCount++
		}
		if t.After(dayAgo) {
			dayCount++
		}
	}

	return map[string]interface{}{
		"account_id":      accountID,
		"publish_count_1h": hourCount,
		"publish_count_24h": dayCount,
		"total_published":  len(records),
	}
}
