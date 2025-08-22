package services

import (
	"chatbot-app/backend/models"
	"chatbot-app/backend/utils/coze"
	"fmt"
)

// CozeService Coze智能体服务
type CozeService struct {
	client *coze.Client
}

// NewCozeService 创建Coze服务实例
func NewCozeService() (*CozeService, error) {
	client, err := coze.New()
	if err != nil {
		return nil, fmt.Errorf("初始化Coze客户端失败: %v", err)
	}

	return &CozeService{
		client: client,
	}, nil
}

// GenerateResponse 生成回复（非流式）
func (s *CozeService) GenerateResponse(message string, history []*models.Message) (string, error) {
	// 如果配置了工作流ID，使用工作流模式
	if s.client.Config.WorkFlowID != "" {
		workflowResp, err := s.client.RunWorkflow(message)
		if err != nil {
			return "", fmt.Errorf("执行Coze工作流失败: %v", err)
		}

		// 从工作流响应中提取结果
		if workflowResp != nil {
			// 这里需要根据实际的coze-go库API来处理响应
			// 暂时返回一个简单的成功消息
			return fmt.Sprintf("Coze工作流执行成功，响应ID: %s", workflowResp.Data), nil
		}
		return "工作流执行完成，但未返回有效结果", nil
	}

	// 使用对话模式（暂不支持非流式）
	return "", fmt.Errorf("Coze非流式对话暂未实现，请使用流式模式")
}

// GenerateStreamResponse 生成流式回复
func (s *CozeService) GenerateStreamResponse(message string, history []*models.Message, userID uint, callback func(chunk string, isEnd bool, err error) bool) error {
	// 创建会话
	conversationID, err := s.client.CreateConversation()
	if err != nil {
		return fmt.Errorf("创建Coze会话失败: %v", err)
	}

	// 如果配置了工作流ID，使用工作流模式
	if s.client.Config.WorkFlowID != "" {
		return s.client.RunWorkflowStream(message, func(eventType string, data interface{}) {
			switch eventType {
			case "message_delta":
				if dataMap, ok := data.(map[string]string); ok {
					if content, exists := dataMap["content"]; exists {
						if !callback(content, false, nil) {
							return
						}
					}
				}
			case "workflow_complated", "workflow_end":
				callback("", true, nil)
			case "workflow_error":
				if dataMap, ok := data.(map[string]string); ok {
					if errorContent, exists := dataMap["content"]; exists {
						callback("", false, fmt.Errorf("工作流错误: %s", errorContent))
						return
					}
				}
				callback("", false, fmt.Errorf("工作流执行失败"))
			}
		})
	}

	// 使用对话流式模式
	return s.client.SendMessageStreamWithCallback(
		conversationID,
		userID,
		history,
		func(eventType string, data interface{}) {
			switch eventType {
			case "message_delta":
				if dataMap, ok := data.(map[string]interface{}); ok {
					if content, exists := dataMap["content"]; exists {
						if contentStr, ok := content.(string); ok {
							if !callback(contentStr, false, nil) {
								return
							}
						}
					}
				}
			case "chat_completed", "conversation_end":
				callback("", true, nil)
			case "chat_failed":
				if dataMap, ok := data.(map[string]interface{}); ok {
					if errorMsg, exists := dataMap["error_msg"]; exists {
						callback("", false, fmt.Errorf("对话失败: %v", errorMsg))
						return
					}
				}
				callback("", false, fmt.Errorf("对话执行失败"))
			}
		},
	)
}

// GetConversationID 创建新的会话ID
func (s *CozeService) GetConversationID() (string, error) {
	return s.client.CreateConversation()
}

// IsWorkflowMode 检查是否为工作流模式
func (s *CozeService) IsWorkflowMode() bool {
	return s.client.Config.WorkFlowID != ""
}

// GetBotID 获取机器人ID
func (s *CozeService) GetBotID() string {
	return s.client.Config.BotID
}

// GetWorkflowID 获取工作流ID
func (s *CozeService) GetWorkflowID() string {
	return s.client.Config.WorkFlowID
}
