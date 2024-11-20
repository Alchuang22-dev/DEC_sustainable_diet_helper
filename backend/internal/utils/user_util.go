package utils

import (
    "time"
    "regexp"
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

// TODO 验证手机号
func IsValidPhoneNumber(phone string) bool {
    // 简单正则示例，适配国际和国内手机号格式
    regex := `^(\+?\d{1,4})?\d{7,10}$`
    re := regexp.MustCompile(regex)
    return re.MatchString(phone)
}

// 生成随机用户名
func GenerateRandomNickname() string {
    rand.Seed(time.Now().UnixNano())
    letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    hash := make([]byte, 6)
    for i := range hash {
        hash[i] = letters[rand.Intn(len(letters))]
    }
    return "User" + string(hash)
}