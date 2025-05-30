package database

import (
	"fmt"
	"log"
	"time"

	"chatbot-app/backend/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取SQL DB实例失败: %w", err)
	}

	// 设置空闲连接池中的最大连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("MySQL连接成功")
	return nil
}
