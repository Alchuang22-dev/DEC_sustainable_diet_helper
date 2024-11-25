// internal/controllers/food_controller.go
package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
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
// @Param lang query string true "语言选择 (zh/en)"
// @Success 200 {array} models.FoodNameResponse
// @Router /foods/names [get]
func (fc *FoodController) GetFoodNames(c *gin.Context) {
    language := c.Query("lang")
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