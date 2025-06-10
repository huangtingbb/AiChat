package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 状态码常量
const (
	// StatusSuccess 成功状态码
	StatusSuccess = 0
	// StatusError 错误状态码
	StatusError = -1
)

// Response API标准响应结构
type Response struct {
	Code      int         `json:"code"`      // 状态码
	Message   string      `json:"message"`   // 消息
	Data      interface{} `json:"data"`      // 数据
	Timestamp int64       `json:"timestamp"` // 时间戳
}

// NewResponse 创建一个新的响应
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewResponse(StatusSuccess, "success", data))
}

// SuccessWithMsg 返回带有自定义消息的成功响应
func SuccessWithMsg(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, NewResponse(StatusSuccess, message, data))
}

// Error 返回错误响应
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, NewResponse(StatusError, message, nil))
}

// Custom 返回带有自定义状态码的响应
func Custom(c *gin.Context, code int, message string) {
	c.JSON(code, NewResponse(code, message, nil))
}

// InvalidParams 返回参数错误响应
func InvalidParams(c *gin.Context, message string) {
	c.JSON(http.StatusOK, NewResponse(-1, message, nil))
}

// Unauthorized 返回未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, NewResponse(401, message, nil))
}

// Forbidden 返回禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, NewResponse(403, message, nil))
}

// NotFound 返回资源不存在响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, NewResponse(404, message, nil))
}
