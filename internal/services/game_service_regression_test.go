package services

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xuchengvcc/restart-life-api/internal/models"
)

type fakeCharacterRepo struct {
	character             *models.Character
	eventHistory          []models.Event
	pendingDecision       *models.Decision
	savedEvents           []models.Event
	savedPendingDecisions []*models.Decision
	clearedDecision       bool
}

func (f *fakeCharacterRepo) Create(ctx context.Context, character *models.Character) error {
	return nil
}

func (f *fakeCharacterRepo) GetByID(ctx context.Context, characterID string) (*models.Character, error) {
	return f.character, nil
}

func (f *fakeCharacterRepo) GetByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	return []*models.Character{}, nil
}

func (f *fakeCharacterRepo) GetActiveByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	return []*models.Character{}, nil
}

func (f *fakeCharacterRepo) Update(ctx context.Context, character *models.Character) error {
	f.character = character
	return nil
}

func (f *fakeCharacterRepo) UpdateAttributes(ctx context.Context, characterID string, attributes *models.CharacterAttributes) error {
	return nil
}

func (f *fakeCharacterRepo) Delete(ctx context.Context, characterID string) error {
	return nil
}

func (f *fakeCharacterRepo) IsOwner(ctx context.Context, characterID string, userID uint) (bool, error) {
	return true, nil
}

func (f *fakeCharacterRepo) SaveGameEvent(ctx context.Context, event *models.Event) error {
	f.savedEvents = append(f.savedEvents, *event)
	return nil
}

func (f *fakeCharacterRepo) GetGameEventHistory(ctx context.Context, characterID string) ([]models.Event, error) {
	return f.eventHistory, nil
}

func (f *fakeCharacterRepo) SavePendingDecision(ctx context.Context, decision *models.Decision) error {
	f.savedPendingDecisions = append(f.savedPendingDecisions, decision)
	f.pendingDecision = decision
	return nil
}

func (f *fakeCharacterRepo) GetPendingDecision(ctx context.Context, characterID string) (*models.Decision, error) {
	return f.pendingDecision, nil
}

func (f *fakeCharacterRepo) ClearPendingDecision(ctx context.Context, characterID string) error {
	f.pendingDecision = nil
	f.clearedDecision = true
	return nil
}

func (f *fakeCharacterRepo) GenerateRandomAttributes() models.CharacterAttributes {
	return models.CharacterAttributes{}
}

func (f *fakeCharacterRepo) CreateCharacterSummaries(characters []*models.Character) []models.CharacterSummary {
	return []models.CharacterSummary{}
}

func TestGetEventHistory_UsesRepositoryData(t *testing.T) {
	repo := &fakeCharacterRepo{
		eventHistory: []models.Event{
			{CharacterID: "char_1", Age: 18, Description: "event 1"},
			{CharacterID: "char_1", Age: 19, Description: "event 2"},
		},
	}
	svc := &gameService{
		characterRepo: repo,
		logger:        logrus.New(),
		redisClient:   redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"}),
	}

	events, err := svc.GetEventHistory(context.Background(), "char_1")
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.Equal(t, "event 1", events[0].Description)
}

func TestBuildGameStateFromDatabase_LoadsEventsAndPendingDecision(t *testing.T) {
	pendingDecision := &models.Decision{
		CharacterID: "char_1",
		Options: &models.DecisionOption{
			Conservative: models.DecisionDetails{DecisionType: 1, OptionText: "c", Consequence: "c"},
			Moderate:     models.DecisionDetails{DecisionType: 2, OptionText: "m", Consequence: "m"},
			Aggressive:   models.DecisionDetails{DecisionType: 3, OptionText: "a", Consequence: "a"},
		},
	}

	repo := &fakeCharacterRepo{
		character: &models.Character{
			CharacterID:   "char_1",
			CharacterName: "test",
			CurrentAge:    20,
			LifeStage:     "adulthood",
			CreatedAt:     1,
			UpdatedAt:     2,
			Attributes:    models.CharacterAttributes{Intelligence: 50},
		},
		eventHistory: []models.Event{
			{CharacterID: "char_1", Age: 19, Description: "history"},
		},
		pendingDecision: pendingDecision,
	}

	svc := &gameService{
		characterRepo: repo,
		logger:        logrus.New(),
		redisClient:   redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"}),
	}

	state, err := svc.buildGameStateFromDatabase(context.Background(), "char_1")
	require.NoError(t, err)
	require.Len(t, state.KeyEvents, 1)
	require.NotNil(t, state.PendingDecision)
	require.Equal(t, "history", state.KeyEvents[0].Description)
}

func TestApplyGameProgressResponse_PersistsEventAndDecision(t *testing.T) {
	repo := &fakeCharacterRepo{
		character: &models.Character{
			CharacterID:   "char_1",
			CharacterName: "test",
			CurrentAge:    21,
			Attributes:    models.CharacterAttributes{},
		},
	}

	svc := &gameService{
		characterRepo: repo,
		logger:        logrus.New(),
		redisClient:   redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"}),
	}

	gameState := &models.GameState{
		CharacterID: "char_1",
		CurrentAge:  21,
		Attributes:  models.CharacterAttributes{},
	}
	character := &models.Character{
		CharacterID: "char_1",
	}

	resp := &AIGameProgressResponse{
		GameEnded:   false,
		HasKeyEvent: true,
		KeyEvent: &models.Event{
			Description: "new event",
			Impact:      "impact",
		},
		Decision: &models.DecisionOption{
			Conservative: models.DecisionDetails{DecisionType: 1, OptionText: "c", Consequence: "c"},
			Moderate:     models.DecisionDetails{DecisionType: 2, OptionText: "m", Consequence: "m"},
			Aggressive:   models.DecisionDetails{DecisionType: 3, OptionText: "a", Consequence: "a"},
		},
		AttributeChanges: map[string]int{
			"intelligence": 1,
		},
		CharacterUpdates: map[string]interface{}{
			"money": float64(100),
		},
		YearDescription: "year desc",
	}

	svc.applyGameProgressResponse(gameState, character, resp)
	require.Len(t, repo.savedEvents, 1)
	require.Len(t, repo.savedPendingDecisions, 1)
	require.Equal(t, int64(100), gameState.Money)
	require.Equal(t, "year desc", gameState.LastYearDescription)
}
