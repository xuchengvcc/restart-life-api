package services

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/repository"
)

type AchievementStatsService interface {
	GetAchievements(ctx context.Context, characterID string, userID uint) (*models.CharacterAchievementsResponse, error)
	GetAchievementCategories(ctx context.Context, userID uint) ([]models.AchievementCategory, error)
	GetCharacterStats(ctx context.Context, characterID string, userID uint) (*models.CharacterStatsResponse, error)
	GetCharacterTimeline(ctx context.Context, characterID string, userID uint) (*models.CharacterTimelineResponse, error)
}

type achievementStatsService struct {
	characterRepo repository.CharacterRepository
	logger        *logrus.Logger
}

func NewAchievementStatsService(characterRepo repository.CharacterRepository, logger *logrus.Logger) AchievementStatsService {
	return &achievementStatsService{
		characterRepo: characterRepo,
		logger:        logger,
	}
}

func (s *achievementStatsService) ensureOwnership(ctx context.Context, characterID string, userID uint) error {
	isOwner, err := s.characterRepo.IsOwner(ctx, characterID, userID)
	if err != nil {
		return fmt.Errorf("check character ownership failed: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("character does not belong to user")
	}
	return nil
}

func clampProgress(progress, max int) int {
	if progress < 0 {
		return 0
	}
	if progress > max {
		return max
	}
	return progress
}

func (s *achievementStatsService) GetAchievements(ctx context.Context, characterID string, userID uint) (*models.CharacterAchievementsResponse, error) {
	if err := s.ensureOwnership(ctx, characterID, userID); err != nil {
		return nil, err
	}

	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get character failed: %w", err)
	}

	events, err := s.characterRepo.GetGameEventHistory(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get event history failed: %w", err)
	}

	million := int64(1_000_000)
	tenThousand := int64(10_000)
	hundredThousand := int64(100_000)

	defs := []struct {
		id          string
		title       string
		description string
		category    string
		progress    int
		max         int
		unlocked    bool
	}{
		{
			id:          "age_18",
			title:       "成年礼",
			description: "角色达到 18 岁",
			category:    "milestone",
			progress:    character.CurrentAge,
			max:         18,
			unlocked:    character.CurrentAge >= 18,
		},
		{
			id:          "age_60",
			title:       "阅历人生",
			description: "角色达到 60 岁",
			category:    "milestone",
			progress:    character.CurrentAge,
			max:         60,
			unlocked:    character.CurrentAge >= 60,
		},
		{
			id:          "intel_80",
			title:       "智者",
			description: "智力达到 80",
			category:    "growth",
			progress:    character.Attributes.Intelligence,
			max:         80,
			unlocked:    character.Attributes.Intelligence >= 80,
		},
		{
			id:          "eq_80",
			title:       "社交达人",
			description: "情商达到 80",
			category:    "social",
			progress:    character.Attributes.EmotionalIntelligence,
			max:         80,
			unlocked:    character.Attributes.EmotionalIntelligence >= 80,
		},
		{
			id:          "wealth_10k",
			title:       "小康生活",
			description: "资产达到 1 万",
			category:    "wealth",
			progress:    int(character.Money / 100),
			max:         int(tenThousand / 100),
			unlocked:    character.Money >= tenThousand,
		},
		{
			id:          "wealth_100k",
			title:       "富足人生",
			description: "资产达到 10 万",
			category:    "wealth",
			progress:    int(character.Money / 1000),
			max:         int(hundredThousand / 1000),
			unlocked:    character.Money >= hundredThousand,
		},
		{
			id:          "wealth_1m",
			title:       "财富自由",
			description: "资产达到 100 万",
			category:    "wealth",
			progress:    int(character.Money / 10000),
			max:         int(million / 10000),
			unlocked:    character.Money >= million,
		},
		{
			id:          "events_10",
			title:       "故事收藏家",
			description: "累计 10 个关键事件",
			category:    "milestone",
			progress:    len(events),
			max:         10,
			unlocked:    len(events) >= 10,
		},
		{
			id:          "completed",
			title:       "人生落幕",
			description: "完成一次完整人生",
			category:    "milestone",
			progress:    boolToInt(character.GameCompleted),
			max:         1,
			unlocked:    character.GameCompleted,
		},
	}

	items := make([]models.AchievementItem, 0, len(defs))
	unlockedCount := 0
	for _, d := range defs {
		item := models.AchievementItem{
			ID:          d.id,
			Title:       d.title,
			Description: d.description,
			Category:    d.category,
			Unlocked:    d.unlocked,
			Progress:    clampProgress(d.progress, d.max),
			MaxProgress: d.max,
		}
		if d.unlocked {
			unlockedCount++
			unlockedAt := character.UpdatedAt
			item.UnlockedAt = &unlockedAt
		}
		items = append(items, item)
	}

	return &models.CharacterAchievementsResponse{
		CharacterID:   character.CharacterID,
		CharacterName: character.CharacterName,
		Items:         items,
		UnlockedCount: unlockedCount,
		TotalCount:    len(items),
	}, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func (s *achievementStatsService) GetAchievementCategories(ctx context.Context, userID uint) ([]models.AchievementCategory, error) {
	characters, err := s.characterRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user characters failed: %w", err)
	}

	// 目前按固定规则定义分类总数，后续可迁移到配置中心或数据库
	base := []models.AchievementCategory{
		{ID: "milestone", Name: "里程碑", Description: "年龄与人生阶段关键节点", TotalCount: 4},
		{ID: "growth", Name: "成长", Description: "角色属性成长类成就", TotalCount: 1},
		{ID: "social", Name: "社交", Description: "社交与人际能力成就", TotalCount: 1},
		{ID: "wealth", Name: "财富", Description: "财富积累相关成就", TotalCount: 3},
	}

	if len(characters) == 0 {
		return base, nil
	}

	countMap := map[string]int{
		"milestone": 0,
		"growth":    0,
		"social":    0,
		"wealth":    0,
	}

	for _, character := range characters {
		if character.CurrentAge >= 18 {
			countMap["milestone"]++
		}
		if character.CurrentAge >= 60 {
			countMap["milestone"]++
		}
		if character.GameCompleted {
			countMap["milestone"]++
		}
		if character.Attributes.Intelligence >= 80 {
			countMap["growth"]++
		}
		if character.Attributes.EmotionalIntelligence >= 80 {
			countMap["social"]++
		}
		if character.Money >= 10_000 {
			countMap["wealth"]++
		}
		if character.Money >= 100_000 {
			countMap["wealth"]++
		}
		if character.Money >= 1_000_000 {
			countMap["wealth"]++
		}
	}

	for i := range base {
		base[i].UnlockedCount = countMap[base[i].ID]
		if base[i].UnlockedCount > base[i].TotalCount {
			base[i].UnlockedCount = base[i].TotalCount
		}
	}

	return base, nil
}

