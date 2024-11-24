// internal/routes/food_routes.go
package routes

// import (
//     "gorm.io/gorm"
//     "github.com/gin-gonic/gin"
//     "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
// )

// // RegisterFoodRoutes 注册食物相关的路由
// func RegisterFoodRoutes(router *gin.Engine, db *gorm.DB) {
//     foodController := controllers.NewFoodController(db)
//     foodGroup := router.Group("/food")
//     {
//         // 创建新的食物
//         foodGroup.POST("/", foodController.CreateFood)
//         // 获取所有食物名称
//         foodGroup.GET("/names", foodController.GetAllFoodNames)
//         // 根据食物名称获取食物信息
//         foodGroup.GET("/info", foodController.GetFoodInfoByName)
//     }
// }
