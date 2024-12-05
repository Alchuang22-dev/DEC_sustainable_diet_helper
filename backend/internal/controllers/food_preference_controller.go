// internal/controllers/food_preference_controller.go
package controllers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
)

type FoodPreferenceController struct {
    DB *gorm.DB
}

// 验证偏好是否存在于配置文件中
func validatePreference(preferenceName string) bool {
    // 读取配置文件
    filePath := filepath.Join("data", "food_preference", "foodPreferences.json")
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return false
    }

    var preferences map[string]interface{}
    if err := json.Unmarshal(data, &preferences); err != nil {
        return false
    }

    _, exists := preferences[preferenceName]
    return exists
}

// AddFoodPreference 添加食物偏好
func (fpc *FoodPreferenceController) AddFoodPreference(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var request struct {
        PreferenceName string `json:"preference_name" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 验证偏好是否存在
    if !validatePreference(request.PreferenceName) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preference name"})
        return
    }

    // 检查是否已存在该偏好
    var existingPreference models.FoodPreference
    result := fpc.DB.Where("user_id = ? AND name = ?", userID, request.PreferenceName).First(&existingPreference)
    if result.Error == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Preference already exists"})
        return
    }

    // 创建新的食物偏好
    preference := models.FoodPreference{
        UserID: userID.(uint),
        Name:   request.PreferenceName,
    }

    if err := fpc.DB.Create(&preference).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add preference"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Food preference added successfully",
        "preference": preference,
    })
}

// DeleteFoodPreference 删除食物偏好
func (fpc *FoodPreferenceController) DeleteFoodPreference(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var request struct {
        PreferenceName string `json:"preference_name" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 验证偏好是否存在
    if !validatePreference(request.PreferenceName) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preference name"})
        return
    }

    result := fpc.DB.Where("user_id = ? AND name = ?", userID, request.PreferenceName).Delete(&models.FoodPreference{})
    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Preference not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Food preference deleted successfully",
    })
}