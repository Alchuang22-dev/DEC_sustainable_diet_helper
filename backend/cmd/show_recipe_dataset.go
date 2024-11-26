package main

import (
    "database/sql"
    "fmt"
    "os"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func ShowFoodsDataset(db *sql.DB) {
    // 1. 创建输出文件
    foodsFile, err := os.Create("foods_dataset.txt")
    if err != nil {
        log.Fatalf("创建foods输出文件失败: %v", err)
    }
    defer foodsFile.Close()

    recipesFile, err := os.Create("recipes_dataset.txt")
    if err != nil {
        log.Fatalf("创建recipes输出文件失败: %v", err)
    }
    defer recipesFile.Close()

    // 2. 查询并输出 foods 表数据
    rows, err := db.Query("SELECT * FROM foods ORDER BY food_id")
    if err != nil {
        log.Fatalf("查询foods表失败: %v", err)
    }
    defer rows.Close()

    // 写入 foods 表头
    fmt.Fprintln(foodsFile, "Food ID\tFood Name")
    fmt.Fprintln(foodsFile, "-----------------")

    // 遍历 foods 结果
    for rows.Next() {
        var foodID int
        var foodName string
        if err := rows.Scan(&foodID, &foodName); err != nil {
            log.Fatalf("读取food行数据失败: %v", err)
        }
        fmt.Fprintf(foodsFile, "%d\t%s\n", foodID, foodName)
    }

    // 3. 查询并输出 recipes 表数据
    rows, err = db.Query(`
        SELECT r.recipe_id, r.recipe_name, r.instructions, r.ghg_total, r.energy_g,
               GROUP_CONCAT(CONCAT(f.food_name, ': ', ri.amount) SEPARATOR ', ') as ingredients
        FROM recipes r
        LEFT JOIN recipe_ingredients ri ON r.recipe_id = ri.recipe_id
        LEFT JOIN foods f ON ri.food_id = f.food_id
        GROUP BY r.recipe_id
        ORDER BY r.recipe_id
    `)
    if err != nil {
        log.Fatalf("查询recipes表失败: %v", err)
    }
    defer rows.Close()

    // 写入 recipes 表头
    fmt.Fprintln(recipesFile, "Recipe ID\tRecipe Name\tIngredients\tInstructions\tTotal GHG\tEnergy (g)")
    fmt.Fprintln(recipesFile, "-----------------------------------------------------------")

    // 遍历 recipes 结果
    for rows.Next() {
        var (
            recipeID                              int
            recipeName, instructions, ingredients string
            ghgTotal, energy                     float64
        )
        if err := rows.Scan(&recipeID, &recipeName, &instructions, &ghgTotal, &energy, &ingredients); err != nil {
            log.Fatalf("读取recipe行数据失败: %v", err)
        }
        
        // 格式化输出到文件
        fmt.Fprintf(recipesFile, "%d\t%s\t%s\t%s\t%.2f\t%.2f\n",
            recipeID, recipeName, ingredients, instructions, ghgTotal, energy)
    }

    // 4. 额外查询一些统计信息并打印到控制台
    var foodCount, recipeCount int
    err = db.QueryRow("SELECT COUNT(*) FROM foods").Scan(&foodCount)
    if err != nil {
        log.Fatalf("统计foods数量失败: %v", err)
    }

    err = db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&recipeCount)
    if err != nil {
        log.Fatalf("统计recipes数量失败: %v", err)
    }

    fmt.Printf("\n数据库统计信息:\n")
    fmt.Printf("总食材数量: %d\n", foodCount)
    fmt.Printf("总食谱数量: %d\n", recipeCount)
}

func main() {
    // 修改连接字符串格式
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
        "root",           // 用户名
        "syx20040412",     // 密码
        "localhost",      // 主机名
        "3306",          // 端口
        "dec_db")        // 数据库名
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("连接数据库失败: %v", err)
    }
    defer db.Close()

    // 查看foods,recipes表，并输出到文件中
    ShowFoodsDataset(db)
    log.Println("食谱数据集生成成功！")
}