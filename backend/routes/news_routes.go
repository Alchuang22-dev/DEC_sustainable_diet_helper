// routes/news_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/shenlayu/se-backend/controllers"
)

func RegisterNewsRoutes(router *gin.Engine, db *gorm.DB) {
    newsController := controllers.NewNewsController(db)
    newsGroup := router.Group("/news")
    {
        // 创建视频新闻
        newsGroup.POST("/video", newsController.CreateVideoNews)
        // 创建常规新闻
        newsGroup.POST("/regular", newsController.CreateRegularNews)
        // 获取新闻详情
        newsGroup.GET("/:id", newsController.GetNewsDetail)
        // 添加评论
        newsGroup.POST("/:id/comment", newsController.AddComment)

		// TODO
		// 删除评论
		// 删除新闻
		// 更新新闻
    }
}