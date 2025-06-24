package constants

// 密码相关常量
const (
	MinPasswordLength = 6
	MaxPasswordLength = 128
	DefaultBcryptCost = 12 // 比默认值稍高，更安全
)

// 用户相关常量
const (
	MaxUsernameLength = 50
	MinUsernameLength = 3
	MaxEmailLength    = 255
	MaxBioLength      = 500
)

// 性别常量
const (
	GenderUnknown = 0
	GenderMale    = 1
	GenderFemale  = 2
	GenderOther   = 3
)

// 用户状态常量
const (
	UserStatusInactive = 0
	UserStatusActive   = 1
	UserStatusSuspended = 2
	UserStatusDeleted   = 3
)

// 用户角色常量
const (
	UserRoleUser  = "user"
	UserRoleAdmin = "admin"
	UserRoleMod   = "moderator"
)
