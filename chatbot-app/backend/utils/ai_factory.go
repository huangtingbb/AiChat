package utils

import (
	"chatbot-app/backend/config"
	"chatbot-app/backend/models"
	"fmt"
)

// 定义支持的AI提供商常量
const (
	ProviderZhipu  = "zhipu"
	ProviderOpenAI = "openai"
	// 可以在这里添加更多提供商
)

// AIClientFactory AI客户端工厂
type AIClientFactory struct {
	// 配置实例
	config *config.Config
}

// NewAIClientFactory 创建AI客户端工厂
func NewAIClientFactory() *AIClientFactory {
	return &AIClientFactory{
		config: config.GetConfig(),
	}
}

// CreateClient 根据模型配置创建AI客户端
func (f *AIClientFactory) CreateClient(model *models.AIModel) (AIClient, error) {
	if model == nil {
		return nil, fmt.Errorf("模型配置不能为空")
	}

	switch model.Provider {
	case ProviderZhipu:
		return f.createZhipuClient(model)
	case ProviderOpenAI:
		// TODO: 实现OpenAI客户端
		return nil, fmt.Errorf("OpenAI客户端暂未实现")
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", model.Provider)
	}
}

// createZhipuClient 创建智谱AI客户端
func (f *AIClientFactory) createZhipuClient(model *models.AIModel) (AIClient, error) {
	// 从配置获取API Key
	apiKey := f.config.AI.ZhipuAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("智谱AI的API Key未配置，请在环境变量ZHIPU_API_KEY中设置")
	}

	// 使用模型配置的URL，如果没有则使用配置中的默认URL
	baseURL := model.URL
	// 创建聊天选项
	options := &ChatCompletionOptions{
		MaxTokens:        model.MaxTokens,
		Temperature:      model.Temperature,
		TopP:             model.TopP,
		PresencePenalty:  model.PresencePenalty,
		FrequencyPenalty: model.FrequencyPenalty,
		Stream:           false, // 暂时不支持流式响应
	}

	// 创建智谱AI客户端
	client := NewZhipuClient(apiKey, baseURL, model.Name, options)

	return client, nil
}

// createOpenAIClient 创建OpenAI客户端（预留）
// func (f *AIClientFactory) createOpenAIClient(model *models.AIModel) (AIClient, error) {
// 	// 从配置获取API Key
// 	apiKey := f.config.AI.OpenAIAPIKey
// 	if apiKey == "" {
// 		return nil, fmt.Errorf("OpenAI的API Key未配置，请在环境变量OPENAI_API_KEY中设置")
// 	}

// 	// 使用模型配置的URL，如果没有则使用配置中的默认URL
// 	baseURL := model.URL
// 	if baseURL == "" {
// 		baseURL = f.config.AI.OpenAIBaseURL
// 	}

// 	// TODO: 实现OpenAI客户端
// 	return nil, fmt.Errorf("OpenAI客户端尚未实现")
// }

// ConvertToClientOptions 将模型配置转换为客户端选项
func ConvertToClientOptions(model *models.AIModel) *ChatCompletionOptions {
	return &ChatCompletionOptions{
		MaxTokens:        model.MaxTokens,
		Temperature:      model.Temperature,
		TopP:             model.TopP,
		PresencePenalty:  model.PresencePenalty,
		FrequencyPenalty: model.FrequencyPenalty,
		Stream:           false,
	}
}

// ConvertHistoryMessages 转换历史消息格式
func ConvertHistoryMessages(history []map[string]string) []Message {
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		messages = append(messages, Message{
			Role:    msg["role"],
			Content: msg["content"],
		})
	}
	return messages
}
