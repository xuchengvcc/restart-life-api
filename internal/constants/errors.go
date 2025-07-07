package constants

import "errors"

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

	// 邮件相关错误代码 (5xxx)
	ErrCodeEmailSendFailed     = 5001
	ErrCodeEmailConfigInvalid  = 5002
	ErrCodeEmailTemplateFailed = 5003
	ErrCodeEmailConnectFailed  = 5004
	ErrCodeEmailAuthFailed     = 5005

	// 验证码相关错误代码 (6xxx)
	ErrCodeVerificationCodeExpired = 6001
	ErrCodeVerificationCodeInvalid = 6002
	ErrCodeVerificationCodeUsed    = 6003
	ErrCodeTooManyRequests         = 6004
	ErrCodeEmailAddressInvalid     = 6005
)

// 响应消息常量
const (
	MsgSuccess            = "success"
	MsgInvalidCredentials = "invalid credentials"
	MsgUserNotFound       = "user not found"
	MsgUserAlreadyExists  = "user already exists"
	MsgEmailAlreadyExists = "email already exists"
	MsgTokenExpired       = "token expired"
	MsgTokenInvalid       = "invalid token"
	MsgPermissionDenied   = "permission denied"
	MsgValidationFailed   = "validation failed"
	MsgInternalError      = "internal server error"
	MsgAccountDisabled    = "account is disabled"
	MsgPasswordIncorrect  = "old password is incorrect"
	MsgResourceNotFound   = "resource not found"
	MsgDuplicateResource  = "resource already exists"

	// 密码相关
	MsgPasswordProcessFailed = "failed to process password"
	MsgPasswordHashFailed    = "failed to hash password"
	MsgPasswordVerifyFailed  = "failed to verify password"

	// 用户操作相关
	MsgUserCreateFailed    = "failed to create user"
	MsgUserUpdateFailed    = "failed to update user"
	MsgUserDeleteFailed    = "failed to delete user"
	MsgProfileUpdateFailed = "failed to update profile"

	// Token相关
	MsgTokenGenerateFailed = "failed to generate tokens"
	MsgTokenRefreshFailed  = "failed to refresh token"
	MsgInvalidRefreshToken = "invalid refresh token"
	MsgInvalidAccessToken  = "invalid access token"
	MsgNotAccessToken      = "token is not an access token"

	// 数据库相关
	MsgDatabaseError       = "database operation failed"
	MsgRecordNotFound      = "record not found"
	MsgDuplicateEntry      = "duplicate entry"
	MsgConstraintViolation = "constraint violation"

	// 密码验证相关
	MsgPasswordTooShort  = "password too short"
	MsgPasswordTooLong   = "password too long"
	MsgPasswordNoVisible = "password must contain visible characters"
	MsgPasswordGenFailed = "failed to generate password"

	// DAO相关错误
	MsgDAOInsertFailed = "failed to insert record"
	MsgDAOUpdateFailed = "failed to update record"
	MsgDAODeleteFailed = "failed to delete record"
	MsgDAOSelectFailed = "failed to select record"
	MsgDAOCountFailed  = "failed to count records"
	MsgDAOGetIDFailed  = "failed to get last insert id"

	// JWT相关错误
	MsgJWTGenerateFailed = "failed to generate jwt token"
	MsgJWTInvalidMethod  = "unexpected signing method"
	MsgJWTParseError     = "failed to parse token"
	MsgJWTInvalidToken   = "invalid token"
	MsgJWTInvalidClaims  = "invalid token claims"

	// 邮件相关错误
	MsgEmailSendFailed            = "failed to send email"
	MsgEmailConfigInvalid         = "invalid email configuration"
	MsgEmailTemplateFailed        = "failed to process email template"
	MsgEmailConnectFailed         = "failed to connect to email server"
	MsgEmailAuthFailed            = "email authentication failed"
	MsgEmailTemplateNotFound      = "email template not found"
	MsgEmailGetWorkdirFailed      = "failed to get current working directory"
	MsgEmailReadTemplateFailed    = "failed to read email template"
	MsgEmailProcessTemplateFailed = "failed to process email template"
	MsgEmailCreateClientFailed    = "failed to create smtp client"
	MsgEmailTLSFailed             = "failed to start tls"
	MsgEmailSetSenderFailed       = "failed to set sender"
	MsgEmailSetRecipientFailed    = "failed to set recipient"
	MsgEmailGetWriterFailed       = "failed to get email writer"
	MsgEmailWriteContentFailed    = "failed to write email content"

	// 验证码相关错误
	MsgVerificationCodeExpired = "verification code has expired"
	MsgVerificationCodeInvalid = "verification code is invalid"
	MsgVerificationCodeUsed    = "verification code has been used"
	MsgTooManyRequests         = "too many requests, please try again later"
	MsgEmailAddressInvalid     = "email address is invalid"
)

