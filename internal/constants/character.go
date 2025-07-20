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

// 角色属性常量
const (
	AttributeIntelligence          = "intelligence"
	AttributeEmotionalIntelligence = "emotional_intelligence"
	AttributeMemory                = "memory"
	AttributeImagination           = "imagination"
	AttributePhysicalFitness       = "physical_fitness"
	AttributeAppearance            = "appearance"
)

// 角色状态常量
const (
	CharacterStatusActive   = 1
	CharacterStatusInactive = 0
	CharacterStatusDeleted  = -1
)

// 生活阶段常量
const (
	LifeStageInfant     = "infant"      // 婴儿 (0-2)
	LifeStageToddler    = "toddler"     // 幼儿 (3-5)
	LifeStageChild      = "child"       // 儿童 (6-12)
	LifeStageTeenager   = "teenager"    // 青少年 (13-17)
	LifeStageYoungAdult = "young_adult" // 青年 (18-30)
	LifeStageAdult      = "adult"       // 成年 (31-60)
	LifeStageElderly    = "elderly"     // 老年 (61+)
)

// GenderType 性别类型
type GenderType int

const (
	GenderUnknown GenderType = 0
	GenderMale    GenderType = 1
	GenderFemale  GenderType = 2
	GenderOther   GenderType = 3
)

var GenderMaps = map[GenderType]string{
	GenderUnknown: "unknown",
	GenderMale:    "male",
	GenderFemale:  "female",
	GenderOther:   "other",
}

func (g GenderType) String() string {
	if str, ok := GenderMaps[g]; ok {
		return str
	}
	return "unknown"
}

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

var RaceMaps = map[RaceType]string{
	RaceUnknown:   "unknown",
	RaceCaucasian: "caucasian",
	RaceAsian:     "asian",
	RaceAfrican:   "african",
	RaceHispanic:  "hispanic",
	RaceNative:    "native",
	RaceMixed:     "mixed",
}

func (r RaceType) String() string {
	if str, ok := RaceMaps[r]; ok {
		return str
	}
	return "unknown"
}
