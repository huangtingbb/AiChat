package database

import (
	"fmt"
	"log"

	"chatbot-app/backend/config"

	"github.com/streadway/amqp"
)

var RabbitMQConn *amqp.Connection
var RabbitMQChannel *amqp.Channel

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ(cfg *config.RabbitMQConfig) error {
	// 创建连接
	var err error
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost)

	RabbitMQConn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("连接RabbitMQ失败: %w", err)
	}

	// 创建通道
	RabbitMQChannel, err = RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("创建RabbitMQ通道失败: %w", err)
	}

	log.Println("RabbitMQ连接成功")
	return nil
}

// CloseRabbitMQ 关闭RabbitMQ连接
func CloseRabbitMQ() {
	if RabbitMQChannel != nil {
		RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}
