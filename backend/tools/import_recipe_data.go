package tools

import (
    "database/sql"
    "fmt"
    "github.com/xuri/excelize/v2"
    _ "github.com/go-sql-driver/mysql"  // 改为 MySQL 驱动
    "strconv"
    "strings"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
    "encoding/json"
)

// ImportData 从Excel导入数据到数据库
func ImportData(excelPath string, dbConfig config.Config) error {
    // 1. 打开Excel文件
    f, err := excelize.OpenFile(excelPath)
    if err != nil {
        return fmt.Errorf("打开Excel文件失败: %v", err)
    }
    defer f.Close()

    // 修改连接字符串格式，直接使用 DBPort 字符串
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
        dbConfig.DBUser, 
        dbConfig.DBPassword, 
        dbConfig.DBHost, 
        dbConfig.DBPort,  // 直接使用 string 类型的端口
        dbConfig.DBName)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return fmt.Errorf("连接数据库失败: %v", err)
    }
    defer db.Close()

    // 3. 创建表（如果不存在）
    err = createTables(db)
    if err != nil {
        return fmt.Errorf("创建表失败: %v", err)
    }

    // 4. 读取Excel数据
    sheetName := f.GetSheetName(0)
    rows, err := f.GetRows(sheetName)
    if err != nil {
        return fmt.Errorf("读取Excel行数据失败: %v", err)
    }

    // 5. 开始事务
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("开始事务失败: %v", err)
    }

    // 6. 插入数据
    for i, row := range rows {
        if i == 0 || i == 1 { // 跳过表头
            continue
        }

        // 解析食谱数据
        recipe, err := parseRecipeFromRow(tx, row)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("解析第%d行数据失败: %v", i+1, err)
        }

        // 插入食谱
        err = insertRecipe(tx, recipe)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("插入第%d行数据失败: %v", i+1, err)
        }
    }

    // 7. 提交事务
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("提交事务失败: %v", err)
    }

    return nil
}

// // createTables 创建必要的数据表
// func createTables(db *sql.DB) error {
//     // 创建recipes表
//     _, err := db.Exec(`
//         CREATE TABLE IF NOT EXISTS recipes (
//             recipe_id SERIAL PRIMARY KEY,
//             recipe_name VARCHAR(200) NOT NULL,
//             img_url TEXT,
//             recipe_description TEXT,
//             cooking_time VARCHAR(50),
//             instructions TEXT,
//             search_word INTEGER[],
//             ghg_production FLOAT,
//             price FLOAT,
//             disposal_amount FLOAT,
//             ghg_disposal FLOAT,
//             ghg_cooking FLOAT,
//             ghg_total FLOAT,
//             dish TEXT,
//             cooking_method TEXT,
//             energy_g FLOAT,
//             fat_g FLOAT,
//             carbohydrates_g FLOAT,
//             zinc_mg FLOAT,
//             folic_acid_µg FLOAT,
//             protein_g FLOAT,
//             total_fiber_g FLOAT,
//             vitamin_a_µg FLOAT,
//             vitamin_c_mg FLOAT,
//             vitamin_e_mg FLOAT,
//             calcium_mg FLOAT,
//             iron_mg FLOAT,
//             potassium_mg FLOAT,
//             magnesium_mg FLOAT,
//             saturated_fat_g FLOAT,
//             cholesterol_g FLOAT,
//             salt_equivalent_g FLOAT
//         )
//     `)
//     if err != nil {
//         return err
//     }

//     /// 2. 创建foods表 - 用于存储食材信息
//     _, err = db.Exec(`
//         CREATE TABLE IF NOT EXISTS foods (
//             food_id SERIAL PRIMARY KEY,
//             food_name VARCHAR(100) NOT NULL UNIQUE
//         )
//     `)
//     if err != nil {
//         return err
//     }

