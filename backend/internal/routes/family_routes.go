// routes/user_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"

    "github.com/gin-gonic/gin"
)

func RegisterFamilyRoutes(router *gin.Engine, db *gorm.DB) {
    familyController := controllers.NewFamilyController(db)

    familyGroup := router.Group("/families")
    {
        // 需要认证的路由
        authGroup := familyGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            // 创建家庭
            authGroup.POST("/create", familyController.CreateFamily)
            // 查看自己的家庭的信息
            authGroup.POST("/details", familyController.FamilyDetails)
            // 查看搜索家庭结果
            authGroup.POST("/search", familyController.SearchFamily)
            // 发送加入家庭请求
            authGroup.POST("/join", familyController.JoinFamily)
            // 批准加入家庭
            authGroup.POST("/admit", familyController.AdmitJoinFamily)
        }
    }
}