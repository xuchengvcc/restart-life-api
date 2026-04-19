package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/repository"
)

// AI 事件生成的响应结构
type AIEventResponse struct {
	HasKeyEvent      bool                   `json:"has_key_event"`
	KeyEvent         *models.Event          `json:"key_event,omitempty"`
	Decision         *models.DecisionOption `json:"decision,omitempty"`
	CharacterUpdates map[string]interface{} `json:"character_updates,omitempty"`
}

// AI 决策结果响应结构
type AIDecisionResponse struct {
	ResultEvent      models.Event           `json:"result_event"`
	AttributeChanges map[string]int         `json:"attribute_changes"`
	CharacterUpdates map[string]interface{} `json:"character_updates"`
}

// AI 游戏结束检查响应结构
type AIGameEndResponse struct {
	GameEnded   bool         `json:"game_ended"`
	DeathCause  *string      `json:"death_cause,omitempty"`
	LifeSummary *string      `json:"life_summary,omitempty"`
	FinalEvent  models.Event `json:"final_event,omitempty"`
}

// AI 统一游戏进程响应结构
type AIGameProgressResponse struct {
	// 游戏是否结束
	GameEnded bool `json:"game_ended"`

	// 如果游戏结束的相关信息
	DeathCause  *string       `json:"death_cause,omitempty"`
	LifeSummary *string       `json:"life_summary,omitempty"`
	FinalEvent  *models.Event `json:"final_event,omitempty"`

	// 如果游戏继续的相关信息
	HasKeyEvent bool                   `json:"has_key_event"`
	KeyEvent    *models.Event          `json:"key_event,omitempty"`
	Decision    *models.DecisionOption `json:"decision,omitempty"`

	// 属性变化（每次调用都应该有，可以为0）
	AttributeChanges map[string]int `json:"attribute_changes"`

	// 角色状态更新（每次调用都应该有）
	CharacterUpdates map[string]interface{} `json:"character_updates"`

	// 这一年的简短描述（每次调用都应该有）
	YearDescription string `json:"year_description"`
}

type GameService interface {
	// 开始或继续游戏，优先返回未完成角色的GameState，否则引导新建角色
	StartOrResumeGame(ctx context.Context, userID uint64) (*models.GameState, error)
	// 推进游戏进程，生成新事件和决策选项
	AdvanceGame(ctx context.Context, characterID string) (*models.GameState, error)
	// 智能推进游戏：自动判断是否需要决策，统一处理推进和决策
	AdvanceGameSmart(ctx context.Context, characterID string, optionType string) (*models.GameState, error)
	// 做出决策，处理用户选择，生成新状态和事件
	MakeDecision(ctx context.Context, characterID string, optionType string) (*models.GameState, error)
	// 获取当前游戏状态
	GetGameState(ctx context.Context, characterID string) (*models.GameState, error)
	// 获取关键事件历史
	GetEventHistory(ctx context.Context, characterID string) ([]models.Event, error)
	// 保存/加载游戏（如有断点续玩需求）
	SaveGame(ctx context.Context, characterID string) error
	LoadGame(ctx context.Context, characterID string) (*models.GameState, error)
}

type gameService struct {
	characterRepo repository.CharacterRepository
	aiServices    map[string]AIService
	logger        *logrus.Logger
	redisClient   *redis.Client
}

func NewGameService(
	characterRepo repository.CharacterRepository,
	aiServices map[string]AIService,
	logger *logrus.Logger,
	redisClient *redis.Client,
) GameService {
	return &gameService{
		characterRepo: characterRepo,
		aiServices:    aiServices,
		logger:        logger,
		redisClient:   redisClient,
	}
}

// StartOrResumeGame 开始或继续游戏
func (s *gameService) StartOrResumeGame(ctx context.Context, userID uint64) (*models.GameState, error) {
	// 查找用户的角色
	characters, err := s.characterRepo.GetByUserID(ctx, uint(userID))
	if err != nil {
		return nil, fmt.Errorf("查找角色失败: %w", err)
	}

	// 查找未完成的角色
	for _, character := range characters {
		if !character.GameCompleted {
			return s.GetGameState(ctx, character.CharacterID)
		}
	}

	// 没有未完成角色，返回错误提示需要新建角色
	return nil, fmt.Errorf("需要创建新角色")
}

// AdvanceGame 推进游戏进程，使用统一的AI调用
func (s *gameService) AdvanceGame(ctx context.Context, characterID string) (*models.GameState, error) {
	// 获取当前游戏状态
	gameState, err := s.GetGameState(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("获取游戏状态失败: %w", err)
	}

	// 获取角色信息
	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("获取角色信息失败: %w", err)
	}

	// 年龄增加1
	gameState.CurrentAge++
	gameState.LifeStage = string(models.GetLifeStageByAge(gameState.CurrentAge))
	gameState.UpdatedAt = time.Now().UnixMilli()
	gameState.TotalPlaytime++

	// 使用统一的AI调用处理游戏进程
	aiResponse, err := s.processGameProgressWithAI(ctx, character, gameState, "", false)
	if err != nil {
		s.logger.WithError(err).Warn("AI处理失败，使用默认逻辑")
		// 降级逻辑
		s.fallbackGameProgress(gameState, character)
	} else {
		// 应用AI响应
		s.applyGameProgressResponse(gameState, character, aiResponse)
	}

	s.logger.WithFields(logrus.Fields{
		"character_id": characterID,
		"age":          gameState.CurrentAge,
		"life_stage":   gameState.LifeStage,
		"game_ended":   !gameState.IsGameActive,
	}).Info("游戏推进完成")

	return gameState, nil
}

