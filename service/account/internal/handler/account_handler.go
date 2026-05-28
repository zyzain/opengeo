package handler

import (
	"context"

	"opengeo/service/account/internal/service"
)

// AccountHandler 账号服务处理器
type AccountHandler struct {
	accountService *service.AccountService
	userService    *service.UserService
	rbacService    *service.RBACService
}

// NewAccountHandler 创建账号服务处理器
func NewAccountHandler(
	accountService *service.AccountService,
	userService *service.UserService,
	rbacService *service.RBACService,
) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		userService:    userService,
		rbacService:    rbacService,
	}
}

// ==================== 用户认证 ====================

// Login 用户登录
func (h *AccountHandler) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	token, refreshToken, user, err := h.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

// Register 用户注册
func (h *AccountHandler) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	token, refreshToken, user, err := h.userService.Register(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		Token:        token,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

// RefreshToken 刷新Token
func (h *AccountHandler) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	token, refreshToken, err := h.userService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// ChangePassword 修改密码
func (h *AccountHandler) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	err := h.userService.ChangePassword(ctx, req.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &ChangePasswordResponse{Success: true}, nil
}

// ==================== 用户管理 ====================

// GetUser 获取用户
func (h *AccountHandler) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, err := h.userService.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateUser 更新用户
func (h *AccountHandler) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	user, err := h.userService.UpdateUser(ctx, req.UserID, req.Email, req.Status)
	if err != nil {
		return nil, err
	}

	return &UpdateUserResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// DeleteUser 删除用户
func (h *AccountHandler) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	err := h.userService.DeleteUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &DeleteUserResponse{Success: true}, nil
}

// ListUsers 列出用户
func (h *AccountHandler) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	users, total, err := h.userService.ListUsers(ctx, req.Page, req.PageSize, req.Keyword, req.TenantID)
	if err != nil {
		return nil, err
	}

	userList := make([]UserInfo, len(users))
	for i, user := range users {
		userList[i] = UserInfo{
			UserID:    user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &ListUsersResponse{
		Users: userList,
		Total: total,
	}, nil
}

// ==================== 账号管理 ====================

// CreateAccount 创建账号
func (h *AccountHandler) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	account, err := h.accountService.CreateAccount(ctx, req.UserID, req.Platform, req.AccountName, req.AccountID)
	if err != nil {
		return nil, err
	}

	return &CreateAccountResponse{
		AccountID:   account.ID,
		UserID:      account.UserID,
		Platform:    account.Platform,
		AccountName: account.AccountName,
	}, nil
}

