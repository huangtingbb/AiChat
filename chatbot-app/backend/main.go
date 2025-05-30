package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/config"
	"chatbot-app/backend/database"
	"chatbot-app/backend/router"
)

func main() {
	// 加载配置
	cfg := config.GetConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	if err := database.InitMySQL(&cfg.Database); err != nil {
		log.Fatalf("MySQL初始化失败: %v", err)
	}

	// // 初始化数据库表
	// if err := database.InitTables(); err != nil {
	// 	log.Fatalf("数据库表初始化失败: %v", err)
	// }

	// 初始化Redis连接
	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}

	// // 初始化Elasticsearch连接
	// if err := database.InitElasticsearch(&cfg.ES); err != nil {
	// 	log.Fatalf("Elasticsearch初始化失败: %v", err)
	// }

	// // 初始化RabbitMQ连接
	// if err := database.InitRabbitMQ(&cfg.RabbitMQ); err != nil {
	// 	log.Fatalf("RabbitMQ初始化失败: %v", err)
	// }
	// defer database.CloseRabbitMQ()

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	router.SetupRouter(r)

	// 启动服务器
	log.Printf("启动服务器在 :%s 端口", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}
