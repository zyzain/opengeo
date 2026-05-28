package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AutoMigrate() error {
	// 使用 Migrator 逐表迁移，忽略已存在的索引错误
	tables := []interface{}{&Tenant{}, &User{}, &Role{}, &Permission{}, &UserRole{}, &RolePermission{}}
	for _, table := range tables {
		if err := s.db.AutoMigrate(table); err != nil {
			// 忽略重复索引错误 (Error 1061)
			if !strings.Contains(err.Error(), "Duplicate key name") {
				return err
			}
		}
	}
	return nil
}

// ==================== Tenant ====================

func (s *Store) CreateTenant(ctx context.Context, tenant *Tenant) error {
	return s.db.WithContext(ctx).Create(tenant).Error
}

func (s *Store) GetTenantByID(ctx context.Context, id int64) (*Tenant, error) {
	var t Tenant
	if err := s.db.WithContext(ctx).First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, err
	}
	return &t, nil
}

func (s *Store) ListTenants(ctx context.Context, page, pageSize int) ([]*Tenant, int64, error) {
	var tenants []*Tenant
	var total int64
	query := s.db.WithContext(ctx).Model(&Tenant{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, 0, err
	}
	return tenants, total, nil
}

// ==================== User ====================

func (s *Store) CreateUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByUsernameInTenant(ctx context.Context, tenantID int64, username string) (*User, error) {
	var user User
	if err := s.db.WithContext(ctx).Where("tenant_id = ? AND username = ?", tenantID, username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *Store) UpdateUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *Store) DeleteUser(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&User{}, id).Error
}

