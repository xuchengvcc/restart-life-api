package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
)

// JWTManager JWT管理器
type JWTManager struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Type     string `json:"type"` // access or refresh
	jwt.RegisteredClaims
}

// GenerateTokenPair 生成访问Token和刷新Token
func (j *JWTManager) GenerateTokenPair(user *models.User) (accessToken, refreshToken string, expiresAt int64, err error) {
	now := time.Now()

	// 生成访问Token
	accessClaims := &Claims{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Type:     constants.TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "restart-life-api",
			Subject:   fmt.Sprintf("%d", user.UserID),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新Token
	refreshClaims := &Claims{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Type:     constants.TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "restart-life-api",
			Subject:   fmt.Sprintf("%d", user.UserID),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresAt = now.Add(j.accessExpiry).Unix()
	return accessToken, refreshToken, expiresAt, nil
}

// ValidateToken 验证Token
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// RefreshToken 刷新Token
func (j *JWTManager) RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, expiresAt int64, err error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.Type != "refresh" {
		return "", "", 0, fmt.Errorf("token is not a refresh token")
	}

	// 创建新的用户信息（仅包含必要字段）
	user := &models.User{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
	}

	return j.GenerateTokenPair(user)
}

// ExtractTokenFromHeader 从Authorization header中提取Token
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > len(constants.BearerPrefix) && authHeader[:len(constants.BearerPrefix)] == constants.BearerPrefix {
		return authHeader[len(constants.BearerPrefix):]
	}
	return ""
}
