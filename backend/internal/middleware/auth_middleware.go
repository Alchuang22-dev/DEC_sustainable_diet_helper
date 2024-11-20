package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
)

// AuthMiddleware 用于验证 JWT 的中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
            c.Abort()
            return
        }

        // 验证格式是否为 Bearer <token>
        if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
            c.Abort()
            return
        }

        tokenString := authHeader[7:] // 去掉 "Bearer " 前缀
        claims, err := utils.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 将用户 ID 存入上下文，供后续处理使用
        c.Set("user_id", claims.Subject)
        c.Next()
    }
}