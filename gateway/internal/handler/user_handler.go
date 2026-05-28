package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// GetUser 获取用户
func (h *Handler) GetUser(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid user id"))
		return
	}

	resp, err := h.userClient.GetUser(ctx, id)
	if err != nil {
		errResponse(c, err, "get user failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// UpdateUser 更新用户
func (h *Handler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid user id"))
		return
	}

	var req struct {
		Email  string `json:"email"`
		Status int32  `json:"status"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.UpdateUser(ctx, id, req.Email, req.Status)
	if err != nil {
		errResponse(c, err, "update user failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// DeleteUser 删除用户
func (h *Handler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid user id"))
		return
	}

	if err := h.userClient.DeleteUser(ctx, id); err != nil {
		errResponse(c, err, "delete user failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListUsers 列出用户
func (h *Handler) ListUsers(ctx context.Context, c *app.RequestContext) {
	page, pageSize := parsePagination(c)
	keyword := c.Query("keyword")

	resp, err := h.userClient.ListUsers(ctx, page, pageSize, keyword)
	if err != nil {
		errResponse(c, err, "list users failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// GetUserRoles 获取用户角色
func (h *Handler) GetUserRoles(ctx context.Context, c *app.RequestContext) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid user id"))
		return
	}

	resp, err := h.userClient.GetUserRoles(ctx, userID)
	if err != nil {
		errResponse(c, err, "get user roles failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// ==================== 角色管理 ====================

// CreateRole 创建角色
func (h *Handler) CreateRole(ctx context.Context, c *app.RequestContext) {
	var req struct {
		TenantID    int64  `json:"tenant_id"`
		Name        string `json:"name" vd:"len($)>0"`
		Description string `json:"description"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.CreateRole(ctx, req.TenantID, req.Name, req.Description)
	if err != nil {
		errResponse(c, err, "create role failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// GetRole 获取角色
func (h *Handler) GetRole(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid role id"))
		return
	}

	resp, err := h.userClient.GetRole(ctx, id)
	if err != nil {
		errResponse(c, err, "get role failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// UpdateRole 更新角色
func (h *Handler) UpdateRole(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid role id"))
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.UpdateRole(ctx, id, req.Name, req.Description)
	if err != nil {
		errResponse(c, err, "update role failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// DeleteRole 删除角色
func (h *Handler) DeleteRole(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid role id"))
		return
	}

	if err := h.userClient.DeleteRole(ctx, id); err != nil {
		errResponse(c, err, "delete role failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListRoles 列出角色
func (h *Handler) ListRoles(ctx context.Context, c *app.RequestContext) {
	tenantID, _ := strconv.ParseInt(c.Query("tenant_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	resp, err := h.userClient.ListRoles(ctx, tenantID, page, pageSize)
	if err != nil {
		errResponse(c, err, "list roles failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// AssignRole 分配角色
func (h *Handler) AssignRole(ctx context.Context, c *app.RequestContext) {
	var req struct {
		UserID int64 `json:"user_id" vd:"$>0"`
		RoleID int64 `json:"role_id" vd:"$>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	if err := h.userClient.AssignRole(ctx, req.UserID, req.RoleID); err != nil {
		errResponse(c, err, "assign role failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// RevokeRole 撤销角色
func (h *Handler) RevokeRole(ctx context.Context, c *app.RequestContext) {
	var req struct {
		UserID int64 `json:"user_id" vd:"$>0"`
		RoleID int64 `json:"role_id" vd:"$>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	if err := h.userClient.RevokeRole(ctx, req.UserID, req.RoleID); err != nil {
		errResponse(c, err, "revoke role failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// GetRolePermissions 获取角色权限
func (h *Handler) GetRolePermissions(ctx context.Context, c *app.RequestContext) {
	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid role id"))
		return
	}

	resp, err := h.userClient.GetRolePermissions(ctx, roleID)
	if err != nil {
		errResponse(c, err, "get role permissions failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// AddRolePermission 添加角色权限
func (h *Handler) AddRolePermission(ctx context.Context, c *app.RequestContext) {
	var req struct {
		RoleID       int64 `json:"role_id" vd:"$>0"`
		PermissionID int64 `json:"permission_id" vd:"$>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	if err := h.userClient.AddRolePermission(ctx, req.RoleID, req.PermissionID); err != nil {
		errResponse(c, err, "add role permission failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// RemoveRolePermission 移除角色权限
func (h *Handler) RemoveRolePermission(ctx context.Context, c *app.RequestContext) {
	var req struct {
		RoleID       int64 `json:"role_id" vd:"$>0"`
		PermissionID int64 `json:"permission_id" vd:"$>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	if err := h.userClient.RemoveRolePermission(ctx, req.RoleID, req.PermissionID); err != nil {
		errResponse(c, err, "remove role permission failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListPermissions 列出权限
func (h *Handler) ListPermissions(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	resp, err := h.userClient.ListPermissions(ctx, page, pageSize)
	if err != nil {
		errResponse(c, err, "list permissions failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// CheckPermission 检查权限
func (h *Handler) CheckPermission(ctx context.Context, c *app.RequestContext) {
	var req struct {
		UserID   int64  `json:"user_id" vd:"$>0"`
		Resource string `json:"resource" vd:"len($)>0"`
		Action   string `json:"action" vd:"len($)>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.CheckPermission(ctx, req.UserID, req.Resource, req.Action)
	if err != nil {
		errResponse(c, err, "check permission failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// ==================== 租户管理 ====================

// CreateTenant 创建租户
func (h *Handler) CreateTenant(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Name   string `json:"name" vd:"len($)>0"`
		Domain string `json:"domain"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.userClient.CreateTenant(ctx, req.Name, req.Domain)
	if err != nil {
		errResponse(c, err, "create tenant failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// GetTenant 获取租户
func (h *Handler) GetTenant(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid tenant id"))
		return
	}

	resp, err := h.userClient.GetTenant(ctx, id)
	if err != nil {
		errResponse(c, err, "get tenant failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// ListTenants 列出租户
func (h *Handler) ListTenants(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	resp, err := h.userClient.ListTenants(ctx, page, pageSize)
	if err != nil {
		errResponse(c, err, "list tenants failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// UpdateTenant 更新租户
func (h *Handler) UpdateTenant(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid tenant id"))
		return
	}

	var req struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
		Plan   string `json:"plan"`
		Status string `json:"status"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{
		"id":     id,
		"name":   req.Name,
		"domain": req.Domain,
		"plan":   req.Plan,
		"status": req.Status,
	}))
}

// DeleteTenant 删除租户
func (h *Handler) DeleteTenant(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid tenant id"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true, "id": id}))
}
