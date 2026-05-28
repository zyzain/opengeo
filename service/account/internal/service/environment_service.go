package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"opengeo/service/account/internal/dal"
)

// EnvironmentService 环境隔离服务
type EnvironmentService struct {
	envRepo *dal.EnvironmentRepository
}

// NewEnvironmentService 创建环境隔离服务
func NewEnvironmentService(envRepo *dal.EnvironmentRepository) *EnvironmentService {
	return &EnvironmentService{envRepo: envRepo}
}

// ==================== 浏览器指纹管理 ====================

// GenerateFingerprint 生成浏览器指纹
func (s *EnvironmentService) GenerateFingerprint(ctx context.Context, platform string) (*dal.BrowserFingerprint, error) {
	fingerprintID, err := generateUniqueID("fp")
	if err != nil {
		return nil, fmt.Errorf("generate fingerprint id: %w", err)
	}

	canvasHash, err := generateRandomHash(16)
	if err != nil {
		return nil, fmt.Errorf("generate canvas hash: %w", err)
	}

	audioHash, err := generateRandomHash(16)
	if err != nil {
		return nil, fmt.Errorf("generate audio hash: %w", err)
	}

	fp := &dal.BrowserFingerprint{
		FingerprintID: fingerprintID,
		Platform:      platform,
		IsUnique:      true,
		Status:        1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 根据平台生成不同的默认值
	switch platform {
	case "windows":
		fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		fp.ScreenWidth = 1920
		fp.ScreenHeight = 1080
		fp.Timezone = "Asia/Shanghai"
		fp.WebGLVendor = "Google Inc. (NVIDIA)"
		fp.WebGLRenderer = "ANGLE (NVIDIA, NVIDIA GeForce GTX 1080 Direct3D11 vs_5_0 ps_5_0, D3D11)"
	case "macos":
		fp.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		fp.ScreenWidth = 2560
		fp.ScreenHeight = 1600
		fp.Timezone = "Asia/Shanghai"
		fp.WebGLVendor = "Apple"
		fp.WebGLRenderer = "Apple M1 Pro"
	case "linux":
		fp.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		fp.ScreenWidth = 1920
		fp.ScreenHeight = 1080
		fp.Timezone = "Asia/Shanghai"
		fp.WebGLVendor = "Mesa"
		fp.WebGLRenderer = "Mesa Intel(R) UHD Graphics 630"
	default:
		fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
		fp.ScreenWidth = 1920
		fp.ScreenHeight = 1080
		fp.Timezone = "Asia/Shanghai"
	}

	fp.Language = "zh-CN,zh;q=0.9,en;q=0.8"
	fp.CanvasHash = canvasHash
	fp.AudioHash = audioHash

	// 检查唯一性
	isUnique, err := s.envRepo.CheckFingerprintUniqueness(ctx, canvasHash, audioHash)
	if err != nil {
		return nil, fmt.Errorf("check uniqueness: %w", err)
	}
	fp.IsUnique = isUnique

	if err := s.envRepo.CreateFingerprint(ctx, fp); err != nil {
		return nil, fmt.Errorf("create fingerprint: %w", err)
	}

	return fp, nil
}

// GetFingerprint 获取指纹
func (s *EnvironmentService) GetFingerprint(ctx context.Context, id int64) (*dal.BrowserFingerprint, error) {
	return s.envRepo.GetFingerprintByID(ctx, id)
}

// ListFingerprints 列出指纹
func (s *EnvironmentService) ListFingerprints(ctx context.Context, platform string, isUnique *bool, page, pageSize int) ([]*dal.BrowserFingerprint, int32, error) {
	return s.envRepo.ListFingerprints(ctx, platform, isUnique, page, pageSize)
}

// DeleteFingerprint 删除指纹
func (s *EnvironmentService) DeleteFingerprint(ctx context.Context, id int64) error {
	return s.envRepo.DeleteFingerprint(ctx, id)
}

// CheckFingerprintUniqueness 检查指纹唯一性
func (s *EnvironmentService) CheckFingerprintUniqueness(ctx context.Context, canvasHash, audioHash string) (bool, error) {
	return s.envRepo.CheckFingerprintUniqueness(ctx, canvasHash, audioHash)
}

// ==================== 代理IP管理 ====================

// AddProxy 添加代理IP
func (s *EnvironmentService) AddProxy(ctx context.Context, ipAddress string, port int32, protocol, username, password, country, city, isp string) (*dal.ProxyIP, error) {
	proxyID, err := generateUniqueID("proxy")
	if err != nil {
		return nil, fmt.Errorf("generate proxy id: %w", err)
	}

	proxy := &dal.ProxyIP{
		ProxyID:     proxyID,
		IPAddress:   ipAddress,
		Port:        port,
		Protocol:    protocol,
		Username:    username,
		Password:    password,
		Country:     country,
		City:        city,
		ISP:         isp,
		IsAvailable: true,
		FailCount:   0,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.envRepo.CreateProxy(ctx, proxy); err != nil {
		return nil, fmt.Errorf("create proxy: %w", err)
	}

	return proxy, nil
}

// GetProxy 获取代理
func (s *EnvironmentService) GetProxy(ctx context.Context, id int64) (*dal.ProxyIP, error) {
	return s.envRepo.GetProxyByID(ctx, id)
}

// ListProxies 列出代理
func (s *EnvironmentService) ListProxies(ctx context.Context, country, protocol string, isAvailable *bool, page, pageSize int) ([]*dal.ProxyIP, int32, error) {
	return s.envRepo.ListProxies(ctx, country, protocol, isAvailable, page, pageSize)
}

// GetAvailableProxy 获取可用代理
func (s *EnvironmentService) GetAvailableProxy(ctx context.Context, country string) (*dal.ProxyIP, error) {
	return s.envRepo.GetAvailableProxy(ctx, country)
}

// MarkProxyFailed 标记代理失败
func (s *EnvironmentService) MarkProxyFailed(ctx context.Context, proxyID int64) error {
	return s.envRepo.MarkProxyFailed(ctx, proxyID)
}

// DeleteProxy 删除代理
func (s *EnvironmentService) DeleteProxy(ctx context.Context, id int64) error {
	return s.envRepo.DeleteProxy(ctx, id)
}

// TestProxy 测试代理可用性
func (s *EnvironmentService) TestProxy(ctx context.Context, proxyID int64) (bool, int32, error) {
	proxy, err := s.envRepo.GetProxyByID(ctx, proxyID)
	if err != nil {
		return false, 0, err
	}

	// TODO: 实现实际的代理测试逻辑
	// 这里返回模拟结果
	_ = proxy

	return true, 50, nil // 可用，延迟50ms
}

// ==================== 环境绑定管理 ====================

// BindAccountEnvironment 绑定账号环境
func (s *EnvironmentService) BindAccountEnvironment(ctx context.Context, accountID, fingerprintID, proxyID int64, envName string) (*dal.AccountEnvironment, error) {
	// 检查账号是否已有绑定
	existing, _ := s.envRepo.GetEnvironmentByAccountID(ctx, accountID)
	if existing != nil {
		// 解绑旧环境
		if err := s.envRepo.UnbindEnvironment(ctx, accountID); err != nil {
			return nil, fmt.Errorf("unbind old environment: %w", err)
		}
	}

	env := &dal.AccountEnvironment{
		AccountID:     accountID,
		FingerprintID: fingerprintID,
		ProxyID:       proxyID,
		EnvName:       envName,
		IsActive:      true,
		LastUsedAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.envRepo.BindEnvironment(ctx, env); err != nil {
		return nil, fmt.Errorf("bind environment: %w", err)
	}

	return env, nil
}

// GetAccountEnvironment 获取账号环境
func (s *EnvironmentService) GetAccountEnvironment(ctx context.Context, accountID int64) (*dal.AccountEnvironment, error) {
	return s.envRepo.GetEnvironmentByAccountID(ctx, accountID)
}

// ListAccountEnvironments 列出账号环境
func (s *EnvironmentService) ListAccountEnvironments(ctx context.Context, accountID int64, isActive *bool, page, pageSize int) ([]*dal.AccountEnvironment, int32, error) {
	return s.envRepo.ListEnvironments(ctx, accountID, isActive, page, pageSize)
}

// UnbindAccountEnvironment 解绑账号环境
func (s *EnvironmentService) UnbindAccountEnvironment(ctx context.Context, accountID int64) error {
	return s.envRepo.UnbindEnvironment(ctx, accountID)
}

// UpdateEnvironmentUsage 更新环境使用时间
func (s *EnvironmentService) UpdateEnvironmentUsage(ctx context.Context, accountID int64) error {
	return s.envRepo.UpdateEnvironmentUsage(ctx, accountID)
}

// AllocateEnvironment 自动分配环境（指纹+代理）
func (s *EnvironmentService) AllocateEnvironment(ctx context.Context, accountID int64, platform, country string) (*dal.AccountEnvironment, error) {
	// 生成新指纹
	fp, err := s.GenerateFingerprint(ctx, platform)
	if err != nil {
		return nil, fmt.Errorf("generate fingerprint: %w", err)
	}

	// 获取可用代理
	proxy, err := s.GetAvailableProxy(ctx, country)
	if err != nil {
		// 没有可用代理，只绑定指纹
		env := &dal.AccountEnvironment{
			AccountID:     accountID,
			FingerprintID: fp.ID,
			ProxyID:       0,
			EnvName:       fmt.Sprintf("env-%d", accountID),
			IsActive:      true,
			LastUsedAt:    time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := s.envRepo.BindEnvironment(ctx, env); err != nil {
			return nil, fmt.Errorf("bind environment: %w", err)
		}
		return env, nil
	}

	// 绑定指纹+代理
	env := &dal.AccountEnvironment{
		AccountID:     accountID,
		FingerprintID: fp.ID,
		ProxyID:       proxy.ID,
		EnvName:       fmt.Sprintf("env-%d", accountID),
		IsActive:      true,
		LastUsedAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.envRepo.BindEnvironment(ctx, env); err != nil {
		return nil, fmt.Errorf("bind environment: %w", err)
	}

	return env, nil
}

// ==================== 辅助函数 ====================

// generateUniqueID 生成唯一ID
func generateUniqueID(prefix string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return prefix + "_" + hex.EncodeToString(bytes), nil
}

// generateRandomHash 生成随机哈希
func generateRandomHash(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
