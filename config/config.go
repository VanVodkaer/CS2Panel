package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// 全局变量 GlobalConfig
var GlobalConfig *Config

// 初始化全局变量
func init() {
	// 从配置文件加载日志配置
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	} else {
		GlobalConfig = cfg
		log.Printf("配置文件加载成功: %+v\n", cfg)
	}

}

// Config 结构体存储应用程序的配置
type Config struct {
	Env struct {
		Mode     string `mapstructure:"mode"`
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"env"`

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Docker struct {
		ImageName  string `mapstructure:"image_name"`
		Tag        string `mapstructure:"tag"`
		Prefix     string `mapstructure:"prefix"`
		MaxRetries int    `mapstructure:"max_retries"`
		RetryDelay int    `mapstructure:"retry_delay"`
	} `mapstructure:"docker"`

	Util struct {
		LogDir        string `mapstructure:"log_dir"`
		LogFileName   string `mapstructure:"log_filename"`
		LogMaxSize    int    `mapstructure:"log_max_size"`
		LogMaxBackups int    `mapstructure:"log_max_backups"`
		LogMaxAge     int    `mapstructure:"log_max_age"`
		LogCompress   bool   `mapstructure:"log_compress"`
	} `mapstructure:"util"`
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	// 初始化 viper
	viper.SetConfigName("config")   // 配置文件名（不含后缀）
	viper.SetConfigType("yaml")     // 配置文件类型
	viper.AddConfigPath("./config") // 配置文件路径（当前目录下的 config 文件夹）
	viper.AutomaticEnv()            // 允许环境变量覆盖配置

	// 设置默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("docker.image_name", "joedwards32/cs2")
	viper.SetDefault("docker.tag", "latest")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 创建配置结构体实例
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("无法解组配置: %w", err)
	}

	return &config, nil
}