func (s *achievementStatsService) GetCharacterStats(ctx context.Context, characterID string, userID uint) (*models.CharacterStatsResponse, error) {
	if err := s.ensureOwnership(ctx, characterID, userID); err != nil {
		return nil, err
	}

	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get character failed: %w", err)
	}

	events, err := s.characterRepo.GetGameEventHistory(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get event history failed: %w", err)
	}

	pendingDecision, err := s.characterRepo.GetPendingDecision(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get pending decision failed: %w", err)
	}

	var lastEvent *models.Event
	if len(events) > 0 {
		le := events[len(events)-1]
		lastEvent = &le
	}

	return &models.CharacterStatsResponse{
		CharacterID:     character.CharacterID,
		CharacterName:   character.CharacterName,
		CurrentAge:      character.CurrentAge,
		LifeStage:       character.LifeStage,
		IsGameActive:    !character.GameCompleted,
		TotalPlaytime:   character.TotalPlaytime,
		Money:           character.Money,
		EventCount:      len(events),
		PendingDecision: pendingDecision != nil,
		Attributes:      character.Attributes,
		LastEvent:       lastEvent,
	}, nil
}

func (s *achievementStatsService) GetCharacterTimeline(ctx context.Context, characterID string, userID uint) (*models.CharacterTimelineResponse, error) {
	if err := s.ensureOwnership(ctx, characterID, userID); err != nil {
		return nil, err
	}

	events, err := s.characterRepo.GetGameEventHistory(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("get event history failed: %w", err)
	}

	timeline := make([]models.TimelineItem, 0, len(events))
	for _, e := range events {
		timeline = append(timeline, models.TimelineItem{
			Age:         e.Age,
			Description: e.Description,
			Impact:      e.Impact,
			CreatedAt:   e.CreatedAt,
		})
	}

	return &models.CharacterTimelineResponse{
		CharacterID: characterID,
		Timeline:    timeline,
		Total:       len(timeline),
	}, nil
}