// GetAccount 获取账号
func (h *AccountHandler) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	account, err := h.accountService.GetAccount(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetAccountResponse{
		AccountID:    account.ID,
		UserID:       account.UserID,
		Platform:     account.Platform,
		AccountName:  account.AccountName,
		AccountIDStr: account.AccountIDStr,
		Status:       account.Status,
		HealthScore:  account.HealthScore,
		CreatedAt:    account.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateAccount 更新账号
func (h *AccountHandler) UpdateAccount(ctx context.Context, req *UpdateAccountRequest) (*UpdateAccountResponse, error) {
	account, err := h.accountService.UpdateAccount(ctx, req.AccountID, req.AccountName, req.Status)
	if err != nil {
		return nil, err
	}

	return &UpdateAccountResponse{
		AccountID:   account.ID,
		AccountName: account.AccountName,
		Status:      account.Status,
		UpdatedAt:   account.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// DeleteAccount 删除账号
func (h *AccountHandler) DeleteAccount(ctx context.Context, req *DeleteAccountRequest) (*DeleteAccountResponse, error) {
	err := h.accountService.DeleteAccount(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	return &DeleteAccountResponse{Success: true}, nil
}

// ListAccounts 列出账号
func (h *AccountHandler) ListAccounts(ctx context.Context, req *ListAccountsRequest) (*ListAccountsResponse, error) {
	accounts, total, err := h.accountService.ListAccounts(ctx, req.UserID, req.Platform, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	accountList := make([]AccountInfo, len(accounts))
	for i, account := range accounts {
		accountList[i] = AccountInfo{
			AccountID:   account.ID,
			UserID:      account.UserID,
			Platform:    account.Platform,
			AccountName: account.AccountName,
			Status:      account.Status,
			HealthScore: account.HealthScore,
			CreatedAt:   account.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return &ListAccountsResponse{
		Accounts: accountList,
		Total:    total,
	}, nil
}

// ==================== 租户管理 ====================

// CreateTenant 创建租户
func (h *AccountHandler) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*CreateTenantResponse, error) {
	t, err := h.rbacService.CreateTenant(ctx, req.Name, req.Domain)
	if err != nil {
		return nil, err
	}
	return &CreateTenantResponse{
		ID:     t.ID,
		Name:   t.Name,
		Domain: t.Domain,
		Status: t.Status,
	}, nil
}

// GetTenant 获取租户
func (h *AccountHandler) GetTenant(ctx context.Context, req *GetTenantRequest) (*GetTenantResponse, error) {
	t, err := h.rbacService.GetTenant(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	return &GetTenantResponse{
		ID:        t.ID,
		Name:      t.Name,
		Domain:    t.Domain,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// ListTenants 列出租户
func (h *AccountHandler) ListTenants(ctx context.Context, req *ListTenantsRequest) (*ListTenantsResponse, error) {
	tenants, total, err := h.rbacService.ListTenants(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	list := make([]TenantInfo, len(tenants))
	for i, t := range tenants {
		list[i] = TenantInfo{
			ID:        t.ID,
			Name:      t.Name,
			Domain:    t.Domain,
			Status:    t.Status,
			CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return &ListTenantsResponse{Tenants: list, Total: total}, nil
}

// ==================== 角色管理 ====================

// CreateRole 创建角色
func (h *AccountHandler) CreateRole(ctx context.Context, req *CreateRoleReq) (*CreateRoleResp, error) {
	r, err := h.rbacService.CreateRole(ctx, req.TenantID, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return &CreateRoleResp{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Description: r.Description,
	}, nil
}

// GetRole 获取角色
func (h *AccountHandler) GetRole(ctx context.Context, req *GetRoleReq) (*GetRoleResp, error) {
	r, err := h.rbacService.GetRole(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	return &GetRoleResp{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateRole 更新角色
func (h *AccountHandler) UpdateRole(ctx context.Context, req *UpdateRoleReq) (*UpdateRoleResp, error) {
	r, err := h.rbacService.UpdateRole(ctx, req.RoleID, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return &UpdateRoleResp{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}, nil
}

// DeleteRole 删除角色
func (h *AccountHandler) DeleteRole(ctx context.Context, req *DeleteRoleReq) (*DeleteRoleResp, error) {
	err := h.rbacService.DeleteRole(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	return &DeleteRoleResp{Success: true}, nil
}

// ListRoles 列出角色
func (h *AccountHandler) ListRoles(ctx context.Context, req *ListRolesReq) (*ListRolesResp, error) {
	roles, total, err := h.rbacService.ListRoles(ctx, req.TenantID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	list := make([]RoleInfo, len(roles))
	for i, r := range roles {
		list[i] = RoleInfo{
			ID:          r.ID,
			TenantID:    r.TenantID,
			Name:        r.Name,
			Description: r.Description,
			CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return &ListRolesResp{Roles: list, Total: total}, nil
}

// AssignRole 分配角色
func (h *AccountHandler) AssignRole(ctx context.Context, req *AssignRoleReq) (*AssignRoleResp, error) {
	err := h.rbacService.AssignRole(ctx, req.UserID, req.RoleID)
	if err != nil {
		return nil, err
	}
	return &AssignRoleResp{Success: true}, nil
}

// RevokeRole 撤销角色
func (h *AccountHandler) RevokeRole(ctx context.Context, req *RevokeRoleReq) (*RevokeRoleResp, error) {
	err := h.rbacService.RevokeRole(ctx, req.UserID, req.RoleID)
	if err != nil {
		return nil, err
	}
	return &RevokeRoleResp{Success: true}, nil
}

// GetUserRoles 获取用户角色
func (h *AccountHandler) GetUserRoles(ctx context.Context, req *GetUserRolesReq) (*GetUserRolesResp, error) {
	roles, err := h.rbacService.GetUserRoles(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	list := make([]RoleInfo, len(roles))
	for i, r := range roles {
		list[i] = RoleInfo{
			ID:          r.ID,
			TenantID:    r.TenantID,
			Name:        r.Name,
			Description: r.Description,
		}
	}
	return &GetUserRolesResp{Roles: list}, nil
}

// ==================== 权限管理 ====================

// ListPermissions 列出权限
func (h *AccountHandler) ListPermissions(ctx context.Context, req *ListPermissionsReq) (*ListPermissionsResp, error) {
	perms, total, err := h.rbacService.ListPermissions(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	list := make([]PermissionInfo, len(perms))
	for i, p := range perms {
		list[i] = PermissionInfo{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Resource:    p.Resource,
			Action:      p.Action,
		}
	}
	return &ListPermissionsResp{Permissions: list, Total: total}, nil
}

// AddRolePermission 添加角色权限
func (h *AccountHandler) AddRolePermission(ctx context.Context, req *AddRolePermissionReq) (*AddRolePermissionResp, error) {
	err := h.rbacService.AddRolePermission(ctx, req.RoleID, req.PermissionID)
	if err != nil {
		return nil, err
	}
	return &AddRolePermissionResp{Success: true}, nil
}

// RemoveRolePermission 移除角色权限
func (h *AccountHandler) RemoveRolePermission(ctx context.Context, req *RemoveRolePermissionReq) (*RemoveRolePermissionResp, error) {
	err := h.rbacService.RemoveRolePermission(ctx, req.RoleID, req.PermissionID)
	if err != nil {
		return nil, err
	}
	return &RemoveRolePermissionResp{Success: true}, nil
}

// GetRolePermissions 获取角色权限
func (h *AccountHandler) GetRolePermissions(ctx context.Context, req *GetRolePermissionsReq) (*GetRolePermissionsResp, error) {
	perms, err := h.rbacService.GetRolePermissions(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	list := make([]PermissionInfo, len(perms))
	for i, p := range perms {
		list[i] = PermissionInfo{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Resource:    p.Resource,
			Action:      p.Action,
		}
	}
	return &GetRolePermissionsResp{Permissions: list}, nil
}

// CheckPermission 检查权限
func (h *AccountHandler) CheckPermission(ctx context.Context, req *CheckPermissionReq) (*CheckPermissionResp, error) {
	allowed, err := h.rbacService.CheckPermission(ctx, req.UserID, req.Resource, req.Action)
	if err != nil {
		return nil, err
	}
	return &CheckPermissionResp{Allowed: allowed}, nil
}

// GetAccountHealth 获取账号健康状态
func (h *AccountHandler) GetAccountHealth(ctx context.Context, req *GetAccountHealthRequest) (*GetAccountHealthResponse, error) {
	health, err := h.accountService.GetAccountHealth(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetAccountHealthResponse{
		AccountID:    health.AccountID,
		HealthScore:  health.HealthScore,
		Status:       health.Status,
		CheckDetails: health.CheckDetails,
		CheckedAt:    health.CheckedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// ==================== 请求/响应模型 ====================

// 认证相关
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type ChangePasswordRequest struct {
	UserID      int64  `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

// 用户相关
type GetUserRequest struct {
	UserID int64 `json:"user_id"`
}

type GetUserResponse struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
}

type UpdateUserRequest struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Status int32  `json:"status"`
}

type UpdateUserResponse struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    int32  `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type DeleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

type ListUsersRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Keyword  string `json:"keyword"`
	TenantID int64  `json:"tenant_id"`
}

type ListUsersResponse struct {
	Users []UserInfo `json:"users"`
	Total int32      `json:"total"`
}

type UserInfo struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
}

// 账号相关
type CreateAccountRequest struct {
	UserID      int64  `json:"user_id"`
	Platform    string `json:"platform"`
	AccountName string `json:"account_name"`
	AccountID   string `json:"account_id"`
}

type CreateAccountResponse struct {
	AccountID   int64  `json:"account_id"`
	UserID      int64  `json:"user_id"`
	Platform    string `json:"platform"`
	AccountName string `json:"account_name"`
}

type GetAccountRequest struct {
	AccountID int64 `json:"account_id"`
}

type GetAccountResponse struct {
	AccountID    int64   `json:"account_id"`
	UserID       int64   `json:"user_id"`
	Platform     string  `json:"platform"`
	AccountName  string  `json:"account_name"`
	AccountIDStr string  `json:"account_id_str"`
	Status       int32   `json:"status"`
	HealthScore  float32 `json:"health_score"`
	CreatedAt    string  `json:"created_at"`
}

type UpdateAccountRequest struct {
	AccountID   int64  `json:"account_id"`
	AccountName string `json:"account_name"`
	Status      int32  `json:"status"`
}

type UpdateAccountResponse struct {
	AccountID   int64  `json:"account_id"`
	AccountName string `json:"account_name"`
	Status      int32  `json:"status"`
	UpdatedAt   string `json:"updated_at"`
}

type DeleteAccountRequest struct {
	AccountID int64 `json:"account_id"`
}

type DeleteAccountResponse struct {
	Success bool `json:"success"`
}

type ListAccountsRequest struct {
	UserID   int64  `json:"user_id"`
	Platform string `json:"platform"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

type ListAccountsResponse struct {
	Accounts []AccountInfo `json:"accounts"`
	Total    int32         `json:"total"`
}

type AccountInfo struct {
	AccountID   int64   `json:"account_id"`
	UserID      int64   `json:"user_id"`
	Platform    string  `json:"platform"`
	AccountName string  `json:"account_name"`
	Status      int32   `json:"status"`
	HealthScore float32 `json:"health_score"`
	CreatedAt   string  `json:"created_at"`
}

type GetAccountHealthRequest struct {
	AccountID int64 `json:"account_id"`
}

type GetAccountHealthResponse struct {
	AccountID    int64   `json:"account_id"`
	HealthScore  float32 `json:"health_score"`
	Status       string  `json:"status"`
	CheckDetails string  `json:"check_details"`
	CheckedAt    string  `json:"checked_at"`
}

// ==================== RBAC 请求/响应 ====================

// Tenant
type CreateTenantRequest struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type CreateTenantResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Status int32  `json:"status"`
}

type GetTenantRequest struct {
	TenantID int64 `json:"tenant_id"`
}

type GetTenantResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ListTenantsRequest struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type ListTenantsResponse struct {
	Tenants []TenantInfo `json:"tenants"`
	Total   int32        `json:"total"`
}

type TenantInfo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Status    int32  `json:"status"`
	CreatedAt string `json:"created_at"`
}

// Role
type CreateRoleReq struct {
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRoleResp struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetRoleReq struct {
	RoleID int64 `json:"role_id"`
}

type GetRoleResp struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type UpdateRoleReq struct {
	RoleID      int64  `json:"role_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRoleResp struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DeleteRoleReq struct {
	RoleID int64 `json:"role_id"`
}

type DeleteRoleResp struct {
	Success bool `json:"success"`
}

type ListRolesReq struct {
	TenantID int64 `json:"tenant_id"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type ListRolesResp struct {
	Roles []RoleInfo `json:"roles"`
	Total int32      `json:"total"`
}

type RoleInfo struct {
	ID          int64  `json:"id"`
	TenantID    int64  `json:"tenant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type AssignRoleReq struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

type AssignRoleResp struct {
	Success bool `json:"success"`
}

type RevokeRoleReq struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

type RevokeRoleResp struct {
	Success bool `json:"success"`
}

type GetUserRolesReq struct {
	UserID int64 `json:"user_id"`
}

type GetUserRolesResp struct {
	Roles []RoleInfo `json:"roles"`
}

// Permission
type ListPermissionsReq struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

type ListPermissionsResp struct {
	Permissions []PermissionInfo `json:"permissions"`
	Total       int32            `json:"total"`
}

type PermissionInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}

type AddRolePermissionReq struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

type AddRolePermissionResp struct {
	Success bool `json:"success"`
}

type RemoveRolePermissionReq struct {
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
}

type RemoveRolePermissionResp struct {
	Success bool `json:"success"`
}

type GetRolePermissionsReq struct {
	RoleID int64 `json:"role_id"`
}

type GetRolePermissionsResp struct {
	Permissions []PermissionInfo `json:"permissions"`
}

type CheckPermissionReq struct {
	UserID   int64  `json:"user_id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type CheckPermissionResp struct {
	Allowed bool `json:"allowed"`
}