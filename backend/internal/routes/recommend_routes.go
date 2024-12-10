// routes/routes.go 中添加
package routes

import (
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRecommendRoutes(router *gin.Engine, db *gorm.DB) {
    controller := &controllers.RecommendController{DB: db}
    
    ingredientGroup := router.Group("/ingredients")
    ingredientGroup.Use(middleware.AuthMiddleware())
    {
        ingredientGroup.POST("/recommend", controller.RecommendIngredients)
        ingredientGroup.POST("/set", controller.SetUserSelectedFoods)
    }

    recipeGroup := router.Group("/recipes")
    recipeGroup.Use(middleware.AuthMiddleware())
    {
        recipeGroup.POST("/recommend", controller.RecommendRecipes)
    }
}
