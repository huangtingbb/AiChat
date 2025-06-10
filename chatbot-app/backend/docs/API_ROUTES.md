# 聊天机器人 API 路由文档

## 基础信息
- **服务器地址**: `localhost:8080`
- **API 基础路径**: `/api`
- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Content-Type**: `application/json`
- **字符编码**: `UTF-8`

## API 路由总览

### 🔧 系统路由
| 方法 | 路径 | 描述 | 认证 | 状态 |
|------|------|------|------|------|
| GET | `/health` | 系统健康检查 | ❌ | ✅ |

### 👤 用户管理
| 方法 | 路径 | 描述 | 认证 | 状态 |
|------|------|------|------|------|
| POST | `/api/user/register` | 用户注册 | ❌ | ✅ |
| POST | `/api/user/login` | 用户登录 | ❌ | ✅ |
| GET | `/api/user/info` | 获取用户信息 | ✅ | ✅ |
| POST | `/api/user/logout` | 用户退出登录 | ✅ | ✅ |

### 💬 聊天管理
| 方法 | 路径 | 描述 | 认证 | 状态 |
|------|------|------|------|------|
| POST | `/api/chat` | 创建聊天会话 | ✅ | ✅ |
| GET | `/api/chat` | 获取聊天会话列表 | ✅ | ✅ |
| GET | `/api/chat/{id}/message` | 获取聊天消息列表 | ✅ | ✅ |
| POST | `/api/chat/{id}/message` | 发送聊天消息 | ✅ | ✅ |

### 🤖 AI 模型管理
| 方法 | 路径 | 描述 | 认证 | 状态 |
|------|------|------|------|------|
| GET | `/api/ai/model` | 获取可用模型列表 | ✅ | ✅ |
| POST | `/api/ai/model/set` | 设置默认模型 | ✅ | ✅ |
| POST | `/api/ai/model/option` | 设置模型参数 | ✅ | ✅ |
| GET | `/api/ai/usage` | 获取模型使用记录 | ✅ | ✅ |

## 🔐 认证说明
- **认证方式**: JWT Bearer Token
- **请求头**: `Authorization: Bearer <your_jwt_token>`
- **Token 获取**: 通过 `/api/user/login` 接口获得
- **Token 有效期**: 根据系统配置（通常24小时）

### 公开接口（无需认证）
- ✅ `GET /health` - 健康检查
- ✅ `POST /api/user/register` - 用户注册
- ✅ `POST /api/user/login` - 用户登录

## 📝 请求与响应示例

### 用户注册
```bash
POST /api/user/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

### 用户登录
```bash
POST /api/user/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

### 创建聊天会话
```bash
POST /api/chat
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "title": "新的聊天会话"
}
```

### 发送消息
```bash
POST /api/chat/1/message
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "content": "你好，请介绍一下你自己",
  "model_id": 1
}
```

## 📊 统一响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 具体数据内容
  },
  "timestamp": 1701234567
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "参数错误：用户名不能为空",
  "data": null,
  "timestamp": 1701234567
}
```

## ⚠️ 状态码说明

| 状态码 | 说明 | 场景 |
|--------|------|------|
| 0 | 成功 | 请求处理成功 |
| 400 | 参数错误 | 请求参数格式错误或缺失 |
| 401 | 未授权 | Token无效或已过期 |
| 403 | 禁止访问 | 权限不足 |
| 404 | 资源不存在 | 请求的资源不存在 |
| 500 | 服务器错误 | 服务器内部错误 |

## 🎯 使用流程示例

### 1. 用户注册并登录
```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","email":"test@example.com"}'

# 2. 用户登录
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### 2. 开始聊天
```bash
# 3. 创建聊天会话
curl -X POST http://localhost:8080/api/chat \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"我的第一个聊天"}'

# 4. 发送消息
curl -X POST http://localhost:8080/api/chat/1/message \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"content":"你好"}'
```

## 🛠️ 开发指南

### 路由命名规范
- ✅ **使用单数形式**: `user`, `chat`, `model`
- ❌ **避免复数形式**: `users`, `chats`, `models`
- ✅ **RESTful 设计**: 遵循标准的REST API设计原则
- ✅ **语义清晰**: 路径能够清楚表达资源和操作

### 错误处理
- 所有错误都返回统一的错误格式
- 包含详细的错误信息和状态码
- 客户端应根据状态码进行相应处理

### 分页说明
对于列表接口，后续可能会支持分页参数：
- `page`: 页码（从1开始）
- `limit`: 每页数量（默认20，最大100）

## 📚 相关文档
- [Swagger API 文档](http://localhost:8080/swagger/index.html)
- [项目 README](../../README.md)
- [数据库设计文档](../database/README.md)

---
**最后更新**: 2025年6月6日  
**文档版本**: 1.0.0 