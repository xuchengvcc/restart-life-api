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

// CountryInfo 国家信息
type CountryInfo struct {
	Code       string  // 国家代码
	Name       string  // 国家名称
	NameCN     string  // 中文名称
	Population float64 // 人口权重（百万）
}

// 全球主要国家列表（按人口权重）
var Countries = []CountryInfo{
	{Code: "CN", Name: "China", NameCN: "中国", Population: 1439.32},
	{Code: "IN", Name: "India", NameCN: "印度", Population: 1380.00},
	{Code: "US", Name: "United States", NameCN: "美国", Population: 331.00},
	{Code: "ID", Name: "Indonesia", NameCN: "印度尼西亚", Population: 273.52},
	{Code: "PK", Name: "Pakistan", NameCN: "巴基斯坦", Population: 220.89},
	{Code: "BR", Name: "Brazil", NameCN: "巴西", Population: 212.56},
	{Code: "NG", Name: "Nigeria", NameCN: "尼日利亚", Population: 206.14},
	{Code: "BD", Name: "Bangladesh", NameCN: "孟加拉国", Population: 164.69},
	{Code: "RU", Name: "Russia", NameCN: "俄罗斯", Population: 145.93},
	{Code: "MX", Name: "Mexico", NameCN: "墨西哥", Population: 128.93},
	{Code: "JP", Name: "Japan", NameCN: "日本", Population: 126.48},
	{Code: "ET", Name: "Ethiopia", NameCN: "埃塞俄比亚", Population: 114.96},
	{Code: "PH", Name: "Philippines", NameCN: "菲律宾", Population: 109.58},
	{Code: "EG", Name: "Egypt", NameCN: "埃及", Population: 102.33},
	{Code: "VN", Name: "Vietnam", NameCN: "越南", Population: 97.34},
	{Code: "CD", Name: "Democratic Republic of the Congo", NameCN: "刚果民主共和国", Population: 89.56},
	{Code: "TR", Name: "Turkey", NameCN: "土耳其", Population: 84.34},
	{Code: "IR", Name: "Iran", NameCN: "伊朗", Population: 83.99},
	{Code: "DE", Name: "Germany", NameCN: "德国", Population: 83.78},
	{Code: "TH", Name: "Thailand", NameCN: "泰国", Population: 69.80},
	{Code: "GB", Name: "United Kingdom", NameCN: "英国", Population: 67.89},
	{Code: "FR", Name: "France", NameCN: "法国", Population: 65.27},
	{Code: "IT", Name: "Italy", NameCN: "意大利", Population: 60.46},
	{Code: "TZ", Name: "Tanzania", NameCN: "坦桑尼亚", Population: 59.73},
	{Code: "ZA", Name: "South Africa", NameCN: "南非", Population: 59.31},
	{Code: "MM", Name: "Myanmar", NameCN: "缅甸", Population: 54.41},
	{Code: "KE", Name: "Kenya", NameCN: "肯尼亚", Population: 53.77},
	{Code: "KR", Name: "South Korea", NameCN: "韩国", Population: 51.27},
	{Code: "CO", Name: "Colombia", NameCN: "哥伦比亚", Population: 50.88},
	{Code: "ES", Name: "Spain", NameCN: "西班牙", Population: 46.75},
	{Code: "UG", Name: "Uganda", NameCN: "乌干达", Population: 45.74},
	{Code: "AR", Name: "Argentina", NameCN: "阿根廷", Population: 45.20},
	{Code: "DZ", Name: "Algeria", NameCN: "阿尔及利亚", Population: 43.85},
	{Code: "SD", Name: "Sudan", NameCN: "苏丹", Population: 43.85},
	{Code: "UA", Name: "Ukraine", NameCN: "乌克兰", Population: 43.73},
	{Code: "IQ", Name: "Iraq", NameCN: "伊拉克", Population: 40.22},
	{Code: "AF", Name: "Afghanistan", NameCN: "阿富汗", Population: 38.93},
	{Code: "PL", Name: "Poland", NameCN: "波兰", Population: 37.85},
	{Code: "CA", Name: "Canada", NameCN: "加拿大", Population: 37.74},
	{Code: "MA", Name: "Morocco", NameCN: "摩洛哥", Population: 36.91},
	{Code: "SA", Name: "Saudi Arabia", NameCN: "沙特阿拉伯", Population: 34.81},
	{Code: "UZ", Name: "Uzbekistan", NameCN: "乌兹别克斯坦", Population: 33.47},
	{Code: "PE", Name: "Peru", NameCN: "秘鲁", Population: 32.97},
	{Code: "MY", Name: "Malaysia", NameCN: "马来西亚", Population: 32.37},
	{Code: "AO", Name: "Angola", NameCN: "安哥拉", Population: 32.87},
	{Code: "MZ", Name: "Mozambique", NameCN: "莫桑比克", Population: 31.26},
	{Code: "GH", Name: "Ghana", NameCN: "加纳", Population: 31.07},
	{Code: "YE", Name: "Yemen", NameCN: "也门", Population: 29.83},
	{Code: "NP", Name: "Nepal", NameCN: "尼泊尔", Population: 29.14},
	{Code: "VE", Name: "Venezuela", NameCN: "委内瑞拉", Population: 28.44},
	{Code: "MG", Name: "Madagascar", NameCN: "马达加斯加", Population: 27.69},
}

// GetCountryByCode 根据国家代码获取国家信息
func GetCountryByCode(code string) *CountryInfo {
	for _, country := range Countries {
		if country.Code == code {
			return &country
		}
	}
	return nil
}

// GetCountryNames 获取所有国家名称列表（用于前端选择）
func GetCountryNames() []string {
	names := make([]string, len(Countries))
	for i, country := range Countries {
		names[i] = country.NameCN
	}
	return names
}

// GetCountryCodes 获取所有国家代码列表
func GetCountryCodes() []string {
	codes := make([]string, len(Countries))
	for i, country := range Countries {
		codes[i] = country.Code
	}
	return codes
}
