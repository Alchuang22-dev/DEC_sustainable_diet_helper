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
    WEIGHT_PREFERENCE_LIKE = 0.3  // 用户偏好权重(like)
    WEIGHT_PREFERENCE_DISLIKE = -0.2  // 用户偏好权重(dislike)
    WEIGHT_FOOD_PREF  = 0.3   // 食物偏好类型权重(like)

)

func (ic *IngredientController) loadFoodPreferences()(Pos_id []uint, Neg_id []uint, err error) {
    data, err := os.ReadFile("data/food_preference/foodPreferences.json")
    if err != nil {
        return nil, nil, err
    }

    var pres_name struct{
        FoodPos []string `json:"food_pos"`
        FoodNeg []string `json:"food_neg"`
    }
    
    if err := json.Unmarshal(data, &pres_name); err != nil {
        return nil, nil, err
    }

    // find id by name 
    var foodPos_id, foodNeg_id []uint
    for _, food := range pres_name.FoodPos {
        var ingredient models.Ingredient
        if err := ic.DB.Where("name = ?", food).First(&ingredient).Error; err != nil {
            return nil, nil, err
        }
        foodPos_id = append(foodPos_id, ingredient.ID)
    }
    for _, food := range pres_name.FoodNeg {
        var ingredient models.Ingredient
        if err := ic.DB.Where("name = ?", food).First(&ingredient).Error; err != nil {
            return nil, nil, err
        }
        foodNeg_id = append(foodNeg_id, ingredient.ID)
    }

    return foodPos_id, foodNeg_id, nil
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
    // 删除用户食材选择历史记录
    ic.DB.Where("select_time < ?", twentyFourHoursAgo).Delete(&models.UserIngredientHistory{})
    // 清楚用户食材偏好历史记录
    ic.DB.Where("update_time < ?", twentyFourHoursAgo).Delete(&models.UserIngredientPreference{})

    
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

    // 基于历史选择
    var historyIngredients []models.UserIngredientHistory
    ic.DB.Where("user_id = ?", userID).Find(&historyIngredients)

    for _, history := range historyIngredients {
        ingredientScores[history.IngredientID] += WEIGHT_HISTORY
    }

    // 基于用户偏好
    var preferences []models.UserIngredientPreference
    ic.DB.Where("user_id = ?", userID).Find(&preferences)

    for _, preference := range preferences {
        if preference.IsLike {
            ingredientScores[preference.IngredientID] += WEIGHT_PREFERENCE_LIKE
        } else {
            ingredientScores[preference.IngredientID] += WEIGHT_PREFERENCE_DISLIKE
        }
    }

    // 基于食物偏好类型
    foodPos_id, foodNeg_id, err := ic.loadFoodPreferences()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food preferences"})
        return
    }

    for _, id := range foodPos_id {
        ingredientScores[id] += WEIGHT_FOOD_PREF
    }

    // 构造分布，采样10个
    var distribution []float64
    for _, score := range ingredientScores {
        distribution = append(distribution, score)
    }
    sampled := sample(distribution, 20)

    // 判断foodNeg_id是否在sampled中
    for _, id := range foodNeg_id {
        if slices.Contains(sampled, id) {
            sampled = append(sampled, id)
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