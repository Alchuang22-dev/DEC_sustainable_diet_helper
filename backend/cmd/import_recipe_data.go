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
            EnFoodName:    record[2],
            Calories:      calories,
            Protein:       protein,
            Fat:           fat,
            Carbohydrates: carbs,
            Sodium:        sodium,
            Price:         price,
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

// cmd/import_recipe_data.go

func importRecipesData(db *gorm.DB, filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    content, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("error reading file: %v", err)
    }

    lines := strings.Split(string(content), "\n")
    if len(lines) < 2 {
        return fmt.Errorf("file has insufficient lines")
    }

    var successful, failed int
    for i, line := range lines[1:] {
        if line == "" {
            continue
        }

        fields := parseCSVLine(line)
        if len(fields) < 5 {
            log.Printf("Line %d: Insufficient fields", i+2)
            failed++
            continue
        }

        // 解析食物名字数组
        foodNamesStr := strings.TrimSpace(fields[1])
        var foodNames models.StringArray
        if err := json.Unmarshal([]byte(foodNamesStr), &foodNames); err != nil {
            log.Printf("Line %d: Error parsing food names: %v", i+2, err)
            log.Printf("Food names string: %s", foodNamesStr)
            failed++
            continue
        }

        // 处理 ingredients 字符串
        ingredientsStr := strings.TrimSpace(fields[2])
        // 移除外部的引号（如果存在）
        ingredientsStr = strings.Trim(ingredientsStr, "\"")
        // 将单引号替换为双引号
        ingredientsStr = strings.ReplaceAll(ingredientsStr, "'", "\"")
        
        // 创建一个临时的 map 来存储解析后的数据
        ingredients := make(models.Ingredients)
        
        // 手动解析 JSON 字符串
        ingredientsStr = strings.Trim(ingredientsStr, "{}")
        pairs := strings.Split(ingredientsStr, ",")
        for _, pair := range pairs {
            pair = strings.TrimSpace(pair)
            kv := strings.Split(pair, ":")
            if len(kv) != 2 {
                continue
            }
            
            // 处理键
            key := strings.Trim(strings.TrimSpace(kv[0]), "\"")
            
            // 处理值
            valueStr := strings.TrimSpace(kv[1])
            value, err := strconv.ParseFloat(valueStr, 64)
            if err != nil {
                log.Printf("Error parsing value for key %s: %v", key, err)
                continue
            }
            
            ingredients[key] = value
        }

        // 创建Recipe记录
        recipe := models.Recipe{
            URL:         strings.TrimSpace(fields[0]),
            RecipeName:  strings.TrimSpace(fields[3]),
            ImageURL:    strings.TrimSpace(fields[4]),
            FoodNames:   foodNames,
            Ingredients: ingredients,
        }

        // 保存到数据库
        if err := recipe.CreateRecipe(db); err != nil {
            log.Printf("Error saving recipe %s: %v", recipe.RecipeName, err)
            failed++
            continue
        }

        // 为每个食物名字建立关联
        for _, foodName := range foodNames {
            var relatedFood models.Food
            result := db.Where("zh_food_name = ? OR en_food_name = ?", 
                foodName, foodName).First(&relatedFood)
                
            if result.Error == nil {
                if err := db.Model(&recipe).Association("Foods").Append(&relatedFood); err != nil {
                    log.Printf("Error associating food %s with recipe %s: %v", 
                        foodName, recipe.RecipeName, err)
                }
            } else {
                log.Printf("Could not find food with name %s for recipe %s", 
                    foodName, recipe.RecipeName)
            }
        }

        successful++
    }

    log.Printf("Recipe import completed. Successful: %d, Failed: %d", successful, failed)
    return nil
}

// 自定义CSV行解析函数 (保持不变)
func parseCSVLine(line string) []string {
    var fields []string
    var currentField strings.Builder
    inQuotes := false
    inBrackets := 0
    inBraces := 0

    for _, ch := range line {
        switch ch {
        case '[':
            inBrackets++
            currentField.WriteRune(ch)
        case ']':
            inBrackets--
            currentField.WriteRune(ch)
        case '{':
            inBraces++
            currentField.WriteRune(ch)
        case '}':
            inBraces--
            currentField.WriteRune(ch)
        case '"':
            inQuotes = !inQuotes
            currentField.WriteRune(ch)
        case ',':
            if inQuotes || inBrackets > 0 || inBraces > 0 {
                currentField.WriteRune(ch)
            } else {
                fields = append(fields, currentField.String())
                currentField.Reset()
            }
        default:
            currentField.WriteRune(ch)
        }
    }

    if currentField.Len() > 0 {
        fields = append(fields, currentField.String())
    }

    return fields
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
    if err := importFoodsData(db, "foods_dataset.csv"); err != nil {
        log.Fatal("Error importing foods data:", err)
    }

    // 导入菜谱数据
    log.Println("Starting recipe data import...")
    if err := importRecipesData(db, "recipes_dataset.csv"); err != nil {
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
        recipe.RecipeName, len(recipe.Foods), len(recipe.FoodNames))
}