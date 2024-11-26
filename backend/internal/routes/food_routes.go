// internal/routes/food_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    // "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
)

func RegisterFoodRoutes(router *gin.Engine, db *gorm.DB) {
    foodController := controllers.NewFoodController(db)
    foodGroup := router.Group("/foods")
    {
        // // 需要认证的路由
        // authGroup := foodGroup.Group("")
        // authGroup.Use(middleware.AuthMiddleware())
        // {
        //     // 创建食物（管理员操作）
        //     authGroup.POST("/create", foodController.CreateFood)
        //     // 更新食物信息（管理员操作）
        //     authGroup.PUT("/:id", foodController.UpdateFood)
        //     // 删除食物（管理员操作）
        //     authGroup.DELETE("/:id", foodController.DeleteFood)
        //     // 获取用户相关的食物分析
        //     authGroup.POST("/analyze", foodController.AnalyzeFood)
        //     // 用户的食物收藏
        //     authGroup.POST("/:id/favourite", foodController.FavoriteFood)
        //     // 取消收藏
        //     authGroup.POST("/:id/cancel_favourite", foodController.CancelFavoriteFood)
        // }

        // 不需要认证的路由
        // 获取食物名称列表
        foodGroup.GET("/names", foodController.GetFoodNames)
		// 计算食物的营养成分和碳排放
		foodGroup.POST("/calculate", foodController.CalculateNutritionAndEmission)
        // // 获取单个食物详情
        // foodGroup.GET("/:id", foodController.GetFoodDetail)
        // // 获取所有食物列表
        // foodGroup.GET("", foodController.GetAllFoods)
        // // 搜索食物
        // foodGroup.GET("/search", foodController.SearchFoods)
    }
}