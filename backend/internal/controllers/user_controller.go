// controllers/user_controller.go
package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
    "strconv"
    "path/filepath"
    "errors"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
    DB *gorm.DB
    Utils utils.UtilsInterface
}

func NewUserController(db *gorm.DB, utils utils.UtilsInterface) *UserController {
    return &UserController{
        DB: db,
        Utils: utils,
    }
}

// WeChatAuth 用户注册/登录
// func (uc *UserController) WeChatAuth(c *gin.Context) {
//     log.Println("WeChatAuth 被调用")

//     code := c.PostForm("code")
//     nickname := c.PostForm("nickname")
    
//     if code == "" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
//         return
//     }

//     // 获取上传的文件
//     file, _ := c.FormFile("avatar") // 头像文件是可选的

//     // 调用微信 API 获取 open_id 和 session_key
//     wechatAPIURL := os.Getenv("WECHAT_API_URL")
//     if wechatAPIURL == "" {
//         wechatAPIURL = "https://api.weixin.qq.com/sns/jscode2session"
//     }

//     // 构建请求 URL
//     wxAPI := fmt.Sprintf(
//         "%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
//         wechatAPIURL, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), code,
//     )

//     resp, err := http.Get(wxAPI)
//     if err != nil {
//         log.Println("调用微信 API 失败:", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call WeChat API"})
//         return
//     }

//     defer resp.Body.Close()

//     var wxResponse struct {
//         OpenID     string `json:"openid"`
//         SessionKey string `json:"session_key"`
//         ErrCode    int    `json:"errcode"`
//         ErrMsg     string `json:"errmsg"`
//     }
//     if err := json.NewDecoder(resp.Body).Decode(&wxResponse); err != nil {
//         log.Println("解析微信 API 响应失败:", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse WeChat API response"})
//         return
//     }

//     // 微信 API 错误
//     if wxResponse.ErrCode != 0 {
//         log.Println("微信 API 返回错误:", wxResponse.ErrMsg)
//         c.JSON(http.StatusUnauthorized, gin.H{"error": wxResponse.ErrMsg})
//         return
//     }

//     // 检查用户是否已存在
//     var user models.User

//     if err := uc.DB.Preload("RefreshTokens").Where("open_id = ?", wxResponse.OpenID).First(&user).Error; err == nil {
//         // 用户已存在，更新 SessionKey
//         user.SessionKey = wxResponse.SessionKey

//         // 更新昵称（如果提供）
//         if nickname != "" && nickname != user.Nickname {
//             user.Nickname = nickname
//         }

//         // 处理头像更新（如果上传了新头像）
//         BaseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
//         if BaseUploadPath == "" {
//             BaseUploadPath = "./uploads"
//         }

//         var relativePath string
//         if file != nil {
//             // 保存用户上传的头像
//             timestamp := time.Now().Unix()
//             relativePath = fmt.Sprintf("avatars/%d_%d.jpg", user.ID, timestamp)
//             savePath := filepath.Join(BaseUploadPath, relativePath)
//             if err := c.SaveUploadedFile(file, savePath); err != nil {
//                 log.Println("保存用户上传的头像失败:", err)
//                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded avatar"})
//                 return
//             }

//             // 删除旧头像文件（如果存在且不是默认头像）
//             if user.AvatarURL != "" && user.AvatarURL != "avatars/default.jpg" {
//                 oldPath := filepath.Join(BaseUploadPath, user.AvatarURL)
//                 if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
//                     log.Printf("无法删除旧头像文件: %s, 错误: %v\n", oldPath, err)
//                 }
//             }

//             // 更新 AvatarURL
//             user.AvatarURL = relativePath
//         }

//         // 如果未上传新头像，但提供了昵称，则仅更新昵称
//         if nickname != "" || file != nil {
//             user.UpdatedAt = time.Now()
//             if err := uc.DB.Save(&user).Error; err != nil {
//                 log.Println("更新用户信息失败:", err)
//                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user information"})
//                 return
//             }
//         }
//     } else if err == gorm.ErrRecordNotFound { // 用户不存在
//         // 用户不存在，创建新用户
//         user = models.User{
//             OpenID:          wxResponse.OpenID,
//             SessionKey:      wxResponse.SessionKey,
//             Nickname:        nickname,
//             FamilyID:        nil,
//             PendingFamilyID: nil,
//             RefreshTokens:   []models.RefreshToken{},
//             CreatedAt:       time.Now(),
//             UpdatedAt:       time.Now(),
//             LikedNews:       []models.News{},
//             FavoritedNews:   []models.News{},
//             DislikedNews:    []models.News{},
//             ViewedNews:      []models.News{},
//         }

