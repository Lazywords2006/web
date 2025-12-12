package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTSecret = "your-secret-key-change-in-production" // 生产环境必须更换
)

// GenerateLicenseKey 生成许可证密钥
// 格式: XXXX-XXXX-XXXX-XXXX-XXXX
func GenerateLicenseKey() (string, error) {
	// 生成20字节随机数据
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// 转换为16进制字符串
	hexStr := strings.ToUpper(hex.EncodeToString(bytes))

	// 格式化为 XXXX-XXXX-XXXX-XXXX-XXXX
	parts := []string{
		hexStr[0:4],
		hexStr[4:8],
		hexStr[8:12],
		hexStr[12:16],
		hexStr[16:20],
	}

	return strings.Join(parts, "-"), nil
}

// GenerateOrderID 生成订单ID
func GenerateOrderID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ORD-" + strings.ToUpper(hex.EncodeToString(bytes)), nil
}

// GenerateJWT 生成JWT令牌
func GenerateJWT(licenseKey, hwid string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"license_key": licenseKey,
		"hwid":        hwid,
		"exp":         expiresAt.Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT 验证JWT令牌
func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// HashPassword 密码加密（简化版，生产环境应使用bcrypt）
func HashPassword(password string) string {
	// TODO: 生产环境使用 bcrypt.GenerateFromPassword
	return password // 临时实现
}

// CheckPassword 验证密码
func CheckPassword(hashedPassword, password string) bool {
	// TODO: 生产环境使用 bcrypt.CompareHashAndPassword
	return hashedPassword == password // 临时实现
}

// GenerateAPIKey 生成API密钥
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk_" + hex.EncodeToString(bytes), nil
}
