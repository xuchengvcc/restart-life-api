package emailservicesimpl

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
)

type NeteaseEmailService struct {
	config *config.EmailConfig
	logger *logrus.Logger
}

func NewNeteaseEmailService(logger *logrus.Logger, config config.EmailConfig) *NeteaseEmailService {
	return &NeteaseEmailService{
		config: &config,
		logger: logger,
	}
}

func (s *NeteaseEmailService) SendVeriCode(to string, code int32) error {
	// 使用模板工具类
	template := GetVerificationCodeTemplate(to, code)

	// 记录发送验证码的日志
	s.logger.WithFields(logrus.Fields{
		"to":       to,
		"code":     "****", // 不记录完整验证码，保护隐私
		"template": template.Template,
	}).Info("starting to send verification code email")

	// 发送模板邮件
	return s.SendTemplatedEmail(to, template)
}

// sendMailSSL 使用SSL发送邮件
func (s *NeteaseEmailService) sendMailSSL(addr, from string, to []string, msg []byte) error {
	// 建立TLS连接
	tlsConfig := &tls.Config{
		ServerName:         s.config.SMTPConfig.Server,
		InsecureSkipVerify: false,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailConnectFailed, err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.config.SMTPConfig.Server)
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailCreateClientFailed, err)
	}
	defer client.Close()

	// 认证
	auth := smtp.PlainAuth("", s.config.SMTPConfig.From, s.config.SMTPConfig.Password, s.config.SMTPConfig.Server)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailAuthFailed, err)
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailSetSenderFailed, err)
	}

	// 设置收件人
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("%w: %v", constants.ErrEmailSetRecipientFailed, err)
		}
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailGetWriterFailed, err)
	}
	defer writer.Close()

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailWriteContentFailed, err)
	}

	return nil
}

// sendMailTLS 使用TLS发送邮件
func (s *NeteaseEmailService) sendMailTLS(addr, from string, to []string, msg []byte) error {
	// 建立普通连接
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailConnectFailed, err)
	}
	defer client.Close()

	// 启动TLS
	tlsConfig := &tls.Config{
		ServerName:         s.config.SMTPConfig.Server,
		InsecureSkipVerify: false,
	}

	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailTLSFailed, err)
	}

	// 认证
	auth := smtp.PlainAuth("", s.config.SMTPConfig.From, s.config.SMTPConfig.Password, s.config.SMTPConfig.Server)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailAuthFailed, err)
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailSetSenderFailed, err)
	}

	// 设置收件人
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("%w: %v", constants.ErrEmailSetRecipientFailed, err)
		}
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailGetWriterFailed, err)
	}
	defer writer.Close()

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailWriteContentFailed, err)
	}

	return nil
}

// ValidateEmailConfig 验证邮件配置是否正确
func (s *NeteaseEmailService) ValidateEmailConfig() error {
	if s.config.SMTPConfig.Server == "" {
		return constants.ErrEmailConfigInvalid
	}
	if s.config.SMTPConfig.Port == "" {
		return constants.ErrEmailConfigInvalid
	}
	if s.config.SMTPConfig.From == "" {
		return constants.ErrEmailConfigInvalid
	}
	if s.config.SMTPConfig.Password == "" {
		return constants.ErrEmailConfigInvalid
	}
	if s.config.Template.VeriCode == "" {
		return constants.ErrEmailConfigInvalid
	}

	// 检查模板文件是否存在
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%w: %v", constants.ErrEmailGetWorkdirFailed, err)
	}

	templatePath := currentDir + "/" + s.config.Template.VeriCode
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return constants.ErrEmailTemplateNotFound
	}

	s.logger.Info("email configuration validation passed")
	return nil
}

// SendTemplatedEmail 发送基于模板的邮件
func (s *NeteaseEmailService) SendTemplatedEmail(to string, template *EmailTemplate) error {
	// 记录发送邮件的日志
	s.logger.WithFields(logrus.Fields{
		"to":       to,
		"subject":  template.Subject,
		"template": template.Template,
	}).Info("starting to send templated email")

	// 处理模板
	message, err := s.ProcessTemplate(template)
	if err != nil {
		s.logger.WithError(err).Error("failed to process email template")
		return fmt.Errorf("%w: %v", constants.ErrEmailProcessTemplateFailed, err)
	}

	// 构建SMTP服务器地址
	smtpAddr := fmt.Sprintf("%s:%s", s.config.SMTPConfig.Server, s.config.SMTPConfig.Port)

	// 构建邮件内容
	msgStr := fmt.Sprintf(`To: %s
From: %s <%s>
Subject: %s
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

%s`, to, s.config.SMTPConfig.User, s.config.SMTPConfig.From, template.Subject, message)

	msg := []byte(msgStr)
	tos := []string{to}

	s.logger.WithFields(logrus.Fields{
		"smtpServer": smtpAddr,
		"from":       s.config.SMTPConfig.From,
		"to":         to,
		"subject":    template.Subject,
		"port":       s.config.SMTPConfig.Port,
	}).Debug("sending email")

	// 根据端口选择发送方式
	if s.config.SMTPConfig.Port == "465" {
		// SSL 连接
		err = s.sendMailSSL(smtpAddr, s.config.SMTPConfig.From, tos, msg)
	} else if s.config.SMTPConfig.Port == "587" {
		// TLS 连接
		err = s.sendMailTLS(smtpAddr, s.config.SMTPConfig.From, tos, msg)
	} else {
		// 普通连接
		auth := smtp.PlainAuth("", s.config.SMTPConfig.From, s.config.SMTPConfig.Password, s.config.SMTPConfig.Server)
		err = smtp.SendMail(smtpAddr, auth, s.config.SMTPConfig.From, tos, msg)
	}

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"smtpServer": smtpAddr,
			"to":         to,
		}).Error("failed to send email")
		return fmt.Errorf("%w: %v", constants.ErrEmailSendFailed, err)
	}

	s.logger.WithField("to", to).Info("email sent successfully")
	return nil
}

// ProcessTemplate 处理邮件模板
func (s *NeteaseEmailService) ProcessTemplate(template *EmailTemplate) (string, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w: %v", constants.ErrEmailGetWorkdirFailed, err)
	}

	// 构建模板文件路径
	templatePath := currentDir + "/" + template.Template

	// 读取HTML模板
	htmlContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", constants.ErrEmailReadTemplateFailed, err)
	}

	// 替换模板中的占位符
	message := string(htmlContent)
	for key, value := range template.Data {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		message = strings.Replace(message, placeholder, fmt.Sprintf("%v", value), -1)
	}

	return message, nil
}
