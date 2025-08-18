package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在白名单中
		//for _, path := range m.options.ExcludePaths {
		//	if c.Request.URL.Path == path {
		//		c.Next()
		//		return
		//	}
		//}

		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "token 无效，请先登录")
			c.Abort()
			return
		}
		token = parts[1]

		// 解析token
		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Unauthorized(c, "token 无效，请先登录")
			c.Abort()
			return
		}

		// 将claims保存到上下文
		c.Set("claims", claims)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
