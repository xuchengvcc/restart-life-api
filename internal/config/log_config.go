package config

import "github.com/spf13/viper"

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

func setLogConfigDefault() {
	viper.SetDefault("logging.level", "debug")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
}
