// controllers/user_controller.go
package controllers

import (
    "net/http"
    "strconv"
    "log"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "DEC/internal/models"
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
    if err := uc.DB.Find(&users).Error; err != nil {
        log.Println("获取所有用户失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
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
    if err := uc.DB.First(&user, id).Error; err != nil {
        log.Println("用户未找到:", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
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
    if err := uc.DB.First(&user, id).Error; err != nil {
        log.Println("用户未找到:", id)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if err := c.ShouldBindJSON(&user); err != nil {
        log.Println("绑定JSON失败:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

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

    if err := uc.DB.Delete(&models.User{}, id).Error; err != nil {
        log.Println("删除用户失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Println("用户删除成功:", id)
    c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}