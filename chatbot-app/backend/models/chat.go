package models

import (
	"time"

	"gorm.io/gorm"
)

// Chat 聊天会话模型
type Chat struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Title     string         `json:"title" gorm:"size:100"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Messages  []Message      `json:"messages"`
}

// Message 聊天消息模型
type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ChatID    uint           `json:"chat_id" gorm:"not null;index"`
	Role      string         `json:"role" gorm:"size:10;not null"` // user 或 assistant
	Content   string         `json:"content" gorm:"type:text;not null"`
	Tokens    int            `json:"tokens" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// TableName 定义表名
func (Chat) TableName() string {
	return "chat"
}

// TableName 定义表名
func (Message) TableName() string {
	return "message"
}
