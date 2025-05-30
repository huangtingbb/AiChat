package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/middleware"
)

// SetupRouter 配置所有路由
func SetupRouter(r *gin.Engine) {
	// 允许跨域
	r.Use(corsMiddleware())

	// 健康检查路由
	r.GET("/health", healthCheckHandler)

	// API路由组
	apiGroup := r.Group("/api")

	// 注册API路由（公开部分）
	userAPI := RegisterAPIRoutes(apiGroup)

	// 创建JWT中间件
	jwtMiddleware := createJWTMiddleware()

	// 需要JWT认证的路由组
	authGroup := apiGroup.Group("")
	authGroup.Use(jwtMiddleware)

	// 注册需要认证的API路由
	RegisterProtectedAPIRoutes(authGroup, userAPI)
}

// corsMiddleware 处理跨域请求
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// healthCheckHandler 健康检查处理器
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// createJWTMiddleware 创建JWT中间件
func createJWTMiddleware() gin.HandlerFunc {
	return middleware.JWTWithOptions(
		middleware.WithExcludePaths([]string{
			"/api/users/register",
			"/api/users/login",
			"/api/health",
		}),
		middleware.WithErrorResponse(gin.H{
			"error": "身份验证失败，请重新登录",
			"code":  http.StatusUnauthorized,
		}),
	)
}