// AdvanceGameSmart 智能推进游戏（统一处理推进和决策）
func (s *gameService) AdvanceGameSmart(ctx context.Context, characterID string, optionType string) (*models.GameState, error) {
	// 获取当前游戏状态
	gameState, err := s.GetGameState(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("获取游戏状态失败: %w", err)
	}

	// 检查是否有待处理的决策
	if gameState.PendingDecision != nil {
		// 有待处理决策，必须提供option_type
		if optionType == "" {
			return nil, fmt.Errorf("有待处理的决策，请提供option_type参数")
		}
		// 处理决策
		return s.MakeDecision(ctx, characterID, optionType)
	}

	// 没有待处理决策，如果提供了option_type则忽略
	if optionType != "" {
		s.logger.WithFields(logrus.Fields{
			"character_id": characterID,
			"option_type":  optionType,
		}).Warn("没有待处理决策，忽略option_type参数")
	}

	// 推进游戏
	return s.AdvanceGame(ctx, characterID)
}

// MakeDecision 做出决策，使用统一的AI调用
func (s *gameService) MakeDecision(ctx context.Context, characterID string, optionType string) (*models.GameState, error) {
	// 获取当前游戏状态
	gameState, err := s.GetGameState(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("获取游戏状态失败: %w", err)
	}

	// 检查是否有待处理的决策
	if gameState.PendingDecision == nil {
		return nil, fmt.Errorf("没有找到待处理的决策")
	}

	// 获取角色信息
	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("获取角色信息失败: %w", err)
	}

	// 使用统一的AI调用处理决策和后续游戏进程
	aiResponse, err := s.processGameProgressWithAI(ctx, character, gameState, optionType, true)
	// 清除待处理的决策
	gameState.PendingDecision = nil
	if clearErr := s.characterRepo.ClearPendingDecision(ctx, characterID); clearErr != nil {
		s.logger.WithError(clearErr).WithField("character_id", characterID).Warn("清除待处理决策失败")
	}

	if err != nil {
		s.logger.WithError(err).Warn("AI决策处理失败，使用默认逻辑")
		// 降级逻辑
		s.fallbackDecisionResult(gameState, character, optionType)
	} else {
		// 应用AI响应
		s.applyGameProgressResponse(gameState, character, aiResponse)
	}

	gameState.UpdatedAt = time.Now().UnixMilli()

	s.logger.WithFields(logrus.Fields{
		"character_id": characterID,
		"age":          gameState.CurrentAge,
		"option_type":  optionType,
		"game_ended":   !gameState.IsGameActive,
	}).Info("决策处理完成")

	return gameState, nil
}

