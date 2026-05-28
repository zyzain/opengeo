package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// GetAICitations 获取AI引用
func (h *Handler) GetAICitations(ctx context.Context, c *app.RequestContext) {
	contentID, _ := strconv.ParseInt(c.Query("content_id"), 10, 64)
	aiModel := c.Query("ai_model")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	result, err := h.monitorClient.GetAICitations(ctx, contentID, aiModel, page, pageSize)
	if err != nil {
		errResponse(c, err, "get AI citations failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetSourceScores 获取信源评分
func (h *Handler) GetSourceScores(ctx context.Context, c *app.RequestContext) {
	channelID, _ := strconv.ParseInt(c.Query("channel_id"), 10, 64)
	accountID, _ := strconv.ParseInt(c.Query("account_id"), 10, 64)

	result, err := h.monitorClient.GetSourceScores(ctx, channelID, accountID)
	if err != nil {
		errResponse(c, err, "get source scores failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetCompetitorMonitors 获取竞品监测
func (h *Handler) GetCompetitorMonitors(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	result, err := h.monitorClient.GetCompetitorMonitors(ctx, userID, page, pageSize)
	if err != nil {
		errResponse(c, err, "get competitor monitors failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GetOurScore 获取我们的评分
func (h *Handler) GetOurScore(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, success(utils.H{
		"score":     92.5,
		"rank":      1,
		"trend":     "up",
		"trend_pct": 5.2,
	}))
}

// GetROIMetrics 获取ROI指标
func (h *Handler) GetROIMetrics(ctx context.Context, c *app.RequestContext) {
	contentID, _ := strconv.ParseInt(c.Query("content_id"), 10, 64)
	channelID, _ := strconv.ParseInt(c.Query("channel_id"), 10, 64)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := h.monitorClient.GetROIMetrics(ctx, contentID, channelID, startDate, endDate)
	if err != nil {
		errResponse(c, err, "get ROI metrics failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}

// GenerateSuggestions 生成优化建议
func (h *Handler) GenerateSuggestions(ctx context.Context, c *app.RequestContext) {
	contentID, _ := strconv.ParseInt(c.Query("content_id"), 10, 64)

	result, err := h.monitorClient.GenerateSuggestions(ctx, contentID)
	if err != nil {
		errResponse(c, err, "generate suggestions failed")
		return
	}

	c.JSON(http.StatusOK, success(result))
}