func (s *Store) ListUsersInTenant(ctx context.Context, tenantID int64, page, pageSize int, keyword string) ([]*User, int64, error) {
	var users []*User
	var total int64
	query := s.db.WithContext(ctx).Model(&User{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *Store) SeedAdmin(ctx context.Context, tenantID int64, hashedPassword string) error {
	var count int64
	s.db.WithContext(ctx).Model(&User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}
	admin := &User{
		TenantID:    tenantID,
		Username:    "admin",
		Password:    hashedPassword,
		Email:       "admin@opengeo.com",
		Status:      1,
		LastLoginAt: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.db.WithContext(ctx).Create(admin).Error
}

// ==================== Role ====================

func (s *Store) CreateRole(ctx context.Context, role *Role) error {
	return s.db.WithContext(ctx).Create(role).Error
}

func (s *Store) GetRoleByID(ctx context.Context, id int64) (*Role, error) {
	var role Role
	if err := s.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (s *Store) UpdateRole(ctx context.Context, role *Role) error {
	return s.db.WithContext(ctx).Save(role).Error
}

func (s *Store) DeleteRole(ctx context.Context, id int64) error {
	s.db.WithContext(ctx).Where("role_id = ?", id).Delete(&RolePermission{})
	s.db.WithContext(ctx).Where("role_id = ?", id).Delete(&UserRole{})
	return s.db.WithContext(ctx).Delete(&Role{}, id).Error
}

func (s *Store) ListRolesInTenant(ctx context.Context, tenantID int64, page, pageSize int) ([]*Role, int64, error) {
	var roles []*Role
	var total int64
	query := s.db.WithContext(ctx).Model(&Role{}).Where("tenant_id = ?", tenantID)
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

// ==================== Permission ====================

func (s *Store) CreatePermission(ctx context.Context, p *Permission) error {
	return s.db.WithContext(ctx).Create(p).Error
}

func (s *Store) GetPermissionByID(ctx context.Context, id int64) (*Permission, error) {
	var p Permission
	if err := s.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, err
	}
	return &p, nil
}

func (s *Store) ListPermissions(ctx context.Context, page, pageSize int) ([]*Permission, int64, error) {
	var perms []*Permission
	var total int64
	query := s.db.WithContext(ctx).Model(&Permission{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&perms).Error; err != nil {
		return nil, 0, err
	}
	return perms, total, nil
}

func (s *Store) GetAllPermissions(ctx context.Context) ([]*Permission, error) {
	var perms []*Permission
	if err := s.db.WithContext(ctx).Order("resource, action").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// ==================== UserRole ====================

func (s *Store) AssignRole(ctx context.Context, userID, roleID int64) error {
	ur := &UserRole{UserID: userID, RoleID: roleID}
	return s.db.WithContext(ctx).Create(ur).Error
}

func (s *Store) RevokeRole(ctx context.Context, userID, roleID int64) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{}).Error
}

func (s *Store) GetUserRoles(ctx context.Context, userID int64) ([]*Role, error) {
	var roles []*Role
	err := s.db.WithContext(ctx).
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// BatchGetUserRoles 批量获取用户角色（解决 N+1 问题）
func (s *Store) BatchGetUserRoles(ctx context.Context, userIDs []int64) (map[int64][]*Role, error) {
	if len(userIDs) == 0 {
		return make(map[int64][]*Role), nil
	}

	type userRoleResult struct {
		UserID int64 `gorm:"column:user_id"`
		Role
	}

	var results []userRoleResult
	err := s.db.WithContext(ctx).
		Table("user_roles").
		Select("user_roles.user_id, roles.*").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id IN ?", userIDs).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	roleMap := make(map[int64][]*Role)
	for _, r := range results {
		role := &Role{
			ID:          r.Role.ID,
			TenantID:    r.Role.TenantID,
			Name:        r.Role.Name,
			Description: r.Role.Description,
			CreatedAt:   r.Role.CreatedAt,
		}
		roleMap[r.UserID] = append(roleMap[r.UserID], role)
	}

	// 确保所有用户都有条目
	for _, uid := range userIDs {
		if _, ok := roleMap[uid]; !ok {
			roleMap[uid] = []*Role{}
		}
	}

	return roleMap, nil
}

// BatchGetRolePermissions 批量获取角色权限（解决 N+1 问题）
func (s *Store) BatchGetRolePermissions(ctx context.Context, roleIDs []int64) (map[int64][]*Permission, error) {
	if len(roleIDs) == 0 {
		return make(map[int64][]*Permission), nil
	}

	type rolePermResult struct {
		RoleID int64 `gorm:"column:role_id"`
		Permission
	}

	var results []rolePermResult
	err := s.db.WithContext(ctx).
		Table("role_permissions").
		Select("role_permissions.role_id, permissions.*").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	permMap := make(map[int64][]*Permission)
	for _, r := range results {
		perm := &Permission{
			ID:          r.Permission.ID,
			Name:        r.Permission.Name,
			Description: r.Permission.Description,
			Resource:    r.Permission.Resource,
			Action:      r.Permission.Action,
			CreatedAt:   r.Permission.CreatedAt,
		}
		permMap[r.RoleID] = append(permMap[r.RoleID], perm)
	}

	for _, rid := range roleIDs {
		if _, ok := permMap[rid]; !ok {
			permMap[rid] = []*Permission{}
		}
	}

	return permMap, nil
}

// ==================== RolePermission ====================

func (s *Store) AddRolePermission(ctx context.Context, roleID, permissionID int64) error {
	rp := &RolePermission{RoleID: roleID, PermissionID: permissionID}
	return s.db.WithContext(ctx).Create(rp).Error
}

func (s *Store) RemoveRolePermission(ctx context.Context, roleID, permissionID int64) error {
	return s.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermission{}).Error
}

func (s *Store) GetRolePermissions(ctx context.Context, roleID int64) ([]*Permission, error) {
	var perms []*Permission
	err := s.db.WithContext(ctx).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&perms).Error
	return perms, err
}

func (s *Store) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&UserRole{}).
		Select("count(*)").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Count(&count).Error
	return count > 0, err
}

// ==================== Seed ====================

func (s *Store) SeedTenant(ctx context.Context) (*Tenant, error) {
	var count int64
	s.db.WithContext(ctx).Model(&Tenant{}).Count(&count)
	if count > 0 {
		var t Tenant
		s.db.WithContext(ctx).First(&t)
		return &t, nil
	}
	tenant := &Tenant{
		Name:      "Default",
		Domain:    "default.opengeo.com",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.db.WithContext(ctx).Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *Store) SeedPermissions(ctx context.Context) error {
	var count int64
	s.db.WithContext(ctx).Model(&Permission{}).Count(&count)
	if count > 0 {
		return nil
	}
	perms := []*Permission{
		{Name: "user:create", Description: "创建用户", Resource: "user", Action: "create"},
		{Name: "user:read", Description: "查看用户", Resource: "user", Action: "read"},
		{Name: "user:update", Description: "更新用户", Resource: "user", Action: "update"},
		{Name: "user:delete", Description: "删除用户", Resource: "user", Action: "delete"},
		{Name: "content:create", Description: "创建内容", Resource: "content", Action: "create"},
		{Name: "content:read", Description: "查看内容", Resource: "content", Action: "read"},
		{Name: "content:update", Description: "更新内容", Resource: "content", Action: "update"},
		{Name: "content:delete", Description: "删除内容", Resource: "content", Action: "delete"},
		{Name: "publish:create", Description: "创建发布", Resource: "publish", Action: "create"},
		{Name: "publish:read", Description: "查看发布", Resource: "publish", Action: "read"},
		{Name: "publish:execute", Description: "执行发布", Resource: "publish", Action: "execute"},
		{Name: "role:create", Description: "创建角色", Resource: "role", Action: "create"},
		{Name: "role:read", Description: "查看角色", Resource: "role", Action: "read"},
		{Name: "role:update", Description: "更新角色", Resource: "role", Action: "update"},
		{Name: "role:delete", Description: "删除角色", Resource: "role", Action: "delete"},
		{Name: "tenant:read", Description: "查看租户", Resource: "tenant", Action: "read"},
		{Name: "tenant:update", Description: "更新租户", Resource: "tenant", Action: "update"},
		{Name: "brand:create", Description: "创建品牌", Resource: "brand", Action: "create"},
		{Name: "brand:read", Description: "查看品牌", Resource: "brand", Action: "read"},
		{Name: "brand:update", Description: "更新品牌", Resource: "brand", Action: "update"},
		{Name: "brand:delete", Description: "删除品牌", Resource: "brand", Action: "delete"},
	}
	for _, p := range perms {
		p.CreatedAt = time.Now()
		if err := s.db.WithContext(ctx).Create(p).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SeedRoles(ctx context.Context, tenantID int64) error {
	var count int64
	s.db.WithContext(ctx).Model(&Role{}).Where("tenant_id = ?", tenantID).Count(&count)
	if count > 0 {
		return nil
	}
	roles := []*Role{
		{TenantID: tenantID, Name: "admin", Description: "系统管理员"},
		{TenantID: tenantID, Name: "operator", Description: "运营人员"},
		{TenantID: tenantID, Name: "viewer", Description: "只读用户"},
	}
	for _, r := range roles {
		r.CreatedAt = time.Now()
		if err := s.db.WithContext(ctx).Create(r).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SeedAdminRoleAssignment(ctx context.Context, userID int64) error {
	var roles []*Role
	s.db.WithContext(ctx).Where("name = ?", "admin").Find(&roles)
	for _, r := range roles {
		var count int64
		s.db.WithContext(ctx).Model(&UserRole{}).Where("user_id = ? AND role_id = ?", userID, r.ID).Count(&count)
		if count == 0 {
			s.db.WithContext(ctx).Create(&UserRole{UserID: userID, RoleID: r.ID})
		}
	}
	return nil
}

func (s *Store) SeedRolePermissions(ctx context.Context) error {
	var count int64
	s.db.WithContext(ctx).Model(&RolePermission{}).Count(&count)
	if count > 0 {
		return nil
	}
	var adminRole Role
	if err := s.db.WithContext(ctx).Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	var perms []*Permission
	s.db.WithContext(ctx).Find(&perms)
	for _, p := range perms {
		s.db.WithContext(ctx).Create(&RolePermission{RoleID: adminRole.ID, PermissionID: p.ID})
	}

	var operatorRole Role
	if err := s.db.WithContext(ctx).Where("name = ?", "operator").First(&operatorRole).Error; err != nil {
		return nil
	}
	operatorPerms := []string{"content:create", "content:read", "content:update", "publish:create", "publish:read", "publish:execute", "user:read", "role:read", "brand:create", "brand:read", "brand:update"}
	for _, name := range operatorPerms {
		var p Permission
		if err := s.db.WithContext(ctx).Where("name = ?", name).First(&p).Error; err == nil {
			s.db.WithContext(ctx).Create(&RolePermission{RoleID: operatorRole.ID, PermissionID: p.ID})
		}
	}

	var viewerRole Role
	if err := s.db.WithContext(ctx).Where("name = ?", "viewer").First(&viewerRole).Error; err != nil {
		return nil
	}
	viewerPerms := []string{"content:read", "publish:read", "user:read", "role:read", "tenant:read", "brand:read"}
	for _, name := range viewerPerms {
		var p Permission
		if err := s.db.WithContext(ctx).Where("name = ?", name).First(&p).Error; err == nil {
			s.db.WithContext(ctx).Create(&RolePermission{RoleID: viewerRole.ID, PermissionID: p.ID})
		}
	}

	return nil
}
