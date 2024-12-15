package models

import (
    "gorm.io/gorm"
    "fmt"
    "math"
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
    ImageUrl      string  `json:"image_url" gorm:"column:image_url"`
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
// 通过食物名称获取食物ID
func FindFoodIDByName(db *gorm.DB, name string) (uint, error) {
    var food Food
    if err := db.Where("zh_food_name = ? OR en_food_name = ?", name, name).First(&food).Error; err != nil {
        return 0, err
    }
    return food.ID, nil
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
type FoodInfoResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
    ImageUrl string `json:"image_url"`
}

// GetAllFoodNames 获取所有食物名称
func GetAllFoodNames(db *gorm.DB, language string) ([]FoodInfoResponse, error) {
    var results []FoodInfoResponse
    query := db.Model(&Food{})
    
    switch language {
    case "zh":
        err := query.Select("id, zh_food_name as name, image_url").Find(&results).Error
        return results, err
    case "en":
        err := query.Select("id, en_food_name as name, image_url").Find(&results).Error
        return results, err
    default:
        return nil, fmt.Errorf("unsupported language: %s", language)
    }
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
    Emission float64 `json:"emission"`
    Calories    float64 `json:"calories"`
    Protein     float64 `json:"protein"`
    Fat         float64 `json:"fat"`
    Carbohydrates       float64 `json:"carbohydrates"`
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
            // 使用 math.Round 保留1位小数
            Emission: math.Round((food.GHG*item.Weight*item.Price/food.Price)*10) / 10,
            Calories: math.Round(food.Calories*item.Weight*10) / 10,
            Protein: math.Round(food.Protein*item.Weight*10) / 10,
            Fat: math.Round(food.Fat*item.Weight*10) / 10,
            Carbohydrates: math.Round(food.Carbohydrates*item.Weight*10) / 10,
            Sodium: math.Round(food.Sodium*item.Weight*10) / 10,
        }

        results = append(results, result)
    }

    return results, nil
}
type Ingredient struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
}

func (Ingredient) TableName() string {
    return "ingredients"
}