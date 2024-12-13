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
    "log"
)

type RecommendController struct {
    DB *gorm.DB
}

// 食材推荐请求结构体
type IngredientRecommendRequest struct {
    UseLastIngredients bool     `json:"use_last_ingredients"`
    LikedIngredients   []uint   `json:"liked_ingredients"`
    DislikedIngredients []uint  `json:"disliked_ingredients"`
}

// 食谱推荐以及设置用户最近选择的食材请求结构体
type RecipeRecommendAndSetUserLastSelectedFoodsRequest struct {
    SelectedIngredients []uint `json:"selected_ingredients"`
    DislikedIngredients []uint `json:"disliked_ingredients"`
}

// 食谱推荐响应结构体
type RecipeRecommendResponse struct {
    // 返回食谱的名称，图片url，食谱id，食谱的组成
    RecommendedRecipes []struct {
        Name string `json:"name"`
        ImageURL string `json:"image_url"`
        RecipeID uint `json:"recipe_id"`
        Ingredients string `json:"ingredients"`
    } `json:"recommended_recipes"`
}

// 食材得分结构体
type foodScore struct {
    id    uint
    score float64
}

// 食材推荐响应结构体
type IngredientRecommendResponse struct {
    RecommendedIngredients []struct {
        ID          uint    `json:"id"`
        Name        string  `json:"name"`
        ImageURL    string  `json:"image_url"`
    } `json:"recommended_ingredients"`
}


// 食材推荐分数常量
const (
    BASE_SCORE       = 0.5   // 基础分数
    WEIGHT_HISTORY   = 0.2   // 历史选择权重
    WEIGHT_PREFERENCE_LIKE = 0.3  // 用户此次传输的偏好权重(like)
    WEIGHT_PREFERENCE_DISLIKE = -0.2  // 用户此次传输的偏好权重(dislike)
    WEIGHT_FOOD_PREF  = 0.3   // 设置界面的食物偏好类型权重(like)

)

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


// 辅助函数：加载食物偏好
func (ic *RecommendController) loadFoodPreferences(preferences []models.FoodPreference) (Pos_id []uint, Neg_id []uint, err error) {
    data, err := os.ReadFile("../../data/food_preference/foodPreferences.json")
    if err != nil {
        return nil, nil, err
    }
    log.Printf("加载食物偏好成功")

    var preferencesMap map[string]struct {
        FoodPos []string `json:"food_pos"`
        FoodNeg []string `json:"food_neg"`
    }
    if err := json.Unmarshal(data, &preferencesMap); err != nil {
        return nil, nil, err
    }
    log.Printf("解析食物偏好成功")
    // 用于存储所有偏好的食物列表
    var allPosFood []string
    var allNegFood []string
    isFirst := true

    // 处理每个用户选择的偏好
    log.Printf("处理每个用户选择的偏好")
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
    log.Printf("处理每个用户选择的偏好成功")

    // 将食物名称转换为ID
    foodPos_id := make([]uint, 0)
    foodNeg_id := make([]uint, 0)

    // 处理positive foods
    for _, food := range allPosFood {
        var foodItem models.Food
        if err := ic.DB.Table("foods").Where("en_food_name = ?", food).First(&foodItem).Error; err != nil {
            continue
        }
        foodPos_id = append(foodPos_id, foodItem.ID)
    }
    log.Printf("处理positive foods成功")

    // 处理 negative foods
    for _, food := range allNegFood {
        var foodItem models.Food
        if err := ic.DB.Table("foods").Where("en_food_name = ?", food).First(&foodItem).Error; err != nil {
            continue
        }
        foodNeg_id = append(foodNeg_id, foodItem.ID)
    }
    log.Printf("处理negative foods成功")
    return foodPos_id, foodNeg_id, nil
}

// 辅助函数：采样食材
func sample(foodScores []foodScore, n int) []uint {
    if len(foodScores) == 0 {
        return nil
    }
    
    total := 0.0
    for _, item := range foodScores {
        total += item.score
    }
    
    result := make([]uint, 0, n)
    remainingScores := make([]foodScore, len(foodScores))
    copy(remainingScores, foodScores)
    
    for i := 0; i < n && len(remainingScores) > 0; i++ {
        if total <= 0 {
            break
        }
        
        r := rand.Float64() * total
        cumSum := 0.0
        
        for j, item := range remainingScores {
            cumSum += item.score
            if cumSum >= r {
                result = append(result, item.id)
                // 更新总分和剩余食材
                total -= item.score
                remainingScores = append(remainingScores[:j], remainingScores[j+1:]...)
                break
            }
        }
    }
    
    return result
}

