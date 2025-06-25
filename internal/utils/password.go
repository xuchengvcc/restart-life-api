package utils

import (
	"crypto/rand"
	"unicode"

	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"golang.org/x/crypto/bcrypt"
)

// PasswordManager 密码管理器
type PasswordManager struct {
	cost int
}

// NewPasswordManager 创建密码管理器
func NewPasswordManager(cost ...int) *PasswordManager {
	costValue := constants.DefaultBcryptCost
	if len(cost) > 0 && cost[0] >= bcrypt.MinCost && cost[0] <= bcrypt.MaxCost {
		costValue = cost[0]
	}

	return &PasswordManager{
		cost: costValue,
	}
}

// HashPassword 生成密码哈希
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	// 验证密码强度
	if err := pm.ValidatePassword(password); err != nil {
		return "", err
	}

	// 生成哈希
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", constants.ErrPasswordHashFailed
	}

	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
func (pm *PasswordManager) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePassword 验证密码强度
func (pm *PasswordManager) ValidatePassword(password string) error {
	if len(password) < constants.MinPasswordLength {
		return constants.ErrPasswordTooShort
	}

	if len(password) > constants.MaxPasswordLength {
		return constants.ErrPasswordTooLong
	}

	// 检查是否包含可见字符
	hasVisible := false
	for _, char := range password {
		if unicode.IsGraphic(char) && !unicode.IsSpace(char) {
			hasVisible = true
			break
		}
	}

	if !hasVisible {
		return constants.ErrPasswordNoVisible
	}

	return nil
}

// GenerateRandomPassword 生成随机密码
func (pm *PasswordManager) GenerateRandomPassword(length int) (string, error) {
	if length < constants.MinPasswordLength {
		length = constants.MinPasswordLength
	}
	if length > constants.MaxPasswordLength {
		length = constants.MaxPasswordLength
	}

	// 定义字符集
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	// 生成随机字节
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", constants.ErrPasswordGenFailed
	}

	// 转换为字符
	for i := 0; i < length; i++ {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}

	return string(bytes), nil
}

// EstimateStrength 估算密码强度（0-100分）
func (pm *PasswordManager) EstimateStrength(password string) int {
	if len(password) == 0 {
		return 0
	}

	score := 0

	// 长度得分 (0-40)
	length := len(password)
	if length >= 8 {
		score += 20
	}
	if length >= 12 {
		score += 10
	}
	if length >= 16 {
		score += 10
	}

	// 字符类型得分 (0-40)
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasLower {
		score += 10
	}
	if hasUpper {
		score += 10
	}
	if hasDigit {
		score += 10
	}
	if hasSpecial {
		score += 10
	}

	// 复杂度得分 (0-20)
	charTypes := 0
	if hasLower {
		charTypes++
	}
	if hasUpper {
		charTypes++
	}
	if hasDigit {
		charTypes++
	}
	if hasSpecial {
		charTypes++
	}

	score += charTypes * 5

	// 确保不超过100分
	if score > 100 {
		score = 100
	}

	return score
}
