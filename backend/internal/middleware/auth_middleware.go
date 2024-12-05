package middleware

import (
    "strings"
    "net/http"
    "strconv"

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

        // 从 Token 中解析 user_id
        userID, err := strconv.ParseUint(claims.Subject, 10, 64) // 假设 Subject 存储的是 user_id
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: invalid user_id"})
            c.Abort()
            return
        }

        // 将用户 ID 存入上下文，供后续使用
        c.Set("user_id", uint(userID))
        c.Next()
    }
}