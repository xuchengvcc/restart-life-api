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
}

// GetErrorMessage 获取错误消息
func GetErrorMessage(code int) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "未知错误"
}
