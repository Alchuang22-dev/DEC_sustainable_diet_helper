// controllers/news_controller.go
package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
)

type NewsController struct {
    DB *gorm.DB
}

func NewNewsController(db *gorm.DB) *NewsController {
    return &NewsController{DB: db}
}

func (nc *NewsController) CreateNews(c *gin.Context) {
    var news models.News
    if err := c.ShouldBindJSON(&news); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 初始化用户列表为空切片，避免返回 null
    news.LikedByUsers = []models.User{}
    news.FavoritedByUsers = []models.User{}
    news.DislikedByUsers = []models.User{}
    news.Tags = []models.Tag{}

    // 检查 NewsType 是否为空
    if news.NewsType == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "NewsType is required"})
        return
    }

    // 检查 AuthorID 是否有效
    var author models.User
    if err := nc.DB.First(&author, news.AuthorID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
        return
    }

    // 根据 NewsType 进行验证
    switch news.NewsType {
    case models.NewsTypeVideo:
        // 如果是视频新闻，需要 Video 字段并且不应包含 Paragraphs 和 Resources
        if news.Video.VideoURL == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Video URL is required for video news"})
            return
        }
        if len(news.Paragraphs) > 0 || len(news.Resources) > 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Paragraphs and Resources are not allowed for video news"})
            return
        }

    case models.NewsTypeRegular:
        // 如果是常规新闻，需要 Paragraphs 或 Resources 字段，不应包含 Video
        if news.Video.VideoURL != "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Video URL is not allowed for regular news"})
            return
        }
        if len(news.Paragraphs) == 0 && len(news.Resources) == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "At least one of Paragraphs or Resources is required for regular news"})
            return
        }

        // 设置 Paragraphs 和 Resources 的 NewsID
        for i := range news.Paragraphs {
            news.Paragraphs[i].NewsID = news.ID
        }
        for i := range news.Resources {
            news.Resources[i].NewsID = news.ID
        }

    case models.NewsTypeExternal:
        // 如果是外部新闻，需要 ExternalLink 字段
        if news.ExternalLink == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "External link is required for external news"})
            return
        }
        if len(news.Paragraphs) > 0 || len(news.Resources) > 0 || news.Video.VideoURL != "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Paragraphs, Resources, and Video are not allowed for external news"})
            return
        }

    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NewsType"})
        return
    }

    // 处理 Tags
    if len(news.Tags) > 0 {
        var existingTags []models.Tag
        var tagNames []string
        for _, tag := range news.Tags {
            tagNames = append(tagNames, tag.Name)
        }

        // 查找现有标签
        nc.DB.Where("name IN ?", tagNames).Find(&existingTags)

        // 创建不存在的标签
        existingTagNames := map[string]bool{}
        for _, tag := range existingTags {
            existingTagNames[tag.Name] = true
        }
        for _, tagName := range tagNames {
            if !existingTagNames[tagName] {
                newTag := models.Tag{Name: tagName}
                if err := nc.DB.Create(&newTag).Error; err != nil {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new tags"})
                    return
                }
                existingTags = append(existingTags, newTag)
            }
        }

        // 关联新闻和标签
        news.Tags = existingTags
    }

    // 插入数据库
    if err := nc.DB.Create(&news).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
        return
    }

    // 返回成功响应
    c.JSON(http.StatusCreated, news)
}

