package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"sync"
	"time"

	"opengeo/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// secretKey JWT密钥（从环境变量读取）
	secretKey []byte
	once      sync.Once
)

// getSecretKey 获取JWT密钥（懒加载，只初始化一次）
func getSecretKey() []byte {
	once.Do(func() {
		key := os.Getenv("JWT_SECRET_KEY")
		if key == "" {
			// 生产环境必须设置环境变量
			// 开发环境自动生成随机密钥
			if os.Getenv("GO_ENV") == "production" {
				panic("JWT_SECRET_KEY environment variable is required in production")
			}
			key = generateRandomKey()
		}
		secretKey = []byte(key)
	})
	return secretKey
}

// generateRandomKey 生成随机密钥
func generateRandomKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic("failed to generate random key: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}

// Claims JWT Claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	TenantID int64  `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID int64, username, email, role string, tenantID int64) (string, error) {
	cfg := config.GetConfig()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.TokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWT.Issuer,
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

// GenerateRefreshToken 生成Refresh Token
func GenerateRefreshToken(userID int64, username string) (string, error) {
	cfg := config.GetConfig()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWT.Issuer,
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

// ValidateToken 验证JWT Token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getSecretKey(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// RefreshToken 刷新Token
func RefreshToken(refreshTokenString string) (string, string, error) {
	claims, err := ValidateToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	newToken, err := GenerateToken(claims.UserID, claims.Username, claims.Email, claims.Role, claims.TenantID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(claims.UserID, claims.Username)
	if err != nil {
		return "", "", err
	}

	return newToken, newRefreshToken, nil
}

// ExtractTokenFromHeader 从Authorization头提取Token
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return authHeader, nil
}
