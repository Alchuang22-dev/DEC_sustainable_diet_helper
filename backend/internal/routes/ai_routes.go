// routes/ai_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"

    "github.com/gin-gonic/gin"
)

func RegisterAIRoutes(router *gin.Engine, db *gorm.DB) {
    aiController := controllers.AIController{
        DB: db,
    }

    aiGroup := router.Group("/ai")
    {
        // 需要认证的路由
        authGroup := aiGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            // 分析图像
            authGroup.POST("/analyze-image", aiController.AnalyzeImage)
        }
    }
}