package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用程序配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
	AI       AIConfig
}

// AIConfig AI配置
type AIConfig struct {
	ZhipuAPIKey   string
	ZhipuBaseURL  string
	OpenAIAPIKey  string
	OpenAIBaseURL string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string // debug, info, warn, error
	Format   string // json, text
	LogFile  string // 日志文件路径
	MaxSize  int    // 日志文件最大大小(MB)
	MaxAge   int    // 日志保留天数
	Compress bool   // 是否压缩旧日志
}

// GetConfig 获取配置，优先从环境变量读取
func GetConfig() *Config {
	godotenv.Load()
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root"),
			DBName:   getEnv("DB_NAME", "chatbot"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", "102914"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Log: LogConfig{
			Level:    getEnv("LOG_LEVEL", "info"),
			Format:   getEnv("LOG_FORMAT", "json"),
			LogFile:  getEnv("LOG_FILE", "logs/app.log"),
			MaxSize:  getEnvAsInt("LOG_MAX_SIZE", 100),
			MaxAge:   getEnvAsInt("LOG_MAX_AGE", 30),
			Compress: getEnvAsBool("LOG_COMPRESS", true),
		},
		AI: AIConfig{
			ZhipuAPIKey: getEnv("ZHIPU_API_KEY", ""),
		},
	}
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int，如果不存在或转换失败则使用默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为bool，如果不存在或转换失败则使用默认值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