// 推荐食材
func (ic *RecommendController) RecommendIngredients(c *gin.Context) {
    log.Printf("开始推荐食材")
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
        return
    }
    log.Printf("用户ID: %v", userID)
    
    var request IngredientRecommendRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    // 添加参数验证
    if len(request.LikedIngredients) > 100 || len(request.DislikedIngredients) > 100 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Too many ingredients selected"})
        return
    }
    log.Printf("request: %v", request)

    // 如果使用上一次的食材，则直接返回上次的食材
    if request.UseLastIngredients {
        var lastSelectedFoods []models.UserLastSelectedFoods
        if err := ic.DB.Where("user_id = ?", userID).Find(&lastSelectedFoods).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last selected foods"})
            return
        }
        
        // 构造推荐食材响应
        recommendedIngredients := make([]struct {
            ID   uint   `json:"id"`
            Name string `json:"name"`
            ImageURL string `json:"image_url"`
        }, 0, len(lastSelectedFoods))

        // 获取每个食材的详细信息
        for _, food := range lastSelectedFoods {
            var foodInfo models.Food
            if err := ic.DB.First(&foodInfo, food.FoodID).Error; err != nil {
                continue
            }
            recommendedIngredients = append(recommendedIngredients, struct {
                ID   uint   `json:"id"`
                Name string `json:"name"`
                ImageURL string `json:"image_url"`
            }{
                ID:   foodInfo.ID,
                Name: foodInfo.EnFoodName,
                ImageURL: foodInfo.ImageUrl,
            })
        }

        response := IngredientRecommendResponse{
            RecommendedIngredients: recommendedIngredients,
        }
        
        c.JSON(http.StatusOK, response)
        return
    }

    // 验证食材ID是否存在
    for _, id := range append(request.LikedIngredients, request.DislikedIngredients...) {
        var food models.Food
        if err := ic.DB.Table("foods").First(&food, id).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
            return
        }
    }
    log.Printf("验证食材ID成功")

    // 1. 清理过期记录
    tx := ic.DB.Begin()
    twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
    if err := tx.Where("user_id = ? AND select_time < ?", userID, twentyFourHoursAgo).Delete(&models.UserIngredientHistory{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "清理历史记录失败"})
        return
    }
    log.Printf("清理历史记录成功")
    if err := tx.Where("user_id = ? AND update_time < ?", userID, twentyFourHoursAgo).Delete(&models.UserIngredientPreference{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "清理偏好记录失败"})
        return
    }
    tx.Commit()
    log.Printf("清理偏好记录成功")

    
    // 2. 更新用户偏好
    tx = ic.DB.Begin()
    for _, ingredientID := range request.LikedIngredients {
        preference := models.UserIngredientPreference{
            UserID: userID.(uint),
            IngredientID: ingredientID,
            IsLike: true,
            UpdateTime: time.Now(),
        }
        if err := tx.Save(&preference).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save preferences"})
            return
        }
    }
    log.Printf("更新用户偏好成功")
    for _, ingredientID := range request.DislikedIngredients {
        preference := models.UserIngredientPreference{
            UserID: userID.(uint),
            IngredientID: ingredientID,
            IsLike: false,
            UpdateTime: time.Now(),
        }
        if err := tx.Save(&preference).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save preferences"})
            return
        }
    }
    tx.Commit()
    log.Printf("更新用户偏好成功")


    // 获取所有食物名称
    foodNameResponse, err := models.GetAllFoodNames(ic.DB, "en")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food names"})
        return
    }
    log.Printf("获取所有食物名称成功")
    // base score
    ingredientScores := make(map[uint]float64)
    for _, food := range foodNameResponse {
        ingredientScores[food.ID] = BASE_SCORE
    }

    // 基于历史选择的食材
    var historyIngredients []models.UserIngredientHistory
    if err := ic.DB.Where("user_id = ?", userID).Find(&historyIngredients).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get history ingredients"})
        return
    }
    log.Printf("获取历史食材成功")
    for _, history := range historyIngredients {
        ingredientScores[history.IngredientID] += WEIGHT_HISTORY
    }

    // 基于用户偏好，用户这次传输进来的偏好以及没有过期的偏好
    var preferences []models.UserIngredientPreference
    if err := ic.DB.Where("user_id = ?", userID).Find(&preferences).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get preferences"})
        return
    }
    log.Printf("获取偏好食材成功")
    for _, preference := range preferences {
        if preference.IsLike {
            ingredientScores[preference.IngredientID] += WEIGHT_PREFERENCE_LIKE
        } else {
            ingredientScores[preference.IngredientID] += WEIGHT_PREFERENCE_DISLIKE
        }
    }
    log.Printf("更新偏好食材成功")

    // 基于食物偏好类型，用户在设置页面的食物偏好

    // 获取设置界面用户的食物偏好类型
    foodPreferences, err := models.GetUserFoodPreferences(ic.DB, userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food preferences"})
        return
    }
    log.Printf("获取设置界面用户的食物偏好类型成功")
    foodPos_id, foodNeg_id, err := ic.loadFoodPreferences(foodPreferences)
    if err != nil {
        log.Printf("获取食物偏好失败: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food preferences"})
        return
    }

    for _, id := range foodPos_id {
        ingredientScores[id] += WEIGHT_FOOD_PREF
    }
    log.Printf("更新食物偏好类型成功")

    // 将map转换为带ID的切片
    foodScores := make([]foodScore, 0, len(ingredientScores))
    for id, score := range ingredientScores {
        foodScores = append(foodScores, foodScore{id: id, score: score})
    }

    if len(foodScores) == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No valid ingredients found"})
        return
    }
    
    sampled := make([]uint, 0, 20)
    usedIDs := make(map[uint]bool)
    
    // 循环直到获得20个有效的食材
    maxAttempts := 100
    attempts := 0
    for len(sampled) < 20 && attempts < maxAttempts {
        newSamples := sample(foodScores, 20-len(sampled))
        attempts++
        
        for _, id := range newSamples {
            if usedIDs[id] || slices.Contains(foodNeg_id, id) {
                continue
            }
            sampled = append(sampled, id)
            usedIDs[id] = true
            
            if len(sampled) >= 20 {
                break
            }
        }
    }
    log.Printf("采样食材成功")
    // 构造推荐食材响应
    recommendedIngredients := make([]struct {
        ID   uint   `json:"id"`
        Name string `json:"name"`
        ImageURL string `json:"image_url"`
    }, 0, len(sampled))
    log.Printf("sampled: %v", sampled)

    for _, id := range sampled {
        var food models.Food
        if err := ic.DB.Table("foods").First(&food, id).Error; err != nil {
            log.Printf("获取食材失败: %v", err)
            continue
        }
        recommendedIngredients = append(recommendedIngredients, struct {
            ID   uint   `json:"id"`
            Name string `json:"name"`
            ImageURL string `json:"image_url"`
        }{
            ID:   food.ID,
            Name: food.EnFoodName,
            ImageURL: food.ImageUrl,
        })
    }
    if len(recommendedIngredients) == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No recommended ingredients"})
        return
    }
    log.Printf("构造推荐食材响应成功")
    response := IngredientRecommendResponse{
        RecommendedIngredients: recommendedIngredients,
    }
    
    c.JSON(http.StatusOK, response)
}

