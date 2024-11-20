package utils

import (
    "os"
    "fmt"
    "net/smtp"
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

/*
TODO
1.	验证码限流：
•	限制每个邮箱在一定时间内只能请求一次验证码（如每分钟）。
2.	清理过期验证码：
•	定期清理数据库中过期的验证码。
3.	邮箱配置：
•	使用可靠的邮件服务，如 AWS SES、SendGrid 或 Gmail API。
4.	日志记录：
•	记录验证码发送和验证的操作日志，以便审计。
*/

// GenerateVerificationCode 生成6位随机验证码
func GenerateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// SendEmail 发送邮件验证码
func SendEmail(toEmail string, code string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	authEmail := os.Getenv("SMTP_EMAIL")
	authPassword := os.Getenv("SMTP_PASSWORD")

	// 设置邮件内容
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your verification code is: %s\nThis code will expire in 5 minutes.", code)
	msg := []byte("Subject: " + subject + "\n\n" + body)

	// 发送邮件
	auth := smtp.PlainAuth("", authEmail, authPassword, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, authEmail, []string{toEmail}, msg)
}