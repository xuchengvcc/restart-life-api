package services

import (
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	emailservicesimpl "github.com/xuchengvcc/restart-life-api/internal/services/email_services_impl"
)

type EmailService interface {
	// SendVeriCode 发送验证码邮件
	SendVeriCode(to string, code int32) error

	// SendNotificationEmail 发送通知邮件（HTML格式）
	SendNotificationEmail(to string, subject string, htmlBody string) error

	// ValidateEmailConfig 验证邮件配置是否正确
	ValidateEmailConfig() error
}

// NewEmailService 创建邮件服务实例
func NewEmailService(logger *logrus.Logger, config config.EmailConfig) EmailService {
	return emailservicesimpl.NewNeteaseEmailService(logger, config)
}
