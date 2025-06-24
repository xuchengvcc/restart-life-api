package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/repository"
	"github.com/xuchengvcc/restart-life-api/internal/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
	GetProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint, req *models.UpdateProfileRequest) (*models.User, error)
	ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error
	ValidateToken(ctx context.Context, token string) (*utils.Claims, error)
}

// authService 认证服务实现
type authService struct {
	userRepo        repository.UserRepository
	jwtManager      *utils.JWTManager
	passwordManager *utils.PasswordManager
	logger          *logrus.Logger
}

// NewAuthService 创建认证服务
func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager *utils.JWTManager,
	passwordManager *utils.PasswordManager,
	logger *logrus.Logger,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		jwtManager:      jwtManager,
		passwordManager: passwordManager,
		logger:          logger,
	}
}

// Register 用户注册
func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.WithError(err).Error("Failed to check username exists")
		return nil, fmt.Errorf(constants.MsgInternalError)
	}
	if exists {
		return nil, fmt.Errorf(constants.MsgUserAlreadyExists)
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.WithError(err).Error("Failed to check email exists")
		return nil, fmt.Errorf(constants.MsgInternalError)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// 生成密码哈希
	passwordHash, err := s.passwordManager.HashPassword(req.Password)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return nil, fmt.Errorf("failed to process password")
	}

	// 创建用户
	user := &models.User{
		Username:     strings.TrimSpace(req.Username),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user")
	}

	// 生成Token
	accessToken, refreshToken, expiresAt, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate tokens")
		return nil, fmt.Errorf("failed to generate tokens")
	}

	// 更新最后登录时间
	user.SetLastLogin()
	err = s.userRepo.UpdateLastLogin(ctx, user.UserID)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to update last login time")
	}

	s.logger.WithFields(logrus.Fields{
		constants.LogFieldUserID:   user.UserID,
		constants.LogFieldUsername: user.Username,
	}).Info("User registered successfully")

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// 支持用户名或邮箱登录
	var user *models.User
	var err error

	loginField := strings.TrimSpace(req.Username)
	if strings.Contains(loginField, "@") {
		// 按邮箱查找
		user, err = s.userRepo.GetByEmail(ctx, strings.ToLower(loginField))
	} else {
		// 按用户名查找
		user, err = s.userRepo.GetByUsername(ctx, loginField)
	}

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"login_field":           loginField,
			constants.LogFieldError: err.Error(),
		}).Warn("User login failed - user not found")
		return nil, fmt.Errorf(constants.MsgInvalidCredentials)
	}

	// 检查用户是否激活
	if !user.IsActive {
		s.logger.WithFields(logrus.Fields{
			constants.LogFieldUserID:   user.UserID,
			constants.LogFieldUsername: user.Username,
		}).Warn("User login failed - account disabled")
		return nil, fmt.Errorf(constants.MsgAccountDisabled)
	}

	// 验证密码
	err = s.passwordManager.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			constants.LogFieldUserID:   user.UserID,
			constants.LogFieldUsername: user.Username,
		}).Warn("User login failed - invalid password")
		return nil, fmt.Errorf(constants.MsgInvalidCredentials)
	}

	// 生成Token
	accessToken, refreshToken, expiresAt, err := s.jwtManager.GenerateTokenPair(user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate tokens")
		return nil, fmt.Errorf("failed to generate tokens")
	}

	// 更新最后登录时间
	user.SetLastLogin()
	err = s.userRepo.UpdateLastLogin(ctx, user.UserID)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to update last login time")
	}

	s.logger.WithFields(logrus.Fields{
		constants.LogFieldUserID:   user.UserID,
		constants.LogFieldUsername: user.Username,
	}).Info("User logged in successfully")

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken 刷新Token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	accessToken, newRefreshToken, expiresAt, err := s.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to refresh token")
		return nil, fmt.Errorf("invalid refresh token")
	}

	// 验证Token并获取用户信息
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		s.logger.WithError(err).Error("Failed to validate new access token")
		return nil, fmt.Errorf("failed to generate new token")
	}

	// 获取完整的用户信息
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user for refresh token")
		return nil, fmt.Errorf("user not found")
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// GetProfile 获取用户档案
func (s *authService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user profile")
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// UpdateProfile 更新用户档案
func (s *authService) UpdateProfile(ctx context.Context, userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user for profile update")
		return nil, fmt.Errorf("user not found")
	}

	// 更新用户信息
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.BirthDate != nil {
		user.BirthDate = req.BirthDate
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.Country != nil {
		user.Country = req.Country
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update user profile")
		return nil, fmt.Errorf("failed to update profile")
	}

	s.logger.WithFields(logrus.Fields{
		constants.LogFieldUserID:   user.UserID,
		constants.LogFieldUsername: user.Username,
	}).Info("User profile updated successfully")

	return user, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user for password change")
		return fmt.Errorf("user not found")
	}

	// 验证旧密码
	err = s.passwordManager.VerifyPassword(user.PasswordHash, req.OldPassword)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			constants.LogFieldUserID:   user.UserID,
			constants.LogFieldUsername: user.Username,
		}).Warn("Password change failed - invalid old password")
		return fmt.Errorf(constants.MsgPasswordIncorrect)
	}

	// 生成新密码哈希
	newPasswordHash, err := s.passwordManager.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash new password")
		return fmt.Errorf("failed to process new password")
	}

	// 更新密码
	user.PasswordHash = newPasswordHash
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update password")
		return fmt.Errorf("failed to update password")
	}

	s.logger.WithFields(logrus.Fields{
		constants.LogFieldUserID:   user.UserID,
		constants.LogFieldUsername: user.Username,
	}).Info("User password changed successfully")

	return nil
}

// ValidateToken 验证Token
func (s *authService) ValidateToken(ctx context.Context, token string) (*utils.Claims, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 检查Token类型
	if claims.Type != constants.TokenTypeAccess {
		return nil, fmt.Errorf("token is not an access token")
	}

	// 验证用户是否仍然存在且激活
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is disabled")
	}

	return claims, nil
}
