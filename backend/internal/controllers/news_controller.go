// controllers/news_controller.go
package controllers

import (
    "net/http"
    "strconv"
    "time"
    "fmt"

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

// 上传新闻
func (nc *NewsController) UploadNews(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 自定义绑定的 JSON 结构
    var newsRequest struct {
        Title        string           `json:"title" binding:"required"`
        Description  string           `json:"description" binding:"required"`
        NewsType     models.NewsType  `json:"news_type" binding:"required"`
        Video        *models.Video    `json:"video,omitempty"`
        Paragraphs   []models.Paragraph `json:"paragraphs,omitempty"`
        Resources    []models.Resource `json:"resources,omitempty"`
        ExternalLink string           `json:"external_link,omitempty"`
        Tags         []string         `json:"tags,omitempty"`
    }

    if err := c.ShouldBindJSON(&newsRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 初始化新闻对象
    news := models.News{
        Title:          newsRequest.Title,
        Description:    newsRequest.Description,
        NewsType:       newsRequest.NewsType,
        AuthorID:       userID.(uint), // 设置作者 ID 为当前用户
        AuthorName:     "",            // 后续通过查询补充
        AuthorAvatar:   "",            // 后续通过查询补充
        UploadTime:     time.Now(),
        ViewCount:      0,
        LikeCount:      0,
        FavoriteCount:  0,
        DislikeCount:   0,
        ShareCount:     0,
        LikedByUsers:   []models.User{},
        FavoritedByUsers: []models.User{},
        DislikedByUsers: []models.User{},
        Tags:           []models.Tag{},
    }

    // 根据 NewsType 进行验证并设置相关字段
    switch news.NewsType {
    case models.NewsTypeVideo:
        if newsRequest.Video == nil || newsRequest.Video.VideoURL == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Video URL is required for video news"})
            return
        }
        news.Video = *newsRequest.Video
    case models.NewsTypeRegular:
        if len(newsRequest.Paragraphs) == 0 && len(newsRequest.Resources) == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "At least one of Paragraphs or Resources is required for regular news"})
            return
        }
        news.Paragraphs = newsRequest.Paragraphs
        news.Resources = newsRequest.Resources
    case models.NewsTypeExternal:
        if newsRequest.ExternalLink == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "External link is required for external news"})
            return
        }
        news.ExternalLink = newsRequest.ExternalLink
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NewsType"})
        return
    }

    // 处理 Tags
    if len(newsRequest.Tags) > 0 {
        var existingTags []models.Tag
        var tagNames []string
        tagNames = append(tagNames, newsRequest.Tags...)

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

    // 补充作者信息
    var author models.User
    if err := nc.DB.First(&author, userID.(uint)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch author information"})
        return
    }
    news.AuthorName = author.Nickname
    news.AuthorAvatar = author.AvatarURL

    // 插入数据库
    if err := nc.DB.Create(&news).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
        return
    }

    // 返回成功响应
    response := struct {
        ID          uint      `json:"id"`
        Title       string    `json:"title"`
        Description string    `json:"description"`
        UploadTime  time.Time `json:"upload_time"`
        ViewCount   int       `json:"view_count"`
        LikeCount   int       `json:"like_count"`
        FavoriteCount int     `json:"favorite_count"`
        DislikeCount int      `json:"dislike_count"`
        ShareCount  int       `json:"share_count"`
        NewsType    models.NewsType `json:"news_type"`
        Tags        []string  `json:"tags"`
    }{
        ID:           news.ID,
        Title:        news.Title,
        Description:  news.Description,
        UploadTime:   news.UploadTime,
        ViewCount:    news.ViewCount,
        LikeCount:    news.LikeCount,
        FavoriteCount: news.FavoriteCount,
        DislikeCount:  news.DislikeCount,
        ShareCount:    news.ShareCount,
        NewsType:     news.NewsType,
        Tags:         newsRequest.Tags,
    }

    c.JSON(http.StatusCreated, response)
}

