package main

import (
	"chatbot-app/backend/config"
	"chatbot-app/backend/database"
	"chatbot-app/backend/models"
	"chatbot-app/backend/services"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载.env文件: %v", err)
	}

	// 获取配置
	cfg := config.GetConfig()
	cozeConfig := config.GetCozeConfig()

	fmt.Println("=== Coze智能体测试 ===")
	fmt.Printf("Coze API URL: %s\n", cozeConfig.APIURL)
	fmt.Printf("Coze Client ID: %s\n", cozeConfig.ClientID)
	fmt.Printf("Coze Bot ID: %s\n", cozeConfig.BotID)
	fmt.Printf("Coze Workflow ID: %s\n", cozeConfig.WorkFlowID)

	// 初始化数据库连接
	if err := database.InitMySQL(&cfg.Database); err != nil {
		log.Printf("警告: MySQL连接失败: %v", err)
	}

	// 初始化Redis连接
	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Printf("警告: Redis连接失败: %v", err)
	}

	// 测试Coze服务创建
	fmt.Println("\n=== 测试Coze服务创建 ===")
	cozeService, err := services.NewCozeService()
	if err != nil {
		log.Fatalf("创建Coze服务失败: %v", err)
	}

	fmt.Printf("Coze服务创建成功!\n")
	fmt.Printf("是否为工作流模式: %v\n", cozeService.IsWorkflowMode())
	fmt.Printf("Bot ID: %s\n", cozeService.GetBotID())
	fmt.Printf("Workflow ID: %s\n", cozeService.GetWorkflowID())

	// 测试会话创建
	fmt.Println("\n=== 测试会话创建 ===")
	conversationID, err := cozeService.GetConversationID()
	if err != nil {
		log.Printf("创建会话失败: %v", err)
	} else {
		fmt.Printf("会话创建成功，ID: %s\n", conversationID)
	}

	// 测试流式响应（如果配置了必要的参数）
	if cozeConfig.ClientID != "" && cozeConfig.BotID != "" {
		fmt.Println("\n=== 测试流式响应 ===")

		// 创建测试历史消息
		history := []*models.Message{
			{
				Role:    "user",
				Content: "你好",
			},
			{
				Role:    "assistant",
				Content: "你好！我是AI助手，有什么可以帮助您的吗？",
			},
		}

		prompt := "请简单介绍一下你自己"
		userID := uint(1)

		fmt.Printf("发送消息: %s\n", prompt)
		fmt.Println("AI回复:")

		// 流式回调函数
		streamCallback := func(chunk string, isEnd bool, err error) bool {
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return false
			}

			if isEnd {
				fmt.Println("\n--- 流式回复完成 ---")
				return true
			}

			// 实时输出文本块
			fmt.Print(chunk)
			return true
		}

		// 调用流式响应
		err = cozeService.GenerateStreamResponse(prompt, history, userID, streamCallback)
		if err != nil {
			log.Printf("流式响应测试失败: %v", err)
		}
	} else {
		fmt.Println("\n=== 跳过流式响应测试 ===")
		fmt.Println("原因: 缺少必要的Coze配置参数")
		fmt.Println("请确保设置了以下环境变量:")
		fmt.Println("- COZE_CLIENT_ID")
		fmt.Println("- COZE_PRIVATE_KEY 或 COZE_PRIVATE_KEY_FILE")
		fmt.Println("- COZE_PUBLIC_KEY_ID")
		fmt.Println("- COZE_BOT_ID")
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("如果看到此消息，说明Coze集成基本正常！")
}
