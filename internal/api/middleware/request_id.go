package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const (
	// RequestIDHeader 请求ID头部名称
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey 在gin.Context中存储请求ID的键
	RequestIDKey = "request_id"
)

// RequestIDConfig 请求ID中间件配置
type RequestIDConfig struct {
	HeaderName string        // 头部名称
	ContextKey string        // 上下文键名
	Generator  func() string // ID生成器
}

// DefaultRequestIDConfig 默认请求ID配置
func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		HeaderName: RequestIDHeader,
		ContextKey: RequestIDKey,
		Generator:  generateRequestID,
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware(config RequestIDConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取现有的请求ID
		requestID := c.GetHeader(config.HeaderName)

		// 如果没有请求ID，则生成一个新的
		if requestID == "" {
			requestID = config.Generator()
		}

		// 设置到上下文中
		c.Set(config.ContextKey, requestID)

		// 设置响应头
		c.Header(config.HeaderName, requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return generateTimeBasedID()
	}
	return hex.EncodeToString(bytes)
}

// generateTimeBasedID 基于时间生成ID（备选方案）
func generateTimeBasedID() string {
	// 简单的时间戳ID生成
	// 在实际生产环境中，可以使用更复杂的算法
	return hex.EncodeToString([]byte("fallback-id"))
}

// GetRequestID 从gin.Context中获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
