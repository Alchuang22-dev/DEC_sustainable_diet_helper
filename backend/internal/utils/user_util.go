package utils

import (
    "math/rand"
    "fmt"
    "os"
    "io"

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

// 复制文件
func CopyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("无法打开源文件: %v", err)
    }
    defer sourceFile.Close()

    destinationFile, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("无法创建目标文件: %v", err)
    }
    defer destinationFile.Close()

    if _, err := io.Copy(destinationFile, sourceFile); err != nil {
        return fmt.Errorf("复制文件失败: %v", err)
    }
    return nil
}