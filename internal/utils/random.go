package utils

import (
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/xuchengvcc/restart-life-api/internal/constants"
)

// AgeRange 年龄段定义
type AgeRange struct {
	MinAge int
	MaxAge int
	Weight float64
}

// 年龄段权重配置（基于真实人口分布）
var ageRanges = []AgeRange{
	{0, 5, 5.0},    // 婴幼儿
	{6, 12, 8.0},   // 儿童
	{13, 17, 6.0},  // 青少年
	{18, 25, 15.0}, // 青年
	{26, 35, 20.0}, // 青壮年
	{36, 45, 18.0}, // 中年
	{46, 55, 12.0}, // 中老年
	{56, 65, 8.0},  // 老年
	{66, 80, 6.0},  // 高龄
	{81, 100, 2.0}, // 超高龄
}

// 预计算的总权重
var ageRangesTotalWeight float64

func init() {
	// 初始化时计算总权重
	for _, ageRange := range ageRanges {
		ageRangesTotalWeight += ageRange.Weight
	}
}

// RandomCountrySelector 随机国家选择器
type RandomCountrySelector struct {
	countries      []constants.CountryInfo
	totalWeight    float64
	weightedRanges []float64
}

var (
	// 全局单例实例
	countrySelector *RandomCountrySelector
	once            sync.Once
)

// GetCountrySelector 获取国家选择器单例
func GetCountrySelector() *RandomCountrySelector {
	once.Do(func() {
		countrySelector = &RandomCountrySelector{
			countries: constants.Countries,
		}
		// 计算总权重和累积权重范围
		countrySelector.calculateWeights()
	})
	return countrySelector
}

// calculateWeights 计算权重
func (r *RandomCountrySelector) calculateWeights() {
	r.totalWeight = 0
	r.weightedRanges = make([]float64, len(r.countries))

	for i, country := range r.countries {
		r.totalWeight += country.Population
		r.weightedRanges[i] = r.totalWeight
	}
}

// SelectRandomCountry 根据人口权重随机选择国家
func (r *RandomCountrySelector) SelectRandomCountry() constants.CountryInfo {
	if len(r.countries) == 0 {
		return constants.CountryInfo{Code: "CN", Name: "China", NameCN: "中国", Population: 1439.32}
	}

	// 使用更好的随机数生成器
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randomValue := rng.Float64() * r.totalWeight

	// 使用标准库的二分查找
	index := sort.Search(len(r.weightedRanges), func(i int) bool {
		return r.weightedRanges[i] >= randomValue
	})

	// 确保索引在有效范围内
	if index >= len(r.countries) {
		index = len(r.countries) - 1
	}

	return r.countries[index]
}

// SelectRandomBirthYear 随机选择出生年份
func SelectRandomBirthYear() int {
	currentYear := time.Now().Year()

	// 生成随机数选择年龄段（使用预计算的总权重）
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randomValue := r.Float64() * ageRangesTotalWeight

	cumulativeWeight := 0.0
	for _, ageRange := range ageRanges {
		cumulativeWeight += ageRange.Weight
		if randomValue <= cumulativeWeight {
			// 在选中的年龄段内随机选择具体年龄
			age := r.Intn(ageRange.MaxAge-ageRange.MinAge+1) + ageRange.MinAge
			return currentYear - age
		}
	}

	// 默认返回25岁对应的出生年份
	return currentYear - 25
}

// GetDefaultCountryCode 获取默认国家代码
func GetDefaultCountryCode() string {
	return "CN" // 默认中国
}

// GetRandomCountryCode 快速获取随机国家代码（直接使用单例）
func GetRandomCountryCode() string {
	return GetCountrySelector().SelectRandomCountry().Code
}

// GetRandomCountryInfo 快速获取随机国家信息（直接使用单例）
func GetRandomCountryInfo() constants.CountryInfo {
	return GetCountrySelector().SelectRandomCountry()
}

// GetAgeRanges 获取年龄段配置（用于前端显示或其他用途）
func GetAgeRanges() []AgeRange {
	return ageRanges
}

// GetAgeRangesTotalWeight 获取年龄段总权重
func GetAgeRangesTotalWeight() float64 {
	return ageRangesTotalWeight
}
