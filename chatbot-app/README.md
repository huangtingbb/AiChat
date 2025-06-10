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
- `POST /api/chat/{id}/message` - 发送消息

### AI模型相关
- `GET /api/ai/model` - 获取可用模型列表
- `POST /api/ai/model/set` - 选择AI模型
- `GET /api/ai/usage` - 获取模型使用记录

### 用户相关
- `POST /api/user/login` - 用户登录
- `POST /api/user/register` - 用户注册
- `GET /api/user/info` - 获取用户信息
- `POST /api/user/logout` - 用户登出

## 待办事项
- [ ] 前端框架选择和搭建
- [ ] 用户界面设计
- [ ] AI模型集成测试
- [ ] 性能优化
- [ ] 单元测试编写
- [ ] 部署配置
