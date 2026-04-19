package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/utils"
)

type mockAuthServiceForHandler struct {
	refreshFn func(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
}

func (m *mockAuthServiceForHandler) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	return nil, nil
}
func (m *mockAuthServiceForHandler) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	return nil, nil
}
func (m *mockAuthServiceForHandler) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	if m.refreshFn == nil {
		return nil, constants.ErrInvalidRefreshToken
	}
	return m.refreshFn(ctx, refreshToken)
}
func (m *mockAuthServiceForHandler) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	return nil, nil
}
func (m *mockAuthServiceForHandler) UpdateProfile(ctx context.Context, userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	return nil, nil
}
func (m *mockAuthServiceForHandler) ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error {
	return nil
}
func (m *mockAuthServiceForHandler) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	return nil
}
func (m *mockAuthServiceForHandler) ResetPasswordWithToken(ctx context.Context, req *models.ResetPasswordWithTokenRequest) error {
	return nil
}
func (m *mockAuthServiceForHandler) ValidateToken(ctx context.Context, token string) (*utils.Claims, error) {
	return nil, nil
}

func TestRefreshToken_ExpiredToken_ReturnsExpiredCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(
		&mockAuthServiceForHandler{
			refreshFn: func(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
				return nil, constants.ErrTokenExpired
			},
		},
		nil,
		logrus.New(),
	)

	body := bytes.NewBufferString(`{"refresh_token":"expired-token"}`)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", body)
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(ctx)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)

	var resp models.APIResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.False(t, resp.Success)
	require.NotNil(t, resp.Error)
	require.Equal(t, models.ErrCodeTokenExpired, resp.Error.Code)
}

func TestRefreshToken_InvalidToken_ReturnsInvalidCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(
		&mockAuthServiceForHandler{
			refreshFn: func(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
				return nil, constants.ErrInvalidRefreshToken
			},
		},
		nil,
		logrus.New(),
	)

	body := bytes.NewBufferString(`{"refresh_token":"invalid-token"}`)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", body)
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(ctx)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)

	var resp models.APIResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.False(t, resp.Success)
	require.NotNil(t, resp.Error)
	require.Equal(t, models.ErrCodeTokenInvalid, resp.Error.Code)
}
