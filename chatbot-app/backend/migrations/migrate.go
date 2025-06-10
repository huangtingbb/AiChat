package main

import (
	"chatbot-app/backend/config"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	// 获取所有SQL文件
	sqlFiles, err := filepath.Glob(filepath.Join(workDir, "*.sql"))
	if err != nil {
		log.Fatalf("查找SQL文件失败: %v", err)
	}

	// 如果在当前目录找不到SQL文件，尝试在上级目录的migrations文件夹中查找
	if len(sqlFiles) == 0 {
		parentDir := filepath.Dir(workDir)
		sqlFiles, err = filepath.Glob(filepath.Join(parentDir, "migrations", "*.sql"))
		if err != nil {
			log.Fatalf("查找SQL文件失败: %v", err)
		}
	}

	if len(sqlFiles) == 0 {
		log.Fatalf("未找到任何SQL文件")
	}

	// 按文件名排序
	sort.Strings(sqlFiles)

	// 读取SQL文件
	// 遍历所有SQL文件并执行
	for _, sqlFile := range sqlFiles {
		log.Printf("正在执行SQL文件: %s", filepath.Base(sqlFile))

		// 读取SQL文件内容
		content, err := ioutil.ReadFile(sqlFile)
		if err != nil {
			log.Printf("警告: 无法读取SQL文件 %s: %v", sqlFile, err)
			continue
		}

		// 分割并执行SQL语句
		statements := strings.Split(string(content), ";")
		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}

			// 跳过创建和使用数据库的语句
			if strings.Contains(strings.ToUpper(statement), "CREATE DATABASE") ||
				strings.Contains(strings.ToUpper(statement), "USE ") {
				log.Printf("跳过语句: %s", statement)
				continue
			}

			_, err = db.Exec(statement)
			if err != nil {
				log.Printf("警告: 执行SQL语句失败: %v\n语句: %s", err, statement)
				// 继续执行其他语句
			} else {
				log.Printf("成功执行SQL语句")
			}
		}
	}

	log.Println("数据库迁移完成")
}
