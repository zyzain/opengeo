package auth

import (
	"context"
	"fmt"
	"time"

	"opengeo/pkg/crypto"
	jwtUtil "opengeo/pkg/jwt"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Login(ctx context.Context, username, password string) (map[string]interface{}, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if user.Status != 1 {
		return nil, fmt.Errorf("user account is disabled")
	}

	if !crypto.CheckPassword(password, user.Password) {
		return nil, fmt.Errorf("invalid username or password")
	}

	roles, _ := s.store.GetUserRoles(ctx, user.ID)
	roleName := "user"
	if len(roles) > 0 {
		roleName = roles[0].Name
	}

	token, err := jwtUtil.GenerateToken(user.ID, user.Username, user.Email, roleName, user.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := jwtUtil.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	user.LastLoginAt = time.Now()
	s.store.UpdateUser(ctx, user)

	return map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
		"user_id":       user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"tenant_id":     user.TenantID,
		"role":          roleName,
	}, nil
}

func (s *Service) Register(ctx context.Context, username, password, email string) (map[string]interface{}, error) {
	existing, _ := s.store.GetUserByUsername(ctx, username)
	if existing != nil {
		return nil, fmt.Errorf("username already exists")
	}

	if err := crypto.ValidatePasswordStrength(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		TenantID:  0,
		Username:  username,
		Password:  hashedPassword,
		Email:     email,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := jwtUtil.GenerateToken(user.ID, user.Username, user.Email, "user", user.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := jwtUtil.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
		"user_id":       user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"tenant_id":     user.TenantID,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (map[string]interface{}, error) {
	newToken, newRefreshToken, err := jwtUtil.RefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return map[string]interface{}{
		"token":         newToken,
		"refresh_token": newRefreshToken,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, id int64) (map[string]interface{}, error) {
	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	roles, _ := s.store.GetUserRoles(ctx, id)
	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}
	return map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"tenant_id":  user.TenantID,
		"status":     user.Status,
		"roles":      roleNames,
		"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int64, email string, status int32) (map[string]interface{}, error) {
	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if email != "" {
		user.Email = email
	}
	if status > 0 {
		user.Status = status
	}
	user.UpdatedAt = time.Now()
	if err := s.store.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return map[string]interface{}{
		"id":     user.ID,
		"email":  user.Email,
		"status": user.Status,
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	user, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return s.store.DeleteUser(ctx, user.ID)
}

func (s *Service) ListUsers(ctx context.Context, page, pageSize int, keyword string) (map[string]interface{}, error) {
	users, total, err := s.store.ListUsersInTenant(ctx, 0, page, pageSize, keyword)
	if err != nil {
		return nil, err
	}

	// 批量获取所有用户的角色（避免 N+1 查询）
	userIDs := make([]int64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	roleMap, _ := s.store.BatchGetUserRoles(ctx, userIDs)

	userList := make([]map[string]interface{}, len(users))
	for i, u := range users {
		roles := roleMap[u.ID]
		roleNames := make([]string, len(roles))
		for j, r := range roles {
			roleNames[j] = r.Name
		}
		userList[i] = map[string]interface{}{
			"id":         u.ID,
			"username":   u.Username,
			"email":      u.Email,
			"tenant_id":  u.TenantID,
			"status":     u.Status,
			"roles":      roleNames,
			"created_at": u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return map[string]interface{}{
		"users": userList,
		"total": total,
	}, nil
}

// ==================== Tenant ====================

func (s *Service) CreateTenant(ctx context.Context, name, domain string) (map[string]interface{}, error) {
	t := &Tenant{
		Name:      name,
		Domain:    domain,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.store.CreateTenant(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}
	return map[string]interface{}{
		"id":     t.ID,
		"name":   t.Name,
		"domain": t.Domain,
		"status": t.Status,
	}, nil
}

func (s *Service) GetTenant(ctx context.Context, id int64) (map[string]interface{}, error) {
	t, err := s.store.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":         t.ID,
		"name":       t.Name,
		"domain":     t.Domain,
		"status":     t.Status,
		"created_at": t.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) ListTenants(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	tenants, total, err := s.store.ListTenants(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(tenants))
	for i, t := range tenants {
		list[i] = map[string]interface{}{
			"id":         t.ID,
			"name":       t.Name,
			"domain":     t.Domain,
			"status":     t.Status,
			"created_at": t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return map[string]interface{}{
		"tenants": list,
		"total":   total,
	}, nil
}

// ==================== Role ====================

func (s *Service) CreateRole(ctx context.Context, tenantID int64, name, description string) (map[string]interface{}, error) {
	r := &Role{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateRole(ctx, r); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return map[string]interface{}{
		"id":          r.ID,
		"tenant_id":   r.TenantID,
		"name":        r.Name,
		"description": r.Description,
	}, nil
}

func (s *Service) GetRole(ctx context.Context, id int64) (map[string]interface{}, error) {
	r, err := s.store.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	perms, _ := s.store.GetRolePermissions(ctx, id)
	permNames := make([]string, len(perms))
	for i, p := range perms {
		permNames[i] = p.Name
	}
	return map[string]interface{}{
		"id":          r.ID,
		"tenant_id":   r.TenantID,
		"name":        r.Name,
		"description": r.Description,
		"permissions": permNames,
		"created_at":  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *Service) UpdateRole(ctx context.Context, id int64, name, description string) (map[string]interface{}, error) {
	r, err := s.store.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		r.Name = name
	}
	if description != "" {
		r.Description = description
	}
	if err := s.store.UpdateRole(ctx, r); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}
	return map[string]interface{}{
		"id":          r.ID,
		"name":        r.Name,
		"description": r.Description,
	}, nil
}

func (s *Service) DeleteRole(ctx context.Context, id int64) error {
	return s.store.DeleteRole(ctx, id)
}

func (s *Service) ListRoles(ctx context.Context, tenantID int64, page, pageSize int) (map[string]interface{}, error) {
	roles, total, err := s.store.ListRolesInTenant(ctx, tenantID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 批量获取所有角色的权限（避免 N+1 查询）
	roleIDs := make([]int64, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.ID
	}

	permMap, _ := s.store.BatchGetRolePermissions(ctx, roleIDs)

	list := make([]map[string]interface{}, len(roles))
	for i, r := range roles {
		perms := permMap[r.ID]
		permNames := make([]string, len(perms))
		for j, p := range perms {
			permNames[j] = p.Name
		}
		list[i] = map[string]interface{}{
			"id":          r.ID,
			"tenant_id":   r.TenantID,
			"name":        r.Name,
			"description": r.Description,
			"permissions": permNames,
			"created_at":  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return map[string]interface{}{
		"roles": list,
		"total": total,
	}, nil
}

func (s *Service) AssignRole(ctx context.Context, userID, roleID int64) error {
	return s.store.AssignRole(ctx, userID, roleID)
}

func (s *Service) RevokeRole(ctx context.Context, userID, roleID int64) error {
	return s.store.RevokeRole(ctx, userID, roleID)
}

func (s *Service) GetUserRoles(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	roles, err := s.store.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(roles))
	for i, r := range roles {
		list[i] = map[string]interface{}{
			"id":          r.ID,
			"name":        r.Name,
			"description": r.Description,
		}
	}
	return list, nil
}

// ==================== Permission ====================

func (s *Service) ListPermissions(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	perms, total, err := s.store.ListPermissions(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(perms))
	for i, p := range perms {
		list[i] = map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"resource":    p.Resource,
			"action":      p.Action,
		}
	}
	return map[string]interface{}{
		"permissions": list,
		"total":       total,
	}, nil
}

func (s *Service) AddRolePermission(ctx context.Context, roleID, permissionID int64) error {
	return s.store.AddRolePermission(ctx, roleID, permissionID)
}

func (s *Service) RemoveRolePermission(ctx context.Context, roleID, permissionID int64) error {
	return s.store.RemoveRolePermission(ctx, roleID, permissionID)
}

func (s *Service) GetRolePermissions(ctx context.Context, roleID int64) ([]map[string]interface{}, error) {
	perms, err := s.store.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(perms))
	for i, p := range perms {
		list[i] = map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"resource":    p.Resource,
			"action":      p.Action,
		}
	}
	return list, nil
}

func (s *Service) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	return s.store.CheckPermission(ctx, userID, resource, action)
}
