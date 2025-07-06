package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/api/middleware"
	"github.com/xuchengvcc/restart-life-api/internal/config"
)

// SetupRoutes 设置所有路由和中间件
func SetupRoutes(cfg *config.Config, container Container) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建Gin引擎
	r := gin.New()

	// 注册全局中间件
	setupMiddleware(r, cfg)

	// 注册健康检查路由
	setupHealthRoutes(r, cfg, container)

	// 注册API路由
	setupAPIRoutes(r, cfg, container)

	logrus.Info("All routes setup completed")
	return r
}

// setupMiddleware 设置中间件
func setupMiddleware(r *gin.Engine, cfg *config.Config) {
	// 请求ID中间件
	r.Use(middleware.RequestIDMiddleware(middleware.DefaultRequestIDConfig()))

	// 异常恢复中间件
	r.Use(middleware.RecoveryMiddleware(middleware.DefaultRecoveryConfig()))

	// 请求日志中间件
	r.Use(middleware.LoggerMiddleware(middleware.DefaultLoggerConfig()))

	// CORS中间件
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           12 * 3600, // 12小时
	}
	r.Use(middleware.CORSMiddleware(corsConfig))

	logrus.Info("Middleware setup completed")
}

// setupAPIRoutes 设置API路由
func setupAPIRoutes(r *gin.Engine, cfg *config.Config, container Container) {
	// 获取处理器和中间件
	authHandler := container.GetAuthHandler()
	authMiddleware := container.GetAuthMiddleware()
	aiHandler := container.GetAIHandler()

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		{
			// 公开路由（不需要认证）
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// 需要认证的路由
			authProtected := auth.Group("")
			authProtected.Use(authMiddleware.RequireAuth())
			{
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.GET("/profile", authHandler.GetProfile)
				authProtected.PUT("/update_profile", authHandler.UpdateProfile)
				authProtected.POST("/change-password", authHandler.ChangePassword)
			}
		}

		// 角色相关路由（需要认证）
		characters := v1.Group("/characters")
		characters.Use(authMiddleware.RequireAuth())
		{
			// TODO: 添加角色路由
			characters.POST("", placeholderHandler("create character"))
			characters.GET("", placeholderHandler("list characters"))
			characters.GET("/:id", placeholderHandler("get character"))
			characters.PUT("/:id", placeholderHandler("update character"))
			characters.DELETE("/:id", placeholderHandler("delete character"))
		}

		// 游戏相关路由（需要认证）
		game := v1.Group("/game")
		game.Use(authMiddleware.RequireAuth())
		{
			// TODO: 添加游戏路由
			game.POST("/start/:character_id", placeholderHandler("start game"))
			game.POST("/advance/:character_id", placeholderHandler("advance game"))
			game.GET("/state/:character_id", placeholderHandler("get game state"))
			game.POST("/decision/:character_id", placeholderHandler("make decision"))
		}

		// 成就相关路由（需要认证）
		achievements := v1.Group("/achievements")
		achievements.Use(authMiddleware.RequireAuth())
		{
			// TODO: 添加成就路由
			achievements.GET("/:character_id", placeholderHandler("get achievements"))
			achievements.GET("/categories", placeholderHandler("get achievement categories"))
		}

		// 关系相关路由（需要认证）
		relationships := v1.Group("/relationships")
		relationships.Use(authMiddleware.RequireAuth())
		{
			// TODO: 添加关系路由
			relationships.GET("/:character_id", placeholderHandler("get relationships"))
			relationships.POST("/:character_id", placeholderHandler("create relationship"))
		}

		// 统计相关路由（需要认证）
		stats := v1.Group("/stats")
		stats.Use(authMiddleware.RequireAuth())
		{
			// TODO: 添加统计路由
			stats.GET("/:character_id", placeholderHandler("get character stats"))
			stats.GET("/:character_id/timeline", placeholderHandler("get timeline"))
		}

		ai := v1.Group("/ai")
		{
			ai.POST("/generate", aiHandler.Generate)
		}
	}

	logrus.Info("API routes setup completed")
}

func setupHealthRoutes(r *gin.Engine, cfg *config.Config, container Container) {
	healthHandler := container.GetHealthHandler()

	// 创建健康检查路由组，并注入配置
	health := r.Group("/health")
	health.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// 健康检查路由
	r.GET("/health", healthHandler.Health)
	r.GET("/ping", healthHandler.Ping)
	r.GET("/ready", healthHandler.Ready)
	r.GET("/version", healthHandler.Version)
	r.GET("/metrics", healthHandler.Metrics)

	// 环境信息路由（带配置注入）
	health.GET("/env", healthHandler.Environment)

	logrus.Info("Health check routes setup completed")
}

// placeholderHandler 占位符处理器，用于未实现的路由
func placeholderHandler(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API endpoint not implemented yet",
			"action":  action,
			"status":  "coming soon",
		})
	}
}

// SetupTestRoutes 设置测试路由（仅在开发环境使用）
func SetupTestRoutes(r *gin.Engine) {
	if gin.Mode() == gin.DebugMode {
		test := r.Group("/test")
		{
			// 测试panic恢复
			test.GET("/panic", func(c *gin.Context) {
				panic("test panic recovery")
			})

			// 测试日志记录
			test.GET("/log", func(c *gin.Context) {
				logrus.Info("Test log message")
				c.JSON(200, gin.H{"message": "log test completed"})
			})

			// 测试请求ID
			test.GET("/request-id", func(c *gin.Context) {
				requestID := middleware.GetRequestID(c)
				c.JSON(200, gin.H{
					"request_id": requestID,
					"message":    "request ID test",
				})
			})
		}

		logrus.Info("Test routes setup completed")
	}
}
