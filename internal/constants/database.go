package constants

import "time"

// 数据库相关常量
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
	MinPageSize     = 1
)

// 时间相关常量
const (
	TokenExpireHours             = 24 // Access Token 过期时间（小时）
	RefreshTokenExpireDays       = 30 // Refresh Token 过期时间（天）
	SessionTimeoutMinutes        = 30 // 会话超时时间（分钟）
	PasswordResetExpireMinute    = 10 // 密码重置Token过期时间（分钟）
	VerificationCodeExpireMinute = 10 // 验证码过期时间（分钟）
)

// 时间格式常量
const (
	TimeFormatDefault  = "2006-01-02 15:04:05"
	TimeFormatISO8601  = time.RFC3339
	TimeFormatDate     = "2006-01-02"
	TimeFormatTime     = "15:04:05"
	TimeFormatDatetime = "2006-01-02T15:04:05Z07:00"
)

// 缓存相关常量
const (
	CacheKeyUserPrefix      = "user:"
	CacheKeyCharacterPrefix = "character:"
	CacheKeySessionPrefix   = "session:"
	CacheDefaultTTL         = 3600 // 默认缓存时间（秒）
	CacheUserTTL            = 1800 // 用户缓存时间（秒）
)

// 数据库表名常量
const (
	TableUsers      = "user_tab"
	TableCharacters = "character_tab"
	TableSessions   = "session_tab"
	TableMigrations = "schema_migrations"
)
