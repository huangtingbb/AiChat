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

### 🤖 Coze智能体集成完成
1. **Coze平台支持**
   - 完整集成Coze智能体平台
   - 支持Coze Bot对话模式和Workflow工作流模式
   - 实现JWT OAuth认证和Token缓存机制
   - 支持流式和非流式两种响应模式

2. **配置管理**
   - 添加`CozeConfig`配置结构体
   - 支持YAML配置文件和环境变量两种配置方式
   - 自动从Redis缓存Token，提高调用效率
   - 配置文件：`config/coze_config.yaml`

3. **服务架构**
   - 创建专用的`CozeService`处理智能体调用
   - 集成到AI工厂模式，支持多AI提供商统一管理
   - 实现`CozeClientAdapter`适配器，符合现有AI接口标准
   - 支持对话历史和上下文管理

4. **数据库支持**
   - 添加Coze模型到数据库迁移文件
   - 支持Coze Bot和Coze Workflow两种模型类型
   - 完善的使用记录和性能监控

5. **流式响应**
   - 支持Coze智能体的流式对话
   - 支持Coze工作流的流式执行
   - 实时事件处理和错误管理
   - 与现有SSE架构完美集成

### 📋 Coze集成完成清单
- ✅ **配置管理**：支持YAML和环境变量两种配置方式
- ✅ **JWT认证**：自动获取和缓存Access Token
- ✅ **Redis缓存**：Token缓存优化，减少API调用
- ✅ **服务架构**：专用CozeService和适配器模式
- ✅ **数据库集成**：自动添加Coze模型到数据库
- ✅ **流式响应**：完整的SSE流式对话支持
- ✅ **错误处理**：完善的错误处理和重试机制
- ✅ **使用统计**：Token消耗和性能监控
- ✅ **测试工具**：提供完整的测试和示例代码
- ✅ **文档完善**：详细的设置指南和API文档

### 🎯 自定义中文验证参数功能完成
1. **中文验证器实现**
   - 创建完整的自定义中文验证系统 `utils/validator.go`
   - 支持中文错误消息和字段名称映射
   - 自动将英文字段名转换为中文提示
   - 集成到Gin框架的验证器系统中

2. **自定义验证规则**
   - `mobile`: 中国手机号验证（11位数字，以1开头）
   - `strong_password`: 强密码验证（必须包含数字和字母）
   - `username_format`: 用户名格式验证（只能包含中文、英文、数字和下划线）
   - 支持动态注册新的验证规则

3. **中文错误消息系统**
   - 内置规则中文化：`required`→"不能为空"，`min`→"长度不能少于X个字符"等
   - 字段名称映射：`username`→"用户名"，`password`→"密码"，`email`→"邮箱"等
   - 友好的错误提示格式，提升用户体验

4. **控制器集成**
   - 修改用户注册和登录接口使用中文验证
   - 修改聊天相关接口使用中文验证
   - 统一的验证错误处理和日志记录
   - 支持获取单个错误或所有错误消息

5. **验证工具函数**
   - `GetValidationError()`: 获取第一个验证错误
   - `GetAllValidationErrors()`: 获取所有验证错误，用分号分隔
   - `ValidateStruct()`: 获取详细的验证错误列表
   - `RegisterCustomValidation()`: 动态注册自定义验证规则

### 📚 完整验证文档
- 创建详细的使用指南 `docs/VALIDATION.md`
- 包含使用方法、示例代码和最佳实践
- 支持自定义验证规则的扩展说明

### 🐛 修复闪烁问题和优化输出速度
1. **优化消息渲染**
   - 修复生成第二个回复时第一个回复闪烁的问题
   - 优化消息渲染key生成策略，只有内容真正变化时才更新
   - 移除不必要的`messageUpdateCount`全局更新依赖
   - 改进流式消息监听逻辑，减少不必要的重新渲染

2. **性能优化**
   - 添加CSS `contain: layout style` 提高渲染性能
   - 优化`will-change`属性使用，减少重排重绘
   - 实现更精确的流式消息状态跟踪
   - 优化滚动行为，减少频繁滚动触发

3. **流式输出速度控制**
   - 在后端添加智能延迟控制，根据内容长度动态调整输出速度
   - 延迟范围：50ms-300ms，让用户能够舒适地阅读内容
   - 避免过快的文本输出影响用户体验

4. **智能滚动优化**
   - 检测用户手动滚动行为，避免自动滚动干扰用户阅读
   - 只有在接近底部时才进行自动滚动
   - 添加"滚动到底部"按钮，方便用户快速回到最新消息
   - 发送消息和切换聊天时强制滚动到底部
   - 修复滚动逻辑冲突，防止程序滚动被误判为用户滚动
   - 添加调试日志和简化滚动函数，确保滚动功能正常工作
   - 优化流式输出时的实时滚动，确保内容输出过程中屏幕跟随滚动
   - 降低流式输出时用户滚动检测的敏感度，避免误判断自动滚动

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

# Coze智能体配置
COZE_API_URL=https://api.coze.cn
COZE_CLIENT_ID=your_coze_client_id
COZE_PRIVATE_KEY=your_coze_private_key
COZE_PRIVATE_KEY_FILE=/path/to/private_key.pem
COZE_PUBLIC_KEY_ID=your_coze_public_key_id
COZE_BOT_ID=your_coze_bot_id
COZE_WORKFLOW_ID=your_coze_workflow_id

