// main.go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"

    "DEC/config"
    "DEC/internal/models"
    "DEC/internal/routes"
)

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

    // 自动迁移
    err = db.AutoMigrate(
        &models.User{},
        &models.VideoNews{},
        &models.RegularNews{},
        &models.Resource{},
        &models.Paragraph{},
        &models.Comment{},
    )
    if err != nil {
        log.Fatal("自动迁移失败:", err)
    }

    // 初始化Gin引擎
    router := gin.Default()

    // 注册用户路由
    routes.RegisterUserRoutes(router, db)

    // 注册新闻路由
    routes.RegisterNewsRoutes(router, db)

    // 启动服务器
    err = router.Run(":8080")
    if err != nil {
        log.Fatal("无法启动服务器:", err)
    }
}