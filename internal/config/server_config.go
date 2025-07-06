package config

import (
	"time"

	"github.com/spf13/viper"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	HTTPSPort    string        `mapstructure:"https_port"`
	Mode         string        `mapstructure:"mode"`
	EnableHTTP   bool          `mapstructure:"enable_http"`
	EnableHTTPS  bool          `mapstructure:"enable_https"`
	SSLCertFile  string        `mapstructure:"ssl_cert_file"`
	SSLKeyFile   string        `mapstructure:"ssl_key_file"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	Version      string        `mapstructure:"version"`
}

func setServerConfigDefault() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.https_port", "8443")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.enable_http", true)
	viper.SetDefault("server.enable_https", true)
	viper.SetDefault("server.ssl_cert_file", "/etc/nginx/asecondchance.cn_bundle.crt")
	viper.SetDefault("server.ssl_key_file", "/etc/nginx/asecondchance.cn.key")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.version", "0.1.0")
}
