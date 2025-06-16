package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerConfig 日志中间件配置
type LoggerConfig struct {
	SkipPaths []string // 跳过记录的路径
}

// DefaultLoggerConfig 默认日志配置
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		SkipPaths: []string{
			"/health",
			"/ping",
			"/metrics",
		},
	}
}

// LoggerMiddleware 请求日志中间件
func LoggerMiddleware(config LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过该路径
		path := c.Request.URL.Path
		if shouldSkipPath(path, config.SkipPaths) {
			c.Next()
			return
		}

		// 记录开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 获取请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		// 获取自定义头部
		platform := c.GetHeader("X-Platform")
		version := c.GetHeader("X-Version")
		requestID := c.GetString("request_id")
		userID := c.GetString("user_id")

		// 构建日志字段
		fields := logrus.Fields{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"body_size":  bodySize,
			"user_agent": userAgent,
		}

		// 添加可选字段
		if platform != "" {
			fields["platform"] = platform
		}
		if version != "" {
			fields["version"] = version
		}
		if requestID != "" {
			fields["request_id"] = requestID
		}
		if userID != "" {
			fields["user_id"] = userID
		}

		// 根据状态码选择日志级别
		entry := logrus.WithFields(fields)

		switch {
		case statusCode >= 500:
			entry.Error("Server error")
		case statusCode >= 400:
			entry.Warn("Client error")
		case statusCode >= 300:
			entry.Info("Redirection")
		default:
			entry.Info("Request completed")
		}
	}
}

// shouldSkipPath 检查是否应该跳过该路径的日志记录
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}
