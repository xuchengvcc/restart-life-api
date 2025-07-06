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
	// config的其他初始化
	cfg.PostInit()

	// 根据配置设置nginx环境
	setNginxEnvironment(cfg)

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

	// 初始化依赖注入容器
	container := NewContainer(cfg, db.DB, redisClient.Client)
	logrus.Info("Dependency injection container initialized")

	// 设置路由
	r := routes.SetupRoutes(cfg, container)

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
		configPath = "configs/docker.yaml"
		if os.Getenv("CONFIG_ENV") == "live" {
			configPath = "configs/live.yaml"
		} else if os.Getenv("CONFIG_ENV") == "test" {
			configPath = "configs/test.yaml"
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

// setNginxEnvironment 根据配置设置nginx环境
func setNginxEnvironment(cfg *config.Config) {
	httpDisableFlag := "/tmp/nginx_disable_http"

	// 根据配置决定是否禁用HTTP
	if !cfg.Server.EnableHTTP {
		// 创建禁用HTTP的标识文件
		if file, err := os.Create(httpDisableFlag); err != nil {
			logrus.WithError(err).Warn("Failed to create nginx HTTP disable flag")
		} else {
			file.Close()
			logrus.Info("HTTP access disabled for nginx (live environment)")
		}
	} else {
		// 删除禁用HTTP的标识文件（如果存在）
		if err := os.Remove(httpDisableFlag); err != nil && !os.IsNotExist(err) {
			logrus.WithError(err).Warn("Failed to remove nginx HTTP disable flag")
		} else if err == nil {
			logrus.Info("HTTP access enabled for nginx (test environment)")
		}
	}
}

// startServer 启动HTTP和/或HTTPS服务器
func startServer(r http.Handler, cfg *config.Config) {
	var servers []*http.Server

	// 根据配置启动HTTP服务器
	if cfg.Server.EnableHTTP {
		httpPort := cfg.Server.Port
		if httpPort == "" {
			httpPort = "8080"
		}

		httpServer := &http.Server{
			Addr:         "0.0.0.0:" + httpPort,
			Handler:      r,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		}

		servers = append(servers, httpServer)

		// 在goroutine中启动HTTP服务器
		go func() {
			logrus.WithFields(logrus.Fields{
				"port":          httpPort,
				"protocol":      "HTTP",
				"read_timeout":  cfg.Server.ReadTimeout,
				"write_timeout": cfg.Server.WriteTimeout,
				"idle_timeout":  cfg.Server.IdleTimeout,
			}).Info("Starting HTTP server")

			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.WithError(err).Fatal("Failed to start HTTP server")
			}
		}()
	} else {
		logrus.Info("HTTP server disabled by configuration")
	}

	// 根据配置启动HTTPS服务器
	if cfg.Server.EnableHTTPS {
		httpsPort := cfg.Server.HTTPSPort
		if httpsPort == "" {
			httpsPort = "8443"
		}

		httpsServer := &http.Server{
			Addr:         "0.0.0.0:" + httpsPort,
			Handler:      r,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		}

		servers = append(servers, httpsServer)

		// 在goroutine中启动HTTPS服务器
		go func() {
			logrus.WithFields(logrus.Fields{
				"port":          httpsPort,
				"protocol":      "HTTPS",
				"cert_file":     cfg.Server.SSLCertFile,
				"key_file":      cfg.Server.SSLKeyFile,
				"read_timeout":  cfg.Server.ReadTimeout,
				"write_timeout": cfg.Server.WriteTimeout,
				"idle_timeout":  cfg.Server.IdleTimeout,
			}).Info("Starting HTTPS server")

			if err := httpsServer.ListenAndServeTLS(cfg.Server.SSLCertFile, cfg.Server.SSLKeyFile); err != nil && err != http.ErrServerClosed {
				logrus.WithError(err).Fatal("Failed to start HTTPS server")
			}
		}()
	} else {
		logrus.Info("HTTPS server disabled by configuration")
	}

	// 检查是否至少启动了一个服务器
	if len(servers) == 0 {
		logrus.Fatal("No servers enabled. Please enable at least one of HTTP or HTTPS")
	}

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down servers...")

	// 5秒超时优雅关闭所有服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 并行关闭所有服务器
	for i, server := range servers {
		go func(idx int, srv *http.Server) {
			if err := srv.Shutdown(ctx); err != nil {
				logrus.WithError(err).WithField("server_index", idx).Error("Server forced to shutdown")
			}
		}(i, server)
	}

	// 等待所有服务器关闭
	<-ctx.Done()
	logrus.Info("All servers stopped gracefully")
}
