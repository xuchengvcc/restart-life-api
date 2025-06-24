package constants

// API版本
const (
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion
)

// API路由常量
const (
	// 认证相关路由
	RouteAuthRegister      = "/auth/register"
	RouteAuthLogin         = "/auth/login"
	RouteAuthLogout        = "/auth/logout"
	RouteAuthRefresh       = "/auth/refresh"
	RouteAuthProfile       = "/auth/profile"
	RouteAuthChangePassword = "/auth/change-password"

	// 角色相关路由
	RouteCharacters           = "/characters"
	RouteCharacterByID        = "/characters/:id"
	RouteCharacterAttributes  = "/characters/:id/attributes"
	RouteCharacterRelationships = "/characters/:id/relationships"

	// 游戏相关路由
	RouteGameStart    = "/game/start"
	RouteGameAction   = "/game/action"
	RouteGameStatus   = "/game/status"
	RouteGameHistory  = "/game/history"

	// 健康检查
	RouteHealth = "/health"
	RouteReady  = "/ready"
)

// HTTP状态码扩展
const (
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429
)

// Content-Type 常量
const (
	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
	ContentTypeText = "text/plain"
	ContentTypeHTML = "text/html"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

// API限制常量
const (
	MaxRequestsPerMinute = 60
	MaxRequestsPerHour   = 1000
	MaxRequestsPerDay    = 10000
)
