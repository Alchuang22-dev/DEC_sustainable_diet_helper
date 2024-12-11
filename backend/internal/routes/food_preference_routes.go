// internal/routes/food_preference_routes.go
package routes


import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
)

func RegisterFoodPreferenceRoutes(router *gin.Engine, db *gorm.DB) {
    fpc := &controllers.FoodPreferenceController{DB: db}
    
    // 添加认证中间件
    authorized := router.Group("/")
    authorized.Use(middleware.AuthMiddleware())
    {
        authorized.POST("/preferences", fpc.AddFoodPreference)
        authorized.DELETE("/preferences", fpc.DeleteFoodPreference)
        authorized.GET("/preferences", fpc.GetUserPreferences)

        // 新增不喜欢的食材偏好路由
        authorized.POST("/disliked_preferences", fpc.AddDislikedFoodPreference)
        authorized.DELETE("/disliked_preferences", fpc.DeleteDislikedFoodPreference)
        authorized.GET("/disliked_preferences", fpc.GetUserDislikedPreferences)
    }
}