package database

import (
	"log"

	"chatbot-app/backend/models"
)

// InitTables 初始化数据库表
func InitTables() error {
	log.Println("正在初始化数据库表...")

	// 使用GORM的AutoMigrate自动迁移表结构
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.Message{},
	); err != nil {
		return err
	}

	log.Println("数据库表初始化完成")
	return nil
}
