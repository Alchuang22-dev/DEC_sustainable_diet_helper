// controllers/user_controller.go
package controllers

import (
    "time"
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

// 发送验证码
func (uc *UserController) SendVerificationCode(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// 检查是否已存在未过期的验证码
	var existingVerification models.EmailVerification
	if err := uc.DB.Where("email = ?", request.Email).First(&existingVerification).Error; err == nil {
		if time.Now().Before(existingVerification.ExpiresAt) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Verification code already sent. Please wait."})
			return
		}
	}

	// 生成新的验证码
	code := utils.GenerateVerificationCode()
	verification := models.EmailVerification{
		Email:     request.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 验证码5分钟有效
	}

	// 保存到数据库
	if err := uc.DB.Save(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save verification code"})
		return
	}

	// 发送验证码邮件
	if err := utils.SendEmail(request.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
}

// 验证验证码并注册
func (uc *UserController) Register(c *gin.Context) {
	var request struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Code     string `json:"code" binding:"required,len=6"` // 验证码
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证验证码
	var verification models.EmailVerification
	if err := uc.DB.Where("email = ? AND code = ?", request.Email, request.Code).First(&verification).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification code"})
		return
	}
	if time.Now().After(verification.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Verification code expired"})
		return
	}

	// 删除已使用的验证码
	uc.DB.Delete(&verification)

	// 检查邮箱是否已注册
	var existingUser models.User
	if err := uc.DB.Where("email = ?", request.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// 如果未提供昵称，生成随机昵称
	nickname := request.Nickname
	if nickname == "" {
		nickname = utils.GenerateRandomNickname()
	}

	// 创建新用户
	user := models.User{
		Nickname: nickname,
		Email:    request.Email,
		Password: hashedPassword,
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 登录
func (uc *UserController) Login(c *gin.Context) {
    var loginRequest struct {
        Email    string `json:"email" binding:"required,email"` // 必须提供邮箱
        Password string `json:"password" binding:"required"`   // 必须提供密码
    }

    // 绑定请求体
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    // 验证用户是否存在
    var user models.User
    if err := uc.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
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

// 更新邮箱
func (uc *UserController) SetEmail(c *gin.Context) {
    log.Println("SetEmail 被调用")
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var request struct {
        Email string `json:"email" binding:"required,email"`
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

    // 检查邮箱唯一性
    var existingUser models.User
    if err := uc.DB.Where("email = ?", request.Email).First(&existingUser).Error; err == nil && existingUser.ID != uint(id) {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
        return
    }

    user.Email = request.Email

    if err := uc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully", "email": user.Email})
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

