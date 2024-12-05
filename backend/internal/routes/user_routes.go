// routes/user_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"

    "github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
    userController := controllers.NewUserController(db)

    userGroup := router.Group("/users")
    {
        // 公共路由
        userGroup.POST("/auth", userController.WeChatAuth) // 注册

        // 需要认证的路由
        authGroup := userGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            authGroup.PUT("/:id/set_nickname", userController.SetNickname) // 更新用户名
            authGroup.PUT("/:id/set_avatar", userController.SetAvatar) // 更新头像
        }

        // TODO
        // 查询用户点赞的新闻
        // 查询用户收藏的新闻
        // 查询用户最近浏览的新闻
    }
}