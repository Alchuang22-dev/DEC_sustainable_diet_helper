package routes

import (
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
)

func RegisterNutritionCarbonRoutes(router *gin.Engine, db *gorm.DB) {
    nutritionCarbonController := &controllers.NutritionCarbonController{DB: db}
    
    // 创建营养和碳排放路由组
    ncGroup := router.Group("/nutrition-carbon")
    {
        // 需要认证的路由
        authGroup := ncGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            // 营养目标相关路由
            authGroup.POST("/nutrition/goals", nutritionCarbonController.SetNutritionGoals)
            authGroup.GET("/nutrition/goals", nutritionCarbonController.GetNutritionGoals)
            authGroup.GET("/nutrition/intakes", nutritionCarbonController.GetActualNutrition)

            // 碳排放目标相关路由
            authGroup.POST("/carbon/goals", nutritionCarbonController.SetCarbonGoals)
            authGroup.GET("/carbon/goals", nutritionCarbonController.GetCarbonGoals)
            authGroup.GET("/carbon/intakes", nutritionCarbonController.GetCarbonIntakes)

            // 共享营养碳排放相关路由
            authGroup.POST("/shared/nutrition-carbon", nutritionCarbonController.SetSharedNutritionCarbonIntake)
        }
    }
} 