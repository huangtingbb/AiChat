package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	// GenerateStreamResponse 生成流式回复
	GenerateStreamResponse(prompt string, history []Message, callback func(chunk string, isEnd bool, err error) bool) error
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
	Id      string `json:"id"`
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

// StreamResponse 流式响应结构
type StreamResponse struct {
	Id      string `json:"id"`
	Created int64  `json:"created"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

// GenerateStreamResponse 实现流式回复
func (c *ZhipuClient) GenerateStreamResponse(prompt string, history []Message, callback func(chunk string, isEnd bool, err error) bool) error {
	// 构建消息列表
	messages := make([]Message, 0, len(history)+1)
	messages = append(messages, history...)
	messages = append(messages, Message{
		Role:    "user",
		Content: prompt,
	})

	// 创建流式请求选项
	streamOptions := &ChatCompletionOptions{
		MaxTokens:        c.Options.MaxTokens,
		Temperature:      c.Options.Temperature,
		TopP:             c.Options.TopP,
		PresencePenalty:  c.Options.PresencePenalty,
		FrequencyPenalty: c.Options.FrequencyPenalty,
		Stream:           true, // 启用流式响应
	}

	// 调用流式聊天补全API
	return c.ChatCompletionStream(c.Model, messages, streamOptions, callback)
}

// ChatCompletionStream 流式聊天补全
func (c *ZhipuClient) ChatCompletionStream(model string, messages []Message, options *ChatCompletionOptions, callback func(chunk string, isEnd bool, err error) bool) error {
	if model == "" {
		model = "glm-4" // 默认使用GLM-4模型
	}

	// 如果没有提供选项，使用默认选项
	if options == nil {
		options = &ChatCompletionOptions{
			MaxTokens:   2048,
			Temperature: 0.75,
			TopP:        0.90,
			Stream:      true,
		}
	}

	// 准备请求数据
	chatReq := ChatRequest{
		Model:            model,
		Messages:         messages,
		MaxTokens:        options.MaxTokens,
		Temperature:      options.Temperature,
		TopP:             options.TopP,
		Stream:           true, // 强制设置为流式
		PresencePenalty:  options.PresencePenalty,
		FrequencyPenalty: options.FrequencyPenalty,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return fmt.Errorf("请求数据序列化失败: %v", err)
	}

	// 生成Token
	token, err := c.generateToken()
	if err != nil {
		return fmt.Errorf("生成Token失败: %v", err)
	}

	// 构建HTTP请求
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// 发送请求
	client := &http.Client{Timeout: time.Second * 300} // 增加超时时间以支持流式响应
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败，状态码: %d，响应: %s", resp.StatusCode, string(respBody))
	}

	// 处理流式响应
	return c.parseStreamResponse(resp.Body, callback)
}

// parseStreamResponse 解析流式响应
func (c *ZhipuClient) parseStreamResponse(body io.Reader, callback func(chunk string, isEnd bool, err error) bool) error {
	const (
		dataPrefix = "data: "
		doneMarker = "[DONE]"
	)

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// 处理数据行
		if strings.HasPrefix(line, dataPrefix) {
			data := strings.TrimPrefix(line, dataPrefix)

			// 检查是否是结束标记
			if data == doneMarker {
				// 通知回调函数流式响应结束
				callback("", true, nil)
				break
			}

			// 解析JSON数据
			var streamResp StreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				// 发送错误到回调函数
				if !callback("", false, fmt.Errorf("解析流式响应失败: %v", err)) {
					break
				}
				continue
			}

			// 提取内容并发送到回调函数
			if len(streamResp.Choices) > 0 {
				content := streamResp.Choices[0].Delta.Content
				if content != "" {
					// 添加延迟控制输出速度，让用户能看清楚内容
					// 根据内容长度调整延迟时间，避免过快的输出
					delay := time.Duration(len(content)*50) * time.Millisecond
					if delay > 300*time.Millisecond {
						delay = 300 * time.Millisecond // 最大延迟300ms
					}
					if delay < 50*time.Millisecond {
						delay = 50 * time.Millisecond // 最小延迟50ms
					}
					time.Sleep(50 * time.Millisecond)

					// 如果回调函数返回false，停止处理
					if !callback(content, false, nil) {
						break
					}
				}

				// 检查是否完成
				if streamResp.Choices[0].FinishReason != nil && *streamResp.Choices[0].FinishReason != "" {
					callback("", true, nil)
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流式响应失败: %v", err)
	}

	return nil
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
