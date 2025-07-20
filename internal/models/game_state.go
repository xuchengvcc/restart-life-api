package models

// GameState 游戏状态
type GameState struct {
	CharacterID     string              `json:"character_id" db:"character_id"`
	CharacterName   string              `json:"character_name" db:"character_name"`
	CurrentAge      int                 `json:"current_age" db:"current_age"`
	LifeStage       string              `json:"life_stage" db:"life_stage"`
	KeyEvents       []Event             `json:"key_events"`
	PendingDecision *Decision           `json:"pending_decision"`
	Attributes      CharacterAttributes `json:"attributes"`
	LastSaveTime    int64               `json:"last_save_time" db:"last_save_time"`
	TotalPlaytime   int                 `json:"total_playtime" db:"total_playtime"`
	IsGameActive    bool                `json:"is_game_active" db:"is_game_active"`
	CreatedAt       int64               `json:"created_at" db:"created_at"`
	UpdatedAt       int64               `json:"updated_at" db:"updated_at"`

	// 角色详细状态
	Education       string `json:"education" db:"education"`
	Career          string `json:"career" db:"career"`
	Location        string `json:"location" db:"location"`
	MaritalStatus   string `json:"marital_status" db:"marital_status"`
	FamilySituation string `json:"family_situation" db:"family_situation"`
	SocialStatus    string `json:"social_status" db:"social_status"`
	HealthStatus    string `json:"health_status" db:"health_status"`
	WealthLevel     string `json:"wealth_level" db:"wealth_level"`
	Relationships   string `json:"relationships" db:"relationships"`
	PersonalGrowth  string `json:"personal_growth" db:"personal_growth"`

	// 财富状态（对应character表的money）
	Money int64 `json:"money" db:"money"`

	// 上一年的简短描述（对应character表的current_activity）
	LastYearDescription string `json:"last_year_description" db:"last_year_description"`
}

// Event 人生事件
type Event struct {
	EventID     uint64 `json:"event_id" db:"event_id"`
	CharacterID string `json:"character_id" db:"character_id"`
	Age         int    `json:"age" db:"age"`
	Description string `json:"description" db:"description"`
	Impact      string `json:"impact" db:"impact"` // 事件影响简要描述
	CreatedAt   int64  `json:"created_at" db:"created_at"`
}

// Decision 决策选择
type Decision struct {
	CharacterID      string          `json:"character_id" db:"character_id"`
	Options          *DecisionOption `json:"options"`
	PreviousOptions  *DecisionOption `json:"previous_options" db:"previous_options"`   // 上一个决策选项
	PreviousDecision *int8           `json:"previous_decision" db:"previous_decision"` // 上一个决策 1: conservative, 2: moderate, 3: aggressive
	CreatedAt        int64           `json:"created_at" db:"created_at"`
	UpdatedAt        int64           `json:"updated_at" db:"updated_at"`
}

// DecisionOption 决策选项
type DecisionOption struct {
	Conservative DecisionDetails `json:"conservative"` // 保守选项
	Moderate     DecisionDetails `json:"moderate"`     // 中庸选项
	Aggressive   DecisionDetails `json:"aggressive"`   // 激进选项
}

type DecisionDetails struct {
	DecisionType int8   `json:"decision_type"` // 1: conservative, 2: moderate, 3: aggressive
	OptionText   string `json:"option_text"`   // 选项描述
	Consequence  string `json:"consequence"`   // 预期后果
}

// 使用现有的 LifeStageType 定义

// GameRequest 游戏操作请求
type StartGameRequest struct {
	CharacterID string `json:"character_id" binding:"required,uuid"`
}

type AdvanceGameRequest struct {
	CharacterID string `json:"character_id" binding:"required,uuid"`
}

type GameProgressRequest struct {
	OptionType string `json:"option_type,omitempty" binding:"omitempty,oneof=conservative moderate aggressive"`
}

type MakeDecisionRequest struct {
	CharacterID string `json:"character_id" binding:"required,uuid"`
	DecisionID  string `json:"decision_id" binding:"required,uuid"`
	OptionType  string `json:"option_type" binding:"required,oneof=conservative moderate aggressive"`
}

type SaveGameRequest struct {
	CharacterID string `json:"character_id" binding:"required,uuid"`
}

type LoadGameRequest struct {
	CharacterID string `json:"character_id" binding:"required,uuid"`
}
