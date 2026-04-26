package utils

import (
	"fmt"
	"time"

	"flec_blog/config"
	"flec_blog/internal/model"

	"github.com/golang-jwt/jwt/v4"
)

// TokenType Token类型
type TokenType string

const (
	AccessToken  TokenType = "access"  // 访问令牌
	RefreshToken TokenType = "refresh" // 刷新令牌
)

// Token过期时间常量（小时）
const (
	AccessTokenExpireHours  = 24 * 7  // Access Token 过期时间: 7天
	RefreshTokenExpireHours = 24 * 30 // Refresh Token 过期时间: 30天
)

// Claims JWT 声明结构
type Claims struct {
	UserID       uint           `json:"user_id"`
	Role         model.UserRole `json:"role"`
	TokenType    TokenType      `json:"token_type"`    // token类型
	TokenVersion uint           `json:"token_version"` // token版本号
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(userID uint, role model.UserRole, tokenVersion uint, cfg *config.JWTConfig) (string, error) {
	return generateToken(userID, role, tokenVersion, AccessToken, AccessTokenExpireHours, cfg)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, role model.UserRole, tokenVersion uint, cfg *config.JWTConfig) (string, error) {
	return generateToken(userID, role, tokenVersion, RefreshToken, RefreshTokenExpireHours, cfg)
}

// generateToken 生成指定类型的token
func generateToken(userID uint, role model.UserRole, tokenVersion uint, tokenType TokenType, expireHours int, cfg *config.JWTConfig) (string, error) {
	claims := &Claims{
		UserID:       userID,
		Role:         role,
		TokenType:    tokenType,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用指定的签名方法创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析并验证 Access Token
func ParseToken(tokenString string, cfg *config.JWTConfig) (*Claims, error) {
	claims, err := parseTokenInternal(tokenString, cfg)
	if err != nil {
		return nil, err
	}

	// 验证必须是Access Token
	if claims.TokenType != AccessToken {
		return nil, fmt.Errorf("invalid token type, expected access token")
	}

	return claims, nil
}

// ParseRefreshToken 解析并验证 Refresh Token
func ParseRefreshToken(tokenString string, cfg *config.JWTConfig) (*Claims, error) {
	claims, err := parseTokenInternal(tokenString, cfg)
	if err != nil {
		return nil, err
	}

	// 验证必须是Refresh Token
	if claims.TokenType != RefreshToken {
		return nil, fmt.Errorf("invalid token type, expected refresh token")
	}

	return claims, nil
}

// parseTokenInternal 内部通用token解析函数（不验证类型）
func parseTokenInternal(tokenString string, cfg *config.JWTConfig) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token并转换为Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
