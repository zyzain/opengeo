package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// CreateAccount 创建账号
func (h *Handler) CreateAccount(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Platform    string `json:"platform" vd:"len($)>0"`
		AccountName string `json:"account_name" vd:"len($)>0"`
		AccountID   string `json:"account_id"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	resp, err := h.accountClient.CreateAccount(ctx, userID, req.Platform, req.AccountName, req.AccountID)
	if err != nil {
		errResponse(c, err, "create account failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// GetAccount 获取账号
func (h *Handler) GetAccount(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid account id"))
		return
	}

	resp, err := h.accountClient.GetAccount(ctx, id)
	if err != nil {
		errResponse(c, err, "get account failed")
		return
	}

	if !checkOwnership(c, resp) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// UpdateAccount 更新账号
func (h *Handler) UpdateAccount(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid account id"))
		return
	}

	existing, err := h.accountClient.GetAccount(ctx, id)
	if err != nil {
		errResponse(c, err, "account not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	var req struct {
		AccountName string `json:"account_name"`
		Status      int32  `json:"status"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.accountClient.UpdateAccount(ctx, id, req.AccountName, req.Status)
	if err != nil {
		errResponse(c, err, "update account failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// DeleteAccount 删除账号
func (h *Handler) DeleteAccount(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid account id"))
		return
	}

	existing, err := h.accountClient.GetAccount(ctx, id)
	if err != nil {
		errResponse(c, err, "account not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	if err := h.accountClient.DeleteAccount(ctx, id); err != nil {
		errResponse(c, err, "delete account failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListAccounts 列出账号
func (h *Handler) ListAccounts(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	platform := c.Query("platform")
	page, pageSize := parsePagination(c)

	resp, err := h.accountClient.ListAccounts(ctx, userID, platform, page, pageSize)
	if err != nil {
		errResponse(c, err, "list accounts failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// GetAccountHealth 获取账号健康
func (h *Handler) GetAccountHealth(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid account id"))
		return
	}

	resp, err := h.accountClient.GetAccountHealth(ctx, id)
	if err != nil {
		errResponse(c, err, "get account health failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// ==================== 账号分组管理 ====================

// CreateAccountGroup 创建账号分组
func (h *Handler) CreateAccountGroup(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Name        string `json:"name" vd:"len($)>0"`
		GroupType   string `json:"group_type"`
		Description string `json:"description"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.accountClient.CreateAccountGroup(ctx, userID, req.Name, req.GroupType, req.Description)
	if err != nil {
		errResponse(c, err, "create group failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetAccountGroup 获取账号分组
func (h *Handler) GetAccountGroup(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid group id"))
		return
	}

	result, err := h.accountClient.GetAccountGroup(ctx, id)
	if err != nil {
		errResponse(c, err, "get group failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// UpdateAccountGroup 更新账号分组
func (h *Handler) UpdateAccountGroup(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid group id"))
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

	result, err := h.accountClient.UpdateAccountGroup(ctx, id, req.Name, req.Description)
	if err != nil {
		errResponse(c, err, "update group failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// DeleteAccountGroup 删除账号分组
func (h *Handler) DeleteAccountGroup(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid group id"))
		return
	}

	if err := h.accountClient.DeleteAccountGroup(ctx, id); err != nil {
		errResponse(c, err, "delete group failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListAccountGroups 列出账号分组
func (h *Handler) ListAccountGroups(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))
	userID := c.GetInt64("user_id")

	result, err := h.accountClient.ListAccountGroups(ctx, userID, page, pageSize)
	if err != nil {
		errResponse(c, err, "list groups failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// AddAccountToGroup 添加账号到分组
func (h *Handler) AddAccountToGroup(ctx context.Context, c *app.RequestContext) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid group id"))
		return
	}

	var req struct {
		AccountID int64 `json:"account_id" vd:"$>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	if err := h.accountClient.AddAccountToGroup(ctx, groupID, req.AccountID); err != nil {
		errResponse(c, err, "add account to group failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// RemoveAccountFromGroup 从分组移除账号
func (h *Handler) RemoveAccountFromGroup(ctx context.Context, c *app.RequestContext) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid group id"))
		return
	}

	accountID, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid account id"))
		return
	}

	if err := h.accountClient.RemoveAccountFromGroup(ctx, groupID, accountID); err != nil {
		errResponse(c, err, "remove account from group failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}
