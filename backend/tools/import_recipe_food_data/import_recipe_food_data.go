// cmd/import_recipe_data.go
package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"  
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "strings"
)

// 将字符串转换为float64的辅助函数
func stringToFloat64(s string) (float64, error) {
    if s == "" {
        return 0, nil
    }
    return strconv.ParseFloat(s, 64)
}

// importFoodsData 导入食物数据
func importFoodsData(db *gorm.DB, filename string) error {
    // 打开CSV文件
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    // 创建CSV reader
    reader := csv.NewReader(file)
    
    // 跳过表头
    if _, err := reader.Read(); err != nil {
        return fmt.Errorf("error reading headers: %v", err)
    }
    
    // 记录成功和失败的数量
    var successful, failed int

    // 逐行读取数据
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error reading row: %v", err)
            failed++
            continue
        }

        // 转换数值字段
        ghg, err := stringToFloat64(record[1])
        calories, err1 := stringToFloat64(record[3])
        protein, err2 := stringToFloat64(record[4])
        fat, err3 := stringToFloat64(record[5])
        carbs, err4 := stringToFloat64(record[6])
        sodium, err5 := stringToFloat64(record[7])
        price, err6 := stringToFloat64(record[8])
        

        if err != nil || err1 != nil || err2 != nil || err3 != nil || 
           err4 != nil || err5 != nil || err6 != nil {
            log.Printf("Error converting number fields for food %s: %v", record[0], err)
            failed++
            continue
        }

        // 创建Food记录
        food := models.Food{
            ZhFoodName:    record[0],
            GHG:           ghg,
            EnFoodName:    strings.ToLower(record[2]),  // 将英文名称转换为小写
            Calories:      calories * 10, // 对于营养的部分，我们数据中的单位是每100g，但是我们的模型中的单位是每1kg
            Protein:       protein * 10,
            Fat:           fat * 10,
            Carbohydrates: carbs * 10,
            Sodium:        sodium * 10,
            Price:         price,
            ImageUrl:      record[9],
        }

        // 保存到数据库
        if err := food.CreateFood(db); err != nil {
            log.Printf("Error saving food %s: %v", food.ZhFoodName, err)
            failed++
            continue
        }
        successful++
    }

    log.Printf("Food import completed. Successful: %d, Failed: %d", successful, failed)
    return nil
}


// importRecipesData 导入食谱数据
func importRecipesData(db *gorm.DB, filename string) error {
    // 打开CSV文件
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("打开文件错误: %v", err)
    }
    defer file.Close()

    // 创建CSV reader
    reader := csv.NewReader(file)
    
    // 跳过表头
    if _, err := reader.Read(); err != nil {
        return fmt.Errorf("读取表头错误: %v", err)
    }

    var successful, failed int
    recipeID := 1

    // 逐行读取数据
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("读取行错误: %v", err)
            failed++
            continue
        }

        // 替换单引号为双引号
        jsonStr := strings.ReplaceAll(record[2], "'", "\"")

        // 解析食材JSON字符串
        var ingredients map[string]float64
        if err := json.Unmarshal([]byte(jsonStr), &ingredients); err != nil {
            log.Printf("解析食材JSON错误 %s: %v", record[3], err)
            failed++
            continue
        }

        // 创建新的 map 存储小写的食材名称
        lowercaseIngredients := make(map[string]float64)
        for foodName, amount := range ingredients {
            lowercaseIngredients[strings.ToLower(foodName)] = amount
        }

        // 将新的 map 转换回 JSON 字符串
        newJsonStr, err := json.Marshal(lowercaseIngredients)
        if err != nil {
            log.Printf("转换食材JSON错误 %s: %v", record[3], err)
            failed++
            continue
        }

        var foods []models.Food
        for foodName := range lowercaseIngredients {
            foodID, err := models.FindFoodIDByName(db, foodName)
            if err != nil {
                log.Printf("找不到食材 %s: %v", foodName, err)
                failed++
                continue
            }
            foods = append(foods, models.Food{Model: gorm.Model{ID: foodID}})
        }

        // 创建Recipe记录
        recipe := models.Recipe{
            URL:         record[0],
            Name:        record[3],
            ImageURL:    fmt.Sprintf("recipes_id_%d", recipeID),
            Ingredients: string(newJsonStr), // 使用转换后的小写JSON字符串
            Foods:       foods,
            Category:    record[5],
        }

        // 保存到数据库
        if err := recipe.CreateRecipe(db); err != nil {
            log.Printf("保存食谱错误 %s: %v", recipe.Name, err)
            failed++
            continue
        }

        successful++
        recipeID++
    }

    log.Printf("食谱导入完成. 成功: %d, 失败: %d", successful, failed)
    return nil
}

func main() {
    // 加载配置
    cfg := config.GetConfig()

    // 构建DSN (Data Source Name)
    dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

    // 连接数据库
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("无法连接到数据库:", err)
    }
    log.Println("数据库连接成功！")

    // 自动迁移数据库表结构
    log.Println("开始迁移数据库表结构...")
    err = db.AutoMigrate(&models.Food{}, &models.Recipe{})
    if err != nil {
        log.Fatal("数据库迁移失败:", err)
    }
    log.Println("数据库表结构迁移完成！")

    // 导入食物数据
    log.Println("Starting food data import...")

    if err := importFoodsData(db, "/app/data/food_dataset/foods_dataset_url.csv"); err != nil {
        log.Fatal("Error importing foods data:", err)
    }

    // 导入菜谱数据
    log.Println("Starting recipe data import...")

    if err := importRecipesData(db, "/app/data/recipes_dataset/recipes_dataset_url.csv"); err != nil {
        log.Fatal("Error importing recipes data:", err)
    }

    log.Println("Data import completed successfully")
    
    // 验证导入结果
    verifyImport(db)
}

// 用于测试导入的函数
func verifyImport(db *gorm.DB) {
    var foodCount, recipeCount int64
    
    db.Model(&models.Food{}).Count(&foodCount)
    db.Model(&models.Recipe{}).Count(&recipeCount)
    
    log.Printf("Verification results:\n")
    log.Printf("Total foods in database: %d\n", foodCount)
    log.Printf("Total recipes in database: %d\n", recipeCount)

    // 验证一些关联
    var recipe models.Recipe
    db.Preload("Foods").First(&recipe)
    log.Printf("Sample recipe '%s' has %d associated foods and %d food names\n", 
        recipe.Name, len(recipe.Foods), len(recipe.Ingredients))
}
