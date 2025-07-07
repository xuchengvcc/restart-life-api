package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/api/middleware"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService             services.AuthService
	verificationCodeService services.VerificationCodeService
	logger                  *logrus.Logger
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(
	authService services.AuthService,
	verificationCodeService services.VerificationCodeService,
	logger *logrus.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService:             authService,
		verificationCodeService: verificationCodeService,
		logger:                  logger,
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
		// 使用 errors.Is 判断错误类型
		if errors.Is(err, constants.ErrUserAlreadyExists) || errors.Is(err, constants.ErrEmailAlreadyExists) {
			response := models.NewErrorResponse(models.ErrCodeUserAlreadyExists, err.Error())
			c.JSON(http.StatusConflict, response)
		} else if errors.Is(err, constants.ErrPasswordProcessFailed) {
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
		if errors.Is(err, constants.ErrInvalidCredentials) {
			response := models.NewErrorResponse(models.ErrCodeInvalidCredentials, "用户名或密码错误")
			c.JSON(http.StatusUnauthorized, response)
		} else if errors.Is(err, constants.ErrAccountDisabled) {
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
		if errors.Is(err, constants.ErrUserNotFound) {
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
		if errors.Is(err, constants.ErrUserNotFound) {
			response := models.NewErrorResponse(models.ErrCodeUserNotFound, "用户不存在")
			c.JSON(http.StatusNotFound, response)
		} else if errors.Is(err, constants.ErrPasswordIncorrect) {
			response := models.NewErrorResponse(models.ErrCodeInvalidCredentials, "原密码错误")
			c.JSON(http.StatusBadRequest, response)
		} else if errors.Is(err, constants.ErrPasswordProcessFailed) {
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

// SendVerificationCode 发送验证码
// @Summary 发送验证码
// @Description 向指定邮箱发送验证码
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.SendVerificationCodeRequest true "发送验证码请求"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 429 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/send-verification-code [post]
func (h *AuthHandler) SendVerificationCode(c *gin.Context) {
	var req models.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 发送重置密码验证码
	err := h.verificationCodeService.SendVerificationCode(
		c.Request.Context(),
		strings.TrimSpace(strings.ToLower(req.Email)),
		models.VerificationCodeTypeResetPassword,
	)

	if err != nil {
		if errors.Is(err, constants.ErrTooManyRequests) {
			response := models.NewErrorResponse(models.ErrCodeTooManyRequests, "发送过于频繁，请稍后再试")
			c.JSON(http.StatusTooManyRequests, response)
		} else if errors.Is(err, constants.ErrEmailAddressInvalid) {
			response := models.NewErrorResponse(models.ErrCodeEmailAddressInvalid, "邮箱地址无效")
			c.JSON(http.StatusBadRequest, response)
		} else if errors.Is(err, constants.ErrEmailSendFailed) {
			response := models.NewErrorResponse(models.ErrCodeEmailSendFailed, "邮件发送失败")
			c.JSON(http.StatusInternalServerError, response)
		} else {
			h.logger.WithError(err).Error("Failed to send verification code")
			response := models.NewErrorResponse(models.ErrCodeInternalError, "发送验证码失败")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := models.NewSuccessResponse(gin.H{"message": "验证码已发送到您的邮箱"})
	c.JSON(http.StatusOK, response)
}

// VerifyCode 验证验证码并获取重置令牌
// @Summary 验证验证码并获取重置令牌
// @Description 验证邮箱验证码是否正确，验证成功后返回用于重置密码的临时令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.VerifyCodeRequest true "验证验证码请求"
// @Success 200 {object} models.APIResponse{data=models.VerifyCodeResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 410 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/verify-code [post]
func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 验证验证码并创建重置令牌
	resetToken, err := h.verificationCodeService.VerifyCodeAndCreateResetToken(
		c.Request.Context(),
		strings.TrimSpace(strings.ToLower(req.Email)),
		req.Code,
	)

	if err != nil {
		if errors.Is(err, constants.ErrVerificationCodeInvalid) {
			response := models.NewErrorResponse(models.ErrCodeVerificationCodeInvalid, "验证码无效")
			c.JSON(http.StatusBadRequest, response)
		} else if errors.Is(err, constants.ErrVerificationCodeExpired) {
			response := models.NewErrorResponse(models.ErrCodeVerificationCodeExpired, "验证码已过期")
			c.JSON(http.StatusGone, response)
		} else if errors.Is(err, constants.ErrVerificationCodeUsed) {
			response := models.NewErrorResponse(models.ErrCodeVerificationCodeUsed, "验证码已使用")
			c.JSON(http.StatusBadRequest, response)
		} else if errors.Is(err, constants.ErrEmailAddressInvalid) {
			response := models.NewErrorResponse(models.ErrCodeEmailAddressInvalid, "邮箱地址无效")
			c.JSON(http.StatusBadRequest, response)
		} else if errors.Is(err, constants.ErrUserNotFound) {
			response := models.NewErrorResponse(models.ErrCodeUserNotFound, "用户不存在")
			c.JSON(http.StatusNotFound, response)
		} else {
			h.logger.WithError(err).Error("Failed to verify verification code")
			response := models.NewErrorResponse(models.ErrCodeInternalError, "验证码验证失败")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	// 返回重置令牌
	responseData := models.VerifyCodeResponse{
		Message:    "验证码验证成功",
		ResetToken: resetToken,
	}
	response := models.NewSuccessResponse(responseData)
	c.JSON(http.StatusOK, response)
}

// ResetPassword 使用令牌重置密码
// @Summary 使用令牌重置密码
// @Description 使用验证验证码后获得的令牌重置用户密码
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordWithTokenRequest true "重置密码请求"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordWithTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 使用令牌重置密码
	err := h.authService.ResetPasswordWithToken(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, constants.ErrTokenInvalid) {
			response := models.NewErrorResponse(models.ErrCodeTokenInvalid, "重置令牌无效或已过期")
			c.JSON(http.StatusUnauthorized, response)
		} else if errors.Is(err, constants.ErrUserNotFound) {
			response := models.NewErrorResponse(models.ErrCodeUserNotFound, "用户不存在")
			c.JSON(http.StatusNotFound, response)
		} else if errors.Is(err, constants.ErrPasswordProcessFailed) {
			response := models.NewErrorResponse(models.ErrCodeValidationFailed, err.Error())
			c.JSON(http.StatusBadRequest, response)
		} else {
			h.logger.WithError(err).Error("Failed to reset password with token")
			response := models.NewErrorResponse(models.ErrCodeInternalError, "重置密码失败")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := models.NewSuccessResponse(gin.H{"message": "密码重置成功"})
	c.JSON(http.StatusOK, response)
}
