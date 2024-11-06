// controllers/news_controller.go
package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/models"
)

type NewsController struct {
    DB *gorm.DB
}

func NewNewsController(db *gorm.DB) *NewsController {
    return &NewsController{DB: db}
}

// 创建视频新闻
func (nc *NewsController) CreateVideoNews(c *gin.Context) {
    var videoNews models.VideoNews
    if err := c.ShouldBindJSON(&videoNews); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    videoNews.NewsType = models.NewsTypeVideo // 设置新闻类型
    if err := nc.DB.Create(&videoNews).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, videoNews)
}

// 创建常规新闻
func (nc *NewsController) CreateRegularNews(c *gin.Context) {
    var regularNews models.RegularNews
    if err := c.ShouldBindJSON(&regularNews); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    regularNews.NewsType = models.NewsTypeRegular // 设置新闻类型
    if err := nc.DB.Create(&regularNews).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, regularNews)
}

// 获取新闻详情（包括评论）
func (nc *NewsController) GetNewsDetail(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 可以根据类型分别处理，这里以视频新闻为例
    var videoNews models.VideoNews
    if err := nc.DB.Preload("Comments").First(&videoNews, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }
    c.JSON(http.StatusOK, videoNews)
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
    if err := nc.DB.Model(&models.BaseNews{}).Select("count(*) > 0").Where("id = ?", comment.NewsID).Find(&newsExists).Error; err != nil || !newsExists {
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