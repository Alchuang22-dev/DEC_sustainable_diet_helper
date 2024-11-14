package models

// 菜谱结构体
type Recipe struct {
    ID          int64        `json:"id"`           // 菜谱ID
    Name        string       `json:"name"`         // 菜谱名称
    Method      string       `json:"method"`       // 制作方法
    VideoURL    string       `json:"video_url"`    // 关联的制作视频链接
    Ingredients []Ingredient `json:"ingredients"`  // 相关联的食物类型（食材）
}

// 食材结构体
type Ingredient struct {
    FoodID int64   `json:"food_id"` // 食物ID
    Amount float64 `json:"amount"`  // 用量
    Unit   string  `json:"unit"`    // 单位（如克、毫升）
}