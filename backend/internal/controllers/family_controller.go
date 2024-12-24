// internal/controllers/family_controller.go
package controllers

import (
	// "fmt"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

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
    if err := fc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    // TODO 改回来
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

    // 生成唯一的家庭 Token
    var familyToken string
    for {
        familyToken = utils.GenerateFamilyToken()
        var existingFamily models.Family
        if err := fc.DB.Where("token = ?", familyToken).First(&existingFamily).Error; err == gorm.ErrRecordNotFound {
            break // Token is unique
        }
    }

    // 创建新家庭
    family := models.Family{
        Name:        request.Name,
        Token:       familyToken,
        Admins:      []models.User{},
        Members:     []models.User{},
        WaitingList: []models.User{},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
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
    // 手动刷新 user 对象
    if err := fc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh user data"})
        return
    }

    // 返回创建的家庭信息
    c.JSON(http.StatusCreated, gin.H{
        "message": "Family created successfully",
        "family": gin.H{
            "id":        family.ID,
            "name":      family.Name,
            "family_id": family.Token,
        },
    })
}

// 获取今日日期
func getStartOfDay(t time.Time, loc *time.Location) time.Time {
    year, month, day := t.In(loc).Date()
    return time.Date(year, month, day, 0, 0, 0, 0, loc)
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
    if err := fc.DB.Preload("Family.Admins").Preload("Family.Members").Preload("Family.WaitingList").Preload("PendingFamily").First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 获取前端传来的时区
    timezone := c.Query("timezone")
    if timezone == "" {
        timezone = "Local" // 默认本地时区
    }
    loc, err := time.LoadLocation(timezone)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timezone"})
        return
    }

    now := time.Now().In(loc)
    today := getStartOfDay(now, loc)

    utcToday := today.UTC()

    // fmt.Println(utcToday)

    if user.PendingFamilyID != nil { // 用户在某个家庭的 waiting list 中
        c.JSON(http.StatusOK, gin.H{
            "status":    "waiting",
            "id":        user.PendingFamily.ID,
            "name":      user.PendingFamily.Name,
            "family_id": user.PendingFamily.Token,
        })
        return
    } else if user.Family != nil { // 用户已在某个家庭中
        // 准备管理员和成员的信息
        admins := make([]gin.H, len(user.Family.Admins))
        for i, admin := range user.Family.Admins {
            admins[i] = gin.H{
                "id":         admin.ID,
                "nickname":   admin.Nickname,
                "avatar_url": admin.AvatarURL,
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

        waiting_members := make([]gin.H, len(user.Family.WaitingList))
        for i, waiting_member := range user.Family.WaitingList {
            waiting_members[i] = gin.H{
                "id":         waiting_member.ID,
                "nickname":   waiting_member.Nickname,
                "avatar_url": waiting_member.AvatarURL,
            }
        }

        // 3. 遍历所有“真实成员”（Members 和 Admins 合并）
        realMembers := user.Family.Members
        realMembers = append(realMembers, user.Family.Admins...)

        memberDailyData := make([]gin.H, 0, len(realMembers))

        var (
            totalCarbonGoalSum   float64
            totalCarbonIntakeSum float64

            totalNutritionGoal = struct {
                Calories, Protein, Fat, Carbohydrates, Sodium float64
            }{}
            totalNutritionIntake = struct {
                Calories, Protein, Fat, Carbohydrates, Sodium float64
            }{}
        )

        for _, m := range realMembers {
            var carbonGoal models.CarbonGoal
            if err := fc.DB.Where("user_id = ? AND DATE(date) = DATE(?)", m.ID, utcToday).First(&carbonGoal).Error; err != nil {
                carbonGoal.Emission = 0
            }

            var carbonIntakes []models.CarbonIntake
            if err := fc.DB.Where("user_id = ? AND DATE(date) = DATE(?)", m.ID, utcToday).Find(&carbonIntakes).Error; err != nil {
                carbonIntakes = nil
            }
            var carbonIntakeSum float64
            for _, ci := range carbonIntakes {
                carbonIntakeSum += ci.Emission
            }

            var nutritionGoal models.NutritionGoal
            if err := fc.DB.Where("user_id = ? AND DATE(date) = DATE(?)", m.ID, utcToday).First(&nutritionGoal).Error; err != nil {
                nutritionGoal = models.NutritionGoal{}
            }

            var nutritionIntakes []models.NutritionIntake
            if err := fc.DB.Where("user_id = ? AND DATE(date) = DATE(?)", m.ID, utcToday).Find(&nutritionIntakes).Error; err != nil {
                nutritionIntakes = nil
            }
            var (
                niCals, niProtein, niFat, niCarbs, niSodium float64
            )
            for _, ni := range nutritionIntakes {
                niCals += ni.Calories
                niProtein += ni.Protein
                niFat += ni.Fat
                niCarbs += ni.Carbohydrates
                niSodium += ni.Sodium
            }

            singleMemberData := gin.H{
                "user_id": m.ID,
                "nickname": m.Nickname,
                "avatar_url": m.AvatarURL,

                "carbon_goal_emission": carbonGoal.Emission,
                "carbon_intake_sum":    carbonIntakeSum,

                "nutrition_goal": gin.H{
                    "calories":      nutritionGoal.Calories,
                    "protein":       nutritionGoal.Protein,
                    "fat":           nutritionGoal.Fat,
                    "carbohydrates": nutritionGoal.Carbohydrates,
                    "sodium":        nutritionGoal.Sodium,
                },
                "nutrition_intake_sum": gin.H{
                    "calories":      niCals,
                    "protein":       niProtein,
                    "fat":           niFat,
                    "carbohydrates": niCarbs,
                    "sodium":        niSodium,
                },
            }
            memberDailyData = append(memberDailyData, singleMemberData)

            totalCarbonGoalSum += carbonGoal.Emission
            totalCarbonIntakeSum += carbonIntakeSum

            totalNutritionGoal.Calories += nutritionGoal.Calories
            totalNutritionGoal.Protein += nutritionGoal.Protein
            totalNutritionGoal.Fat += nutritionGoal.Fat
            totalNutritionGoal.Carbohydrates += nutritionGoal.Carbohydrates
            totalNutritionGoal.Sodium += nutritionGoal.Sodium

            totalNutritionIntake.Calories += niCals
            totalNutritionIntake.Protein += niProtein
            totalNutritionIntake.Fat += niFat
            totalNutritionIntake.Carbohydrates += niCarbs
            totalNutritionIntake.Sodium += niSodium
        }

        familySumData := gin.H{
            "carbon_goal_sum":   totalCarbonGoalSum,
            "carbon_intake_sum": totalCarbonIntakeSum,
            "nutrition_goal_sum": gin.H{
                "calories":      totalNutritionGoal.Calories,
                "protein":       totalNutritionGoal.Protein,
                "fat":           totalNutritionGoal.Fat,
                "carbohydrates": totalNutritionGoal.Carbohydrates,
                "sodium":        totalNutritionGoal.Sodium,
            },
            "nutrition_intake_sum": gin.H{
                "calories":      totalNutritionIntake.Calories,
                "protein":       totalNutritionIntake.Protein,
                "fat":           totalNutritionIntake.Fat,
                "carbohydrates": totalNutritionIntake.Carbohydrates,
                "sodium":        totalNutritionIntake.Sodium,
            },
        }

        c.JSON(http.StatusOK, gin.H{
            "status":           "family",
            "id":               user.Family.ID,
            "name":             user.Family.Name,
            "family_id":        user.Family.Token,
            "member_count":     user.Family.MemberCount,
            "admins":           admins,
            "members":          members,
            "waiting_members":  waiting_members,
            "today_date":       today.Format("2006-01-02"),
            "member_daily_data": memberDailyData,
            "family_sums":       familySumData,
        })
        return
    } else {
        c.JSON(http.StatusOK, gin.H{
            "status": "empty",
        })
        return
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

    // 获取目标家庭
    var family models.Family
    if err := fc.DB.Preload("WaitingList").First(&family, adminUser.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
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
    if err := fc.DB.First(&user, userID).Error; err != nil {
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
    if *user.FamilyID != family.ID {
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
    if isTargetAdmin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The target user is already an admin of the family"})
        return
    }

    // 检查用户是否是家庭成员
    if *user.FamilyID != family.ID {
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
        // 从成员移除
        if err := tx.Model(&family).Association("Members").Delete(&user); err != nil {
            return err
        }
        // 添加到管理员
        if err := tx.Model(&family).Association("Admins").Append(&user); err != nil {
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
// Admin 退出家庭后顺延
func (fc *FamilyController) LeaveFamily(c *gin.Context) {
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

    if user.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    var family models.Family
    if err := fc.DB.Preload("Admins").Preload("Members").First(&family, *user.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 开始事务
    if err := fc.DB.Transaction(func(tx *gorm.DB) error {
        // 1. 删除用户提出的菜品（FamilyDish 表中 ProposerUserID 为该用户的记录）
        if err := tx.Where("proposer_user_id = ?", user.ID).Delete(&models.FamilyDish{}).Error; err != nil {
            return err
        }

        // 2. 更新菜品状态（删除孤立菜品）
        var remainingDishes []models.FamilyDish
        if err := tx.Where("family_id = ?", family.ID).Find(&remainingDishes).Error; err != nil {
            return err
        }
        for _, dish := range remainingDishes {
            var proposerCount int64
            if err := tx.Model(&models.FamilyDish{}).
                Where("family_id = ? AND dish_id = ?", dish.FamilyID, dish.DishID).
                Count(&proposerCount).Error; err != nil {
                return err
            }
            // 删除没有提议者的菜品
            if proposerCount == 0 {
                if err := tx.Delete(&dish).Error; err != nil {
                    return err
                }
            }
        }

        if err := tx.Model(&family).Association("Admins").Delete(&user); err != nil {
            return err
        }

        if err := tx.Model(&family).Association("Members").Delete(&user); err != nil {
            return err
        }

        family.MemberCount--
        if err := tx.Save(&family).Error; err != nil {
            return err
        }

        user.FamilyID = nil
        if err := tx.Save(&user).Error; err != nil {
            return err
        }

        adminCount := tx.Model(&family).Association("Admins").Count()

        if adminCount == 0 && family.MemberCount > 0 {
            var members []models.User
            if err := tx.Model(&family).Association("Members").Find(&members); err != nil {
                return err
            }

            if len(members) > 0 {
                // 随机指派用户为管理员
                newAdmin := members[rand.Intn(len(members))]

                if err := tx.Model(&family).Association("Members").Delete(&newAdmin); err != nil {
                    return err
                }

                if err := tx.Model(&family).Association("Admins").Append(&newAdmin); err != nil {
                    return err
                }
            }
        }

        if family.MemberCount == 0 {
            if err := tx.Delete(&family).Error; err != nil {
                return err
            }
        }

        return nil
    }); err != nil {
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave family"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Successfully left the family"})
}

// 踢出家庭
func (fc *FamilyController) DeleteFamilyMember(c *gin.Context) {
    adminUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var request struct {
        UserID uint `json:"user_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var adminUser models.User
    if err := fc.DB.First(&adminUser, adminUserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Admin user not found"})
        return
    }

    if adminUser.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    var family models.Family
    if err := fc.DB.Preload("Admins").Preload("Members").First(&family, *adminUser.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

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

    if user.FamilyID == nil || *user.FamilyID != family.ID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The user is not in your family"})
        return
    }

    if user.ID == adminUserID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot remove yourself"})
        return
    }

    // 移除用户并减少家庭成员计数
    if err := fc.DB.Transaction(func(tx *gorm.DB) error {
        // 删除用户提出的菜品
        if err := tx.Where("proposer_user_id = ?", user.ID).Delete(&models.FamilyDish{}).Error; err != nil {
            return err
        }

        // 更新菜品列表，删除没有提议者的菜品
        var remainingDishes []models.FamilyDish
        if err := tx.Where("family_id = ?", family.ID).Find(&remainingDishes).Error; err != nil {
            return err
        }
        for _, dish := range remainingDishes {
            var proposerCount int64
            if err := tx.Model(&models.FamilyDish{}).
                Where("family_id = ? AND dish_id = ?", dish.FamilyID, dish.DishID).
                Count(&proposerCount).Error; err != nil {
                return err
            }
            if proposerCount == 0 {
                if err := tx.Delete(&dish).Error; err != nil {
                    return err
                }
            }
        }

        // 从家庭中移除用户
        if err := tx.Model(&family).Association("Admins").Delete(&user); err != nil {
            return err
        }
        if err := tx.Model(&family).Association("Members").Delete(&user); err != nil {
            return err
        }

        // 减少家庭成员计数
        family.MemberCount--
        if err := tx.Save(&family).Error; err != nil {
            return err
        }

        // 将用户的 FamilyID 设置为 nil
        user.FamilyID = nil
        return tx.Save(&user).Error
    }); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from family"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Successfully removed user from family"})
}

// 解散家庭
func (fc *FamilyController) BreakFamily(c *gin.Context) {
    adminUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var adminUser models.User
    if err := fc.DB.Preload("Family.Admins").Preload("Family.Members").First(&adminUser, adminUserID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if adminUser.Family == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    family := adminUser.Family

    // 检查用户是否是家庭管理员
    isAdmin := false
    for _, admin := range family.Admins {
        if admin.ID == adminUserID {
            isAdmin = true
            break
        }
    }

    if !isAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to dissolve this family"})
        return
    }

    if err := fc.DB.Transaction(func(tx *gorm.DB) error {
        // 删除家庭中所有菜品
        if err := tx.Where("family_id = ?", family.ID).Delete(&models.FamilyDish{}).Error; err != nil {
            return err
        }

        // 从所有用户中解除家庭关联
        if err := tx.Model(&models.User{}).Where("family_id = ?", family.ID).Update("family_id", nil).Error; err != nil {
            return err
        }

        // 删除家庭的所有关联：Admins、Members、WaitingList
        if err := tx.Model(&family).Association("Admins").Clear(); err != nil {
            return err
        }
        if err := tx.Model(&family).Association("Members").Clear(); err != nil {
            return err
        }
        if err := tx.Model(&family).Association("WaitingList").Clear(); err != nil {
            return err
        }

        // 删除家庭记录
        if err := tx.Delete(&family).Error; err != nil {
            return err
        }

        return nil
    }); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dissolve the family"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Family dissolved successfully",
        "family_id": family.ID,
    })
}

// AddDesiredDish 处理添加想吃的菜请求
func (fc *FamilyController) AddDesiredDish(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    type AddDesiredDishRequest struct {
        DishID        uint  `json:"dish_id" binding:"required"`
        LevelOfDesire *uint `json:"level_of_desire" binding:"required,oneof=0 1 2"`
    }

    var request AddDesiredDishRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var user models.User
    if err := fc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if user.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    var family models.Family
    if err := fc.DB.Preload("Dishes").First(&family, *user.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    // 检查菜品是否已被该用户提出
    var existingFamilyDish models.FamilyDish
    if err := fc.DB.Where("family_id = ? AND dish_id = ? AND proposer_user_id = ?", family.ID, request.DishID, user.ID).First(&existingFamilyDish).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have already desired this dish"})
        return
    }

    // 允许同一个菜被多个成员提出，因此不需要检查其他成员的记录
    // 创建 FamilyDish 记录
    familyDish := models.FamilyDish{
        FamilyID:       family.ID,
        DishID:         request.DishID,
        LevelOfDesire:  *request.LevelOfDesire, // 解引用指针
        ProposerUserID: user.ID,
    }

    if err := fc.DB.Create(&familyDish).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add desired dish"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Desired dish added successfully"})
}

// 获取所有想吃菜品，按想吃程度排序
func (fc *FamilyController) GetDesiredDishes(c *gin.Context) {
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

    if user.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    var family models.Family
    if err := fc.DB.Preload("Dishes.Proposer").First(&family, *user.FamilyID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve family"})
        return
    }

    type GetDesiredDishesResponse struct {
        DishID        uint   `json:"dish_id"`
        LevelOfDesire uint   `json:"level_of_desire"`
        ProposerUser  models.User   `json:"proposer_user"`
    }

    var response []GetDesiredDishesResponse
    for _, fd := range family.Dishes {
        response = append(response, GetDesiredDishesResponse{
            DishID:        fd.DishID,
            LevelOfDesire: fd.LevelOfDesire,
            ProposerUser:  fd.Proposer,
        })
    }

    // 按 LevelOfDesire 降序排序
    sort.Slice(response, func(i, j int) bool {
        return response[i].LevelOfDesire > response[j].LevelOfDesire
    })

    c.JSON(http.StatusOK, response)
}

// 删除想吃菜品
func (fc *FamilyController) DeleteDesiredDish(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    type DeleteFamilyDishRequest struct {
        DishID uint `json:"dish_id" binding:"required"`
    }

    var request DeleteFamilyDishRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var user models.User
    if err := fc.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if user.FamilyID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You are not part of any family"})
        return
    }

    var familyDish models.FamilyDish
    if err := fc.DB.Where("family_id = ? AND dish_id = ? AND proposer_user_id = ?", *user.FamilyID, request.DishID, user.ID).First(&familyDish).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Desired dish not found or not proposed by you"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete desired dish"})
        }
        return
    }

    if err := fc.DB.Delete(&familyDish).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete desired dish"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Desired dish deleted successfully"})
}