// 以 id 为索引，获取单个新闻的具体信息
func (nc *NewsController) GetNewsDetail(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userIDInt, ok := userID.(uint)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 查询新闻基础信息
    var news models.News
    if err := nc.DB.Preload("Tags").
        Preload("Author").
        First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 查询用户是否点赞、收藏、点踩了该新闻
    var isLiked, isFavorited, isDisliked int64
    if err := nc.DB.Table("user_likes_news").
        Where("user_id = ? AND news_id = ?", userIDInt, newsID).
        Count(&isLiked).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like status"})
        return
    }
    if err := nc.DB.Table("user_favorites_news").
        Where("user_id = ? AND news_id = ?", userIDInt, newsID).
        Count(&isFavorited).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check favorite status"})
        return
    }
    if err := nc.DB.Table("user_dislikes_news").
        Where("user_id = ? AND news_id = ?", userIDInt, newsID).
        Count(&isDisliked).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check dislike status"})
        return
    }

    // 查询所有一级评论
    var comments []models.Comment
    if err := nc.DB.Where("news_id = ? AND is_reply = ?", newsID, false).
        Preload("Replies").
        Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
        return
    }

    // 根据新闻类型加载特定资源
    var video models.Video
    var paragraphs []models.Paragraph
    var resources []models.Resource
    switch news.NewsType {
    case models.NewsTypeVideo:
        if err := nc.DB.Where("news_id = ?", newsID).First(&video).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load video details"})
            return
        }
    case models.NewsTypeRegular:
        if err := nc.DB.Where("news_id = ?", newsID).Find(&paragraphs).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load paragraphs"})
            return
        }
        if err := nc.DB.Where("news_id = ?", newsID).Find(&resources).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load resources"})
            return
        }
    }

    // 构建响应
    response := gin.H{
        "id":             news.ID,
        "title":          news.Title,
        "description":    news.Description,
        "upload_time":    news.UploadTime,
        "view_count":     news.ViewCount,
        "like_count":     news.LikeCount,
        "favorite_count": news.FavoriteCount,
        "dislike_count":  news.DislikeCount,
        "share_count":    news.ShareCount,
        "news_type":      news.NewsType,
        "author_name":    news.AuthorName,
        "author_avatar":  news.AuthorAvatar,
        "tags":           news.Tags,
        "is_liked":       isLiked,              // 该用户是否点赞了新闻
        "is_favorited":   isFavorited,
        "is_disliked":    isDisliked,
        "comments":       comments,
    }

    // 根据 NewsType 动态添加字段
    switch news.NewsType {
    case models.NewsTypeVideo:
        response["video"] = video
    case models.NewsTypeRegular:
        response["paragraphs"] = paragraphs
        response["resources"] = resources
    case models.NewsTypeExternal:
        response["external_link"] = news.ExternalLink
    }

    // 返回响应
    c.JSON(http.StatusOK, response)
}

// 以某些条件为约束，查看多个新闻预览
func (nc *NewsController) GetNewsPreviews(c *gin.Context) {
    // 解析查询参数
    category := c.Query("category")
    page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
    if err != nil || page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }
    size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
    if err != nil || size < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
        return
    }
    sortBy := c.DefaultQuery("sort_by", "upload_time") // 默认按上传时间排序
    order := c.DefaultQuery("order", "desc")          // 默认降序

    // 校验排序字段和顺序
    allowedSortFields := map[string]bool{
        "upload_time": true,
        "view_count":  true,
        "like_count":  true,
    }
    if !allowedSortFields[sortBy] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort_by parameter"})
        return
    }
    if order != "asc" && order != "desc" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameter"})
        return
    }

    // 计算分页偏移量
    offset := (page - 1) * size

    // 定义预览响应结构
    type NewsPreview struct {
        ID             uint       `json:"id"`
        Title          string     `json:"title"`
        Description    string     `json:"description"`
        UploadTime     time.Time  `json:"upload_time"`
        ViewCount      int        `json:"view_count"`
        ShareCount     int        `json:"share_count"`
        LikeCount      int        `json:"like_count"`
        FavoriteCount  int        `json:"favorite_count"`
        DislikeCount   int        `json:"dislike_count"`
        NewsType       string     `json:"news_type"`
        AuthorName     string     `json:"author_name"`
        AuthorAvatar   string     `json:"author_avatar"`
        Tags           []string   `json:"tags"`
        ExtraField     string     `json:"extra_field,omitempty"` // 根据 NewsType 动态字段 类似于预览图片 TODO
    }

    // 查询新闻列表
    var news []models.News
    query := nc.DB.Preload("Tags").Preload("Author").
        Order(fmt.Sprintf("%s %s", sortBy, order)).
        Offset(offset).
        Limit(size)
    if category != "" {
        query = query.Where("type = ?", category)
    }

    if err := query.Find(&news).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news previews"})
        return
    }

    // 转换为响应格式
    previews := make([]NewsPreview, len(news))
    for i, n := range news {
        tags := []string{}
        for _, tag := range n.Tags {
            tags = append(tags, tag.Name)
        }

        preview := NewsPreview{
            ID:            n.ID,
            Title:         n.Title,
            Description:   n.Description,
            UploadTime:    n.UploadTime,
            ViewCount:     n.ViewCount,
            ShareCount:    n.ShareCount,
            LikeCount:     n.LikeCount,
            FavoriteCount: n.FavoriteCount,
            DislikeCount:  n.DislikeCount,
            NewsType:      string(n.NewsType),
            AuthorName:    n.AuthorName,
            AuthorAvatar:  n.AuthorAvatar,
            Tags:          tags,
        }

        // 根据 NewsType 动态添加字段
        // TODO
        switch n.NewsType {
        case models.NewsTypeVideo:
            preview.ExtraField = "Video-specific data" // 示例
        case models.NewsTypeRegular:
            preview.ExtraField = "Regular-specific data" // 示例
        case models.NewsTypeExternal:
            preview.ExtraField = n.ExternalLink // 示例：返回外部链接
        }

        previews[i] = preview
    }

    // 返回响应
    c.JSON(http.StatusOK, gin.H{
        "total": len(previews),
        "page":  page,
        "size":  size,
        "data":  previews,
    })
}

