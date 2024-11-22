// routes/news_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
)

// TODO 修改，把id去掉，
func RegisterNewsRoutes(router *gin.Engine, db *gorm.DB) {
    newsController := controllers.NewNewsController(db)
    newsGroup := router.Group("/news")
    {
        // 需要认证的路由
        authGroup := newsGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            // 创建新闻
            authGroup.POST("/upload", newsController.UploadNews)
            // 点赞新闻
            authGroup.POST("/:id/like", newsController.LikeNews)
            // 取消点赞新闻
            authGroup.POST("/:id/cancel_like", newsController.CancelLikeNews)
            // 收藏新闻
            authGroup.POST("/:id/favourite", newsController.FavoriteNews)
            // 取消收藏新闻
            authGroup.POST("/:id/cancel_favourite", newsController.CancelFavoriteNews)
            // 点踩新闻
            authGroup.POST("/:id/dislike", newsController.DislikeNews)
            // 取消点踩新闻
            authGroup.POST("/:id/cancel_dislike", newsController.CancelDislikeNews)
            // 添加评论
            authGroup.POST("/:id/comment", newsController.AddComment)
            // 浏览新闻（可选，需要记录谁浏览了）
            authGroup.POST("/:id/view", newsController.ViewNews)
        }

        // 不需要认证的路由
        newsGroup.GET("/:id", newsController.GetNewsDetail)
        // 其他公开的路由
    }
}