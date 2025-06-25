package config

import (
	"time"

	"github.com/spf13/viper"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func setMySQLConfigDefault() {
	viper.SetDefault("database.mysql.host", "localhost")
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.database", "restart_life_dev")
	viper.SetDefault("database.mysql.username", "root")
	viper.SetDefault("database.mysql.password", "password")
	viper.SetDefault("database.mysql.charset", "utf8mb4")
	viper.SetDefault("database.mysql.parse_time", true)
	viper.SetDefault("database.mysql.max_open_conns", 10)
	viper.SetDefault("database.mysql.max_idle_conns", 5)
	viper.SetDefault("database.mysql.conn_max_lifetime", "300s")
}
