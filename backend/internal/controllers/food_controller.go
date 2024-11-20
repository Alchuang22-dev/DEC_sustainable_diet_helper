// // controllers/food_controller.go
package controllers

// import (
//     // "net/http"
//     // "strconv"

//     // "github.com/gin-gonic/gin"
//     // "gorm.io/gorm"
//     // "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
// )

// type FoodController struct {
// 	DB *gorm.DB
// }
// // 创建新的食物控制器
// func NewFoodController(db *gorm.DB) *FoodController {
// 	return &FoodController{DB: db}
// }
// // 创建新的食物

// // 获取食物的所有名称，用于前端搜索框

// // 根据食物名称，返回食物的信息（名称，碳排放，营养成分）

// // 根据食物集合，返回推荐的菜谱

// // 对于没有定义食物集合的部分，则按照某种算法进行推荐（暂时不考虑实现，等确定了算法再说）