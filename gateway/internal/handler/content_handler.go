package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// CreateContent 创建内容
func (h *Handler) CreateContent(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Title        string `json:"title" vd:"len($)>0"`
		Body         string `json:"body" vd:"len($)>0"`
		ContentType  string `json:"content_type"`
		SchemaMarkup string `json:"schema_markup"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	resp, err := h.contentClient.CreateContent(ctx, userID, req.Title, req.Body, req.ContentType, req.SchemaMarkup)
	if err != nil {
		errResponse(c, err, "create content failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// GetContent 获取内容
func (h *Handler) GetContent(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid content id"))
		return
	}

	resp, err := h.contentClient.GetContent(ctx, id)
	if err != nil {
		errResponse(c, err, "get content failed")
		return
	}

	if !checkOwnership(c, resp) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// UpdateContent 更新内容
func (h *Handler) UpdateContent(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid content id"))
		return
	}

	existing, err := h.contentClient.GetContent(ctx, id)
	if err != nil {
		errResponse(c, err, "content not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	var req struct {
		Title        string `json:"title"`
		Body         string `json:"body"`
		SchemaMarkup string `json:"schema_markup"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.contentClient.UpdateContent(ctx, id, req.Title, req.Body, req.SchemaMarkup)
	if err != nil {
		errResponse(c, err, "update content failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// DeleteContent 删除内容
func (h *Handler) DeleteContent(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid content id"))
		return
	}

	existing, err := h.contentClient.GetContent(ctx, id)
	if err != nil {
		errResponse(c, err, "content not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	if err := h.contentClient.DeleteContent(ctx, id); err != nil {
		errResponse(c, err, "delete content failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListContents 列出内容
func (h *Handler) ListContents(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	page, pageSize := parsePagination(c)
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	contentType := c.Query("content_type")

	resp, err := h.contentClient.ListContents(ctx, userID, page, pageSize, status, contentType)
	if err != nil {
		errResponse(c, err, "list contents failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// OptimizeContent 优化内容
func (h *Handler) OptimizeContent(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid content id"))
		return
	}

	var req struct {
		OptimizationType string `json:"optimization_type"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.contentClient.OptimizeContent(ctx, id, req.OptimizationType)
	if err != nil {
		errResponse(c, err, "optimize content failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// PublishContent 发布内容
func (h *Handler) PublishContent(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid content id"))
		return
	}

	var req struct {
		ChannelIDs []int64 `json:"channel_ids"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	resp, err := h.contentClient.PublishContent(ctx, id, req.ChannelIDs)
	if err != nil {
		errResponse(c, err, "publish content failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}

// CheckCompliance 检查合规性
func (h *Handler) CheckCompliance(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid content id"))
		return
	}

	resp, err := h.contentClient.CheckCompliance(ctx, id)
	if err != nil {
		errResponse(c, err, "check compliance failed")
		return
	}

	c.JSON(http.StatusOK, success(resp))
}
