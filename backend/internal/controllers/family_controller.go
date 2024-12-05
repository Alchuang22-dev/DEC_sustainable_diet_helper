// internal/controllers/family_controller.go
package controllers

import (
	"net/http"
    "strconv"
    "time"
    // "fmt"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
    if user.FamilyID != nil {
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
    user.FamilyID = &family.ID
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
            "family_id": family.Token,
        },
    })
}

// 查看自己的家庭的信息, 如果自己不在家庭或在 waiting list 也要相应显示
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

    if user.PendingFamilyID != nil { // 用户在某个家庭的 waiting list 中
        c.JSON(http.StatusOK, gin.H{
            "status":  "waiting",
            "id":      user.PendingFamily.ID,
            "name":    user.PendingFamily.Name,
            "family_id":   user.PendingFamily.Token,
        })
    } else if user.Family != nil { // 用户已在某个家庭中
        // 准备管理员和成员的信息
        admins := make([]gin.H, len(user.Family.Admins))
        for i, admin := range user.Family.Admins {
            admins[i] = gin.H{
                "id":         admin.ID,
                "nickname":   admin.Nickname,
                "avatar_url": admin.AvatarURL, // TODO 应该传递图片，暂时用avatar_url替代
            }
        }

        members := make([]gin.H, len(user.Family.Members))
        for i, member := range user.Family.Members {
            members[i] = gin.H{
                "id":         member.ID,
                "nickname":   member.Nickname,
                "avatar_url": member.AvatarURL,
            }
        }

        // 返回家庭信息
        c.JSON(http.StatusOK, gin.H{
            "status":  "family",
            "id":      user.Family.ID,
            "name":    user.Family.Name,
            "family_id":   user.Family.Token,
            "member_count": user.Family.MemberCount,
            "admins":  admins,
            "members": members,
        })
    } else {
        c.JSON(http.StatusOK, gin.H{
            "status":  "empty",
        })
    }
}

// 查看搜索家庭结果
func (fc *FamilyController) SearchFamily(c *gin.Context) {
    _, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取 token 参数
    token := c.Query("family_id")
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
        "family_id":       family.Token,
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
    if user.FamilyID != nil {
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

    if adminUser.FamilyID == nil {
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

    if user.FamilyID != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User has been in a family"})
        return // TODO 错误处理
    }

    // 更新用户的 FamilyID 和 PendingFamilyID 字段
    user.FamilyID = &family.ID
    user.PendingFamilyID = nil
    if err := fc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user's family information"})
        return
    }

    // 从等待列表中移除用户并添加到成员列表（使用事务保证原子性）
    err := fc.DB.Transaction(func(tx *gorm.DB) error {
        // 从等待列表中移除用户
        if err := tx.Model(&family).Association("WaitingList").Delete(&user); err != nil {
            return err
        }

        // 将用户添加到成员列表
        if err := tx.Model(&family).Association("Members").Append(&user); err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update family membership"})
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
func (fc *FamilyController) RejectJoinFamily(c *gin.Context) {
    // 从 JWT 中解析当前用户 ID
    adminUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取被拒绝用户的 ID
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

    if adminUser.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    // 检查是否为家庭管理员
    isAdmin := false
    for _, admin := range adminUser.Family.Admins {
        if admin.ID == adminUserID {
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
    if err := fc.DB.Preload("WaitingList").First(&family, adminUser.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 检查被拒绝用户是否在家庭的等待列表中
    var user models.User
    if err := fc.DB.First(&user, request.UserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if user.PendingFamilyID == nil || *user.PendingFamilyID != family.ID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User is not in the waiting list of your family"})
        return
    }

    // 从等待列表中移除用户
    if err := fc.DB.Model(&family).Association("WaitingList").Delete(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from waiting list"})
        return
    }

    // 清除用户的 PendingFamilyID 字段
    user.PendingFamilyID = nil
    if err := fc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user's pending family information"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "User's join request rejected successfully",
        "family_id": family.ID,
        "user_id":   user.ID,
    })
}

// 取消申请
func (fc *FamilyController) CancelJoinFamily(c *gin.Context) {
    // 从 JWT 中解析当前用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取当前用户
    var user models.User
    if err := fc.DB.Preload("PendingFamily").First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 检查用户是否有待处理的家庭申请
    if user.PendingFamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not requested to join any family"})
        return
    }

    // 获取用户申请的家庭
    var family models.Family
    if err := fc.DB.Preload("WaitingList").First(&family, *user.PendingFamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 从等待列表中移除用户
    if err := fc.DB.Model(&family).Association("WaitingList").Delete(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from waiting list"})
        return
    }

    // 清除用户的 PendingFamilyID 字段
    user.PendingFamilyID = nil
    if err := fc.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user's pending family information"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Join request canceled successfully",
        "family_id": family.ID,
        "user_id":   user.ID,
    })
}

// 查看当前试图加入的家庭信息
func (fc *FamilyController) PendingFamilyDetails(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var user models.User
    if err := fc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if user.PendingFamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not requested to join any family"})
        return
    }

    var family models.Family
    if err := fc.DB.First(&family, user.PendingFamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":       family.ID,
        "name":     family.Name,
        "token":    family.Token,
    })
}

