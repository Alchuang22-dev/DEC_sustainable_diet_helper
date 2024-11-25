package utils

import (
    "errors"
    "strconv"
    "time"
    "fmt"

    "github.com/golang-jwt/jwt/v4"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
)

// 生成 JWT
func GenerateJWT(userID uint) (string, error) {
    claims := &jwt.RegisteredClaims{
        Subject:   strconv.Itoa(int(userID)),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTConfig.Expiration)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(config.JWTSecretKey)
}

// 验证 JWT
func ValidateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
    // 解析并验证 JWT
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 检查签名方法是否为 HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.JWTSecretKey), nil // 确保密钥为字节数组
    })

    if err != nil {
        // 返回更具体的错误原因
        if validationErr, ok := err.(*jwt.ValidationError); ok {
            switch {
            case validationErr.Errors&jwt.ValidationErrorExpired != 0:
                return nil, errors.New("token has expired")
            case validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0:
                return nil, errors.New("invalid token signature")
            default:
                return nil, errors.New("invalid token")
            }
        }
        return nil, err
    }

    // 验证 Claims 并确保 Token 有效
    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token claims")
}