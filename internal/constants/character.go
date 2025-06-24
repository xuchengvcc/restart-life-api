package constants

// 角色相关常量
const (
	MaxCharacterNameLength = 100
	MinCharacterAge        = 0
	MaxCharacterAge        = 200
	MinAttributeValue      = 0
	MaxAttributeValue      = 100
	MaxCharactersPerUser   = 10 // 每个用户最多角色数
)

// 种族常量（示例，可根据游戏设定调整）
const (
	RaceUnknown   = 0
	RaceCaucasian = 1
	RaceAsian     = 2
	RaceAfrican   = 3
	RaceHispanic  = 4
	RaceNative    = 5
	RaceMixed     = 6
)

// 角色属性常量
const (
	AttributeIntelligence         = "intelligence"
	AttributeEmotionalIntelligence = "emotional_intelligence"
	AttributeMemory               = "memory"
	AttributeImagination          = "imagination"
	AttributePhysicalFitness      = "physical_fitness"
	AttributeAppearance           = "appearance"
)

// 角色状态常量
const (
	CharacterStatusActive   = 1
	CharacterStatusInactive = 0
	CharacterStatusDeleted  = -1
)

// 生活阶段常量
const (
	LifeStageInfant     = "infant"     // 婴儿 (0-2)
	LifeStageToddler    = "toddler"    // 幼儿 (3-5)
	LifeStageChild      = "child"      // 儿童 (6-12)
	LifeStageTeenager   = "teenager"   // 青少年 (13-17)
	LifeStageYoungAdult = "young_adult" // 青年 (18-30)
	LifeStageAdult      = "adult"      // 成年 (31-60)
	LifeStageElderly    = "elderly"    // 老年 (61+)
)
