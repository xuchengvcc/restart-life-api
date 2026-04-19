package main

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/api/handlers"
	"github.com/xuchengvcc/restart-life-api/internal/api/middleware"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/dao"
	"github.com/xuchengvcc/restart-life-api/internal/job"
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
	authHandler      *handlers.AuthHandler
	aiHandler        *handlers.AIHandler
	healthHandler    *handlers.HealthHandler
	characterHandler *handlers.CharacterHandler
	gameHandler      *handlers.GameHandler
	achievementHandler *handlers.AchievementHandler
	statsHandler       *handlers.StatsHandler

	// 服务
	authService             services.AuthService
	emailService            services.EmailService
	verificationCodeService services.VerificationCodeService
	characterService        services.CharacterService
	gameService             services.GameService
	achievementStatsService services.AchievementStatsService
	aiServices              map[string]services.AIService

	// 仓库
	userRepository             repository.UserRepository
	verificationCodeRepository repository.VerificationCodeRepository
	characterRepository        repository.CharacterRepository

	// 定时任务
	jobManager *job.JobManager

	// DAO
	userDAO             dao.UserDAO
	verificationCodeDAO dao.VerificationCodeDAO
	characterDAO        dao.CharacterDAO
}

// NewContainer 创建并初始化容器
func NewContainer(cfg *config.Config, db *sql.DB, redis *redis.Client) *Container {
	container := &Container{
		cfg:   cfg,
		db:    db,
		redis: redis,
	}
	container.initUtils()
	container.initDAOs()
	container.initRepositories()
	container.initServices()
	container.initJobs()
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
	c.verificationCodeDAO = *dao.NewVerificationCodeDAO(c.redis, c.db)
	c.characterDAO = dao.NewCharacterDAO(c.db)
}

// initRepositories 初始化仓库
func (c *Container) initRepositories() {
	c.userRepository = repository.NewUserRepository(c.userDAO)
	c.verificationCodeRepository = repository.NewVerificationCodeRepository(&c.verificationCodeDAO)
	c.characterRepository = repository.NewCharacterRepository(c.characterDAO)
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

	// 角色服务
	c.characterService = services.NewCharacterService(c.characterRepository)

	// AI服务
	c.aiServices = services.NewAIServices(c.cfg.AI, c.logger)

	// 游戏服务
	c.gameService = services.NewGameService(c.characterRepository, c.aiServices, c.logger, c.redis)
	c.achievementStatsService = services.NewAchievementStatsService(c.characterRepository, c.logger)
}

// initJobs 初始化定时任务
func (c *Container) initJobs() {
	// 创建任务管理器
	c.jobManager = job.NewJobManager(c.logger)

	// 创建SSL证书监控任务
	sslCertJob := job.NewSSLCertMonitorJob(c.logger, c.emailService, &c.cfg.Email, c.redis)

	// 添加任务到管理器
	if err := c.jobManager.AddJob(sslCertJob); err != nil {
		c.logger.WithError(err).Error("failed to add SSL certificate monitor job")
	}
} // initMiddlewares 初始化中间件
func (c *Container) initMiddlewares() {
	c.authMiddleware = middleware.NewAuthMiddleware(c.authService, c.logger)
}

// initHandlers 初始化处理器
func (c *Container) initHandlers() {
	c.authHandler = handlers.NewAuthHandler(c.authService, c.verificationCodeService, c.logger)

	// AI处理器
	c.aiHandler = handlers.NewAIHandler(c.aiServices, c.logger)

	// 角色处理器
	c.characterHandler = handlers.NewCharacterHandler(c.characterService, c.logger)

	// 游戏处理器
	c.gameHandler = handlers.NewGameHandler(c.gameService, c.logger)
	c.achievementHandler = handlers.NewAchievementHandler(c.achievementStatsService, c.logger)
	c.statsHandler = handlers.NewStatsHandler(c.achievementStatsService, c.logger)

	c.healthHandler = handlers.NewHealthHandler("restart-life-api", c.db, c.redis)
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

// GetCharacterHandler 获取角色处理器
func (c *Container) GetCharacterHandler() *handlers.CharacterHandler {
	return c.characterHandler
}

// GetGameHandler 获取游戏处理器
func (c *Container) GetGameHandler() *handlers.GameHandler {
	return c.gameHandler
}

// GetAchievementHandler 获取成就处理器
func (c *Container) GetAchievementHandler() *handlers.AchievementHandler {
	return c.achievementHandler
}

// GetStatsHandler 获取统计处理器
func (c *Container) GetStatsHandler() *handlers.StatsHandler {
	return c.statsHandler
}

// GetJobManager 获取定时任务管理器
func (c *Container) GetJobManager() *job.JobManager {
	return c.jobManager
}
