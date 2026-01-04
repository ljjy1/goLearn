package service

import (
	"errors"
	"homework4/internal/app/mysql"
	"homework4/internal/middleware/auth"
	"homework4/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{db: mysql.DB}
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" example:"testuser"` // 用户名 必传，3-20个字符
	Nickname string `json:"nickname" binding:"required,min=1,max=64" example:"测试用户"`     // 昵称 必传，1-64个字符
	Password string `json:"password" binding:"required,min=6" example:"123456"`          // 密码 必传，至少6个字符
	Email    string `json:"email" binding:"omitempty,email" example:"test@example.com"`  // 邮箱，可选，格式校验
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"` // 用户名 必传
	Password string `json:"password" binding:"required" example:"123456"`   // 密码 必传
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT 令牌
	UserID   uint   `json:"user_id" example:"1"`                                     // 用户ID
	Username string `json:"username" example:"testuser"`                             // 用户名
	Nickname string `json:"nickname" example:"测试用户"`                                 // 昵称
}

/**
 * @Description: 注册用户
 * @param req
 * @return (*models.User, error)
 */
func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	var existingUser models.User
	//判断用户名是否已经存在
	result := s.db.Where(&models.User{Username: req.Username}).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("用户名已存在")
	}
	//加密明文密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	//创建用户
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

/**
 * @Description: 用户登录
 * @param req
 * @return (*LoginResponse, error)
 */
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user models.User
	//根据用户名查询用户
	result := s.db.Where(&models.User{Username: req.Username}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, result.Error
	}

	//密码对比
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	//生成token
	token := auth.GenerateToken(user.ID, user.Username, user.Nickname)

	//返回token信息和用户信息
	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}
