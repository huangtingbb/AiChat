package models

import (
	"time"

	"gorm.io/gorm"
)

// AIModel AI模型配置
type AIModel struct {
	Id               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"size:50;not null;uniqueIndex"` // 模型名称，如zhipu-glm-4
	DisplayName      string         `json:"display_name" gorm:"size:100;not null"`    // 显示名称，如智谱GLM-4
	Provider         string         `json:"provider" gorm:"size:50;not null"`         // 提供商，如zhipu、openai
	Type             string         `json:"type" gorm:"size:20;not null"`             // 类型，如chat、image
	URL              string         `json:"url" gorm:"size:255"`                      // API请求URL
	MaxTokens        int            `json:"max_tokens" gorm:"default:2048"`           // 最大Token数
	Temperature      float64        `json:"temperature" gorm:"default:0.7"`           // 温度参数
	TopP             float64        `json:"top_p" gorm:"default:0.9"`                 // Top-P参数
	PresencePenalty  float64        `json:"presence_penalty" gorm:"default:0"`        // 重复惩罚
	FrequencyPenalty float64        `json:"frequency_penalty" gorm:"default:0"`       // 频率惩罚
	Enabled          bool           `json:"enabled" gorm:"default:true"`              // 是否启用
	IsDefault        bool           `json:"is_default" gorm:"default:false"`          // 是否为默认模型
	ApiParameters    string         `json:"api_parameters" gorm:"type:json"`          // API参数(JSON格式)
	Description      string         `json:"description" gorm:"type:text"`             // 模型描述
	Class            string         `json:"class" gorm:"size:10"`                     // 大分类 workflow、bot、bigmodal
	ClassId          string         `json:"class_id" gorm:"size:25"`                  // coze大分类的id,对应workflow_id,bot_id等
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-"`
}

// AIModelUsage AI模型使用记录
type AIModelUsage struct {
	Id               uint           `json:"id" gorm:"primaryKey"`
	UserId           uint           `json:"user_id" gorm:"not null;index"`
	ModelId          uint           `json:"model_id" gorm:"not null;index"`
	MessageId        uint           `json:"message_id" gorm:"index"`
	Prompt           string         `json:"prompt" gorm:"type:text"`
	Response         string         `json:"response" gorm:"type:text"`
	PromptTokens     int            `json:"prompt_tokens" gorm:"default:0"`
	CompletionTokens int            `json:"completion_tokens" gorm:"default:0"`
	TotalTokens      int            `json:"total_tokens" gorm:"default:0"`
	Duration         int            `json:"duration" gorm:"default:0"`      // 耗时(毫秒)
	Status           string         `json:"status" gorm:"size:20;not null"` // 状态: success, error
	ErrorMsg         string         `json:"error_msg" gorm:"type:text"`
	Cost             float64        `json:"cost" gorm:"default:0"` // 计费金额
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `json:"-"`
}

// TableName 定义表名
func (AIModel) TableName() string {
	return "ai_model"
}

// TableName 定义表名
func (AIModelUsage) TableName() string {
	return "ai_model_usage"
}