// 推荐菜谱
func (ic *RecommendController) RecommendRecipes(c *gin.Context) {
    log.Printf("开始推荐菜谱")
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
        return
    }
    log.Printf("用户ID: %v", userID)
    // 清楚1周以上的历史数据
    tx := ic.DB.Begin()
    oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)
    if err := tx.Where("user_id = ? AND select_time < ?", userID, oneWeekAgo).Delete(&models.UserRecipeHistory{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear history data"})
        return
    }
    tx.Commit()
    log.Printf("清除历史数据成功")

    // 获取用户历史菜谱
    var historyRecipes []models.UserRecipeHistory
    if err := ic.DB.Where("user_id = ?", userID).Find(&historyRecipes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get history recipes"})
        return
    }
    log.Printf("获取历史菜谱成功")

    // 获取请求体
    var request RecipeRecommendAndSetUserLastSelectedFoodsRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    log.Printf("获取请求体成功")

    // 验证SelectedIngredients不为空
    if len(request.SelectedIngredients) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No ingredients selected"})
        return
    }
    
    // 获取设置界面用户的食物偏好类型
    foodPreferences, err := models.GetUserFoodPreferences(ic.DB, userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food preferences"})
        return
    }
    log.Printf("获取设置界面用户的食物偏好类型成功")
    
    _, foodNeg_id, err := ic.loadFoodPreferences(foodPreferences)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get food preferences"})
        return
    }

    // 合并foodNeg_id 与 request.DislikedIngredients
    foodNeg_id = append(foodNeg_id, request.DislikedIngredients...)



    // 对于每一个食材id，获取包含该食材的食谱
    var recipes []models.Recipe
    maxAttempts := 20
    for _, ingredientID := range request.SelectedIngredients {
        // 跳过负面食材
        if slices.Contains(foodNeg_id, ingredientID) {
            continue
        }

        // 验证ingredientID是否存在
        var food models.Food
        if err := ic.DB.First(&food, ingredientID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
            return
        }

        recipeIDs, err := models.GetRecipeIDsByFoodID(ic.DB, ingredientID)
        if err != nil {
            log.Printf("获取包含食材%v的食谱失败: %v", ingredientID, err)
            continue
        }
        log.Printf("获取包含食材%v的食谱成功", ingredientID)

        if len(recipeIDs) > 0 {
            attempts := 0
            validRecipe := false

            availableRecipes := make([]uint, len(recipeIDs))
            copy(availableRecipes, recipeIDs)

            for attempts < maxAttempts && !validRecipe && len(availableRecipes) > 1 {
                randomIndex := rand.Intn(len(availableRecipes))
                recipeID := availableRecipes[randomIndex]
                recipe, err := models.GetRecipeByID(ic.DB, recipeID)
                if err != nil {
                    log.Printf("获取食谱失败: recipeID=%v, error=%v", recipeID, err)
                    availableRecipes = append(availableRecipes[:randomIndex], availableRecipes[randomIndex+1:]...)
                    attempts++
                    continue
                }
                // 获取食谱关联的食材
                ingredientIds, err := models.GetIngredientIDsByRecipeID(ic.DB, recipe.ID)
                if err != nil {
                    log.Printf("获取食谱关联的食材失败: %v", err)
                    availableRecipes = append(availableRecipes[:randomIndex], availableRecipes[randomIndex+1:]...)
                    attempts++
                    continue
                }

                containsNeg := false
                for _, id := range ingredientIds {
                    if slices.Contains(foodNeg_id, id) {
                        containsNeg = true
                        break
                    }
                }

                if !containsNeg {
                    validRecipe = true
                    recipes = append(recipes, *recipe)
                }

                availableRecipes = append(availableRecipes[:randomIndex], availableRecipes[randomIndex+1:]...)
                attempts++
            }

            if !validRecipe {
                log.Printf("无法找到有效的食谱, 选取最后一个可能的食谱")
                // 数组长度检查
                if len(availableRecipes) == 0 {
                    log.Printf("没有可用的食谱")
                    continue
                }
                recipe, err := models.GetRecipeByID(ic.DB, availableRecipes[0])
                if err != nil {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recipe"})
                    return
                }
                recipes = append(recipes, *recipe)
            }
        }
        log.Printf("获取包含食材%v的食谱成功", ingredientID)
    }

    // 构造推荐菜谱响应
    recommendedRecipes := make([]struct {
        Name        string `json:"name"`
        ImageURL    string `json:"image_url"`
        RecipeID    uint   `json:"recipe_id"`
        Ingredients string `json:"ingredients"`
    }, 0, len(recipes))

    for _, recipe := range recipes {
        recommendedRecipes = append(recommendedRecipes, struct {
            Name        string `json:"name"`
            ImageURL    string `json:"image_url"`
            RecipeID    uint   `json:"recipe_id"`
            Ingredients string `json:"ingredients"`
        }{
            Name:        recipe.Name,
            ImageURL:    recipe.ImageURL,
            RecipeID:    recipe.ID,
            Ingredients: recipe.Ingredients,
        })
    }
    log.Printf("构造推荐菜谱响应成功")

    response := RecipeRecommendResponse{
        RecommendedRecipes: recommendedRecipes,
    }

    c.JSON(http.StatusOK, response)
}

// set user selected foods
func (ic *RecommendController) SetUserSelectedFoods(c *gin.Context) {
    log.Printf("开始设置用户选择的食材")
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
        return
    }

    // 获取请求体
    var request RecipeRecommendAndSetUserLastSelectedFoodsRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    // 验证SelectedIngredients不为空
    if len(request.SelectedIngredients) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No ingredients selected"})
        return
    }
    
    // 验证SelectedIngredients中的食材ID是否存在
    for _, ingredientID := range request.SelectedIngredients {
        var food models.Food
        if err := ic.DB.First(&food, ingredientID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
            return
        }
    }

    // 设置用户最近选择的食材
    tx := ic.DB.Begin()
    for _, ingredientID := range request.SelectedIngredients {
        tx.Create(&models.UserLastSelectedFoods{UserID: userID.(uint), FoodID: ingredientID, SelectTime: time.Now()})
    }
    tx.Commit()
    log.Printf("设置用户最近选择的食材成功")

    c.JSON(http.StatusOK, gin.H{"message": "User selected foods set successfully"})
}