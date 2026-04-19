package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xuchengvcc/restart-life-api/internal/models"
)

type mockAchievementStatsService struct {
	achievementResp       *models.CharacterAchievementsResponse
	achievementErr        error
	categoriesResp        []models.AchievementCategory
	categoriesErr         error
	statsResp             *models.CharacterStatsResponse
	statsErr              error
	timelineResp          *models.CharacterTimelineResponse
	timelineErr           error
	lastCharacterID       string
	lastUserID            uint
	getAchievementsCalled bool
	getCategoriesCalled   bool
	getStatsCalled        bool
	getTimelineCalled     bool
}

func (m *mockAchievementStatsService) GetAchievements(ctx context.Context, characterID string, userID uint) (*models.CharacterAchievementsResponse, error) {
	m.getAchievementsCalled = true
	m.lastCharacterID = characterID
	m.lastUserID = userID
	return m.achievementResp, m.achievementErr
}

func (m *mockAchievementStatsService) GetAchievementCategories(ctx context.Context, userID uint) ([]models.AchievementCategory, error) {
	m.getCategoriesCalled = true
	m.lastUserID = userID
	return m.categoriesResp, m.categoriesErr
}

func (m *mockAchievementStatsService) GetCharacterStats(ctx context.Context, characterID string, userID uint) (*models.CharacterStatsResponse, error) {
	m.getStatsCalled = true
	m.lastCharacterID = characterID
	m.lastUserID = userID
	return m.statsResp, m.statsErr
}

func (m *mockAchievementStatsService) GetCharacterTimeline(ctx context.Context, characterID string, userID uint) (*models.CharacterTimelineResponse, error) {
	m.getTimelineCalled = true
	m.lastCharacterID = characterID
	m.lastUserID = userID
	return m.timelineResp, m.timelineErr
}

func setupHandlerTestContext(method, url string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, url, nil)
	return ctx, recorder
}

func decodeAPIResponse(t *testing.T, recorder *httptest.ResponseRecorder) models.APIResponse {
	t.Helper()
	var response models.APIResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	return response
}

func TestAchievementHandler_GetCharacterAchievements_Success(t *testing.T) {
	ctx, recorder := setupHandlerTestContext(http.MethodGet, "/api/v1/achievements/char_1")
	ctx.Set("user_id", uint(7))
	ctx.Params = gin.Params{{Key: "character_id", Value: "char_1"}}

	mockSvc := &mockAchievementStatsService{
		achievementResp: &models.CharacterAchievementsResponse{
			CharacterID:   "char_1",
			CharacterName: "tester",
			Items:         []models.AchievementItem{{ID: "age_18", Unlocked: true}},
			UnlockedCount: 1,
			TotalCount:    1,
		},
	}

	handler := NewAchievementHandler(mockSvc, logrus.New())
	handler.GetCharacterAchievements(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.True(t, mockSvc.getAchievementsCalled)
	require.Equal(t, "char_1", mockSvc.lastCharacterID)
	require.Equal(t, uint(7), mockSvc.lastUserID)

	response := decodeAPIResponse(t, recorder)
	require.True(t, response.Success)
	require.NotNil(t, response.Data)
}

func TestAchievementHandler_GetCharacterAchievements_Unauthorized(t *testing.T) {
	ctx, recorder := setupHandlerTestContext(http.MethodGet, "/api/v1/achievements/char_1")
	ctx.Params = gin.Params{{Key: "character_id", Value: "char_1"}}

	mockSvc := &mockAchievementStatsService{}
	handler := NewAchievementHandler(mockSvc, logrus.New())
	handler.GetCharacterAchievements(ctx)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.False(t, mockSvc.getAchievementsCalled)
}

func TestAchievementHandler_GetCharacterAchievements_ServiceError(t *testing.T) {
	ctx, recorder := setupHandlerTestContext(http.MethodGet, "/api/v1/achievements/char_1")
	ctx.Set("user_id", uint(7))
	ctx.Params = gin.Params{{Key: "character_id", Value: "char_1"}}

	mockSvc := &mockAchievementStatsService{achievementErr: errors.New("repo failed")}
	handler := NewAchievementHandler(mockSvc, logrus.New())
	handler.GetCharacterAchievements(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	response := decodeAPIResponse(t, recorder)
	require.False(t, response.Success)
	require.NotNil(t, response.Error)
}

func TestAchievementHandler_GetAchievementCategories_Success(t *testing.T) {
	ctx, recorder := setupHandlerTestContext(http.MethodGet, "/api/v1/achievements/categories")
	ctx.Set("user_id", uint(8))

	mockSvc := &mockAchievementStatsService{
		categoriesResp: []models.AchievementCategory{
			{ID: "milestone", TotalCount: 4, UnlockedCount: 2},
		},
	}

	handler := NewAchievementHandler(mockSvc, logrus.New())
	handler.GetAchievementCategories(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.True(t, mockSvc.getCategoriesCalled)
	require.Equal(t, uint(8), mockSvc.lastUserID)
}

func TestStatsHandler_GetCharacterStats_Success(t *testing.T) {
	ctx, recorder := setupHandlerTestContext(http.MethodGet, "/api/v1/stats/char_2")
	ctx.Set("user_id", uint(9))
	ctx.Params = gin.Params{{Key: "character_id", Value: "char_2"}}

	mockSvc := &mockAchievementStatsService{
		statsResp: &models.CharacterStatsResponse{
			CharacterID:    "char_2",
			CharacterName:  "stats",
			PendingDecision: true,
		},
	}

	handler := NewStatsHandler(mockSvc, logrus.New())
	handler.GetCharacterStats(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.True(t, mockSvc.getStatsCalled)
	require.Equal(t, "char_2", mockSvc.lastCharacterID)
	require.Equal(t, uint(9), mockSvc.lastUserID)
}

func TestStatsHandler_GetCharacterTimeline_Success(t *testing.T) {
	ctx, recorder := setupHandlerTestContext(http.MethodGet, "/api/v1/stats/char_3/timeline")
	ctx.Set("user_id", uint(10))
	ctx.Params = gin.Params{{Key: "character_id", Value: "char_3"}}

	mockSvc := &mockAchievementStatsService{
		timelineResp: &models.CharacterTimelineResponse{
			CharacterID: "char_3",
			Timeline:    []models.TimelineItem{{Age: 20, Description: "event"}},
			Total:       1,
		},
	}

	handler := NewStatsHandler(mockSvc, logrus.New())
	handler.GetCharacterTimeline(ctx)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.True(t, mockSvc.getTimelineCalled)
	require.Equal(t, "char_3", mockSvc.lastCharacterID)
	require.Equal(t, uint(10), mockSvc.lastUserID)
}

func TestStatsHandler_GetCharacterTimeline_ServiceError(t *testing.T) {
	ctx, recorder := setupHandlerTestContext(http.MethodGet, "/api/v1/stats/char_3/timeline")
	ctx.Set("user_id", uint(10))
	ctx.Params = gin.Params{{Key: "character_id", Value: "char_3"}}

	mockSvc := &mockAchievementStatsService{timelineErr: errors.New("timeline failed")}
	handler := NewStatsHandler(mockSvc, logrus.New())
	handler.GetCharacterTimeline(ctx)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	response := decodeAPIResponse(t, recorder)
	require.False(t, response.Success)
	require.NotNil(t, response.Error)
}
