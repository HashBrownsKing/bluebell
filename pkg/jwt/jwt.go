package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JwtSecret  = []byte("一闪一闪亮晶晶") // 密钥
	ExpireTime = 24 * time.Hour    // Toker过期时间
)

var (
	ErrorInvalidSignMethod = errors.New("非法的签名算法")
	ErrorInvalidToken      = errors.New("非法的token")
)

type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成JWT token
func GenerateToken(userID int64, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireTime)),
			Issuer:    "bluebell",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecret)
}

// 根据给定的 JWT 字符串，解析出数据。(解释token)
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(
		tokenString,     // 前端传入的token字符串
		&CustomClaims{}, // 空的自定义的Claims(用来装填解析后的数据)
		func(token *jwt.Token) (interface{}, error) { // 关键回调函数：提供密钥 + 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrorInvalidSignMethod
			}
			return JwtSecret, nil
		},
	)
	if err != nil { // 解析失败
		return nil, err
	}
	// // 验证：解析结果是我们定义的 CustomClaims 类型 + Token 状态合法
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrorInvalidToken
}
