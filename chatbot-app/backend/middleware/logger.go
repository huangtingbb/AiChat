package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"chatbot-app/backend/utils"

	"github.com/gin-gonic/gin"
)

// responseBodyWriter 用于捕获响应体
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器以捕获响应体
		responseBody := &bytes.Buffer{}
		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(startTime)

		// 获取用户Id（如果存在）
		userId := ""
		if claims, exists := c.Get("claims"); exists {
			if userClaims, ok := claims.(*utils.Claims); ok {
				userId = fmt.Sprintf("%d", userClaims.UserId)
			}
		}

		// 构建日志字段
		logFields := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency":    latency.String(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"user_id":    userId,
		}

		// 添加查询参数（如果存在）
		if len(c.Request.URL.RawQuery) > 0 {
			logFields["query"] = c.Request.URL.RawQuery
		}

		// 添加请求体（仅对非GET请求）
		if c.Request.Method != "GET" && len(requestBody) > 0 && len(requestBody) < 1024 {
			logFields["request_body"] = string(requestBody)
		}

		// 添加响应体（仅对错误状态码）
		if c.Writer.Status() >= 400 && responseBody.Len() > 0 && responseBody.Len() < 1024 {
			logFields["response_body"] = responseBody.String()
		}

		// 记录日志
		if c.Writer.Status() >= 500 {
			utils.LogError("服务器错误", nil, logFields)
		} else if c.Writer.Status() >= 400 {
			utils.LogWarn("客户端错误", logFields)
		} else {
			utils.LogInfo("请求处理完成", logFields)
		}
	}
}