// 预定义错误变量 - 用于 errors.Is 判断
var (
	// 认证相关错误
	ErrInvalidCredentials = errors.New(MsgInvalidCredentials)
	ErrUserNotFound       = errors.New(MsgUserNotFound)
	ErrUserAlreadyExists  = errors.New(MsgUserAlreadyExists)
	ErrEmailAlreadyExists = errors.New(MsgEmailAlreadyExists)
	ErrTokenExpired       = errors.New(MsgTokenExpired)
	ErrTokenInvalid       = errors.New(MsgTokenInvalid)
	ErrAccountDisabled    = errors.New(MsgAccountDisabled)
	ErrPasswordIncorrect  = errors.New(MsgPasswordIncorrect)

	// 密码相关错误
	ErrPasswordProcessFailed = errors.New(MsgPasswordProcessFailed)
	ErrPasswordHashFailed    = errors.New(MsgPasswordHashFailed)
	ErrPasswordVerifyFailed  = errors.New(MsgPasswordVerifyFailed)

	// 用户操作相关错误
	ErrUserCreateFailed    = errors.New(MsgUserCreateFailed)
	ErrUserUpdateFailed    = errors.New(MsgUserUpdateFailed)
	ErrProfileUpdateFailed = errors.New(MsgProfileUpdateFailed)

	// Token相关错误
	ErrTokenGenerateFailed = errors.New(MsgTokenGenerateFailed)
	ErrInvalidRefreshToken = errors.New(MsgInvalidRefreshToken)
	ErrNotAccessToken      = errors.New(MsgNotAccessToken)

	// 业务相关错误
	ErrValidationFailed = errors.New(MsgValidationFailed)
	ErrInternalError    = errors.New(MsgInternalError)
	ErrPermissionDenied = errors.New(MsgPermissionDenied)

	// 密码验证错误
	ErrPasswordTooShort  = errors.New(MsgPasswordTooShort)
	ErrPasswordTooLong   = errors.New(MsgPasswordTooLong)
	ErrPasswordNoVisible = errors.New(MsgPasswordNoVisible)
	ErrPasswordGenFailed = errors.New(MsgPasswordGenFailed)

	// DAO相关错误
	ErrDAOInsertFailed = errors.New(MsgDAOInsertFailed)
	ErrDAOUpdateFailed = errors.New(MsgDAOUpdateFailed)
	ErrDAODeleteFailed = errors.New(MsgDAODeleteFailed)
	ErrDAOSelectFailed = errors.New(MsgDAOSelectFailed)
	ErrDAOCountFailed  = errors.New(MsgDAOCountFailed)
	ErrDAOGetIDFailed  = errors.New(MsgDAOGetIDFailed)

	// JWT相关错误
	ErrJWTGenerateFailed = errors.New(MsgJWTGenerateFailed)
	ErrJWTInvalidMethod  = errors.New(MsgJWTInvalidMethod)
	ErrJWTParseError     = errors.New(MsgJWTParseError)
	ErrJWTInvalidToken   = errors.New(MsgJWTInvalidToken)
	ErrJWTInvalidClaims  = errors.New(MsgJWTInvalidClaims)

	// 邮件相关错误
	ErrEmailSendFailed            = errors.New(MsgEmailSendFailed)
	ErrEmailConfigInvalid         = errors.New(MsgEmailConfigInvalid)
	ErrEmailTemplateFailed        = errors.New(MsgEmailTemplateFailed)
	ErrEmailConnectFailed         = errors.New(MsgEmailConnectFailed)
	ErrEmailAuthFailed            = errors.New(MsgEmailAuthFailed)
	ErrEmailTemplateNotFound      = errors.New(MsgEmailTemplateNotFound)
	ErrEmailGetWorkdirFailed      = errors.New(MsgEmailGetWorkdirFailed)
	ErrEmailReadTemplateFailed    = errors.New(MsgEmailReadTemplateFailed)
	ErrEmailProcessTemplateFailed = errors.New(MsgEmailProcessTemplateFailed)
	ErrEmailCreateClientFailed    = errors.New(MsgEmailCreateClientFailed)
	ErrEmailTLSFailed             = errors.New(MsgEmailTLSFailed)
	ErrEmailSetSenderFailed       = errors.New(MsgEmailSetSenderFailed)
	ErrEmailSetRecipientFailed    = errors.New(MsgEmailSetRecipientFailed)
	ErrEmailGetWriterFailed       = errors.New(MsgEmailGetWriterFailed)
	ErrEmailWriteContentFailed    = errors.New(MsgEmailWriteContentFailed)

	// 验证码相关错误
	ErrVerificationCodeExpired = errors.New(MsgVerificationCodeExpired)
	ErrVerificationCodeInvalid = errors.New(MsgVerificationCodeInvalid)
	ErrVerificationCodeUsed    = errors.New(MsgVerificationCodeUsed)
	ErrTooManyRequests         = errors.New(MsgTooManyRequests)
	ErrEmailAddressInvalid     = errors.New(MsgEmailAddressInvalid)
)