// 添加评论
func (nc *NewsController) AddComment(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 定义请求体绑定结构，避免前端上传不必要的字段
    var commentRequest struct {
        NewsID   uint   `json:"news_id" binding:"required"`
        Content  string `json:"content" binding:"required"`
        IsReply  bool   `json:"is_reply"`
        ParentID *uint  `json:"parent_id,omitempty"`
    }

    if err := c.ShouldBindJSON(&commentRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 检查 newsID 是否有效
    var newsExists bool
    if err := nc.DB.Model(&models.News{}).Select("count(*) > 0").Where("id = ?", commentRequest.NewsID).Find(&newsExists).Error; err != nil || !newsExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid newsID"})
        return
    }

    // 检查 is_reply 字段逻辑
    if commentRequest.IsReply {
        if commentRequest.ParentID == nil || *commentRequest.ParentID == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ParentID is required for a reply"})
            return
        }

        var parentComment models.Comment
        if err := nc.DB.First(&parentComment, *commentRequest.ParentID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment not found"})
            return
        }

        // 确保 parentComment 的 newsID 与当前 comment 的 newsID 一致
        if parentComment.NewsID != commentRequest.NewsID {
            c.JSON(http.StatusBadRequest, gin.H{"error": "NewsID does not match the parent comment's newsID"})
            return
        }
    } else {
        // 如果不是回复评论，则 parentID 必须为 nil 或 0
        if commentRequest.ParentID != nil && *commentRequest.ParentID != 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ParentID must be 0 for a top-level comment"})
            return
        }
    }

    // 创建评论
    comment := models.Comment{
        NewsID:     commentRequest.NewsID,
        Content:    commentRequest.Content,
        UserID:     userID.(uint), // 从 JWT 获取的用户 ID
        IsReply:    commentRequest.IsReply,
        ParentID:   commentRequest.ParentID,
        PublishTime: time.Now(),
        LikeCount:  0,
        Replies: []models.Comment{},
    }

    if err := nc.DB.Create(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Comment added successfully",
        "comment": comment,
    })
}

// 用户点赞新闻
func (nc *NewsController) LikeNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查新闻是否存在
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查用户是否已点赞
    var likeExists int64
    if err := nc.DB.Table("user_likes_news").
        Where("user_id = ? AND news_id = ?", userID, newsID).
        Count(&likeExists).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like status"})
        return
    }
    if likeExists > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already liked this news"})
        return
    }

    // 建立点赞关系
    if err := nc.DB.Exec("INSERT INTO user_likes_news (user_id, news_id) VALUES (?, ?)", userID, newsID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like news"})
        return
    }

    // 增加点赞数
    if err := nc.DB.Model(&news).Update("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "News liked successfully",
        "like_count": news.LikeCount + 1,
    })
}

// 取消用户点赞新闻
func (nc *NewsController) CancelLikeNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查新闻是否存在
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查用户是否已点赞
    var likeExists int64
    if err := nc.DB.Table("user_likes_news").
        Where("user_id = ? AND news_id = ?", userID, newsID).
        Count(&likeExists).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like status"})
        return
    }
    if likeExists == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not liked this news"})
        return
    }

    // 删除点赞关系
    if err := nc.DB.Exec("DELETE FROM user_likes_news WHERE user_id = ? AND news_id = ?", userID, newsID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel like"})
        return
    }

    // 减少点赞数
    if err := nc.DB.Model(&news).Update("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "News like canceled successfully",
        "like_count": news.LikeCount - 1,
    })
}

