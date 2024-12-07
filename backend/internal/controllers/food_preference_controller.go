// internal/controllers/food_preference_controller.go
package controllers

import (
    "encoding/json"
    "net/http"
    "path/filepath"
    "os"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"

    "log"
)

type FoodPreferenceController struct {
    DB *gorm.DB
}

// 验证偏好是否存在于配置文件中
func validatePreference(preferenceName string) bool {
    // 读取配置文件
    filePath := filepath.Join("data", "food_preference", "foodPreferences.json")
    data, err := os.ReadFile(filePath)
    if err != nil {
        log.Printf("读取配置文件失败: %v\n", err)
        return false
    }

    var preferences map[string]interface{}
    if err := json.Unmarshal(data, &preferences); err != nil {
        log.Printf("解析配置文件失败: %v\n", err)
        return false
    }

    _, exists := preferences[preferenceName]
    log.Printf("偏好 %s 存在: %v\n", preferenceName, exists)
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
        // log详细的错误信息
        log.Printf("错误信息: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 验证偏好是否存在 
    if !validatePreference(request.PreferenceName) {
        log.Printf("偏好不存在: %s\n", request.PreferenceName)
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
        "preference": preference.Name,
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
		"preference": request.PreferenceName,
    })
}
// 添加新的方法到 FoodPreferenceController
func (fpc *FoodPreferenceController) GetUserPreferences(c *gin.Context) {
    // 从上下文获取用户ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var preferences []models.FoodPreference
    if err := fpc.DB.Where("user_id = ?", userID).Find(&preferences).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user preferences"})
        return
    }

	response := make([]gin.H, len(preferences))
	for i, pref := range preferences {
		response[i] = gin.H{
			"id":   pref.ID,
			"name": pref.Name,
		}
	}

    c.JSON(http.StatusOK, response)
}