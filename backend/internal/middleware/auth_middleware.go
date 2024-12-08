package middleware

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils" // 替换为你的 utils 包路径
)

// AuthMiddleware 用于验证 JWT 的中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            c.Abort()
            return
        }

        // 假设 Bearer <token>
        var tokenString string
        _, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        claims, err := utils.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            c.Abort()
            return
        }

        userID, err := strconv.Atoi(claims.Subject)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
            c.Abort()
            return
        }

        // 设置用户 ID 到上下文
        c.Set("user_id", uint(userID))

        c.Next()
    }
}