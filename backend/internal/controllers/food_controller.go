// internal/controllers/food_controller.go
package controllers

// import (
//     "net/http"

//     "github.com/gin-gonic/gin"
//     "gorm.io/gorm"
// 	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
// )

// // FoodController 结构体
// type FoodController struct {
//     DB *gorm.DB
// }

// // NewFoodController 构造函数
// func NewFoodController(db *gorm.DB) *FoodController {
//     return &FoodController{DB: db}
// }

// // CreateFoodRequest 创建食物的请求结构体
// type CreateFoodRequest struct {
//     Name          string           `json:"name" binding:"required"`
//     Weight        float64          `json:"weight" binding:"required"`
//     Price         float64          `json:"price" binding:"required"`
//     TransportMode string           `json:"transport_mode" binding:"required,oneof=陆运 海运 空运"`
//     Location      string           `json:"location" binding:"required,oneof=本地 外地"`
//     Emissions     float64          `json:"emissions" binding:"required"`
//     Nutrition     NutritionRequest `json:"nutrition" binding:"required"`
// }

// // NutritionRequest 营养信息请求结构体
// type NutritionRequest struct {
//     Calories      float64 `json:"calories" binding:"required"`
//     Protein       float64 `json:"protein" binding:"required"`
//     Fat           float64 `json:"fat" binding:"required"`
//     Carbohydrates float64 `json:"carbohydrates" binding:"required"`
// }

// // CreateFood 创建新的食物
// func (fc *FoodController) CreateFood(c *gin.Context) {
//     var req CreateFoodRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
//     // 检查 userID 是否有效
//     var userExists bool
//     if err := nc.DB.Model(&models.User{}).Select("count(*) > 0").Where("id = ?", comment.UserID).Find(&userExists).Error; err != nil || !userExists {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
//         return
//     }

//     // 创建营养记录
//     nutrition := models.Nutrition{
//         Calories:      req.Nutrition.Calories,
//         Protein:       req.Nutrition.Protein,
//         Fat:           req.Nutrition.Fat,
//         Carbohydrates: req.Nutrition.Carbohydrates,
//     }

//     if err := fc.DB.Create(&nutrition).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create nutrition data"})
//         return
//     }

//     // 创建食物记录
//     food := models.Food{
//         Name:          req.Name,
//         Weight:        req.Weight,
//         Price:         req.Price,
//         TransportMode: req.TransportMode,
//         Location:      req.Location,
//         Emissions:     req.Emissions,
//         NutritionID:   nutrition.ID,
//     }

//     if err := fc.DB.Create(&food).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food data"})
//         return
//     }

//     // 预加载营养信息
//     if err := fc.DB.Preload("Nutrition").First(&food, food.ID).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load food data"})
//         return
//     }

//     c.JSON(http.StatusCreated, gin.H{"message": "Food created successfully", "food": food})
// }

// // GetAllFoodNames 获取所有食物名称
// func (fc *FoodController) GetAllFoodNames(c *gin.Context) {
//     var foods []models.Food
//     if err := fc.DB.Select("name").Find(&foods).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch food names"})
//         return
//     }

//     // 提取名称
//     var names []string
//     for _, food := range foods {
//         names = append(names, food.Name)
//     }

//     c.JSON(http.StatusOK, gin.H{"food_names": names})
// }

// // GetFoodInfoByName 获取食物信息
// func (fc *FoodController) GetFoodInfoByName(c *gin.Context) {
//     name := c.Query("name")
//     if name == "" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
//         return
//     }

//     var food models.Food
//     if err := fc.DB.Preload("Nutrition").Where("name = ?", name).First(&food).Error; err != nil {
//         if err == gorm.ErrRecordNotFound {
//             c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
//         } else {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch food data"})
//         }
//         return
//     }

//     // 返回所需的信息
//     response := gin.H{
//         "name":      food.Name,
//         "emissions": food.Emissions,
//         "nutrition": gin.H{
//             "calories":      food.Nutrition.Calories,
//             "protein":       food.Nutrition.Protein,
//             "fat":           food.Nutrition.Fat,
//             "carbohydrates": food.Nutrition.Carbohydrates,
//         },
//     }

//     c.JSON(http.StatusOK, response)
// }
