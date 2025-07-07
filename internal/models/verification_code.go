package models

// VerificationCode 验证码模型（简化版）
type VerificationCode struct {
	Code string `json:"code"` // 验证码
	Type string `json:"type"` // 验证码类型
}

// VerificationCodeType 验证码类型常量
const (
	VerificationCodeTypeRegister      = "register"
	VerificationCodeTypeResetPassword = "reset_password"
	VerificationCodeTypeChangeEmail   = "change_email"
)
