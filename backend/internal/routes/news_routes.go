// routes/news_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
)

func RegisterNewsRoutes(router *gin.Engine, db *gorm.DB) {
    newsController := controllers.NewNewsController(db)
    newsGroup := router.Group("/news")
    {
        // 需要认证的路由
        authGroup := newsGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            authGroup.POST("/upload_image", newsController.UploadImage)      // 上传单张图片
            authGroup.POST("/create_draft", newsController.CreateDraft)     // 创建草稿
            // 更新草稿
            authGroup.PUT("/drafts/:id", newsController.UpdateDraft)
            // 删除草稿
            authGroup.DELETE("/drafts/:id", newsController.DeleteDraft)
            // 将草稿转换为新闻
            authGroup.POST("/convert_draft", newsController.ConvertDraftToNews)
            // 获取自己的新闻 ID 列表
            authGroup.GET("/my_news", newsController.GetMyNews)
            // 获取自己的草稿 ID 列表
            authGroup.GET("/my_drafts", newsController.GetMyDrafts)
            // 预览新闻
            authGroup.POST("/preview_news", newsController.PreviewNews)
            // 预览草稿
            authGroup.POST("/preview_drafts", newsController.PreviewDrafts)
            // 详细查看新闻
            authGroup.GET("/details/news/:id", newsController.GetNewsDetails)
            // 详细查看草稿
            authGroup.GET("/details/draft/:id", newsController.GetDraftDetails)
            // 删除新闻
            authGroup.DELETE("/:id", newsController.DeleteNews)
            
            // 获取新闻 ID
            authGroup.GET("/paginated/view_count", newsController.GetNewsByViewCount) // 观看量降序
            authGroup.GET("/paginated/like_count", newsController.GetNewsByLikeCount) // 点赞量降序
            authGroup.GET("/paginated/upload_time", newsController.GetNewsByUploadTime) // 时间由旧到新

            // 评论相关
            authGroup.POST("/comments", newsController.AddComment)      // 添加评论
            authGroup.DELETE("/comments/:id", newsController.DeleteComment) // 删除评论
            // 点赞相关
            authGroup.POST("/:id/like", newsController.LikeNews)            // 点赞新闻
            authGroup.DELETE("/:id/like", newsController.CancelLikeNews)   // 取消点赞新闻
            // 收藏相关
            authGroup.POST("/:id/favorite", newsController.FavoriteNews)          // 收藏新闻
            authGroup.DELETE("/:id/favorite", newsController.CancelFavoriteNews)  // 取消收藏新闻
            // 点踩相关
            authGroup.POST("/:id/dislike", newsController.DislikeNews)            // 点踩新闻
            authGroup.DELETE("/:id/dislike", newsController.CancelDislikeNews)   // 取消点踩新闻
            // 浏览记录
            authGroup.POST("/:id/view", newsController.ViewNews) // 浏览新闻

            authGroup.GET("/:id/status", newsController.GetUserNewsStatus) // 返回用户对新闻的过往交互

            authGroup.POST("/:id/comment_like", newsController.LikeComment) // 点赞评论
            authGroup.DELETE("/:id/comment_like", newsController.CancelLikeComment) // 取消点赞评论

            authGroup.GET("/search", newsController.SearchNews)
        }
    }
}