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
        userGroup.POST("/send_code", userController.SendVerificationCode) // 发送验证码
        userGroup.POST("/register", userController.Register) // 注册
        userGroup.POST("/login", userController.Login) // 登录

        // 需要认证的路由
        authGroup := userGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            authGroup.PUT("/:id/nickname", userController.SetNickname) // 更新用户名
            authGroup.PUT("/:id/phone_number", userController.SetEmail) // 更新手机号
            authGroup.PUT("/:id/password", userController.SetPassword) // 更新密码
        }
    }
}