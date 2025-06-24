package constants

// 游戏核心常量
const (
	// 游戏开始年龄范围
	MinStartAge = 0
	MaxStartAge = 18

	// 游戏结束条件
	MaxLifespan = 120
	MinLifespan = 0

	// 属性变化限制
	MaxAttributeChange = 10  // 单次属性变化最大值
	MinAttributeChange = -10 // 单次属性变化最小值
)

// 事件类型常量
const (
	EventTypeLife      = "life"      // 生活事件
	EventTypeEducation = "education" // 教育事件
	EventTypeCareer    = "career"    // 职业事件
	EventTypeRelation  = "relation"  // 关系事件
	EventTypeHealth    = "health"    // 健康事件
	EventTypeRandom    = "random"    // 随机事件
)

// 事件影响类型
const (
	EffectTypePositive = "positive" // 正面影响
	EffectTypeNegative = "negative" // 负面影响
	EffectTypeNeutral  = "neutral"  // 中性影响
)

// 游戏难度常量
const (
	DifficultyEasy   = "easy"
	DifficultyNormal = "normal"
	DifficultyHard   = "hard"
	DifficultyInsane = "insane"
)

// 生活模式常量
const (
	LifeModeRealistic = "realistic" // 现实模式
	LifeModeFantasy   = "fantasy"   // 幻想模式
	LifeModeCustom    = "custom"    // 自定义模式
)

// 国家/地区常量（示例）
const (
	CountryChina  = "China"
	CountryUSA    = "USA"
	CountryJapan  = "Japan"
	CountryUK     = "UK"
	CountryFrance = "France"
	CountryGermany = "Germany"
)

// 教育程度常量
const (
	EducationNone       = "none"
	EducationPrimary    = "primary"
	EducationSecondary  = "secondary"
	EducationHighSchool = "high_school"
	EducationCollege    = "college"
	EducationBachelor   = "bachelor"
	EducationMaster     = "master"
	EducationPhD        = "phd"
)

// 婚姻状况常量
const (
	MaritalStatusSingle   = "single"
	MaritalStatusDating   = "dating"
	MaritalStatusEngaged  = "engaged"
	MaritalStatusMarried  = "married"
	MaritalStatusDivorced = "divorced"
	MaritalStatusWidowed  = "widowed"
)
