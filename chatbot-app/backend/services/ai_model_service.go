package services

import (
	"errors"
	"time"

	"chatbot-app/backend/database"
	"chatbot-app/backend/models"
)

// AIModelService AI模型服务
type AIModelService struct{}

// GetAllModelList 获取所有启用的AI模型
func (s *AIModelService) GetAllModelList(chatType string) ([]models.AIModel, error) {
	var aiModelList []models.AIModel
	query := database.DB.Where("enabled = ?", true)
	if chatType != "" {
		query = query.Where("type = ?", chatType)
	}
	if err := query.Find(&aiModelList).Error; err != nil {
		return nil, err
	}
	return aiModelList, nil
}

// GetModelById 根据Id获取AI模型
func (s *AIModelService) GetModelById(id uint) (*models.AIModel, error) {
	var aiModel models.AIModel
	if err := database.DB.First(&aiModel, id).Error; err != nil {
		return nil, errors.New("模型不存在")
	}
	return &aiModel, nil
}

// GetModelByName 根据名称获取AI模型
func (s *AIModelService) GetModelByName(name string) (*models.AIModel, error) {
	var aiModel models.AIModel
	if err := database.DB.Where("name = ?", name).First(&aiModel).Error; err != nil {
		return nil, errors.New("模型不存在")
	}
	return &aiModel, nil
}

// GetDefaultModel 获取默认AI模型
func (s *AIModelService) GetDefaultModel(chatType string) (*models.AIModel, error) {
	var aiModel models.AIModel
	query := database.DB.Where("is_default = ? AND enabled = ?", true, true)
	if chatType != "" {
		query = query.Where("type = ?", chatType)
	}
	if err := query.First(&aiModel).Error; err != nil {
		return nil, errors.New("未找到默认模型")
	}
	return &aiModel, nil
}

// RecordModelUsage 记录模型使用情况
func (s *AIModelService) RecordModelUsage(usage *models.AIModelUsage) error {
	return database.DB.Create(usage).Error
}

// GetModelUsageByUser 获取用户的模型使用记录
func (s *AIModelService) GetModelUsageByUser(userId uint) ([]models.AIModelUsage, error) {
	var usages []models.AIModelUsage
	if err := database.DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}

// CreateModelUsageFromResponse 从响应创建使用记录
func (s *AIModelService) CreateModelUsageFromResponse(
	userId uint,
	modelId uint,
	messageId uint,
	prompt string,
	response string,
	promptTokens int,
	completionTokens int,
	duration int,
) *models.AIModelUsage {
	return &models.AIModelUsage{
		UserId:           userId,
		ModelId:          modelId,
		MessageId:        messageId,
		Prompt:           prompt,
		Response:         response,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		Duration:         duration,
		Status:           "success",
		CreatedAt:        time.Now(),
	}
}

// CreateModelUsageError 创建错误使用记录
func (s *AIModelService) CreateModelUsageError(
	userId uint,
	modelId uint,
	prompt string,
	errorMsg string,
) *models.AIModelUsage {
	return &models.AIModelUsage{
		UserId:    userId,
		ModelId:   modelId,
		Prompt:    prompt,
		Status:    "error",
		ErrorMsg:  errorMsg,
		CreatedAt: time.Now(),
	}
}
