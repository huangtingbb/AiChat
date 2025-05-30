package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"size:50;not null;unique"`
	Password  string         `json:"-" gorm:"size:100;not null"`
	Email     string         `json:"email" gorm:"size:100;unique"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Role      string         `json:"role" gorm:"size:20;default:'user'"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 定义表名
func (User) TableName() string {
	return "user"
}