// 获取新闻详情（包括评论、特定资源和关联用户列表）
func (nc *NewsController) GetNewsDetail(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 查询基础新闻信息并加载公共关联字段
    var news models.News
    if err := nc.DB.
        Preload("Comments").
        Preload("LikedByUsers").
        Preload("FavoritedByUsers").
        Preload("DislikedByUsers").
        First(&news, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 根据新闻类型加载特定字段
    switch news.NewsType {
    case models.NewsTypeVideo:
        if err := nc.DB.Preload("Video").First(&news, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load video news details"})
            return
        }
    case models.NewsTypeRegular:
        if err := nc.DB.Preload("Paragraphs").Preload("Resources").First(&news, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load regular news details"})
            return
        }
    }

    // 确保关联字段不为 null
    if news.LikedByUsers == nil {
        news.LikedByUsers = []models.User{}
    }
    if news.FavoritedByUsers == nil {
        news.FavoritedByUsers = []models.User{}
    }
    if news.DislikedByUsers == nil {
        news.DislikedByUsers = []models.User{}
    }
    if news.Comments == nil {
        news.Comments = []models.Comment{}
    }
    if news.NewsType == models.NewsTypeRegular {
        if news.Paragraphs == nil {
            news.Paragraphs = []models.Paragraph{}
        }
        if news.Resources == nil {
            news.Resources = []models.Resource{}
        }
    } else if news.NewsType == models.NewsTypeVideo && news.Video.ID == 0 {
        news.Video = models.Video{}
    }

    // 返回新闻详情
    c.JSON(http.StatusOK, news)
}

// 添加评论
func (nc *NewsController) AddComment(c *gin.Context) {
    var comment models.Comment
    if err := c.ShouldBindJSON(&comment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 检查 newsID 是否有效
    var newsExists bool
    if err := nc.DB.Model(&models.News{}).Select("count(*) > 0").Where("id = ?", comment.NewsID).Find(&newsExists).Error; err != nil || !newsExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid newsID"})
        return
    }

    // 检查 userID 是否有效
    var userExists bool
    if err := nc.DB.Model(&models.User{}).Select("count(*) > 0").Where("id = ?", comment.UserID).Find(&userExists).Error; err != nil || !userExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
        return
    }

    // 检查 is_reply 字段，如果是回复评论，则 parentID 不应为 0
    if comment.IsReply {
        var parentComment models.Comment
        if err := nc.DB.First(&parentComment, comment.ParentID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment not found"})
            return
        }

        // 确保 parentComment 的 newsID 与当前 comment 的 newsID 一致
        if parentComment.NewsID != comment.NewsID {
            c.JSON(http.StatusBadRequest, gin.H{"error": "newsID does not match the parent comment's newsID"})
            return
        }
    } else {
        // 如果不是回复评论，则 parentID 必须为 0
        if comment.ParentID != nil && *comment.ParentID != 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "parentID must be 0 for a top-level comment"})
            return
        }
    }

    // 创建评论
    if err := nc.DB.Create(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, comment)
}

// TODO
// 删除评论
func (nc *NewsController) DeleteComment(c *gin.Context) {

}

// 删除新闻
func (nc *NewsController) DeleteNews(c *gin.Context) {

}

// 更新新闻
func (nc *NewsController) UpdateNews(c *gin.Context) {

}

// 用户点赞新闻
func (nc *NewsController) LikeNews(c *gin.Context) {
    userID, err := strconv.Atoi(c.GetHeader("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 获取用户和新闻记录
    var user models.User
    if err := nc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
    }
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "News not found"})
        return
    }

    // 使用原生 SQL 检查是否已点赞
    var likeExists int64
    nc.DB.Table("user_likes_news").Where("user_id = ? AND news_id = ?", userID, newsID).Count(&likeExists)
    if likeExists > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already liked this news"})
        return
    }

    // 建立点赞关系
    if err := nc.DB.Model(&user).Association("LikedNews").Append(&news); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like news"})
        return
    }
    
    // 获取最新的点赞数
    likeCount := nc.DB.Model(&news).Association("LikedByUsers").Count()

    c.JSON(http.StatusOK, gin.H{
        "message": "News liked successfully",
        "like_count": likeCount,
    })
}

// 取消用户点赞新闻
func (nc *NewsController) CancelLikeNews(c *gin.Context) {
    userID, err := strconv.Atoi(c.GetHeader("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 获取用户和新闻记录
    var user models.User
    if err := nc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
    }
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "News not found"})
        return
    }

    // 使用原生 SQL 检查是否已点赞
    var likeExists int64
    nc.DB.Table("user_likes_news").Where("user_id = ? AND news_id = ?", userID, newsID).Count(&likeExists)
    if likeExists == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not liked this news"})
        return
    }

    // 删除点赞关系
    if err := nc.DB.Model(&user).Association("LikedNews").Delete(&news); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel like"})
        return
    }

    // 获取最新的点赞数
    likeCount := nc.DB.Model(&news).Association("LikedByUsers").Count()

    c.JSON(http.StatusOK, gin.H{
        "message": "News like canceled successfully",
        "like_count": likeCount,
    })
}

