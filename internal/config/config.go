package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	AI       AIConfig       `mapstructure:"ai"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	// 设置默认值
	setDefaults()

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 支持嵌套字段用下划线，如 DATABASE_MYSQL_PORT
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 手动处理数组环境变量
	processAIProvidersFromEnv(&config)

	return &config, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setDefaults()

	var config Config
	viper.Unmarshal(&config)

	return &config
}

// setDefaults 设置默认配置值
func setDefaults() {
	setServerConfigDefault()
	setMySQLConfigDefault()
	setRedisConfigDefault()
	setAuthConfigDefault()
	setCORConfigDefault()
	setLogConfigDefault()
	setAIConfigDefault()
}

func (c *Config) PostInit() {
	c.AI.PostInit()
}

// processAIProvidersFromEnv 从环境变量处理AI Providers配置
func processAIProvidersFromEnv(config *Config) {
	// 处理数组形式的环境变量 AI_PROVIDERS_0_*, AI_PROVIDERS_1_*, ...
	for i := 0; ; i++ {
		nameKey := "AI_PROVIDERS_" + strconv.Itoa(i) + "_NAME"
		apiKeyKey := "AI_PROVIDERS_" + strconv.Itoa(i) + "_API_KEY"

		name := os.Getenv(nameKey)
		apiKey := os.Getenv(apiKeyKey)

		// 如果没有找到这个索引的配置，退出循环
		if name == "" && apiKey == "" {
			break
		}

		// 确保providers数组有足够的元素
		for len(config.AI.Providers) <= i {
			config.AI.Providers = append(config.AI.Providers, &Provider{})
		}

		// 更新配置
		if name != "" {
			config.AI.Providers[i].Name = name
		}
		if apiKey != "" {
			config.AI.Providers[i].APIKey = apiKey
		}

		// 处理model_ids数组
		for j := 0; ; j++ {
			modelKey := "AI_PROVIDERS_" + strconv.Itoa(i) + "_MODEL_IDS_" + strconv.Itoa(j)
			modelID := os.Getenv(modelKey)

			if modelID == "" {
				break
			}

			// 确保model_ids数组有足够的元素
			for len(config.AI.Providers[i].ModelIDs) <= j {
				config.AI.Providers[i].ModelIDs = append(config.AI.Providers[i].ModelIDs, "")
			}

			config.AI.Providers[i].ModelIDs[j] = modelID
		}
	}
}
