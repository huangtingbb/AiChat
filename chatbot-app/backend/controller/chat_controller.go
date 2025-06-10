package controller

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/models"
	"chatbot-app/backend/services"
	"chatbot-app/backend/utils"
)

// ChatController 聊天控制器
type ChatController struct {
	chatService    services.ChatService
	aiService      *services.AiService
	aiModelService *services.AIModelService
}

// NewChatController 创建聊天控制器
func NewChatController() *ChatController {
	return &ChatController{
		chatService:    services.ChatService{},
		aiService:      services.NewAiService(),
		aiModelService: &services.AIModelService{},
	}
}

// CreateChat 创建聊天会话
// @Summary 创建聊天会话
// @Description 为当前用户创建一个新的聊天会话
// @Tags 聊天
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body object{title=string} true "聊天会话标题"
// @Success 200 {object} utils.Response{data=object} "创建成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/chat [post]
func (controller *ChatController) CreateChat(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogWarn("创建聊天会话参数验证失败", map[string]interface{}{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		})
		utils.InvalidParams(c, err.Error())
		return
	}

	// 从JWT中获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	utils.LogInfo("创建聊天会话请求", map[string]interface{}{
		"user_id": userID,
		"title":   req.Title,
		"ip":      c.ClientIP(),
	})

	chat, err := controller.chatService.CreateChat(userID, req.Title)
	if err != nil {
		utils.LogError("创建聊天会话失败", err, map[string]interface{}{
			"user_id": userID,
			"title":   req.Title,
		})
		utils.Error(c, err.Error())
		return
	}

	utils.LogInfo("聊天会话创建成功", map[string]interface{}{
		"user_id": userID,
		"chat_id": chat.ID,
		"title":   chat.Title,
	})

	utils.Success(c, chat)
}

// GetUserChatList 获取用户聊天会话列表
// @Summary 获取聊天会话列表
// @Description 获取当前用户的所有聊天会话列表
// @Tags 聊天
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=object{chats=array}} "聊天会话列表"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/chat [get]
func (controller *ChatController) GetUserChatList(c *gin.Context) {
	// 从JWT中获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	chats, err := controller.chatService.GetUserChatList(userID)
	if err != nil {
		utils.LogError("获取用户聊天列表失败", err, map[string]interface{}{
			"user_id": userID,
		})
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"chats": chats})
}

// GetChatMessageList 获取聊天消息列表
// @Summary 获取聊天消息列表
// @Description 获取指定聊天会话的所有消息
// @Tags 聊天
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path integer true "聊天会话ID"
// @Success 200 {object} utils.Response{data=object{messages=array}} "消息列表"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "聊天会话不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/chat/{id}/message [get]
func (controller *ChatController) GetChatMessageList(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		utils.LogWarn("获取聊天消息：无效的聊天ID", map[string]interface{}{
			"chat_id": chatIDStr,
			"error":   err.Error(),
		})
		utils.InvalidParams(c, "无效的聊天ID")
		return
	}

	// 从JWT中获取用户ID
	claims, _ := c.Get("claims")
	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	// 验证聊天会话是否属于当前用户
	_, err = controller.chatService.GetChatByID(uint(chatID), userID)
	if err != nil {
		utils.LogWarn("获取聊天消息：无权访问", map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err.Error(),
		})
		utils.NotFound(c, err.Error())
		return
	}

	messages, err := controller.chatService.GetChatMessages(uint(chatID))
	if err != nil {
		utils.LogError("获取聊天消息失败", err, map[string]interface{}{
			"chat_id": chatID,
			"user_id": userID,
		})
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"messages": messages})
}

