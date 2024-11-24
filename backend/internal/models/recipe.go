package models


// Recipe 代表一个食谱的结构体
type Recipe struct {
    RecipeID           int       `json:"recipe_id" db:"recipe_id"`
    RecipeName         string    `json:"recipe_name" db:"recipe_name"`
    ImgURL             string    `json:"img_url" db:"img_url"`
    RecipeDescription  string    `json:"recipe_description" db:"recipe_description"`
    CookingTime        string    `json:"cooking_time" db:"cooking_time"`
    Instructions       string    `json:"instructions" db:"instructions"`
    SearchWord         []int64    `json:"SearchWorld" db:"search_word"`
    ModifiedIngredients []Ingredient   `json:"modified_ingredients" db:"modified_ingredients"`
    GHGProduction      float64   `json:"GHG_production" db:"ghg_production"`
    Price              float64   `json:"price" db:"price"`
    DisposalAmount     DisposalInfo   `json:"disposal_amount" db:"disposal_amount"`
    GHGDisposal        float64   `json:"GHG_disposal" db:"ghg_disposal"`
    GHGCooking         float64   `json:"GHG_cooking" db:"ghg_cooking"`
    GHGTotals          float64   `json:"GHG_total" db:"ghg_total"`
    Dish               string    `json:"dish" db:"dish"`
    CookingMethod      string    `json:"cooking_method" db:"cooking_method"`
    Energy             float64   `json:"energy_g" db:"energy_g"`
    Fat                float64   `json:"fat_g" db:"fat_g"`
    Carbohydrates      float64   `json:"carbohydrates_g" db:"carbohydrates_g"`
    Zinc               float64   `json:"zinc_mg" db:"zinc_mg"`
    FolicAcid          float64   `json:"folic_acid_µg" db:"folic_acid_µg"`
    Protein            float64   `json:"protein_g" db:"protein_g"`
    TotalFiber         float64   `json:"total_fiber_g" db:"total_fiber_g"`
    VitaminA           float64   `json:"vitamin_a_µg" db:"vitamin_a_µg"`
    VitaminC           float64   `json:"vitamin_c_mg" db:"vitamin_c_mg"`
    VitaminE           float64   `json:"vitamin_e_mg" db:"vitamin_e_mg"`
    Calcium            float64   `json:"calcium_mg" db:"calcium_mg"`
    Iron               float64   `json:"iron_mg" db:"iron_mg"`
    Potassium          float64   `json:"potassium_mg" db:"potassium_mg"`
    Magnesium          float64   `json:"magnesium_mg" db:"magnesium_mg"`
    SaturatedFat       float64   `json:"saturated_fat_g" db:"saturated_fat_g"`
    Cholesterol        float64   `json:"cholesterol_g" db:"cholesterol_g"`
    SaltEquivalent     float64   `json:"salt_equivalent_g" db:"salt_equivalent_g"`
}

// 食材结构体
type Ingredient struct {
    FoodName string   `json:"food_name"` // 食物名称
    FoodID   int      `json:"food_id"`   // 食物ID
    Amount float64 `json:"amount"`  // 用量
}

type DisposalInfo struct {
    Refuse         float64 `json:"refuse"`
    FoodLoss       float64 `json:"food_loss"`
    OverRemoval    float64 `json:"over_removal"`
    DirectDisposal float64 `json:"direct_disposal"`
    Leftover       float64 `json:"leftover"`
}