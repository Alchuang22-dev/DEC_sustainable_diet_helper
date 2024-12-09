// internal/models/food_preference.go
package models

import "gorm.io/gorm"

type FoodPreference struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    UserID uint   `json:"user_id"`
    Name   string `gorm:"size:50;not null" json:"name"` // 如 "highProtein"
}

type DislikedFoodPreference struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    UserID uint   `json:"user_id"`
    FoodID uint   `json:"food_id"`
}

// 获取用户设置界面食物偏好类型
func GetUserFoodPreferences(db *gorm.DB, userID uint) ([]FoodPreference, error) {
    var preferences []FoodPreference
    if err := db.Where("user_id = ?", userID).Find(&preferences).Error; err != nil {
        return nil, err
    }
    return preferences, nil
}