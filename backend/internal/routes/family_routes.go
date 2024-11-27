// routes/user_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    // "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"

    "github.com/gin-gonic/gin"
)

func RegisterFamilyRoutes(router *gin.Engine, db *gorm.DB) {
    familyController := controllers.NewUserController(db)

    fanilyGroup := router.Group("/families")
    {
        // fanilyGroup.POST("/auth", familyController.)
    }
}