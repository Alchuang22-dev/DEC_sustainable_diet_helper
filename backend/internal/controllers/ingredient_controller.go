// internal/controllers/ingredient_controller.go
package controllers

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
    "encoding/json"
    "os"
    "slices"
    "math/rand"
)

type IngredientController struct {
    DB *gorm.DB
}

// 请求结构体
type IngredientRecommendRequest struct {
    UseLastIngredients bool     `json:"use_last_ingredients"`
    LikedIngredients   []uint   `json:"liked_ingredients"`
    DislikedIngredients []uint  `json:"disliked_ingredients"`
}

// 响应结构体
type IngredientRecommendResponse struct {
    RecommendedIngredients []struct {
        ID          uint    `json:"id"`
        Name        string  `json:"name"`
    } `json:"recommended_ingredients"`
}

const (
    BASE_SCORE       = 0.5   // 基础分数
    WEIGHT_HISTORY   = 0.2   // 历史选择权重
    WEIGHT_PREFERENCE_LIKE = 0.3  // 用户此次传输的偏好权重(like)
    WEIGHT_PREFERENCE_DISLIKE = -0.2  // 用户此次传输的偏好权重(dislike)
    WEIGHT_FOOD_PREF  = 0.3   // 设置界面的食物偏好类型权重(like)

)

// internal/controllers/ingredient_controller.go

func (ic *IngredientController) loadFoodPreferences(preferences []models.FoodPreference) (Pos_id []uint, Neg_id []uint, err error) {
    data, err := os.ReadFile("data/food_preference/foodPreferences.json")
    if err != nil {
        return nil, nil, err
    }

    var preferencesMap map[string]struct {
        FoodPos []string `json:"food_pos"`
        FoodNeg []string `json:"food_neg"`
    }
    if err := json.Unmarshal(data, &preferencesMap); err != nil {
        return nil, nil, err
    }

    // 用于存储所有偏好的食物列表
    var allPosFood []string
    var allNegFood []string
    isFirst := true

    // 处理每个用户选择的偏好
    for _, pref := range preferences {
        if prefData, exists := preferencesMap[pref.Name]; exists {
            if isFirst {
                // 第一个偏好直接作为基准
                allPosFood = prefData.FoodPos
                allNegFood = prefData.FoodNeg
                isFirst = false
            } else {
                // 对positive food取交集
                allPosFood = intersection(allPosFood, prefData.FoodPos)
                // 对negative food取并集
                allNegFood = union(allNegFood, prefData.FoodNeg)
            }
        }
    }

    // 将食物名称转换为ID
    foodPos_id := make([]uint, 0)
    foodNeg_id := make([]uint, 0)

    // 处理positive foods
    for _, food := range allPosFood {
        var ingredient models.Ingredient
        if err := ic.DB.Where("name = ?", food).First(&ingredient).Error; err != nil {
            continue // 跳过找不到的食材
        }
        foodPos_id = append(foodPos_id, ingredient.ID)
    }

    // 处理negative foods
    for _, food := range allNegFood {
        var ingredient models.Ingredient
        if err := ic.DB.Where("name = ?", food).First(&ingredient).Error; err != nil {
            continue // 跳过找不到的食材
        }
        foodNeg_id = append(foodNeg_id, ingredient.ID)
    }

    return foodPos_id, foodNeg_id, nil
}

// 辅助函数：计算两个字符串切片的交集
func intersection(a, b []string) []string {
    set := make(map[string]bool)
    var result []string

    for _, item := range a {
        set[item] = true
    }

    for _, item := range b {
        if set[item] {
            result = append(result, item)
        }
    }

    return result
}

// 辅助函数：计算两个字符串切片的并集
func union(a, b []string) []string {
    set := make(map[string]bool)
    var result []string

    for _, item := range a {
        set[item] = true
    }
    for _, item := range b {
        set[item] = true
    }

    for item := range set {
        result = append(result, item)
    }

    return result
}

