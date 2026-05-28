package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// CreateSchedule 创建调度
func (h *Handler) CreateSchedule(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ScheduleName   string `json:"schedule_name" vd:"len($)>0"`
		ScheduleType   string `json:"schedule_type" vd:"len($)>0"`
		CronExpression string `json:"cron_expression"`
		Config         string `json:"config"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	userID := c.GetInt64("user_id")
	result, err := h.scheduleClient.CreateSchedule(ctx, userID, req.ScheduleName, req.ScheduleType, req.CronExpression, req.Config)
	if err != nil {
		errResponse(c, err, "create schedule failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetSchedule 获取调度
func (h *Handler) GetSchedule(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid schedule id"))
		return
	}

	result, err := h.scheduleClient.GetSchedule(ctx, id)
	if err != nil {
		errResponse(c, err, "get schedule failed")
		return
	}

	if !checkOwnership(c, result) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// UpdateSchedule 更新调度
func (h *Handler) UpdateSchedule(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid schedule id"))
		return
	}

	existing, err := h.scheduleClient.GetSchedule(ctx, id)
	if err != nil {
		errResponse(c, err, "schedule not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	var req struct {
		ScheduleName   string `json:"schedule_name"`
		CronExpression string `json:"cron_expression"`
		Config         string `json:"config"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid request parameters"))
		return
	}

	result, err := h.scheduleClient.UpdateSchedule(ctx, id, req.ScheduleName, req.CronExpression, req.Config)
	if err != nil {
		errResponse(c, err, "update schedule failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// DeleteSchedule 删除调度
func (h *Handler) DeleteSchedule(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid schedule id"))
		return
	}

	existing, err := h.scheduleClient.GetSchedule(ctx, id)
	if err != nil {
		errResponse(c, err, "schedule not found")
		return
	}
	if !checkOwnership(c, existing) {
		c.JSON(http.StatusForbidden, fail(403, "forbidden"))
		return
	}

	if err := h.scheduleClient.DeleteSchedule(ctx, id); err != nil {
		errResponse(c, err, "delete schedule failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// ListSchedules 列出调度
func (h *Handler) ListSchedules(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	page, pageSize := parsePagination(c)

	result, err := h.scheduleClient.ListSchedules(ctx, userID, page, pageSize)
	if err != nil {
		errResponse(c, err, "list schedules failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// EnableSchedule 启用调度
func (h *Handler) EnableSchedule(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid schedule id"))
		return
	}

	if err := h.scheduleClient.EnableSchedule(ctx, id); err != nil {
		errResponse(c, err, "enable schedule failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}

// DisableSchedule 禁用调度
func (h *Handler) DisableSchedule(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(400, "invalid schedule id"))
		return
	}

	if err := h.scheduleClient.DisableSchedule(ctx, id); err != nil {
		errResponse(c, err, "disable schedule failed")
		return
	}

	c.JSON(http.StatusOK, success(utils.H{"success": true}))
}
