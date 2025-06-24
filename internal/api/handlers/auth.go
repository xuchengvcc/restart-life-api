package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/api/middleware"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService services.AuthService
	logger      *logrus.Logger
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService services.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册请求"
// @Success 201 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid register request")
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 清理输入数据
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	authResponse, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if strings.Contains(err.Error(), "already exists") {
			response := models.NewErrorResponse(models.ErrCodeUserAlreadyExists, err.Error())
			c.JSON(http.StatusConflict, response)
		} else if strings.Contains(err.Error(), "password") {
			response := models.NewErrorResponse(models.ErrCodeValidationFailed, err.Error())
			c.JSON(http.StatusBadRequest, response)
		} else {
			h.logger.WithError(err).Error("Registration failed")
			response := models.NewErrorResponse(models.ErrCodeInternalError, "Registration failed")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := models.NewSuccessResponse(authResponse)
	c.JSON(http.StatusCreated, response)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，支持用户名或邮箱登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录请求"
// @Success 200 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid login request")
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	authResponse, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			response := models.NewErrorResponse(models.ErrCodeInvalidCredentials, "用户名或密码错误")
			c.JSON(http.StatusUnauthorized, response)
		} else if strings.Contains(err.Error(), "disabled") {
			response := models.NewErrorResponse(models.ErrCodePermissionDenied, "账户已被禁用")
			c.JSON(http.StatusUnauthorized, response)
		} else {
			h.logger.WithError(err).Error("Login failed")
			response := models.NewErrorResponse(models.ErrCodeInternalError, "登录失败")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := models.NewSuccessResponse(authResponse)
	c.JSON(http.StatusOK, response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口（客户端清除Token）
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	username, _ := middleware.GetUsername(c)

	h.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"username": username,
	}).Info("User logged out")

	response := models.NewSuccessResponse(gin.H{"message": "退出登录成功"})
	c.JSON(http.StatusOK, response)
}

// RefreshToken 刷新Token
// @Summary 刷新访问Token
// @Description 使用刷新Token获取新的访问Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "刷新Token请求"
// @Success 200 {object} models.APIResponse{data=models.AuthResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "刷新Token不能为空", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	authResponse, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response := models.NewErrorResponse(models.ErrCodeTokenInvalid, "刷新Token无效或已过期")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	response := models.NewSuccessResponse(authResponse)
	c.JSON(http.StatusOK, response)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user profile")
		response := models.NewErrorResponse(models.ErrCodeUserNotFound, "用户不存在")
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := models.NewSuccessResponse(user)
	c.JSON(http.StatusOK, response)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的个人信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.UpdateProfileRequest true "更新用户信息请求"
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response := models.NewErrorResponse(models.ErrCodeUserNotFound, "用户不存在")
			c.JSON(http.StatusNotFound, response)
		} else {
			h.logger.WithError(err).Error("Failed to update user profile")
			response := models.NewErrorResponse(models.ErrCodeInternalError, "更新用户信息失败")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := models.NewSuccessResponse(user)
	c.JSON(http.StatusOK, response)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.MustGetUserID(c)

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response := models.NewErrorResponse(models.ErrCodeUserNotFound, "用户不存在")
			c.JSON(http.StatusNotFound, response)
		} else if strings.Contains(err.Error(), "incorrect") {
			response := models.NewErrorResponse(models.ErrCodeInvalidCredentials, "原密码错误")
			c.JSON(http.StatusBadRequest, response)
		} else if strings.Contains(err.Error(), "password") {
			response := models.NewErrorResponse(models.ErrCodeValidationFailed, err.Error())
			c.JSON(http.StatusBadRequest, response)
		} else {
			h.logger.WithError(err).Error("Failed to change password")
			response := models.NewErrorResponse(models.ErrCodeInternalError, "修改密码失败")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := models.NewSuccessResponse(gin.H{"message": "密码修改成功"})
	c.JSON(http.StatusOK, response)
}
