package models

import (
	"github.com/xuchengvcc/restart-life-api/internal/constants"
)

// APIError API错误结构
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, details ...string) *APIResponse {
	apiError := &APIError{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		apiError.Details = details[0]
	}

	return &APIResponse{
		Success: false,
		Error:   apiError,
	}
}

// 错误代码引用 - 使用constants包中的定义
var (
	ErrCodeInvalidCredentials = constants.ErrCodeInvalidCredentials
	ErrCodeTokenExpired       = constants.ErrCodeTokenExpired
	ErrCodeTokenInvalid       = constants.ErrCodeTokenInvalid
	ErrCodeUserNotFound       = constants.ErrCodeUserNotFound
	ErrCodeUserAlreadyExists  = constants.ErrCodeUserAlreadyExists
	ErrCodePermissionDenied   = constants.ErrCodePermissionDenied
	ErrCodeValidationFailed   = constants.ErrCodeValidationFailed
	ErrCodeInternalError      = constants.ErrCodeInternalError

	// 邮件相关错误代码
	ErrCodeEmailSendFailed     = constants.ErrCodeEmailSendFailed
	ErrCodeEmailConfigInvalid  = constants.ErrCodeEmailConfigInvalid
	ErrCodeEmailTemplateFailed = constants.ErrCodeEmailTemplateFailed
	ErrCodeEmailConnectFailed  = constants.ErrCodeEmailConnectFailed
	ErrCodeEmailAuthFailed     = constants.ErrCodeEmailAuthFailed

	// 验证码相关错误代码
	ErrCodeVerificationCodeExpired = constants.ErrCodeVerificationCodeExpired
	ErrCodeVerificationCodeInvalid = constants.ErrCodeVerificationCodeInvalid
	ErrCodeVerificationCodeUsed    = constants.ErrCodeVerificationCodeUsed
	ErrCodeTooManyRequests         = constants.ErrCodeTooManyRequests
	ErrCodeEmailAddressInvalid     = constants.ErrCodeEmailAddressInvalid
)

// 错误消息映射
var ErrorMessages = map[int]string{
	ErrCodeInvalidCredentials: "用户名或密码错误",
	ErrCodeTokenExpired:       "Token已过期",
	ErrCodeTokenInvalid:       "Token无效",
	ErrCodeUserNotFound:       "用户不存在",
	ErrCodeUserAlreadyExists:  "用户已存在",
	ErrCodePermissionDenied:   "权限不足",
	ErrCodeValidationFailed:   "数据验证失败",
	ErrCodeInternalError:      "内部服务器错误",

	// 邮件相关错误
	ErrCodeEmailSendFailed:     "邮件发送失败",
	ErrCodeEmailConfigInvalid:  "邮件配置无效",
	ErrCodeEmailTemplateFailed: "邮件模板处理失败",
	ErrCodeEmailConnectFailed:  "邮件服务器连接失败",
	ErrCodeEmailAuthFailed:     "邮件服务器认证失败",

	// 验证码相关错误
	ErrCodeVerificationCodeExpired: "验证码已过期",
	ErrCodeVerificationCodeInvalid: "验证码无效",
	ErrCodeVerificationCodeUsed:    "验证码已使用",
	ErrCodeTooManyRequests:         "请求过于频繁",
	ErrCodeEmailAddressInvalid:     "邮箱地址无效",
}

// GetErrorMessage 获取错误消息
func GetErrorMessage(code int) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "未知错误"
}
