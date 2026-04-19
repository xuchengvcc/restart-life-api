package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/xuchengvcc/restart-life-api/internal/config"
)

// HealthHandler handles health endpoints.
type HealthHandler struct {
	startTime time.Time
	version   string
	db        *sql.DB
	redis     *redis.Client
	requests  int64
	errors    int64
}

type DependencyCheckDetail struct {
	Status    string `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(version string, db *sql.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
		db:        db,
		redis:     redisClient,
	}
}

type HealthResponse struct {
	Status       string                           `json:"status"`
	Timestamp    int64                            `json:"timestamp"`
	Service      string                           `json:"service"`
	Version      string                           `json:"version"`
	Uptime       string                           `json:"uptime"`
	Checks       map[string]string                `json:"checks,omitempty"`
	CheckDetails map[string]DependencyCheckDetail `json:"check_details,omitempty"`
}

type PingResponse struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type ReadyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Ready   bool   `json:"ready"`
}

type VersionResponse struct {
	Service   string `json:"service"`
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	GitCommit string `json:"git_commit,omitempty"`
}

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

func (h *HealthHandler) dependencyChecks() (map[string]string, map[string]DependencyCheckDetail, bool) {
	checks := make(map[string]string)
	details := make(map[string]DependencyCheckDetail)
	healthy := true

	dbStart := time.Now()
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer dbCancel()
	if h.db == nil {
		checks["database"] = "not_configured"
		details["database"] = DependencyCheckDetail{Status: "not_configured", LatencyMS: 0}
		healthy = false
	} else if err := h.db.PingContext(dbCtx); err != nil {
		checks["database"] = "unhealthy"
		details["database"] = DependencyCheckDetail{Status: "unhealthy", LatencyMS: time.Since(dbStart).Milliseconds(), Error: err.Error()}
		healthy = false
	} else {
		checks["database"] = "healthy"
		details["database"] = DependencyCheckDetail{Status: "healthy", LatencyMS: time.Since(dbStart).Milliseconds()}
	}

	redisStart := time.Now()
	redisCtx, redisCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer redisCancel()
	if h.redis == nil {
		checks["redis"] = "not_configured"
		details["redis"] = DependencyCheckDetail{Status: "not_configured", LatencyMS: 0}
		healthy = false
	} else if err := h.redis.Ping(redisCtx).Err(); err != nil {
		checks["redis"] = "unhealthy"
		details["redis"] = DependencyCheckDetail{Status: "unhealthy", LatencyMS: time.Since(redisStart).Milliseconds(), Error: err.Error()}
		healthy = false
	} else {
		checks["redis"] = "healthy"
		details["redis"] = DependencyCheckDetail{Status: "healthy", LatencyMS: time.Since(redisStart).Milliseconds()}
	}

	return checks, details, healthy
}

func (h *HealthHandler) recordRequest(statusCode int) {
	atomic.AddInt64(&h.requests, 1)
	if statusCode >= http.StatusBadRequest {
		atomic.AddInt64(&h.errors, 1)
	}
}

// Health returns liveness + dependency checks.
func (h *HealthHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)
	checks, checkDetails, healthy := h.dependencyChecks()

	response := HealthResponse{
		Status:       "healthy",
		Timestamp:    time.Now().Unix(),
		Service:      "restart-life-api",
		Version:      h.version,
		Uptime:       uptime.String(),
		Checks:       checks,
		CheckDetails: checkDetails,
	}

	if !healthy {
		response.Status = "unhealthy"
	}

	statusCode := http.StatusOK
	if response.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	h.recordRequest(statusCode)
	c.JSON(statusCode, response)
}

// Ping returns a basic pong.
func (h *HealthHandler) Ping(c *gin.Context) {
	response := PingResponse{
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	}

	h.recordRequest(http.StatusOK)
	c.JSON(http.StatusOK, response)
}

// Ready checks whether required dependencies are ready.
func (h *HealthHandler) Ready(c *gin.Context) {
	_, _, ready := h.dependencyChecks()
	message := "Service is ready to accept requests"
	if !ready {
		message = "Service dependencies are not ready"
	}

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

	h.recordRequest(statusCode)
	c.JSON(statusCode, response)
}

// Version returns service version info.
func (h *HealthHandler) Version(c *gin.Context) {
	response := VersionResponse{
		Service:   "restart-life-api",
		Version:   h.version,
		BuildTime: time.Now().Format(time.RFC3339),
		GoVersion: "1.23.8",
	}

	h.recordRequest(http.StatusOK)
	c.JSON(http.StatusOK, response)
}

// Metrics returns basic runtime metrics.
func (h *HealthHandler) Metrics(c *gin.Context) {
	uptime := time.Since(h.startTime)
	checks, details, healthy := h.dependencyChecks()
	requests := atomic.LoadInt64(&h.requests)
	errors := atomic.LoadInt64(&h.errors)
	errorRate := 0.0
	if requests > 0 {
		errorRate = float64(errors) / float64(requests)
	}

	metrics := map[string]interface{}{
		"uptime_seconds":     uptime.Seconds(),
		"start_time":         h.startTime.Unix(),
		"current_time":       time.Now().Unix(),
		"version":            h.version,
		"request_count":      requests,
		"error_count":        errors,
		"error_rate":         errorRate,
		"dependencies":       checks,
		"dependency_details": details,
		"ready":              healthy,
	}

	h.recordRequest(http.StatusOK)
	c.JSON(http.StatusOK, metrics)
}

// Environment returns runtime environment info.
func (h *HealthHandler) Environment(c *gin.Context) {
	cfg, exists := c.Get("config")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration not available"})
		return
	}

	config := cfg.(*config.Config)
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

	c.Header("X-Environment", environment)
	c.Header("X-Enable-HTTP", fmt.Sprintf("%t", config.Server.EnableHTTP))
	c.Header("X-Enable-HTTPS", fmt.Sprintf("%t", config.Server.EnableHTTPS))

	h.recordRequest(http.StatusOK)
	c.JSON(http.StatusOK, response)
}
