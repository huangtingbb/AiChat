# ChatGPT类似应用

这是一个类似ChatGPT的聊天应用，具有用户认证和聊天功能。

## 技术栈

### 后端
- Go语言
- Gin框架 - Web框架
- GORM - ORM库
- MySQL - 数据库
- Redis - 缓存
- Elasticsearch - 全文搜索
- RabbitMQ - 消息队列
- JWT - 用户认证

### 前端
- Vue 3 - 前端框架
- Pinia - 状态管理
- Vue Router - 路由管理
- Element Plus - UI组件库
- Axios - HTTP客户端

## 功能特性

- 用户注册和登录
- JWT认证
- 聊天会话管理
- 消息历史记录
- AI模型集成

## 项目结构

```
chatbot-app/
├── backend/           # Go后端
│   ├── api/           # API控制器
│   ├── config/        # 配置
│   ├── database/      # 数据库连接
│   ├── middleware/    # 中间件
│   ├── models/        # 数据模型
│   ├── services/      # 业务逻辑
│   └── utils/         # 工具函数
├── frontend/          # Vue前端
│   ├── public/        # 静态资源
│   └── src/           # 源代码
│       ├── api/       # API请求
│       ├── assets/    # 资源文件
│       ├── components/# 组件
│       ├── router/    # 路由
│       ├── store/     # 状态管理
│       ├── utils/     # 工具函数
│       └── views/     # 页面
└── docs/              # 文档
```

## 运行指南

### 后端

1. 确保安装了Go环境和MySQL、Redis、Elasticsearch、RabbitMQ
2. 进入backend目录
3. 安装依赖：`go mod tidy`
4. 运行服务器：`go run main.go`

### 前端

1. 确保安装了Node.js环境
2. 进入frontend目录
3. 安装依赖：`npm install`
4. 开发模式运行：`npm run dev`
5. 构建生产版本：`npm run build`

## 环境要求

- Go 1.20+
- Node.js 16+
- MySQL 8.0+
- Redis 6.0+
- Elasticsearch 8.0+
- RabbitMQ 3.8+ 