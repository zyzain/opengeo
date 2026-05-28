package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Tenant 租户模型
type Tenant struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	Domain    string    `json:"domain" gorm:"size:256"`
	Status    int32     `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User 用户模型（归属租户）
type User struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID    int64      `json:"tenant_id" gorm:"index;not null;default:0"`
	Username    string     `json:"username" gorm:"uniqueIndex:idx_tenant_username;size:64;not null"`
	Password    string     `json:"-" gorm:"size:256;not null"`
	Email       string     `json:"email" gorm:"size:128;index"`
	Status      int32      `json:"status" gorm:"default:1"`
	LastLoginAt time.Time  `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Role 角色模型（租户隔离）
type Role struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID    int64     `json:"tenant_id" gorm:"index;not null;default:0"`
	Name        string    `json:"name" gorm:"size:64;not null"`
	Description string    `json:"description" gorm:"size:256"`
	CreatedAt   time.Time `json:"created_at"`
}

// Permission 权限模型
type Permission struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:64;not null"`
	Description string    `json:"description" gorm:"size:256"`
	Resource    string    `json:"resource" gorm:"size:128"`
	Action      string    `json:"action" gorm:"size:64"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRole 用户角色关联
type UserRole struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey"`
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       int64 `gorm:"primaryKey"`
	PermissionID int64 `gorm:"primaryKey"`
}

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", username)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", email)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List 列出用户
func (r *UserRepository) List(ctx context.Context, page, pageSize int32, keyword string, tenantID int64) ([]*User, int32, error) {
	var users []*User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})

	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, int32(total), nil
}

// Account 账号模型
type Account struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64      `json:"user_id" gorm:"not null;index:idx_account_list"`
	Platform      string     `json:"platform" gorm:"size:64;not null;index:idx_account_list"`
	AccountName   string     `json:"account_name" gorm:"size:128;not null"`
	AccountIDStr  string     `json:"account_id_str" gorm:"size:128"`
	Status        int32      `json:"status" gorm:"default:1;index:idx_account_check"`
	HealthScore   float32    `json:"health_score" gorm:"default:100"`
	LastCheckTime *time.Time `json:"last_check_time" gorm:"index:idx_account_check"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AccountRepository 账号仓储
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository 创建账号仓储
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create 创建账号
func (r *AccountRepository) Create(ctx context.Context, account *Account) error {
	if err := r.db.WithContext(ctx).Create(account).Error; err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

// GetByID 根据ID获取账号
func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*Account, error) {
	var account Account
	if err := r.db.WithContext(ctx).First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &account, nil
}

// Update 更新账号
func (r *AccountRepository) Update(ctx context.Context, account *Account) error {
	if err := r.db.WithContext(ctx).Save(account).Error; err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}

// Delete 删除账号
func (r *AccountRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&Account{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}

// List 列出账号
func (r *AccountRepository) List(ctx context.Context, userID int64, platform string, page, pageSize int32) ([]*Account, int32, error) {
	var accounts []*Account
	var total int64

	query := r.db.WithContext(ctx).Model(&Account{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count accounts: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&accounts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, int32(total), nil
}

// AccountHealth 账号健康状态
type AccountHealth struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID       int64     `json:"account_id" gorm:"index;not null"`
	HealthScore     float32   `json:"health_score"`
	Status          string    `json:"status" gorm:"size:32"`
	CheckType       string    `json:"check_type" gorm:"size:32"`
	CheckDetails    string    `json:"check_details" gorm:"type:text"`
	AlertLevel      string    `json:"alert_level" gorm:"size:16"`
	AlertSent       bool      `json:"alert_sent" gorm:"default:false"`
	AlertChannels   string    `json:"alert_channels" gorm:"size:256"`
	PublishPaused   bool      `json:"publish_paused" gorm:"default:false"`
	PauseReason     string    `json:"pause_reason" gorm:"size:512"`
	CheckedAt       time.Time `json:"checked_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// AlertRecord 告警记录
