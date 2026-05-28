package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	"opengeo/gateway/internal/client"
	"opengeo/gateway/internal/dal"
	"opengeo/pkg/config"
	"opengeo/pkg/errcode"
)

func success(data interface{}) utils.H {
	return utils.H{
		"code":    0,
		"message": "success",
		"data":    data,
	}
}

func fail(code int, msg string) utils.H {
	return utils.H{
		"code":    code,
		"message": msg,
		"data":    nil,
	}
}

// safeError 返回安全的错误信息（不泄露内部细节）
func safeError(err error, fallback string) string {
	errMsg := err.Error()

	// AppError 直接返回消息
	if appErr, ok := err.(*errcode.AppError); ok {
		return appErr.Message
	}

	// 已知的业务错误直接返回
	knownErrors := []string{
		"invalid username or password",
		"username already exists",
		"user account is disabled",
		"token expired",
		"invalid token",
		"missing authorization header",
		"channel not found",
		"account not found",
		"platform not found",
		"rate limit exceeded",
	}

	for _, known := range knownErrors {
		if strings.Contains(errMsg, known) {
			return errMsg
		}
	}

	// 数据库错误返回通用信息
	if strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "Duplicate") {
		return "data already exists"
	}
	if strings.Contains(errMsg, "not found") {
		return "resource not found"
	}
	if strings.Contains(errMsg, "foreign key") {
		return "invalid reference"
	}

	// 其他错误返回通用信息
	return fallback
}

// errToHTTPStatus 将错误映射到 HTTP 状态码
func errToHTTPStatus(err error) int {
	if appErr, ok := err.(*errcode.AppError); ok {
		switch {
		case appErr.Code >= 20000 && appErr.Code < 30000:
			return http.StatusBadRequest
		case appErr.Code == errcode.Unauthorized:
			return http.StatusUnauthorized
		case appErr.Code == errcode.Forbidden:
			return http.StatusForbidden
		case appErr.Code == errcode.NotFound:
			return http.StatusNotFound
		case appErr.Code == errcode.AlreadyExists:
			return http.StatusConflict
		case appErr.Code == errcode.RateLimitExceeded:
			return http.StatusTooManyRequests
		default:
			return http.StatusInternalServerError
		}
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "not found") {
		return http.StatusNotFound
	}
	if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "duplicate") {
		return http.StatusConflict
	}
	if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "required") {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

// errResponse 统一错误响应
func errResponse(c *app.RequestContext, err error, fallbackMsg string) {
	status := errToHTTPStatus(err)
	code := int(status)
	msg := safeError(err, fallbackMsg)

	if appErr, ok := err.(*errcode.AppError); ok {
		code = int(appErr.Code)
	}

	c.JSON(status, fail(code, msg))
}

// getDefaultPageSize 获取默认分页大小
func getDefaultPageSize() string {
	cfg := config.GetConfig()
	return strconv.Itoa(cfg.Pagination.DefaultPageSize)
}

// parsePagination 解析分页参数并校验上限
func parsePagination(c *app.RequestContext) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", getDefaultPageSize()))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// checkOwnership 校验资源归属权
// admin 角色自动放行，其他用户只能操作自己的资源
func checkOwnership(c *app.RequestContext, resource map[string]interface{}) bool {
	role, _ := c.Get("role")
	if roleStr, ok := role.(string); ok && roleStr == "admin" {
		return true
	}

	currentUserID, _ := c.Get("user_id")
	uid, ok := currentUserID.(int64)
	if !ok {
		return false
	}

	if ownerID, exists := resource["user_id"]; exists {
		if id, ok := ownerID.(int64); ok {
			return id == uid
		}
	}

	return false
}

// Handler HTTP处理器
type Handler struct {
	userClient     *client.UserClient
	contentClient  *client.ContentClient
	knowledgeClient *client.KnowledgeClient
	publishClient  *client.PublishClient
	accountClient  *client.AccountClient
	scheduleClient *client.ScheduleClient
	monitorClient  *client.MonitorClient
	systemClient   *client.SystemClient

	// DAL层
	fpRepo              *dal.BrowserFingerprintRepository
	proxyRepo           *dal.ProxyIPRepository
	tplRepo             *dal.ContentTemplateRepository
	staggerStrategyRepo *dal.StaggerStrategyRepository
	staggerConfigRepo   *dal.StaggerConfigRepository
	brandRepo           *dal.BrandRepository
}

// NewHandler 创建HTTP处理器
func NewHandler(
	userClient *client.UserClient,
	contentClient *client.ContentClient,
	knowledgeClient *client.KnowledgeClient,
	publishClient *client.PublishClient,
	accountClient *client.AccountClient,
	scheduleClient *client.ScheduleClient,
	monitorClient *client.MonitorClient,
	systemClient *client.SystemClient,
	fpRepo *dal.BrowserFingerprintRepository,
	proxyRepo *dal.ProxyIPRepository,
	tplRepo *dal.ContentTemplateRepository,
	staggerStrategyRepo *dal.StaggerStrategyRepository,
	staggerConfigRepo *dal.StaggerConfigRepository,
	brandRepo *dal.BrandRepository,
) *Handler {
	return &Handler{
		userClient:          userClient,
		contentClient:       contentClient,
		knowledgeClient:     knowledgeClient,
		publishClient:       publishClient,
		accountClient:       accountClient,
		scheduleClient:      scheduleClient,
		monitorClient:       monitorClient,
		systemClient:        systemClient,
		fpRepo:              fpRepo,
		proxyRepo:           proxyRepo,
		tplRepo:             tplRepo,
		staggerStrategyRepo: staggerStrategyRepo,
		staggerConfigRepo:   staggerConfigRepo,
		brandRepo:           brandRepo,
	}
}

// Health 健康检查
func (h *Handler) Health(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, success(utils.H{
		"status": "healthy",
	}))
}

// Ready 就绪检查
func (h *Handler) Ready(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, success(utils.H{
		"status": "ready",
	}))
}
