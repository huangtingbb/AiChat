package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"chatbot-app/backend/services"
	"chatbot-app/backend/utils"
)

// ChatAPI 聊天API
type ChatAPI struct {
	chatService services.ChatService
}

// NewChatAPI 创建聊天API
func NewChatAPI() *ChatAPI {
	return &ChatAPI{
		chatService: services.ChatService{},
	}
}

// CreateChatHandler 创建聊天会话处理器
func (api *ChatAPI) CreateChatHandler(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT中获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userClaims := claims.(*utils.Claims)

	chat, err := api.chatService.CreateChat(userClaims.UserID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}

// GetUserChatsHandler 获取用户聊天会话列表处理器
func (api *ChatAPI) GetUserChatsHandler(c *gin.Context) {
	// 从JWT中获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userClaims := claims.(*utils.Claims)

	chats, err := api.chatService.GetUserChats(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// GetChatMessagesHandler 获取聊天消息列表处理器
func (api *ChatAPI) GetChatMessagesHandler(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的聊天ID"})
		return
	}

	// 从JWT中获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userClaims := claims.(*utils.Claims)

	// 验证聊天会话是否属于当前用户
	_, err = api.chatService.GetChatByID(uint(chatID), userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	messages, err := api.chatService.GetChatMessages(uint(chatID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendMessageHandler 发送消息处理器
func (api *ChatAPI) SendMessageHandler(c *gin.Context) {
	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的聊天ID"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT中获取用户ID
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userClaims := claims.(*utils.Claims)

	// 验证聊天会话是否属于当前用户
	_, err = api.chatService.GetChatByID(uint(chatID), userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 保存用户消息
	userMessage, err := api.chatService.AddMessage(uint(chatID), "user", req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: 调用AI模型生成回复
	botResponse := "这是AI助手的回复，实际功能需要对接真实的AI服务"

	// 保存AI回复
	botMessage, err := api.chatService.AddMessage(uint(chatID), "assistant", botResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_message": userMessage,
		"bot_message":  botMessage,
	})
}

// RegisterRoutes 注册路由
func (api *ChatAPI) RegisterRoutes(router *gin.RouterGroup) {
	chats := router.Group("/chats")
	{
		chats.POST("", api.CreateChatHandler)
		chats.GET("", api.GetUserChatsHandler)
		chats.GET("/:id/messages", api.GetChatMessagesHandler)
		chats.POST("/:id/messages", api.SendMessageHandler)
	}
}
