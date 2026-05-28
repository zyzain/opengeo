package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/service/account/internal/dal"
)

// AccountGroupService 账号分组服务
type AccountGroupService struct {
	groupRepo   *dal.AccountGroupRepository
	accountRepo *dal.AccountRepository
}

// NewAccountGroupService 创建账号分组服务
func NewAccountGroupService(
	groupRepo *dal.AccountGroupRepository,
	accountRepo *dal.AccountRepository,
) *AccountGroupService {
	return &AccountGroupService{
		groupRepo:   groupRepo,
		accountRepo: accountRepo,
	}
}

// CreateGroup 创建分组
func (s *AccountGroupService) CreateGroup(ctx context.Context, userID int64, name, groupType, description string, parentID *int64) (*dal.AccountGroup, error) {
	// 验证分组类型
	if !isValidGroupType(groupType) {
		return nil, fmt.Errorf("invalid group type: %s", groupType)
	}

	group := &dal.AccountGroup{
		UserID:      userID,
		Name:        name,
		ParentID:    parentID,
		GroupType:   groupType,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

// GetGroup 获取分组
func (s *AccountGroupService) GetGroup(ctx context.Context, groupID int64) (*dal.AccountGroup, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	return group, nil
}

// UpdateGroup 更新分组
func (s *AccountGroupService) UpdateGroup(ctx context.Context, groupID int64, name, description string) (*dal.AccountGroup, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	if name != "" {
		group.Name = name
	}
	if description != "" {
		group.Description = description
	}
	group.UpdatedAt = time.Now()

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	return group, nil
}

// DeleteGroup 删除分组
func (s *AccountGroupService) DeleteGroup(ctx context.Context, groupID int64) error {
	if err := s.groupRepo.Delete(ctx, groupID); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

// ListGroups 列出分组
func (s *AccountGroupService) ListGroups(ctx context.Context, userID int64, groupType string, parentID *int64) ([]*dal.AccountGroup, error) {
	groups, err := s.groupRepo.List(ctx, userID, groupType, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	return groups, nil
}

// AddAccountToGroup 添加账号到分组
func (s *AccountGroupService) AddAccountToGroup(ctx context.Context, accountID, groupID int64) error {
	// 验证账号是否存在
	_, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// 验证分组是否存在
	_, err = s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	if err := s.groupRepo.AddAccountToGroup(ctx, accountID, groupID); err != nil {
		return fmt.Errorf("failed to add account to group: %w", err)
	}

	return nil
}

// RemoveAccountFromGroup 从分组移除账号
func (s *AccountGroupService) RemoveAccountFromGroup(ctx context.Context, accountID, groupID int64) error {
	if err := s.groupRepo.RemoveAccountFromGroup(ctx, accountID, groupID); err != nil {
		return fmt.Errorf("failed to remove account from group: %w", err)
	}

	return nil
}

// GetGroupAccounts 获取分组下的账号
func (s *AccountGroupService) GetGroupAccounts(ctx context.Context, groupID int64) ([]*dal.Account, error) {
	accounts, err := s.groupRepo.GetGroupAccounts(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group accounts: %w", err)
	}

	return accounts, nil
}

// GetAccountGroups 获取账号所属的分组
func (s *AccountGroupService) GetAccountGroups(ctx context.Context, accountID int64) ([]*dal.AccountGroup, error) {
	groups, err := s.groupRepo.GetAccountGroups(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account groups: %w", err)
	}

	return groups, nil
}

// GetGroupStats 获取分组统计
func (s *AccountGroupService) GetGroupStats(ctx context.Context, groupID int64) (*dal.GroupStats, error) {
	stats, err := s.groupRepo.GetGroupStats(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group stats: %w", err)
	}

	return stats, nil
}

// GetGroupHierarchy 获取分组层级结构
func (s *AccountGroupService) GetGroupHierarchy(ctx context.Context, userID int64) ([]*GroupTreeNode, error) {
	// 获取所有根分组
	rootGroups, err := s.groupRepo.List(ctx, userID, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list root groups: %w", err)
	}

	// 构建树形结构
	var nodes []*GroupTreeNode
	for _, group := range rootGroups {
		node := &GroupTreeNode{
			Group: group,
		}

		// 递归获取子分组
		children, err := s.buildGroupTree(ctx, group.ID)
		if err != nil {
			return nil, err
		}
		node.Children = children

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// buildGroupTree 递归构建分组树
func (s *AccountGroupService) buildGroupTree(ctx context.Context, parentID int64) ([]*GroupTreeNode, error) {
	children, err := s.groupRepo.List(ctx, 0, "", &parentID)
	if err != nil {
		return nil, err
	}

	var nodes []*GroupTreeNode
	for _, child := range children {
		node := &GroupTreeNode{
			Group: child,
		}

		grandChildren, err := s.buildGroupTree(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		node.Children = grandChildren

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// GroupTreeNode 分组树节点
type GroupTreeNode struct {
	Group    *dal.AccountGroup  `json:"group"`
	Children []*GroupTreeNode   `json:"children,omitempty"`
}

// isValidGroupType 验证分组类型
func isValidGroupType(groupType string) bool {
	validTypes := map[string]bool{
		dal.AccountGroupTypeAuthority:    true,
		dal.AccountGroupTypeProfessional: true,
		dal.AccountGroupTypeEcology:      true,
	}
	return validTypes[groupType]
}