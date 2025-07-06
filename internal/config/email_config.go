package config

import (
	"strconv"

	"github.com/spf13/viper"
)

type EmailConfig struct {
	SMTPConfig SMTPConfig     `mapstructure:"smtp"`
	Template   TemplateConfig `mapstructure:"template"`
}

type SMTPConfig struct {
	Server   string `mapstructure:"server"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type TemplateConfig struct {
	VeriCode string `mapstructure:"vericode"`
}

func setEmailConfigDefault() {
	viper.SetDefault("email.smtp.server", "smtp.163.com")
	viper.SetDefault("email.smtp.port", "465")
	viper.SetDefault("email.smtp.user", "A Second Chance Official")
	viper.SetDefault("email.smtp.password", "your-email-auth-password")
	viper.SetDefault("email.smtp.from", "your-email@163.com")
	viper.SetDefault("email.template.vericode", "template/vericode.html")
}

func (c *EmailConfig) GetVeriCodeTemplatePath() string {
	return c.Template.VeriCode
}

// GetSMTPAddr 获取SMTP服务器地址
func (c *EmailConfig) GetSMTPAddr() string {
	return c.SMTPConfig.Server + ":" + c.SMTPConfig.Port
}

// GetPortInt 获取端口号的整数值
func (c *EmailConfig) GetPortInt() int {
	port, err := strconv.Atoi(c.SMTPConfig.Port)
	if err != nil {
		return 465 // 默认端口
	}
	return port
}
