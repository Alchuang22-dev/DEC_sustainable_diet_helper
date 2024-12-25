package utils

import (
    "github.com/golang-jwt/jwt/v4"
)

type UtilsInterface interface {
    GenerateAccessToken(userID uint) (string, error)
    GenerateRefreshToken(userID uint) (string, error)
    CopyFile(src, dst string) error
	ValidateToken(tokenString string) (*jwt.RegisteredClaims, error)
}

type UtilsImpl struct{}

func (u UtilsImpl) GenerateAccessToken(userID uint) (string, error) {
    return GenerateAccessToken(userID)
}

func (u UtilsImpl) GenerateRefreshToken(userID uint) (string, error) {
    return GenerateRefreshToken(userID)
}

func (u UtilsImpl) CopyFile(src, dst string) error {
    return CopyFile(src, dst)
}

func (u UtilsImpl) ValidateToken(src string) (*jwt.RegisteredClaims, error) {
    return ValidateToken(src)
}