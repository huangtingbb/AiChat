package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/utils"
)

// AuthMiddleware 定义认证中间件接口
type AuthMiddleware interface {
	// Middleware 返回Gin中间件函数
	Middleware() gin.HandlerFunc
}

// JWTAuthMiddleware JWT认证中间件
type JWTAuthMiddleware struct {
	// 自定义选项
	options JWTOptions
}

// JWTOptions JWT中间件选项
type JWTOptions struct {
	// 路径白名单，这些路径不需要JWT验证
	ExcludePaths []string
	// 解析token失败的错误响应
	ErrorResponse gin.H
	// 用于获取token的函数
	TokenExtractor func(*gin.Context) (string, error)
}

// NewJWTAuth 创建新的JWT认证中间件
func NewJWTAuth(options ...func(*JWTOptions)) *JWTAuthMiddleware {
	// 默认选项
	opts := JWTOptions{
		ExcludePaths: []string{},
		ErrorResponse: gin.H{
			"error": "未授权",
			"code":  http.StatusUnauthorized,
		},
		TokenExtractor: ExtractTokenFromHeader,
	}

	// 应用自定义选项
	for _, option := range options {
		option(&opts)
	}

	return &JWTAuthMiddleware{
		options: opts,
	}
}

// WithExcludePaths 设置不需要JWT验证的路径
func WithExcludePaths(paths []string) func(*JWTOptions) {
	return func(o *JWTOptions) {
		o.ExcludePaths = paths
	}
}

// WithErrorResponse 设置错误响应
func WithErrorResponse(response gin.H) func(*JWTOptions) {
	return func(o *JWTOptions) {
		o.ErrorResponse = response
	}
}

// WithTokenExtractor 设置token提取函数
func WithTokenExtractor(extractor func(*gin.Context) (string, error)) func(*JWTOptions) {
	return func(o *JWTOptions) {
		o.TokenExtractor = extractor
	}
}

// Middleware 实现AuthMiddleware接口
func (m *JWTAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在白名单中
		for _, path := range m.options.ExcludePaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// 获取token
		token, err := m.options.TokenExtractor(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, m.options.ErrorResponse)
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, m.options.ErrorResponse)
			c.Abort()
			return
		}

		// 将claims保存到上下文
		c.Set("claims", claims)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// ExtractTokenFromHeader 从Authorization头提取token
func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrMissingToken
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrInvalidTokenFormat
	}

	return parts[1], nil
}

// ExtractTokenFromQuery 从查询参数中提取token
func ExtractTokenFromQuery(paramName string) func(*gin.Context) (string, error) {
	return func(c *gin.Context) (string, error) {
		token := c.Query(paramName)
		if token == "" {
			return "", ErrMissingToken
		}
		return token, nil
	}
}

// ExtractTokenFromCookie 从Cookie中提取token
func ExtractTokenFromCookie(cookieName string) func(*gin.Context) (string, error) {
	return func(c *gin.Context) (string, error) {
		cookie, err := c.Cookie(cookieName)
		if err != nil || cookie == "" {
			return "", ErrMissingToken
		}
		return cookie, nil
	}
}

// 定义错误类型
var (
	ErrMissingToken       = NewAuthError("未提供token")
	ErrInvalidTokenFormat = NewAuthError("token格式错误")
)

// AuthError 认证错误
type AuthError struct {
	Message string
}

// Error 实现error接口
func (e AuthError) Error() string {
	return e.Message
}

// NewAuthError 创建认证错误
func NewAuthError(message string) AuthError {
	return AuthError{Message: message}
}
