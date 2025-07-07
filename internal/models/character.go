package models

import "time"

// Character 游戏角色模型
type Character struct {
	CharacterID   string `json:"character_id" db:"character_id"`
	UserID        uint   `json:"user_id" db:"user_id"`
	CharacterName string `json:"character_name" db:"character_name"`
	BirthCountry  string `json:"birth_country" db:"birth_country"`
	BirthYear     int    `json:"birth_year" db:"birth_year"`
	CurrentAge    int    `json:"current_age" db:"current_age"`
	Gender        int    `json:"gender" db:"gender"` // 0:未知 1:男 2:女 3:其他
	Race          int    `json:"race" db:"race"`     // 0:未知 1:白人 2:黄种人 3:黑人
	IsActive      bool   `json:"is_active" db:"is_active"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`

	// 角色属性
	Attributes CharacterAttributes `json:"attributes"`

	// 角色扩展信息
	EducationLevel      *string `json:"education_level" db:"education_level"`
	MaritalStatus       *string `json:"marital_status" db:"marital_status"`
	CurrentCountry      *string `json:"current_country" db:"current_country"`
	CurrentLocation     *string `json:"current_location" db:"current_location"`
	CurrentActivity     *string `json:"current_activity" db:"current_activity"`
	Personality         *string `json:"personality" db:"personality"`
	Career              *string `json:"career" db:"career"`
	SkillTendency       *string `json:"skill_tendency" db:"skill_tendency"`
	FamilyBackground    *string `json:"family_background" db:"family_background"`
	SocialRelationships *string `json:"social_relationships" db:"social_relationships"`
	CareerDesc          *string `json:"career_desc" db:"career_desc"`
	EducationDesc       *string `json:"education_desc" db:"education_desc"`

	// 游戏状态
	LifeStage      string `json:"life_stage" db:"life_stage"`
	CurrentStatus  string `json:"current_status" db:"current_status"`
	HappinessLevel int    `json:"happiness_level" db:"happiness_level"`
	HealthLevel    int    `json:"health_level" db:"health_level"`
	Money          int64  `json:"money" db:"money"`

	// 游戏进度
	TotalPlaytime int     `json:"total_playtime" db:"total_playtime"`
	GameCompleted bool    `json:"game_completed" db:"game_completed"`
	FinalAge      *int    `json:"final_age" db:"final_age"`
	DeathCause    *string `json:"death_cause" db:"death_cause"`
	Summary       *string `json:"summary" db:"summary"` // 游戏总结，简要描述角色的一生经历、成就等
}

func (c *Character) SetDefaultValues() {
	now := time.Now().UnixMilli()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.IsActive = true
	c.CurrentAge = 0
	c.LifeStage = string(LifeStageBirth)
	c.CurrentStatus = "healthy"
	c.HappinessLevel = 50
	c.HealthLevel = 100
	c.Money = 0
	c.TotalPlaytime = 0
	c.GameCompleted = false
}

// CharacterAttributes 角色属性
type CharacterAttributes struct {
	Intelligence          int `json:"intelligence" db:"intelligence"`                     // 智力 0-100
	EmotionalIntelligence int `json:"emotional_intelligence" db:"emotional_intelligence"` // 情商 0-100
	Memory                int `json:"memory" db:"memory"`                                 // 记忆力 0-100
	Imagination           int `json:"imagination" db:"imagination"`                       // 想象力 0-100
	PhysicalFitness       int `json:"physical_fitness" db:"physical_fitness"`             // 体质 0-100
	Appearance            int `json:"appearance" db:"appearance"`                         // 外貌 0-100
}

// CreateCharacterRequest 创建角色请求
type CreateCharacterRequest struct {
	CharacterName string `json:"character_name" binding:"required"`
	BirthCountry  string `json:"birth_country" binding:"required"`
	BirthYear     int    `json:"birth_year" binding:"required,min=1800,max=2050"`
	Gender        int    `json:"gender" binding:"required,min=0,max=3"`
	Race          int    `json:"race" binding:"required,min=0,max=10"`
}

