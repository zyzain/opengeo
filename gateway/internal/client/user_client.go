package client

import (
	"context"

	"opengeo/gateway/internal/auth"
)

type UserClient struct {
	svc *auth.Service
}

func NewUserClient(svc *auth.Service) *UserClient {
	return &UserClient{svc: svc}
}

func (c *UserClient) Login(ctx context.Context, username, password string) (map[string]interface{}, error) {
	return c.svc.Login(ctx, username, password)
}

func (c *UserClient) Register(ctx context.Context, username, password, email string) (map[string]interface{}, error) {
	return c.svc.Register(ctx, username, password, email)
}

func (c *UserClient) RefreshToken(ctx context.Context, refreshToken string) (map[string]interface{}, error) {
	return c.svc.RefreshToken(ctx, refreshToken)
}

func (c *UserClient) GetUser(ctx context.Context, id int64) (map[string]interface{}, error) {
	return c.svc.GetUser(ctx, id)
}

func (c *UserClient) UpdateUser(ctx context.Context, id int64, email string, status int32) (map[string]interface{}, error) {
	return c.svc.UpdateUser(ctx, id, email, status)
}

func (c *UserClient) DeleteUser(ctx context.Context, id int64) error {
	return c.svc.DeleteUser(ctx, id)
}

func (c *UserClient) ListUsers(ctx context.Context, page, pageSize int, keyword string) (map[string]interface{}, error) {
	return c.svc.ListUsers(ctx, page, pageSize, keyword)
}

func (c *UserClient) CreateTenant(ctx context.Context, name, domain string) (map[string]interface{}, error) {
	return c.svc.CreateTenant(ctx, name, domain)
}

func (c *UserClient) GetTenant(ctx context.Context, id int64) (map[string]interface{}, error) {
	return c.svc.GetTenant(ctx, id)
}

func (c *UserClient) ListTenants(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	return c.svc.ListTenants(ctx, page, pageSize)
}

func (c *UserClient) CreateRole(ctx context.Context, tenantID int64, name, description string) (map[string]interface{}, error) {
	return c.svc.CreateRole(ctx, tenantID, name, description)
}

func (c *UserClient) GetRole(ctx context.Context, id int64) (map[string]interface{}, error) {
	return c.svc.GetRole(ctx, id)
}

func (c *UserClient) UpdateRole(ctx context.Context, id int64, name, description string) (map[string]interface{}, error) {
	return c.svc.UpdateRole(ctx, id, name, description)
}

func (c *UserClient) DeleteRole(ctx context.Context, id int64) error {
	return c.svc.DeleteRole(ctx, id)
}

func (c *UserClient) ListRoles(ctx context.Context, tenantID int64, page, pageSize int) (map[string]interface{}, error) {
	return c.svc.ListRoles(ctx, tenantID, page, pageSize)
}

func (c *UserClient) AssignRole(ctx context.Context, userID, roleID int64) error {
	return c.svc.AssignRole(ctx, userID, roleID)
}

func (c *UserClient) RevokeRole(ctx context.Context, userID, roleID int64) error {
	return c.svc.RevokeRole(ctx, userID, roleID)
}

func (c *UserClient) GetUserRoles(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	return c.svc.GetUserRoles(ctx, userID)
}

func (c *UserClient) ListPermissions(ctx context.Context, page, pageSize int) (map[string]interface{}, error) {
	return c.svc.ListPermissions(ctx, page, pageSize)
}

func (c *UserClient) AddRolePermission(ctx context.Context, roleID, permissionID int64) error {
	return c.svc.AddRolePermission(ctx, roleID, permissionID)
}

func (c *UserClient) RemoveRolePermission(ctx context.Context, roleID, permissionID int64) error {
	return c.svc.RemoveRolePermission(ctx, roleID, permissionID)
}

func (c *UserClient) GetRolePermissions(ctx context.Context, roleID int64) ([]map[string]interface{}, error) {
	return c.svc.GetRolePermissions(ctx, roleID)
}

func (c *UserClient) CheckPermission(ctx context.Context, userID int64, resource, action string) (map[string]interface{}, error) {
	allowed, err := c.svc.CheckPermission(ctx, userID, resource, action)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"allowed": allowed}, nil
}
