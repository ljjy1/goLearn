package auth

import (
	"context"
	"fmt"
	"homework4/config"
	"homework4/internal/app/redis"
	"homework4/internal/utils/jwt"
	
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"homework4/internal/middleware/response"
)

const (
	AuthUserKey = "auth_user"
)

type AuthUser struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(response.NewUnauthorizedError("未提供认证令牌"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(response.NewUnauthorizedError("认证令牌格式错误"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := jwt.ParseToken(tokenString, config.Cfg.JWT.Secret)
		if err != nil {
			c.Error(response.NewUnauthorizedError("认证令牌无效"))
			c.Abort()
			return
		}

		storedToken, err := RedisGetToken(claims.UserID)
		if err != nil {
			c.Error(response.NewUnauthorizedError("认证令牌过期"))
			c.Abort()
			return
		}

		if storedToken != tokenString {
			c.Error(response.NewUnauthorizedError("认证令牌不匹配"))
			c.Abort()
			return
		}

		c.Set(AuthUserKey, AuthUser{
			UserID:   claims.UserID,
			Username: claims.Username,
			Nickname: claims.Nickname,
		})

		c.Next()
	}
}

/**
 * @description: redis存储用户token到Redis
 * @param {string} token token字符串
 * @param {uint} userID 用户ID
 * @return {error} 错误信息
 */
func RedisStoreToken(token string, userID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("token:%d", userID)
	return redis.RedisClient.Set(ctx, key, token, time.Duration(config.Cfg.JWT.Expires)*time.Hour).Err()
}

/**
 * @description: 获取当前登录用户信息
 */
func GetCurrentAuthUser(c *gin.Context) AuthUser {
	return c.MustGet(AuthUserKey).(AuthUser)
}

/**
 * @description: 从Redis获取用户token
 * @param {uint} userID 用户ID
 * @return {string} token
 */
func RedisGetToken(userID uint) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("token:%d", userID)
	return redis.RedisClient.Get(ctx, key).Result()
}

/**
 * @description: 删除redis用户token
 * @param {uint} userID 用户ID
 * @return {error} 错误信息
 */
func RedisDeleteToken(userID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("token:%d", userID)
	return redis.RedisClient.Del(ctx, key).Err()
}

/**
 * @description: 生成JWT token
 * @param {uint} userID 用户ID
 * @param {string} username 登录账号
 * @return {string} token
 */
func GenerateToken(userID uint, username string, nickname string) string {
	token, _ := jwt.GenerateToken(userID, username, nickname, config.Cfg.JWT.Secret, config.Cfg.JWT.Expires)
	RedisStoreToken(token, userID)
	return token
}