func (ic *IngredientController) RecommendIngredients(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    var request IngredientRecommendRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    // 1. 清理过期记录
    twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
    // 删除特定用户的过期食材选择历史记录
    ic.DB.Where("user_id = ? AND select_time < ?", userID, twentyFourHoursAgo).Delete(&models.UserIngredientHistory{})
    // 清理特定用户的过期食材偏好历史记录
    ic.DB.Where("user_id = ? AND update_time < ?", userID, twentyFourHoursAgo).Delete(&models.UserIngredientPreference{})

    
    // 2. 更新用户偏好
    for _, ingredientID := range request.LikedIngredients {
        preference := models.UserIngredientPreference{
            UserID: userID.(uint),
            IngredientID: ingredientID,
            IsLike: true,
            UpdateTime: time.Now(),
        }
        ic.DB.Save(&preference)
    }
    
    for _, ingredientID := range request.DislikedIngredients {
        preference := models.UserIngredientPreference{
            UserID: userID.(uint),
            IngredientID: ingredientID,
            IsLike: false,
            UpdateTime: time.Now(),
        }
        ic.DB.Save(&preference)
    }

    ingredientScores := make(map[uint]float64)

    // 获取所有食物名称
    foodNameResponse, err := models.GetAllFoodNames(ic.DB, "en")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food names"})
        return
    }
    
    // base score
    for _, food := range foodNameResponse {
        ingredientScores[food.ID] = BASE_SCORE
    }

    // 基于历史选择的食材
    var historyIngredients []models.UserIngredientHistory
    ic.DB.Where("user_id = ?", userID).Find(&historyIngredients)

    for _, history := range historyIngredients {
        ingredientScores[history.IngredientID] += WEIGHT_HISTORY
    }

    // 基于用户偏好，用户这次传输进来的偏好以及没有过期的偏好
    var preferences []models.UserIngredientPreference
    ic.DB.Where("user_id = ?", userID).Find(&preferences)

    for _, preference := range preferences {
        if preference.IsLike {
            ingredientScores[preference.IngredientID] += WEIGHT_PREFERENCE_LIKE
        } else {
            ingredientScores[preference.IngredientID] += WEIGHT_PREFERENCE_DISLIKE
        }
    }

    // 基于食物偏好类型，用户在设置页面的食物偏好
    // 获取设置界面用户的食物偏好类型
    foodPreferences, err := models.GetUserFoodPreferences(ic.DB, userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food preferences"})
        return
    }

    foodPos_id, foodNeg_id, err := ic.loadFoodPreferences(foodPreferences)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food preferences"})
        return
    }

    for _, id := range foodPos_id {
        ingredientScores[id] += WEIGHT_FOOD_PREF
    }

    // 构造分布，采样直到得到20个非负面食材
    var distribution []float64
    for _, score := range ingredientScores {
        distribution = append(distribution, score)
    }
    
    sampled := make([]uint, 0, 20)
    usedIDs := make(map[uint]bool)
    
    // 循环直到获得20个有效的食材
    for len(sampled) < 20 {
        // 采样新的食材
        newSamples := sample(distribution, 20-len(sampled))
        
        // 过滤掉重复的和负面的食材
        for _, id := range newSamples {
            // 跳过已使用的ID
            if usedIDs[id] {
                continue
            }
            // 跳过负面食材
            if slices.Contains(foodNeg_id, id) {
                continue
            }
            
            sampled = append(sampled, id)
            usedIDs[id] = true
            
            if len(sampled) >= 20 {
                break
            }
        }
    }

    // 构造推荐食材响应
    recommendedIngredients := make([]struct {
        ID   uint   `json:"id"`
        Name string `json:"name"`
    }, 0, len(sampled))

    for _, id := range sampled {
        var ingredient models.Ingredient
        if err := ic.DB.First(&ingredient, id).Error; err != nil {
            continue
        }
        recommendedIngredients = append(recommendedIngredients, struct {
            ID   uint   `json:"id"`
            Name string `json:"name"`
        }{ID: id, Name: ingredient.Name})
    }

    response := IngredientRecommendResponse{
        RecommendedIngredients: recommendedIngredients,
    }
    

    c.JSON(http.StatusOK, response)
}

func sample(distribution []float64, n int) []uint {
    total := 0.0
    for _, score := range distribution {
        total += score
    }
    
    // 归一化并采样
    result := make([]uint, 0, n)
    for i := 0; i < n; i++ {
        r := rand.Float64() * total
        sum := 0.0
        for id, score := range distribution {
            sum += score
            if sum >= r {
                result = append(result, uint(id))
                break
            }
        }
    }
    return result
}