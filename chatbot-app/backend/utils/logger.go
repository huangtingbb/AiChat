package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"chatbot-app/backend/config"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger 初始化日志系统
func InitLogger(cfg *config.LogConfig) error {
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.DateTime,
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.DateTime,
		})
	}

	// 创建日志目录
	if cfg.LogFile != "" {
		logDir := filepath.Dir(cfg.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// 创建日志文件
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		// 同时输出到控制台和文件
		Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	return nil
}

// LogInfo 记录信息日志
func LogInfo(msg string, fields ...map[string]interface{}) {
	entry := Logger.WithTime(time.Now())
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Info(msg)
}

// LogError 记录错误日志
func LogError(msg string, err error, fields ...map[string]interface{}) {
	entry := Logger.WithTime(time.Now())
	if err != nil {
		entry = entry.WithError(err)
	}
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Error(msg)
}

// LogWarn 记录警告日志
func LogWarn(msg string, fields ...map[string]interface{}) {
	entry := Logger.WithTime(time.Now())
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Warn(msg)
}

// LogDebug 记录调试日志
func LogDebug(msg string, fields ...map[string]interface{}) {
	entry := Logger.WithTime(time.Now())
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Debug(msg)
}
