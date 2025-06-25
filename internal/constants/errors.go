package constants

// 错误代码常量
const (
	// 认证相关错误代码 (1xxx)
	ErrCodeInvalidCredentials = 1001
	ErrCodeTokenExpired       = 1002
	ErrCodeTokenInvalid       = 1003
	ErrCodeUserNotFound       = 1004
	ErrCodeUserAlreadyExists  = 1005
	ErrCodePermissionDenied   = 1006

	// 业务相关错误代码 (2xxx)
	ErrCodeValidationFailed  = 2001
	ErrCodeInternalError     = 2002
	ErrCodeResourceNotFound  = 2003
	ErrCodeDuplicateResource = 2004
	ErrCodeInvalidParameter  = 2005

	// 角色相关错误代码 (3xxx)
	ErrCodeCharacterNotFound     = 3001
	ErrCodeCharacterLimitReached = 3002
	ErrCodeInvalidCharacterData  = 3003

	// 游戏相关错误代码 (4xxx)
	ErrCodeGameSessionNotFound = 4001
	ErrCodeGameActionFailed    = 4002
	ErrCodeInvalidGameState    = 4003
)

// 响应消息常量
const (
	MsgSuccess            = "success"
	MsgInvalidCredentials = "invalid credentials"
	MsgUserNotFound       = "user not found"
	MsgUserAlreadyExists  = "user already exists"
	MsgTokenExpired       = "token expired"
	MsgTokenInvalid       = "invalid token"
	MsgPermissionDenied   = "permission denied"
	MsgValidationFailed   = "validation failed"
	MsgInternalError      = "internal server error"
	MsgAccountDisabled    = "account is disabled"
	MsgPasswordIncorrect  = "old password is incorrect"
	MsgResourceNotFound   = "resource not found"
	MsgDuplicateResource  = "resource already exists"
)
