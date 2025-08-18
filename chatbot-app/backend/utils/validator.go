package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

// ValidateError 验证错误结构
type ValidateError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// InitValidator 初始化验证器
func InitValidator() error {
	// 创建中文翻译器
	zhLocal := zh.New()
	uni := ut.New(zhLocal, zhLocal)
	var ok bool
	trans, ok = uni.GetTranslator("zh")
	if !ok {
		return fmt.Errorf("创建中文翻译器失败")
	}

	// 获取gin的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v

		// 注册中文翻译
		if err := zh_translations.RegisterDefaultTranslations(validate, trans); err != nil {
			return fmt.Errorf("注册中文翻译失败: %w", err)
		}

		// 注册自定义标签名称翻译
		registerCustomTranslations()

		// 注册字段名称获取函数
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			// 优先使用json标签作为字段名
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name == "" {
				name = fld.Name
			}
			return name
		})

		// 初始化自定义验证规则
		initCustomValidations()
	}

	return nil
}

// registerCustomTranslations 注册自定义翻译
func registerCustomTranslations() {
	// 注册自定义错误消息
	translations := []struct {
		tag         string
		translation string
	}{
		{"required", "{0}不能为空"},
		{"min", "{0}长度不能少于{1}个字符"},
		{"max", "{0}长度不能超过{1}个字符"},
		{"email", "{0}格式不正确"},
		{"len", "{0}长度必须为{1}个字符"},
		{"numeric", "{0}必须是数字"},
		{"alphanum", "{0}只能包含字母和数字"},
		{"alpha", "{0}只能包含字母"},
		{"gte", "{0}必须大于等于{1}"},
		{"lte", "{0}必须小于等于{1}"},
		{"gt", "{0}必须大于{1}"},
		{"lt", "{0}必须小于{1}"},
		{"oneof", "{0}必须是以下值之一: {1}"},
	}

	for _, t := range translations {
		registerTranslation(t.tag, t.translation)
	}
}

// registerTranslation 注册单个翻译
func registerTranslation(tag, translation string) {
	validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, getFieldName(fe), fe.Param())
		return t
	})
}

// getFieldName 获取字段的中文名称
func getFieldName(fe validator.FieldError) string {
	fieldName := fe.Field()

	// 字段名称映射
	fieldMap := map[string]string{
		"username": "用户名",
		"password": "密码",
		"email":    "邮箱",
		"content":  "内容",
		"model_id": "模型ID",
		"title":    "标题",
		"phone":    "手机号",
		"nickname": "昵称",
		"avatar":   "头像",
	}

	if chineseName, exists := fieldMap[strings.ToLower(fieldName)]; exists {
		return chineseName
	}

	return fieldName
}

// ValidateStruct 验证结构体并返回中文错误消息
func ValidateStruct(obj interface{}) []ValidateError {
	var errors []ValidateError

	err := validate.Struct(obj)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidateError{
				Field:   err.Field(),
				Message: err.Translate(trans),
			})
		}
	}

	return errors
}

// GetValidationError 获取第一个验证错误消息
func GetValidationError(obj interface{}) string {
	errors := ValidateStruct(obj)
	if len(errors) > 0 {
		return errors[0].Message
	}
	return ""
}

// GetAllValidationErrors 获取所有验证错误消息，用分号分隔
func GetAllValidationErrors(obj interface{}) string {
	errors := ValidateStruct(obj)
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Message)
	}

	return strings.Join(messages, "; ")
}

// ValidateStructWithTagMessages 验证结构体并使用结构体标签中的自定义错误消息
func ValidateStructWithTagMessages(obj interface{}) []ValidateError {
	var errors []ValidateError

	err := validate.Struct(obj)
	if err != nil {
		// 获取结构体类型信息
		val := reflect.ValueOf(obj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		structType := val.Type()

		for _, err := range err.(validator.ValidationErrors) {
			// validator返回的字段名可能是JSON标签名，需要查找实际的结构体字段
			fieldName := err.Field()

			// 先尝试直接按字段名查找
			field, found := structType.FieldByName(fieldName)

			// 如果没找到，尝试按JSON标签名查找
			if !found {
				field, found = findStructFieldByJSONName(structType, fieldName)
			}

			var message string
			if found {
				// 查找对应验证规则的自定义消息
				msgTag := "msg_" + err.Tag()
				if customMsg := field.Tag.Get(msgTag); customMsg != "" {
					message = customMsg
				} else {
					// 使用默认翻译
					message = err.Translate(trans)
				}
			} else {
				// 使用默认翻译
				message = err.Translate(trans)
			}

			errors = append(errors, ValidateError{
				Field:   fieldName,
				Message: message,
			})
		}
	}

	return errors
}

// GetValidationErrorWithTagMessages 获取第一个验证错误消息（使用结构体标签）
func GetValidationErrorWithTagMessages(obj interface{}) string {
	errors := ValidateStructWithTagMessages(obj)
	if len(errors) > 0 {
		return errors[0].Message
	}
	return ""
}

// GetAllValidationErrorsWithTagMessages 获取所有验证错误消息（使用结构体标签），用分号分隔
func GetAllValidationErrorsWithTagMessages(obj interface{}) string {
	errors := ValidateStructWithTagMessages(obj)
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Message)
	}

	return strings.Join(messages, "; ")
}

// RegisterFieldName 注册字段中文名称
func RegisterFieldName(fieldName, chineseName string) {
	// 这个函数可以用来动态注册字段名称
	// 可以在控制器中调用来注册特定的字段名称
}

// CustomValidationFunc 自定义验证函数类型
type CustomValidationFunc func(fl validator.FieldLevel) bool

// RegisterCustomValidation 注册自定义验证规则
func RegisterCustomValidation(tag string, fn CustomValidationFunc, message string) error {
	if validate == nil {
		return fmt.Errorf("验证器未初始化")
	}

	// 注册验证函数
	if err := validate.RegisterValidation(tag, validator.Func(fn)); err != nil {
		return err
	}

	// 注册翻译
	registerTranslation(tag, message)

	return nil
}

// 初始化时注册一些常用的自定义验证规则
func initCustomValidations() {
	// 手机号验证
	RegisterCustomValidation("mobile", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		// 简单的中国手机号验证规则
		if len(phone) != 11 {
			return false
		}
		return strings.HasPrefix(phone, "1")
	}, "{0}格式不正确")

	// 强密码验证（至少包含数字和字母）
	RegisterCustomValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		hasNumber := strings.ContainsAny(password, "0123456789")
		hasLetter := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		return hasNumber && hasLetter
	}, "{0}必须包含数字和字母")

	// 用户名验证（只能包含中文、英文、数字和下划线）
	RegisterCustomValidation("username_format", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		for _, r := range username {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '_' ||
				(r >= 0x4e00 && r <= 0x9fff)) { // 中文字符范围
				return false
			}
		}
		return true
	}, "{0}只能包含中文、英文、数字和下划线")
}

// findStructFieldByJSONName 根据JSON标签名查找结构体字段
func findStructFieldByJSONName(structType reflect.Type, jsonName string) (reflect.StructField, bool) {
	// 遍历所有字段
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// 获取json标签
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		// 解析json标签，去掉选项（如omitempty）
		jsonFieldName := strings.SplitN(jsonTag, ",", 2)[0]

		// 如果匹配，返回字段
		if jsonFieldName == jsonName {
			return field, true
		}
	}

	return reflect.StructField{}, false
}
