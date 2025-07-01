package config

import (
	"fmt"
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

// validateConfig 验证配置参数
func validateConfig(config *Config) error {
	// 检查阿里云AccessKey配置
	if config.Aliyun.AccessKeyId == "" || config.Aliyun.AccessKeyId == "${ACCESS_KEY_ID}" {
		return fmt.Errorf("阿里云AccessKeyId未配置")
	}
	if config.Aliyun.AccessKeySecret == "" || config.Aliyun.AccessKeySecret == "${ACCESS_KEY_SECRET}" {
		return fmt.Errorf("阿里云AccessKeySecret未配置")
	}
	if config.Aliyun.RegionId == "" || config.Aliyun.RegionId == "${REGION_ID}" {
		return fmt.Errorf("阿里云RegionId未配置")
	}

	// 检查服务器端口配置
	if config.Server.Port == "" || config.Server.Port == "${PORT}" {
		return fmt.Errorf("服务器端口未配置")
	}
	// 可以添加端口号格式验证，但由于端口可能来自环境变量，这里只做简单检查
	if len(config.Server.Port) > 5 {
		return fmt.Errorf("服务器端口号格式不正确")
	}

	return nil
}

// bindEnvVariables 绑定环境变量
func bindEnvVariables() error {
	envVars := map[string]string{
		"server.port":              "PORT",
		"aliyun.access_key_id":     "ACCESS_KEY_ID",
		"aliyun.access_key_secret": "ACCESS_KEY_SECRET",
		"aliyun.region_id":         "REGION_ID",
	}

	for configKey, envKey := range envVars {
		err := viper.BindEnv(configKey, envKey)
		if err != nil {
			return fmt.Errorf("绑定环境变量失败 %s: %w", configKey, err)
		}
	}

	return nil
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		// 如果 .env 文件不存在，不返回错误，继续使用环境变量
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("加载.env文件失败: %w", err)
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
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 允许环境变量覆盖配置文件
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DNS_UPDATE") // 环境变量前缀

	// 设置环境变量映射
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 绑定环境变量
	if err := bindEnvVariables(); err != nil {
		return nil, err
	}

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 设置默认值
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Aliyun.RegionId == "" {
		config.Aliyun.RegionId = "cn-hangzhou"
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}
