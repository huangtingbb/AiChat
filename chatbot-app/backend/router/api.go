package router

import (
	"github.com/gin-gonic/gin"

	"chatbot-app/backend/api"
)

// RegisterAPIRoutes 注册所有API路由
func RegisterAPIRoutes(apiGroup *gin.RouterGroup) *api.UserAPI {
	// 注册用户API路由（公开部分）
	userAPI := api.NewUserAPI()
	userAPI.RegisterRoutes(apiGroup)
	
	return userAPI
}

// RegisterProtectedAPIRoutes 注册需要认证的API路由
func RegisterProtectedAPIRoutes(authGroup *gin.RouterGroup, userAPI *api.UserAPI) {
	// 注册聊天API路由
	chatAPI := api.NewChatAPI()
	chatAPI.RegisterRoutes(authGroup)

	// 注册需要认证的用户API路由
	userAPI.RegisterProtectedRoutes(authGroup)
	
	// 可以在这里添加更多需要认证的API路由
} 