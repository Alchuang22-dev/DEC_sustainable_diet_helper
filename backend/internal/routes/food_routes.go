// routes/news_routes.go
package routes

// import (
//     "gorm.io/gorm"
//     "github.com/gin-gonic/gin"
//     "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
// )

// func RegisterFoodRoutes(router *gin.Engine, db *gorm.DB) {
//     newsController := controllers.NewNewsController(db)
//     newsGroup := router.Group("/news")
//     {
//         // 创建新闻
//         newsGroup.POST("/", newsController.CreateNews)
//         // 获取新闻详情
//         newsGroup.GET("/:id", newsController.GetNewsDetail)
//         // 添加评论
//         newsGroup.POST("/:id/comment", newsController.AddComment)

// 		// TODO
// 		// 删除评论
// 		// 删除新闻
// 		// 更新新闻

//         // 点赞新闻
//         newsGroup.POST("/:id/like", newsController.LikeNews)
//         // 取消点赞新闻
//         newsGroup.POST("/:id/cancel_like", newsController.CancelLikeNews)

//         // 收藏新闻
//         newsGroup.POST("/:id/favourite", newsController.FavoriteNews)
//         // 取消收藏新闻
//         newsGroup.POST("/:id/cancel_favourite", newsController.CancelFavoriteNews)

//         // 点踩新闻
//         newsGroup.POST("/:id/dislike", newsController.DislikeNews)
//         // 取消点踩新闻
//         newsGroup.POST("/:id/cancel_dislike", newsController.CancelDislikeNews)

//         newsGroup.POST("/:id/view", newsController.ViewNews)
//     }
// }