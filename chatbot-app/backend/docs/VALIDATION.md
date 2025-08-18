# 自定义中文验证参数使用指南

## 概述

本项目实现了一套完整的自定义中文验证系统，支持中文错误消息和自定义验证规则。

## 功能特性

### 1. 中文错误消息
- 自动将验证错误翻译为中文
- 支持字段名称的中文映射
- 提供友好的错误提示

### 2. 自定义验证规则
- `mobile`: 中国手机号验证
- `strong_password`: 强密码验证（必须包含数字和字母）
- `username_format`: 用户名格式验证（只能包含中文、英文、数字和下划线）

### 3. 内置规则中文化
- `required`: 不能为空
- `min`: 长度不能少于X个字符
- `max`: 长度不能超过X个字符
- `email`: 格式不正确
- 等等...

## 使用方法

### 1. 在结构体中定义验证规则

```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50,username_format" validate:"required,min=3,max=50,username_format"`
    Password string `json:"password" binding:"required,min=6,max=50,strong_password" validate:"required,min=6,max=50,strong_password"`
    Email    string `json:"email" binding:"required,email" validate:"required,email"`
    Phone    string `json:"phone" binding:"omitempty,mobile" validate:"omitempty,mobile"`
}
```

### 2. 在控制器中使用验证

```go
func (controller *UserController) Register(c *gin.Context) {
    var req RegisterRequest
    
    // Gin 内置验证
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.InvalidParams(c, err.Error())
        return
    }
    
    // 自定义验证器验证
    if validationError := utils.GetValidationError(req); validationError != "" {
        utils.InvalidParams(c, validationError)
        return
    }
    
    // 处理业务逻辑...
}
```

## 验证器工具函数

### 1. `GetValidationError(obj interface{}) string`
获取第一个验证错误消息

### 2. `GetAllValidationErrors(obj interface{}) string`
获取所有验证错误消息，用分号分隔

### 3. `ValidateStruct(obj interface{}) []ValidateError`
获取详细的验证错误列表

## 字段名称映射

系统会自动将英文字段名映射为中文：

- `username` → 用户名
- `password` → 密码
- `email` → 邮箱
- `content` → 内容
- `model_id` → 模型ID
- `title` → 标题
- `phone` → 手机号
- `nickname` → 昵称
- `avatar` → 头像

## 错误消息示例

### 成功的验证
```json
{
    "code": 0,
    "message": "注册成功",
    "data": {...}
}
```

### 验证失败
```json
{
    "code": -1,
    "message": "用户名只能包含中文、英文、数字和下划线",
    "data": null
}
```

## 添加自定义验证规则

### 1. 注册验证规则

```go
// 在 utils/validator.go 的 initCustomValidations() 函数中添加
RegisterCustomValidation("custom_rule", func(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    // 自定义验证逻辑
    return true // 或 false
}, "{0}自定义错误消息")
```

### 2. 在结构体中使用

```go
type MyStruct struct {
    Field string `json:"field" validate:"required,custom_rule"`
}
```

## 注意事项

1. 确保在 `main.go` 中调用了 `utils.InitValidator()`
2. 结构体字段需要同时使用 `binding` 和 `validate` 标签
3. 自定义验证规则需要在应用启动时注册
4. 字段名称映射可以通过修改 `getFieldName` 函数来扩展

## 完整示例

参考 `controller/user_controller.go` 中的 `Register` 和 `Login` 方法，以及 `controller/chat_controller.go` 中的相关方法。
