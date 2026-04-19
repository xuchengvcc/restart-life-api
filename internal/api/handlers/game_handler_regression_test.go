package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

type mockGameService struct {
	receivedUserID uint64
}

func (m *mockGameService) StartOrResumeGame(ctx context.Context, userID uint64) (*models.GameState, error) {
	m.receivedUserID = userID
	return &models.GameState{
		CharacterID:   "char_1",
		CharacterName: "test",
		IsGameActive:  true,
	}, nil
}

func (m *mockGameService) AdvanceGame(ctx context.Context, characterID string) (*models.GameState, error) {
	return &models.GameState{}, nil
}

func (m *mockGameService) AdvanceGameSmart(ctx context.Context, characterID string, optionType string) (*models.GameState, error) {
	return &models.GameState{}, nil
}

func (m *mockGameService) MakeDecision(ctx context.Context, characterID string, optionType string) (*models.GameState, error) {
	return &models.GameState{}, nil
}

func (m *mockGameService) GetGameState(ctx context.Context, characterID string) (*models.GameState, error) {
	return &models.GameState{}, nil
}

func (m *mockGameService) GetEventHistory(ctx context.Context, characterID string) ([]models.Event, error) {
	return []models.Event{}, nil
}

func (m *mockGameService) SaveGame(ctx context.Context, characterID string) error {
	return nil
}

func (m *mockGameService) LoadGame(ctx context.Context, characterID string) (*models.GameState, error) {
	return &models.GameState{}, nil
}

var _ services.GameService = (*mockGameService)(nil)

func TestStartOrResumeGame_UserIDTypeCompatibility(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		userID         interface{}
		expectedUserID uint64
	}{
		{name: "uint", userID: uint(7), expectedUserID: 7},
		{name: "uint64", userID: uint64(8), expectedUserID: 8},
		{name: "int", userID: int(9), expectedUserID: 9},
		{name: "string", userID: "10", expectedUserID: 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/game/start-or-resume", nil)
			ctx.Set("user_id", tc.userID)

			mockSvc := &mockGameService{}
			handler := NewGameHandler(mockSvc, logrus.New())
			handler.StartOrResumeGame(ctx)

			require.Equal(t, http.StatusOK, recorder.Code)
			require.Equal(t, tc.expectedUserID, mockSvc.receivedUserID)
		})
	}
}

