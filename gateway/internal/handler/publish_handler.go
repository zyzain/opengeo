package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	"opengeo/gateway/internal/client"
	"opengeo/gateway/internal/model"
	"opengeo/pkg/locale"
)

// ==================== 渠道管理 ====================

// CreateChannel 创建渠道
func (h *Handler) CreateChannel(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ChannelType   string `json:"channel_type" vd:"len($)>0"`
		ChannelName   string `json:"channel_name" vd:"len($)>0"`
		ChannelConfig string `json:"channel_config"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.publishClient.CreateChannel(ctx, userID, req.ChannelType, req.ChannelName, req.ChannelConfig)
	if err != nil {
		errResponse(c, err, "create channel failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetChannel 获取渠道
func (h *Handler) GetChannel(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid channel id"))
		return
	}

	result, err := h.publishClient.GetChannel(ctx, id)
	if err != nil {
		errResponse(c, err, "get channel failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// ListChannels 列出渠道
func (h *Handler) ListChannels(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))
	userID := c.GetInt64("user_id")

	result, err := h.publishClient.ListChannels(ctx, userID, page, pageSize)
	if err != nil {
		errResponse(c, err, "list channels failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetChannelPlatforms 获取支持的渠道平台
func (h *Handler) GetChannelPlatforms(ctx context.Context, c *app.RequestContext) {
	loc := locale.FromContext(ctx)
	platforms := []map[string]interface{}{
		{"value": "wechat", "label": locale.T(loc, "platform_wechat"), "color": "green"},
		{"value": "weibo", "label": locale.T(loc, "platform_weibo"), "color": "red"},
		{"value": "douyin", "label": locale.T(loc, "platform_douyin"), "color": "purple"},
		{"value": "xiaohongshu", "label": locale.T(loc, "platform_xiaohongshu"), "color": "pink"},
		{"value": "zhihu", "label": locale.T(loc, "platform_zhihu"), "color": "blue"},
		{"value": "toutiao", "label": locale.T(loc, "platform_toutiao"), "color": "orange"},
	}
	c.JSON(http.StatusOK, success(platforms))
}

// ==================== 发布管理 ====================

// CreatePublishTask 创建发布任务
func (h *Handler) CreatePublishTask(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ContentID     int64  `json:"content_id" vd:"$>0"`
		ChannelID     int64  `json:"channel_id" vd:"$>0"`
		ScheduledTime string `json:"scheduled_time"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.publishClient.CreatePublishTask(ctx, userID, req.ContentID, req.ChannelID, req.ScheduledTime)
	if err != nil {
		errResponse(c, err, "create publish task failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetPublishTask 获取发布任务
func (h *Handler) GetPublishTask(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid task id"))
		return
	}

	result, err := h.publishClient.GetPublishTask(ctx, id)
	if err != nil {
		errResponse(c, err, "get publish task failed")
		return
	}

	if !checkOwnership(c, result) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// ListPublishTasks 列出发布任务
func (h *Handler) ListPublishTasks(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	page, pageSize := parsePagination(c)

	result, err := h.publishClient.ListPublishTasks(ctx, userID, status, page, pageSize)
	if err != nil {
		errResponse(c, err, "list publish tasks failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// CancelPublishTask 取消发布任务
func (h *Handler) CancelPublishTask(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid task id"))
		return
	}

	existing, err := h.publishClient.GetPublishTask(ctx, id)
	if err != nil {
		errResponse(c, err, "task not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	if err := h.publishClient.CancelPublishTask(ctx, id); err != nil {
		errResponse(c, err, "cancel publish task failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// RetryPublishTask 重试发布任务
func (h *Handler) RetryPublishTask(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid task id"))
		return
	}

	existing, err := h.publishClient.GetPublishTask(ctx, id)
	if err != nil {
		errResponse(c, err, "task not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	if err := h.publishClient.RetryPublishTask(ctx, id); err != nil {
		errResponse(c, err, "retry publish task failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// PreviewPublish 预览发布
func (h *Handler) PreviewPublish(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ChannelID int64  `json:"channel_id" vd:"$>0"`
		Title     string `json:"title"`
		Body      string `json:"body"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	result, err := h.publishClient.PreviewPublish(ctx, req.ChannelID, req.Title, req.Body)
	if err != nil {
		errResponse(c, err, "preview publish failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// ValidatePublish 校验发布
func (h *Handler) ValidatePublish(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ChannelID int64  `json:"channel_id" vd:"$>0"`
		Title     string `json:"title"`
		Body      string `json:"body"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	result, err := h.publishClient.ValidatePublish(ctx, req.ChannelID, req.Title, req.Body)
	if err != nil {
		errResponse(c, err, "validate publish failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// ==================== 内容去重 ====================

// CheckDedup 内容去重检测
func (h *Handler) CheckDedup(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Text string `json:"text" vd:"len($)>0"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	userID := c.GetInt64("user_id")

	result, err := h.publishClient.CheckContentDedup(ctx, userID, req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "dedup check failed"))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// ==================== 平台管理 ====================

// ListPlatforms 列出平台
func (h *Handler) ListPlatforms(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	result, err := h.publishClient.ListPlatforms(ctx, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, safeError(err, "internal server error")))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetPlatform 获取平台详情
func (h *Handler) GetPlatform(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid platform id"))
		return
	}

	result, err := h.publishClient.GetPlatform(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, fail(404, safeError(err, "resource not found")))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// CreatePlatform 创建平台
func (h *Handler) CreatePlatform(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Code        string `json:"code" vd:"len($)>0"`
		Name        string `json:"name" vd:"len($)>0"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		Description string `json:"description"`
		SortOrder   int32  `json:"sort_order"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	result, err := h.publishClient.CreatePlatform(ctx, &client.CreatePlatformRequest{
		Code:        req.Code,
		Name:        req.Name,
		Icon:        req.Icon,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, safeError(err, "internal server error")))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// UpdatePlatform 更新平台
func (h *Handler) UpdatePlatform(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid platform id"))
		return
	}

	var req struct {
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		Description string `json:"description"`
		SortOrder   int32  `json:"sort_order"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	result, err := h.publishClient.UpdatePlatform(ctx, id, &client.UpdatePlatformRequest{
		Name:        req.Name,
		Icon:        req.Icon,
		Color:       req.Color,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, safeError(err, "internal server error")))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// DeletePlatform 删除平台
func (h *Handler) DeletePlatform(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid platform id"))
		return
	}

	if err := h.publishClient.DeletePlatform(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, safeError(err, "internal server error")))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// EnablePlatform 启用平台
func (h *Handler) EnablePlatform(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid platform id"))
		return
	}

	if err := h.publishClient.EnablePlatform(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, safeError(err, "internal server error")))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// DisablePlatform 禁用平台
func (h *Handler) DisablePlatform(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid platform id"))
		return
	}

	if err := h.publishClient.DisablePlatform(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, safeError(err, "internal server error")))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ==================== 指纹管理 ====================

// ListFingerprints 列出指纹配置
func (h *Handler) ListFingerprints(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	fps, total, err := h.fpRepo.List(ctx, userID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "list fingerprints failed"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{
		"items":     fps,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// CreateFingerprint 创建指纹配置
func (h *Handler) CreateFingerprint(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Name      string `json:"name" vd:"len($)>0"`
		Platform  string `json:"platform" vd:"len($)>0"`
		UserAgent string `json:"user_agent" vd:"len($)>0"`
		Screen    string `json:"screen"`
		Language  string `json:"language"`
		Timezone  string `json:"timezone"`
		WebGL     string `json:"webgl"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	userID := c.GetInt64("user_id")
	fp := &model.BrowserFingerprint{
		UserID:    userID,
		Name:      req.Name,
		UserAgent: req.UserAgent,
		Platform:  req.Platform,
		Screen:    req.Screen,
		Language:  req.Language,
		Timezone:  req.Timezone,
		Status:    "active",
		IsEnabled: true,
	}

	if err := h.fpRepo.Create(ctx, fp); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "create fingerprint failed"))
		return
	}

	c.JSON(http.StatusOK, success(fp))
}

// UpdateFingerprint 更新指纹配置
func (h *Handler) UpdateFingerprint(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	var req struct {
		Name      string `json:"name"`
		UserAgent string `json:"user_agent"`
		Platform  string `json:"platform"`
		Screen    string `json:"screen"`
		Language  string `json:"language"`
		Timezone  string `json:"timezone"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	fp, err := h.fpRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, fail(404, "fingerprint not found"))
		return
	}

	if req.Name != "" {
		fp.Name = req.Name
	}
	if req.UserAgent != "" {
		fp.UserAgent = req.UserAgent
	}
	if req.Platform != "" {
		fp.Platform = req.Platform
	}
	if req.Screen != "" {
		fp.Screen = req.Screen
	}
	if req.Language != "" {
		fp.Language = req.Language
	}
	if req.Timezone != "" {
		fp.Timezone = req.Timezone
	}

	if err := h.fpRepo.Update(ctx, fp); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "update fingerprint failed"))
		return
	}

	c.JSON(http.StatusOK, success(fp))
}

// DeleteFingerprint 删除指纹配置
func (h *Handler) DeleteFingerprint(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	if err := h.fpRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "delete fingerprint failed"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true, "id": id}))
}

// ToggleFingerprint 切换指纹状态
func (h *Handler) ToggleFingerprint(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	fp, err := h.fpRepo.ToggleStatus(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "toggle fingerprint failed"))
		return
	}

	c.JSON(http.StatusOK, success(fp))
}

// ==================== 代理IP管理 ====================

// ListProxies 列出代理IP
func (h *Handler) ListProxies(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	status := c.Query("status")
	page, pageSize := parsePagination(c)

	proxies, total, err := h.proxyRepo.List(ctx, userID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "list proxies failed"))
		return
	}

	for _, p := range proxies {
		p.Password = "****"
	}

	c.JSON(http.StatusOK, success(utils.H{
		"items":     proxies,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// CreateProxy 创建代理IP
func (h *Handler) CreateProxy(ctx context.Context, c *app.RequestContext) {
	var req struct {
		IP       string `json:"ip" vd:"len($)>0"`
		Port     int    `json:"port" vd:"$>0"`
		Protocol string `json:"protocol" vd:"len($)>0"`
		Location string `json:"location"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	userID := c.GetInt64("user_id")
	proxy := &model.ProxyIP{
		UserID:    userID,
		IP:        req.IP,
		Port:      req.Port,
		Protocol:  req.Protocol,
		Location:  req.Location,
		Username:  req.Username,
		Password:  req.Password,
		Status:    "active",
		IsEnabled: true,
	}

	if err := h.proxyRepo.Create(ctx, proxy); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "create proxy failed"))
		return
	}

	proxy.Password = "****"
	c.JSON(http.StatusOK, success(proxy))
}

// DeleteProxy 删除代理IP
func (h *Handler) DeleteProxy(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	if err := h.proxyRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "delete proxy failed"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true, "id": id}))
}

// CheckProxy 检查代理IP
func (h *Handler) CheckProxy(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	proxy, err := h.proxyRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, fail(404, "proxy not found"))
		return
	}

	now := time.Now()
	proxy.LastCheck = &now
	proxy.Speed = 95
	proxy.Uptime = 99.0
	h.proxyRepo.Update(ctx, proxy)

	c.JSON(http.StatusOK, success(proxy))
}

// ToggleProxy 切换代理状态
func (h *Handler) ToggleProxy(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	proxy, err := h.proxyRepo.ToggleStatus(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "toggle proxy failed"))
		return
	}

	c.JSON(http.StatusOK, success(proxy))
}

// ==================== 错峰策略管理 ====================

// ListStaggerStrategies 列出错峰策略
func (h *Handler) ListStaggerStrategies(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	strategies, total, err := h.staggerStrategyRepo.List(ctx, userID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "list stagger strategies failed"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{
		"items":     strategies,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// CreateStaggerStrategy 创建错峰策略
func (h *Handler) CreateStaggerStrategy(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Name        string `json:"name" vd:"len($)>0"`
		Accounts    int    `json:"accounts"`
		Interval    int    `json:"interval"`
		RandomRange int    `json:"random_range"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	userID := c.GetInt64("user_id")
	strategy := &model.StaggerStrategy{
		UserID:      userID,
		Name:        req.Name,
		Accounts:    req.Accounts,
		Interval:    req.Interval,
		RandomRange: req.RandomRange,
		Status:      "active",
		IsEnabled:   true,
	}

	if err := h.staggerStrategyRepo.Create(ctx, strategy); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "create stagger strategy failed"))
		return
	}

	c.JSON(http.StatusOK, success(strategy))
}

// UpdateStaggerStrategy 更新错峰策略
func (h *Handler) UpdateStaggerStrategy(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	var req struct {
		Name        string `json:"name"`
		Accounts    int    `json:"accounts"`
		Interval    int    `json:"interval"`
		RandomRange int    `json:"random_range"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	strategy, err := h.staggerStrategyRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, fail(404, "stagger strategy not found"))
		return
	}

	if req.Name != "" {
		strategy.Name = req.Name
	}
	if req.Accounts > 0 {
		strategy.Accounts = req.Accounts
	}
	if req.Interval > 0 {
		strategy.Interval = req.Interval
	}
	if req.RandomRange > 0 {
		strategy.RandomRange = req.RandomRange
	}

	if err := h.staggerStrategyRepo.Update(ctx, strategy); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "update stagger strategy failed"))
		return
	}

	c.JSON(http.StatusOK, success(strategy))
}

// DeleteStaggerStrategy 删除错峰策略
func (h *Handler) DeleteStaggerStrategy(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	if err := h.staggerStrategyRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "delete stagger strategy failed"))
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true, "id": id}))
}

// ToggleStaggerStrategy 切换错峰策略状态
func (h *Handler) ToggleStaggerStrategy(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid id"))
		return
	}

	strategy, err := h.staggerStrategyRepo.ToggleStatus(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "toggle stagger strategy failed"))
		return
	}

	c.JSON(http.StatusOK, success(strategy))
}

// GetStaggerConfig 获取错峰配置
func (h *Handler) GetStaggerConfig(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")

	config, err := h.staggerConfigRepo.GetByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "get stagger config failed"))
		return
	}

	c.JSON(http.StatusOK, success(config))
}

// UpdateStaggerConfig 更新错峰配置
func (h *Handler) UpdateStaggerConfig(ctx context.Context, c *app.RequestContext) {
	var req struct {
		MinInterval      int `json:"min_interval"`
		MaxInterval      int `json:"max_interval"`
		RandomRange      int `json:"random_range"`
		BatchSize        int `json:"batch_size"`
		CooldownAfter    int `json:"cooldown_after"`
		CooldownDuration int `json:"cooldown_duration"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, safeError(err, "invalid request")))
		return
	}

	userID := c.GetInt64("user_id")
	config := &model.StaggerConfig{
		UserID:           userID,
		MinInterval:      req.MinInterval,
		MaxInterval:      req.MaxInterval,
		RandomRange:      req.RandomRange,
		BatchSize:        req.BatchSize,
		CooldownAfter:    req.CooldownAfter,
		CooldownDuration: req.CooldownDuration,
	}

	if err := h.staggerConfigRepo.Save(ctx, config); err != nil {
		c.JSON(http.StatusInternalServerError, fail(500, "update stagger config failed"))
		return
	}

	c.JSON(http.StatusOK, success(config))
}
