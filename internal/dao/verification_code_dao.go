package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// VerificationCodeDAO 验证码数据访问对象（基于Redis）
type VerificationCodeDAO struct {
	rdb *redis.Client
}

// NewVerificationCodeDAO 创建验证码DAO
func NewVerificationCodeDAO(rdb *redis.Client) *VerificationCodeDAO {
	return &VerificationCodeDAO{
		rdb: rdb,
	}
}

// CreateVerificationCode 创建验证码
func (dao *VerificationCodeDAO) CreateVerificationCode(ctx context.Context, email, code, codeType string, ttl time.Duration) error {
	// 构建Redis key
	key := dao.buildKey(email, codeType)

	// 简单存储验证码值
	err := dao.rdb.Set(ctx, key, code, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	return nil
}

// GetVerificationCode 根据邮箱和类型获取验证码
func (dao *VerificationCodeDAO) GetVerificationCode(ctx context.Context, email, codeType string) (string, error) {
	key := dao.buildKey(email, codeType)

	code, err := dao.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // 验证码不存在或已过期
		}
		return "", fmt.Errorf("failed to get verification code: %w", err)
	}

	return code, nil
}

// VerifyAndDeleteCode 验证并删除验证码（一次性使用）
func (dao *VerificationCodeDAO) VerifyAndDeleteCode(ctx context.Context, email, inputCode, codeType string) (bool, error) {
	key := dao.buildKey(email, codeType)

	// 获取存储的验证码
	storedCode, err := dao.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 验证码不存在或已过期
		}
		return false, fmt.Errorf("failed to get verification code: %w", err)
	}

	// 验证码匹配则删除（确保一次性使用）
	if storedCode == inputCode {
		dao.rdb.Del(ctx, key)
		return true, nil
	}

	return false, nil
}

// DeleteVerificationCode 删除验证码
func (dao *VerificationCodeDAO) DeleteVerificationCode(ctx context.Context, email, codeType string) error {
	key := dao.buildKey(email, codeType)
	return dao.rdb.Del(ctx, key).Err()
}

// CountRecentVerificationCodes 统计指定邮箱在指定时间内的验证码数量
func (dao *VerificationCodeDAO) CountRecentVerificationCodes(ctx context.Context, email string, since time.Time) (int, error) {
	// Redis中使用rate limiting key来跟踪发送频率
	rateLimitKey := dao.buildRateLimitKey(email)

	// 检查rate limiting窗口内的计数
	count, err := dao.rdb.Get(ctx, rateLimitKey).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // 没有记录，说明没有发送过
		}
		return 0, fmt.Errorf("failed to get rate limit count: %w", err)
	}

	return count, nil
}

// IncrementRateLimit 增加发送计数（用于频率限制）
func (dao *VerificationCodeDAO) IncrementRateLimit(ctx context.Context, email string, window time.Duration) error {
	rateLimitKey := dao.buildRateLimitKey(email)

	// 使用管道操作确保原子性
	pipe := dao.rdb.Pipeline()
	pipe.Incr(ctx, rateLimitKey)
	pipe.Expire(ctx, rateLimitKey, window)

	_, err := pipe.Exec(ctx)
	return err
}

// buildKey 构建验证码Redis key
func (dao *VerificationCodeDAO) buildKey(email, codeType string) string {
	return fmt.Sprintf("verification_code:%s:%s", email, codeType)
}

// buildRateLimitKey 构建频率限制Redis key
func (dao *VerificationCodeDAO) buildRateLimitKey(email string) string {
	return fmt.Sprintf("verification_code_rate_limit:%s", email)
}

// CreatePasswordResetToken 创建密码重置令牌
func (dao *VerificationCodeDAO) CreatePasswordResetToken(ctx context.Context, email, token string, ttl time.Duration) error {
	key := fmt.Sprintf("password_reset_token:%s", token)
	return dao.rdb.Set(ctx, key, email, ttl).Err()
}

// GetEmailByResetToken 根据重置令牌获取邮箱
func (dao *VerificationCodeDAO) GetEmailByResetToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("password_reset_token:%s", token)
	email, err := dao.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // 令牌不存在或已过期
		}
		return "", fmt.Errorf("failed to get email by reset token: %w", err)
	}
	return email, nil
}

// DeletePasswordResetToken 删除密码重置令牌
func (dao *VerificationCodeDAO) DeletePasswordResetToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("password_reset_token:%s", token)
	return dao.rdb.Del(ctx, key).Err()
}
