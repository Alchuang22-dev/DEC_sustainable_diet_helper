// controllers/user_controller.go
package controllers

import (
    "net/http"
    "strconv"
    "log"
    "fmt"
    "encoding/json"

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

func (uc *UserController) WeChatAuth(c *gin.Context) {
    log.Println("WeChatAuth 被调用")

    // 定义请求体结构
    var authRequest struct {
        Code      string `json:"code" binding:"required"`        // 微信登录凭证
        Nickname  string `json:"nickname"`                      // 可选用户昵称
        AvatarURL string `json:"avatar_url"`                    // 可选用户头像
    }

    // 绑定请求体
    if err := c.ShouldBindJSON(&authRequest); err != nil {
        log.Println("绑定JSON失败:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // 调用微信 API 获取 openid 和 session_key
    wxAPI := fmt.Sprintf(
        "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        "your_app_id", "your_app_secret", authRequest.Code,
    )

    resp, err := http.Get(wxAPI)
    if err != nil {
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse WeChat API response"})
        return
    }

    // 微信 API 错误
    if wxResponse.ErrCode != 0 {
        log.Println("微信API返回错误:", wxResponse.ErrMsg)
        c.JSON(http.StatusUnauthorized, gin.H{"error": wxResponse.ErrMsg})
        return
    }

    // 检查用户是否已存在
    var user models.User
    if err := uc.DB.Where("openid = ?", wxResponse.OpenID).First(&user).Error; err == nil {
        // 用户已存在，更新 SessionKey
        user.SessionKey = wxResponse.SessionKey
        if err := uc.DB.Save(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user session"})
            return
        }
    } else if err == gorm.ErrRecordNotFound {
        // 用户不存在，注册新用户
        user = models.User{
            OpenID:     wxResponse.OpenID,
            SessionKey: wxResponse.SessionKey,
            Nickname:   authRequest.Nickname,
            AvatarURL:  authRequest.AvatarURL,
        }
        if user.Nickname == "" {
            user.Nickname = utils.GenerateRandomNickname()
        }

        if err := uc.DB.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }
    } else {
        // 其他错误
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // 生成 JWT
    token, err := utils.GenerateJWT(user.ID)
    if err != nil {
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