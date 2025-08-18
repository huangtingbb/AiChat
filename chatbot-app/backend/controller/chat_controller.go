package controller

import (
	"encoding/json"
	"strconv"
	"strings"

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
		Title string `json:"title" binding:"required" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogWarn("创建聊天会话参数验证失败", map[string]interface{}{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		})
		utils.InvalidParams(c, err.Error())
		return
	}

	// 使用自定义验证器进行验证
	if validationError := utils.GetValidationError(req); validationError != "" {
		utils.LogWarn("创建聊天会话参数验证失败", map[string]interface{}{
			"error": validationError,
			"ip":    c.ClientIP(),
		})
		utils.InvalidParams(c, validationError)
		return
	}

	// 从JWT中获取用户Id
	claims, exists := c.Get("claims")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	userClaims := claims.(*utils.Claims)
	userId := userClaims.UserId

	utils.LogInfo("创建聊天会话请求", map[string]interface{}{
		"user_id": userId,
		"title":   req.Title,
		"ip":      c.ClientIP(),
	})

	chat, err := controller.chatService.CreateChat(userId, req.Title)
	if err != nil {
		utils.LogError("创建聊天会话失败", err, map[string]interface{}{
			"user_id": userId,
			"title":   req.Title,
		})
		utils.Error(c, err.Error())
		return
	}

	utils.LogInfo("聊天会话创建成功", map[string]interface{}{
		"user_id": userId,
		"chat_id": chat.Id,
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
	// 从JWT中获取用户Id
	claims, exists := c.Get("claims")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	userClaims := claims.(*utils.Claims)
	userId := userClaims.UserId

	chats, err := controller.chatService.GetUserChatList(userId)
	if err != nil {
		utils.LogError("获取用户聊天列表失败", err, map[string]interface{}{
			"user_id": userId,
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
// @Param id path integer true "聊天会话Id"
// @Success 200 {object} utils.Response{data=object{messages=array}} "消息列表"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "聊天会话不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/chat/{id}/message [get]
func (controller *ChatController) GetChatMessageList(c *gin.Context) {
	chatIdStr := c.Param("id")
	chatId, err := strconv.ParseUint(chatIdStr, 10, 64)
	if err != nil {
		utils.LogWarn("获取聊天消息：无效的聊天Id", map[string]interface{}{
			"chat_id": chatIdStr,
			"error":   err.Error(),
		})
		utils.InvalidParams(c, "无效的聊天Id")
		return
	}

	// 从JWT中获取用户Id
	claims, _ := c.Get("claims")
	userClaims := claims.(*utils.Claims)
	userId := userClaims.UserId

	// 验证聊天会话是否属于当前用户
	_, err = controller.chatService.GetChatById(uint(chatId), userId)
	if err != nil {
		utils.LogWarn("获取聊天消息：无权访问", map[string]interface{}{
			"user_id": userId,
			"chat_id": chatId,
			"error":   err.Error(),
		})
		utils.NotFound(c, err.Error())
		return
	}

	messages, err := controller.chatService.GetChatMessages(uint(chatId))
	if err != nil {
		utils.LogError("获取聊天消息失败", err, map[string]interface{}{
			"chat_id": chatId,
			"user_id": userId,
		})
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"messages": messages})
}

// SendMessage 发送消息（流式响应）
// @Summary 发送聊天消息（流式响应）
// @Description 在指定聊天会话中发送消息并获取AI流式回复
// @Tags 聊天
// @Accept json
// @Produce text/event-stream
// @Security Bearer
// @Param id path integer true "聊天会话Id"
// @Param body body object{content=string,model_id=integer} true "消息内容与可选模型Id"
// @Success 200 {string} string "Server-Sent Events流式响应"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "聊天会话不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/chat/{id}/message [post]
func (controller *ChatController) SendMessage(c *gin.Context) {
	chatIdStr := c.Param("id")
	chatId, err := strconv.ParseUint(chatIdStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "无效的聊天Id")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required" validate:"required"`
		ModelId uint   `json:"model_id" binding:"required" validate:"required"` // 模型Id
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	// 使用自定义验证器进行验证
	if validationError := utils.GetValidationError(req); validationError != "" {
		utils.LogWarn("发送消息参数验证失败", map[string]interface{}{
			"error": validationError,
			"ip":    c.ClientIP(),
		})
		utils.InvalidParams(c, validationError)
		return
	}

	selectedModel, err := controller.aiModelService.GetModelById(req.ModelId)

	if err != nil {
		utils.Error(c, "模型不存在，请重新选择")
	}

	// 从JWT中获取用户Id
	userId := c.GetUint("userId")

	utils.LogInfo("发送消息请求", map[string]interface{}{
		"user_id":  userId,
		"chat_id":  chatId,
		"model_id": selectedModel.Id,
		//"content":  req.Content[:min(len(req.Content), 100)], // 只记录前100个字符
	})

	// 验证聊天会话是否属于当前用户
	chat, err := controller.chatService.GetChatById(uint(chatId), userId)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	// 保存用户消息
	userMessage, err := controller.chatService.AddMessage(chat, "user", req.Content)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	// 获取历史消息作为上下文
	messages, err := controller.chatService.GetChatMessages(uint(chatId))
	if err != nil {
		utils.Error(c, "获取聊天历史失败")
		return
	}

	// 将消息转换为适合AI服务的格式
	var history []map[string]string
	for _, msg := range messages {
		if msg.Id == userMessage.Id {
			continue // 跳过刚刚添加的用户消息
		}
		history = append(history, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	// 设置流式响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 首先发送用户消息信息
	userMessageData, _ := json.Marshal(gin.H{
		"type":    "user_message",
		"message": userMessage,
	})
	c.SSEvent("message", string(userMessageData))
	c.Writer.Flush()

	// 发送流式开始信号
	startData, _ := json.Marshal(gin.H{
		"type":  "stream_start",
		"model": selectedModel.DisplayName,
	})
	c.SSEvent("message", string(startData))
	c.Writer.Flush()

	// 用于收集完整AI响应和相关数据
	var fullResponse strings.Builder
	var botMessage *models.Message

	// 定义流式回调函数
	streamCallback := func(chunk string, isEnd bool, err error) bool {
		if err != nil {
			utils.LogError("AI流式回复失败", err, map[string]interface{}{
				"user_id": userId,
				"chat_id": chatId,
			})

			// 发送错误消息
			errorData, _ := json.Marshal(gin.H{
				"type":  "error",
				"error": "AI生成回复失败: " + err.Error(),
			})
			c.SSEvent("message", string(errorData))
			c.Writer.Flush()
			return false
		}

		if isEnd {
			// 流式响应结束，保存完整回复
			response := fullResponse.String()

			// 保存AI回复到数据库
			savedBotMessage, saveErr := controller.chatService.AddMessageWithModelMetadata(
				uint(chatId),
				"assistant",
				response,
				selectedModel.Id,
				"", // 元数据稍后更新
			)
			if saveErr != nil {
				utils.LogError("保存AI回复失败", saveErr, map[string]interface{}{
					"user_id": userId,
					"chat_id": chatId,
				})

				// 发送保存错误消息
				errorData, _ := json.Marshal(gin.H{
					"type":  "error",
					"error": "保存回复失败",
				})
				c.SSEvent("message", string(errorData))
				c.Writer.Flush()
				return false
			}

			botMessage = savedBotMessage

			// 发送流式结束信号
			endData, _ := json.Marshal(gin.H{
				"type":       "stream_end",
				"message_id": botMessage.Id,
				"full_text":  response,
			})
			c.SSEvent("message", string(endData))
			c.Writer.Flush()

			return true
		}

		// 发送文本块
		fullResponse.WriteString(chunk)

		chunkData, _ := json.Marshal(gin.H{
			"type": "stream_chunk",
			"text": chunk,
		})
		c.SSEvent("message", string(chunkData))
		c.Writer.Flush()

		// 检查客户端是否断开连接
		select {
		case <-c.Request.Context().Done():
			utils.LogInfo("客户端断开连接，停止流式传输", map[string]interface{}{
				"user_id": userId,
				"chat_id": chatId,
			})
			return false
		default:
			return true
		}
	}

	// 调用AI服务生成流式回复
	err = controller.aiService.GenerateStreamResponse(selectedModel, req.Content, history, userId, streamCallback)
	if err != nil {
		utils.LogError("启动AI流式回复失败", err, map[string]interface{}{
			"user_id": userId,
			"chat_id": chatId,
		})

		// 发送启动错误消息
		errorData, _ := json.Marshal(gin.H{
			"type":  "error",
			"error": "AI回复失败: " + err.Error(),
		})
		c.SSEvent("message", string(errorData))
		c.Writer.Flush()
		return
	}

	// 记录完成日志
	utils.LogInfo("AI流式对话完成", map[string]interface{}{
		"user_id": userId,
		"chat_id": chatId,
		"model":   selectedModel.DisplayName,
		"message_id": func() uint {
			if botMessage != nil {
				return botMessage.Id
			}
			return 0
		}(),
		"response_length": fullResponse.Len(),
	})
}
