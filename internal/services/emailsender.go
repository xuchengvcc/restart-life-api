package services

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	emailservicesimpl "github.com/xuchengvcc/restart-life-api/internal/services/email_services_impl"
)

type EmailSender interface {
	// SendVeriCode 发送验证码邮件
	SendVeriCode(to string, code int32) error

	// ValidateEmailConfig 验证邮件配置是否正确
	ValidateEmailConfig() error
}

// NewEmailService 创建邮件服务实例
func NewEmailService(logger *logrus.Logger, config config.EmailConfig) EmailSender {
	return emailservicesimpl.NewNeteaseEmailService(logger, config)
}

// ValidateEmailAddress 简单的邮箱地址验证
func ValidateEmailAddress(email string) error {
	if email == "" {
		return fmt.Errorf("邮箱地址不能为空")
	}
	// 这里可以添加更复杂的邮箱验证逻辑
	// 比如使用正则表达式验证邮箱格式
	return nil
}
