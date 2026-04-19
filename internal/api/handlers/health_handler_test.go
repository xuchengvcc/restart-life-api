package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHealth_WithMissingDependencies_ReturnsServiceUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler := NewHealthHandler("restart-life-api", nil, nil)
	handler.Health(ctx)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)

	var resp HealthResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, "unhealthy", resp.Status)
	require.Equal(t, "not_configured", resp.Checks["database"])
	require.Equal(t, "not_configured", resp.Checks["redis"])
}

func TestReady_WithMissingDependencies_ReturnsServiceUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health/ready", nil)

	handler := NewHealthHandler("restart-life-api", nil, nil)
	handler.Ready(ctx)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)

	var resp ReadyResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.False(t, resp.Ready)
	require.Equal(t, "not ready", resp.Status)
}

func TestPing_ReturnsPong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health/ping", nil)

	handler := NewHealthHandler("restart-life-api", nil, nil)
	handler.Ping(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp PingResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Equal(t, "pong", resp.Message)
}

func TestMetrics_ContainsObservabilityFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHealthHandler("restart-life-api", nil, nil)

	// produce a failed health request to increment error counters
	healthRecorder := httptest.NewRecorder()
	healthCtx, _ := gin.CreateTestContext(healthRecorder)
	healthCtx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)
	handler.Health(healthCtx)
	require.Equal(t, http.StatusServiceUnavailable, healthRecorder.Code)

	metricsRecorder := httptest.NewRecorder()
	metricsCtx, _ := gin.CreateTestContext(metricsRecorder)
	metricsCtx.Request = httptest.NewRequest(http.MethodGet, "/health/metrics", nil)
	handler.Metrics(metricsCtx)
	require.Equal(t, http.StatusOK, metricsRecorder.Code)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(metricsRecorder.Body.Bytes(), &payload))
	require.Contains(t, payload, "request_count")
	require.Contains(t, payload, "error_count")
	require.Contains(t, payload, "error_rate")
	require.Contains(t, payload, "dependency_details")
}
