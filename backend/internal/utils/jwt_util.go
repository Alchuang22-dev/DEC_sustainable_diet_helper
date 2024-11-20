package utils

import (
    "errors"
    "strconv"
    "time"

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
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return config.JWTSecretKey, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}