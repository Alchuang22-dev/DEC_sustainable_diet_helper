// main.go
package main

import (
    "os"
    "log"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gin-contrib/cors"
    "time"

    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/routes"
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
        &models.News{},
        &models.Resource{},
        &models.Paragraph{},
        &models.Video{},
        &models.Comment{},
        &models.Food{},
        &models.Recipe{},
        &models.Family{},
    )
    if err != nil {
        log.Fatal("自动迁移失败:", err)
    }

    // 初始化Gin引擎
    router := gin.Default()
    router.MaxMultipartMemory = 8 << 20 

    // 配置静态文件服务
    BaseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
    if BaseUploadPath == "" {
        BaseUploadPath = "./upload" // 默认路径
    }
    router.Static("/static", BaseUploadPath)

    // 配置CORS
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // 允许的前端域名
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},                            // 允许的HTTP方法
        AllowHeaders:     []string{"*"},      // 允许的HTTP头部
        AllowCredentials: false,  // 允许跨域请求发送Cookie等凭证 当 AllowCredentials 设置为 true 时，你也需要确保 AllowOrigins 不是 "*". 如果你的前端需要发送 cookies 或其他凭证（如授权头），必须显式列出允许的域名。
        MaxAge:           12 * time.Hour,  // 最大缓存时长
    }))

    // 注册用户路由
    routes.RegisterUserRoutes(router, db)

    // 注册新闻路由
    routes.RegisterNewsRoutes(router, db)

    // 注册食物路由
    routes.RegisterFoodRoutes(router, db)

    // 注册家庭路由
    routes.RegisterFamilyRoutes(router, db)

    // 注册食材偏好路由
    routes.RegisterFoodPreferenceRoutes(router, db)

    // 启动服务器
    err = router.Run(":8080")
    if err != nil {
        log.Fatal("无法启动服务器:", err)
    }
}