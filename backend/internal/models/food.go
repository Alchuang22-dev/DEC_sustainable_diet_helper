package models

// 食物结构体
type Food struct {
    ID             int64     `json:"id"`             // 食物ID
    Name           string    `json:"name"`           // 食物名称（用户输入类型，下拉选取）
    Weight         float64   `json:"weight"`         // 用户购买的重量
    Price          float64   `json:"price"`          // 用户购买的价格
    TransportMode  string    `json:"transport_mode"` // 用户购买的方式（陆运，海运，空运）
    Location       string    `json:"location"`       // 用户购买的地理位置（本地，外地）
    Emissions      float64   `json:"emissions"`      // 食物的碳排放
    Nutrition      Nutrition `json:"nutrition"`      // 食物的营养成分
}

// 营养结构体（根据实际数据库字段重新定义）
type Nutrition struct {
    Calories      float64 `json:"calories"`       // 卡路里
    Protein       float64 `json:"protein"`        // 蛋白质
    Fat           float64 `json:"fat"`            // 脂肪
    Carbohydrates float64 `json:"carbohydrates"`  // 碳水化合物
    // 可根据需要添加更多营养成分字段
}

// 用户食物集合结构体
type UserFoodCollection struct {
    Foods          []Food    `json:"foods"`           // 关联的食物类型
    TotalEmissions float64   `json:"total_emissions"` // 总碳排放
    TotalNutrition Nutrition `json:"total_nutrition"` // 总营养成分
    Recipes        []Recipe  `json:"recipes"`         // 关联的菜谱
}

