package services

import (
	"errors"
	"strings"
	"sync"
	"time"

	"chatbot-app/backend/models"
	"chatbot-app/backend/utils"
)

// ModelOption 模型选项
type ModelOption struct {
	Model       string  // 模型名称
	MaxTokens   int     // 最大token数
	Temperature float64 // 温度参数
	TopP        float64 // 多样性参数
}

// AiService AI服务
type AiService struct {
	// AI客户端工厂
	clientFactory *utils.AIClientFactory
	// 当前使用的模型
	model string
	// 模型选项
	options map[string]*utils.ChatCompletionOptions
	// AI模型服务
	modelService *AIModelService
	// 互斥锁
	mu sync.RWMutex
}

// NewAiService 创建新的AI服务实例
func NewAiService() *AiService {
	// 创建AI模型服务
	modelService := &AIModelService{}

	// 创建客户端工厂（配置从环境变量读取在工厂内部处理）
	clientFactory := utils.NewAIClientFactory()

	// 创建服务
	service := &AiService{
		clientFactory: clientFactory,
		model:         "",
		options:       make(map[string]*utils.ChatCompletionOptions),
		modelService:  modelService,
	}

	// 从数据库加载默认模型
	service.loadDefaultModelFromDB()

	return service
}

// loadDefaultModelFromDB 从数据库加载默认模型
func (s *AiService) loadDefaultModelFromDB() {
	// 从数据库获取默认模型
	defaultModel, err := s.modelService.GetDefaultModel()
	if err == nil {
		s.mu.Lock()
		s.model = defaultModel.Name

		// 更新选项
		s.options[defaultModel.Name] = utils.ConvertToClientOptions(defaultModel)
		s.mu.Unlock()
	} else {
		// 如果没有默认模型，获取所有可用模型并选择第一个
		modelList, err := s.modelService.GetAllModelList()
		if err == nil && len(modelList) > 0 {
			s.mu.Lock()
			s.model = modelList[0].Name

			// 更新选项
			s.options[modelList[0].Name] = utils.ConvertToClientOptions(&modelList[0])
			s.mu.Unlock()
		}
	}
}

// GenerateResponse 生成AI回复
func (s *AiService) GenerateResponse(aiModel *models.AIModel, prompt string, history []map[string]string, userID uint) (string, *models.AIModelUsage, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", nil, errors.New("提问内容不能为空")
	}

	// 开始计时
	startTime := time.Now()

	// 转换历史消息格式
	aiHistory := utils.ConvertHistoryMessages(history)

	// 获取AI客户端
	client, err := s.clientFactory.CreateClient(aiModel)
	if err != nil {
		// 创建错误记录
		errorUsage := s.modelService.CreateModelUsageError(userID, aiModel.ID, prompt, "获取AI客户端失败: "+err.Error())
		if recordErr := s.modelService.RecordModelUsage(errorUsage); recordErr != nil {
			// 记录日志
		}
		return "", errorUsage, errors.New("获取AI客户端失败: " + err.Error())
	}

	// 调用AI生成回复
	response, err := client.GenerateResponse(prompt, aiHistory)
	if err != nil {
		// 创建错误记录
		errorUsage := s.modelService.CreateModelUsageError(userID, aiModel.ID, prompt, "AI生成回复失败: "+err.Error())
		if recordErr := s.modelService.RecordModelUsage(errorUsage); recordErr != nil {
			// 记录日志
		}
		return "", errorUsage, errors.New("AI生成回复失败: " + err.Error())
	}

	// 计算耗时
	duration := int(time.Since(startTime).Milliseconds())

	// 创建使用记录
	// 这里的Token计算是简化的，实际项目中应该从API响应中获取
	promptTokens := len(prompt) / 4       // 简化的Token估算
	completionTokens := len(response) / 4 // 简化的Token估算
	usage := s.modelService.CreateModelUsageFromResponse(
		userID,
		aiModel.ID,
		0, // 消息ID后续设置
		prompt,
		response,
		promptTokens,
		completionTokens,
		duration,
	)

	// 记录使用情况
	if err := s.modelService.RecordModelUsage(usage); err != nil {
		// 记录日志，但不影响返回结果
	}

	return response, usage, nil
}

// GenerateStreamResponse 生成流式AI回复
//func (s *AiService) GenerateStreamResponse(prompt string, history []map[string]string, userID uint, callback func(string, error) bool) error {
//	// 这个方法展示了如何实现流式回复的框架
//	// 实际项目中需要根据具体的API和需求来实现
//	// 这里仅作为示例
//
//	// 检查输入
//	if strings.TrimSpace(prompt) == "" {
//		return errors.New("提问内容不能为空")
//	}
//
//	// 模拟流式响应
//	// 实际项目中，这里应该调用支持流式响应的AI API
//	response, _, err := s.GenerateResponse(prompt, history, userID)
//	if err != nil {
//		return err
//	}
//
//	// 模拟分段发送响应
//	// 实际项目中，这里应该处理流式API的响应
//	words := strings.Split(response, " ")
//	for _, word := range words {
//		// 调用回调函数，发送每个词
//		// 如果回调返回false，表示客户端已断开连接，停止发送
//		if !callback(word+" ", nil) {
//			break
//		}
//	}
//
//	return nil
//}
