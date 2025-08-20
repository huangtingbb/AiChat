# Coze智能体集成设置指南

## 概述

本项目已完整集成Coze智能体平台，支持两种模式：
- **Bot对话模式**：与Coze智能体进行自然对话
- **Workflow工作流模式**：执行复杂的多步骤任务流程

## 配置步骤

### 1. 获取Coze开发者账号

1. 访问 [Coze开发者平台](https://www.coze.cn/)
2. 注册并登录账号
3. 创建应用获取必要的凭证

### 2. 配置方式

#### 方式一：环境变量配置

在 `.env` 文件中添加以下配置：

```bash
# Coze智能体配置
COZE_API_URL=https://api.coze.cn
COZE_CLIENT_ID=your_coze_client_id
COZE_PRIVATE_KEY=your_coze_private_key
COZE_PUBLIC_KEY_ID=your_coze_public_key_id
COZE_BOT_ID=your_coze_bot_id
COZE_WORKFLOW_ID=your_coze_workflow_id  # 可选，用于工作流模式
```

#### 方式二：YAML配置文件

创建 `config/coze_config.yaml` 文件：

```yaml
coze:
  api_url: "https://api.coze.cn"
  client_id: "your_client_id"
  private_key: |
    -----BEGIN PRIVATE KEY-----
    your_private_key_content_here
    -----END PRIVATE KEY-----
  public_key_id: "your_public_key_id"
  bot_id: "your_bot_id"
  workflow_id: "your_workflow_id"  # 可选
```

### 3. 配置参数说明

| 参数 | 必需 | 说明 |
|------|------|------|
| `api_url` | 是 | Coze API地址，通常为 `https://api.coze.cn` |
| `client_id` | 是 | 应用的Client ID |
| `private_key` | 是 | JWT签名私钥（PEM格式） |
| `public_key_id` | 是 | 公钥ID |
| `bot_id` | 是 | 要使用的智能体ID |
| `workflow_id` | 否 | 工作流ID，设置后将使用工作流模式 |

## 使用方法

### 1. 在数据库中添加Coze模型

系统启动时会自动运行迁移，添加以下模型：

- **Coze智能体**：用于对话模式
- **Coze工作流**：用于工作流模式

### 2. 在前端选择Coze模型

用户可以在聊天界面中选择Coze相关的AI模型进行对话。

### 3. 编程方式调用

```go
// 创建Coze服务
cozeService, err := services.NewCozeService()
if err != nil {
    log.Fatal(err)
}

// 流式对话
err = cozeService.GenerateStreamResponse(
    "你好，请介绍一下你自己",
    historyMessages,
    userID,
    func(chunk string, isEnd bool, err error) bool {
        if err != nil {
            log.Printf("错误: %v", err)
            return false
        }
        if isEnd {
            fmt.Println("对话完成")
            return true
        }
        fmt.Print(chunk) // 实时输出
        return true
    },
)
```

## 功能特性

### ✅ 已实现功能

- JWT OAuth认证和Token缓存
- 流式对话响应
- 工作流执行
- 错误处理和重连
- Redis缓存优化
- 使用记录统计

### 🔄 工作原理

1. **Token管理**：自动获取和缓存JWT Token，有效期14分钟
2. **会话管理**：为每个对话创建独立的会话ID
3. **流式传输**：支持Server-Sent Events实时响应
4. **模式切换**：根据配置自动选择对话或工作流模式

### 📊 监控和日志

- Token获取和缓存状态
- API调用响应时间
- 错误日志和重试机制
- 使用量统计和记录

## 测试验证

运行测试程序验证配置：

```bash
cd chatbot-app/backend
go run examples/coze_example.go
```

正确配置后应该看到：
- Coze服务创建成功
- 会话ID创建成功
- 流式响应测试（如果配置完整）

## 故障排除

### 常见错误

1. **"未配置私钥"**
   - 检查 `COZE_PRIVATE_KEY` 或配置文件中的私钥设置
   - 确保私钥格式正确（PEM格式）

2. **"获取AccessToken失败"**
   - 检查Client ID和私钥是否匹配
   - 确认公钥ID是否正确

3. **"创建Coze会话失败"**
   - 检查Bot ID是否有效
   - 确认API URL是否可访问

### 调试技巧

1. 启用详细日志输出
2. 检查Redis连接状态
3. 验证网络连接和防火墙设置
4. 查看Coze开发者平台的API调用日志

## 安全建议

1. **私钥保护**：不要将私钥提交到版本控制系统
2. **环境隔离**：生产和测试环境使用不同的凭证
3. **访问控制**：限制API密钥的访问权限
4. **定期轮换**：定期更新API凭证

## 更新日志

- **v1.0.0**：初始版本，支持基本对话和工作流
- **v1.1.0**：添加Token缓存和错误重试
- **v1.2.0**：完善流式响应和监控
