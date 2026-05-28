package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TenantRepository 租户仓储
type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *TenantRepository) GetByID(ctx context.Context, id int64) (*Tenant, error) {
	var t Tenant
	if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
		return nil, fmt.Errorf("tenant not found: %d", id)
	}
	return &t, nil
}

func (r *TenantRepository) List(ctx context.Context, page, pageSize int32) ([]*Tenant, int32, error) {
	var tenants []*Tenant
	var total int64
	query := r.db.WithContext(ctx).Model(&Tenant{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, 0, err
	}
	return tenants, int32(total), nil
}

// RoleRepository 角色仓储
type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) GetByID(ctx context.Context, id int64) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, fmt.Errorf("role not found: %d", id)
	}
	return &role, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	r.db.WithContext(ctx).Where("role_id = ?", id).Delete(&RolePermission{})
	r.db.WithContext(ctx).Where("role_id = ?", id).Delete(&UserRole{})
	return r.db.WithContext(ctx).Delete(&Role{}, id).Error
}

func (r *RoleRepository) List(ctx context.Context, tenantID int64, page, pageSize int32) ([]*Role, int32, error) {
	var roles []*Role
	var total int64
	query := r.db.WithContext(ctx).Model(&Role{}).Where("tenant_id = ?", tenantID)
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, int32(total), nil
}

// PermissionRepository 权限仓储
type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) List(ctx context.Context, page, pageSize int32) ([]*Permission, int32, error) {
	var perms []*Permission
	var total int64
	query := r.db.WithContext(ctx).Model(&Permission{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(int(offset)).Limit(int(pageSize)).Order("resource, action").Find(&perms).Error; err != nil {
		return nil, 0, err
	}
	return perms, int32(total), nil
}

// UserRoleRepository 用户角色关联仓储
type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) Assign(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).Create(&UserRole{UserID: userID, RoleID: roleID}).Error
}

func (r *UserRoleRepository) Revoke(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error
}

func (r *UserRoleRepository) GetUserRoles(ctx context.Context, userID int64) ([]*Role, error) {
	var roles []*Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *UserRoleRepository) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserRole{}).
		Select("count(*)").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Count(&count).Error
	return count > 0, err
}

// RolePermissionRepository 角色权限关联仓储
type RolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

func (r *RolePermissionRepository) Add(ctx context.Context, roleID, permissionID int64) error {
	return r.db.WithContext(ctx).Create(&RolePermission{RoleID: roleID, PermissionID: permissionID}).Error
}

func (r *RolePermissionRepository) Remove(ctx context.Context, roleID, permissionID int64) error {
	return r.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermission{}).Error
}

func (r *RolePermissionRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*Permission, error) {
	var perms []*Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&perms).Error
	return perms, err
}

// SeedTenant 创建默认租户
func SeedTenant(ctx context.Context, db *gorm.DB) *Tenant {
	var count int64
	db.Model(&Tenant{}).Count(&count)
	if count > 0 {
		var t Tenant
		db.First(&t)
		return &t
	}
	t := &Tenant{
		Name:      "Default",
		Domain:    "default.opengeo.com",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(t)
	return t
}