// GetGameState 获取游戏状态（优先从Redis获取）
func (s *gameService) GetGameState(ctx context.Context, characterID string) (*models.GameState, error) {
	// 1. 尝试从Redis获取
	if gameState := s.getGameStateFromRedis(characterID); gameState != nil {
		return gameState, nil
	}

	// 2. 从Database重建
	gameState, err := s.buildGameStateFromDatabase(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// 3. 写入Redis缓存
	s.saveGameStateToRedis(characterID, gameState)

	return gameState, nil
}

// GetEventHistory 获取事件历史
func (s *gameService) GetEventHistory(ctx context.Context, characterID string) ([]models.Event, error) {
	events, err := s.characterRepo.GetGameEventHistory(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("获取事件历史失败: %w", err)
	}
	return events, nil
}

// SaveGame 保存游戏
func (s *gameService) SaveGame(ctx context.Context, characterID string) error {
	gameState, err := s.GetGameState(ctx, characterID)
	if err != nil {
		return fmt.Errorf("获取游戏状态失败: %w", err)
	}

	gameState.LastSaveTime = time.Now().UnixMilli()

	// 清除Redis缓存，确保下次获取最新状态
	s.clearGameStateFromRedis(characterID)

	// 保存后立即回读进行轻量一致性校验，避免状态漂移被静默吞掉。
	reloadedState, err := s.LoadGame(ctx, characterID)
	if err != nil {
		return fmt.Errorf("save integrity reload failed: %w", err)
	}
	if err := s.verifyGameStateConsistency(gameState, reloadedState); err != nil {
		return fmt.Errorf("save integrity check failed: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"character_id": characterID,
	}).Info("游戏保存完成")

	return nil
}

// LoadGame 加载游戏
func (s *gameService) LoadGame(ctx context.Context, characterID string) (*models.GameState, error) {
	return s.GetGameState(ctx, characterID)
}

// 私有辅助方法

func (s *gameService) verifyGameStateConsistency(expected, actual *models.GameState) error {
	if expected == nil || actual == nil {
		return fmt.Errorf("game state is nil")
	}
	if expected.CharacterID != actual.CharacterID {
		return fmt.Errorf("character_id mismatch: expected=%s actual=%s", expected.CharacterID, actual.CharacterID)
	}
	if expected.CurrentAge != actual.CurrentAge {
		return fmt.Errorf("current_age mismatch: expected=%d actual=%d", expected.CurrentAge, actual.CurrentAge)
	}
	if expected.IsGameActive != actual.IsGameActive {
		return fmt.Errorf("is_game_active mismatch: expected=%t actual=%t", expected.IsGameActive, actual.IsGameActive)
	}
	if (expected.PendingDecision == nil) != (actual.PendingDecision == nil) {
		return fmt.Errorf("pending_decision presence mismatch")
	}
	return nil
}

func (s *gameService) getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (s *gameService) getWealthLevel(money int64) string {
	switch {
	case money < 0:
		return "负债"
	case money < 10000:
		return "贫困"
	case money < 100000:
		return "小康"
	case money < 1000000:
		return "富裕"
	default:
		return "富豪"
	}
}

func (s *gameService) getRelationships(age int) string {
	switch {
	case age < 18:
		return "父母、家人"
	case age < 30:
		return "朋友、同学"
	case age < 60:
		return "同事、配偶、子女"
	default:
		return "子女、孙辈"
	}
}

func (s *gameService) getPersonalGrowth(age int) string {
	switch {
	case age < 6:
		return "快乐成长，学习基本技能"
	case age < 18:
		return "接受教育，形成价值观"
	case age < 30:
		return "探索人生，建立事业"
	case age < 60:
		return "成熟稳重，承担责任"
	default:
		return "享受人生，传承智慧"
	}
}

func (s *gameService) clampAttributes(attrs *models.CharacterAttributes) {
	attrs.Intelligence = s.clampValue(attrs.Intelligence)
	attrs.EmotionalIntelligence = s.clampValue(attrs.EmotionalIntelligence)
	attrs.Memory = s.clampValue(attrs.Memory)
	attrs.Imagination = s.clampValue(attrs.Imagination)
	attrs.PhysicalFitness = s.clampValue(attrs.PhysicalFitness)
	attrs.Appearance = s.clampValue(attrs.Appearance)
}

func (s *gameService) clampValue(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

// AI 集成方法

// processGameProgressWithAI 统一的AI游戏进程处理
func (s *gameService) processGameProgressWithAI(ctx context.Context, character *models.Character, gameState *models.GameState, optionType string, isDecision bool) (*AIGameProgressResponse, error) {
	// 构建统一的AI提示词
	prompt := s.buildUnifiedGamePrompt(character, gameState, optionType, isDecision)

	// 调用AI服务
	aiText, err := s.callAIService(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI调用失败: %w", err)
	}

	// 解析AI响应
	response, err := s.parseGameProgressResponse(aiText)
	if err != nil {
		return nil, fmt.Errorf("AI响应解析失败: %w, AI response: %s", err, aiText)
	}

	return response, nil
}

// applyGameProgressResponse 应用AI游戏进程响应
func (s *gameService) applyGameProgressResponse(gameState *models.GameState, character *models.Character, response *AIGameProgressResponse) {
	// 如果游戏结束
	if response.GameEnded {
		gameState.IsGameActive = false
		character.GameCompleted = true
		finalAge := gameState.CurrentAge
		character.FinalAge = &finalAge
		character.DeathCause = response.DeathCause
		character.Summary = response.LifeSummary

		// 添加结局事件
		if response.FinalEvent != nil {
			endEvent := *response.FinalEvent
			endEvent.EventID = 0
			endEvent.CharacterID = character.CharacterID
			endEvent.Age = gameState.CurrentAge
			endEvent.CreatedAt = time.Now().UnixMilli()
			gameState.KeyEvents = append(gameState.KeyEvents, endEvent)
			if err := s.characterRepo.SaveGameEvent(context.Background(), &endEvent); err != nil {
				s.logger.WithError(err).Error("保存结局事件失败")
			}
		}

		if err := s.characterRepo.ClearPendingDecision(context.Background(), character.CharacterID); err != nil {
			s.logger.WithError(err).Warn("清除待处理决策失败")
		}

		// 更新角色到数据库
		err := s.characterRepo.Update(context.Background(), character)
		if err != nil {
			s.logger.WithError(err).Error("更新角色游戏完成状态失败")
		}
		return
	}

	// 应用属性变化（每次调用都应该有）
	if response.AttributeChanges != nil {
		if val, ok := response.AttributeChanges["intelligence"]; ok {
			gameState.Attributes.Intelligence += val
		}
		if val, ok := response.AttributeChanges["emotional_intelligence"]; ok {
			gameState.Attributes.EmotionalIntelligence += val
		}
		if val, ok := response.AttributeChanges["memory"]; ok {
			gameState.Attributes.Memory += val
		}
		if val, ok := response.AttributeChanges["imagination"]; ok {
			gameState.Attributes.Imagination += val
		}
		if val, ok := response.AttributeChanges["physical_fitness"]; ok {
			gameState.Attributes.PhysicalFitness += val
		}
		if val, ok := response.AttributeChanges["appearance"]; ok {
			gameState.Attributes.Appearance += val
		}
		// 确保属性在合理范围内
		s.clampAttributes(&gameState.Attributes)
	}

	// 应用角色状态变化（每次调用都应该有）
	if response.CharacterUpdates != nil {
		// 更新财富状态
		if moneyChange, ok := response.CharacterUpdates["money"]; ok {
			if moneyVal, ok := moneyChange.(float64); ok {
				gameState.Money += int64(moneyVal)
				character.Money = gameState.Money
			} else if moneyVal, ok := moneyChange.(int64); ok {
				gameState.Money += moneyVal
				character.Money = gameState.Money
			}
		}

		// 更新其他状态字段
		if healthLevel, ok := response.CharacterUpdates["health_level"]; ok {
			if healthVal, ok := healthLevel.(float64); ok {
				character.HealthLevel = int(healthVal)
			}
		}

		if happinessLevel, ok := response.CharacterUpdates["happiness_level"]; ok {
			if happinessVal, ok := happinessLevel.(float64); ok {
				character.HappinessLevel = int(happinessVal)
			}
		}

		// 记录状态更新日志
		s.logger.WithFields(logrus.Fields{
			"character_id": character.CharacterID,
			"updates":      response.CharacterUpdates,
		}).Debug("应用角色状态更新")
	}

	// 更新年度描述
	if response.YearDescription != "" {
		gameState.LastYearDescription = response.YearDescription
		// 同时更新到character的current_activity字段
		description := response.YearDescription
		character.CurrentActivity = &description
	}

	// 添加关键事件（包括决策结果事件）
	if response.HasKeyEvent && response.KeyEvent != nil {
		event := *response.KeyEvent
		event.EventID = 0
		event.CharacterID = character.CharacterID
		event.Age = gameState.CurrentAge
		event.CreatedAt = time.Now().UnixMilli()
		gameState.KeyEvents = append(gameState.KeyEvents, event)
		if err := s.characterRepo.SaveGameEvent(context.Background(), &event); err != nil {
			s.logger.WithError(err).Error("保存关键事件失败")
		}
	}

	// 设置新的决策选项（如果有）
	if response.Decision != nil {
		decision := &models.Decision{
			CharacterID: character.CharacterID,
			Options:     response.Decision,
			CreatedAt:   time.Now().UnixMilli(),
			UpdatedAt:   time.Now().UnixMilli(),
		}
		gameState.PendingDecision = decision
		if err := s.characterRepo.SavePendingDecision(context.Background(), decision); err != nil {
			s.logger.WithError(err).Error("保存待处理决策失败")
		}
	} else {
		gameState.PendingDecision = nil
		if err := s.characterRepo.ClearPendingDecision(context.Background(), character.CharacterID); err != nil {
			s.logger.WithError(err).Warn("清除待处理决策失败")
		}
	}

	// 更新角色到数据库
	err := s.characterRepo.Update(context.Background(), character)
	if err != nil {
		s.logger.WithError(err).Error("更新角色信息失败")
	}

	// 清除Redis缓存（下次GetGameState时会重新构建）
	s.clearGameStateFromRedis(character.CharacterID)
}

// callAIService 调用AI服务
func (s *gameService) callAIService(ctx context.Context, prompt string) (string, error) {
	// 优先使用 Hunyuan 服务
	if hunyuan, exists := s.aiServices["hunyuan"]; exists {
		return hunyuan.GenerateText(ctx, prompt)
	}

	// 备选方案：使用 Gemini 服务
	if gemini, exists := s.aiServices["gemini"]; exists {
		return gemini.GenerateText(ctx, prompt)
	}

	return "", fmt.Errorf("没有可用的AI服务")
}

// buildUnifiedGamePrompt 构建统一的游戏AI提示词
func (s *gameService) buildUnifiedGamePrompt(character *models.Character, gameState *models.GameState, optionType string, isDecision bool) string {
	baseInfo := fmt.Sprintf(`你是一个人生模拟游戏的智能引擎。请根据角色信息，进行以下判断和生成：

角色基本信息：
- 姓名：%s
- 当前年龄：%d
- 人生阶段：%s
- 出生国家：%s
- 出生年份：%d

当前属性：
- 智力：%d
- 情商：%d
- 记忆力：%d
- 想象力：%d
- 体质：%d
- 外貌：%d

当前状态：
- 教育背景：%s
- 职业情况：%s
- 居住地点：%s
- 婚姻状况：%s
- 家庭情况：%s
- 社会地位：%s
- 健康状况：%s
- 财富水平：%s
- 人际关系：%s
- 个人成长：%s`,
		character.CharacterName,
		gameState.CurrentAge,
		gameState.LifeStage,
		character.BirthCountry,
		character.BirthYear,
		gameState.Attributes.Intelligence,
		gameState.Attributes.EmotionalIntelligence,
		gameState.Attributes.Memory,
		gameState.Attributes.Imagination,
		gameState.Attributes.PhysicalFitness,
		gameState.Attributes.Appearance,
		gameState.Education,
		gameState.Career,
		gameState.Location,
		gameState.MaritalStatus,
		gameState.FamilySituation,
		gameState.SocialStatus,
		gameState.HealthStatus,
		gameState.WealthLevel,
		gameState.Relationships,
		gameState.PersonalGrowth,
	)

	// 添加历史事件信息
	historyInfo := s.buildHistoryInfo(gameState)

	// 添加当前决策信息（如果有）
	currentDecisionInfo := s.buildCurrentDecisionInfo(gameState, optionType, isDecision)

	var taskDescription string
	var responseFormat string

	if isDecision {
		// 决策处理模式
		var optionDesc string
		switch optionType {
		case "conservative":
			optionDesc = "保守选择"
		case "moderate":
			optionDesc = "中庸选择"
		case "aggressive":
			optionDesc = "激进选择"
		}

		taskDescription = fmt.Sprintf(`

任务：角色做出了"%s"，请：
1. 首先判断游戏是否应该结束（考虑年龄、体质、意外等因素）
2. 如果游戏结束，生成死亡原因、人生总结和最终事件
3. 如果游戏继续，必须生成：
   - 决策结果的属性变化（可以为0，但必须提供所有属性的变化值）
   - 角色状态的变化（教育、职业、健康等状态的更新）
   - 这一年的简短描述（角色主要活动和状态）
   - 谨慎决定是否生成关键事件
4. 必须为下一个人生阶段提供新的三个决策选项（保守、中庸、激进）`, optionDesc)

		responseFormat = `

请以JSON格式返回：
{
  "game_ended": true/false,
  "death_cause": "死亡原因（如果游戏结束）",
  "life_summary": "人生总结（如果游戏结束）",
  "final_event": {
    "description": "最终事件描述",
    "impact": "最终影响"
  },
  "has_key_event": true/false,
  "key_event": {
    "description": "决策结果或新事件描述",
    "impact": "事件影响描述"
  },
  "attribute_changes": {
    "intelligence": 数值变化,
    "emotional_intelligence": 数值变化,
    "memory": 数值变化,
    "imagination": 数值变化,
    "physical_fitness": 数值变化,
    "appearance": 数值变化
  },
  "character_updates": {
    "education_desc": "教育经历更新（如有变化）",
    "career_desc": "职业履历更新（如有变化）",
    "current_status": "当前状态（如healthy, sick等）",
    "happiness_level": 数值,
    "health_level": 数值,
    "money": 数值变化
  },
  "year_description": "这一年的简短描述，概括角色的主要活动和状态",
  "decision": {
    "conservative": {
      "decision_type": 1,
      "option_text": "具体的保守选项描述",
      "consequence": "预期后果描述"
    },
    "moderate": {
      "decision_type": 2,
      "option_text": "具体的中庸选项描述",
      "consequence": "预期后果描述"
    },
    "aggressive": {
      "decision_type": 3,
      "option_text": "具体的激进选项描述",
      "consequence": "预期后果描述"
    }
  }
}`
	} else {
		// 游戏推进模式
		taskDescription = `

任务：角色年龄刚刚增加了1岁，请：
1. 首先判断游戏是否应该结束（年龄、体质、意外等因素）
2. 如果游戏结束，生成死亡原因、人生总结和最终事件
3. 如果游戏继续，必须生成：
   - 属性的自然变化（可以为0，但必须提供所有属性的变化值）
   - 角色状态的变化（教育、职业、健康等状态的更新）
   - 这一年的简短描述（角色主要活动和状态）
   - 谨慎决定是否生成关键事件（大多数年龄推进无需事件）
4. 如果生成关键事件，必须同时提供相应的三个决策选项`

		responseFormat = `

请以JSON格式返回：
{
  "game_ended": true/false,
  "death_cause": "死亡原因（如果游戏结束）",
  "life_summary": "人生总结（如果游戏结束）",
  "final_event": {
    "description": "最终事件描述",
    "impact": "最终影响"
  },
  "has_key_event": true/false,
  "key_event": {
    "description": "关键事件描述",
    "impact": "事件影响"
  },
  "attribute_changes": {
    "intelligence": 数值变化,
    "emotional_intelligence": 数值变化,
    "memory": 数值变化,
    "imagination": 数值变化,
    "physical_fitness": 数值变化,
    "appearance": 数值变化
  },
  "character_updates": {
    "education_desc": "教育经历更新（如有变化）",
    "career_desc": "职业履历更新（如有变化）",
    "current_status": "当前状态（如healthy, sick等）",
    "happiness_level": 数值,
    "health_level": 数值,
    "money": 数值变化
  },
  "year_description": "这一年的简短描述，概括角色的主要活动和状态",
  "decision": {
    "conservative": {
      "decision_type": 1,
      "option_text": "具体的保守选项描述",
      "consequence": "预期后果描述"
    },
    "moderate": {
      "decision_type": 2,
      "option_text": "具体的中庸选项描述",
      "consequence": "预期后果描述"
    },
    "aggressive": {
      "decision_type": 3,
      "option_text": "具体的激进选项描述",
      "consequence": "预期后果描述"
    }
  }
}`
	}

	rules := `

规则：
1. 85岁后死亡概率增加，但不是绝对
2. 体质<20时健康风险增加
3. 关键事件(key_event)只在真正对人生有重大影响时才生成，不是每年都有
4. 关键事件包括：重要人生节点、重大决策结果、意外事件等
5. 必须每次都提供attribute_changes（即使变化为0也要明确给出）
6. 必须每次都提供character_updates（角色状态的变化）
7. 必须每次都提供year_description（这一年的简短描述）
8. 必须根据角色当前状态和历史背景生成具体、有意义的决策选项
9. 三个决策选项要有明确区别，保守注重稳定，中庸寻求平衡，激进追求突破
10. 属性变化要合理，可以为负值
11. 保持适度的随机性和现实性
12. 大多数年龄推进不需要生成事件，只更新年龄即可
13. 决策选项必须与当前角色状态和人生阶段相关`

	return baseInfo + historyInfo + currentDecisionInfo + taskDescription + rules + responseFormat
}

// buildHistoryInfo 构建历史事件信息
func (s *gameService) buildHistoryInfo(gameState *models.GameState) string {
	if len(gameState.KeyEvents) == 0 {
		return "\n\n历史关键事件：无"
	}

	historyInfo := "\n\n历史关键事件："
	for i, event := range gameState.KeyEvents {
		historyInfo += fmt.Sprintf("\n%d. [%d岁] %s - %s", i+1, event.Age, event.Description, event.Impact)
	}
	return historyInfo
}

// buildCurrentDecisionInfo 构建当前决策信息
func (s *gameService) buildCurrentDecisionInfo(gameState *models.GameState, optionType string, isDecision bool) string {
	if !isDecision || gameState.PendingDecision == nil {
		return ""
	}

	info := "\n\n当前决策背景："
	if gameState.PendingDecision.Options != nil {
		info += "\n可选项："
		info += fmt.Sprintf("\n- 保守选择: %s (%s)",
			gameState.PendingDecision.Options.Conservative.OptionText,
			gameState.PendingDecision.Options.Conservative.Consequence)
		info += fmt.Sprintf("\n- 中庸选择: %s (%s)",
			gameState.PendingDecision.Options.Moderate.OptionText,
			gameState.PendingDecision.Options.Moderate.Consequence)
		info += fmt.Sprintf("\n- 激进选择: %s (%s)",
			gameState.PendingDecision.Options.Aggressive.OptionText,
			gameState.PendingDecision.Options.Aggressive.Consequence)
	}

	var selectedOption string
	switch optionType {
	case "conservative":
		selectedOption = "保守选择"
	case "moderate":
		selectedOption = "中庸选择"
	case "aggressive":
		selectedOption = "激进选择"
	}

	if selectedOption != "" {
		info += fmt.Sprintf("\n角色选择了：%s", selectedOption)
	}

	return info
}

// parseGameProgressResponse 解析统一的AI游戏进程响应
func (s *gameService) parseGameProgressResponse(aiText string) (*AIGameProgressResponse, error) {
	clean := sanitizeAIJSON(aiText)
	var response AIGameProgressResponse
	err := json.Unmarshal([]byte(clean), &response)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	return &response, nil
}

// sanitizeAIJSON 去掉模型返回中的 Markdown 代码块包裹，并截取首尾 JSON 范围
func sanitizeAIJSON(aiText string) string {
	text := strings.TrimSpace(aiText)
	if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimPrefix(text, "```JSON")
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSpace(text)
		if idx := strings.LastIndex(text, "```"); idx >= 0 {
			text = strings.TrimSpace(text[:idx])
		}
	}

	// 再次保险：截取首个 '{' 到最后一个 '}' 的部分
	if start := strings.Index(text, "{"); start >= 0 {
		if end := strings.LastIndex(text, "}"); end >= start {
			text = text[start : end+1]
		}
	}

	return text
}

// 降级处理方法

// fallbackGameProgress 游戏推进的降级处理
func (s *gameService) fallbackGameProgress(gameState *models.GameState, character *models.Character) {
	// 简单的年龄判断
	if gameState.CurrentAge >= 85 {
		gameState.IsGameActive = false
		character.GameCompleted = true
		finalAge := gameState.CurrentAge
		character.FinalAge = &finalAge
		deathCause := "自然死亡"
		character.DeathCause = &deathCause
		summary := fmt.Sprintf("%s享年%d岁，完成了精彩的人生。", character.CharacterName, finalAge)
		character.Summary = &summary

		err := s.characterRepo.Update(context.Background(), character)
		if err != nil {
			s.logger.WithError(err).Error("更新角色游戏完成状态失败")
		}
		return
	}

	// 只在关键年龄节点生成模板事件
	if s.shouldGenerateKeyEvent(gameState.CurrentAge) {
		event := s.generateTemplateEvent(gameState, character)
		gameState.KeyEvents = append(gameState.KeyEvents, event)
		if err := s.characterRepo.SaveGameEvent(context.Background(), &event); err != nil {
			s.logger.WithError(err).Error("保存模板事件失败")
		}

		decision := s.generateTemplateDecision(gameState, character)
		gameState.PendingDecision = decision
		if err := s.characterRepo.SavePendingDecision(context.Background(), decision); err != nil {
			s.logger.WithError(err).Error("保存模板待处理决策失败")
		}
	} else {
		gameState.PendingDecision = nil
		if err := s.characterRepo.ClearPendingDecision(context.Background(), character.CharacterID); err != nil {
			s.logger.WithError(err).Warn("清除待处理决策失败")
		}
	}
	// 大多数年龄推进不生成任何事件，只更新年龄

	// 更新角色到数据库
	err := s.characterRepo.Update(context.Background(), character)
	if err != nil {
		s.logger.WithError(err).Error("更新角色信息失败")
	}

	// 清除Redis缓存（下次GetGameState时会重新构建）
	s.clearGameStateFromRedis(character.CharacterID)
}

// fallbackDecisionResult 决策结果的降级处理
func (s *gameService) fallbackDecisionResult(gameState *models.GameState, character *models.Character, optionType string) {
	if err := s.characterRepo.ClearPendingDecision(context.Background(), character.CharacterID); err != nil {
		s.logger.WithError(err).Warn("清除待处理决策失败")
	}

	// 先检查游戏是否结束
	if gameState.CurrentAge >= 85 {
		gameState.IsGameActive = false
		character.GameCompleted = true
		finalAge := gameState.CurrentAge
		character.FinalAge = &finalAge
		deathCause := "自然死亡"
		character.DeathCause = &deathCause
		summary := fmt.Sprintf("%s享年%d岁，完成了精彩的人生。", character.CharacterName, finalAge)
		character.Summary = &summary

		err := s.characterRepo.Update(context.Background(), character)
		if err != nil {
			s.logger.WithError(err).Error("更新角色游戏完成状态失败")
		}
		return
	}

	// 应用模板决策结果
	s.applyTemplateDecisionResult(gameState, character, optionType)

	// 更新角色到数据库
	err := s.characterRepo.Update(context.Background(), character)
	if err != nil {
		s.logger.WithError(err).Error("更新角色信息失败")
	}

	// 清除Redis缓存（下次GetGameState时会重新构建）
	s.clearGameStateFromRedis(character.CharacterID)
}

func (s *gameService) applyAIResponse(gameState *models.GameState, character *models.Character, response *AIEventResponse) {
	// 如果有关键事件，添加到游戏状态
	if response.HasKeyEvent && response.KeyEvent != nil {
		event := *response.KeyEvent
		event.EventID = 0 // 数据库自增
		event.CharacterID = character.CharacterID
		event.Age = gameState.CurrentAge
		event.CreatedAt = time.Now().UnixMilli()
		gameState.KeyEvents = append(gameState.KeyEvents, event)
	}

	// 设置决策选项
	if response.Decision != nil {
		decision := &models.Decision{
			CharacterID: character.CharacterID,
			Options:     response.Decision,
			CreatedAt:   time.Now().UnixMilli(),
			UpdatedAt:   time.Now().UnixMilli(),
		}
		gameState.PendingDecision = decision
	}
}

// applyAIDecisionResult 应用AI生成的决策结果
func (s *gameService) applyAIDecisionResult(gameState *models.GameState, character *models.Character, response *AIDecisionResponse) {
	// 添加结果事件
	event := response.ResultEvent
	event.EventID = 0
	event.CharacterID = character.CharacterID
	event.Age = gameState.CurrentAge
	event.CreatedAt = time.Now().UnixMilli()
	gameState.KeyEvents = append(gameState.KeyEvents, event)

	// 应用属性变化
	if response.AttributeChanges != nil {
		if val, ok := response.AttributeChanges["intelligence"]; ok {
			gameState.Attributes.Intelligence += val
		}
		if val, ok := response.AttributeChanges["emotional_intelligence"]; ok {
			gameState.Attributes.EmotionalIntelligence += val
		}
		if val, ok := response.AttributeChanges["memory"]; ok {
			gameState.Attributes.Memory += val
		}
		if val, ok := response.AttributeChanges["imagination"]; ok {
			gameState.Attributes.Imagination += val
		}
		if val, ok := response.AttributeChanges["physical_fitness"]; ok {
			gameState.Attributes.PhysicalFitness += val
		}
		if val, ok := response.AttributeChanges["appearance"]; ok {
			gameState.Attributes.Appearance += val
		}
	}

	// 确保属性在合理范围内
	s.clampAttributes(&gameState.Attributes)
}

// 模板降级方法

// generateTemplateEventAndDecision 模板事件和决策生成（AI失败时的降级方案）
func (s *gameService) generateTemplateEventAndDecision(gameState *models.GameState, character *models.Character) {
	age := gameState.CurrentAge

	// 只在特定年龄生成关键事件
	shouldGenerateEvent := s.shouldGenerateKeyEvent(age)

	if shouldGenerateEvent {
		// 生成模板事件
		event := s.generateTemplateEvent(gameState, character)
		gameState.KeyEvents = append(gameState.KeyEvents, event)

		// 生成模板决策
		decision := s.generateTemplateDecision(gameState, character)
		gameState.PendingDecision = decision
	}
}

// shouldGenerateKeyEvent 判断是否应该生成关键事件
func (s *gameService) shouldGenerateKeyEvent(age int) bool {
	// 只在真正的关键年龄节点生成事件
	keyAges := []int{6, 12, 18, 22, 30, 40, 50, 60, 70}
	for _, keyAge := range keyAges {
		if age == keyAge {
			return true
		}
	}

	// 其他年龄不生成关键事件，减少存储压力
	return false
}

// generateTemplateEvent 生成模板事件（只在关键节点）
func (s *gameService) generateTemplateEvent(gameState *models.GameState, character *models.Character) models.Event {
	age := gameState.CurrentAge
	var description string

	switch age {
	case 6:
		description = fmt.Sprintf("%s开始上小学，人生的求学之路正式开始。", character.CharacterName)
	case 12:
		description = fmt.Sprintf("%s升入中学，开始面临更多的学业挑战。", character.CharacterName)
	case 18:
		description = fmt.Sprintf("%s成年了，面临人生的重要选择：继续求学还是步入社会？", character.CharacterName)
	case 22:
		description = fmt.Sprintf("%s大学毕业，准备踏入职场开始事业。", character.CharacterName)
	case 30:
		description = fmt.Sprintf("%s步入三十岁，开始思考事业和家庭的平衡。", character.CharacterName)
	case 40:
		description = fmt.Sprintf("%s到了不惑之年，在事业上面临新的机遇和挑战。", character.CharacterName)
	case 50:
		description = fmt.Sprintf("%s进入知天命的年纪，开始重新审视人生的意义。", character.CharacterName)
	case 60:
		description = fmt.Sprintf("%s临近退休，开始为人生的下半场做准备。", character.CharacterName)
	case 70:
		description = fmt.Sprintf("%s步入古稀之年，享受着人生的智慧和宁静。", character.CharacterName)
	default:
		description = fmt.Sprintf("%s在%d岁时迎来了人生的重要时刻。", character.CharacterName, age)
	}

	return models.Event{
		EventID:     0, // 数据库自增
		CharacterID: character.CharacterID,
		Age:         age,
		Description: description,
		Impact:      "人生重要节点",
		CreatedAt:   time.Now().UnixMilli(),
	}
}

// generateTemplateDecision 生成模板决策（只在关键节点）
func (s *gameService) generateTemplateDecision(gameState *models.GameState, character *models.Character) *models.Decision {
	age := gameState.CurrentAge

	// 只在关键年龄节点生成决策
	var conservative, moderate, aggressive models.DecisionDetails

	switch age {
	case 18:
		conservative = models.DecisionDetails{
			DecisionType: 1,
			OptionText:   "选择稳定的专业，按部就班读大学",
			Consequence:  "稳定发展，风险较小",
		}
		moderate = models.DecisionDetails{
			DecisionType: 2,
			OptionText:   "选择自己感兴趣的专业",
			Consequence:  "可能有更好的发展前景",
		}
		aggressive = models.DecisionDetails{
			DecisionType: 3,
			OptionText:   "放弃大学，直接创业或工作",
			Consequence:  "高风险高回报的选择",
		}
	case 22:
		conservative = models.DecisionDetails{
			DecisionType: 1,
			OptionText:   "找一份稳定的工作",
			Consequence:  "收入稳定，生活安逸",
		}
		moderate = models.DecisionDetails{
			DecisionType: 2,
			OptionText:   "选择有发展潜力的公司",
			Consequence:  "可能获得更好的职业发展",
		}
		aggressive = models.DecisionDetails{
			DecisionType: 3,
			OptionText:   "自主创业",
			Consequence:  "可能获得巨大成功或失败",
		}
	default:
		conservative = models.DecisionDetails{
			DecisionType: 1,
			OptionText:   "稳扎稳打，按部就班",
			Consequence:  "稳定发展，风险较小",
		}
		moderate = models.DecisionDetails{
			DecisionType: 2,
			OptionText:   "适度冒险，寻求突破",
			Consequence:  "有机会获得更好发展",
		}
		aggressive = models.DecisionDetails{
			DecisionType: 3,
			OptionText:   "大胆尝试，勇敢追梦",
			Consequence:  "可能获得巨大成功或失败",
		}
	}

	decision := &models.Decision{
		CharacterID: character.CharacterID,
		Options: &models.DecisionOption{
			Conservative: conservative,
			Moderate:     moderate,
			Aggressive:   aggressive,
		},
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	return decision
}

// applyTemplateDecisionResult 应用模板决策结果
func (s *gameService) applyTemplateDecisionResult(gameState *models.GameState, character *models.Character, optionType string) {
	// 根据选择类型应用不同的效果
	switch optionType {
	case "conservative":
		// 保守选择：稳定发展
		gameState.Attributes.Intelligence += 1
	case "moderate":
		// 中庸选择：平衡发展
		gameState.Attributes.Intelligence += 2
		gameState.Attributes.EmotionalIntelligence += 1
	case "aggressive":
		// 激进选择：高风险高回报
		gameState.Attributes.Intelligence += 3
		gameState.Attributes.EmotionalIntelligence += 2
		gameState.Attributes.PhysicalFitness += 1
	}

	// 生成决策后的事件
	resultEvent := s.generateTemplateDecisionResultEvent(gameState, character, optionType)
	gameState.KeyEvents = append(gameState.KeyEvents, resultEvent)
	if err := s.characterRepo.SaveGameEvent(context.Background(), &resultEvent); err != nil {
		s.logger.WithError(err).Error("保存模板决策结果事件失败")
	}

	// 确保属性在合理范围内
	s.clampAttributes(&gameState.Attributes)
}

// generateTemplateDecisionResultEvent 生成模板决策结果事件
func (s *gameService) generateTemplateDecisionResultEvent(gameState *models.GameState, character *models.Character, optionType string) models.Event {
	var description string

	switch optionType {
	case "conservative":
		description = fmt.Sprintf("%s选择了稳健的道路，虽然进展缓慢但很稳定。", character.CharacterName)
	case "moderate":
		description = fmt.Sprintf("%s选择了适中的方案，取得了不错的进展。", character.CharacterName)
	case "aggressive":
		description = fmt.Sprintf("%s选择了冒险的道路，虽然有风险但也有不错的收获。", character.CharacterName)
	default:
		description = fmt.Sprintf("%s做出了选择，人生继续前进。", character.CharacterName)
	}

	return models.Event{
		EventID:     0, // 数据库自增
		CharacterID: character.CharacterID,
		Age:         gameState.CurrentAge,
		Description: description,
		Impact:      fmt.Sprintf("选择了%s策略", optionType),
		CreatedAt:   time.Now().UnixMilli(),
	}
}

// Redis 缓存相关方法

const (
	gameStateKeyPrefix = "game_state:"
	gameStateTTL       = 30 * time.Minute // 30分钟
)

// getGameStateFromRedis 从Redis获取GameState
func (s *gameService) getGameStateFromRedis(characterID string) *models.GameState {
	key := gameStateKeyPrefix + characterID
	data, err := s.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		if err != redis.Nil {
			s.logger.WithError(err).WithField("character_id", characterID).Warn("从Redis获取GameState失败")
		}
		return nil
	}

	var gameState models.GameState
	if err := json.Unmarshal([]byte(data), &gameState); err != nil {
		s.logger.WithError(err).WithField("character_id", characterID).Error("解析Redis中的GameState失败")
		// 清除损坏的缓存
		s.redisClient.Del(context.Background(), key)
		return nil
	}

	s.logger.WithField("character_id", characterID).Debug("从Redis成功获取GameState")
	return &gameState
}

// saveGameStateToRedis 保存GameState到Redis
func (s *gameService) saveGameStateToRedis(characterID string, gameState *models.GameState) {
	key := gameStateKeyPrefix + characterID
	data, err := json.Marshal(gameState)
	if err != nil {
		s.logger.WithError(err).WithField("character_id", characterID).Error("序列化GameState失败")
		return
	}

	err = s.redisClient.Set(context.Background(), key, data, gameStateTTL).Err()
	if err != nil {
		s.logger.WithError(err).WithField("character_id", characterID).Error("保存GameState到Redis失败")
		return
	}

	// 延长TTL，保持用户活跃状态
	s.redisClient.Expire(context.Background(), key, gameStateTTL)

	s.logger.WithField("character_id", characterID).Debug("成功保存GameState到Redis")
}

// clearGameStateFromRedis 清除Redis中的GameState缓存
func (s *gameService) clearGameStateFromRedis(characterID string) {
	key := gameStateKeyPrefix + characterID
	err := s.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		s.logger.WithError(err).WithField("character_id", characterID).Warn("清除Redis缓存失败")
		return
	}

	s.logger.WithField("character_id", characterID).Debug("成功清除Redis缓存")
}

