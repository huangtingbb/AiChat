package api

import (
	"chatbot-app/backend/services"
	"chatbot-app/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserAPI 用户API
type UserAPI struct {
	userService services.UserService
}

// NewUserAPI 创建用户API
func NewUserAPI() *UserAPI {
	return &UserAPI{
		userService: services.UserService{},
	}
}

// RegisterHandler 注册处理器
func (api *UserAPI) RegisterHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6,max=50"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := api.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "user": user})
}

// LoginHandler 登录处理器
func (api *UserAPI) LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := api.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user":    user,
	})
}

// GetUserInfoHandler 获取用户信息处理器
func (api *UserAPI) GetUserInfoHandler(c *gin.Context) {
	// 从JWT中获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userClaims := claims.(*utils.Claims)

	user, err := api.userService.GetUserByID(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// LogoutHandler 退出登录处理器
func (api *UserAPI) LogoutHandler(c *gin.Context) {
	// 从JWT中获取用户ID和token
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userClaims := claims.(*utils.Claims)

	// 这里可以添加token黑名单逻辑
	// 例如将token添加到Redis中并设置过期时间
	// api.userService.InvalidateToken(token)

	c.JSON(http.StatusOK, gin.H{
		"message": "退出成功",
		"user_id": userClaims.UserID,
	})
}

// RegisterRoutes 注册路由
func (api *UserAPI) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		// 公开接口，无需认证
		users.POST("/register", api.RegisterHandler)
		users.POST("/login", api.LoginHandler)
	}

	// 需要中间件保护的用户接口在main.go中使用authGroup进行注册
}

// RegisterProtectedRoutes 注册需要认证的路由
func (api *UserAPI) RegisterProtectedRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/info", api.GetUserInfoHandler)
		users.POST("/logout", api.LogoutHandler) // 添加退出登录路由
		// 其他需要认证的用户API路由
	}
}
