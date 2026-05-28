package service

import (
	"context"
	"fmt"
	"time"

	"opengeo/service/account/internal/dal"
)

// RBACService 角色权限服务
type RBACService struct {
	tenantRepo     *dal.TenantRepository
	roleRepo       *dal.RoleRepository
	permRepo       *dal.PermissionRepository
	userRoleRepo   *dal.UserRoleRepository
	rolePermRepo   *dal.RolePermissionRepository
}

func NewRBACService(
	tenantRepo *dal.TenantRepository,
	roleRepo *dal.RoleRepository,
	permRepo *dal.PermissionRepository,
	userRoleRepo *dal.UserRoleRepository,
	rolePermRepo *dal.RolePermissionRepository,
) *RBACService {
	return &RBACService{
		tenantRepo:   tenantRepo,
		roleRepo:     roleRepo,
		permRepo:     permRepo,
		userRoleRepo: userRoleRepo,
		rolePermRepo: rolePermRepo,
	}
}

// ==================== Tenant ====================

func (s *RBACService) CreateTenant(ctx context.Context, name, domain string) (*dal.Tenant, error) {
	t := &dal.Tenant{
		Name:      name,
		Domain:    domain,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.tenantRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}
	return t, nil
}

func (s *RBACService) GetTenant(ctx context.Context, id int64) (*dal.Tenant, error) {
	return s.tenantRepo.GetByID(ctx, id)
}

func (s *RBACService) ListTenants(ctx context.Context, page, pageSize int32) ([]*dal.Tenant, int32, error) {
	return s.tenantRepo.List(ctx, page, pageSize)
}

// ==================== Role ====================

func (s *RBACService) CreateRole(ctx context.Context, tenantID int64, name, description string) (*dal.Role, error) {
	r := &dal.Role{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := s.roleRepo.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return r, nil
}

func (s *RBACService) GetRole(ctx context.Context, id int64) (*dal.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

func (s *RBACService) UpdateRole(ctx context.Context, id int64, name, description string) (*dal.Role, error) {
	r, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		r.Name = name
	}
	if description != "" {
		r.Description = description
	}
	if err := s.roleRepo.Update(ctx, r); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}
	return r, nil
}

func (s *RBACService) DeleteRole(ctx context.Context, id int64) error {
	return s.roleRepo.Delete(ctx, id)
}

func (s *RBACService) ListRoles(ctx context.Context, tenantID int64, page, pageSize int32) ([]*dal.Role, int32, error) {
	return s.roleRepo.List(ctx, tenantID, page, pageSize)
}

// ==================== Permission ====================

func (s *RBACService) ListPermissions(ctx context.Context, page, pageSize int32) ([]*dal.Permission, int32, error) {
	return s.permRepo.List(ctx, page, pageSize)
}

// ==================== UserRole ====================

func (s *RBACService) AssignRole(ctx context.Context, userID, roleID int64) error {
	return s.userRoleRepo.Assign(ctx, userID, roleID)
}

func (s *RBACService) RevokeRole(ctx context.Context, userID, roleID int64) error {
	return s.userRoleRepo.Revoke(ctx, userID, roleID)
}

func (s *RBACService) GetUserRoles(ctx context.Context, userID int64) ([]*dal.Role, error) {
	return s.userRoleRepo.GetUserRoles(ctx, userID)
}

func (s *RBACService) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	return s.userRoleRepo.CheckPermission(ctx, userID, resource, action)
}

// ==================== RolePermission ====================

func (s *RBACService) AddRolePermission(ctx context.Context, roleID, permissionID int64) error {
	return s.rolePermRepo.Add(ctx, roleID, permissionID)
}

func (s *RBACService) RemoveRolePermission(ctx context.Context, roleID, permissionID int64) error {
	return s.rolePermRepo.Remove(ctx, roleID, permissionID)
}

func (s *RBACService) GetRolePermissions(ctx context.Context, roleID int64) ([]*dal.Permission, error) {
	return s.rolePermRepo.GetRolePermissions(ctx, roleID)
}
