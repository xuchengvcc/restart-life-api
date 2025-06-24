package constants

// HTTP Headers
const (
	HeaderAuthorization  = "Authorization"
	HeaderContentType    = "Content-Type"
	HeaderRequestID      = "X-Request-ID"
	HeaderUserAgent      = "User-Agent"
	HeaderAcceptLanguage = "Accept-Language"
	HeaderCacheControl   = "Cache-Control"
	HeaderETag           = "ETag"
	HeaderIfNoneMatch    = "If-None-Match"
)

// JWT相关常量
const (
	BearerPrefix     = "Bearer "
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// 上下文键常量
const (
	ContextKeyRequestID = "request_id"
	ContextKeyUserID    = "user_id"
	ContextKeyClaims    = "claims"
	ContextKeyUser      = "user"
	ContextKeyLogger    = "logger"
)

// 日志字段常量
const (
	LogFieldUserID    = "user_id"
	LogFieldUsername  = "username"
	LogFieldRequestID = "request_id"
	LogFieldMethod    = "method"
	LogFieldPath      = "path"
	LogFieldIP        = "ip"
	LogFieldUserAgent = "user_agent"
	LogFieldLatency   = "latency"
	LogFieldStatus    = "status"
	LogFieldError     = "error"
	LogFieldAction    = "action"
	LogFieldResource  = "resource"
)
