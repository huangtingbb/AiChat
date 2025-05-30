package database

import (
	"fmt"
	"log"

	"chatbot-app/backend/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

// InitElasticsearch 初始化Elasticsearch连接
func InitElasticsearch(cfg *config.ElasticsearchConfig) error {
	esCfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port),
		},
		Username: cfg.User,
		Password: cfg.Password,
	}

	var err error
	ESClient, err = elasticsearch.NewClient(esCfg)
	if err != nil {
		return fmt.Errorf("创建Elasticsearch客户端失败: %w", err)
	}

	// 测试连接
	res, err := ESClient.Info()
	if err != nil {
		return fmt.Errorf("连接Elasticsearch失败: %w", err)
	}
	defer res.Body.Close()

	log.Println("Elasticsearch连接成功")
	return nil
}
