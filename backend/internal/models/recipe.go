// internal/models/recipe.go
package models

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "gorm.io/gorm"
)
// Ingredients 用于存储食材成分的map
type Ingredients map[string]float64

// Value 实现driver.Valuer接口
func (i Ingredients) Value() (driver.Value, error) {
    return json.Marshal(i)
}

// Scan 实现sql.Scanner接口
func (i *Ingredients) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return nil
    }
    return json.Unmarshal(bytes, &i)
}

// StringArray 用于存储食物名字数组
type StringArray []string

// Value 实现driver.Valuer接口
func (a StringArray) Value() (driver.Value, error) {
    return json.Marshal(a)
}

// Scan 实现sql.Scanner接口
func (a *StringArray) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return nil
    }
    return json.Unmarshal(bytes, &a)
}

type Recipe struct {
    gorm.Model
    URL         string      `json:"url" gorm:"column:url;not null"`
    RecipeName  string      `json:"recipe_name" gorm:"column:recipe_name;not null"`
    FoodNames   StringArray `json:"food_names" gorm:"type:json"`
    Ingredients Ingredients `json:"ingredients" gorm:"type:json"`
    ImageURL    string      `json:"image_url" gorm:"column:image_url"` // 新增图片URL字段
    Foods       []Food      `json:"foods" gorm:"many2many:food_recipes;"`
}

// TableName 指定表名
func (Recipe) TableName() string {
    return "recipes"
}

// CreateRecipe 创建菜谱
func (r *Recipe) CreateRecipe(db *gorm.DB) error {
    return db.Create(r).Error
}

// GetRecipeByID 通过ID获取菜谱
func GetRecipeByID(db *gorm.DB, id uint) (*Recipe, error) {
    var recipe Recipe
    err := db.First(&recipe, id).Error
    return &recipe, err
}

// GetAllRecipes 获取所有菜谱
func GetAllRecipes(db *gorm.DB) ([]Recipe, error) {
    var recipes []Recipe
    err := db.Find(&recipes).Error
    return recipes, err
}

// GetRecipesByFoodID 获取与特定食物相关的所有菜谱
func GetRecipesByFoodID(db *gorm.DB, foodID uint) ([]Recipe, error) {
    var recipes []Recipe
    err := db.Joins("JOIN food_recipes ON food_recipes.recipe_id = recipes.id").
           Where("food_recipes.food_id = ?", foodID).
           Find(&recipes).Error
    return recipes, err
}

// UpdateRecipe 更新菜谱信息
func (r *Recipe) UpdateRecipe(db *gorm.DB) error {
    return db.Save(r).Error
}

// DeleteRecipe 删除菜谱
func (r *Recipe) DeleteRecipe(db *gorm.DB) error {
    return db.Delete(r).Error
}

// GetRecipesByFoodName 通过食物名称获取相关菜谱
func GetRecipesByFoodName(db *gorm.DB, foodName string) ([]Recipe, error) {
    var recipes []Recipe
    err := db.Where("food_names @> ?", fmt.Sprintf("[\"%s\"]", foodName)).Find(&recipes).Error
    return recipes, err
}

// GetRecipeWithFoods 获取菜谱及其关联的食物信息
func GetRecipeWithFoods(db *gorm.DB, recipeID uint) (*Recipe, error) {
    var recipe Recipe
    err := db.Preload("Foods").First(&recipe, recipeID).Error
    return &recipe, err
}