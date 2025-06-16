package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RecoveryConfig 恢复中间件配置
type RecoveryConfig struct {
	EnableStackTrace bool   // 是否启用堆栈跟踪
	LogLevel         string // 日志级别
}

// DefaultRecoveryConfig 默认恢复配置
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		EnableStackTrace: true,
		LogLevel:         "error",
	}
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// RecoveryMiddleware 异常恢复中间件
func RecoveryMiddleware(config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := debug.Stack()

				// 获取请求信息
				method := c.Request.Method
				path := c.Request.URL.Path
				clientIP := c.ClientIP()
				userAgent := c.Request.UserAgent()
				requestID := c.GetString("request_id")
				userID := c.GetString("user_id")

				// 构建日志字段
				fields := logrus.Fields{
					"error":      fmt.Sprintf("%v", err),
					"method":     method,
					"path":       path,
					"client_ip":  clientIP,
					"user_agent": userAgent,
				}

				// 添加可选字段
				if requestID != "" {
					fields["request_id"] = requestID
				}
				if userID != "" {
					fields["user_id"] = userID
				}

				// 添加堆栈跟踪
				if config.EnableStackTrace {
					fields["stack"] = string(stack)
				}

				// 记录错误日志
				logrus.WithFields(fields).Error("Panic recovered")

				// 返回统一错误响应
				errorResponse := ErrorResponse{
					Success: false,
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "服务器内部错误，请稍后重试",
				}

				// 在开发环境下提供更多错误信息
				if gin.Mode() == gin.DebugMode {
					errorResponse.Details = fmt.Sprintf("%v", err)
				}

				c.JSON(http.StatusInternalServerError, errorResponse)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// CustomRecoveryMiddleware 自定义恢复中间件，支持自定义错误处理函数
func CustomRecoveryMiddleware(handler func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				handler(c, err)
			}
		}()

		c.Next()
	}
}
