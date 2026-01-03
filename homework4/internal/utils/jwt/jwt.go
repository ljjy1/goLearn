package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	//用户ID
	UserID uint `json:"user_id"`
	//登录账号
	Username string `json:"username"`
	//昵称
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

/**
 * @description: 生成JWT token
 * @param {uint} userID 用户ID
 * @param {string} username 登录账号
 * @param {string} nickname 昵称
 * @param {string} secret 密钥
 * @param {uint8} expires 过期时间(小时)
 * @return {string} token
 */
func GenerateToken(userID uint, username string, nickname string, secret string, expires uint8) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expires) * time.Hour)
	jwtSecret := []byte(secret)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "homework4",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

/**
 * @description: 解析JWT token
 * @param {string} tokenString token字符串
 * @param {string} secret 密钥
 * @return {*Claims} 解析后的claims
 */
func ParseToken(tokenString string, secret string) (*Claims, error) {
	jwtSecret := []byte(secret)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