//         // 如果未提供昵称，生成随机昵称
//         if user.Nickname == "" {
//             user.Nickname = utils.GenerateRandomNickname()
//         }

//         // 创建用户，获取 user.ID
//         if err := uc.DB.Create(&user).Error; err != nil {
//             log.Println("创建用户失败:", err)
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
//             return
//         }

//         // 处理头像
//         BaseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
//         if BaseUploadPath == "" {
//             BaseUploadPath = "./uploads"
//         }

//         var relativePath string
//         if file != nil {
//             // 保存用户上传的头像
//             timestamp := time.Now().Unix()
//             relativePath = fmt.Sprintf("avatars/%d_%d.jpg", user.ID, timestamp)
//             savePath := filepath.Join(BaseUploadPath, relativePath)
//             if err := c.SaveUploadedFile(file, savePath); err != nil {
//                 log.Println("保存用户上传的头像失败:", err)
//                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded avatar"})
//                 return
//             }
//         } else {
//             // 使用默认头像
//             relativePath = "avatars/default.jpg"
//             defaultAvatarPath := filepath.Join(BaseUploadPath, "avatars", "default.jpg")
//             newAvatarPath := filepath.Join(BaseUploadPath, relativePath)
//             if err := utils.CopyFile(defaultAvatarPath, newAvatarPath); err != nil {
//                 log.Printf("复制默认头像失败: %v\n", err)
//                 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default avatar"})
//                 return
//             }
//         }

//         // 设置头像 URL
//         user.AvatarURL = relativePath

//         // 更新用户的 AvatarURL
//         if err := uc.DB.Save(&user).Error; err != nil {
//             log.Println("更新用户头像失败:", err)
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user avatar"})
//             return
//         }
//     } else {
//         log.Println("查询数据库时发生错误:", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//         return
//     }

//     // 生成 Access Token
//     accessToken, err := utils.GenerateAccessToken(user.ID)
//     if err != nil {
//         log.Println("生成 Access Token 失败:", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
//         return
//     }

//     // 生成 Refresh Token
//     refreshToken, err := utils.GenerateRefreshToken(user.ID)
//     if err != nil {
//         log.Println("生成 Refresh Token 失败:", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
//         return
//     }

//     // 存储 Refresh Token 到数据库
//     newRefreshToken := models.RefreshToken{
//         Token:     refreshToken,
//         UserID:    user.ID,
//         ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
//         Revoked:   false,
//     }
//     if err := uc.DB.Create(&newRefreshToken).Error; err != nil {
//         log.Println("存储 Refresh Token 失败:", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
//         return
//     }