// TODO
// 更改某个家庭成员为 member
func (fc *FamilyController) SetMember(c *gin.Context) {
    adminUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var request struct {
        UserID  uint `json:"user_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var adminUser models.User
    if err := fc.DB.First(&adminUser, adminUserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if adminUser.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    var family models.Family
    if err := fc.DB.Preload("Admins").Preload("Members").First(&family, adminUser.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 检查用户是否是管理员
    isAdmin := false
    for _, admin := range family.Admins {
        if admin.ID == adminUserID {
            isAdmin = true
            break
        }
    }
    if !isAdmin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not an admin of this family"})
        return
    }

    var user models.User
    if err := fc.DB.First(&user, request.UserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 检查目标用户是否是管理员
    isTargetAdmin := false
    for _, admin := range family.Admins {
        if admin.ID == user.ID {
            isTargetAdmin = true
            break
        }
    }
    if !isTargetAdmin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The target user is not an admin of the family"})
        return
    }

    // 检查用户是否是家庭成员
    if user.FamilyID != &family.ID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The user is not in your family"})
        return
    }

    // 防止管理员对自身权限进行修改
    if user.ID == adminUserID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot change your own role"})
        return
    }

    // 使用事务操作，确保原子性
    err := fc.DB.Transaction(func(tx *gorm.DB) error {
        // 从管理员移除
        if err := tx.Model(&family).Association("Admins").Delete(&user); err != nil {
            return err
        }
        // 添加到成员
        if err := tx.Model(&family).Association("Members").Append(&user); err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Successfully set user to member",
        "family_id": family.ID,
        "user_id":   user.ID,
    })
}

// 更改某个家庭成员为 admin
func (fc *FamilyController) SetAdmin(c *gin.Context) {
    adminUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var request struct {
        UserID  uint `json:"user_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var adminUser models.User
    if err := fc.DB.First(&adminUser, adminUserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if adminUser.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    var family models.Family
    if err := fc.DB.Preload("Admins").Preload("Members").First(&family, adminUser.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 检查用户是否是管理员
    isAdmin := false
    for _, admin := range family.Admins {
        if admin.ID == adminUserID {
            isAdmin = true
            break
        }
    }
    if !isAdmin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not an admin of this family"})
        return
    }

    var user models.User
    if err := fc.DB.First(&user, request.UserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 检查目标用户是否是管理员
    isTargetAdmin := false
    for _, admin := range family.Admins {
        if admin.ID == user.ID {
            isTargetAdmin = true
            break
        }
    }
    if !isTargetAdmin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The target user is not an admin of the family"})
        return
    }

    // 检查用户是否是家庭成员
    if user.FamilyID != &family.ID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The user is not in your family"})
        return
    }

    // 防止管理员对自身权限进行修改
    if user.ID == adminUserID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot change your own role"})
        return
    }

    // 使用事务操作，确保原子性
    err := fc.DB.Transaction(func(tx *gorm.DB) error {
        // 从管理员移除
        if err := tx.Model(&family).Association("Admins").Delete(&user); err != nil {
            return err
        }
        // 添加到成员
        if err := tx.Model(&family).Association("Members").Append(&user); err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Successfully set user to member",
        "family_id": family.ID,
        "user_id":   user.ID,
    })
}

// 退出家庭
func (fc* FamilyController) LeaveFamily(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var user models.User
    if err := fc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, "User not found")
        return
    }

    if user.FamilyID == nil {
        c.JSON(http.StatusBadRequest, "You are not part of any family")
        return
    }

    var family models.Family
    if err := fc.DB.Preload("Admins").Preload("Members").First(&family, &user.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 检查用户是不是唯一家庭成员。如果是，自动解散家庭
    if family.MemberCount > 1 {
        isAdmin := false
        isMember := false

        for _, admin := range family.Admins {
            if admin.ID == userID {
                isAdmin = true
                break
            }
        }
        for _, member := range family.Members {
            if member.ID == userID {
                isMember = true
                break
            }
        }

        if isAdmin && isMember {
            // TODO 处理错误
            c.JSON(http.StatusInternalServerError, gin.H{"error": "The user is currently both an admin and a member"})
            return
        }
        if !isAdmin && !isMember {
            // TODO 处理错误
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "The user is currently neighther an admin or a member"})
            return
        }
        if isAdmin {
            if err := fc.DB.Model(&family).Association("Admins").Delete(&user); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from admin list"})
                return
            }
        }
        if isAdmin {
            if err := fc.DB.Model(&family).Association("Members").Delete(&user); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from admin list"})
                return
            }
        }

        user.FamilyID = nil
        if err := fc.DB.Save(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to update user's family information"})
            return
        }
    } else if family.MemberCount == 1 {
        // 直接解散家庭
        
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "There is no member or admin in the family"})
        return
    }
}

// 踢出家庭
func (fc* FamilyController) DeleteFamilyMember(c *gin.Context) {
}

// 解散家庭
func (fc* FamilyController) BreakFamily(c *gin.Context) {
}