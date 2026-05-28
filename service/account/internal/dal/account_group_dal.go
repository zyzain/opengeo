package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AccountGroup 账号分组模型
type AccountGroup struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int64     `json:"user_id" gorm:"index;not null"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	ParentID    *int64    `json:"parent_id" gorm:"index"`
	GroupType   string    `json:"group_type" gorm:"size:32"` // authority, professional, ecology
	Description string    `json:"description" gorm:"size:256"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AccountGroupType 分组类型常量
const (
	AccountGroupTypeAuthority    = "authority"    // 权威背书层
	AccountGroupTypeProfessional = "professional" // 专业认证层
	AccountGroupTypeEcology      = "ecology"      // 生态渗透层
)

// AccountGroupRelation 账号分组关联
type AccountGroupRelation struct {
	AccountID int64 `json:"account_id" gorm:"primaryKey"`
	GroupID   int64 `json:"group_id" gorm:"primaryKey"`
}

// AccountGroupRepository 账号分组仓储
type AccountGroupRepository struct {
	db *gorm.DB
}

// NewAccountGroupRepository 创建账号分组仓储
func NewAccountGroupRepository(db *gorm.DB) *AccountGroupRepository {
	return &AccountGroupRepository{db: db}
}

// Create 创建分组
func (r *AccountGroupRepository) Create(ctx context.Context, group *AccountGroup) error {
	if err := r.db.WithContext(ctx).Create(group).Error; err != nil {
		return fmt.Errorf("failed to create account group: %w", err)
	}
	return nil
}

// GetByID 根据ID获取分组
func (r *AccountGroupRepository) GetByID(ctx context.Context, id int64) (*AccountGroup, error) {
	var group AccountGroup
	if err := r.db.WithContext(ctx).First(&group, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account group not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get account group: %w", err)
	}
	return &group, nil
}

// Update 更新分组
func (r *AccountGroupRepository) Update(ctx context.Context, group *AccountGroup) error {
	if err := r.db.WithContext(ctx).Save(group).Error; err != nil {
		return fmt.Errorf("failed to update account group: %w", err)
	}
	return nil
}

// Delete 删除分组
func (r *AccountGroupRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&AccountGroup{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete account group: %w", err)
	}
	return nil
}

// List 列出分组
func (r *AccountGroupRepository) List(ctx context.Context, userID int64, groupType string, parentID *int64) ([]*AccountGroup, error) {
	var groups []*AccountGroup

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if groupType != "" {
		query = query.Where("group_type = ?", groupType)
	}
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if err := query.Order("created_at DESC").Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("failed to list account groups: %w", err)
	}

	return groups, nil
}

// AddAccountToGroup 添加账号到分组
func (r *AccountGroupRepository) AddAccountToGroup(ctx context.Context, accountID, groupID int64) error {
	relation := &AccountGroupRelation{
		AccountID: accountID,
		GroupID:   groupID,
	}
	if err := r.db.WithContext(ctx).Create(relation).Error; err != nil {
		return fmt.Errorf("failed to add account to group: %w", err)
	}
	return nil
}

// RemoveAccountFromGroup 从分组移除账号
func (r *AccountGroupRepository) RemoveAccountFromGroup(ctx context.Context, accountID, groupID int64) error {
	if err := r.db.WithContext(ctx).Where("account_id = ? AND group_id = ?", accountID, groupID).Delete(&AccountGroupRelation{}).Error; err != nil {
		return fmt.Errorf("failed to remove account from group: %w", err)
	}
	return nil
}

// GetGroupAccounts 获取分组下的账号
func (r *AccountGroupRepository) GetGroupAccounts(ctx context.Context, groupID int64) ([]*Account, error) {
	var accounts []*Account

	err := r.db.WithContext(ctx).
		Joins("JOIN account_group_relations ON accounts.id = account_group_relations.account_id").
		Where("account_group_relations.group_id = ?", groupID).
		Find(&accounts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get group accounts: %w", err)
	}

	return accounts, nil
}

// GetAccountGroups 获取账号所属的分组
func (r *AccountGroupRepository) GetAccountGroups(ctx context.Context, accountID int64) ([]*AccountGroup, error) {
	var groups []*AccountGroup

	err := r.db.WithContext(ctx).
		Joins("JOIN account_group_relations ON account_groups.id = account_group_relations.group_id").
		Where("account_group_relations.account_id = ?", accountID).
		Find(&groups).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get account groups: %w", err)
	}

	return groups, nil
}

// GetGroupStats 获取分组统计
func (r *AccountGroupRepository) GetGroupStats(ctx context.Context, groupID int64) (*GroupStats, error) {
	var stats GroupStats

	// 获取分组信息
	group, err := r.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	stats.GroupID = groupID
	stats.GroupName = group.Name
	stats.GroupType = group.GroupType

	// 统计账号数量
	var accountCount int64
	err = r.db.WithContext(ctx).
		Model(&AccountGroupRelation{}).
		Where("group_id = ?", groupID).
		Count(&accountCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count accounts: %w", err)
	}
	stats.AccountCount = int32(accountCount)

	// 统计健康分数平均值
	var avgHealth float32
	err = r.db.WithContext(ctx).
		Model(&Account{}).
		Select("AVG(health_score)").
		Joins("JOIN account_group_relations ON accounts.id = account_group_relations.account_id").
		Where("account_group_relations.group_id = ?", groupID).
		Scan(&avgHealth).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate avg health: %w", err)
	}
	stats.AvgHealthScore = avgHealth

	return &stats, nil
}

// GroupStats 分组统计
type GroupStats struct {
	GroupID        int64   `json:"group_id"`
	GroupName      string  `json:"group_name"`
	GroupType      string  `json:"group_type"`
	AccountCount   int32   `json:"account_count"`
	AvgHealthScore float32 `json:"avg_health_score"`
}