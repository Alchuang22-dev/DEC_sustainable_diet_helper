// controllers/user_controller.go
package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
    DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
    return &UserController{DB: db}
}

// 微信注册/登录
func (uc *UserController) WeChatAuth(c *gin.Context) {
    log.Println("WeChatAuth 被调用")

    // 定义请求体结构
    var authRequest struct {
        Code     string `json:"code" binding:"required"` // 微信登录凭证
        Nickname string `json:"nickname"`               // 可选用户昵称
    }

    // 绑定请求体
    if err := c.ShouldBindJSON(&authRequest); err != nil {
        log.Println("绑定JSON失败:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 调用微信 API 获取 open_id 和 session_key
    // 获取微信 API URL，优先使用环境变量
    wechatAPIURL := os.Getenv("WECHAT_API_URL")
    if wechatAPIURL == "" {
        wechatAPIURL = "https://api.weixin.qq.com/sns/jscode2session"
    }

    // 构建请求 URL
    wxAPI := fmt.Sprintf(
        "%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        wechatAPIURL, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), authRequest.Code,
    )

    resp, err := http.Get(wxAPI)
    if err != nil {
        log.Println("调用微信 API 失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call WeChat API"})
        return
    }
    defer resp.Body.Close()

    var wxResponse struct {
        OpenID     string `json:"open_id"`
        SessionKey string `json:"session_key"`
        ErrCode    int    `json:"errcode"`
        ErrMsg     string `json:"errmsg"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&wxResponse); err != nil {
        log.Println("解析微信 API 响应失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse WeChat API response"})
        return
    }

    // 微信 API 错误
    if wxResponse.ErrCode != 0 {
        log.Println("微信 API 返回错误:", wxResponse.ErrMsg)
        c.JSON(http.StatusUnauthorized, gin.H{"error": wxResponse.ErrMsg})
        return
    }

    // 检查用户是否已存在
    var user models.User

    if err := uc.DB.Where("open_id = ?", wxResponse.OpenID).First(&user).Error; err == nil {
        // 用户已存在，更新 SessionKey
        user.SessionKey = wxResponse.SessionKey
        if err := uc.DB.Save(&user).Error; err != nil {
            log.Println("更新用户 SessionKey 失败:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user session"})
            return
        }
    } else if err == gorm.ErrRecordNotFound { // 表为空或没有对应条目
        // 用户不存在，创建新用户
        user = models.User{
            OpenID:     wxResponse.OpenID,
            SessionKey: wxResponse.SessionKey,
            Nickname:   authRequest.Nickname,
            FamilyID:   0,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
            LikedNews:  []models.News{},
            FavoritedNews:  []models.News{},
            DislikedNews:   []models.News{},
            ViewedNews: []models.News{},
        }

        // 如果未提供昵称，生成随机昵称
        if user.Nickname == "" {
            user.Nickname = utils.GenerateRandomNickname()
        }

        // 初始化头像路径为默认值
        BaseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
        if BaseUploadPath == "" {
            BaseUploadPath = "./uploads" // 默认路径
        }
        timestamp := time.Now().Unix()
        relativePath := fmt.Sprintf("avatars/%d_%d.jpg", user.ID, timestamp)
        user.AvatarURL = relativePath

        // 创建用户
        if err := uc.DB.Create(&user).Error; err != nil {
            log.Println("创建用户失败:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }

        // 拷贝默认头像到新用户头像路径
        defaultAvatarPath := fmt.Sprintf("%s/avatars/default.jpg", BaseUploadPath) // 默认头像路径
        newAvatarPath := fmt.Sprintf("%s/%s", BaseUploadPath, relativePath)
        if err := utils.CopyFile(defaultAvatarPath, newAvatarPath); err != nil {
            log.Printf("复制默认头像失败: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default avatar"})
            return
        }
    } else {
        log.Println("查询数据库时发生错误:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // 生成 JWT
    token, err := utils.GenerateJWT(user.ID)
    if err != nil {
        log.Println("生成 JWT 失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // 返回成功响应
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":         user.ID,
            "nickname":   user.Nickname,
            "avatar_url": user.AvatarURL,
        },
    })
}

// 设置用户名
func (uc *UserController) SetNickname(c *gin.Context) {
    log.Println("SetNickname 被调用")

    // 从上下文中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 解析请求体
    var request struct {
        Nickname string `json:"nickname" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 查询用户并更新昵称
    var user models.User
    if err := uc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.Nickname = request.Nickname
    user.UpdatedAt = time.Now()
    if err := uc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nickname"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Nickname updated successfully", "nickname": user.Nickname})
}

// 设置头像
func (uc *UserController) SetAvatar(c *gin.Context) {
    log.Println("SetAvatar 被调用")

    // 获取用户 ID
    id, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 检查用户是否存在
    var user models.User
    if err := uc.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 获取上传的文件
    file, err := c.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
        return
    }

    // 获取基本路径
    BaseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
    if BaseUploadPath == "" {
        BaseUploadPath = "./uploads" // 默认路径
    }

    // 保存文件到服务器
    timestamp := time.Now().Unix()
    savePath := fmt.Sprintf("%s/avatars/%d_%d.jpg", BaseUploadPath, user.ID, timestamp) // 文件路径
    if err := c.SaveUploadedFile(file, savePath); err != nil {
        log.Println("文件保存失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // 删除旧头像文件
    if user.AvatarURL != "" {
        oldPath := fmt.Sprintf("%s/%s", BaseUploadPath, user.AvatarURL)
        if err := os.Remove(oldPath); err != nil {
            log.Printf("无法删除旧头像文件: %s, 错误: %v\n", oldPath, err)
        }
    }

    fmt.Printf("\n\ntimestamp")
    fmt.Println(timestamp)

    // 更新用户头像路径
    relativePath := fmt.Sprintf("avatars/%d_%d", user.ID, timestamp)
    user.AvatarURL = relativePath
    user.UpdatedAt = time.Now()
    if err := uc.DB.Save(&user).Error; err != nil {
        log.Println("更新用户头像失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
        return
    }

    var x models.User
    uc.DB.Where("id=1").First((&x))
    fmt.Println((x))
    fmt.Println(relativePath)

    log.Println("用户头像设置成功:", relativePath)
    c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully", "avatar_url": relativePath})
}