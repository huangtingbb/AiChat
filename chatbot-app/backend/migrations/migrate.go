package main

import (
	"chatbot-app/backend/config"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("开始数据库迁移...")

	// 从配置文件读取数据库配置
	cfg := config.GetConfig()
	dbUser := cfg.Database.User
	dbPass := cfg.Database.Password
	dbHost := cfg.Database.Host
	dbPort := cfg.Database.Port
	dbName := cfg.Database.DBName

	log.Printf("使用数据库配置: 用户=%s, 主机=%s, 端口=%s, 数据库=%s", dbUser, dbHost, dbPort, dbName)

	// 连接MySQL（不指定数据库名，因为我们要先创建数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPass, dbHost, dbPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接到MySQL: %v", err)
	}
	defer db.Close()

	// 检查连接
	if err := db.Ping(); err != nil {
		log.Fatalf("无法ping MySQL: %v", err)
	}

	log.Println("成功连接到MySQL")

	// 创建数据库（如果不存在）
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}
	log.Printf("确保数据库 %s 存在", dbName)

	// 切换到指定数据库
	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		log.Fatalf("切换到数据库 %s 失败: %v", dbName, err)
	}
	log.Printf("已切换到数据库 %s", dbName)

	// 获取当前工作目录，而不是可执行文件的目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作目录失败: %v", err)
	}
	log.Printf("当前工作目录: %s", workDir)

	// SQL文件路径
	sqlFilePath := filepath.Join(workDir, "001_init_schema.sql")

	// 检查文件是否存在
	if _, err := os.Stat(sqlFilePath); os.IsNotExist(err) {
		// 如果在当前目录找不到，尝试在上级目录的migrations文件夹中查找
		parentDir := filepath.Dir(workDir)
		alternatePath := filepath.Join(parentDir, "migrations", "001_init_schema.sql")
		log.Printf("在当前目录找不到SQL文件，尝试查找: %s", alternatePath)

		if _, err := os.Stat(alternatePath); os.IsNotExist(err) {
			log.Fatalf("SQL文件不存在: %s 或 %s", sqlFilePath, alternatePath)
		}
		sqlFilePath = alternatePath
	}

	log.Printf("找到SQL文件: %s", sqlFilePath)

	// 读取SQL文件
	content, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("无法读取SQL文件: %v", err)
	}

	// 执行SQL语句
	// 由于GO的数据库驱动一次只能执行一条SQL语句，我们需要分割多条语句
	statements := strings.Split(string(content), ";")
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		// 跳过创建和使用数据库的语句，因为我们已经处理过了
		if strings.Contains(strings.ToUpper(statement), "CREATE DATABASE") ||
			strings.Contains(strings.ToUpper(statement), "USE ") {
			log.Printf("跳过语句: %s", statement)
			continue
		}

		_, err = db.Exec(statement)
		if err != nil {
			log.Printf("警告: 执行SQL语句失败: %v\n语句: %s", err, statement)
			// 不要因为一条语句失败就退出，继续执行其他语句
		} else {
			log.Printf("成功执行SQL语句")
		}
	}

	log.Println("数据库迁移完成")
}
