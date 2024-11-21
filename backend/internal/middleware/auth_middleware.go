package middleware

import (
    "strings"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils" // 替换为你的 utils 包路径
)

// AuthMiddleware 用于验证 JWT 的中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取 Authorization Header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
            c.Abort()
            return
        }

        // 验证 Bearer 格式
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
            c.Abort()
            return
        }

        // 提取 Token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // 验证 Token
        claims, err := utils.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 将用户 ID 存入上下文，供后续使用
        userID := claims.Subject // 确认你的 ValidateJWT 是否将 user_id 存入 Subject
        if userID == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing user_id"})
            c.Abort()
            return
        }

        c.Set("user_id", userID) // 存储为字符串类型
        c.Next()
    }
}