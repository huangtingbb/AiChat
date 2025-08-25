package controller

import (
	"chatbot-app/backend/services"
	"chatbot-app/backend/utils"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: services.UserService{},
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body object{username=string,password=string,email=string} true "注册信息"
// @Success 200 {object} utils.Response{data=object{user=object}} "注册成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/user/register [post]
func (controller *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50" msg_required:"请输入用户名" msg_min:"用户名至少需3个字符" msg_max:"用户名不能超过50个字符"`
		Password string `json:"password" binding:"required,min=6,max=50" msg_required:"请设置密码" msg_min:"密码至少需6位数" msg_max:"密码不能超过50位"`
		Email    string `json:"email" binding:"required,email" msg_required:"请输入邮箱地址" msg_email:"邮箱格式不正确"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := utils.GetValidationErrorWithTagMessages(req, err)
		utils.LogWarn("用户注册参数验证失败", map[string]interface{}{
			"error": validationError,
			"ip":    c.ClientIP(),
		})
		utils.InvalidParams(c, validationError)
		return
	}

	utils.LogInfo("用户注册请求", map[string]interface{}{
		"username": req.Username,
		"email":    req.Email,
		"ip":       c.ClientIP(),
	})

	user, err := controller.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		utils.LogError("用户注册失败", err, map[string]interface{}{
			"username": req.Username,
			"email":    req.Email,
			"ip":       c.ClientIP(),
		})
		utils.Error(c, err.Error())
		return
	}

	utils.LogInfo("用户注册成功", map[string]interface{}{
		"user_id":  user.Id,
		"username": user.Username,
		"ip":       c.ClientIP(),
	})

	utils.SuccessWithMsg(c, "注册成功", user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取认证令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body object{username=string,password=string} true "登录信息"
// @Success 200 {object} utils.Response{data=object{token=string,user=object}} "登录成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/user/login [post]
func (controller *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" form:"username" binding:"required" msg_required:"请输入用户名"`
		Password string `json:"password" form:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&req); err != nil {
		// 使用自定义验证器进行验证
		validationError := utils.GetValidationErrorWithTagMessages(req, err)
		utils.LogWarn("用户登录参数验证失败", map[string]interface{}{
			"error": validationError,
			"ip":    c.ClientIP(),
		})
		utils.InvalidParams(c, validationError)
		return
	}

	user, err := controller.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	// 生成JWT
	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		utils.Error(c, "生成token失败")
		return
	}
	utils.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
		"user":  user,
	})
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=object{user=object}} "用户信息"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/user/info [get]
func (controller *UserController) GetUserInfo(c *gin.Context) {
	// 从JWT中获取用户Id
	claims, exists := c.Get("claims")
	if !exists {
		utils.LogWarn("获取用户信息失败：未授权", map[string]interface{}{
			"ip": c.ClientIP(),
		})
		utils.Unauthorized(c, "未授权")
		return
	}

	userClaims := claims.(*utils.Claims)

	user, err := controller.userService.GetUserById(userClaims.UserId)
	if err != nil {
		utils.LogWarn(err.Error(), map[string]interface{}{
			"user_id": userClaims.UserId,
			"error":   err.Error(),
			"ip":      c.ClientIP(),
		})
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"user": user})
}

// Logout 退出登录
// @Summary 退出登录
// @Description 使当前用户退出登录状态
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=object{user_id=integer}} "退出成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/user/logout [post]
func (controller *UserController) Logout(c *gin.Context) {
	// 从JWT中获取用户Id
	claims, exists := c.Get("claims")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	userClaims := claims.(*utils.Claims)

	utils.LogInfo("用户退出登录", map[string]interface{}{
		"user_id": userClaims.UserId,
		"ip":      c.ClientIP(),
	})

	utils.SuccessWithMsg(c, "退出成功", gin.H{
		"user_id": userClaims.UserId,
	})
}
