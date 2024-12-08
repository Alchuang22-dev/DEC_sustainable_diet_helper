package models

import (
    "time"
    "gorm.io/gorm"
)

// MealType 餐次类型
type MealType string

const (
    Breakfast MealType = "breakfast"
    Lunch     MealType = "lunch"
    Dinner    MealType = "dinner"
    Other     MealType = "other"
)

// NutritionGoal 营养目标模型
type NutritionGoal struct {
    gorm.Model
    UserID        uint      `json:"user_id" gorm:"not null"`
    User          User      `json:"user" gorm:"foreignKey:UserID"`
    Date          time.Time `json:"date" gorm:"not null;index"` // 添加索引以优化日期查询
    Calories      float64   `json:"calories"`
    Protein       float64   `json:"protein"`
    Fat           float64   `json:"fat"`
    Carbohydrates float64   `json:"carbohydrates"`
    Sodium        float64   `json:"sodium"`
}

// NutritionIntake 营养实际摄入模型
type NutritionIntake struct {
    gorm.Model
    UserID        uint      `json:"user_id" gorm:"not null"`
    User          User      `json:"user" gorm:"foreignKey:UserID"`
    Date          time.Time `json:"date" gorm:"not null;index"` // 添加索引以优化日期查询
    MealType      MealType  `json:"meal_type" gorm:"not null"`
    Calories      float64   `json:"calories"`
    Protein       float64   `json:"protein"`
    Fat           float64   `json:"fat"`
    Carbohydrates float64   `json:"carbohydrates"`
    Sodium        float64   `json:"sodium"`
}

// TableName 指定营养目标表名
func (NutritionGoal) TableName() string {
    return "nutrition_goals"
}

// TableName 指定营养摄入表名
func (NutritionIntake) TableName() string {
    return "nutrition_intakes"
}