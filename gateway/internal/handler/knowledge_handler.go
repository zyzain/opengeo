package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// CreateEntity 创建知识实体
func (h *Handler) CreateEntity(ctx context.Context, c *app.RequestContext) {
	var req struct {
		EntityName     string `json:"entity_name" vd:"len($)>0"`
		EntityType     string `json:"entity_type" vd:"len($)>0"`
		EntityData     string `json:"entity_data"`
		AuthorityLinks string `json:"authority_links"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.knowledgeClient.CreateEntity(ctx, userID, req.EntityName, req.EntityType, req.EntityData, req.AuthorityLinks)
	if err != nil {
		errResponse(c, err, "create entity failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetEntity 获取知识实体
func (h *Handler) GetEntity(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid entity id"))
		return
	}

	result, err := h.knowledgeClient.GetEntity(ctx, id)
	if err != nil {
		errResponse(c, err, "get entity failed")
		return
	}

	if !checkOwnership(c, result) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// UpdateEntity 更新知识实体
func (h *Handler) UpdateEntity(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid entity id"))
		return
	}

	existing, err := h.knowledgeClient.GetEntity(ctx, id)
	if err != nil {
		errResponse(c, err, "entity not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	var req struct {
		EntityName     string `json:"entity_name"`
		EntityType     string `json:"entity_type"`
		EntityData     string `json:"entity_data"`
		AuthorityLinks string `json:"authority_links"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	result, err := h.knowledgeClient.UpdateEntity(ctx, id, req.EntityName, req.EntityType, req.EntityData, req.AuthorityLinks)
	if err != nil {
		errResponse(c, err, "update entity failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// DeleteEntity 删除知识实体
func (h *Handler) DeleteEntity(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid entity id"))
		return
	}

	existing, err := h.knowledgeClient.GetEntity(ctx, id)
	if err != nil {
		errResponse(c, err, "entity not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	if err := h.knowledgeClient.DeleteEntity(ctx, id); err != nil {
		errResponse(c, err, "delete entity failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListEntities 列出知识实体
func (h *Handler) ListEntities(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	entityType := c.Query("entity_type")
	page, pageSize := parsePagination(c)

	result, err := h.knowledgeClient.ListEntities(ctx, userID, entityType, page, pageSize)
	if err != nil {
		errResponse(c, err, "list entities failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// SearchEntities 搜索知识实体
func (h *Handler) SearchEntities(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	keyword := c.Query("keyword")

	result, err := h.knowledgeClient.SearchEntities(ctx, userID, keyword)
	if err != nil {
		errResponse(c, err, "search entities failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}