// buildGameStateFromDatabase 从数据库重建GameState
func (s *gameService) buildGameStateFromDatabase(ctx context.Context, characterID string) (*models.GameState, error) {
	// 从数据库获取角色信息
	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("获取角色信息失败: %w", err)
	}

	// 构建游戏状态
	gameState := &models.GameState{
		CharacterID:   characterID,
		CharacterName: character.CharacterName,
		CurrentAge:    character.CurrentAge,
		LifeStage:     character.LifeStage,
		Attributes:    character.Attributes,
		IsGameActive:  !character.GameCompleted,
		LastSaveTime:  character.UpdatedAt,
		TotalPlaytime: character.TotalPlaytime,
		CreatedAt:     character.CreatedAt,
		UpdatedAt:     character.UpdatedAt,

		// 角色详细状态
		Education:       s.getStringValue(character.EducationDesc),
		Career:          s.getStringValue(character.CareerDesc),
		Location:        s.getStringValue(character.CurrentCountry),
		MaritalStatus:   s.getStringValue(character.MaritalStatus),
		FamilySituation: s.getStringValue(character.FamilyBackground),
		SocialStatus:    s.getStringValue(character.SocialRelationships),
		HealthStatus:    fmt.Sprintf("健康等级: %d", character.HealthLevel),
		WealthLevel:     s.getWealthLevel(character.Money),
		Relationships:   s.getRelationships(character.CurrentAge),
		PersonalGrowth:  s.getPersonalGrowth(character.CurrentAge),

		// 财富状态
		Money: character.Money,

		// 从character表的current_activity字段加载年度描述
		LastYearDescription: s.getStringValue(character.CurrentActivity),
	}

	events, err := s.characterRepo.GetGameEventHistory(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("加载事件历史失败: %w", err)
	}
	gameState.KeyEvents = events

	pendingDecision, err := s.characterRepo.GetPendingDecision(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("加载待处理决策失败: %w", err)
	}
	gameState.PendingDecision = pendingDecision

	s.logger.WithField("character_id", characterID).Debug("从数据库重建GameState成功")
	return gameState, nil
}
