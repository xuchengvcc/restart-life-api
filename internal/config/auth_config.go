package config

import (
	"time"

	"github.com/spf13/viper"
)

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret     string        `mapstructure:"jwt_secret"`
	JWTExpiry     time.Duration `mapstructure:"jwt_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

func setAuthConfigDefault() {
	viper.SetDefault("auth.jwt_secret", "your-dev-jwt-secret-key")
	viper.SetDefault("auth.jwt_expiry", "24h")
	viper.SetDefault("auth.refresh_expiry", "168h")
}

func setCORConfigDefault() {
	viper.SetDefault("cors.allow_origins", []string{"http://localhost:8080"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Authorization", "Accept"})
	viper.SetDefault("cors.allow_credentials", true)
}
