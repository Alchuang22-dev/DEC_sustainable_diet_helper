// internal/models/food_preference.go
package models

type FoodPreference struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    UserID uint   `json:"user_id"`
    Name   string `gorm:"size:50;not null" json:"name"` // 如 "highProtein"
}