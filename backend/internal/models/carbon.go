package models

import (
    "time"
    "gorm.io/gorm"
)

// CarbonGoal 碳排放目标模型
type CarbonGoal struct {
    gorm.Model
    UserID   uint      `json:"user_id" gorm:"not null"`
    User     User      `json:"user" gorm:"foreignKey:UserID"`
    Date     time.Time `json:"date" gorm:"not null;index"` // 添加索引以优化日期查询
    Emission float64   `json:"emission"`                   // 碳排放目标值
}

// CarbonIntake 碳排放实际值模型
type CarbonIntake struct {
    gorm.Model
    UserID   uint      `json:"user_id" gorm:"not null"`
    User     User      `json:"user" gorm:"foreignKey:UserID"`
    Date     time.Time `json:"date" gorm:"not null;index"` // 添加索引以优化日期查询
    MealType MealType  `json:"meal_type" gorm:"not null"`
    Emission float64   `json:"emission"`                   // 实际碳排放值
}

// TableName 指定碳排放目标表名
func (CarbonGoal) TableName() string {
    return "carbon_goals"
}

// TableName 指定碳排放实际值表名
func (CarbonIntake) TableName() string {
    return "carbon_intakes"
}