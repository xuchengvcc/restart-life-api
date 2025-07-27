package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityConfig 安全配置
type SecurityConfig struct {
	ContentTypeOptions    string // X-Content-Type-Options
	FrameOptions          string // X-Frame-Options
	XSSProtection         string // X-XSS-Protection
	ContentSecurityPolicy string // Content-Security-Policy
	CacheControl          string // Cache-Control
	ReferrerPolicy        string // Referrer-Policy
	PermissionsPolicy     string // Permissions-Policy
}

// DefaultSecurityConfig 默认安全配置
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		ContentTypeOptions:    "nosniff",
		FrameOptions:          "DENY",
		XSSProtection:         "1; mode=block",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'",
		CacheControl:          "no-cache, no-store, must-revalidate, private",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		PermissionsPolicy:     "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()",
	}
}

// SecurityMiddleware 安全头中间件
func SecurityMiddleware(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Content-Type-Options: 防止MIME类型嗅探
		if config.ContentTypeOptions != "" {
			c.Header("X-Content-Type-Options", config.ContentTypeOptions)
		}

		// X-Frame-Options: 防止点击劫持
		if config.FrameOptions != "" {
			c.Header("X-Frame-Options", config.FrameOptions)
		}

		// X-XSS-Protection: XSS保护
		if config.XSSProtection != "" {
			c.Header("X-XSS-Protection", config.XSSProtection)
		}

		// Content-Security-Policy: 内容安全策略
		if config.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", config.ContentSecurityPolicy)
		}

		// Cache-Control: 缓存控制
		if config.CacheControl != "" {
			c.Header("Cache-Control", config.CacheControl)
		}

		// Referrer-Policy: 引用策略
		if config.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", config.ReferrerPolicy)
		}

		// Permissions-Policy: 权限策略
		if config.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", config.PermissionsPolicy)
		}

		// 针对API响应的特殊缓存策略
		if isAPIRequest(c) {
			// API响应使用no-cache，但允许重新验证
			c.Header("Cache-Control", "no-cache, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// isAPIRequest 判断是否为API请求
func isAPIRequest(c *gin.Context) bool {
	path := c.Request.URL.Path
	return len(path) >= 4 && path[:4] == "/api"
}

// ProductionSecurityConfig 生产环境安全配置
func ProductionSecurityConfig() SecurityConfig {
	return SecurityConfig{
		ContentTypeOptions:    "nosniff",
		FrameOptions:          "SAMEORIGIN", // 生产环境可以使用SAMEORIGIN
		XSSProtection:         "1; mode=block",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'",
		CacheControl:          "no-cache, no-store, must-revalidate, private",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		PermissionsPolicy:     "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()",
	}
}
