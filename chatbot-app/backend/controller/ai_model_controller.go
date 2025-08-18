package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/models"
	"chatbot-app/backend/services"
	"chatbot-app/backend/utils"
)

// AIModelController AI模型控制器
type AIModelController struct {
	aiService      *services.AiService
	aiModelService *services.AIModelService
}

// NewAIModelController 创建AI模型控制器
func NewAIModelController() *AIModelController {
	return &AIModelController{
		aiService:      services.NewAiService(),
		aiModelService: &services.AIModelService{},
	}
}

// GetAvailableModelList 获取可用模型列表
// @Summary 获取可用模型列表
// @Description 获取系统中所有可用的AI模型列表及当前使用的模型
// @Tags AI模型
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=object{models=array,current_model=string}} "模型列表"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/ai/model [get]
func (controller *AIModelController) GetAvailableModelList(c *gin.Context) {
	modelList, err := controller.aiModelService.GetAllModelList()
	if err != nil {
		utils.LogError("获取模型列表失败", err)
		utils.Error(c, "获取模型列表失败: "+err.Error())
		return
	}

	defaultModel, _ := controller.aiModelService.GetDefaultModel()

	utils.Success(c, gin.H{
		"models":        modelList,
		"current_model": defaultModel,
	})
}

// GetModelUsageHandler 获取模型使用情况
// @Summary 获取模型使用情况
// @Description 获取当前用户的AI模型使用记录
// @Tags AI模型
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query integer false "页码" default(1)
// @Param limit query integer false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=object{usage=array,total=integer}} "使用记录"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/ai/usage [get]
func (controller *AIModelController) GetModelUsageHandler(c *gin.Context) {
	// 从JWT中获取用户Id
	claims, exists := c.Get("claims")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	userClaims := claims.(*utils.Claims)
	userId := userClaims.UserId

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	utils.LogInfo("获取模型使用记录请求", map[string]interface{}{
		"user_id": userId,
		"page":    page,
		"limit":   limit,
	})

	// 获取用户的模型使用记录
	usageList, err := controller.aiModelService.GetModelUsageByUser(userId)
	if err != nil {
		utils.LogError("获取模型使用记录失败", err, map[string]interface{}{
			"user_id": userId,
		})
		utils.Error(c, "获取使用记录失败: "+err.Error())
		return
	}

	// 简单分页处理
	total := len(usageList)
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		usageList = []models.AIModelUsage{}
	} else {
		if end > total {
			end = total
		}
		usageList = usageList[start:end]
	}

	utils.Success(c, gin.H{
		"usage": usageList,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
