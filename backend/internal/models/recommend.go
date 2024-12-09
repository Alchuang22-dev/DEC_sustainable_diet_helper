// internal/models/ingredient_preference.go
package models

import "time"

// UserRecipeHistory 用户菜谱选择历史
type UserRecipeHistory struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    RecipeID  uint      `gorm:"not null"`
    SelectTime time.Time `gorm:"not null"`
}

// UserIngredientHistory 用户食材选择历史
type UserIngredientHistory struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      uint      `gorm:"not null"`
    IngredientID uint     `gorm:"not null"`
    SelectTime  time.Time `gorm:"not null"`
}

// UserIngredientPreference 用户食材偏好
type UserIngredientPreference struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      uint      `gorm:"not null"`
    IngredientID uint     `gorm:"not null"`
    IsLike      bool      `gorm:"not null"` // true表示喜欢，false表示不喜欢
    UpdateTime  time.Time `gorm:"not null"`
}
