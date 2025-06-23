package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/api/routes"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/database"
)

func main() {
	// 加载配置
	cfg, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// 初始化日志
	initLogger(cfg)

	// 初始化 MySQL 连接
	db, err := database.InitMySQLFromConfig(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MySQL")
	}
	defer db.DB.Close()
	logrus.Info("MySQL connected successfully")

	// 初始化 Redis 连接
	redisClient, err := database.InitRedisFromConfig(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to Redis")
	}
	defer redisClient.Client.Close()
	logrus.Info("Redis connected successfully")

	// 设置路由
	r := routes.SetupRoutes(cfg)

	// 在开发环境下添加测试路由
	routes.SetupTestRoutes(r)

	// 启动服务器
	startServer(r, cfg)
}

// loadConfig 加载配置
func loadConfig() (*config.Config, error) {
	// 尝试从环境变量获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 默认配置文件路径
		configPath = "configs/development.yaml"
		if os.Getenv("CONFIG_ENV") == "live" {
			configPath = "configs/live.yaml"
		} else if os.Getenv("CONFIG_ENV") == "test" {
			configPath = "configs/live.yaml"
		}
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logrus.WithField("config_path", configPath).Warn("Config file not found, using environment variables")
		return config.LoadFromEnv(), nil
	}

	// 从文件加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	logrus.WithField("config_path", configPath).Info("Configuration loaded successfully")
	return cfg, nil
}

// initLogger 初始化日志配置
func initLogger(cfg *config.Config) {
	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		logrus.WithError(err).Warn("Invalid log level, using info level")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志格式
	switch cfg.Logging.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// 设置输出
	switch cfg.Logging.Output {
	case "stdout":
		logrus.SetOutput(os.Stdout)
	case "stderr":
		logrus.SetOutput(os.Stderr)
	default:
		logrus.SetOutput(os.Stdout)
	}

	logrus.WithFields(logrus.Fields{
		"level":  cfg.Logging.Level,
		"format": cfg.Logging.Format,
		"output": cfg.Logging.Output,
	}).Info("Logger initialized")
}

// startServer 启动HTTP服务器
func startServer(r http.Handler, cfg *config.Config) {
	// 获取端口配置
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 在goroutine中启动服务器
	go func() {
		logrus.WithFields(logrus.Fields{
			"port":          port,
			"read_timeout":  cfg.Server.ReadTimeout,
			"write_timeout": cfg.Server.WriteTimeout,
			"idle_timeout":  cfg.Server.IdleTimeout,
		}).Info("Starting HTTP server")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// 5秒超时优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Server forced to shutdown")
	}

	logrus.Info("Server stopped gracefully")
}
