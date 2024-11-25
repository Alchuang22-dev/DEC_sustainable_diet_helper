package models

import (
    "gorm.io/gorm"
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
