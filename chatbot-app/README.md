# 聊天机器人项目

## 项目概述
这是一个基于Go后端和Vue3前端的聊天机器人应用，支持多种AI模型对话。

## 技术栈

### 后端
- **语言**: Go 1.24
- **框架**: Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **日志**: Logrus
- **认证**: JWT
- **文档**: Swagger

### 前端
- **框架**: Vue3
- **构建工具**: Vite
- **UI组件**: 待选择

## 项目结构

```
chatbot-app/
├── backend/                 # 后端代码
│   ├── config/             # 配置管理
│   ├── controller/         # 控制器层
│   ├── database/           # 数据库连接
│   ├── docs/               # API文档
│   ├── logs/               # 日志文件
│   ├── middleware/         # 中间件
│   ├── migrations/         # 迁移文件
│   ├── models/             # 数据模型
│   ├── router/             # 路由配置
│   ├── services/           # 业务逻辑层
│   ├── utils/              # 工具包
│   │   ├── ai_client.go    # AI客户端工具
│   │   ├── ai_factory.go   # AI客户端工厂
│   │   ├── jwt.go          # JWT工具
│   │   ├── logger.go       # 日志工具
│   │   ├── password.go     # 密码工具
│   │   └── response.go     # 响应工具
│   ├── .env                # 环境变量
│   └── main.go             # 程序入口
├── frontend/               # 前端代码
└── README.md               # 项目说明
```

## 最新更新 (2025-01-27)

### 🚀 流式AI回复功能完成
1. **流式响应实现**
   - 扩展`AIClient`接口，添加`GenerateStreamResponse`方法
   - 实现智谱AI的Server-Sent Events (SSE) 流式响应
   - 支持实时接收和处理AI生成的文本块
   - 完整的错误处理和连接管理

2. **流式回调机制**
   - 统一的流式回调函数签名：`func(chunk string, isEnd bool, err error) bool`
   - 支持实时文本输出和完整响应收集
   - 客户端可通过返回`false`中断流式传输
   - 自动记录使用情况和Token消耗

3. **智谱AI GLM-Z1模型支持**
   - 支持智谱AI最新的推理模型
   - 完整的请求头设置（Accept: text/event-stream）
   - 正确解析SSE格式的响应数据
   - JWT Token自动生成和管理

### AI服务架构重构
1. **工具代码重构**
   - 将AI客户端相关代码从`services/ai_clients`移动到`utils`目录
   - 创建统一的AI客户端接口`AIClient`
   - 实现智谱AI客户端`ZhipuClient`
   - 添加AI客户端工厂`AIClientFactory`

2. **服务层优化**
   - 重写`AiService`，使用utils中的工具
   - 简化AI模型服务`AIModelService`
   - 移除重复的代码和方法
   - 统一消息格式转换

3. **控制器完善**
   - 添加`SelectModel`方法用于选择AI模型
   - 添加`GetModelUsageHandler`方法获取使用记录
   - 完善API文档注释

4. **代码结构优化**
   - 删除旧的`ai_clients`目录
   - 统一使用utils中的工具函数
   - 改进错误处理和日志记录

### 技术改进
- **模块化设计**: 将AI相关工具集中到utils目录，便于维护和扩展
- **接口抽象**: 使用接口设计，支持多种AI提供商
- **工厂模式**: 使用工厂模式创建AI客户端，便于扩展新的AI服务
- **配置管理**: 支持从环境变量读取API密钥和基础URL
- **流式处理**: 支持Server-Sent Events流式响应，提升用户体验

### 环境变量配置
```bash
# 智谱AI配置
ZHIPU_API_KEY=your_api_key_here
ZHIPU_BASE_URL=https://open.bigmodel.cn/api/paas/v4
```

## 环境变量详细配置

### 配置文件位置
在`backend`目录下创建`.env`文件，包含以下配置：

```bash
# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=chatbot

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password_here
REDIS_DB=0

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=logs/app.log
LOG_MAX_SIZE=100
LOG_MAX_AGE=30
LOG_COMPRESS=true

# AI配置
# 智谱AI配置（必填）
ZHIPU_API_KEY=your_zhipu_api_key_here
ZHIPU_BASE_URL=https://open.bigmodel.cn/api/paas/v4

# OpenAI配置（预留）
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_BASE_URL=https://api.openai.com/v1

# JWT密钥
JWT_SECRET=your_jwt_secret_key_here
```

