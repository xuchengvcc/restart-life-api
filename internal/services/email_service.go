package services

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	emailservicesimpl "github.com/xuchengvcc/restart-life-api/internal/services/email_services_impl"
)

type EmailService interface {
	// SendVeriCode 发送验证码邮件
	SendVeriCode(to string, code int32) error

	// ValidateEmailConfig 验证邮件配置是否正确
	ValidateEmailConfig() error
}

// NewEmailService 创建邮件服务实例
func NewEmailService(logger *logrus.Logger, config config.EmailConfig) EmailService {
	return emailservicesimpl.NewNeteaseEmailService(logger, config)
}

// ValidateEmailAddress 邮箱地址验证
func ValidateEmailAddress(email string) error {
	if email == "" {
		return fmt.Errorf("邮箱地址不能为空")
	}

	// 使用正则表达式验证邮箱格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("邮箱地址格式无效")
	}

	return nil
}
