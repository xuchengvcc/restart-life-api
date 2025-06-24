package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
	"github.com/xuchengvcc/restart-life-api/internal/utils"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	authService services.AuthService
	logger      *logrus.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService services.AuthService, logger *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.respondUnauthorized(c, "Authorization header is required")
			return
		}

		// 提取Token
		token := utils.ExtractTokenFromHeader(authHeader)
		if token == "" {
			m.respondUnauthorized(c, "Invalid authorization header format")
			return
		}

		// 验证Token
		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.logger.WithError(err).Warn("Invalid token provided")
			m.respondUnauthorized(c, "Invalid or expired token")
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuth 可选认证的中间件
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 提取Token
		token := utils.ExtractTokenFromHeader(authHeader)
		if token == "" {
			c.Next()
			return
		}

		// 验证Token
		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.logger.WithError(err).Debug("Optional auth failed")
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id, true
		}
	}
	return 0, false
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name, true
		}
	}
	return "", false
}

// GetEmail 从上下文中获取邮箱
func GetEmail(c *gin.Context) (string, bool) {
	if email, exists := c.Get("email"); exists {
		if e, ok := email.(string); ok {
			return e, true
		}
	}
	return "", false
}

// GetClaims 从上下文中获取Claims
func GetClaims(c *gin.Context) (*utils.Claims, bool) {
	if claims, exists := c.Get("claims"); exists {
		if c, ok := claims.(*utils.Claims); ok {
			return c, true
		}
	}
	return nil, false
}

// MustGetUserID 从上下文中获取用户ID（必须存在）
func MustGetUserID(c *gin.Context) uint {
	userID, exists := GetUserID(c)
	if !exists {
		panic("user_id not found in context")
	}
	return userID
}

// respondUnauthorized 返回未授权响应
func (m *AuthMiddleware) respondUnauthorized(c *gin.Context, message string) {
	response := models.NewErrorResponse(models.ErrCodeTokenInvalid, message)
	c.JSON(http.StatusUnauthorized, response)
	c.Abort()
}

// CheckPermission 权限检查中间件
func (m *AuthMiddleware) CheckPermission(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以根据需要实现权限检查逻辑
		// 目前简单实现：只要用户已认证就有权限
		userID, exists := GetUserID(c)
		if !exists {
			m.respondUnauthorized(c, "Authentication required")
			return
		}

		// 在实际应用中，这里可以查询用户的角色和权限
		// 然后检查是否具有所需的权限
		m.logger.WithFields(logrus.Fields{
			"user_id":     userID,
			"permissions": strings.Join(requiredPermissions, ","),
		}).Debug("Permission check passed")

		c.Next()
	}
}
