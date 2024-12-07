// utils/jwt_util.go
package utils

import (
    "errors"
    "strconv"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
)

// GenerateAccessToken 生成 Access Token
func GenerateAccessToken(userID uint) (string, error) {
    claims := &jwt.RegisteredClaims{
        Subject:   strconv.Itoa(int(userID)),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTConfig.AccessTokenExpiration)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(config.JWTSecretKey)
}

// GenerateRefreshToken 生成 Refresh Token
func GenerateRefreshToken(userID uint) (string, error) {
    claims := &jwt.RegisteredClaims{
        Subject:   strconv.Itoa(int(userID)),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTConfig.RefreshTokenExpiration)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        // 可选：添加一个唯一标识符（JTI）用于撤销
        // ID: uuid.New().String(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(config.JWTSecretKey)
}

// ValidateToken 验证任意 Token（Access 或 Refresh）
func ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 检查签名方法是否为 HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return config.JWTSecretKey, nil
    })

    if err != nil {
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

    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token claims")
}