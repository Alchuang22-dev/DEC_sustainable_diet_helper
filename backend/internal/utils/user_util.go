package utils

import (
    "math/rand"

    "golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedPassword), err
}

// CheckPassword 验证密码
func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// 生成随机用户名
func GenerateRandomNickname() string {
    letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    hash := make([]byte, 6)
    for i := range hash {
        hash[i] = letters[rand.Intn(len(letters))]
    }
    return "User" + string(hash)
}