//     // 3. 创建ingredients表 - 关联recipes和foods
//     _, err = db.Exec(`
//         CREATE TABLE IF NOT EXISTS ingredients (
//             recipe_id VARCHAR(50) REFERENCES recipes(recipe_id),
//             food_id INTEGER REFERENCES foods(food_id),
//             amount FLOAT,
//             PRIMARY KEY (recipe_id, food_id)
//         )
//     `)
//     return err
// }
func createTables(db *sql.DB) error {
    // 创建recipes表
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS recipes (
            recipe_id INT AUTO_INCREMENT PRIMARY KEY,
            recipe_name VARCHAR(200) NOT NULL,
            img_url TEXT,
            recipe_description TEXT,
            cooking_time VARCHAR(50),
            instructions TEXT,
            search_word JSON,  -- MySQL 中使用 JSON 类型存储数组
            ghg_production FLOAT,
            price FLOAT,
            disposal_amount JSON,
            ghg_disposal FLOAT,
            ghg_cooking FLOAT,
            ghg_total FLOAT,
            dish TEXT,
            cooking_method TEXT,
            energy_g FLOAT,
            fat_g FLOAT,
            carbohydrates_g FLOAT,
            zinc_mg FLOAT,
            folic_acid_µg FLOAT,
            protein_g FLOAT,
            total_fiber_g FLOAT,
            vitamin_a_µg FLOAT,
            vitamin_c_mg FLOAT,
            vitamin_e_mg FLOAT,
            calcium_mg FLOAT,
            iron_mg FLOAT,
            potassium_mg FLOAT,
            magnesium_mg FLOAT,
            saturated_fat_g FLOAT,
            cholesterol_g FLOAT,
            salt_equivalent_g FLOAT
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `)
    if err != nil {
        return err
    }

    // 创建foods表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS foods (
            food_id INT AUTO_INCREMENT PRIMARY KEY,
            food_name VARCHAR(100) NOT NULL UNIQUE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `)
    if err != nil {
        return err
    }

    // 创建recipe_ingredients关联表
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS recipe_ingredients (
            recipe_id INT,
            food_id INT,
            amount FLOAT,
            PRIMARY KEY (recipe_id, food_id),  -- 复合主键，防止重复关联
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id) ON DELETE CASCADE,  -- 删除食谱时自动删除关联的食材
            FOREIGN KEY (food_id) REFERENCES foods(food_id) ON DELETE RESTRICT  -- 防止删除正在使用的食材
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `)
    if err != nil {
        return err
    }
    
    return nil  // 修正：返回 nil 而不是 err
}
func getOrCreateFoodID(tx *sql.Tx, foodName string) (int64, error) {
    var foodID int64
    
    // 先尝试查找已存在的食材
    err := tx.QueryRow(`
        SELECT food_id FROM foods WHERE food_name = ?  -- 使用 ? 而不是 $1
    `, foodName).Scan(&foodID)
    
    if err == sql.ErrNoRows {
        // 食材不存在，创建新的食材记录
        result, err := tx.Exec(`
            INSERT INTO foods (food_name) 
            VALUES (?)  -- 使用 ? 而不是 $1
        `, foodName)
        
        if err != nil {
            return 0, fmt.Errorf("创建食材记录失败: %v", err)
        }

        // 获取新插入记录的 ID
        foodID, err = result.LastInsertId()
        if err != nil {
            return 0, fmt.Errorf("获取食材ID失败: %v", err)
        }
    } else if err != nil {
        return 0, fmt.Errorf("查询食材失败: %v", err)
    }
    
    return foodID, nil
}

func parseRecipeFromRow(tx *sql.Tx,row []string) (*models.Recipe, error) {
    recipe := &models.Recipe{}
    var err error
    
    // 基本字段解析
    recipe.RecipeName = row[1]
    recipe.ImgURL = row[2]
    recipe.RecipeDescription = row[7]
    recipe.CookingTime = row[9]
    recipe.Instructions = row[10]
 
    // 解析搜索词数组和食材名称
    searchWordStr := row[15] // 食材名称
    searchWordStr = strings.Trim(searchWordStr, "[]")
    searchWordStr = strings.Trim(searchWordStr, "'")  // 移除单引号
    searchWordStrs := strings.Split(searchWordStr, "', '")  // 处理带单引号的分隔

    recipe.SearchWord = make([]int64, 0)
    for _, str := range searchWordStrs {
        str = strings.TrimSpace(str)
        if str == "" {
            continue
        }
        // 获取或创建食材ID
        foodID, err := getOrCreateFoodID(tx, str)
        if err != nil {
            return nil, fmt.Errorf("处理食材'%s'失败: %v", str, err) 
        }
        recipe.SearchWord = append(recipe.SearchWord, foodID)
    }
 
    dataStr := row[16]

    // 首先定义一个辅助的 map 结构来解析原始数据
    var rawData map[string]float64
    // 将单引号替换为双引号
    dataStr = strings.ReplaceAll(dataStr, "'", "\"")
    err = json.Unmarshal([]byte(dataStr), &rawData)
    if err != nil {
        // 如果解析失败，打印出问题的字符串以便调试
        return nil, fmt.Errorf("解析原始数据失败: %v\n数据内容: %s", err, dataStr)
    }

    // 准备 SQL 查询语句
    query := "SELECT food_id FROM foods WHERE food_name = ?"

    // 转换为 Ingredient 切片
    var modifiedIngredients []models.Ingredient
    for foodName, amount := range rawData {
        // 查询数据库获取 food_id
        var foodID int64
        err := tx.QueryRow(query, foodName).Scan(&foodID)
        if err != nil {
            if err == sql.ErrNoRows {
                // 如果在数据库中找不到对应的食材，设置 food_id 为 -1
                foodID, err = getOrCreateFoodID(tx, foodName)
                if err != nil{
                    return nil, fmt.Errorf("查询食材ID失败: %v", err)
                }
            } else {
                return nil, fmt.Errorf("查询食材ID失败: %v", err)
            }
        }

        // 创建 Ingredient 对象并添加到切片中
        ingredient := models.Ingredient{
            FoodName: foodName,
            FoodID:   int(foodID),  // 将 int64 转换为 int
            Amount:   amount,
        }
        modifiedIngredients = append(modifiedIngredients, ingredient)
    }

    recipe.ModifiedIngredients = modifiedIngredients
    
 
    // 解析浮点数字段
    if recipe.GHGProduction, err = strconv.ParseFloat(row[18], 64); err != nil {
        return nil, fmt.Errorf("解析GHG_production失败: %v", err)
    }
    if recipe.Price, err = strconv.ParseFloat(row[19], 64); err != nil {
        return nil, fmt.Errorf("解析price失败: %v", err)
    }
    // if recipe.DisposalAmount, err = strconv.ParseFloat(row[20], 64); err != nil {
    //     return nil, fmt.Errorf("解析disposal_amount失败: %v", err)
    // }
    dataStr = row[20]
    // 将单引号替换为双引号
    dataStr = strings.ReplaceAll(dataStr, "'", "\"")
    var disposalAmount models.DisposalInfo
    err = json.Unmarshal([]byte(dataStr), &disposalAmount)
    if err != nil {
        return nil, fmt.Errorf("解析disposal_amount失败: %v", err)
    }
    recipe.DisposalAmount = disposalAmount


    if recipe.GHGDisposal, err = strconv.ParseFloat(row[21], 64); err != nil {
        return nil, fmt.Errorf("解析GHG_disposal失败: %v", err)
    }
    if recipe.GHGCooking, err = strconv.ParseFloat(row[22], 64); err != nil {
        return nil, fmt.Errorf("解析GHG_cooking失败: %v", err)
    }
    if recipe.GHGTotals, err = strconv.ParseFloat(row[23], 64); err != nil {
        return nil, fmt.Errorf("解析GHG_total失败: %v", err)
    }
 
    recipe.Dish = row[24]
    recipe.CookingMethod = row[27]
 
    if recipe.Energy, err = strconv.ParseFloat(row[28], 64); err != nil {
        return nil, fmt.Errorf("解析energy_g失败: %v", err)
    }
    if recipe.Fat, err = strconv.ParseFloat(row[29], 64); err != nil {
        return nil, fmt.Errorf("解析fat_g失败: %v", err)
    }
    if recipe.Carbohydrates, err = strconv.ParseFloat(row[30], 64); err != nil {
        return nil, fmt.Errorf("解析carbohydrates_g失败: %v", err)
    }
    if recipe.Zinc, err = strconv.ParseFloat(row[31], 64); err != nil {
        return nil, fmt.Errorf("解析zinc_mg失败: %v", err)
    }
    if recipe.FolicAcid, err = strconv.ParseFloat(row[32], 64); err != nil {
        return nil, fmt.Errorf("解析folic_acid_µg失败: %v", err)
    }
    if recipe.Protein, err = strconv.ParseFloat(row[33], 64); err != nil {
        return nil, fmt.Errorf("解析protein_g失败: %v", err)
    }
    if recipe.TotalFiber, err = strconv.ParseFloat(row[34], 64); err != nil {
        return nil, fmt.Errorf("解析total_fiber_g失败: %v", err)
    }
    if recipe.VitaminA, err = strconv.ParseFloat(row[35], 64); err != nil {
        return nil, fmt.Errorf("解析vitamin_a_µg失败: %v", err)
    }
    if recipe.VitaminC, err = strconv.ParseFloat(row[36], 64); err != nil {
        return nil, fmt.Errorf("解析vitamin_c_mg失败: %v", err)
    }
    if recipe.VitaminE, err = strconv.ParseFloat(row[37], 64); err != nil {
        return nil, fmt.Errorf("解析vitamin_e_mg失败: %v", err)
    }
    if recipe.Calcium, err = strconv.ParseFloat(row[38], 64); err != nil {
        return nil, fmt.Errorf("解析calcium_mg失败: %v", err)
    }
    if recipe.Iron, err = strconv.ParseFloat(row[39], 64); err != nil {
        return nil, fmt.Errorf("解析iron_mg失败: %v", err)
    }
    if recipe.Potassium, err = strconv.ParseFloat(row[40], 64); err != nil {
        return nil, fmt.Errorf("解析potassium_mg失败: %v", err)
    }
    if recipe.Magnesium, err = strconv.ParseFloat(row[41], 64); err != nil {
        return nil, fmt.Errorf("解析magnesium_mg失败: %v", err)
    }
    if recipe.SaturatedFat, err = strconv.ParseFloat(row[42], 64); err != nil {
        return nil, fmt.Errorf("解析saturated_fat_g失败: %v", err)
    }
    if recipe.Cholesterol, err = strconv.ParseFloat(row[43], 64); err != nil {
        return nil, fmt.Errorf("解析cholesterol_g失败: %v", err)
    }
    if recipe.SaltEquivalent, err = strconv.ParseFloat(row[44], 64); err != nil {
        return nil, fmt.Errorf("解析salt_equivalent_g失败: %v", err)
    }
 
    return recipe, nil
 }


 func insertRecipe(tx *sql.Tx, recipe *models.Recipe) error {
    // 将 SearchWord 数组转换为 JSON 字符串
    searchWordJSON, err := json.Marshal(recipe.SearchWord)
    if err != nil {
        return fmt.Errorf("转换 search_word 失败: %v", err)
    }

    // 将 DisposalAmount 结构体转换为 JSON 字符串
    disposalAmountJSON, err := json.Marshal(recipe.DisposalAmount)
    if err != nil {
        return fmt.Errorf("转换 disposal_amount 失败: %v", err)
    }

    result, err := tx.Exec(`
        INSERT INTO recipes (
            recipe_name, img_url, recipe_description, cooking_time,
            instructions, search_word, ghg_production, price, disposal_amount, ghg_disposal,
            ghg_cooking, ghg_total, dish, cooking_method, energy_g, fat_g,
            carbohydrates_g, zinc_mg, folic_acid_µg, protein_g, total_fiber_g,
            vitamin_a_µg, vitamin_c_mg, vitamin_e_mg, calcium_mg, iron_mg,
            potassium_mg, magnesium_mg, saturated_fat_g, cholesterol_g,
            salt_equivalent_g
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
                    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        recipe.RecipeName, recipe.ImgURL, recipe.RecipeDescription,
        recipe.CookingTime, recipe.Instructions, searchWordJSON,
        recipe.GHGProduction, recipe.Price, disposalAmountJSON,  // 使用 JSON 字符串而不是结构体
        recipe.GHGDisposal, recipe.GHGCooking, recipe.GHGTotals,
        recipe.Dish, recipe.CookingMethod, recipe.Energy, recipe.Fat,
        recipe.Carbohydrates, recipe.Zinc, recipe.FolicAcid, recipe.Protein,
        recipe.TotalFiber, recipe.VitaminA, recipe.VitaminC, recipe.VitaminE,
        recipe.Calcium, recipe.Iron, recipe.Potassium, recipe.Magnesium,
        recipe.SaturatedFat, recipe.Cholesterol, recipe.SaltEquivalent,
    )
    if err != nil {
        return fmt.Errorf("插入食谱失败: %v", err)
    }

    // 获取插入的 ID
    recipeID, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("获取插入ID失败: %v", err)
    }

    // 插入食材信息
    for _, ingredient := range recipe.ModifiedIngredients {    
        _, err = tx.Exec(`
            INSERT INTO recipe_ingredients (recipe_id, food_id, amount)
            VALUES (?, ?, ?)`,
            recipeID, ingredient.FoodID, ingredient.Amount,
        )
        if err != nil {
            return fmt.Errorf("插入食材关联失败: %v", err)
        }
    }

    return nil
}