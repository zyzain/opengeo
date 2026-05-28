package errcode

import "fmt"

// ErrorCode 错误码
type ErrorCode int32

// 全局错误码定义
const (
	// 通用错误码
	Success         ErrorCode = 0
	InternalError   ErrorCode = 10001
	InvalidParams   ErrorCode = 10002
	Unauthorized    ErrorCode = 10003
	Forbidden       ErrorCode = 10004
	NotFound        ErrorCode = 10005
	AlreadyExists   ErrorCode = 10006
	RateLimitExceeded ErrorCode = 10007

	// 用户相关错误码 (2xxxx)
	UserNotFound      ErrorCode = 20001
	UserAlreadyExists ErrorCode = 20002
	InvalidPassword   ErrorCode = 20003
	UserDisabled      ErrorCode = 20004

	// 内容相关错误码 (3xxxx)
	ContentNotFound     ErrorCode = 30001
	ContentAlreadyExists ErrorCode = 30002
	ContentInvalid      ErrorCode = 30003
	ContentOptimizationFailed ErrorCode = 30004

	// 账号相关错误码 (4xxxx)
	AccountNotFound     ErrorCode = 40001
	AccountAlreadyExists ErrorCode = 40002
	AccountDisabled     ErrorCode = 40003
	AccountHealthCheckFailed ErrorCode = 40004

	// 发布相关错误码 (5xxxx)
	PublishTaskNotFound ErrorCode = 50001
	PublishTaskFailed   ErrorCode = 50002
	ChannelNotFound     ErrorCode = 50003
	ChannelDisabled     ErrorCode = 50004

	// 调度相关错误码 (6xxxx)
	ScheduleNotFound ErrorCode = 60001
	ScheduleDisabled ErrorCode = 60002

	// 监测相关错误码 (7xxxx)
	CitationNotFound ErrorCode = 70001
	ScoreNotFound    ErrorCode = 70002

	// 系统相关错误码 (8xxxx)
	ConfigNotFound ErrorCode = 80001
	PluginNotFound ErrorCode = 80002
	WebhookNotFound ErrorCode = 80003
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New 创建应用错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewWithDetail 创建带详情的应用错误
func NewWithDetail(code ErrorCode, message, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// Wrap 包装错误
func Wrap(code ErrorCode, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: err.Error(),
	}
}

// 预定义错误
var (
	ErrInternal     = New(InternalError, "内部错误")
	ErrInvalidParams = New(InvalidParams, "参数错误")
	ErrUnauthorized = New(Unauthorized, "未授权")
	ErrForbidden    = New(Forbidden, "禁止访问")
	ErrNotFound     = New(NotFound, "资源不存在")
	ErrAlreadyExists = New(AlreadyExists, "资源已存在")
	ErrRateLimitExceeded = New(RateLimitExceeded, "请求频率超限")

	ErrUserNotFound      = New(UserNotFound, "用户不存在")
	ErrUserAlreadyExists = New(UserAlreadyExists, "用户已存在")
	ErrInvalidPassword   = New(InvalidPassword, "密码错误")
	ErrUserDisabled      = New(UserDisabled, "用户已禁用")

	ErrContentNotFound     = New(ContentNotFound, "内容不存在")
	ErrContentAlreadyExists = New(ContentAlreadyExists, "内容已存在")
	ErrContentInvalid      = New(ContentInvalid, "内容无效")
	ErrContentOptimizationFailed = New(ContentOptimizationFailed, "内容优化失败")

	ErrAccountNotFound     = New(AccountNotFound, "账号不存在")
	ErrAccountAlreadyExists = New(AccountAlreadyExists, "账号已存在")
	ErrAccountDisabled     = New(AccountDisabled, "账号已禁用")
	ErrAccountHealthCheckFailed = New(AccountHealthCheckFailed, "账号健康检查失败")

	ErrPublishTaskNotFound = New(PublishTaskNotFound, "发布任务不存在")
	ErrPublishTaskFailed   = New(PublishTaskFailed, "发布任务失败")
	ErrChannelNotFound     = New(ChannelNotFound, "渠道不存在")
	ErrChannelDisabled     = New(ChannelDisabled, "渠道已禁用")

	ErrScheduleNotFound = New(ScheduleNotFound, "调度不存在")
	ErrScheduleDisabled = New(ScheduleDisabled, "调度已禁用")

	ErrCitationNotFound = New(CitationNotFound, "AI引用不存在")
	ErrScoreNotFound    = New(ScoreNotFound, "评分不存在")

	ErrConfigNotFound  = New(ConfigNotFound, "配置不存在")
	ErrPluginNotFound  = New(PluginNotFound, "插件不存在")
	ErrWebhookNotFound = New(WebhookNotFound, "Webhook不存在")
)