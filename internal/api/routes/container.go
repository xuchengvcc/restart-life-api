package routes

import (
	"github.com/xuchengvcc/restart-life-api/internal/api/handlers"
	"github.com/xuchengvcc/restart-life-api/internal/api/middleware"
)

// Container 依赖注入容器接口
type Container interface {
	GetAuthHandler() *handlers.AuthHandler
	GetAuthMiddleware() *middleware.AuthMiddleware
}
