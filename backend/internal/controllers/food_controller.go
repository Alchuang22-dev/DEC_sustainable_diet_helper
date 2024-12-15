// internal/controllers/food_controller.go
package controllers

import (
    "net/http"
    "log"
    "strconv"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"fmt"
)

type FoodController struct {
    DB *gorm.DB
}

func NewFoodController(db *gorm.DB) *FoodController {
    return &FoodController{DB: db}
}


// 不需要认证的路由处理方法

// GetFoodNames godoc
// @Summary 获取食物名称列表
// @Description 获取所有食物的名称列表，支持中文和英文
// @Tags foods
// @Accept json
// @Produce json
// @Param none
// @Success 200 {array} models.FoodInfoResponse
// @Router /foods/names [get]
func (fc *FoodController) GetFoodNames(c *gin.Context) {
    language := c.Query("lang")
    if language == "" {
        language = "zh"
    }
    
    if language != "zh" && language != "en" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid language parameter. Use 'zh' or 'en'",
        })
        return
    }
    
    names, err := models.GetAllFoodNames(fc.DB, language)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, names)
}

// CalculateNutritionAndEmission godoc
// @Summary 计算食物的营养成分和碳排放
// @Description 根据食物ID、价格和重量计算营养成分和碳排放
// @Tags foods
// @Accept json
// @Produce json
// @Param items body []models.FoodCalculateItem true "食物计算请求"
// @Success 200 {array} models.FoodCalculateResult
// @Router /foods/calculate [post]
func (fc *FoodController) CalculateNutritionAndEmission(c *gin.Context) {
    log.Printf("开始处理 CalculateNutritionAndEmission")
    var items []models.FoodCalculateItem

    // 绑定请求数据
    if err := c.ShouldBindJSON(&items); err != nil {
        log.Printf("Invalid request format: " + err.Error())
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format: " + err.Error(),
        })
        return
    }
    log.Printf("绑定请求数据成功")

    // 验证输入数据
    for _, item := range items {
        if item.Weight <= 0 {
            log.Printf("Invalid weight for food ID " + strconv.Itoa(int(item.ID)) + ": weight must be positive")
            c.JSON(http.StatusBadRequest, gin.H{
                "error": fmt.Sprintf("Invalid weight for food ID %d: weight must be positive", item.ID),
            })
            return
        }
        if item.Price <= 0 {
            log.Printf("Invalid price for food ID " + strconv.Itoa(int(item.ID)) + ": price must be positive")
            c.JSON(http.StatusBadRequest, gin.H{
                "error": fmt.Sprintf("Invalid price for food ID %d: price must be positive", item.ID),
            })
            return
        }
    }
    log.Printf("验证输入数据成功")

    // 计算结果
    results, err := models.CalculateFoodNutritionAndEmission(fc.DB, items)
    if err != nil {
        log.Printf("Error calculating nutrition and emission: " + err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    log.Printf("计算结果成功")
    log.Printf("results: %v", results)
    c.JSON(http.StatusOK, results)
}