// internal/routes/food_preference_routes.go
package routes


import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
)

func RegisterIngredientControllerRoutes(router *gin.Engine, db *gorm.DB) {
    ic := &controllers.IngredientController{DB: db}
	// 添加认证中间件
    authorized := router.Group("/")
    authorized.Use(middleware.AuthMiddleware())
    {
        authorized.POST("/ingredients/recommend", ic.RecommendIngredients)
    }
}