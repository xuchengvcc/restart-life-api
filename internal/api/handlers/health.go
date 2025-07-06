package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuchengvcc/restart-life-api/internal/config"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp int64             `json:"timestamp"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// PingResponse Ping响应
type PingResponse struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// ReadyResponse 就绪检查响应
type ReadyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Ready   bool   `json:"ready"`
}

// VersionResponse 版本信息响应
type VersionResponse struct {
	Service   string `json:"service"`
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	GitCommit string `json:"git_commit,omitempty"`
}

// EnvironmentResponse 环境信息响应
type EnvironmentResponse struct {
	Environment string `json:"environment"`
	Mode        string `json:"mode"`
	EnableHTTP  bool   `json:"enable_http"`
	EnableHTTPS bool   `json:"enable_https"`
	Port        string `json:"port"`
	HTTPSPort   string `json:"https_port"`
	Version     string `json:"version"`
	Timestamp   int64  `json:"timestamp"`
}

// Health 健康检查接口
// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Unix(),
		Service:   "restart-life-api",
		Version:   h.version,
		Uptime:    uptime.String(),
	}

	// 可以在这里添加更多的健康检查
	// 例如数据库连接、Redis连接等
	checks := make(map[string]string)

	// TODO: 添加数据库健康检查
	// if err := h.db.HealthCheck(); err != nil {
	//     checks["database"] = "unhealthy"
	//     response.Status = "unhealthy"
	// } else {
	//     checks["database"] = "healthy"
	// }

	// TODO: 添加Redis健康检查
	// if err := h.redis.HealthCheck(); err != nil {
	//     checks["redis"] = "unhealthy"
	//     response.Status = "unhealthy"
	// } else {
	//     checks["redis"] = "healthy"
	// }

	if len(checks) > 0 {
		response.Checks = checks
	}

	// 根据健康状态返回相应的HTTP状态码
	statusCode := http.StatusOK
	if response.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Ping 基础连通性检查
// @Summary Ping检查
// @Description 基础连通性检查
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	response := PingResponse{
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, response)
}

// Ready 服务就绪检查
// @Summary 就绪检查
// @Description 检查服务是否准备好接受请求
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} ReadyResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	// 检查服务是否准备就绪
	// 这里可以添加更复杂的就绪检查逻辑
	ready := true
	message := "Service is ready to accept requests"

	// TODO: 添加就绪检查逻辑
	// 例如检查必要的依赖服务是否可用

	response := ReadyResponse{
		Status:  "ready",
		Message: message,
		Ready:   ready,
	}

	statusCode := http.StatusOK
	if !ready {
		response.Status = "not ready"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Version 版本信息
// @Summary 版本信息
// @Description 获取服务版本信息
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} VersionResponse
// @Router /version [get]
func (h *HealthHandler) Version(c *gin.Context) {
	response := VersionResponse{
		Service:   "restart-life-api",
		Version:   h.version,
		BuildTime: time.Now().Format(time.RFC3339),
		GoVersion: "1.23.8",
		// GitCommit: 可以在构建时注入Git提交哈希
	}

	c.JSON(http.StatusOK, response)
}

// Metrics 基础指标信息（为Prometheus等监控系统准备）
// @Summary 基础指标
// @Description 获取基础指标信息
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	uptime := time.Since(h.startTime)

	metrics := map[string]interface{}{
		"uptime_seconds": uptime.Seconds(),
		"start_time":     h.startTime.Unix(),
		"current_time":   time.Now().Unix(),
		"version":        h.version,
		// TODO: 添加更多指标
		// "active_connections": getActiveConnections(),
		// "memory_usage": getMemoryUsage(),
		// "cpu_usage": getCPUUsage(),
	}

	c.JSON(http.StatusOK, metrics)
}

// Environment 环境信息接口
// @Summary 环境信息检查
// @Description 获取当前服务的环境配置信息
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} EnvironmentResponse
// @Router /health/env [get]
func (h *HealthHandler) Environment(c *gin.Context) {
	// 从上下文获取配置（需要在路由中设置）
	cfg, exists := c.Get("config")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration not available"})
		return
	}

	config := cfg.(*config.Config)

	// 判断环境类型
	environment := "test"
	if !config.Server.EnableHTTP {
		environment = "live"
	}

	response := EnvironmentResponse{
		Environment: environment,
		Mode:        config.Server.Mode,
		EnableHTTP:  config.Server.EnableHTTP,
		EnableHTTPS: config.Server.EnableHTTPS,
		Port:        config.Server.Port,
		HTTPSPort:   config.Server.HTTPSPort,
		Version:     h.version,
		Timestamp:   time.Now().Unix(),
	}

	// 添加响应头以供nginx使用
	c.Header("X-Environment", environment)
	c.Header("X-Enable-HTTP", fmt.Sprintf("%t", config.Server.EnableHTTP))
	c.Header("X-Enable-HTTPS", fmt.Sprintf("%t", config.Server.EnableHTTPS))

	c.JSON(http.StatusOK, response)
}