### 必需配置项
- `ZHIPU_API_KEY`: 智谱AI的API密钥，在[智谱AI开放平台](https://open.bigmodel.cn/)获取
- `DB_PASSWORD`: 数据库密码
- `REDIS_PASSWORD`: Redis密码（如果有设置密码）
- `JWT_SECRET`: JWT签名密钥，建议使用长随机字符串

## 流式AI回复使用指南

### 功能特性
- ✅ **实时响应**: 支持Server-Sent Events (SSE)流式输出
- ✅ **智谱AI集成**: 完整支持智谱AI GLM系列模型
- ✅ **错误处理**: 完善的错误处理和重连机制
- ✅ **使用记录**: 自动记录Token消耗和响应时间
- ✅ **连接控制**: 支持客户端主动中断传输

### 使用示例
```go
// 创建AI服务实例
aiService := services.NewAiService()

// 获取AI模型
aiModel, _ := aiModelService.GetModelByName("glm-4-plus")

// 定义流式回调函数
streamCallback := func(chunk string, isEnd bool, err error) bool {
    if err != nil {
        log.Printf("错误: %v", err)
        return false // 停止接收
    }
    
    if isEnd {
        fmt.Println("\n流式回复完成")
        return true
    }
    
    // 实时输出文本块
    fmt.Print(chunk)
    return true // 继续接收
}

// 调用流式AI回复
err := aiService.GenerateStreamResponse(
    aiModel, 
    "你的问题", 
    history, 
    userID, 
    streamCallback
)
```

### 运行示例
```bash
# 设置环境变量
export ZHIPU_API_KEY="your_api_key_here"

# 运行流式回复示例
cd chatbot-app/backend
go run examples/stream_example.go
```

### API集成
流式AI回复可以轻松集成到Web API中：

```go
// 在控制器中使用流式响应
func (controller *ChatController) StreamMessage(c *gin.Context) {
    // 设置SSE响应头
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    
    // 定义流式回调
    streamCallback := func(chunk string, isEnd bool, err error) bool {
        if err != nil {
            c.SSEvent("error", err.Error())
            return false
        }
        
        if isEnd {
            c.SSEvent("end", "")
            return true
        }
        
        c.SSEvent("data", chunk)
        c.Writer.Flush()
        return true
    }
    
    // 调用流式AI服务
    err := controller.aiService.GenerateStreamResponse(
        selectedModel, prompt, history, userID, streamCallback
    )
    
    if err != nil {
        c.SSEvent("error", err.Error())
    }
}
```

### 支持的AI模型
- **GLM-4**: 智谱AI通用对话模型
- **GLM-4-Plus**: 智谱AI增强版对话模型
- **GLM-4-Air**: 智谱AI轻量版模型
- **GLM-Z1**: 智谱AI最新推理模型（推荐）

### 性能优化
- **连接池**: 使用HTTP连接池复用连接
- **超时控制**: 设置合理的超时时间（5分钟）
- **内存管理**: 流式处理减少内存占用
- **并发支持**: 支持多个并发流式请求

### AI模型配置说明
- **数据库驱动**: 系统从数据库`ai_model`表读取AI模型配置
- **动态切换**: 支持运行时动态切换AI模型
- **API Key管理**: 模型的API Key从环境变量读取，提高安全性
- **Base URL配置**: 如果模型表中没有配置URL，则使用环境变量中的默认URL
- **多提供商支持**: 支持智谱AI、OpenAI等多种AI提供商（OpenAI待实现）

## 运行说明

### 后端启动
```bash
cd chatbot-app/backend
go mod tidy
go run main.go
```

### 日志配置
日志文件默认保存在 `logs/app.log`，支持以下配置：
- 日志级别：debug, info, warn, error
- 日志格式：json, text
- 文件大小限制：100MB
- 保留天数：30天
- 自动压缩旧日志

### API文档
```bash
swag init -g chatbot-app/backend/main.go -o chatbot-app/backend/docs
```
启动后访问：http://localhost:8080/swagger/index.html

## 开发规范

### 日志记录规范
1. 使用结构化日志，包含相关上下文信息
2. 错误日志必须包含错误详情和相关参数
3. 用户操作日志包含用户ID和IP地址
4. 敏感信息不记录到日志中

### 错误处理规范
1. 统一使用utils包中的响应函数
2. 记录详细的错误日志
3. 返回用户友好的错误信息

## API接口

### 聊天相关
- `POST /api/chat` - 创建聊天会话
- `GET /api/chat` - 获取聊天会话列表
- `GET /api/chat/{id}/message` - 获取聊天消息
- `POST /api/chat/{id}/message` - 发送消息（流式响应）

#### 流式聊天API说明
`POST /api/chat/{id}/message` 接口现在使用Server-Sent Events (SSE)流式响应，提供实时的AI回复体验。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/chat/1/message \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_jwt_token" \
  -H "Accept: text/event-stream" \
  -d '{"content": "你好，请介绍一下你自己", "model_id": 1}'
```

**响应事件类型**：
- `user_message` - 用户消息已保存
- `stream_start` - 开始AI流式回复
- `stream_chunk` - AI回复文本块
- `stream_end` - AI回复完成
- `error` - 错误信息

**响应示例**：
```
data: {"type":"user_message","message":{"id":1,"content":"你好"}}

data: {"type":"stream_start","model":"智谱GLM-4-Plus"}

data: {"type":"stream_chunk","text":"你好！"}

data: {"type":"stream_chunk","text":"我是智谱AI助手"}

data: {"type":"stream_end","message_id":2,"full_text":"你好！我是智谱AI助手..."}
```

### 前端流式响应集成

前端已完全支持SSE流式响应，提供实时的用户体验：

#### 功能特性
- ✅ **实时渲染**: 文本内容实时更新和Markdown渲染
- ✅ **流式指示器**: 带动画的"正在输入"提示
- ✅ **自动滚动**: 消息自动滚动到底部
- ✅ **错误处理**: 完善的错误提示和重连机制
- ✅ **连接控制**: 支持中断流式传输

#### 技术实现
- **Vue3 + Pinia**: 响应式状态管理
- **Fetch Streams API**: 原生流式数据处理
- **marked.js**: 实时Markdown渲染
- **Element Plus**: 现代化UI组件

#### 使用方法
```javascript
// 在store中发送流式消息
await chatStore.sendUserMessage(content, modelId)

// 自动处理以下SSE事件:
// - user_message: 用户消息保存确认
// - stream_start: 开始AI回复
// - stream_chunk: 实时文本块
// - stream_end: 回复完成
// - error: 错误处理
```

#### 测试工具
项目提供了独立的测试页面：
```bash
# 访问测试页面
open chatbot-app/frontend/src/test-stream.html
```

该测试页面可以：
- 直接测试SSE流式响应
- 调试不同的消息内容
- 查看详细的事件日志
- 验证JWT Token有效性

### AI模型相关
- `GET /api/ai/model` - 获取可用模型列表
- `POST /api/ai/model/set` - 选择AI模型
- `GET /api/ai/usage` - 获取模型使用记录

### 用户相关
- `POST /api/user/login` - 用户登录
- `POST /api/user/register` - 用户注册
- `GET /api/user/info` - 获取用户信息
- `POST /api/user/logout` - 用户登出

## 🎉 完整流式AI聊天系统已完成

### 完整功能链路
1. **前端发起**: Vue3聊天界面发送消息
2. **后端接收**: Gin处理HTTP请求，验证JWT
3. **AI调用**: 通过智谱AI客户端发起流式请求
4. **流式传输**: Server-Sent Events实时传输文本块
5. **前端渲染**: 实时Markdown渲染和UI更新
6. **数据存储**: 完整对话保存到MySQL数据库

### 技术栈总览
```
前端: Vue3 + Pinia + Element Plus + marked.js
后端: Go + Gin + GORM + JWT
数据库: MySQL 8.0 + Redis
AI服务: 智谱AI GLM系列模型
通信: SSE (Server-Sent Events) + REST API
```

### 快速开始
```bash
# 1. 启动后端服务
cd chatbot-app/backend
export ZHIPU_API_KEY="your_api_key_here"
go run main.go

# 2. 启动前端服务
cd chatbot-app/frontend
npm install
npm run dev

# 3. 访问应用
open http://localhost:3000
```

### 体验流式对话
1. 注册/登录账户
2. 创建新对话
3. 发送消息并实时观看AI回复
4. 支持Markdown格式、代码高亮、表格等
5. 自动保存对话历史

## 待办事项
- [ ] 前端框架选择和搭建
- [ ] 用户界面设计
- [ ] AI模型集成测试
- [ ] 性能优化
- [ ] 单元测试编写
- [ ] 部署配置

## 最近修复

### SSE消息渲染问题修复 (2024-12)

**问题描述：**
使用SSE协议请求后端发送消息接口时，消息可以正常返回，前端也能接收到消息，但是没有正确渲染到页面。

**修复内容：**

1. **响应式数据处理优化**
   - 移除了不必要的`reactive()`包装，Pinia已自动处理响应式
   - 修改消息数组更新方式，确保Vue能正确检测到变化
   - 使用直接对象替换而不是属性修改来触发响应式更新

2. **SSE数据解析改进**
   - 修复了SSE数据解析逻辑，正确处理Gin框架的SSE格式
   - 改进了data字段的提取方式：`data: jsondata` -> 正确trim处理

3. **消息ID生成优化**
   - 使用更精确的ID生成算法，避免ID冲突
   - 临时消息ID格式：`temp_user_timestamp_randomstring` 和 `temp_bot_timestamp_randomstring`

4. **Vue监听器优化**
   - 添加了`flush: 'post'`选项确保DOM更新后执行
   - 增加了对store中messages数组的直接监听
   - 改进了日志输出，便于调试

**核心修复代码：**

```javascript
// handleSSEEvent中的消息更新方式
case 'stream_chunk':
  const currentMessage = this.messages[botMessageIndex]
  const updatedMessage = {
    ...currentMessage,
    content: currentMessage.content + eventData.text,
    isStreaming: true
  }
  // 直接替换消息对象来触发响应式更新
  this.messages[botMessageIndex] = updatedMessage
  this.messageUpdateCount++
  break
```

**测试验证：**
- SSE消息能够正确解析和渲染
- 流式文本能够实时显示在界面上
- 消息滚动和UI更新正常工作
