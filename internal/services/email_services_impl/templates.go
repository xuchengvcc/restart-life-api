package emailservicesimpl

import (
	"strconv"
	"time"
)

// EmailTemplate 邮件模板结构
type EmailTemplate struct {
	Subject  string
	Template string
	Data     map[string]interface{}
}

// GetVerificationCodeTemplate 获取验证码邮件模板
func GetVerificationCodeTemplate(to string, code int32) *EmailTemplate {
	return &EmailTemplate{
		Subject:  "验证码 - A Second Chance",
		Template: "template/vericode.html",
		Data: map[string]interface{}{
			"To":   to,
			"Code": strconv.Itoa(int(code)),
			"Time": time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

// GetWelcomeTemplate 获取欢迎邮件模板（可扩展）
func GetWelcomeTemplate(to string, username string) *EmailTemplate {
	return &EmailTemplate{
		Subject:  "欢迎加入 A Second Chance",
		Template: "template/welcome.html", // 可以添加更多模板
		Data: map[string]interface{}{
			"To":       to,
			"Username": username,
			"Time":     time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

// GetPasswordResetTemplate 获取密码重置邮件模板（可扩展）
func GetPasswordResetTemplate(to string, resetLink string) *EmailTemplate {
	return &EmailTemplate{
		Subject:  "密码重置 - A Second Chance",
		Template: "template/password_reset.html", // 可以添加更多模板
		Data: map[string]interface{}{
			"To":        to,
			"ResetLink": resetLink,
			"Time":      time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}
