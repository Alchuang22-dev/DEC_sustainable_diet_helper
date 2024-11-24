package main

import (
    "log"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/tools" 
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
)

func main() {
	// 加载配置
	cfg := config.GetConfig()

    err := tools.ImportData("/Users/sunyuanxu/Downloads/recipe_data_with_Eng_name.xlsx", cfg)
    if err != nil {
        log.Fatalf("数据导入失败: %v", err)
    }

    log.Println("数据导入成功！")
}