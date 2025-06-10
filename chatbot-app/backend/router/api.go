package router

import (
	"chatbot-app/backend/controller"
	"chatbot-app/backend/middleware"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/utils"
)

// SetupRouter 配置所有路由
func SetupRouter(r *gin.Engine) {
	// 健康检查路由
	// @Summary 健康检查
	// @Description 检查API服务是否正常运行
	// @Tags 系统
	// @Accept json
	// @Produce json
	// @Success 200 {object} utils.Response "返回服务状态"
	// @Router /health [get]
	r.GET("/health", healthCheckHandler)

	userController := controller.NewUserController()

	//不需要认证的路由
	r.POST("/api/user/login", userController.Login)
	r.POST("/api/user/register", userController.Register)

	// API路由组
	api := r.Group("/api", middleware.Auth())
	{
		user := api.Group("/user")
		{
			user.GET("/info", userController.GetUserInfo)
			user.POST("/logout", userController.Logout)
		}
		chat := api.Group("/chat")
		chatController := controller.NewChatController()
		{
			chat.POST("", chatController.CreateChat)
			chat.GET("", chatController.GetUserChatList)
			chat.GET("/:id/message", chatController.GetChatMessageList)
			chat.POST("/:id/message", chatController.SendMessage)
		}
		// AI模型相关路由
		ai := api.Group("/ai")
		aiModelController := controller.NewAIModelController()
		{
			ai.GET("/model", aiModelController.GetAvailableModelList)
			ai.GET("/model_usage", aiModelController.GetModelUsageHandler)
		}
	}
}

// healthCheckHandler 健康检查
// @Summary 健康检查
// @Description 检查API服务是否正常运行
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "返回服务状态"
// @Router /health [get]
func healthCheckHandler(c *gin.Context) {
	utils.Success(c, gin.H{
		"status": "ok",
	})
}