type AlertRecord struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID   int64     `json:"account_id" gorm:"index;not null"`
	HealthID    int64     `json:"health_id" gorm:"index"`
	AlertType   string    `json:"alert_type" gorm:"size:32"`
	Channel     string    `json:"channel" gorm:"size:32"`
	Title       string    `json:"title" gorm:"size:256"`
	Content     string    `json:"content" gorm:"type:text"`
	Success     bool      `json:"success" gorm:"default:false"`
	ErrorMsg    string    `json:"error_msg" gorm:"size:512"`
	SentAt      time.Time `json:"sent_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// HealthCheckResult 健康检测结果
type HealthCheckResult struct {
	AccountID    int64              `json:"account_id"`
	HealthScore  float32            `json:"health_score"`
	Status       string             `json:"status"`
	AlertLevel   string             `json:"alert_level"`
	PublishPaused bool              `json:"publish_paused"`
	PauseReason  string             `json:"pause_reason"`
	Checks       []CheckItem        `json:"checks"`
}

// CheckItem 检查项
type CheckItem struct {
	Name      string  `json:"name"`
	Passed    bool    `json:"passed"`
	Score     float32 `json:"score"`
	Detail    string  `json:"detail"`
}

// GetHealth 获取账号健康状态
func (r *AccountRepository) GetHealth(ctx context.Context, accountID int64) (*AccountHealth, error) {
	var health AccountHealth
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountID).Order("checked_at DESC").First(&health).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &AccountHealth{
				AccountID:   accountID,
				HealthScore: 100,
				Status:      "normal",
				AlertLevel:  "none",
				CheckedAt:   time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get account health: %w", err)
	}
	return &health, nil
}

// SaveHealth 保存健康检测结果
func (r *AccountRepository) SaveHealth(ctx context.Context, health *AccountHealth) error {
	if err := r.db.WithContext(ctx).Create(health).Error; err != nil {
		return fmt.Errorf("failed to save health: %w", err)
	}
	return nil
}

// GetHealthHistory 获取健康历史
func (r *AccountRepository) GetHealthHistory(ctx context.Context, accountID int64, limit int) ([]*AccountHealth, error) {
	var history []*AccountHealth
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountID).
		Order("checked_at DESC").Limit(limit).Find(&history).Error; err != nil {
		return nil, fmt.Errorf("failed to get health history: %w", err)
	}
	return history, nil
}

// GetAccountsToCheck 获取需要检测的账号
func (r *AccountRepository) GetAccountsToCheck(ctx context.Context, intervalMinutes int, limit int) ([]*Account, error) {
	var accounts []*Account
	cutoff := time.Now().Add(-time.Duration(intervalMinutes) * time.Minute)
	if err := r.db.WithContext(ctx).Where("status = 1 AND (last_check_time IS NULL OR last_check_time < ?)", cutoff).
		Limit(limit).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get accounts to check: %w", err)
	}
	return accounts, nil
}

// UpdateAccountHealthScore 更新账号健康分数
func (r *AccountRepository) UpdateAccountHealthScore(ctx context.Context, accountID int64, score float32, paused bool) error {
	if err := r.db.WithContext(ctx).Model(&Account{}).Where("id = ?", accountID).
		Updates(map[string]interface{}{
			"health_score":    score,
			"last_check_time": time.Now(),
			"updated_at":      time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update account health score: %w", err)
	}
	return nil
}

// AlertRecordRepository 告警记录仓储
type AlertRecordRepository struct {
	db *gorm.DB
}

// NewAlertRecordRepository 创建告警记录仓储
func NewAlertRecordRepository(db *gorm.DB) *AlertRecordRepository {
	return &AlertRecordRepository{db: db}
}

// Create 创建告警记录
func (r *AlertRecordRepository) Create(ctx context.Context, record *AlertRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("failed to create alert record: %w", err)
	}
	return nil
}

// ListByAccountID 列出账号告警记录
func (r *AlertRecordRepository) ListByAccountID(ctx context.Context, accountID int64, page, pageSize int) ([]*AlertRecord, int32, error) {
	var records []*AlertRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&AlertRecord{}).Where("account_id = ?", accountID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count alert records: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sent_at DESC").Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list alert records: %w", err)
	}

	return records, int32(total), nil
}