package models

import (
	"time"
)

// User 用户模型
type User struct {
	UserID       uint    `json:"user_id" db:"user_id"`
	Username     string  `json:"username" db:"username"`
	Email        string  `json:"email" db:"email"`
	PasswordHash string  `json:"-" db:"password_hash"`
	CreatedAt    int64   `json:"created_at" db:"created_at"`
	UpdatedAt    int64   `json:"updated_at" db:"updated_at"`
	LastLogin    *int64  `json:"last_login" db:"last_login"`
	IsActive     bool    `json:"is_active" db:"is_active"`
	AvatarURL    *string `json:"avatar_url" db:"avatar_url"`
	Bio          *string `json:"bio" db:"bio"`
	BirthDate    *string `json:"birth_date" db:"birth_date"`
	Gender       *int    `json:"gender" db:"gender"`
	Country      *string `json:"country" db:"country"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Bio       *string `json:"bio"`
	BirthDate *string `json:"birth_date"`
	Gender    *int    `json:"gender"`
	Country   *string `json:"country"`
	AvatarURL *string `json:"avatar_url"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// SendVerificationCodeRequest 发送验证码请求
type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// VerifyCodeResponse 验证验证码响应
type VerifyCodeResponse struct {
	Message    string `json:"message"`
	ResetToken string `json:"reset_token"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPasswordWithTokenRequest 使用令牌重置密码请求
type ResetPasswordWithTokenRequest struct {
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// TokenClaims JWT Token声明
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Type     string `json:"type"` // access or refresh
}

// GetCreatedAtTime 获取创建时间
func (u *User) GetCreatedAtTime() time.Time {
	return time.Unix(0, u.CreatedAt*int64(time.Millisecond))
}

// GetUpdatedAtTime 获取更新时间
func (u *User) GetUpdatedAtTime() time.Time {
	return time.Unix(0, u.UpdatedAt*int64(time.Millisecond))
}

// GetLastLoginTime 获取最后登录时间
func (u *User) GetLastLoginTime() *time.Time {
	if u.LastLogin == nil {
		return nil
	}
	t := time.Unix(0, *u.LastLogin*int64(time.Millisecond))
	return &t
}

// SetTimestamps 设置时间戳
func (u *User) SetTimestamps() {
	now := time.Now().UnixMilli()
	if u.UserID == 0 {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

// SetLastLogin 设置最后登录时间
func (u *User) SetLastLogin() {
	now := time.Now().UnixMilli()
	u.LastLogin = &now
}
