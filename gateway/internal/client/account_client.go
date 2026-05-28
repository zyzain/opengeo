package client

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AccountClient struct {
	db *gorm.DB
}

func NewAccountClient(db *gorm.DB) *AccountClient {
	return &AccountClient{db: db}
}

func (c *AccountClient) CreateAccount(ctx context.Context, userID int64, platform, accountName, accountID string) (map[string]interface{}, error) {
	account := map[string]interface{}{
		"user_id":      userID,
		"platform":     platform,
		"account_name": accountName,
		"account_id":   accountID,
		"status":       1,
		"health_score": 100,
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}
	if err := c.db.WithContext(ctx).Table("accounts").Create(account).Error; err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return account, nil
}

func (c *AccountClient) GetAccount(ctx context.Context, id int64) (map[string]interface{}, error) {
	var account map[string]interface{}
	if err := c.db.WithContext(ctx).Table("accounts").Where("id = ?", id).First(&account).Error; err != nil {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

func (c *AccountClient) UpdateAccount(ctx context.Context, id int64, accountName string, status int32) (map[string]interface{}, error) {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if accountName != "" {
		updates["account_name"] = accountName
	}
	if status > 0 {
		updates["status"] = status
	}
	c.db.WithContext(ctx).Table("accounts").Where("id = ?", id).Updates(updates)
	var account map[string]interface{}
	c.db.WithContext(ctx).Table("accounts").Where("id = ?", id).First(&account)
	return account, nil
}

func (c *AccountClient) DeleteAccount(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("accounts").Where("id = ?", id).Delete(nil).Error
}

func (c *AccountClient) ListAccounts(ctx context.Context, userID int64, platform string, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("accounts")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *AccountClient) GetAccountHealth(ctx context.Context, id int64) (map[string]interface{}, error) {
	var health map[string]interface{}
	if err := c.db.WithContext(ctx).Table("account_health").Where("account_id = ?", id).Order("checked_at DESC").First(&health).Error; err != nil {
		return map[string]interface{}{
			"account_id":   id,
			"health_score": 100,
			"status":       "normal",
		}, nil
	}
	return health, nil
}

// 账号分组
func (c *AccountClient) CreateAccountGroup(ctx context.Context, userID int64, name, groupType, description string) (map[string]interface{}, error) {
	group := map[string]interface{}{
		"user_id":     userID,
		"name":        name,
		"group_type":  groupType,
		"description": description,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	if err := c.db.WithContext(ctx).Table("account_groups").Create(group).Error; err != nil {
		return nil, fmt.Errorf("create group: %w", err)
	}
	return group, nil
}

func (c *AccountClient) GetAccountGroup(ctx context.Context, id int64) (map[string]interface{}, error) {
	var group map[string]interface{}
	if err := c.db.WithContext(ctx).Table("account_groups").Where("id = ?", id).First(&group).Error; err != nil {
		return nil, fmt.Errorf("group not found")
	}
	return group, nil
}

func (c *AccountClient) UpdateAccountGroup(ctx context.Context, id int64, name, description string) (map[string]interface{}, error) {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	c.db.WithContext(ctx).Table("account_groups").Where("id = ?", id).Updates(updates)
	var group map[string]interface{}
	c.db.WithContext(ctx).Table("account_groups").Where("id = ?", id).First(&group)
	return group, nil
}

func (c *AccountClient) DeleteAccountGroup(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Table("account_groups").Where("id = ?", id).Delete(nil).Error
}

func (c *AccountClient) ListAccountGroups(ctx context.Context, userID int64, page, pageSize int) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var total int64

	query := c.db.WithContext(ctx).Table("account_groups")
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&items)

	return map[string]interface{}{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, nil
}

func (c *AccountClient) AddAccountToGroup(ctx context.Context, groupID, accountID int64) error {
	return c.db.WithContext(ctx).Table("account_group_relations").Create(map[string]interface{}{
		"account_id": accountID,
		"group_id":   groupID,
	}).Error
}

func (c *AccountClient) RemoveAccountFromGroup(ctx context.Context, groupID, accountID int64) error {
	return c.db.WithContext(ctx).Table("account_group_relations").Where("account_id = ? AND group_id = ?", accountID, groupID).Delete(nil).Error
}
