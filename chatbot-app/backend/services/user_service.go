package services

import (
	"errors"
	"time"

	"chatbot-app/backend/database"
	"chatbot-app/backend/models"
	"chatbot-app/backend/utils"
)

// UserService 用户服务
type UserService struct{}

// Register 用户注册
func (s *UserService) Register(username, password, email string) (*models.User, error) {
	// 检查用户是否已存在
	var existingUser models.User
	result := database.DB.Where("username = ? OR email = ?", username, email).First(&existingUser)
	if result.RowsAffected > 0 {
		return nil, errors.New("用户名或邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username:  username,
		Password:  hashedPassword,
		Email:     email,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

// GetUserById 根据Id获取用户
func (s *UserService) GetUserById(id uint) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
