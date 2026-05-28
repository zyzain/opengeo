package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	"opengeo/gateway/internal/model"
)

// ==================== 系统管理 ====================

// GetSystemConfigs 获取系统配置
func (h *Handler) GetSystemConfigs(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	result, err := h.systemClient.GetSystemConfigs(ctx, page, pageSize)
	if err != nil {
		errResponse(c, err, "get system configs failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// UpdateSystemConfig 更新系统配置
func (h *Handler) UpdateSystemConfig(ctx context.Context, c *app.RequestContext) {
	key := c.Param("key")

	var req struct {
		ConfigValue string `json:"config_value"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	result, err := h.systemClient.UpdateSystemConfig(ctx, key, req.ConfigValue)
	if err != nil {
		errResponse(c, err, "update system config failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// ==================== 插件管理 ====================

// ListPlugins 列出插件
func (h *Handler) ListPlugins(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	result, err := h.systemClient.ListPlugins(ctx, page, pageSize)
	if err != nil {
		errResponse(c, err, "list plugins failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// InstallPlugin 安装插件
func (h *Handler) InstallPlugin(ctx context.Context, c *app.RequestContext) {
	var req struct {
		PluginName string `json:"plugin_name" vd:"len($)>0"`
		PluginType string `json:"plugin_type" vd:"len($)>0"`
		Version    string `json:"version"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{
		"plugin_name": req.PluginName,
		"plugin_type": req.PluginType,
		"version":     req.Version,
		"is_enabled":  true,
		"message":     "插件安装成功",
	}))
}

// UpdatePlugin 更新插件
func (h *Handler) UpdatePlugin(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid plugin id"))
		return
	}

	var req struct {
		IsEnabled bool `json:"is_enabled"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{
		"id":         id,
		"is_enabled": req.IsEnabled,
	}))
}

// DeletePlugin 删除插件
func (h *Handler) DeletePlugin(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid plugin id"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true, "id": id}))
}

// ==================== Webhook管理 ====================

// ListWebhooks 列出Webhook
func (h *Handler) ListWebhooks(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	result, err := h.systemClient.ListWebhooks(ctx, userID, page, pageSize)
	if err != nil {
		errResponse(c, err, "list webhooks failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// CreateWebhook 创建Webhook
func (h *Handler) CreateWebhook(ctx context.Context, c *app.RequestContext) {
	var req struct {
		WebhookName string   `json:"webhook_name" vd:"len($)>0"`
		URL         string   `json:"url" vd:"len($)>0"`
		Events      []string `json:"events"`
		Secret      string   `json:"secret"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.systemClient.CreateWebhook(ctx, userID, req.WebhookName, req.URL, req.Secret, req.Events)
	if err != nil {
		errResponse(c, err, "create webhook failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// UpdateWebhook 更新Webhook
func (h *Handler) UpdateWebhook(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid webhook id"))
		return
	}

	var req struct {
		WebhookName string   `json:"webhook_name"`
		URL         string   `json:"url"`
		Events      []string `json:"events"`
		Secret      string   `json:"secret"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	result, err := h.systemClient.UpdateWebhook(ctx, id, req.WebhookName, req.URL, req.Secret, req.Events)
	if err != nil {
		errResponse(c, err, "update webhook failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// DeleteWebhook 删除Webhook
func (h *Handler) DeleteWebhook(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid webhook id"))
		return
	}

	if err := h.systemClient.DeleteWebhook(ctx, id); err != nil {
		errResponse(c, err, "delete webhook failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// TestWebhook 测试Webhook
func (h *Handler) TestWebhook(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid webhook id"))
		return
	}

	result, err := h.systemClient.TestWebhook(ctx, id)
	if err != nil {
		errResponse(c, err, "test webhook failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetWebhookHistory 获取Webhook触发历史
func (h *Handler) GetWebhookHistory(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid webhook id"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	result, err := h.systemClient.GetWebhookHistory(ctx, id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "get webhook history failed"))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// ==================== 模板管理 ====================

// ListTemplates 列出模板
func (h *Handler) ListTemplates(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	templateType := c.Query("template_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	tpls, total, err := h.tplRepo.List(ctx, userID, templateType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "list templates failed"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{
		"items":     tpls,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// CreateTemplate 创建模板
func (h *Handler) CreateTemplate(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Name         string `json:"name" vd:"len($)>0"`
		Description  string `json:"description"`
		TemplateType string `json:"template_type" vd:"len($)>0"`
		ChannelType  string `json:"channel_type"`
		Content      string `json:"content" vd:"len($)>0"`
		Variables    string `json:"variables"`
		IsPublic     *bool  `json:"is_public"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	userID := c.GetInt64("user_id")
	isPublic := false
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	tpl := &model.ContentTemplate{
		UserID:       userID,
		Name:         req.Name,
		TemplateType: req.TemplateType,
		ChannelType:  req.ChannelType,
		Content:      req.Content,
		Variables:    req.Variables,
		Description:  req.Description,
		IsPublic:     isPublic,
		IsEnabled:    true,
	}

	if err := h.tplRepo.Create(ctx, tpl); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "create template failed"))
		return
	}

	c.JSON(http.StatusOK, success(tpl))
}

// GetTemplate 获取模板
func (h *Handler) GetTemplate(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid template id"))
		return
	}

	tpl, err := h.tplRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, fail(404, "template not found"))
		return
	}

	c.JSON(http.StatusOK, success(tpl))
}