// 用户收藏新闻
func (nc *NewsController) FavoriteNews(c *gin.Context) {
    userID, err := strconv.Atoi(c.GetHeader("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查用户和新闻是否存在
    var user models.User
    if err := nc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 使用原生 SQL 检查是否已收藏
    var favoriteExists int64
    nc.DB.Table("user_favorites_news").Where("user_id = ? AND news_id = ?", userID, newsID).Count(&favoriteExists)
    if favoriteExists > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already favorited this news"})
        return
    }

    // 建立收藏关系
    if err := nc.DB.Model(&user).Association("FavoritedNews").Append(&news); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to favorite news"})
        return
    }

    // 获取最新的收藏数
    favoriteCount := nc.DB.Model(&news).Association("FavoritedByUsers").Count()

    c.JSON(http.StatusOK, gin.H{
        "message": "News favorited successfully",
        "favorite_count": favoriteCount,
    })
}

// 取消用户收藏新闻
func (nc *NewsController) CancelFavoriteNews(c *gin.Context) {
    userID, err := strconv.Atoi(c.GetHeader("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查用户和新闻是否存在
    var user models.User
    if err := nc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查是否已收藏
    var favoriteExists int64
    nc.DB.Table("user_favorites_news").Where("user_id = ? AND news_id = ?", userID, newsID).Count(&favoriteExists)
    if favoriteExists == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not favorited this news"})
        return
    }

    // 删除收藏关系
    if err := nc.DB.Model(&user).Association("FavoritedNews").Delete(&news); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel favorite"})
        return
    }

    // 获取最新的收藏数
    favoriteCount := nc.DB.Model(&news).Association("FavoritedByUsers").Count()

    c.JSON(http.StatusOK, gin.H{
        "message": "News favorite canceled successfully",
        "favorite_count": favoriteCount,
    })
}

// 用户点踩新闻
func (nc *NewsController) DislikeNews(c *gin.Context) {
    userID, err := strconv.Atoi(c.GetHeader("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查用户和新闻是否存在
    var user models.User
    if err := nc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 使用原生 SQL 检查是否已点踩
    var dislikeExists int64
    nc.DB.Table("user_dislikes_news").Where("user_id = ? AND news_id = ?", userID, newsID).Count(&dislikeExists)
    if dislikeExists > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already disliked this news"})
        return
    }

    // 建立点踩关系
    if err := nc.DB.Model(&user).Association("DislikedNews").Append(&news); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike news"})
        return
    }

    // 获取最新的点踩数
    dislikeCount := nc.DB.Model(&news).Association("DislikedByUsers").Count()

    c.JSON(http.StatusOK, gin.H{
        "message": "News disliked successfully",
        "dislike_count": dislikeCount,
    })
}

// 取消用户点踩新闻
func (nc *NewsController) CancelDislikeNews(c *gin.Context) {
    userID, err := strconv.Atoi(c.GetHeader("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查用户和新闻是否存在
    var user models.User
    if err := nc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查是否已点踩
    var dislikeExists int64
    nc.DB.Table("user_dislikes_news").Where("user_id = ? AND news_id = ?", userID, newsID).Count(&dislikeExists)
    if dislikeExists == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not disliked this news"})
        return
    }

    // 删除点踩关系
    if err := nc.DB.Model(&user).Association("DislikedNews").Delete(&news); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel dislike"})
        return
    }

    // 获取最新的点踩数
    dislikeCount := nc.DB.Model(&news).Association("DislikedByUsers").Count()

    c.JSON(http.StatusOK, gin.H{
        "message": "News dislike canceled successfully",
        "dislike_count": dislikeCount,
    })
}

// 用户浏览新闻
func (nc *NewsController) ViewNews(c *gin.Context) {
    userID, err := strconv.Atoi(c.GetHeader("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查用户和新闻是否存在
    var user models.User
    if err := nc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 获取当前浏览记录数
    const maxViewed = 200
    viewedCount := nc.DB.Model(&user).Association("ViewedNews").Count()

    // 如果浏览记录已满，则删除最早的记录
    if viewedCount >= maxViewed {
        var oldestNews models.News
        // 查找用户最早浏览的新闻记录
        if err := nc.DB.Table("user_viewed_news").
            Where("user_id = ?", userID).
            Order("created_at ASC").
            Limit(1).
            Select("news_id").
            Scan(&oldestNews).Error; err == nil {
            // 从关联中移除最早的记录
            nc.DB.Model(&user).Association("ViewedNews").Delete(&oldestNews)
        }
    }

    // 添加新的浏览记录
    if err := nc.DB.Model(&user).Association("ViewedNews").Append(&news); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record news view"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "News view recorded successfully",
    })
}