// SendMessage 发送消息
// @Summary 发送聊天消息
// @Description 在指定聊天会话中发送消息并获取AI回复
// @Tags 聊天
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path integer true "聊天会话ID"
// @Param body body object{content=string,model_id=integer} true "消息内容与可选模型ID"
// @Success 200 {object} utils.Response{data=object{user_message=object,bot_message=object,metrics=object}} "发送成功，返回用户消息、AI回复及指标"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "聊天会话不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/chat/{id}/message [post]
func (controller *ChatController) SendMessage(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		utils.LogWarn("发送消息：无效的聊天ID", map[string]interface{}{
			"chat_id": chatIDStr,
			"error":   err.Error(),
		})
		utils.InvalidParams(c, "无效的聊天ID")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
		ModelID uint   `json:"model_id"` // 模型ID
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogWarn("发送消息参数验证失败", map[string]interface{}{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		})
		utils.InvalidParams(c, err.Error())
		return
	}

	// 从JWT中获取用户ID
	userID := c.GetUint("userID")

	utils.LogInfo("发送消息请求", map[string]interface{}{
		"user_id":  userID,
		"chat_id":  chatID,
		"model_id": req.ModelID,
		//"content":  req.Content[:min(len(req.Content), 100)], // 只记录前100个字符
	})

	// 验证聊天会话是否属于当前用户
	_, err = controller.chatService.GetChatByID(uint(chatID), userID)
	if err != nil {
		utils.LogWarn("发送消息：无权访问聊天会话", map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
			"error":   err.Error(),
		})
		utils.NotFound(c, err.Error())
		return
	}

	// 保存用户消息
	userMessage, err := controller.chatService.AddMessage(uint(chatID), "user", req.Content)
	if err != nil {
		utils.LogError("保存用户消息失败", err, map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
		})
		utils.Error(c, err.Error())
		return
	}

	// 获取历史消息作为上下文
	messages, err := controller.chatService.GetChatMessages(uint(chatID))
	if err != nil {
		utils.LogError("获取聊天历史失败", err, map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
		})
		utils.Error(c, "获取聊天历史失败")
		return
	}

	// 将消息转换为适合AI服务的格式
	var history []map[string]string
	for _, msg := range messages {
		if msg.ID == userMessage.ID {
			continue // 跳过刚刚添加的用户消息
		}
		history = append(history, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	// 如果请求中指定了模型ID，临时切换模型
	var selectedModel *models.AIModel
	if req.ModelID > 0 {
		selectedModel, err = controller.aiModelService.GetModelByID(req.ModelID)
		if err == nil {
			utils.LogInfo("临时切换AI模型", map[string]interface{}{
				"user_id":    userID,
				"model_id":   req.ModelID,
				"model_name": selectedModel.Name,
			})
		} else {
			utils.LogWarn("指定的模型不存在，使用默认模型", map[string]interface{}{
				"model_id": req.ModelID,
				"error":    err.Error(),
			})
		}
	} else {
		// 获取默认模型
		selectedModel, _ = controller.aiModelService.GetDefaultModel()
	}

	// 调用AI服务生成回复
	botResponse, usage, err := controller.aiService.GenerateResponse(selectedModel, req.Content, history, userID)
	if err != nil {
		utils.LogError("AI生成回复失败", err, map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
		})
		utils.Error(c, "AI生成回复失败: "+err.Error())
		return
	}

	// 构建元数据
	metadata := map[string]interface{}{
		"prompt_tokens":     usage.PromptTokens,
		"completion_tokens": usage.CompletionTokens,
		"total_tokens":      usage.TotalTokens,
		"duration_ms":       usage.Duration,
	}
	metadataJSON, _ := json.Marshal(metadata)

	// 保存AI回复
	botMessage, err := controller.chatService.AddMessageWithModelMetadata(
		uint(chatID),
		"assistant",
		botResponse,
		usage.ModelID,
		string(metadataJSON),
	)
	if err != nil {
		utils.LogError("保存AI回复失败", err, map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
		})
		utils.Error(c, err.Error())
		return
	}

	// 更新使用记录中的消息ID
	usage.MessageID = botMessage.ID
	if err := controller.aiModelService.RecordModelUsage(usage); err != nil {
		utils.LogWarn("记录模型使用情况失败", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
	}

	var modelName string
	if selectedModel != nil {
		modelName = selectedModel.DisplayName
	} else {
		modelName = "未知模型"
	}

	utils.LogInfo("AI对话完成", map[string]interface{}{
		"user_id":           userID,
		"chat_id":           chatID,
		"model":             modelName,
		"duration_ms":       usage.Duration,
		"prompt_tokens":     usage.PromptTokens,
		"completion_tokens": usage.CompletionTokens,
		"total_tokens":      usage.TotalTokens,
	})

	utils.Success(c, gin.H{
		"user_message": userMessage,
		"bot_message":  botMessage,
		"metrics": gin.H{
			"model":             modelName,
			"duration_ms":       usage.Duration,
			"prompt_tokens":     usage.PromptTokens,
			"completion_tokens": usage.CompletionTokens,
			"total_tokens":      usage.TotalTokens,
		},
	})
}
