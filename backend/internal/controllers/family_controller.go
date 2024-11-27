// internal/controllers/family_controller.go
package controllers

import (
	"net/http"
	"time"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "fmt"
    "strconv"
)

type FamilyController struct {
    DB *gorm.DB
}

func NewFamilyController(db *gorm.DB) *FamilyController {
    return &FamilyController{DB: db}
}

// 创建家庭
func (fc *FamilyController) CreateFamily(c *gin.Context) {
    // 从 JWT 中解析用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 检查用户是否已属于某个家庭
    var user models.User
    if err := fc.DB.Preload("Family").First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    if user.FamilyID != 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User already belongs to a family"})
        return
    }

    // 解析请求体
    var request struct {
        Name string `json:"name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 创建家庭 Token（唯一标识）
    familyToken := utils.GenerateFamilyToken()

    // 创建新家庭
    family := models.Family{
        Name:  request.Name,
        Token: familyToken,
        Admins: []models.User{},
        Members:    []models.User{},
        WaitingList: []models.User{},
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        MemberCount: 1,
    }

    // 插入新家庭记录
    if err := fc.DB.Create(&family).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create family"})
        return
    }

    // 将当前用户设置为家庭管理员
    if err := fc.DB.Model(&family).Association("Admins").Append(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user as admin"})
        return
    }

    // 将家庭 ID 绑定到用户
    user.FamilyID = family.ID
    if err := fc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate user with family"})
        return
    }

    // 返回创建的家庭信息
    c.JSON(http.StatusCreated, gin.H{
        "message": "Family created successfully",
        "family": gin.H{
            "id":    family.ID,
            "name":  family.Name,
            "token": family.Token,
        },
    })
}

// 查看自己的家庭的信息
func (fc *FamilyController) FamilyDetails(c *gin.Context) {
    // 从 JWT 中解析用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 查询用户信息并预加载家庭信息
    var user models.User
    if err := fc.DB.Preload("Family.Admins").Preload("Family.Members").First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 检查用户是否属于家庭
    if user.Family == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User does not belong to a family"})
        return
    }

    // 准备管理员和成员的信息
    admins := make([]gin.H, len(user.Family.Admins))
    for i, admin := range user.Family.Admins {
        admins[i] = gin.H{
            "nickname":   admin.Nickname,
            "avatar_url": admin.AvatarURL, // TODO 应该传递图片，暂时用avatar_url替代
        }
    }

    members := make([]gin.H, len(user.Family.Members))
    for i, member := range user.Family.Members {
        members[i] = gin.H{
            "nickname":   member.Nickname,
            "avatar_url": member.AvatarURL,
        }
    }

    // 返回家庭信息
    c.JSON(http.StatusOK, gin.H{
        "id":      user.Family.ID,
        "name":    user.Family.Name,
        "token":   user.Family.Token,
        "member_count": user.Family.MemberCount,
        "admins":  admins,
        "members": members,
    })
}

// 查看搜索家庭结果
func (fc *FamilyController) SearchFamily(c *gin.Context) {
    _, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取 token 参数
    token := c.Query("token")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
        return
    }

    // 查找对应的家庭
    var family models.Family
    if err := fc.DB.Where("token = ?", token).First(&family).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Family not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search family"})
        return
    }

    // 返回家庭信息
    c.JSON(http.StatusOK, gin.H{
        "id":          family.ID,
        "name":        family.Name,
        "token":       family.Token,
        "member_count": family.MemberCount,
    })
}

// 发送加入家庭请求
func (fc *FamilyController) JoinFamily(c *gin.Context) {
    // 从 JWT 中解析当前用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取家庭 ID
    familyID, err := strconv.Atoi(c.Param("id"))
    if err != nil || familyID <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid family ID"})
        return
    }

    // 检查家庭是否存在
    var family models.Family
    if err := fc.DB.First(&family, familyID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Family not found"})
        return
    }

    // 获取当前用户
    var user models.User
    if err := fc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 检查是否已属于某个家庭
    if user.FamilyID != 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are already a member of a family"})
        return
    }

    // 检查是否已经申请加入其他家庭
    if user.PendingFamilyID != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already requested to join another family"})
        return
    }

    // 更新用户的 PendingFamilyID 字段
    user.PendingFamilyID = &family.ID
    if err := fc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pending family ID"})
        return
    }

    // 将用户添加到家庭的等待列表
    if err := fc.DB.Model(&family).Association("WaitingList").Append(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to the waiting list"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Join request sent successfully",
    })
}

// 批准加入家庭
func (fc *FamilyController) AdmitJoinFamily(c *gin.Context) {
    // 从 JWT 中解析当前用户 ID
    adminUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取被批准用户的 ID
    var request struct {
        UserID uint `json:"user_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 获取当前用户的 FamilyID
    var adminUser models.User
    if err := fc.DB.Preload("Family.Admins").First(&adminUser, adminUserID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin user"})
        return
    }

    if adminUser.FamilyID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    // 检查是否为家庭管理员
    isAdmin := false
    for _, admin := range adminUser.Family.Admins {
        if admin.ID == adminUserID.(uint) {
            isAdmin = true
            break
        }
    }
    if !isAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "You are not an admin of this family"})
        return
    }

    // 获取目标家庭
    var family models.Family
    if err := fc.DB.Preload("WaitingList").Preload("Members").First(&family, adminUser.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 检查被批准用户是否在家庭的等待列表中
    var user models.User
    if err := fc.DB.First(&user, request.UserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if user.PendingFamilyID == nil || *user.PendingFamilyID != family.ID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User is not in the waiting list of your family"})
        return // TODO 错误处理
    }

    if user.FamilyID != 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User has been in a family"})
        return // TODO 错误处理
    }

    // 更新用户的 FamilyID 和 PendingFamilyID 字段
    user.FamilyID = family.ID
    user.PendingFamilyID = nil
    if err := fc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user's family information"})
        return
    }

    // 从等待列表中移除用户并添加到成员列表
    if err := fc.DB.Model(&family).Association("WaitingList").Delete(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from waiting list"})
        return
    }
    if err := fc.DB.Model(&family).Association("Members").Append(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to members"})
        return
    }

    // 更新家庭成员计数
    family.MemberCount++
    if err := fc.DB.Save(&family).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update family member count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "User successfully admitted to the family",
        "family_id": family.ID,
        "user_id":   user.ID,
    })
}

// 拒绝加入家庭


// 取消申请


// TODO
// 更改某个家庭成员权限


// 退出家庭


// 解散家庭