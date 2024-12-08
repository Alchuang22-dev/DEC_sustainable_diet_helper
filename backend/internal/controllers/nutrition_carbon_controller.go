package controllers

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "fmt"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
)

type NutritionCarbonController struct {
    DB *gorm.DB
}
// nutrition_goal_request 结构体
type NutritionGoalRequest struct {
    Date          time.Time `json:"date" binding:"required"`
    Calories      float64   `json:"calories" binding:"min=0"`
    Protein       float64   `json:"protein" binding:"min=0"`
    Fat           float64   `json:"fat" binding:"min=0"`
    Carbohydrates float64   `json:"carbohydrates" binding:"min=0"`
    Sodium        float64   `json:"sodium" binding:"min=0"`
}

// carbon_goal_request 结构体
type CarbonGoalRequest struct {
    Date time.Time `json:"date" binding:"required"`
    Emission float64 `json:"emission" binding:"min=0"`
}

// 验证日期,需要保证起始是今天，且连续往后
func validateDate(data []time.Time) (bool, error) {
    today := time.Now().Truncate(24 * time.Hour)
    for i, date := range data {
        if date.Before(today) {
            return false, fmt.Errorf("日期不能早于今天")
        }
        if i > 0 && date.Sub(data[i-1]) != time.Hour * 24 {
            return false, fmt.Errorf("日期不连续")
        }
    }
    return true, nil
}
// 设置营养目标
func (nc *NutritionCarbonController) SetNutritionGoals(c *gin.Context){
    userID, err := c.Get("user_id")
    if !err {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 请求体
    var goals []NutritionGoalRequest

    if err := c.ShouldBindJSON(&goals); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 验证日期
    var dates []time.Time
    for _, goal := range goals {
        dates = append(dates, goal.Date)
    }
    if ok, err := validateDate(dates); !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 开启事务
    tx := nc.DB.Begin()

    // 删除一周前的目标
    if err := tx.Where("user_id = ? AND date < ?", userID, time.Now().AddDate(0, 0, -7)).Delete(&models.NutritionGoal{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除一周前的目标失败"})
        return
    }

    // 创建新目标
    for _, goal := range goals {
        newGoal := models.NutritionGoal{
            UserID: userID.(uint),
            Date: goal.Date,
            Calories: goal.Calories,
            Protein: goal.Protein,
            Fat: goal.Fat,
            Carbohydrates: goal.Carbohydrates,
            Sodium: goal.Sodium,
        }
        if err := tx.Create(&newGoal).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目标失败"})
            return
        }
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "目标设置成功"})
}

// 设置碳排放目标
func (nc *NutritionCarbonController) SetCarbonGoals(c *gin.Context) {
    userID, err := c.Get("user_id")
    if !err {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 请求体
    var goals []CarbonGoalRequest

    if err := c.ShouldBindJSON(&goals); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 验证日期
    var dates []time.Time
    for _, goal := range goals {
        dates = append(dates, goal.Date)
    }
    if ok, err := validateDate(dates); !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 开启事务
    tx := nc.DB.Begin()

    // 删除一周前的目标 
    if err := tx.Where("user_id = ? AND date < ?", userID, time.Now().AddDate(0, 0, -7)).Delete(&models.CarbonGoal{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除一周前的目标失败"})
        return
    }

    // 创建新目标
    for _, goal := range goals {
        newGoal := models.CarbonGoal{
            UserID: userID.(uint),
            Date: goal.Date,
            Emission: goal.Emission,
        }
        if err := tx.Create(&newGoal).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目标失败"})
            return
        }
    } 

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "目标设置成功"})
}   

