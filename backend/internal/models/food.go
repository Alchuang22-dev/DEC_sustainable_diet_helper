package models

import (
    "gorm.io/gorm"
    "fmt"
)

type Food struct {
    gorm.Model
    ZhFoodName    string  `json:"zh_food_name" gorm:"column:zh_food_name;not null"`
    EnFoodName    string  `json:"en_food_name" gorm:"column:en_food_name;not null"`
    GHG           float64 `json:"ghg" gorm:"column:ghg"`
    Calories      float64 `json:"calories" gorm:"column:calories"`
    Protein       float64 `json:"protein" gorm:"column:protein"`
    Fat           float64 `json:"fat" gorm:"column:fat"`
    Carbohydrates float64 `json:"carbohydrates" gorm:"column:carbohydrates"`
    Sodium        float64 `json:"sodium" gorm:"column:sodium"`
    Price         float64 `json:"price" gorm:"column:price"`
    Recipes       []Recipe `json:"recipes" gorm:"many2many:food_recipes;"`
}

// TableName 指定表名
func (Food) TableName() string {
    return "foods"
}

// CreateFood 创建食物记录
func (f *Food) CreateFood(db *gorm.DB) error {
    return db.Create(f).Error
}

// GetFoodByID 通过ID获取食物
func GetFoodByID(db *gorm.DB, id uint) (*Food, error) {
    var food Food
    err := db.First(&food, id).Error
    return &food, err
}

// GetAllFoods 获取所有食物
func GetAllFoods(db *gorm.DB) ([]Food, error) {
    var foods []Food
    err := db.Find(&foods).Error
    return foods, err
}

// UpdateFood 更新食物信息
func (f *Food) UpdateFood(db *gorm.DB) error {
    return db.Save(f).Error
}

// DeleteFood 删除食物
func (f *Food) DeleteFood(db *gorm.DB) error {
    return db.Delete(f).Error
}

// FoodNameResponse 定义返回的食物名称结构
type FoodNameResponse struct {
    ID      uint   `json:"id"`
    ZhName string `json:"zh_name"` // 导出字段
    EnName string `json:"en_name"` // 导出字段
}

// GetAllFoodNames 获取所有食物名称
func GetAllFoodNames(db *gorm.DB) ([]FoodNameResponse, error) {
    var results []FoodNameResponse
    query := db.Model(&Food{})
    err := query.Select("id, zh_food_name as zh_name, en_food_name as en_name").Find(&results).Error
    return results, err
}

// FoodCalculateRequest 定义单个食物的计算请求
type FoodCalculateItem struct {
    ID     uint    `json:"id" binding:"required"`
    Price  float64 `json:"price" binding:"required"`
    Weight float64 `json:"weight" binding:"required"` // 单位：kg
}

// FoodCalculateResponse 定义单个食物的计算响应
type FoodCalculateResult struct {
    ID          uint    `json:"id"`
    CO2Emission float64 `json:"co2_emission"`
    Calories    float64 `json:"calories"`
    Protein     float64 `json:"protein"`
    Fat         float64 `json:"fat"`
    Carbs       float64 `json:"carbs"`
    Sodium      float64 `json:"sodium"`
}

// CalculateFoodNutrition 计算食物的营养成分和碳排放
func CalculateFoodNutritionAndEmission(db *gorm.DB, items []FoodCalculateItem) ([]FoodCalculateResult, error) {
    var results []FoodCalculateResult

    // 获取所有相关的食物ID
    var foodIDs []uint
    for _, item := range items {
        foodIDs = append(foodIDs, item.ID)
    }

    // 从数据库获取食物信息
    var foods []Food
    if err := db.Where("id IN ?", foodIDs).Find(&foods).Error; err != nil {
        return nil, err
    }

    // 创建食物ID到食物信息的映射
    foodMap := make(map[uint]Food)
    for _, food := range foods {
        foodMap[food.ID] = food
    }

    // 计算每个食物的结果
    for _, item := range items {
        food, exists := foodMap[item.ID]
        if !exists {
            return nil, fmt.Errorf("food with id %d not found", item.ID)
        }

        // 计算结果
        result := FoodCalculateResult{
            ID: item.ID,
            // CO2Emission = (GHG * weight * item.price) / food.price
            CO2Emission: (food.GHG * item.Weight * item.Price) / food.Price,
            // 其他营养成分直接按照重量比例计算
            Calories:    food.Calories * item.Weight,
            Protein:     food.Protein * item.Weight,
            Fat:         food.Fat * item.Weight,
            Carbs:      food.Carbohydrates * item.Weight,
            Sodium:      food.Sodium * item.Weight,
        }

        results = append(results, result)
    }

    return results, nil
}