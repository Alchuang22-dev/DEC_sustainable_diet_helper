// routes/routes.go 中添加
package routes

import (
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterIngredientRoutes(router *gin.Engine, db *gorm.DB) {
    controller := &controllers.IngredientController{DB: db}
    
    ingredientGroup := router.Group("/ingredients")
    ingredientGroup.Use(middleware.AuthMiddleware())
    {
        ingredientGroup.POST("/recommend", controller.RecommendIngredients)
    }
}