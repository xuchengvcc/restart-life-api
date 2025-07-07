package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/repository"
)

// VerificationCodeService 验证码服务接口
type VerificationCodeService interface {
	// SendVerificationCode 发送验证码
	SendVerificationCode(ctx context.Context, email, codeType string) error
	// VerifyCode 验证验证码
	VerifyCode(ctx context.Context, email, code, codeType string) error
	// VerifyCodeAndCreateResetToken 验证验证码并创建重置令牌（两步密码重置第一步）
	VerifyCodeAndCreateResetToken(ctx context.Context, email, code string) (string, error)
	// GetEmailByResetToken 根据重置令牌获取邮箱（两步密码重置第二步验证）
	GetEmailByResetToken(ctx context.Context, token string) (string, error)
	// DeletePasswordResetToken 删除密码重置令牌
	DeletePasswordResetToken(ctx context.Context, token string) error
}

// verificationCodeService 验证码服务实现
type verificationCodeService struct {
	codeRepo     repository.VerificationCodeRepository
	emailService EmailService
	logger       *logrus.Logger
}

// NewVerificationCodeService 创建验证码服务
func NewVerificationCodeService(
	codeRepo repository.VerificationCodeRepository,
	emailService EmailService,
	logger *logrus.Logger,
) VerificationCodeService {
	return &verificationCodeService{
		codeRepo:     codeRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// SendVerificationCode 发送验证码
func (s *verificationCodeService) SendVerificationCode(ctx context.Context, email, codeType string) error {
	// 验证邮箱地址
	if err := ValidateEmailAddress(email); err != nil {
		return constants.ErrEmailAddressInvalid
	}

	// 检查频率限制 - 1分钟内最多发送1次
	recentCount, err := s.codeRepo.CountRecentVerificationCodes(ctx, email, time.Now().Add(-time.Minute))
	if err != nil {
		s.logger.WithError(err).Error("Failed to count recent verification codes")
		return constants.ErrInternalError
	}

	if recentCount > 0 {
		return constants.ErrTooManyRequests
	}

	// 生成6位数字验证码
	code, err := s.generateVerificationCode()
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate verification code")
		return constants.ErrInternalError
	}

	// 存储验证码到Redis，10分钟有效期
	if err := s.codeRepo.CreateVerificationCode(ctx, email, code, codeType, constants.VerificationCodeExpireMinute*time.Minute); err != nil {
		s.logger.WithError(err).Error("Failed to create verification code")
		return constants.ErrInternalError
	}

	// 增加频率限制计数
	if err := s.codeRepo.IncrementRateLimit(ctx, email, time.Minute); err != nil {
		s.logger.WithError(err).Error("Failed to increment rate limit")
		// 这里不返回错误，因为验证码已经创建成功
	}

	// 发送邮件
	codeInt, _ := big.NewInt(0).SetString(code, 10)
	if err := s.emailService.SendVeriCode(email, int32(codeInt.Int64())); err != nil {
		s.logger.WithError(err).WithField("email", email).Error("Failed to send verification code email")
		return constants.ErrEmailSendFailed
	}

	s.logger.WithFields(logrus.Fields{
		"email": email,
		"type":  codeType,
	}).Info("Verification code sent successfully")

	return nil
}

// VerifyCode 验证验证码
func (s *verificationCodeService) VerifyCode(ctx context.Context, email, code, codeType string) error {
	// 验证邮箱地址
	if err := ValidateEmailAddress(email); err != nil {
		return constants.ErrEmailAddressInvalid
	}

	// 验证并删除验证码（确保一次性使用）
	isValid, err := s.codeRepo.VerifyAndDeleteCode(ctx, email, code, codeType)
	if err != nil {
		s.logger.WithError(err).Error("Failed to verify verification code")
		return constants.ErrInternalError
	}

	if !isValid {
		return constants.ErrVerificationCodeInvalid
	}

	s.logger.WithFields(logrus.Fields{
		"email": email,
		"type":  codeType,
	}).Info("Verification code verified successfully")

	return nil
}

// VerifyCodeAndCreateResetToken 验证验证码并创建重置令牌（两步密码重置第一步）
func (s *verificationCodeService) VerifyCodeAndCreateResetToken(ctx context.Context, email, code string) (string, error) {
	// 验证邮箱地址
	if err := ValidateEmailAddress(email); err != nil {
		return "", constants.ErrEmailAddressInvalid
	}

	// 验证并删除验证码
	isValid, err := s.codeRepo.VerifyAndDeleteCode(ctx, email, code, models.VerificationCodeTypeResetPassword)
	if err != nil {
		s.logger.WithError(err).Error("Failed to verify verification code")
		return "", constants.ErrInternalError
	}

	if !isValid {
		return "", constants.ErrVerificationCodeInvalid
	}

	// 生成重置令牌
	resetToken, err := s.generateResetToken()
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate reset token")
		return "", constants.ErrInternalError
	}

	// 创建重置令牌，有效期10分钟
	if err := s.codeRepo.CreatePasswordResetToken(ctx, email, resetToken, constants.PasswordResetExpireMinute*time.Minute); err != nil {
		s.logger.WithError(err).Error("Failed to create password reset token")
		return "", constants.ErrInternalError
	}

	s.logger.WithFields(logrus.Fields{
		"email": email,
	}).Info("Password reset token created successfully")

	return resetToken, nil
}

// GetEmailByResetToken 根据重置令牌获取邮箱（两步密码重置第二步验证）
func (s *verificationCodeService) GetEmailByResetToken(ctx context.Context, token string) (string, error) {
	return s.codeRepo.GetEmailByResetToken(ctx, token)
}

// DeletePasswordResetToken 删除密码重置令牌
func (s *verificationCodeService) DeletePasswordResetToken(ctx context.Context, token string) error {
	return s.codeRepo.DeletePasswordResetToken(ctx, token)
}

// generateVerificationCode 生成6位数字验证码
func (s *verificationCodeService) generateVerificationCode() (string, error) {
	// 生成100000到999999之间的随机数
	min := big.NewInt(100000)
	max := big.NewInt(999999)
	diff := new(big.Int).Sub(max, min)
	diff.Add(diff, big.NewInt(1))

	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return "", err
	}

	n.Add(n, min)
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// generateResetToken 生成重置令牌
func (s *verificationCodeService) generateResetToken() (string, error) {
	// 生成32位随机字符串
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 32)

	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		token[i] = charset[num.Int64()]
	}

	return string(token), nil
}
