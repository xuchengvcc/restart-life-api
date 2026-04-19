package models

// AchievementItem 成就项
type AchievementItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Unlocked    bool   `json:"unlocked"`
	Progress    int    `json:"progress"`
	MaxProgress int    `json:"max_progress"`
	UnlockedAt  *int64 `json:"unlocked_at,omitempty"`
}

// AchievementCategory 成就分类
type AchievementCategory struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	TotalCount    int    `json:"total_count"`
	UnlockedCount int    `json:"unlocked_count"`
}

// CharacterAchievementsResponse 角色成就响应
type CharacterAchievementsResponse struct {
	CharacterID   string            `json:"character_id"`
	CharacterName string            `json:"character_name"`
	Items         []AchievementItem `json:"items"`
	UnlockedCount int               `json:"unlocked_count"`
	TotalCount    int               `json:"total_count"`
}

// CharacterStatsResponse 角色统计响应
type CharacterStatsResponse struct {
	CharacterID    string              `json:"character_id"`
	CharacterName  string              `json:"character_name"`
	CurrentAge     int                 `json:"current_age"`
	LifeStage      string              `json:"life_stage"`
	IsGameActive   bool                `json:"is_game_active"`
	TotalPlaytime  int                 `json:"total_playtime"`
	Money          int64               `json:"money"`
	EventCount     int                 `json:"event_count"`
	PendingDecision bool               `json:"pending_decision"`
	Attributes     CharacterAttributes `json:"attributes"`
	LastEvent      *Event              `json:"last_event,omitempty"`
}

// TimelineItem 时间线项
type TimelineItem struct {
	Age         int    `json:"age"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	CreatedAt   int64  `json:"created_at"`
}

// CharacterTimelineResponse 角色时间线响应
type CharacterTimelineResponse struct {
	CharacterID string         `json:"character_id"`
	Timeline    []TimelineItem `json:"timeline"`
	Total       int            `json:"total"`
}

