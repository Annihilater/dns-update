package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Aliyun AliyunConfig `mapstructure:"aliyun"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// AliyunConfig 阿里云配置
type AliyunConfig struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	RegionId        string `mapstructure:"region_id"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		// 如果 .env 文件不存在，不返回错误，继续使用环境变量
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// 设置配置文件信息
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 如果提供了配置路径，就使用它
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}

	// 默认在当前目录和configs目录下查找
	viper.AddConfigPath(".")
	viper.AddConfigPath("configs")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 允许环境变量覆盖配置文件
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DNS_UPDATE") // 环境变量前缀

	// 设置环境变量映射
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 绑定环境变量
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("aliyun.access_key_id", "ACCESS_KEY_ID")
	viper.BindEnv("aliyun.access_key_secret", "ACCESS_KEY_SECRET")
	viper.BindEnv("aliyun.region_id", "REGION_ID")

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 设置默认值
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Aliyun.RegionId == "" {
		config.Aliyun.RegionId = "cn-hangzhou"
	}

	return &config, nil
}
