package config

// Config 应用程序配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	ES       ElasticsearchConfig
	RabbitMQ RabbitMQConfig
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

// ElasticsearchConfig Elasticsearch配置
type ElasticsearchConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

// GetConfig 获取配置
func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Mode: "release",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "root",
			DBName:   "chatbot",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "102914",
			DB:       0,
		},
		ES: ElasticsearchConfig{
			Host:     "localhost",
			Port:     "9200",
			User:     "elastic",
			Password: "password",
		},
		RabbitMQ: RabbitMQConfig{
			Host:     "localhost",
			Port:     "5672",
			User:     "guest",
			Password: "guest",
			VHost:    "/",
		},
	}
}