// UpdateTemplate 更新模板
func (h *Handler) UpdateTemplate(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid template id"))
		return
	}

	var req struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		TemplateType string `json:"template_type"`
		ChannelType  string `json:"channel_type"`
		Content      string `json:"content"`
		Variables    string `json:"variables"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	tpl, err := h.tplRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, fail(404, "template not found"))
		return
	}

	if req.Name != "" {
		tpl.Name = req.Name
	}
	if req.Description != "" {
		tpl.Description = req.Description
	}
	if req.TemplateType != "" {
		tpl.TemplateType = req.TemplateType
	}
	if req.ChannelType != "" {
		tpl.ChannelType = req.ChannelType
	}
	if req.Content != "" {
		tpl.Content = req.Content
	}
	if req.Variables != "" {
		tpl.Variables = req.Variables
	}

	if err := h.tplRepo.Update(ctx, tpl); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "update template failed"))
		return
	}

	c.JSON(http.StatusOK, success(tpl))
}

// DeleteTemplate 删除模板
func (h *Handler) DeleteTemplate(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid template id"))
		return
	}

	if err := h.tplRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "delete template failed"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true, "id": id}))
}

// ==================== 套餐与权限定义 ====================

// GetPlans 获取套餐配置
func (h *Handler) GetPlans(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, success(utils.H{
		"items": []map[string]interface{}{
			{
				"id":           "starter",
				"label":        "入门版",
				"color":        "blue",
				"max_users":    5,
				"max_storage":  2,
				"max_contents": 500,
				"price":        0,
			},
			{
				"id":           "professional",
				"label":        "专业版",
				"color":        "green",
				"max_users":    20,
				"max_storage":  10,
				"max_contents": 5000,
				"price":        299,
			},
			{
				"id":           "enterprise",
				"label":        "企业版",
				"color":        "purple",
				"max_users":    100,
				"max_storage":  100,
				"max_contents": 50000,
				"price":        999,
			},
		},
	}))
}

// GetPermissionDefinitions 获取权限定义
func (h *Handler) GetPermissionDefinitions(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, success(utils.H{
		"groups": []map[string]interface{}{
			{
				"title": "内容管理",
				"permissions": []map[string]interface{}{
					{"id": "content:create", "label": "创建内容"},
					{"id": "content:read", "label": "查看内容"},
					{"id": "content:update", "label": "编辑内容"},
					{"id": "content:delete", "label": "删除内容"},
					{"id": "content:publish", "label": "发布内容"},
					{"id": "content:optimize", "label": "AI优化内容"},
				},
			},
			{
				"title": "账号管理",
				"permissions": []map[string]interface{}{
					{"id": "account:create", "label": "创建账号"},
					{"id": "account:read", "label": "查看账号"},
					{"id": "account:update", "label": "编辑账号"},
					{"id": "account:delete", "label": "删除账号"},
				},
			},
			{
				"title": "发布管理",
				"permissions": []map[string]interface{}{
					{"id": "publish:create", "label": "创建发布任务"},
					{"id": "publish:read", "label": "查看发布任务"},
					{"id": "publish:cancel", "label": "取消发布任务"},
					{"id": "publish:retry", "label": "重试发布任务"},
				},
			},
			{
				"title": "监测管理",
				"permissions": []map[string]interface{}{
					{"id": "monitor:read", "label": "查看监测数据"},
					{"id": "monitor:configure", "label": "配置监测规则"},
				},
			},
			{
				"title": "系统管理",
				"permissions": []map[string]interface{}{
					{"id": "system:config", "label": "系统配置"},
					{"id": "system:user", "label": "用户管理"},
					{"id": "system:role", "label": "角色管理"},
					{"id": "system:plugin", "label": "插件管理"},
				},
			},
		},
	}))
}
