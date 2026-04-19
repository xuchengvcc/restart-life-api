package services

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xuchengvcc/restart-life-api/internal/models"
)

type fakeAchievementRepo struct {
	isOwner            bool
	isOwnerErr         error
	character          *models.Character
	characterErr       error
	characters         []*models.Character
	charactersErr      error
	events             []models.Event
	eventsErr          error
	pendingDecision    *models.Decision
	pendingDecisionErr error
}

func (f *fakeAchievementRepo) Create(ctx context.Context, character *models.Character) error { return nil }
func (f *fakeAchievementRepo) GetByID(ctx context.Context, characterID string) (*models.Character, error) {
	if f.characterErr != nil {
		return nil, f.characterErr
	}
	return f.character, nil
}
func (f *fakeAchievementRepo) GetByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	if f.charactersErr != nil {
		return nil, f.charactersErr
	}
	return f.characters, nil
}
func (f *fakeAchievementRepo) GetActiveByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	return []*models.Character{}, nil
}
func (f *fakeAchievementRepo) Update(ctx context.Context, character *models.Character) error { return nil }
func (f *fakeAchievementRepo) UpdateAttributes(ctx context.Context, characterID string, attributes *models.CharacterAttributes) error {
	return nil
}
func (f *fakeAchievementRepo) Delete(ctx context.Context, characterID string) error { return nil }
func (f *fakeAchievementRepo) IsOwner(ctx context.Context, characterID string, userID uint) (bool, error) {
	if f.isOwnerErr != nil {
		return false, f.isOwnerErr
	}
	return f.isOwner, nil
}
func (f *fakeAchievementRepo) SaveGameEvent(ctx context.Context, event *models.Event) error { return nil }
func (f *fakeAchievementRepo) GetGameEventHistory(ctx context.Context, characterID string) ([]models.Event, error) {
	if f.eventsErr != nil {
		return nil, f.eventsErr
	}
	return f.events, nil
}
func (f *fakeAchievementRepo) SavePendingDecision(ctx context.Context, decision *models.Decision) error {
	return nil
}
func (f *fakeAchievementRepo) GetPendingDecision(ctx context.Context, characterID string) (*models.Decision, error) {
	if f.pendingDecisionErr != nil {
		return nil, f.pendingDecisionErr
	}
	return f.pendingDecision, nil
}
func (f *fakeAchievementRepo) ClearPendingDecision(ctx context.Context, characterID string) error { return nil }
func (f *fakeAchievementRepo) GenerateRandomAttributes() models.CharacterAttributes {
	return models.CharacterAttributes{}
}
func (f *fakeAchievementRepo) CreateCharacterSummaries(characters []*models.Character) []models.CharacterSummary {
	return []models.CharacterSummary{}
}

func TestAchievementStatsService_GetAchievements(t *testing.T) {
	repo := &fakeAchievementRepo{
		isOwner: true,
		character: &models.Character{
			CharacterID:   "char_1",
			CharacterName: "tester",
			CurrentAge:    61,
			UpdatedAt:     100,
			GameCompleted: true,
			Money:         1_000_000,
			Attributes: models.CharacterAttributes{
				Intelligence:          85,
				EmotionalIntelligence: 88,
			},
		},
		events: make([]models.Event, 12),
	}
	svc := NewAchievementStatsService(repo, logrus.New())

	resp, err := svc.GetAchievements(context.Background(), "char_1", 1)
	require.NoError(t, err)
	require.Equal(t, "char_1", resp.CharacterID)
	require.Equal(t, len(resp.Items), resp.TotalCount)
	require.GreaterOrEqual(t, resp.UnlockedCount, 1)
}

func TestAchievementStatsService_GetAchievements_NotOwner(t *testing.T) {
	repo := &fakeAchievementRepo{isOwner: false}
	svc := NewAchievementStatsService(repo, logrus.New())

	resp, err := svc.GetAchievements(context.Background(), "char_1", 1)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestAchievementStatsService_GetAchievementCategories(t *testing.T) {
	repo := &fakeAchievementRepo{
		characters: []*models.Character{
			{
				CurrentAge:    70,
				GameCompleted: true,
				Money:         1_000_000,
				Attributes: models.CharacterAttributes{
					Intelligence:          90,
					EmotionalIntelligence: 91,
				},
			},
		},
	}
	svc := NewAchievementStatsService(repo, logrus.New())

	categories, err := svc.GetAchievementCategories(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, categories, 4)
	for _, c := range categories {
		require.LessOrEqual(t, c.UnlockedCount, c.TotalCount)
	}
}

func TestAchievementStatsService_GetCharacterStats(t *testing.T) {
	repo := &fakeAchievementRepo{
		isOwner: true,
		character: &models.Character{
			CharacterID:   "char_2",
			CharacterName: "stats",
			CurrentAge:    22,
			LifeStage:     "young_adult",
			TotalPlaytime: 123,
			Money:         999,
			Attributes:    models.CharacterAttributes{Intelligence: 50},
		},
		events: []models.Event{
			{Description: "first"},
			{Description: "last"},
		},
		pendingDecision: &models.Decision{CharacterID: "char_2"},
	}
	svc := NewAchievementStatsService(repo, logrus.New())

	stats, err := svc.GetCharacterStats(context.Background(), "char_2", 1)
	require.NoError(t, err)
	require.Equal(t, "char_2", stats.CharacterID)
	require.Equal(t, 2, stats.EventCount)
	require.True(t, stats.PendingDecision)
	require.NotNil(t, stats.LastEvent)
	require.Equal(t, "last", stats.LastEvent.Description)
}

func TestAchievementStatsService_GetCharacterTimeline(t *testing.T) {
	repo := &fakeAchievementRepo{
		isOwner: true,
		events: []models.Event{
			{Age: 18, Description: "event_1", Impact: "none", CreatedAt: 1},
			{Age: 19, Description: "event_2", Impact: "up", CreatedAt: 2},
		},
	}
	svc := NewAchievementStatsService(repo, logrus.New())

	resp, err := svc.GetCharacterTimeline(context.Background(), "char_3", 1)
	require.NoError(t, err)
	require.Equal(t, "char_3", resp.CharacterID)
	require.Equal(t, 2, resp.Total)
	require.Equal(t, "event_1", resp.Timeline[0].Description)
}

func TestAchievementStatsService_GetCharacterTimeline_OwnerCheckError(t *testing.T) {
	repo := &fakeAchievementRepo{
		isOwnerErr: errors.New("db failure"),
	}
	svc := NewAchievementStatsService(repo, logrus.New())

	resp, err := svc.GetCharacterTimeline(context.Background(), "char_3", 1)
	require.Error(t, err)
	require.Nil(t, resp)
}
