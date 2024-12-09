package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// UserShare 结构体
type UserShare struct {
    UserID  uint `json:"user_id" binding:"required"`
    Ratio float64 `json:"ratio" binding:"required,min=0,max=1"`
}

// SharedNutritionIntakeRequest 共享营养摄入请求
type SharedNutritionCarbonIntakeRequest struct {
    Date          time.Time   `json:"date" binding:"required"`
    MealType      models.MealType    `json:"meal_type" binding:"required"`
    Calories      float64     `json:"calories" binding:"min=0"`
    Protein       float64     `json:"protein" binding:"min=0"`
    Fat           float64     `json:"fat" binding:"min=0"`
    Carbohydrates float64     `json:"carbohydrates" binding:"min=0"`
    Sodium        float64     `json:"sodium" binding:"min=0"`
    Emission      float64     `json:"emission" binding:"min=0"`
    UserShares    []UserShare `json:"user_shares" binding:"required,min=1"`
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

    // 处理每个目标
    for _, goal := range goals {
        // 检查该日期是否已存在目标
        var existingGoal models.NutritionGoal
        result := tx.Where("user_id = ? AND date = ?", userID, goal.Date).First(&existingGoal)
        
        if result.Error == nil {
            // 目标已存在，更新目标
            existingGoal.Calories = goal.Calories
            existingGoal.Protein = goal.Protein
            existingGoal.Fat = goal.Fat
            existingGoal.Carbohydrates = goal.Carbohydrates
            existingGoal.Sodium = goal.Sodium
            
            if err := tx.Save(&existingGoal).Error; err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"error": "更新目标失败"})
                return
            }
        } else {
            // 目标不存在，创建新目标
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

    // 处理每个目标
    for _, goal := range goals {
        // 检查该日期是否已存在目标
        var existingGoal models.CarbonGoal
        result := tx.Where("user_id = ? AND date = ?", userID, goal.Date).First(&existingGoal)
        
        if result.Error == nil {
            // 目标已存在，更新目标
            existingGoal.Emission = goal.Emission
            
            if err := tx.Save(&existingGoal).Error; err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"error": "更新目标失败"})
                return
            }
        } else {
            // 目标不存在，创建新目标
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

// validateUserShares 验证用户分摊信息
func (nc *NutritionCarbonController) validateUserShares(currentUserID uint, shares []UserShare) (bool, error) {
    // 获取当前用户家庭信息
    var user models.User
    if err := nc.DB.Preload("Family.Members").First(&user, currentUserID).Error; err != nil {
        return false, fmt.Errorf("获取用户信息失败")
    }

    if user.Family == nil {
        return false, fmt.Errorf("用户不属于任何家庭")
    }

    // 验证所有用户是否属于同一个家庭
    familyMembers := make(map[uint]bool)
    for _, member := range user.Family.Members {
        familyMembers[member.ID] = true
    }

    // 先验证每个比例值是否有效
    for _, share := range shares {
        if share.Ratio <= 0 || share.Ratio > 1 {
            return false, fmt.Errorf("无效的请求数据")
        }
    }

    // 再验证比例总和是否为1
    var totalRatio float64
    for _, share := range shares {
        if !familyMembers[share.UserID] {
            return false, fmt.Errorf("用户 %d 不属于同一个家庭", share.UserID)
        }
        totalRatio += share.Ratio
    }

    // 允许有0.00001的误差
    if totalRatio < 0.99999 || totalRatio > 1.00001 {
        return false, fmt.Errorf("分摊比例之和必须等于1")
    }

    return true, nil
}

// SetSharedNutritionCarbonIntake 设置共享营养碳排放
func (nc *NutritionCarbonController) SetSharedNutritionCarbonIntake(c *gin.Context) {
    userID, err := c.Get("user_id")
    if !err {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var req SharedNutritionCarbonIntakeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    // 验证用户分摊信息
    if valid, err := nc.validateUserShares(userID.(uint), req.UserShares); !valid {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 开启事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "开启事务失败"})
        return
    }

    // 为每个用户创建营养和碳摄入记录
    for _, share := range req.UserShares {
        nutritionIntake := models.NutritionIntake{
            UserID:        share.UserID,
            Date:         req.Date,
            MealType:     req.MealType,
            Calories:     req.Calories * share.Ratio,
            Protein:      req.Protein * share.Ratio,
            Fat:         req.Fat * share.Ratio,
            Carbohydrates: req.Carbohydrates * share.Ratio,
            Sodium:       req.Sodium * share.Ratio,
        }

        carbonIntake := models.CarbonIntake{
            UserID: share.UserID,
            Date: req.Date,
            MealType: req.MealType,
            Emission: req.Emission * share.Ratio,
        }

        if err := tx.Create(&nutritionIntake).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "创建营养摄入记录失败"})
            return
        }

        if err := tx.Create(&carbonIntake).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "创建碳排放记录失败"})
            return
        }
    }

    tx.Commit()
    c.JSON(http.StatusOK, gin.H{"message": "共享营养碳排放记录创建成功"})
}