// 获取营养目标
func (nc *NutritionCarbonController) GetNutritionGoals(c *gin.Context) {
    userID, err := c.Get("user_id")
    if !err {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 删除一周前的目标
    if err := nc.DB.Where("user_id = ? AND date < ?", userID, time.Now().AddDate(0, 0, -7)).Delete(&models.NutritionGoal{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除一周前的目标失败"})
        return
    }

    // 计算时间范围
    now := time.Now()
    startDate := now.AddDate(0, 0, -6)
    endDate := now.AddDate(0, 0, 1)

    // 创建一个包含8天的目标数组(从6天前到明天)
    goals := make([]models.NutritionGoal, 8)
    
    // 初始化每一天的基础数据
    for i := 0; i < 8; i++ {
        goals[i] = models.NutritionGoal{
            UserID: userID.(uint),
            Date: now.AddDate(0, 0, i-6), // 从6天前开始,到明天
            Calories: 0,
            Protein: 0,
            Fat: 0,
            Carbohydrates: 0,
            Sodium: 0,
        }
    }

    // 查询存在的目标数据
    var existingGoals []models.NutritionGoal
    if err := nc.DB.Where("user_id = ? AND date BETWEEN ? AND ?", 
        userID, startDate, endDate).
        Order("date").
        Find(&existingGoals).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取营养目标失败"})
        return
    }

    // 用存在的数据覆盖对应日期的默认值
    for _, existingGoal := range existingGoals {
        dayDiff := existingGoal.Date.Sub(startDate).Hours() / 24
        if dayIndex := int(dayDiff); dayIndex >= 0 && dayIndex < 8 { 
            goals[dayIndex] = existingGoal
        }
    }

    c.JSON(http.StatusOK, gin.H{"data": goals})
}

// 获取碳排放目标
func (nc *NutritionCarbonController) GetCarbonGoals(c *gin.Context) {
    userID, err := c.Get("user_id")
    if !err {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 删除一周前的目标
    if err := nc.DB.Where("user_id = ? AND date < ?", userID, time.Now().AddDate(0, 0, -7)).Delete(&models.CarbonGoal{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除一周前的目标失败"})
        return
    }

    // 计算时间范围
    now := time.Now()
    startDate := now.AddDate(0, 0, -6)
    endDate := now.AddDate(0, 0, 1)

    // 创建一个包含8天的目标数组(从6天前到明天)
    goals := make([]models.CarbonGoal, 8)

    // 初始化每一天的基础数据
    for i := 0; i < 8; i++ {
        goals[i] = models.CarbonGoal{
            UserID: userID.(uint),
            Date: now.AddDate(0, 0, i-6), // 从6天前开始,到明天
            Emission: 0,
        }
    }

    // 查询存在的目标数据
    var existingGoals []models.CarbonGoal
    if err := nc.DB.Where("user_id = ? AND date BETWEEN ? AND ?", 
        userID, startDate, endDate).
        Order("date").
        Find(&existingGoals).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取碳排放目标失败"})
        return
    }

    // 用存在的数据覆盖对应日期的默认值
    for _, existingGoal := range existingGoals {
        dayDiff := existingGoal.Date.Sub(startDate).Hours() / 24
        if dayIndex := int(dayDiff); dayIndex >= 0 && dayIndex < 8 {
            goals[dayIndex] = existingGoal
        }
    }

    c.JSON(http.StatusOK, gin.H{"data": goals})
}

// 获取实际营养摄入
func (nc *NutritionCarbonController) GetActualNutrition(c *gin.Context) {
    userID, err := c.Get("user_id")
    if !err {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 删除7天前的摄入记录
    if err := nc.DB.Where("user_id = ? AND date < ?", userID, time.Now().AddDate(0, 0, -7)).Delete(&models.NutritionIntake{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除7天前的摄入记录失败"})
        return
    }

    // 计算时间范围
    now := time.Now()
    startDate := now.AddDate(0, 0, -6)
    endDate := now

    // 创建一个包含7天x4餐的默认摄入记录
    defaultIntakes := make([]models.NutritionIntake, 7*4)
    mealTypes := []models.MealType{models.Breakfast, models.Lunch, models.Dinner, models.Other}
    
    // 初始化每天每餐的默认值
    for i := 0; i < 7; i++ {
        currentDate := startDate.AddDate(0, 0, i)
        for j, mealType := range mealTypes {
            defaultIntakes[i*4+j] = models.NutritionIntake{
                UserID:        userID.(uint),
                Date:         currentDate,
                MealType:     mealType,
                Calories:     0,
                Protein:      0,
                Fat:         0,
                Carbohydrates: 0,
                Sodium:       0,
            }
        }
    }

    // 查询实际摄入记录
    var actualIntakes []models.NutritionIntake
    if err := nc.DB.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
        Order("date, meal_type").
        Find(&actualIntakes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取营养摄入记录失败"})
        return
    }

    // 累加实际摄入值到对应的默认记录中
    for _, actual := range actualIntakes {
        dayDiff := actual.Date.Sub(startDate).Hours() / 24
        dayIndex := int(dayDiff)
        
        var mealIndex int
        switch actual.MealType {
        case models.Breakfast:
            mealIndex = 0
        case models.Lunch:
            mealIndex = 1
        case models.Dinner:
            mealIndex = 2
        case models.Other:
            mealIndex = 3
        }
        
        index := dayIndex*4 + mealIndex
        if index >= 0 && index < len(defaultIntakes) {
            defaultIntakes[index].Calories += actual.Calories
            defaultIntakes[index].Protein += actual.Protein
            defaultIntakes[index].Fat += actual.Fat
            defaultIntakes[index].Carbohydrates += actual.Carbohydrates
            defaultIntakes[index].Sodium += actual.Sodium
        }
    }

    c.JSON(http.StatusOK, defaultIntakes)
}

// 获取实际碳排放
func (nc *NutritionCarbonController) GetCarbonIntakes(c *gin.Context) {
    userID, err := c.Get("user_id")
    if !err {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 删除7天前的碳排放记录
    if err := nc.DB.Where("user_id = ? AND date < ?", userID, time.Now().AddDate(0, 0, -7)).Delete(&models.CarbonIntake{}).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "删除7天前的碳排放记录失败"})
        return
    }

    // 计算时间范围
    now := time.Now()
    startDate := now.AddDate(0, 0, -6)
    endDate := now

    // 创建一个包含7天x4餐的默认碳排放记录
    defaultIntakes := make([]models.CarbonIntake, 7*4)
    mealTypes := []models.MealType{models.Breakfast, models.Lunch, models.Dinner, models.Other}

    // 初始化每天每餐的默认值
    for i := 0; i < 7; i++ {
        currentDate := startDate.AddDate(0, 0, i)
        for j, mealType := range mealTypes {
            defaultIntakes[i*4+j] = models.CarbonIntake{
                UserID: userID.(uint),
                Date: currentDate,
                MealType: mealType,
                Emission: 0,
            }
        }
    }

    // 查询实际碳排放记录
    var actualIntakes []models.CarbonIntake
    if err := nc.DB.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
        Order("date, meal_type").
        Find(&actualIntakes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取碳排放记录失败"})
        return
    }

    // 累加实际摄入值到对应的默认记录中
    for _, actual := range actualIntakes {
        dayDiff := actual.Date.Sub(startDate).Hours() / 24
        dayIndex := int(dayDiff)
        
        var mealIndex int
        switch actual.MealType {
        case models.Breakfast:
            mealIndex = 0
        case models.Lunch:
            mealIndex = 1
        case models.Dinner:
            mealIndex = 2
        case models.Other:
            mealIndex = 3
        }
        
        index := dayIndex*4 + mealIndex
        if index >= 0 && index < len(defaultIntakes) {
            defaultIntakes[index].Emission += actual.Emission
        }
    }

    c.JSON(http.StatusOK, defaultIntakes)
}