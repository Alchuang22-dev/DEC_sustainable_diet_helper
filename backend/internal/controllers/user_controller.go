// controllers/user_controller.go
package controllers

import (
    "net/http"
    "strconv"
    "log"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
)

type UserController struct {
    DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
    return &UserController{DB: db}
}

// 注册
func (uc *UserController) Register(c *gin.Context) {
    log.Println("Register 被调用")

    // 定义请求体结构
    var registerRequest struct {
        Nickname    string `json:"nickname"`                              // 昵称可选
        PhoneNumber string `json:"phone_number" binding:"required"`       // 手机号必填
        Password    string `json:"password" binding:"required,min=6"`     // 密码必填，长度至少 6 位
    }

    // 绑定请求体
    if err := c.ShouldBindJSON(&registerRequest); err != nil {
        log.Println("绑定JSON失败:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 验证手机号格式（简单示例）
    if !utils.IsValidPhoneNumber(registerRequest.PhoneNumber) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
        return
    }

    // 检查手机号是否已注册
    var existingUser models.User
    if err := uc.DB.Where("phone_number = ?", registerRequest.PhoneNumber).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
        return
    }

    // 加密密码
    hashedPassword, err := utils.HashPassword(registerRequest.Password)
    if err != nil {
        log.Println("密码加密失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }

    // 如果未提供昵称，生成随机昵称
    nickname := registerRequest.Nickname
    if nickname == "" {
        nickname = utils.GenerateRandomNickname()
    }

    // 创建新用户
    user := models.User{
        Nickname:    nickname,
        PhoneNumber: registerRequest.PhoneNumber,
        Password:    hashedPassword, // 存储加密后的密码
        LikedNews:     []models.News{},
        FavoritedNews: []models.News{},
        DislikedNews:  []models.News{},
        ViewedNews:    []models.News{},
    }

    if err := uc.DB.Create(&user).Error; err != nil {
        log.Println("创建用户失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    log.Println("用户创建成功:", user)
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 登录
func (uc *UserController) Login(c *gin.Context) {
    var loginRequest struct {
        PhoneNumber string `json:"phone_number" binding:"required"` // 必须提供手机号
        Password    string `json:"password" binding:"required"`    // 必须提供密码
    }

    // 绑定请求体
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // 验证手机号格式（可选）
    if !utils.IsValidPhoneNumber(loginRequest.PhoneNumber) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
        return
    }

    // 验证用户是否存在
    var user models.User
    if err := uc.DB.Where("phone_number = ?", loginRequest.PhoneNumber).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // 用户不存在
        return
    }

    // 验证密码
    if err := utils.CheckPassword(user.Password, loginRequest.Password); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // 密码不匹配
        return
    }

    // 生成 JWT
    token, err := utils.GenerateJWT(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // 返回成功响应
    c.JSON(http.StatusOK, gin.H{"token": token})
}

// 设置用户名
func (uc *UserController) SetNickname(c *gin.Context) {
    log.Println("SetNickname 被调用")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var request struct {
        Nickname string `json:"nickname" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var user models.User
    if err := uc.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.Nickname = request.Nickname

    if err := uc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nickname"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Nickname updated successfully", "nickname": user.Nickname})
}

// 设置手机号
func (uc *UserController) SetPhoneNumber(c *gin.Context) {
    log.Println("SetPhoneNumber 被调用")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var request struct {
        PhoneNumber string `json:"phone_number" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 验证手机号格式
    if !utils.IsValidPhoneNumber(request.PhoneNumber) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
        return
    }

    var user models.User
    if err := uc.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 检查手机号唯一性
    var existingUser models.User
    if err := uc.DB.Where("phone_number = ?", request.PhoneNumber).First(&existingUser).Error; err == nil && existingUser.ID != uint(id) {
        c.JSON(http.StatusConflict, gin.H{"error": "Phone number already in use"})
        return
    }

    user.PhoneNumber = request.PhoneNumber

    if err := uc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update phone number"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Phone number updated successfully", "phone_number": user.PhoneNumber})
}

// 设置密码
func (uc *UserController) SetPassword(c *gin.Context) {
    log.Println("SetPassword 被调用")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var request struct {
        Password string `json:"password" binding:"required,min=6"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var user models.User
    if err := uc.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 加密密码
    hashedPassword, err := utils.HashPassword(request.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }

    user.Password = hashedPassword

    if err := uc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
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