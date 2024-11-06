// controllers/user_controller.go
package controllers

import (
    "net/http"
    "strconv"
    "log"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
)

type UserController struct {
    DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
    return &UserController{DB: db}
}

// 获取所有用户
func (uc *UserController) GetAllUsers(c *gin.Context) {
    log.Println("GetAllUsers 被调用")
    var users []models.User
    if err := uc.DB.Preload("LikedNews").
        Preload("FavoritedNews").
        Preload("DislikedNews").
        Preload("ViewedNews").
        Find(&users).Error; err != nil {
        log.Println("获取所有用户失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 初始化空列表，避免返回 null
    for i := range users {
        if users[i].LikedNews == nil {
            users[i].LikedNews = []models.News{}
        }
        if users[i].FavoritedNews == nil {
            users[i].FavoritedNews = []models.News{}
        }
        if users[i].DislikedNews == nil {
            users[i].DislikedNews = []models.News{}
        }
        if users[i].ViewedNews == nil {
            users[i].ViewedNews = []models.News{}
        }
    }

    c.JSON(http.StatusOK, users)
}

// 创建新用户
func (uc *UserController) CreateUser(c *gin.Context) {
    log.Println("CreateUser 被调用")
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        log.Println("绑定JSON失败:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 初始化空关联列表
    user.LikedNews = []models.News{}
    user.FavoritedNews = []models.News{}
    user.DislikedNews = []models.News{}
    user.ViewedNews = []models.News{}

    if err := uc.DB.Create(&user).Error; err != nil {
        log.Println("创建用户失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Println("用户创建成功:", user)
    c.JSON(http.StatusCreated, user)
}

// 根据ID获取用户
func (uc *UserController) GetUserByID(c *gin.Context) {
    log.Println("GetUserByID 被调用")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        log.Println("无效的用户ID:", c.Param("id"))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var user models.User
    if err := uc.DB.Preload("LikedNews").
        Preload("FavoritedNews").
        Preload("DislikedNews").
        Preload("ViewedNews").
        First(&user, id).Error; err != nil {
        log.Println("用户未找到:", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 确保空列表初始化
    if user.LikedNews == nil {
        user.LikedNews = []models.News{}
    }
    if user.FavoritedNews == nil {
        user.FavoritedNews = []models.News{}
    }
    if user.ViewedNews == nil {
        user.ViewedNews = []models.News{}
    }

    c.JSON(http.StatusOK, user)
}

// 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
    log.Println("UpdateUser 被调用")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        log.Println("无效的用户ID:", c.Param("id"))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var user models.User
    // 预加载关联数据
    if err := uc.DB.Preload("LikedNews").
        Preload("FavoritedNews").
        Preload("ViewedNews").
        Preload("DislikedNews").
        First(&user, id).Error; err != nil {
        log.Println("用户未找到:", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 绑定更新的数据，但不修改关联字段
    var updatedUser models.User
    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        log.Println("绑定JSON失败:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 更新非关联字段
    user.Name = updatedUser.Name
    user.Email = updatedUser.Email

    if err := uc.DB.Save(&user).Error; err != nil {
        log.Println("更新用户失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Println("用户更新成功:", user)
    c.JSON(http.StatusOK, user)
}

// 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
    log.Println("DeleteUser 被调用")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        log.Println("无效的用户ID:", c.Param("id"))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var user models.User
    if err := uc.DB.First(&user, id).Error; err != nil {
        log.Println("用户未找到:", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 删除用户的关联记录
    if err := uc.DB.Model(&user).Association("LikedNews").Clear(); err != nil {
        log.Println("清除点赞记录失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear liked news"})
        return
    }
    if err := uc.DB.Model(&user).Association("FavoritedNews").Clear(); err != nil {
        log.Println("清除收藏记录失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear favorited news"})
        return
    }
    if err := uc.DB.Model(&user).Association("DislikedNews").Clear(); err != nil {
        log.Println("清除点踩记录失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear disliked news"})
        return
    }
    if err := uc.DB.Model(&user).Association("ViewedNews").Clear(); err != nil {
        log.Println("清除浏览记录失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear viewed news"})
        return
    }

    // 删除用户记录
    if err := uc.DB.Delete(&user).Error; err != nil {
        log.Println("删除用户失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Println("用户删除成功:", id)
    c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}