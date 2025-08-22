package coze

import (
	"chatbot-app/backend/config"
	"chatbot-app/backend/database"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coze-dev/coze-go"
)

const (
	COZE_TOKEN_KEY       = "coze:access_token"
	TOKEN_EXPIRE_MINUTES = 14
)

type Client struct {
	Config     *config.CozeConfig
	Api        *coze.CozeAPI
	BotID      string // 动态设置的BotID
	WorkflowID string // 动态设置的WorkflowID
}

func New() (*Client, error) {
	cozeConv := &Client{
		Config: config.GetCozeConfig(),
	}
	token, err := GetToken()
	if err != nil {
		return nil, fmt.Errorf("获取Coze Token失败: %v", err)
	}
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	cozeApi := coze.NewCozeAPI(coze.NewTokenAuth(token), coze.WithBaseURL(cozeConv.Config.APIURL), coze.WithHttpClient(httpClient))
	cozeConv.Api = &cozeApi

	// 设置默认的BotID和WorkflowID
	cozeConv.BotID = cozeConv.Config.BotID
	cozeConv.WorkflowID = cozeConv.Config.WorkFlowID

	return cozeConv, nil
}

// NewWithParams 创建带有指定参数的Coze客户端
func NewWithParams(botID, workflowID string) (*Client, error) {
	client, err := New()
	if err != nil {
		return nil, err
	}

	// 覆盖动态参数
	if botID != "" {
		client.BotID = botID
	}
	if workflowID != "" {
		client.WorkflowID = workflowID
	}

	return client, nil
}

func GetToken() (string, error) {
	ctx := context.Background()

	// 首先尝试从Redis获取缓存的token
	if database.RedisClient != nil {
		cachedToken, err := database.RedisClient.Get(ctx, COZE_TOKEN_KEY).Result()
		if err == nil && cachedToken != "" {
			fmt.Println("使用缓存的token")
			return cachedToken, nil
		}
	}

	// Redis中没有token或者Redis不可用，从API获取新token
	fmt.Println("从API获取新token")

	cozeConfig := config.GetCozeConfig()
	fmt.Println(cozeConfig)
	// 优先使用配置文件中的私钥字符串，如果为空则尝试读取文件
	var jwtOauthPrivateKey string
	if cozeConfig.PrivateKey != "" {
		jwtOauthPrivateKey = cozeConfig.PrivateKey
	} else {
		return "", fmt.Errorf("未配置私钥")
	}

	// 验证必需的配置参数
	if cozeConfig.ClientID == "" {
		return "", fmt.Errorf("Coze_Client_ID未配置，请设置COZE_CLIENT_ID环境变量")
	}
	if cozeConfig.PublicKeyID == "" {
		return "", fmt.Errorf("COZE_PUBLIC_KEY_ID环境变量未配置，请设置COZE_PUBLIC_KEY_ID环境变量")
	}

	oauth, err := coze.NewJWTOAuthClient(coze.NewJWTOAuthClientParam{
		PrivateKeyPEM: jwtOauthPrivateKey,
		ClientID:      cozeConfig.ClientID,
		PublicKey:     cozeConfig.PublicKeyID,
	}, coze.WithAuthBaseURL(cozeConfig.APIURL))

	if err != nil {
		// 提供更详细的错误信息
		if strings.Contains(err.Error(), "private key") {
			return "", fmt.Errorf("私钥格式错误: %v\n建议检查:\n1. 私钥是否为PEM格式\n2. 私钥是否完整包含头尾标记\n3. 私钥内容是否正确", err)
		}
		return "", fmt.Errorf("创建JWT OAuth客户端失败: %v", err)
	}

	resp, err := oauth.GetAccessToken(ctx, nil)
	fmt.Println("resp", resp)
	if err != nil {
		return "", fmt.Errorf("获取AccessToken失败: %v", err)
	}

	// 将新token存储到Redis，设置14分钟过期时间
	if database.RedisClient != nil {
		expiration := time.Duration(TOKEN_EXPIRE_MINUTES) * time.Minute
		err = database.RedisClient.Set(ctx, COZE_TOKEN_KEY, resp.AccessToken, expiration).Err()
		if err != nil {
			fmt.Printf("警告: 无法将token存储到Redis: %v\n", err)
		} else {
			fmt.Printf("token已存储到Redis，有效期: %d分钟\n", TOKEN_EXPIRE_MINUTES)
		}
	}

	return resp.AccessToken, nil
}
