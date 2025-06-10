package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionOptions 聊天补全选项
type ChatCompletionOptions struct {
	MaxTokens        int     `json:"max_tokens,omitempty"`        // 生成文本的最大长度
	Temperature      float64 `json:"temperature,omitempty"`       // 温度参数，控制生成文本的随机性
	TopP             float64 `json:"top_p,omitempty"`             // 控制生成文本的多样性
	Stream           bool    `json:"stream,omitempty"`            // 是否启用流式响应
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`  // 影响模型不重复生成已经出现过的token的可能性
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"` // 影响模型不重复生成某些频繁出现的token的可能性
}

// AIClient AI客户端接口
type AIClient interface {
	// GenerateResponse 生成回复
	GenerateResponse(prompt string, history []Message) (string, error)
}

// ZhipuClient 智谱AI客户端
type ZhipuClient struct {
	BaseURL string
	ApiKey  string
	Model   string
	Options *ChatCompletionOptions
}

// NewZhipuClient 创建新的智谱AI客户端
func NewZhipuClient(apiKey, baseURL, model string, options *ChatCompletionOptions) *ZhipuClient {
	if options == nil {
		options = &ChatCompletionOptions{
			MaxTokens:   2048,
			Temperature: 0.75,
			TopP:        0.90,
			Stream:      false,
		}
	}

	return &ZhipuClient{
		BaseURL: baseURL,
		ApiKey:  apiKey,
		Model:   model,
		Options: options,
	}
}

// GenerateResponse 实现AIClient接口，生成回复
func (c *ZhipuClient) GenerateResponse(prompt string, history []Message) (string, error) {
	// 构建消息列表
	messages := make([]Message, 0, len(history)+1)
	messages = append(messages, history...)
	messages = append(messages, Message{
		Role:    "user",
		Content: prompt,
	})

	// 调用聊天补全API
	response, err := c.ChatCompletion(c.Model, messages, c.Options)
	if err != nil {
		return "", err
	}

	// 提取回复内容
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("未收到有效回复")
}

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	Stream           bool      `json:"stream,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// 生成JWT Token
func (c *ZhipuClient) generateToken() (string, error) {
	apiKey := c.ApiKey
	if len(apiKey) == 0 {
		return "", fmt.Errorf("ApiKey不能为空")
	}
	// 分割 API Key
	parts := make([]string, 2)
	parts = strings.Split(apiKey, ".")

	id := parts[0]
	secret := parts[1]

	// 获取当前时间（毫秒）
	now := time.Now().UnixNano() / int64(time.Millisecond)

	// 创建 payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"api_key":   id,
		"exp":       now + int64(3600*1000),
		"timestamp": now,
	})

	// 设置头部信息
	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"

	// 签名并获取完整 token
	return token.SignedString([]byte(secret))
}

// ChatCompletion 聊天补全
func (c *ZhipuClient) ChatCompletion(model string, messages []Message, options *ChatCompletionOptions) (*ChatResponse, error) {
	if model == "" {
		model = "glm-4" // 默认使用GLM-4模型
	}

	// 如果没有提供选项，使用默认选项
	if options == nil {
		options = &ChatCompletionOptions{
			MaxTokens:   2048,
			Temperature: 0.75,
			TopP:        0.90,
			Stream:      false,
		}
	}

	// 准备请求数据
	chatReq := ChatRequest{
		Model:            model,
		Messages:         messages,
		MaxTokens:        options.MaxTokens,
		Temperature:      options.Temperature,
		TopP:             options.TopP,
		Stream:           options.Stream,
		PresencePenalty:  options.PresencePenalty,
		FrequencyPenalty: options.FrequencyPenalty,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("请求数据序列化失败: %v", err)
	}

	//生成Token
	token, err := c.generateToken()
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %v", err)
	}

	// 构建HTTP请求
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// 发送请求
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &chatResp, nil
}