// 用户收藏新闻
func (nc *NewsController) FavoriteNews(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查新闻是否存在
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查是否已收藏
    var favoriteExists int64
    if err := nc.DB.Table("user_favorites_news").
        Where("user_id = ? AND news_id = ?", userID, newsID).
        Count(&favoriteExists).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check favorite status"})
        return
    }
    if favoriteExists > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already favorited this news"})
        return
    }

    // 建立收藏关系
    if err := nc.DB.Exec("INSERT INTO user_favorites_news (user_id, news_id) VALUES (?, ?)", userID, newsID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to favorite news"})
        return
    }

    // 更新收藏数
    if err := nc.DB.Model(&news).Update("favorite_count", gorm.Expr("favorite_count + 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorite count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":        "News favorited successfully",
        "favorite_count": news.FavoriteCount + 1,
    })
}

// 取消用户收藏新闻
func (nc *NewsController) CancelFavoriteNews(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查新闻是否存在
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查是否已收藏
    var favoriteExists int64
    if err := nc.DB.Table("user_favorites_news").
        Where("user_id = ? AND news_id = ?", userID, newsID).
        Count(&favoriteExists).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check favorite status"})
        return
    }
    if favoriteExists == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not favorited this news"})
        return
    }

    // 删除收藏关系
    if err := nc.DB.Exec("DELETE FROM user_favorites_news WHERE user_id = ? AND news_id = ?", userID, newsID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel favorite"})
        return
    }

    // 更新收藏数
    if err := nc.DB.Model(&news).Update("favorite_count", gorm.Expr("favorite_count - 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorite count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":        "News favorite canceled successfully",
        "favorite_count": news.FavoriteCount - 1,
    })
}

// 用户点踩新闻
func (nc *NewsController) DislikeNews(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查新闻是否存在
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查是否已点踩
    var dislikeExists int64
    if err := nc.DB.Table("user_dislikes_news").
        Where("user_id = ? AND news_id = ?", userID, newsID).
        Count(&dislikeExists).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check dislike status"})
        return
    }
    if dislikeExists > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already disliked this news"})
        return
    }

    // 建立点踩关系
    if err := nc.DB.Exec("INSERT INTO user_dislikes_news (user_id, news_id) VALUES (?, ?)", userID, newsID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike news"})
        return
    }

    // 更新点踩数
    if err := nc.DB.Model(&news).Update("dislike_count", gorm.Expr("dislike_count + 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dislike count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":        "News disliked successfully",
        "dislike_count":  news.DislikeCount + 1,
    })
}

// 取消用户点踩新闻
func (nc *NewsController) CancelDislikeNews(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查新闻是否存在
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 检查是否已点踩
    var dislikeExists int64
    if err := nc.DB.Table("user_dislikes_news").
        Where("user_id = ? AND news_id = ?", userID, newsID).
        Count(&dislikeExists).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check dislike status"})
        return
    }
    if dislikeExists == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not disliked this news"})
        return
    }

    // 删除点踩关系
    if err := nc.DB.Exec("DELETE FROM user_dislikes_news WHERE user_id = ? AND news_id = ?", userID, newsID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel dislike"})
        return
    }

    // 更新点踩数
    if err := nc.DB.Model(&news).Update("dislike_count", gorm.Expr("dislike_count - 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dislike count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":        "News dislike canceled successfully",
        "dislike_count":  news.DislikeCount - 1,
    })
}

// 用户浏览新闻
func (nc *NewsController) ViewNews(c *gin.Context) {
    // 从 JWT 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 检查新闻是否存在
    var news models.News
    if err := nc.DB.First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 限制浏览记录最大数量
    const maxViewed = 200
    var viewedCount int64
    if err := nc.DB.Table("user_viewed_news").
        Where("user_id = ?", userID).
        Count(&viewedCount).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch view count"})
        return
    }

    // 如果浏览记录已满，则删除最早的记录
    if viewedCount >= maxViewed {
        var oldestNewsID uint
        if err := nc.DB.Table("user_viewed_news").
            Where("user_id = ?", userID).
            Order("created_at ASC").
            Limit(1).
            Select("news_id").
            Scan(&oldestNewsID).Error; err == nil {
            nc.DB.Exec("DELETE FROM user_viewed_news WHERE user_id = ? AND news_id = ?", userID, oldestNewsID)
        }
    }

    // 添加新的浏览记录
    if err := nc.DB.Exec("INSERT INTO user_viewed_news (user_id, news_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE created_at = NOW()", userID, newsID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record news view"})
        return
    }

    // 更新新闻的浏览计数
    if err := nc.DB.Model(&news).Update("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update view count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "News view recorded successfully",
    })
}



// TODO
// 删除评论
func (nc *NewsController) DeleteComment(c *gin.Context) {

}