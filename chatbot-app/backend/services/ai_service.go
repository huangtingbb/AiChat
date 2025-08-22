package services

import (
	"errors"
	"fmt"
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
func (s *AiService) GenerateResponse(aiModel *models.AIModel, prompt string, history []map[string]string, userId uint) (string, *models.AIModelUsage, error) {
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
		errorUsage := s.modelService.CreateModelUsageError(userId, aiModel.Id, prompt, "获取AI客户端失败: "+err.Error())
		if recordErr := s.modelService.RecordModelUsage(errorUsage); recordErr != nil {
			// 记录日志
		}
		return "", errorUsage, errors.New("获取AI客户端失败: " + err.Error())
	}

	// 调用AI生成回复
	response, err := client.GenerateResponse(prompt, aiHistory)
	if err != nil {
		// 创建错误记录
		errorUsage := s.modelService.CreateModelUsageError(userId, aiModel.Id, prompt, "AI生成回复失败: "+err.Error())
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
		userId,
		aiModel.Id,
		0, // 消息Id后续设置
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
func (s *AiService) GenerateStreamResponse(aiModel *models.AIModel, prompt string, history []map[string]string, userId uint, callback func(chunk string, isEnd bool, err error) bool) error {
	// 检查输入
	if strings.TrimSpace(prompt) == "" {
		return errors.New("提问内容不能为空")
	}

	// 如果没有指定模型，获取默认模型
	if aiModel == nil {
		defaultModel, err := s.modelService.GetDefaultModel()
		if err != nil {
			return errors.New("获取默认模型失败: " + err.Error())
		}
		aiModel = defaultModel
	}

	// 开始计时
	startTime := time.Now()

	// 检查是否为Coze提供商，使用专门的服务
	if aiModel.Provider == "coze" {
		return s.handleCozeStreamResponse(aiModel, prompt, history, userId, callback, startTime)
	}

	// 转换历史消息格式
	aiHistory := utils.ConvertHistoryMessages(history)

	// 获取AI客户端
	client, err := s.clientFactory.CreateClient(aiModel)
	if err != nil {
		// 创建错误记录
		errorUsage := s.modelService.CreateModelUsageError(userId, aiModel.Id, prompt, "获取AI客户端失败: "+err.Error())
		if recordErr := s.modelService.RecordModelUsage(errorUsage); recordErr != nil {
			// 记录日志
		}
		return errors.New("获取AI客户端失败: " + err.Error())
	}

	// 用于收集完整响应以记录使用情况
	var fullResponse strings.Builder
	var hasError bool

	// 定义内部回调函数，包装用户回调并处理使用记录
	internalCallback := func(chunk string, isEnd bool, err error) bool {
		if err != nil {
			hasError = true
			// 创建错误记录
			errorUsage := s.modelService.CreateModelUsageError(userId, aiModel.Id, prompt, "AI流式生成回复失败: "+err.Error())
			if recordErr := s.modelService.RecordModelUsage(errorUsage); recordErr != nil {
				// 记录日志
			}
			return callback("", false, err)
		}

		if isEnd {
			// 流式响应结束，记录使用情况（如果没有错误）
			if !hasError {
				duration := int(time.Since(startTime).Milliseconds())
				response := fullResponse.String()

				// 创建使用记录（简化的Token计算）
				promptTokens := len(prompt) / 4       // 简化的Token估算
				completionTokens := len(response) / 4 // 简化的Token估算
				usage := s.modelService.CreateModelUsageFromResponse(
					userId,
					aiModel.Id,
					0, // 消息Id后续设置
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
			}

			return callback("", true, nil) // 通知用户回调流式响应结束
		}

		// 收集响应内容
		fullResponse.WriteString(chunk)

		// 调用用户回调
		return callback(chunk, false, nil)
	}

	// 调用AI客户端的流式生成回复
	return client.GenerateStreamResponse(prompt, aiHistory, internalCallback)
}

// handleCozeStreamResponse 处理Coze流式响应
func (s *AiService) handleCozeStreamResponse(aiModel *models.AIModel, prompt string, history []map[string]string, userId uint, callback func(chunk string, isEnd bool, err error) bool, startTime time.Time) error {
	// 创建Coze服务，传入模型配置
	cozeService, err := NewCozeService(aiModel)
	if err != nil {
		// 创建错误记录
		errorUsage := s.modelService.CreateModelUsageError(userId, aiModel.Id, prompt, "创建Coze服务失败: "+err.Error())
		if recordErr := s.modelService.RecordModelUsage(errorUsage); recordErr != nil {
			// 记录日志
		}
		return fmt.Errorf("创建Coze服务失败: %v", err)
	}

	// 转换历史消息格式为Coze格式
	cozeHistory := make([]*models.Message, 0, len(history))
	for _, msg := range history {
		cozeHistory = append(cozeHistory, &models.Message{
			Role:    msg["role"],
			Content: msg["content"],
		})
	}

	// 用于收集完整响应以记录使用情况
	var fullResponse strings.Builder
	var hasError bool
	var usageRecorded bool // 添加标志防止重复记录

	// 定义内部回调函数，包装用户回调并处理使用记录
	internalCallback := func(chunk string, isEnd bool, err error) bool {
		if err != nil {
			hasError = true
			// 创建错误记录（只记录一次）
			if !usageRecorded {
				errorUsage := s.modelService.CreateModelUsageError(userId, aiModel.Id, prompt, "Coze流式生成回复失败: "+err.Error())
				if recordErr := s.modelService.RecordModelUsage(errorUsage); recordErr != nil {
					// 记录日志
				}
				usageRecorded = true
			}
			return callback("", false, err)
		}

		if isEnd {
			// 流式响应结束，记录使用情况（如果没有错误且未记录过）
			if !hasError && !usageRecorded {
				duration := int(time.Since(startTime).Milliseconds())
				response := fullResponse.String()

				// 创建使用记录（简化的Token计算）
				promptTokens := len(prompt) / 4       // 简化的Token估算
				completionTokens := len(response) / 4 // 简化的Token估算
				usage := s.modelService.CreateModelUsageFromResponse(
					userId,
					aiModel.Id,
					0, // 消息Id后续设置
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
				usageRecorded = true
			}

			return callback("", true, nil) // 通知用户回调流式响应结束
		}

		// 收集响应内容
		fullResponse.WriteString(chunk)

		// 调用用户回调
		return callback(chunk, false, nil)
	}

	// 调用Coze服务的流式生成回复
	return cozeService.GenerateStreamResponse(prompt, cozeHistory, userId, internalCallback)
}
