package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey     string
	ExpireTime    time.Duration
	RefreshTime   time.Duration
	Issuer        string
}

// GetDefaultJWTConfig 获取默认JWT配置
func GetDefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey:   "mes-system-secret-key-2024",
		ExpireTime:  24 * time.Hour,
		RefreshTime: 7 * 24 * time.Hour,
		Issuer:      "mes-system",
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(config *JWTConfig, userID uint, username, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(config.ExpireTime)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// ParseToken 解析JWT令牌
func ParseToken(config *JWTConfig, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新JWT令牌
func RefreshToken(config *JWTConfig, tokenString string) (string, error) {
	claims, err := ParseToken(config, tokenString)
	if err != nil {
		return "", err
	}

	// 检查是否在刷新时间范围内
	if time.Until(claims.ExpiresAt.Time) > config.RefreshTime {
		return "", errors.New("token is not eligible for refresh")
	}

	return GenerateToken(config, claims.UserID, claims.Username, claims.Role)
}