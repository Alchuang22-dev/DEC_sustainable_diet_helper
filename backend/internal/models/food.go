// food.go
package models

import (
    "time"
)

type FoodType struct {
    ID          uint          `gorm:"primaryKey" json:"id"`
    Name        string        `gorm:"unique;not null" json:"name"`
    CreatedAt   time.Time     `json:"created_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
    OpenFoods   []OpenFoodData
    UserEntries []UserFoodEntry
}

type OpenFoodData struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    FoodTypeID    uint      `json:"food_type_id"`
    FoodName      string    `gorm:"not null" json:"food_name"`
    CarbonEmission float64   `json:"carbon_emission"` // 单位 kg CO₂e
    OtherAttributes string   `json:"other_attributes"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type UserFoodEntry struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `json:"user_id"`
    FoodTypeID  uint      `json:"food_type_id"`
    TotalWeight float64   `json:"total_weight"` // 单位 kg
    TotalPrice  float64   `json:"total_price"`  // 单位 RMB
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
