package main

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/api/handlers"
	"github.com/xuchengvcc/restart-life-api/internal/api/middleware"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/dao"
	"github.com/xuchengvcc/restart-life-api/internal/repository"
	"github.com/xuchengvcc/restart-life-api/internal/services"
	"github.com/xuchengvcc/restart-life-api/internal/utils"
)

// Container 依赖注入容器实现
type Container struct {
	cfg   *config.Config
	db    *sql.DB
	redis *redis.Client

	// 工具
	jwtManager      *utils.JWTManager
	passwordManager *utils.PasswordManager
	logger          *logrus.Logger

	// 中间件
	authMiddleware *middleware.AuthMiddleware

	// 处理器
	authHandler   *handlers.AuthHandler
	aiHandler     *handlers.AIHandler
	healthHandler *handlers.HealthHandler

	// 服务
	authService             services.AuthService
	emailService            services.EmailService
	verificationCodeService services.VerificationCodeService
	aiServices              map[string]services.AIService

	// 仓库
	userRepository             repository.UserRepository
	verificationCodeRepository repository.VerificationCodeRepository

	// DAO
	userDAO             dao.UserDAO
	verificationCodeDAO dao.VerificationCodeDAO
}

// NewContainer 创建新的依赖注入容器
func NewContainer(cfg *config.Config, db *sql.DB, redisClient *redis.Client) *Container {
	container := &Container{
		cfg:   cfg,
		db:    db,
		redis: redisClient,
	}

	// 初始化组件
	container.initUtils()
	container.initDAOs()
	container.initRepositories()
	container.initServices()
	container.initMiddlewares()
	container.initHandlers()

	return container
}

// initUtils 初始化工具类
func (c *Container) initUtils() {
	c.logger = logrus.New()
	c.jwtManager = utils.NewJWTManager(c.cfg.Auth.JWTSecret, c.cfg.Auth.JWTExpiry, c.cfg.Auth.RefreshExpiry)
	c.passwordManager = utils.NewPasswordManager()
}

// initDAOs 初始化数据访问对象
func (c *Container) initDAOs() {
	c.userDAO = dao.NewUserDAO(c.db)
	c.verificationCodeDAO = *dao.NewVerificationCodeDAO(c.redis)
}

// initRepositories 初始化仓库
func (c *Container) initRepositories() {
	c.userRepository = repository.NewUserRepository(c.userDAO)
	c.verificationCodeRepository = repository.NewVerificationCodeRepository(&c.verificationCodeDAO)
}

// initServices 初始化服务
func (c *Container) initServices() {
	// 邮件服务
	c.emailService = services.NewEmailService(c.logger, c.cfg.Email)

	// 验证码服务 - 暂时跳过，因为依赖邮件服务
	c.verificationCodeService = services.NewVerificationCodeService(
		c.verificationCodeRepository,
		c.emailService,
		c.logger,
	)

	// 认证服务
	c.authService = services.NewAuthService(
		c.userRepository,
		c.verificationCodeService,
		c.jwtManager,
		c.passwordManager,
		c.logger,
	)

	// AI服务
	c.aiServices = services.NewAIServices(c.cfg.AI, c.logger)
}

// initMiddlewares 初始化中间件
func (c *Container) initMiddlewares() {
	c.authMiddleware = middleware.NewAuthMiddleware(c.authService, c.logger)
}

// initHandlers 初始化处理器
func (c *Container) initHandlers() {
	c.authHandler = handlers.NewAuthHandler(c.authService, c.verificationCodeService, c.logger)

	// AI处理器
	c.aiHandler = handlers.NewAIHandler(c.aiServices, c.logger)

	c.healthHandler = handlers.NewHealthHandler("restart-life-api")
}

// GetAuthMiddleware 获取认证中间件
func (c *Container) GetAuthMiddleware() *middleware.AuthMiddleware {
	return c.authMiddleware
}

// GetAuthHandler 获取认证处理器
func (c *Container) GetAuthHandler() *handlers.AuthHandler {
	return c.authHandler
}

// GetAIHandler 获取AI处理器
func (c *Container) GetAIHandler() *handlers.AIHandler {
	return c.aiHandler
}

// GetHealthHandler 获取健康检查处理器
func (c *Container) GetHealthHandler() *handlers.HealthHandler {
	return c.healthHandler
}
