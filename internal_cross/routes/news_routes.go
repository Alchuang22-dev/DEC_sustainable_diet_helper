// routes/news_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/YourUsername/YourProject/internal/controllers"
)

func RegisterNewsRoutes(router *gin.Engine, db *gorm.DB) {
    newsController := controllers.NewNewsController(db)
    newsGroup := router.Group("/news")
    {
        // 创建新闻
        newsGroup.POST("/", newsController.CreateNews)
        // 获取所有新闻
        newsGroup.GET("/", newsController.GetAllNews)
        // 获取新闻详情
        newsGroup.GET("/:id", newsController.GetNewsDetail)
        // 更新新闻
        newsGroup.PUT("/:id", newsController.UpdateNews)
        // 删除新闻
        newsGroup.DELETE("/:id", newsController.DeleteNews)

        // 添加评论
        newsGroup.POST("/:id/comment", newsController.AddComment)
        // 删除评论
        newsGroup.DELETE("/:id/comment/:commentId", newsController.DeleteComment)
        // 获取新闻的所有评论
        newsGroup.GET("/:id/comments", newsController.GetComments)

        // 点赞新闻
        newsGroup.POST("/:id/like", newsController.LikeNews)
        // 取消点赞新闻
        newsGroup.POST("/:id/cancel_like", newsController.CancelLikeNews)

        // 收藏新闻
        newsGroup.POST("/:id/favorite", newsController.FavoriteNews)
        // 取消收藏新闻
        newsGroup.POST("/:id/cancel_favorite", newsController.CancelFavoriteNews)

        // 点踩新闻
        newsGroup.POST("/:id/dislike", newsController.DislikeNews)
        // 取消点踩新闻
        newsGroup.POST("/:id/cancel_dislike", newsController.CancelDislikeNews)

        // 浏览新闻
        newsGroup.POST("/:id/view", newsController.ViewNews)

        // 根据标签获取新闻
        newsGroup.GET("/tags/:tag", newsController.GetNewsByTag)
        // 搜索新闻
        newsGroup.GET("/search", newsController.SearchNews)

		// 关注作者
        newsGroup.POST("/:id/follow_author", newsController.FollowAuthor)
        // 取消关注作者
        newsGroup.POST("/:id/unfollow_author", newsController.UnfollowAuthor)
        // 分享新闻
        newsGroup.POST("/:id/share", newsController.ShareNews)
    }
}
