package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// BrowserFingerprint 浏览器指纹
type BrowserFingerprint struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	FingerprintID string    `json:"fingerprint_id" gorm:"uniqueIndex;size:64;not null"`
	UserAgent     string    `json:"user_agent" gorm:"size:512"`
	Platform      string    `json:"platform" gorm:"size:32"`
	Language      string    `json:"language" gorm:"size:32"`
	ScreenWidth   int32     `json:"screen_width"`
	ScreenHeight  int32     `json:"screen_height"`
	Timezone      string    `json:"timezone" gorm:"size:64"`
	WebGLVendor   string    `json:"webgl_vendor" gorm:"size:128"`
	WebGLRenderer string    `json:"webgl_renderer" gorm:"size:256"`
	Fonts         string    `json:"fonts" gorm:"type:text"`
	Plugins       string    `json:"plugins" gorm:"type:text"`
	CanvasHash    string    `json:"canvas_hash" gorm:"size:128"`
	AudioHash     string    `json:"audio_hash" gorm:"size:128"`
	IsUnique      bool      `json:"is_unique" gorm:"default:true"`
	Status        int32     `json:"status" gorm:"default:1"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ProxyIP 代理IP
type ProxyIP struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProxyID     string     `json:"proxy_id" gorm:"uniqueIndex;size:64;not null"`
	IPAddress   string     `json:"ip_address" gorm:"size:64;not null"`
	Port        int32      `json:"port"`
	Protocol    string     `json:"protocol" gorm:"size:16"`
	Username    string     `json:"username" gorm:"size:64"`
	Password    string     `json:"password" gorm:"size:128"`
	Country     string     `json:"country" gorm:"size:64"`
	City        string     `json:"city" gorm:"size:64"`
	ISP         string     `json:"isp" gorm:"size:64"`
	Latency     int32      `json:"latency"`
	IsAvailable bool       `json:"is_available" gorm:"default:true"`
	FailCount   int32      `json:"fail_count" gorm:"default:0"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Status      int32      `json:"status" gorm:"default:1"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AccountEnvironment 账号环境绑定
type AccountEnvironment struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID     int64     `json:"account_id" gorm:"index;not null"`
	FingerprintID int64     `json:"fingerprint_id" gorm:"index;not null"`
	ProxyID       int64     `json:"proxy_id" gorm:"index"`
	EnvName       string    `json:"env_name" gorm:"size:128"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	LastUsedAt    time.Time `json:"last_used_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// EnvironmentRepository 环境仓储
type EnvironmentRepository struct {
	db *gorm.DB
}

// NewEnvironmentRepository 创建环境仓储
func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

// ==================== 浏览器指纹 ====================

// CreateFingerprint 创建浏览器指纹
func (r *EnvironmentRepository) CreateFingerprint(ctx context.Context, fp *BrowserFingerprint) error {
	if err := r.db.WithContext(ctx).Create(fp).Error; err != nil {
		return fmt.Errorf("failed to create fingerprint: %w", err)
	}
	return nil
}

// GetFingerprintByID 根据ID获取指纹
func (r *EnvironmentRepository) GetFingerprintByID(ctx context.Context, id int64) (*BrowserFingerprint, error) {
	var fp BrowserFingerprint
	if err := r.db.WithContext(ctx).First(&fp, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("fingerprint not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get fingerprint: %w", err)
	}
	return &fp, nil
}

// GetFingerprintByFingerprintID 根据指纹ID获取
func (r *EnvironmentRepository) GetFingerprintByFingerprintID(ctx context.Context, fingerprintID string) (*BrowserFingerprint, error) {
	var fp BrowserFingerprint
	if err := r.db.WithContext(ctx).Where("fingerprint_id = ?", fingerprintID).First(&fp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("fingerprint not found: %s", fingerprintID)
		}
		return nil, fmt.Errorf("failed to get fingerprint: %w", err)
	}
	return &fp, nil
}

// ListFingerprints 列出指纹
func (r *EnvironmentRepository) ListFingerprints(ctx context.Context, platform string, isUnique *bool, page, pageSize int) ([]*BrowserFingerprint, int32, error) {
	var fps []*BrowserFingerprint
	var total int64

	query := r.db.WithContext(ctx).Model(&BrowserFingerprint{})
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if isUnique != nil {
		query = query.Where("is_unique = ?", *isUnique)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count fingerprints: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&fps).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list fingerprints: %w", err)
	}

	return fps, int32(total), nil
}

// UpdateFingerprint 更新指纹
func (r *EnvironmentRepository) UpdateFingerprint(ctx context.Context, fp *BrowserFingerprint) error {
	if err := r.db.WithContext(ctx).Save(fp).Error; err != nil {
		return fmt.Errorf("failed to update fingerprint: %w", err)
	}
	return nil
}

// DeleteFingerprint 删除指纹
func (r *EnvironmentRepository) DeleteFingerprint(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&BrowserFingerprint{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete fingerprint: %w", err)
	}
	return nil
}

// CheckFingerprintUniqueness 检查指纹唯一性
func (r *EnvironmentRepository) CheckFingerprintUniqueness(ctx context.Context, canvasHash, audioHash string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&BrowserFingerprint{}).
		Where("canvas_hash = ? OR audio_hash = ?", canvasHash, audioHash).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check uniqueness: %w", err)
	}
	return count == 0, nil
}

// ==================== 代理IP ====================

// CreateProxy 创建代理IP
func (r *EnvironmentRepository) CreateProxy(ctx context.Context, proxy *ProxyIP) error {
	if err := r.db.WithContext(ctx).Create(proxy).Error; err != nil {
		return fmt.Errorf("failed to create proxy: %w", err)
	}
	return nil
}

// GetProxyByID 根据ID获取代理
func (r *EnvironmentRepository) GetProxyByID(ctx context.Context, id int64) (*ProxyIP, error) {
	var proxy ProxyIP
	if err := r.db.WithContext(ctx).First(&proxy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("proxy not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get proxy: %w", err)
	}
	return &proxy, nil
}

// ListProxies 列出代理
func (r *EnvironmentRepository) ListProxies(ctx context.Context, country, protocol string, isAvailable *bool, page, pageSize int) ([]*ProxyIP, int32, error) {
	var proxies []*ProxyIP
	var total int64

	query := r.db.WithContext(ctx).Model(&ProxyIP{})
	if country != "" {
		query = query.Where("country = ?", country)
	}
	if protocol != "" {
		query = query.Where("protocol = ?", protocol)
	}
	if isAvailable != nil {
		query = query.Where("is_available = ?", *isAvailable)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count proxies: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("latency ASC").Find(&proxies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list proxies: %w", err)
	}

	return proxies, int32(total), nil
}

// GetAvailableProxy 获取可用代理（延迟最低）
func (r *EnvironmentRepository) GetAvailableProxy(ctx context.Context, country string) (*ProxyIP, error) {
	var proxy ProxyIP
	query := r.db.WithContext(ctx).Where("is_available = ? AND fail_count < 3", true)
	if country != "" {
		query = query.Where("country = ?", country)
	}
	if err := query.Order("latency ASC").First(&proxy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no available proxy")
		}
		return nil, fmt.Errorf("failed to get available proxy: %w", err)
	}
	return &proxy, nil
}

// UpdateProxy 更新代理
func (r *EnvironmentRepository) UpdateProxy(ctx context.Context, proxy *ProxyIP) error {
	if err := r.db.WithContext(ctx).Save(proxy).Error; err != nil {
		return fmt.Errorf("failed to update proxy: %w", err)
	}
	return nil
}

// MarkProxyFailed 标记代理失败
func (r *EnvironmentRepository) MarkProxyFailed(ctx context.Context, proxyID int64) error {
	if err := r.db.WithContext(ctx).Model(&ProxyIP{}).
		Where("id = ?", proxyID).
		Updates(map[string]interface{}{
			"fail_count":   gorm.Expr("fail_count + 1"),
			"is_available": gorm.Expr("fail_count + 1 < 3"),
			"updated_at":   time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to mark proxy failed: %w", err)
	}
	return nil
}

// DeleteProxy 删除代理
func (r *EnvironmentRepository) DeleteProxy(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&ProxyIP{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete proxy: %w", err)
	}
	return nil
}

// ==================== 环境绑定 ====================

// BindEnvironment 绑定账号环境
func (r *EnvironmentRepository) BindEnvironment(ctx context.Context, env *AccountEnvironment) error {
	if err := r.db.WithContext(ctx).Create(env).Error; err != nil {
		return fmt.Errorf("failed to bind environment: %w", err)
	}
	return nil
}

// GetEnvironmentByAccountID 根据账号ID获取环境
func (r *EnvironmentRepository) GetEnvironmentByAccountID(ctx context.Context, accountID int64) (*AccountEnvironment, error) {
	var env AccountEnvironment
	if err := r.db.WithContext(ctx).Where("account_id = ? AND is_active = ?", accountID, true).First(&env).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("environment not found for account: %d", accountID)
		}
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}
	return &env, nil
}

// ListEnvironments 列出环境绑定
func (r *EnvironmentRepository) ListEnvironments(ctx context.Context, accountID int64, isActive *bool, page, pageSize int) ([]*AccountEnvironment, int32, error) {
	var envs []*AccountEnvironment
	var total int64

	query := r.db.WithContext(ctx).Model(&AccountEnvironment{})
	if accountID > 0 {
		query = query.Where("account_id = ?", accountID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count environments: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&envs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list environments: %w", err)
	}

	return envs, int32(total), nil
}

// UnbindEnvironment 解绑账号环境
func (r *EnvironmentRepository) UnbindEnvironment(ctx context.Context, accountID int64) error {
	if err := r.db.WithContext(ctx).Model(&AccountEnvironment{}).
		Where("account_id = ?", accountID).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to unbind environment: %w", err)
	}
	return nil
}

// UpdateEnvironmentUsage 更新环境使用时间
func (r *EnvironmentRepository) UpdateEnvironmentUsage(ctx context.Context, accountID int64) error {
	if err := r.db.WithContext(ctx).Model(&AccountEnvironment{}).
		Where("account_id = ? AND is_active = ?", accountID, true).
		Update("last_used_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to update environment usage: %w", err)
	}
	return nil
}