//     // 返回成功响应，包括 Access Token 和 Refresh Token
//     c.JSON(http.StatusOK, gin.H{
//         "access_token":  accessToken,
//         "refresh_token": refreshToken,
//         "user": gin.H{
//             "id":         user.ID,
//             "nickname":   user.Nickname,
//             "avatar_url": user.AvatarURL,
//         },
//     })
// }
func (uc *UserController) WeChatAuth(c *gin.Context) {
    log.Println("WeChatAuth 被调用")

    var request struct {
        Code string `json:"code" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    code := request.Code

    // 调用微信 API 获取 open_id 和 session_key
    wechatAPIURL := os.Getenv("WECHAT_API_URL")
    if wechatAPIURL == "" {
        wechatAPIURL = "https://api.weixin.qq.com/sns/jscode2session"
    }

    // 构建请求 URL
    wxAPI := fmt.Sprintf(
        "%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        wechatAPIURL, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), code,
    )

    resp, err := http.Get(wxAPI)
    if err != nil {
        log.Println("调用微信 API 失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call WeChat API"})
        return
    }

    defer resp.Body.Close()

    var wxResponse struct {
        OpenID     string `json:"openid"`
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

    if err := uc.DB.Preload("RefreshTokens").Where("open_id = ?", wxResponse.OpenID).First(&user).Error; err == nil {
        // 用户已存在，更新 SessionKey
        user.SessionKey = wxResponse.SessionKey
    } else if err == gorm.ErrRecordNotFound { // 用户不存在
        // 用户不存在，创建新用户
        user = models.User{
            OpenID:          wxResponse.OpenID,
            SessionKey:      wxResponse.SessionKey,
            Nickname:        utils.GenerateRandomNickname(),
            FamilyID:        nil,
            PendingFamilyID: nil,
            RefreshTokens:   []models.RefreshToken{},
            CreatedAt:       time.Now(),
            UpdatedAt:       time.Now(),
            LikedNews:       []models.News{},
            FavoritedNews:   []models.News{},
            DislikedNews:    []models.News{},
            ViewedNews:      []models.News{},
        }

        // 创建用户，获取 user.ID
        if err := uc.DB.Create(&user).Error; err != nil {
            log.Println("创建用户失败:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }

        // 处理头像
        BaseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
        if BaseUploadPath == "" {
            BaseUploadPath = "./uploads"
        }

        // 使用默认头像
        relativePath := "avatars/default.jpg"
        defaultAvatarPath := filepath.Join(BaseUploadPath, "avatars", "default.jpg")
        newAvatarPath := filepath.Join(BaseUploadPath, relativePath)
        if err := uc.Utils.CopyFile(defaultAvatarPath, newAvatarPath); err != nil {
            log.Printf("复制默认头像失败: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default avatar"})
            return
        }

        user.AvatarURL = relativePath

        // 更新用户的 AvatarURL
        if err := uc.DB.Save(&user).Error; err != nil {
            log.Println("更新用户头像失败:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user avatar"})
            return
        }
    } else {
        log.Println("查询数据库时发生错误:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // 生成 Access Token
    accessToken, err := uc.Utils.GenerateAccessToken(user.ID)
    if err != nil {
        log.Println("生成 Access Token 失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
        return
    }

    // 生成 Refresh Token
    refreshToken, err := uc.Utils.GenerateRefreshToken(user.ID)
    if err != nil {
        log.Println("生成 Refresh Token 失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
        return
    }

    // 存储 Refresh Token 到数据库
    newRefreshToken := models.RefreshToken{
        Token:     refreshToken,
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
        Revoked:   false,
    }
    if err := uc.DB.Create(&newRefreshToken).Error; err != nil {
        log.Println("存储 Refresh Token 失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
        return
    }

    // 返回成功响应，包括 Access Token 和 Refresh Token
    c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
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
    relativePath := fmt.Sprintf("avatars/%d_%d.jpg", user.ID, timestamp) // 文件的相对路径
    savePath := fmt.Sprintf("%s/%s", BaseUploadPath, relativePath)       // 文件的完整路径
    if err := c.SaveUploadedFile(file, savePath); err != nil {
        log.Println("文件保存失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // 删除旧头像文件（如果存在）
    if user.AvatarURL != "" && user.AvatarURL != "avatars/default.jpg" {
        oldPath := fmt.Sprintf("%s/%s", BaseUploadPath, user.AvatarURL)
        if err := os.Remove(oldPath); err != nil {
            log.Printf("无法删除旧头像文件: %s, 错误: %v\n", oldPath, err)
        }
    }

    // 更新用户头像路径
    user.AvatarURL = relativePath
    user.UpdatedAt = time.Now()
    if err := uc.DB.Save(&user).Error; err != nil {
        log.Println("更新用户头像失败:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
        return
    }

    log.Println("用户头像设置成功:", relativePath)
    c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully", "avatar_url": relativePath})
}

// RefreshTokenHandler 处理刷新 Access Token 的请求
func (uc *UserController) RefreshTokenHandler(c *gin.Context) {
    type RefreshTokenRequest struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    var req RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 验证 Refresh Token
    claims, err := uc.Utils.ValidateToken(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // 获取用户 ID
    userID, err := strconv.Atoi(claims.Subject)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
        return
    }

    // 查找 Refresh Token 在数据库中
    var storedRefreshToken models.RefreshToken
    if err := uc.DB.Where("token = ?", req.RefreshToken).First(&storedRefreshToken).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
        return
    }

    // 检查 Refresh Token 是否已过期或被撤销
    if storedRefreshToken.ExpiresAt.Before(time.Now()) || storedRefreshToken.Revoked {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is expired or revoked"})
        return
    }

    // 可选：验证用户存在
    var user models.User
    if err := uc.DB.First(&user, storedRefreshToken.UserID).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    // 生成新的 Access Token
    newAccessToken, err := uc.Utils.GenerateAccessToken(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
        return
    }

    // 可选：生成新的 Refresh Token，并撤销旧的
    newRefreshToken, err := uc.Utils.GenerateRefreshToken(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
        return
    }

    // 创建新的 Refresh Token 记录
    newRT := models.RefreshToken{
        Token:     newRefreshToken,
        UserID:    storedRefreshToken.UserID,
        ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
        Revoked:   false,
    }
    if err := uc.DB.Create(&newRT).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new refresh token"})
        return
    }

    // 撤销旧的 Refresh Token
    storedRefreshToken.Revoked = true
    if err := uc.DB.Save(&storedRefreshToken).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke old refresh token"})
        return
    }

    // 返回新的 Access Token 和 Refresh Token
    c.JSON(http.StatusOK, gin.H{
        "access_token":  newAccessToken,
        "refresh_token": newRefreshToken,
    })
}

// LogoutHandler 处理用户登出请求
func (uc *UserController) LogoutHandler(c *gin.Context) {
    type LogoutRequest struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    var req LogoutRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 验证 Refresh Token
    claims, err := uc.Utils.ValidateToken(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        return
    }

    // 获取用户 ID
    _, err = strconv.Atoi(claims.Subject)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
        return
    }

    // 查找 Refresh Token 在数据库中
    var storedRefreshToken models.RefreshToken
    if err := uc.DB.Where("token = ?", req.RefreshToken).First(&storedRefreshToken).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
        return
    }

    // 检查 Refresh Token 是否已过期或被撤销
    if storedRefreshToken.ExpiresAt.Before(time.Now()) || storedRefreshToken.Revoked {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is expired or revoked"})
        return
    }

    // 标记 Refresh Token 为已撤销
    storedRefreshToken.Revoked = true
    if err := uc.DB.Save(&storedRefreshToken).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke refresh token"})
        return
    }

    // 撤销所有关联的 Refresh Tokens（如果需要）
    if err := uc.DB.Model(&models.RefreshToken{}).
        Where("user_id = ? AND revoked = ?", storedRefreshToken.UserID, false).
        Update("revoked", true).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke all refresh tokens"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (uc *UserController) UserBasicDetails(c *gin.Context) {
    userID, _ := c.Get("user_id")

    var user models.User
    if err := uc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 计算已注册的天数
    location, _ := time.LoadLocation("Local")

    // 将用户创建时间和当前时间转换为服务器时区
    createdDate := user.CreatedAt.In(location).Truncate(24 * time.Hour)
    currentDate := time.Now().In(location).Truncate(24 * time.Hour)

    // 计算天数差
    registeredDays := int(currentDate.Sub(createdDate).Hours()/24)
    fmt.Println(registeredDays)

    c.JSON(http.StatusOK, gin.H{
        "id":             user.ID,
        "nickname":       user.Nickname,
        "avatar_url":     user.AvatarURL,
        "registered_days": registeredDays,
    })
}

// GetMyFavoritedNews 获取用户收藏的新闻 ID
func (uc *UserController) GetMyFavoritedNews(c *gin.Context) {
    // 1. 获取 user_id
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. 查询 user 并预加载 FavoritedNews
    var user models.User
    if err := uc.DB.Preload("FavoritedNews").First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
        }
        return
    }

    // 3. 收集 ID
    favoritedNewsIDs := make([]uint, 0, len(user.FavoritedNews))
    for _, news := range user.FavoritedNews {
        favoritedNewsIDs = append(favoritedNewsIDs, news.ID)
    }

    // 4. 返回 ID 列表
    c.JSON(http.StatusOK, gin.H{
        "news_ids": favoritedNewsIDs,
    })
}

// GetMyLikedNews 获取用户点赞的新闻 ID
func (uc *UserController) GetMyLikedNews(c *gin.Context) {
    // 1. 获取 user_id
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. 在数据库中找到此 user，并预加载 LikedNews
    var user models.User
    if err := uc.DB.Preload("LikedNews").First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
        }
        return
    }

    // 3. 收集 user.LikedNews 中所有 ID
    likedNewsIDs := make([]uint, 0, len(user.LikedNews))
    for _, news := range user.LikedNews {
        likedNewsIDs = append(likedNewsIDs, news.ID)
    }

    // 4. 返回 ID 列表
    c.JSON(http.StatusOK, gin.H{
        "news_ids": likedNewsIDs,
    })
}

// GetMyViewedNews 获取用户最近看过的新闻 ID
func (uc *UserController) GetMyViewedNews(c *gin.Context) {
    // 1. 获取 user_id
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. 查询 user 并预加载 ViewedNews
    var user models.User
    if err := uc.DB.Preload("ViewedNews").First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
        }
        return
    }

    // 3. 收集 ID
    viewedNewsIDs := make([]uint, 0, len(user.ViewedNews))
    for _, news := range user.ViewedNews {
        viewedNewsIDs = append(viewedNewsIDs, news.ID)
    }

    // 4. 返回 ID 列表
    c.JSON(http.StatusOK, gin.H{
        "news_ids": viewedNewsIDs,
    })
}

func (uc *UserController) GetUserProfile(c *gin.Context) {
    // 获取路径参数中的用户 ID
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // 查找用户信息
    var user models.User
    if err := uc.DB.First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
        return
    }

    // 查找用户创建的新闻
    var news []struct {
        ID    uint   `json:"id"`
        Title string `json:"title"`
    }
    if err := uc.DB.Model(&models.News{}).Select("id, title").Where("author_id = ?", user.ID).Find(&news).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's news"})
        return
    }

    // 返回用户信息和新闻列表
    c.JSON(http.StatusOK, gin.H{
        "nickname":  user.Nickname,
        "avatar_url": user.AvatarURL,
        "news":      news,
    })
}