// UpdateCharacterRequest 更新角色请求
type UpdateCharacterRequest struct {
	CharacterID         string  `json:"character_id" binding:"required,uuid"`
	CharacterName       *string `json:"character_name,omitempty" binding:"omitempty,min=1,max=100"`
	EducationLevel      *string `json:"education_level,omitempty"`
	MaritalStatus       *string `json:"marital_status,omitempty"`
	CurrentCountry      *string `json:"current_country,omitempty"`
	CurrentLocation     *string `json:"current_location,omitempty"`
	CurrentActivity     *string `json:"current_activity,omitempty"`
	Personality         *string `json:"personality,omitempty"`
	Career              *string `json:"career,omitempty"`
	SkillTendency       *string `json:"skill_tendency,omitempty"`
	FamilyBackground    *string `json:"family_background,omitempty"`
	SocialRelationships *string `json:"social_relationships,omitempty"`
	CareerDesc          *string `json:"career_desc,omitempty"`
	EducationDesc       *string `json:"education_desc,omitempty"`
}

// UpdateCharacterAttributesRequest 更新角色属性请求
type UpdateCharacterAttributesRequest struct {
	Intelligence          *int `json:"intelligence,omitempty" binding:"omitempty,min=0,max=100"`
	EmotionalIntelligence *int `json:"emotional_intelligence,omitempty" binding:"omitempty,min=0,max=100"`
	Memory                *int `json:"memory,omitempty" binding:"omitempty,min=0,max=100"`
	Imagination           *int `json:"imagination,omitempty" binding:"omitempty,min=0,max=100"`
	PhysicalFitness       *int `json:"physical_fitness,omitempty" binding:"omitempty,min=0,max=100"`
	Appearance            *int `json:"appearance,omitempty" binding:"omitempty,min=0,max=100"`
}

// CharacterListResponse 角色列表响应
type CharacterListResponse struct {
	Characters []CharacterSummary `json:"characters"`
	Total      int                `json:"total"`
}

// CharacterSummary 角色概要信息（用于列表显示）
type CharacterSummary struct {
	CharacterID   string              `json:"character_id"`
	CharacterName string              `json:"character_name"`
	CurrentAge    int                 `json:"current_age"`
	Gender        int                 `json:"gender"`
	Race          int                 `json:"race"`
	LifeStage     string              `json:"life_stage"`
	Attributes    CharacterAttributes `json:"attributes"`
	CreatedAt     int64               `json:"created_at"`
	IsActive      bool                `json:"is_active"`
	Summary       *string             `json:"summary,omitempty"` // 游戏总结，简要描述角色的一生经历、成就等
}

// GenderType 性别类型
type GenderType int

const (
	GenderUnknown GenderType = 0
	GenderMale    GenderType = 1
	GenderFemale  GenderType = 2
	GenderOther   GenderType = 3
)

// RaceType 种族类型
type RaceType int

const (
	RaceUnknown   RaceType = 0
	RaceCaucasian RaceType = 1 // 白人
	RaceAsian     RaceType = 2 // 黄种人
	RaceAfrican   RaceType = 3 // 黑人
	RaceHispanic  RaceType = 4 // 拉丁裔
	RaceNative    RaceType = 5 // 原住民
	RaceMixed     RaceType = 6 // 混血
)

// LifeStageType 人生阶段类型
type LifeStageType string

const (
	LifeStageBirth       LifeStageType = "birth"       // 出生
	LifeStageChildhood   LifeStageType = "childhood"   // 童年（0-12）
	LifeStageAdolescence LifeStageType = "adolescence" // 青春期（13-17）
	LifeStageAdulthood   LifeStageType = "adulthood"   // 成年（18-59）
	LifeStageOldAge      LifeStageType = "old_age"     // 老年（60+）
	LifeStageDeath       LifeStageType = "death"       // 死亡
)

// GetGenderString 获取性别字符串
func (g GenderType) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	case GenderOther:
		return "other"
	default:
		return "unknown"
	}
}

// GetRaceString 获取种族字符串
func (r RaceType) String() string {
	switch r {
	case RaceCaucasian:
		return "caucasian"
	case RaceAsian:
		return "asian"
	case RaceAfrican:
		return "african"
	case RaceHispanic:
		return "hispanic"
	case RaceNative:
		return "native"
	case RaceMixed:
		return "mixed"
	default:
		return "unknown"
	}
}

// GetLifeStageByAge 根据年龄获取人生阶段
func GetLifeStageByAge(age int) LifeStageType {
	switch {
	case age == 0:
		return LifeStageBirth
	case age >= 1 && age <= 12:
		return LifeStageChildhood
	case age >= 13 && age <= 17:
		return LifeStageAdolescence
	case age >= 18 && age <= 59:
		return LifeStageAdulthood
	default:
		return LifeStageOldAge
	}
}
