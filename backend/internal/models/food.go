package models

// 食物结构体
type Food struct {
    ID             int64     `json:"id"`             // 食物ID
    Name_zh           string    `json:"zh_name"`           // 食物中文名称
    Name_en           string    `json:"en_name"`           // 食物英文名称
    Price          float64   `json:"price"`          // 用户购买的价格
    Emissions      float64   `json:"emissions"`      // 食物的碳排放
    Nutrition      Nutrition `json:"nutrition"`      // 食物的营养成分
}

// 营养结构体（根据实际数据库字段重新定义）
type Nutrition struct {
    Calories      float64 `json:"calories"`       // 卡路里
    Protein       float64 `json:"protein"`        // 蛋白质
    Fat           float64 `json:"fat"`            // 脂肪
    Carbohydrates float64 `json:"carbohydrates"`  // 碳水化合物
    Sodium float64 `json:"sodium"`  // 碳水化合物
}

