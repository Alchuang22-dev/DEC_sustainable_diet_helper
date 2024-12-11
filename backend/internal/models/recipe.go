package models

import (
	"gorm.io/gorm"
	"encoding/json"
	"fmt"
)

// Recipe 食谱模型
type Recipe struct {
	gorm.Model
	URL         string  `json:"url" gorm:"column:url;not null"`         // 食谱来源网址
	Name        string  `json:"name" gorm:"column:name;not null"`       // 食谱名称
	ImageURL    string  `json:"image_url" gorm:"column:image_url"`      // 食谱图片URL
	Ingredients string  `json:"ingredients" gorm:"column:ingredients"`   // 原料组成(JSON格式存储)
	Foods       []Food  `json:"foods" gorm:"many2many:food_recipes;"`   // 关联的食物
}

// RecipeIngredient 用于JSON序列化和反序列化的结构体
type RecipeIngredient struct {
	FoodName string  `json:"food_name"`
	Weight   float64 `json:"weight"` // 单位：克
}

// TableName 指定表名
func (Recipe) TableName() string {
	return "recipes"
}

// SetIngredients 设置配料表
func (r *Recipe) SetIngredients(ingredients map[string]float64) error {
	jsonData, err := json.Marshal(ingredients)
	if err != nil {
		return err
	}
	r.Ingredients = string(jsonData)
	return nil
}

// GetIngredients 获取配料表
func (r *Recipe) GetIngredients() (map[string]float64, error) {
	var ingredients map[string]float64
	if err := json.Unmarshal([]byte(r.Ingredients), &ingredients); err != nil {
		return nil, err
	}
	return ingredients, nil
}

// CreateRecipe 创建食谱记录
func (r *Recipe) CreateRecipe(db *gorm.DB) error {
	return db.Create(r).Error
}

// GetRecipeByID 通过ID获取食谱
func GetRecipeByID(db *gorm.DB, id uint) (*Recipe, error) {
	var recipe Recipe
	err := db.Preload("Foods").First(&recipe, id).Error
	return &recipe, err
}

// GetAllRecipes 获取所有食谱
func GetAllRecipes(db *gorm.DB) ([]Recipe, error) {
	var recipes []Recipe
	err := db.Preload("Foods").Find(&recipes).Error
	return recipes, err
}

// UpdateRecipe 更新食谱信息
func (r *Recipe) UpdateRecipe(db *gorm.DB) error {
	return db.Save(r).Error
}

// DeleteRecipe 删除食谱
func (r *Recipe) DeleteRecipe(db *gorm.DB) error {
	return db.Delete(r).Error
}

// GetRecipeIDsByFoodID 获取包含指定食材的所有食谱ID
func GetRecipeIDsByFoodID(db *gorm.DB, foodID uint) ([]uint, error) {
    var recipeIDs []uint
    
    // 使用food_recipes关联表查询
    err := db.Table("food_recipes").
        Select("recipe_id").
        Where("food_id = ?", foodID).
        Pluck("recipe_id", &recipeIDs).
        Error
    
    if err != nil {
        return nil, fmt.Errorf("查询食谱ID失败: %v", err)
    }
    
    // 添加空结果检查
    if len(recipeIDs) == 0 {
        return nil, fmt.Errorf("未找到包含食材ID %d 的食谱", foodID)
    }
    
    return recipeIDs, nil
}