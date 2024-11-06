// routes/user_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/controllers"

    "github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
    userController := controllers.NewUserController(db)

    userGroup := router.Group("/users")
    {
        userGroup.GET("/", userController.GetAllUsers)
        userGroup.POST("/", userController.CreateUser)
        userGroup.GET("/:id", userController.GetUserByID)
        userGroup.PUT("/:id", userController.UpdateUser)
        userGroup.DELETE("/:id", userController.DeleteUser)
    }
}