package models

// AIGamePrompt AI游戏提示词结构
type AIGamePrompt struct {
	SystemPrompt string        `json:"system_prompt"`
	Character    AICharacter   `json:"character"`
	CurrentState AIGameContext `json:"current_state"`
	RequestType  string        `json:"request_type"` // "advance_age", "make_decision", "generate_event"
}

// AICharacter AI角色描述
type AICharacter struct {
	Name          string              `json:"name"`
	Age           int                 `json:"age"`
	Gender        string              `json:"gender"`
	Race          string              `json:"race"`
	BirthCountry  string              `json:"birth_country"`
	BirthYear     int                 `json:"birth_year"`
	LifeStage     string              `json:"life_stage"`
	Attributes    CharacterAttributes `json:"attributes"`
	CurrentStatus AICharacterStatus   `json:"current_status"`
}

// AICharacterStatus AI角色当前状态
type AICharacterStatus struct {
	Education       string `json:"education"`
	Career          string `json:"career"`
	Location        string `json:"location"`
	MaritalStatus   string `json:"marital_status"`
	FamilySituation string `json:"family_situation"`
	SocialStatus    string `json:"social_status"`
	HealthStatus    string `json:"health_status"`
	WealthLevel     string `json:"wealth_level"`
	Relationships   string `json:"relationships"`
	PersonalGrowth  string `json:"personal_growth"`
}

// AIGameContext AI游戏上下文
type AIGameContext struct {
	KeyEvents       []string `json:"key_events"`       // 关键事件
	CurrentDecision string   `json:"current_decision"` // 当前待决策事件
	LifeGoals       []string `json:"life_goals"`       // 人生目标
	Challenges      []string `json:"challenges"`       // 当前挑战
}

// AIGameResponse AI游戏响应
type AIGameResponse struct {
	Success          bool               `json:"success"`
	EventTitle       string             `json:"event_title"`
	EventDescription string             `json:"event_description"`
	Decision         AIDecisionResponse `json:"decision"`
	AttributeChanges map[string]int     `json:"attribute_changes"`
	StatusChanges    map[string]string  `json:"status_changes"`
	NextAge          int                `json:"next_age"`
	NextLifeStage    string             `json:"next_life_stage"`
	GameEnd          bool               `json:"game_end"`
	GameEndReason    string             `json:"game_end_reason"`
}

// AIDecisionResponse AI决策响应
type AIDecisionResponse struct {
	Question string             `json:"question"`
	Context  string             `json:"context"`
	Options  []AIDecisionOption `json:"options"`
}

// AIDecisionOption AI决策选项
type AIDecisionOption struct {
	Type         string `json:"type"`         // conservative, moderate, aggressive
	Text         string `json:"text"`         // 选项描述（简短）
	Consequence  string `json:"consequence"`  // 预期后果描述
	Requirements string `json:"requirements"` // 前置条件
}

// AIDecisionRequest AI决策请求
type AIDecisionRequest struct {
	Character      AICharacter   `json:"character"`
	CurrentState   AIGameContext `json:"current_state"`
	DecisionID     string        `json:"decision_id"`
	SelectedOption string        `json:"selected_option"` // conservative, moderate, aggressive
}

// AIDecisionResult AI决策结果
type AIDecisionResult struct {
	Success           bool              `json:"success"`
	ResultDescription string            `json:"result_description"`
	AttributeChanges  map[string]int    `json:"attribute_changes"`
	StatusChanges     map[string]string `json:"status_changes"`
	NextEvents        []string          `json:"next_events"`
	LifeImpact        string            `json:"life_impact"`
}
