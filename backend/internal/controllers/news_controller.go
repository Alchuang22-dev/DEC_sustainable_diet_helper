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

// 创建视频新闻
func (nc *NewsController) CreateVideoNews(c *gin.Context) {
    var videoNews models.VideoNews
    if err := c.ShouldBindJSON(&videoNews); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
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

// TODO 这里好像没给 news 对应地方加？
// 添加评论
func (nc *NewsController) AddComment(c *gin.Context) {
    var comment models.Comment
    if err := c.ShouldBindJSON(&comment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
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