# JWT密钥
JWT_SECRET=your_jwt_secret_key_here
```

### 必需配置项
- `ZHIPU_API_KEY`: 智谱AI的API密钥，在[智谱AI开放平台](https://open.bigmodel.cn/)获取
- `DB_PASSWORD`: 数据库密码
- `REDIS_PASSWORD`: Redis密码（如果有设置密码）
- `JWT_SECRET`: JWT签名密钥，建议使用长随机字符串

### Coze智能体配置项
- `COZE_CLIENT_ID`: Coze应用的Client ID，在[Coze开发者平台](https://www.coze.cn/)获取
- `COZE_PRIVATE_KEY`: Coze应用的私钥（PEM格式）
- `COZE_PRIVATE_KEY_FILE`: 私钥文件路径（与COZE_PRIVATE_KEY二选一）
- `COZE_PUBLIC_KEY_ID`: Coze应用的公钥ID
- `COZE_BOT_ID`: 要使用的Coze智能体ID
- `COZE_WORKFLOW_ID`: 要使用的Coze工作流ID（可选，用于工作流模式）

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
    userId, 
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
        selectedModel, prompt, history, userId, streamCallback
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
- **Coze智能体**: Coze平台的智能体对话模式
- **Coze工作流**: Coze平台的工作流模式，适合复杂的多步骤任务

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
- **多提供商支持**: 支持智谱AI、Coze、OpenAI等多种AI提供商（OpenAI待实现）

### Coze智能体使用示例
```go
// 创建Coze服务实例
cozeService, err := services.NewCozeService()
if err != nil {
    log.Fatal("创建Coze服务失败:", err)
}

// 使用流式对话
err = cozeService.GenerateStreamResponse(
    "请帮我分析一下这个问题", 
    historyMessages, 
    userID, 
    func(chunk string, isEnd bool, err error) bool {
        if err != nil {
            log.Printf("错误: %v", err)
            return false
        }
        
        if isEnd {
            fmt.Println("\n对话完成")
            return true
        }
        
        fmt.Print(chunk) // 实时输出
        return true
    },
)
```

### Coze配置文件示例
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
  workflow_id: "your_workflow_id"  # 可选，用于工作流模式
```

### Coze快速开始

1. **复制配置模板**：
   ```bash
   cp config/coze_config.example.yaml config/coze_config.yaml
   ```

2. **填写配置信息**：
   - 访问 [Coze开发者平台](https://www.coze.cn/) 获取凭证
   - 编辑 `config/coze_config.yaml` 填入实际值

3. **测试配置**：
   ```bash
   cd chatbot-app/backend
   go run examples/coze_example.go
   ```

4. **使用Coze模型**：
   - 启动项目后，在前端选择"Coze智能体"或"Coze工作流"模型
   - 开始与Coze智能体对话

详细设置指南请参考：[docs/COZE_SETUP.md](docs/COZE_SETUP.md)

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
3. 用户操作日志包含用户Id和IP地址
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

3. **消息Id生成优化**
   - 使用更精确的Id生成算法，避免Id冲突
   - 临时消息Id格式：`temp_user_timestamp_randomstring` 和 `temp_bot_timestamp_randomstring`

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

### 流式响应闪屏优化 (2024-12)

**问题描述：**
文本可以逐步出现，但页面在文本显示时一直闪屏，影响用户体验。

**优化内容：**

1. **滚动行为优化**
   - 添加滚动节流机制，避免频繁滚动操作
   - 使用smooth滚动替代瞬间跳转
   - 设置50ms的滚动延迟和100ms的防抖

2. **监听器性能优化**
   - 合并多个重复的watch监听器
   - 添加防抖机制减少触发频率
   - 优化日志输出，减少控制台性能影响

3. **Store更新频率控制**
   - 限制messageUpdateCount的更新频率（每100ms最多一次）
   - 减少不必要的强制响应式更新
   - 优化文本块处理的日志输出频率

4. **CSS渲染性能优化**
   - 添加`will-change`属性优化关键元素
   - 启用硬件加速（`transform: translateZ(0)`）
   - 使用CSS Containment（`contain: layout style`）
   - 优化动画参数，减少视觉干扰

**优化效果：**
- 消除了文本流式显示时的页面闪屏
- 提升了滚动的流畅度
- 减少了CPU和GPU的渲染负担
- 保持了流式文本的实时性

**核心优化代码：**

```javascript
// 滚动节流优化
let scrollTimer = null
const scrollToBottom = async () => {
  if (scrollTimer) clearTimeout(scrollTimer)
  scrollTimer = setTimeout(async () => {
    await nextTick()
    if (messagesContainer.value) {
      messagesContainer.value.scrollTo({
        top: messagesContainer.value.scrollHeight,
        behavior: 'smooth'
      })
    }
  }, 50)
}

// Store更新频率控制
if (!this._lastUpdateTime || Date.now() - this._lastUpdateTime > 100) {
  this.messageUpdateCount++
  this._lastUpdateTime = Date.now()
}
```

```css
/* CSS性能优化 */
.chat-messages {
  will-change: scroll-position;
  transform: translateZ(0);
}

.message-content {
  will-change: contents;
  contain: layout style;
}
```
