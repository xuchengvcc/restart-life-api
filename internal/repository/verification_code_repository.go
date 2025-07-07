package repository

import (
	"context"
	"time"

	"github.com/xuchengvcc/restart-life-api/internal/dao"
)

// VerificationCodeRepository 验证码仓库接口
type VerificationCodeRepository interface {
	CreateVerificationCode(ctx context.Context, email, code, codeType string, ttl time.Duration) error
	GetVerificationCode(ctx context.Context, email, codeType string) (string, error)
	VerifyAndDeleteCode(ctx context.Context, email, inputCode, codeType string) (bool, error)
	DeleteVerificationCode(ctx context.Context, email, codeType string) error
	CountRecentVerificationCodes(ctx context.Context, email string, since time.Time) (int, error)
	IncrementRateLimit(ctx context.Context, email string, window time.Duration) error
	CreatePasswordResetToken(ctx context.Context, email, token string, ttl time.Duration) error
	GetEmailByResetToken(ctx context.Context, token string) (string, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
}

// verificationCodeRepository 验证码仓库实现
type verificationCodeRepository struct {
	dao *dao.VerificationCodeDAO
}

// NewVerificationCodeRepository 创建验证码仓库
func NewVerificationCodeRepository(dao *dao.VerificationCodeDAO) VerificationCodeRepository {
	return &verificationCodeRepository{
		dao: dao,
	}
}

// CreateVerificationCode 创建验证码
func (r *verificationCodeRepository) CreateVerificationCode(ctx context.Context, email, code, codeType string, ttl time.Duration) error {
	return r.dao.CreateVerificationCode(ctx, email, code, codeType, ttl)
}

// GetVerificationCode 根据邮箱和类型获取验证码
func (r *verificationCodeRepository) GetVerificationCode(ctx context.Context, email, codeType string) (string, error) {
	return r.dao.GetVerificationCode(ctx, email, codeType)
}

// VerifyAndDeleteCode 验证并删除验证码
func (r *verificationCodeRepository) VerifyAndDeleteCode(ctx context.Context, email, inputCode, codeType string) (bool, error) {
	return r.dao.VerifyAndDeleteCode(ctx, email, inputCode, codeType)
}

// DeleteVerificationCode 删除验证码
func (r *verificationCodeRepository) DeleteVerificationCode(ctx context.Context, email, codeType string) error {
	return r.dao.DeleteVerificationCode(ctx, email, codeType)
}

// CountRecentVerificationCodes 统计指定邮箱在指定时间内的验证码数量
func (r *verificationCodeRepository) CountRecentVerificationCodes(ctx context.Context, email string, since time.Time) (int, error) {
	return r.dao.CountRecentVerificationCodes(ctx, email, since)
}

// IncrementRateLimit 增加发送计数（用于频率限制）
func (r *verificationCodeRepository) IncrementRateLimit(ctx context.Context, email string, window time.Duration) error {
	return r.dao.IncrementRateLimit(ctx, email, window)
}

// CreatePasswordResetToken 创建密码重置令牌
func (r *verificationCodeRepository) CreatePasswordResetToken(ctx context.Context, email, token string, ttl time.Duration) error {
	return r.dao.CreatePasswordResetToken(ctx, email, token, ttl)
}

// GetEmailByResetToken 根据重置令牌获取邮箱
func (r *verificationCodeRepository) GetEmailByResetToken(ctx context.Context, token string) (string, error) {
	return r.dao.GetEmailByResetToken(ctx, token)
}

// DeletePasswordResetToken 删除密码重置令牌
func (r *verificationCodeRepository) DeletePasswordResetToken(ctx context.Context, token string) error {
	return r.dao.DeletePasswordResetToken(ctx, token)
}
