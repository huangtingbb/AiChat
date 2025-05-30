package middleware

import (
	"github.com/gin-gonic/gin"
)

// JWT 返回JWT中间件
// 保留此函数以保持向后兼容
func JWT() gin.HandlerFunc {
	jwtAuth := NewJWTAuth()
	return jwtAuth.Middleware()
}

// JWTWithOptions 创建自定义选项的JWT中间件
func JWTWithOptions(options ...func(*JWTOptions)) gin.HandlerFunc {
	jwtAuth := NewJWTAuth(options...)
	return jwtAuth.Middleware()
}
