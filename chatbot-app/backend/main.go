package main

import (
	"chatbot-app/backend/config"
	"chatbot-app/backend/database"
	_ "chatbot-app/backend/docs" // 导入swagger文档
	"chatbot-app/backend/middleware"
	"chatbot-app/backend/router"
	"chatbot-app/backend/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           聊天机器人API
// @version         1.0
// @description     这是一个聊天机器人的API文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description                 Bearer JWT token for authentication

func main() {

	// 加载配置
	cfg := config.GetConfig()

	// 初始化日志系统
	if err := utils.InitLogger(&cfg.Log); err != nil {
		panic("日志系统初始化失败: " + err.Error())
	}

	utils.LogInfo("应用程序启动", map[string]interface{}{
		"version": "1.0",
		"mode":    cfg.Server.Mode,
		"port":    cfg.Server.Port,
	})

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化验证器
	if err := utils.InitValidator(); err != nil {
		utils.LogError("验证器初始化失败", err)
		panic(err)
	}
	utils.LogInfo("验证器初始化成功")

	// 初始化数据库连接
	if err := database.InitMySQL(&cfg.Database); err != nil {
		utils.LogError("MySQL初始化失败", err)
		panic(err)
	}
	utils.LogInfo("MySQL连接成功")

	// 初始化Redis连接
	if err := database.InitRedis(&cfg.Redis); err != nil {
		utils.LogError("Redis初始化失败", err)
		panic(err)
	}
	utils.LogInfo("Redis连接成功")

	// 创建Gin引擎
	r := gin.New()
	// default 默认包含Recovery、 Logger 中间件
	//r = gin.Default()

	// 添加中间件
	r.Use(gin.Recovery())
	r.Use(middleware.CorsMiddleware())
	r.Use(gin.Logger())

	// 设置路由
	router.SetupRouter(r)

	// 添加Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	utils.LogInfo("启动HTTP服务器", map[string]interface{}{
		"port": cfg.Server.Port,
	})

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		utils.LogError("无法启动服务器", err)
		panic(err)
	}
}
