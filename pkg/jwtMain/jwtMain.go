package jwtMain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义声明
type MyClaims struct {
	UserPhone int64  `json:"userPhone"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

var secret = []byte("fun-256-bit-secret")

// GenerateToken 签发
func GenerateToken(userPhone int64, username string, expirationTime time.Time) (string, error) {
	claims := MyClaims{
		UserPhone: userPhone,
		Username:  username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "myapp",
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// ParseToken 解析并校验
func ParseToken(tokenStr string) (*MyClaims, error) {
	
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
