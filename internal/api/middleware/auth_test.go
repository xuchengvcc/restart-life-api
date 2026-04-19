package middleware

import (
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

type mockAuthService struct {
	validateFn func(ctx context.Context, token string) (*utils.Claims, error)
}

func (m *mockAuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	return nil, nil
}
func (m *mockAuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	return nil, nil
}
func (m *mockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	return nil, nil
}
func (m *mockAuthService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	return nil, nil
}
func (m *mockAuthService) UpdateProfile(ctx context.Context, userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	return nil, nil
}
func (m *mockAuthService) ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error {
	return nil
}
func (m *mockAuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	return nil
}
func (m *mockAuthService) ResetPasswordWithToken(ctx context.Context, req *models.ResetPasswordWithTokenRequest) error {
	return nil
}
func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (*utils.Claims, error) {
	if m.validateFn == nil {
		return nil, constants.ErrTokenInvalid
	}
	return m.validateFn(ctx, token)
}

func TestRequireAuth_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mw := NewAuthMiddleware(&mockAuthService{}, logrus.New())
	router.GET("/protected", mw.RequireAuth(), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)

	var body models.APIResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.False(t, body.Success)
	require.NotNil(t, body.Error)
	require.Equal(t, models.ErrCodeTokenInvalid, body.Error.Code)
}

func TestRequireAuth_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mw := NewAuthMiddleware(&mockAuthService{
		validateFn: func(ctx context.Context, token string) (*utils.Claims, error) {
			return nil, constants.ErrTokenExpired
		},
	}, logrus.New())
	router.GET("/protected", mw.RequireAuth(), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)

	var body models.APIResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.False(t, body.Success)
	require.NotNil(t, body.Error)
	require.Equal(t, models.ErrCodeTokenExpired, body.Error.Code)
}

func TestRequireAuth_CaseInsensitiveBearerAccepted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mw := NewAuthMiddleware(&mockAuthService{
		validateFn: func(ctx context.Context, token string) (*utils.Claims, error) {
			require.Equal(t, "valid-token", token)
			return &utils.Claims{
				UserID:   42,
				Username: "tester",
				Email:    "tester@example.com",
				Type:     constants.TokenTypeAccess,
			}, nil
		},
	}, logrus.New())

	router.GET("/protected", mw.RequireAuth(), func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		require.True(t, ok)
		require.Equal(t, uint(42), userID)
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "bearer valid-token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNoContent, resp.Code)
}
