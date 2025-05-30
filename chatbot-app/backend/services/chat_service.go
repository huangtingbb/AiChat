package services

import (
	"chatbot-app/backend/database"
	"chatbot-app/backend/models"
	"errors"
	"fmt"
	"time"
)

// ChatService 聊天服务
type ChatService struct{}

// CreateChat 创建聊天会话
func (s *ChatService) CreateChat(userID uint, title string) (*models.Chat, error) {
	chat := &models.Chat{
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(chat).Error; err != nil {
		return nil, err
	}

	return chat, nil
}

// GetUserChats 获取用户的所有聊天会话
func (s *ChatService) GetUserChats(userID uint) ([]models.Chat, error) {
	var chats []models.Chat
	if err := database.DB.Where("user_id = ?", userID).Order("id DESC").Find(&chats).Error; err != nil {
		return nil, err
	}
	fmt.Println(chats)

	return chats, nil
}

// GetChatByID 根据ID获取聊天会话
func (s *ChatService) GetChatByID(chatID, userID uint) (*models.Chat, error) {
	var chat models.Chat
	if err := database.DB.Where("id = ? AND user_id = ?", chatID, userID).First(&chat).Error; err != nil {
		return nil, errors.New("聊天会话不存在")
	}
	return &chat, nil
}

// AddMessage 添加消息到聊天会话
func (s *ChatService) AddMessage(chatID uint, role, content string) (*models.Message, error) {
	// 验证聊天会话是否存在
	var chat models.Chat
	if err := database.DB.First(&chat, chatID).Error; err != nil {
		return nil, errors.New("聊天会话不存在")
	}

	// 创建消息
	message := &models.Message{
		ChatID:    chatID,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

// GetChatMessages 获取聊天会话的所有消息
func (s *ChatService) GetChatMessages(chatID uint) ([]models.Message, error) {
	var messages []models.Message
	if err := database.DB.Where("chat_id = ?", chatID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
