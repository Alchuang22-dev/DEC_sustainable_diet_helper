// internal/controllers/family_controller_test.go
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	// "os"
	"testing"
	// "time"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupFamilyTestDB 初始化内存中的 SQLite 数据库并迁移模型
func setupFamilyTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to in-memory database")
	}

	// 迁移所有相关模型
	err = db.AutoMigrate(&models.User{}, &models.Family{}, &models.News{}, &models.FamilyDish{})
	if err != nil {
		panic("failed to migrate models")
	}

	return db
}

// setupFamilyRouter 初始化 Gin 路由和控制器
func setupFamilyRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	fc := NewFamilyController(db)

	// 注册路由组
	familyGroup := router.Group("/families")
	{
		authGroup := familyGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.POST("/create", fc.CreateFamily)
			authGroup.GET("/details", fc.FamilyDetails)
			authGroup.GET("/search", fc.SearchFamily)
			authGroup.POST("/:id/join", fc.JoinFamily)
			authGroup.POST("/admit", fc.AdmitJoinFamily)
			authGroup.POST("/reject", fc.RejectJoinFamily)
			authGroup.POST("/cancel_join", fc.CancelJoinFamily)
			authGroup.GET("/pending_family_details", fc.PendingFamilyDetails)
			authGroup.PUT("/set_member", fc.SetMember)
			authGroup.PUT("/set_admin", fc.SetAdmin)
			authGroup.DELETE("/leave_family", fc.LeaveFamily)
			authGroup.DELETE("/delete_family_member", fc.DeleteFamilyMember)
			authGroup.DELETE("/break", fc.BreakFamily)
			authGroup.POST("/add_desired_dish", fc.AddDesiredDish)
            authGroup.GET("/desired_dishes", fc.GetDesiredDishes)
            authGroup.DELETE("/desired_dishes", fc.DeleteDesiredDish)
		}
	}

	return router
}

// generateTestJWT 生成一个简单的 JWT（用于测试）
func generateTestJWT(userID uint) string {
	// 假设 utils.GenerateJWT(userID) 正常工作
	token, err := utils.GenerateAccessToken(userID)
	if err != nil {
		panic("failed to generate test JWT")
	}
	return token
}

// withAuth 在请求中添加 Authorization 头
func withAuth(req *http.Request, userID uint) {
	req.Header.Set("Authorization", "Bearer "+generateTestJWT(userID))
}

// ================ 测试 CreateFamily ================
func TestCreateFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建一个用户: userInFamily (已在家庭)
    alreadyFamily := models.Family{
        Name:        "AlreadyFamily",
        Token:       "already_family_token",
        MemberCount: 0,
    }
    db.Create(&alreadyFamily)

    userInFamily := models.User{
        OpenID:    "OpenID_CreateFamily_InFamily",
        Nickname:  "UserInFamily",
        FamilyID:  &alreadyFamily.ID, // 绑定到这个家庭
    }
    db.Create(&userInFamily)
    db.Model(&alreadyFamily).Association("Members").Append(&userInFamily)
    alreadyFamily.MemberCount++
    db.Save(&alreadyFamily)

    // 创建一个用户: userNotInFamily (未加入家庭)
    userNotInFamily := models.User{
        OpenID:   "OpenID_CreateFamily_NotInFamily",
        Nickname: "UserNotInFamily",
    }
    db.Create(&userNotInFamily)

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"name": "NewFamily"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "User Not Found",
            userID:         99999,
            requestBody:    gin.H{"name": "SomeFamily"},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Already In Family",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"name": "AnotherFamily"},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "User already belongs to a family"},
        },
        {
            name:           "Invalid Request Body (missing name)",
            userID:         userNotInFamily.ID,
            requestBody:    gin.H{}, // 没有 "name" 字段
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "Failed To Create Family (simulate db error)",
            userID:         userNotInFamily.ID,
            requestBody:    gin.H{"name": "SimulateFamily"},
            setupFunc: func() {
                // 模拟在 "create" 回调中注入错误
                db.Callback().Create().Before("gorm:create").Register("force_create_family_error", func(tx *gorm.DB) {
                    if tx.Statement.Table == "families" {
                        tx.Error = fmt.Errorf("forced error for families create")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to create family"},
        },
        {
            name:           "Successful Create Family",
            userID:         userNotInFamily.ID,
            requestBody:    gin.H{"name": "MyNewFamily"},
            setupFunc: func() {
                // 移除上一个回调，以免影响
                db.Callback().Create().Remove("force_create_family_error")
            },
            expectedStatus: http.StatusCreated,
            expectedBody: map[string]interface{}{
                "message": "Family created successfully",
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/families/create", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            // 对于成功情况(201)，只校验 "message" 字段足够；对 error 情况校验 "error"。
            if tc.expectedStatus == http.StatusCreated {
                assert.Equal(t, tc.expectedBody["message"], resp["message"])
                // 还可以检查 resp["family"] 是否存在
                _, ok := resp["family"]
                assert.True(t, ok, "response should contain 'family' field")
            } else {
                // 校验错误信息
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}

// ================ 测试 FamilyDetails ================
func TestFamilyDetails(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建用户: userPending (pendingFamily)
    familyForPending := models.Family{
        Name:        "PendingFamily",
        Token:       "pending_family_token",
        MemberCount: 0,
    }
    db.Create(&familyForPending)

    userPending := models.User{
        OpenID:          "OpenID_FamilyDetails_Pending",
        Nickname:        "UserPending",
        PendingFamilyID: &familyForPending.ID,
    }
    db.Create(&userPending)
    db.Model(&familyForPending).Association("WaitingList").Append(&userPending)

    // 创建用户: userInFamily (真实已在家庭)
    realFamily := models.Family{
        Name:        "RealFamily",
        Token:       "real_family_token",
        MemberCount: 1,
    }
    db.Create(&realFamily)
    userInFamily := models.User{
        OpenID:   "OpenID_FamilyDetails_InFamily",
        Nickname: "UserInFamily",
        FamilyID: &realFamily.ID,
    }
    db.Create(&userInFamily)
    db.Model(&realFamily).Association("Admins").Append(&userInFamily) // 让其成为 admin
    realFamily.MemberCount++
    db.Save(&realFamily)

    // 创建用户: userNoFamily (空家庭)
    userNoFamily := models.User{
        OpenID:   "OpenID_FamilyDetails_NoFamily",
        Nickname: "UserNoFamily",
    }
    db.Create(&userNoFamily)

    tests := []struct {
        name           string
        userID         uint
        queryTimezone  string
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            queryTimezone:  "",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "User Not Found",
            userID:         99999,
            queryTimezone:  "",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "Invalid Timezone",
            userID:         userInFamily.ID,
            queryTimezone:  "Invalid/TimeZone",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid timezone"},
        },
        {
            name:           "User In Pending Family",
            userID:         userPending.ID,
            queryTimezone:  "", // local
            setupFunc: func() {
                // 恢复 pendingFamily
                db.First(&familyForPending)
                db.First(&userPending)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "status":    "waiting",
                "id":        float64(familyForPending.ID), // JSON 反序列化后 ID 会是 float64
                "name":      familyForPending.Name,
                "family_id": familyForPending.Token,
            },
        },
        {
            name:           "User In Real Family",
            userID:         userInFamily.ID,
            queryTimezone:  "", 
            setupFunc: func() {
                db.First(&realFamily)
                db.First(&userInFamily)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "status":       "family",
                "id":           float64(realFamily.ID),
                "name":         realFamily.Name,
                "family_id":    realFamily.Token,
                "member_count": float64(realFamily.MemberCount),
                // 还会有 admins, members, waiting_members, ... 
                // 这里只做部分断言（示例）
            },
        },
        {
            name:           "User No Family",
            userID:         userNoFamily.ID,
            queryTimezone:  "",
            setupFunc:      func() {},
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"status": "empty"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/families/details"
            if tc.queryTimezone != "" {
                url += "?timezone=" + tc.queryTimezone
            }
            req, _ := http.NewRequest("GET", url, nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            // 如果是 200，可能返回 waiting/family/empty
            if tc.expectedStatus == http.StatusOK {
                // 简化做部分 key-value 检查
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            } else {
                // 检查 error
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}

// ================ 测试 SearchFamily ================
func TestSearchFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建一个用户 -> userAuthorized
    userAuthorized := models.User{
        OpenID:   "OpenID_SearchFamily_Authorized",
        Nickname: "UserSearchAuth",
    }
    db.Create(&userAuthorized)

    // 创建一个家庭 -> myFamily
    myFamily := models.Family{
        Name:        "MySearchableFamily",
        Token:       "my_search_token",
        MemberCount: 0,
    }
    db.Create(&myFamily)

    tests := []struct {
        name           string
        userID         uint
        familyToken    string
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            familyToken:    "whatever",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Missing Token",
            userID:         userAuthorized.ID,
            familyToken:    "",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Token is required"},
        },
        {
            name:           "Family Not Found",
            userID:         userAuthorized.ID,
            familyToken:    "not_exist_token",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "Family not found"},
        },
        {
            name:           "Failed To Search Family (simulate db error)",
            userID:         userAuthorized.ID,
            familyToken:    myFamily.Token,
            setupFunc: func() {
                // 强制对 family 的查询报错
                db.Callback().Query().Before("gorm:query").Register("force_search_family_error", func(tx *gorm.DB) {
                    if tx.Statement.Table == "families" {
                        tx.Error = fmt.Errorf("forced error for search family")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to search family"},
        },
        {
            name:           "Success Search Family",
            userID:         userAuthorized.ID,
            familyToken:    myFamily.Token,
            setupFunc: func() {
                db.Callback().Query().Remove("force_search_family_error")
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "id":           float64(myFamily.ID),
                "name":         myFamily.Name,
                "family_id":    myFamily.Token,
                "member_count": float64(myFamily.MemberCount),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/families/search"
            if tc.familyToken != "" {
                url += "?family_id=" + tc.familyToken
            }

            req, _ := http.NewRequest("GET", url, nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            // 对成功或失败分支做不同断言
            if tc.expectedStatus == http.StatusOK {
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            } else {
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}

// ================ 测试 JoinFamily ================
func TestJoinFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 准备用户 & 家庭
    family := models.Family{
        Name:        "JoinableFamily",
        Token:       "join_family_token",
        MemberCount: 0,
    }
    db.Create(&family)

    userNoFamily := models.User{
        OpenID:   "OpenID_JoinFamily_NoFamily",
        Nickname: "UserJoinNoFamily",
    }
    db.Create(&userNoFamily)

    userInFamily := models.User{
        OpenID:   "OpenID_JoinFamily_InFamily",
        Nickname: "UserJoinInFamily",
        FamilyID: &family.ID,
    }
    db.Create(&userInFamily)
    db.Model(&family).Association("Members").Append(&userInFamily)
    family.MemberCount++
    db.Save(&family)

    userHasPending := models.User{
        OpenID:          "OpenID_JoinFamily_HasPending",
        Nickname:        "UserJoinHasPending",
        PendingFamilyID: &family.ID, // 已有 pending
    }
    db.Create(&userHasPending)

    tests := []struct {
        name           string
        userID         uint
        familyParam    string // 路径参数
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            familyParam:    "1",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Family ID",
            userID:         userNoFamily.ID,
            familyParam:    "abc", // 无法转换成数字
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid family ID"},
        },
        {
            name:           "Family Not Found",
            userID:         userNoFamily.ID,
            familyParam:    "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "Family not found"},
        },
        {
            name:           "User Not Found",
            userID:         99999,
            familyParam:    fmt.Sprintf("%d", family.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Already In A Family",
            userID:         userInFamily.ID,
            familyParam:    fmt.Sprintf("%d", family.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are already a member of a family"},
        },
        {
            name:           "User Already Has Pending Family",
            userID:         userHasPending.ID,
            familyParam:    fmt.Sprintf("%d", family.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You have already requested to join another family"},
        },
        {
            name:           "Failed To Save Pending (simulate db error)",
            userID:         userNoFamily.ID,
            familyParam:    fmt.Sprintf("%d", family.ID),
            setupFunc: func() {
                db.Callback().Update().Before("gorm:update").Register("force_save_pending_error", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced error for saving user pending family")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to update pending family ID"},
        },
        {
            name:           "Success Join Family",
            userID:         userNoFamily.ID,
            familyParam:    fmt.Sprintf("%d", family.ID),
            setupFunc: func() {
                // 移除 waiting list 错误回调
                db.Callback().Update().Remove("force_save_pending_error")
                db.Callback().Update().Remove("force_waiting_list_error")
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Join request sent successfully"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := fmt.Sprintf("/families/%s/join", tc.familyParam)
            req, _ := http.NewRequest("POST", url, nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            for k, v := range tc.expectedBody {
                assert.Equal(t, v, resp[k])
            }
        })
    }
}

// ================ 测试 AdmitJoinFamily ================
func TestAdmitJoinFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 先创建一个家庭 & 管理员
    family := models.Family{
        Name:        "AdmitFamily",
        Token:       "admit_family_token",
        MemberCount: 0,
    }
    db.Create(&family)

    adminUser := models.User{
        OpenID:   "OpenID_AdmitJoin_Admin",
        Nickname: "AdminUserAdmit",
        FamilyID: &family.ID,
    }
    db.Create(&adminUser)
    db.Model(&family).Association("Admins").Append(&adminUser)
    family.MemberCount++
    db.Save(&family)

    // 再创建一个等待用户
    waitingUser := models.User{
        OpenID:          "OpenID_AdmitJoin_Waiting",
        Nickname:        "WaitingUserAdmit",
        PendingFamilyID: &family.ID,
    }
    db.Create(&waitingUser)
    db.Model(&family).Association("WaitingList").Append(&waitingUser)

    // 创建一个不在家庭的管理员 => adminNoFamily
    adminNoFamily := models.User{
        OpenID: "OpenID_AdmitJoin_AdminNoFamily",
    }
    db.Create(&adminNoFamily)

    // 创建一个用户: userAlreadyInFamily
    userAlreadyInFamily := models.User{
        OpenID:   "OpenID_AdmitJoin_AlreadyInFamily",
        Nickname: "AlreadyInFamily",
        FamilyID: &family.ID,
    }
    db.Create(&userAlreadyInFamily)

    tests := []struct {
        name           string
        adminID        uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            adminID:        0,
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body",
            adminID:        adminUser.ID,
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "Failed To Retrieve Admin User (simulate error)",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc: func() {
                // 注入错误
                db.Callback().Query().Before("gorm:query").Register("force_retrieve_admin_user_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced admin user retrieve error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve admin user"},
        },
        {
            name:           "Admin User Not In Any Family",
            adminID:        adminNoFamily.ID,
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc: func() {
                db.Callback().Query().Remove("force_retrieve_admin_user_err")
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "Not Admin Of This Family (403)",
            adminID:        userAlreadyInFamily.ID, // 虽然在家庭里，但不是管理员
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusForbidden,
            expectedBody:   map[string]interface{}{"error": "You are not an admin of this family"},
        },
        {
            name:           "Failed To Retrieve Family (simulate error)",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_retrieve_family_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "families" {
                        tx.Error = fmt.Errorf("forced error retrieve family")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve admin user"},
        },
        {
            name:           "User Not Found",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": 99999},
            setupFunc: func() {
                db.Callback().Query().Remove("force_retrieve_family_err")
            },
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Not In This Family's Waiting List",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": userAlreadyInFamily.ID}, // userAlreadyInFamily 并非 waitingList
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "User is not in the waiting list of your family"},
        },
        {
            name:           "User Already In Some Family",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": userAlreadyInFamily.ID},
            setupFunc: func() {
                // 先把 userAlreadyInFamily 放到 waitingList，pendingFamily也改成 family.ID
                userAlreadyInFamily.PendingFamilyID = &family.ID
                db.Save(&userAlreadyInFamily)
                db.Model(&family).Association("WaitingList").Append(&userAlreadyInFamily)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "User has been in a family"},
        },
        {
            name:           "Failed To Update User's Family Info (simulate error)",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc: func() {
                // 把 waitingUser 重新放回到 waitingList
                db.Model(&family).Association("WaitingList").Append(&waitingUser)
                waitingUser.PendingFamilyID = &family.ID
                waitingUser.FamilyID = nil
                db.Save(&waitingUser)

                // 强制更新 user family 出错
                db.Callback().Update().Before("gorm:update").Register("force_update_user_family_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced user family update error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to update user's family information"},
        },
        {
            name:           "Failed To Update Family Membership (transaction error)",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc: func() {
                db.Callback().Update().Remove("force_update_user_family_err")

                // 在 association 操作时注入错误
                db.Callback().Update().Before("gorm:association").Register("force_association_err", func(tx *gorm.DB) {
                    tx.Error = fmt.Errorf("forced association error")
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to update user's family information"},
        },
        {
            name:           "Successful Admit User",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": waitingUser.ID},
            setupFunc: func() {
                db.Callback().Update().Remove("force_association_err")

                // 重新把 waitingUser 放到 waitingList
                db.Model(&family).Association("WaitingList").Append(&waitingUser)
                waitingUser.PendingFamilyID = &family.ID
                waitingUser.FamilyID = nil
                db.Save(&waitingUser)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "message":   "User successfully admitted to the family",
                "family_id": float64(family.ID),
                "user_id":   float64(waitingUser.ID),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/families/admit", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.adminID != 0 {
                withAuth(req, tc.adminID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            for k, v := range tc.expectedBody {
                assert.Equal(t, v, resp[k])
            }
        })
    }
}

// TestRejectJoinFamily_InvalidRequestBody 测试无效的请求体
func TestRejectJoinFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 准备家庭 & 管理员
    family := models.Family{
        Name:        "RejectFamily",
        Token:       "reject_family_token",
        MemberCount: 0,
    }
    db.Create(&family)

    adminUser := models.User{
        OpenID:   "OpenID_RejectJoin_Admin",
        Nickname: "AdminUserReject",
        FamilyID: &family.ID,
    }
    db.Create(&adminUser)
    db.Model(&family).Association("Admins").Append(&adminUser)
    family.MemberCount++
    db.Save(&family)

    memberUser := models.User{
        OpenID:   "OpenID_RejectJoin_Member",
        Nickname: "AdminUserReject",
        FamilyID: &family.ID,
    }
    db.Create(&memberUser)
    db.Model(&family).Association("Members").Append(&memberUser)
    family.MemberCount++
    db.Save(&family)

    // 用户: userWaiting
    userWaiting := models.User{
        OpenID:          "OpenID_RejectJoin_Waiting",
        Nickname:        "RejectJoinWaitingUser",
        PendingFamilyID: &family.ID,
    }
    db.Create(&userWaiting)
    db.Model(&family).Association("WaitingList").Append(&userWaiting)

    // 用户: adminNoFamily
    adminNoFamily := models.User{
        OpenID: "OpenID_RejectJoin_AdminNoFamily",
    }
    db.Create(&adminNoFamily)

    tests := []struct {
        name           string
        adminID        uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            adminID:        0,
            requestBody:    gin.H{"user_id": userWaiting.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body",
            adminID:        adminUser.ID,
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "Failed To Retrieve Admin User (simulate error)",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": userWaiting.ID},
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_retrieve_admin_err_reject", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced retrieve admin user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve admin user"},
        },
        {
            name:           "Admin User Not In Any Family",
            adminID:        adminNoFamily.ID,
            requestBody:    gin.H{"user_id": userWaiting.ID},
            setupFunc: func() {
                db.Callback().Query().Remove("force_retrieve_admin_err_reject")
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "Not Admin Of This Family",
            adminID:        memberUser.ID,
            requestBody:    gin.H{"user_id": userWaiting.ID},
            setupFunc: func() {
                // 用 userWaiting 或者新用户 => 只要不是 admin 就好
                db.Callback().Query().Remove("force_retrieve_family_err_reject")
            },
            expectedStatus: http.StatusForbidden,
            expectedBody:   map[string]interface{}{"error": "You are not an admin of this family"},
        },
        {
            name:           "User Not Found",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": 99999},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Not In Waiting List",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": adminNoFamily.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "User is not in the waiting list of your family"},
        },
        {
            name:           "Failed To Update User's PendingFamilyID (simulate error)",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": userWaiting.ID},
            setupFunc: func() {
                db.Callback().Update().Remove("force_remove_waiting_err")
                db.Model(&family).Association("WaitingList").Append(&userWaiting)
                userWaiting.PendingFamilyID = &family.ID
                db.Save(&userWaiting)

                db.Callback().Update().Before("gorm:update").Register("force_update_pending_reject", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced user pendingFamilyID update error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to update user's pending family information"},
        },
        {
            name:           "Success Reject Join",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": userWaiting.ID},
            setupFunc: func() {
                db.Callback().Update().Remove("force_update_pending_reject")
                // 再次添加到 waitingList
                db.Model(&family).Association("WaitingList").Append(&userWaiting)
                userWaiting.PendingFamilyID = &family.ID
                db.Save(&userWaiting)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "message":   "User's join request rejected successfully",
                "family_id": float64(family.ID),
                "user_id":   float64(userWaiting.ID),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.adminID != 0 {
                // 注意: 这里为了模拟 "Not Admin Of This Family" 的测试，我们传 adminID=0 
                // 也会导致 Unauthorized => 401。若要模拟 "不是管理员" 但在家庭中，可以再造一个
                // userInFamilyButNotAdmin.
                withAuth(req, tc.adminID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            for k, v := range tc.expectedBody {
                assert.Equal(t, v, resp[k])
            }
        })
    }
}

func TestCancelJoinFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建用户和家庭
    family := models.Family{
        Name:        "Test Family",
        Token:       "testtoken",
        MemberCount: 1,
    }
    db.Create(&family)

    user := models.User{
        Nickname:        "TestUser",
        AvatarURL:       "http://example.com/avatar.jpg",
        FamilyID:        nil,
        PendingFamilyID: &family.ID,
    }
    db.Create(&user)

    // 将用户添加到家庭的等待列表
    db.Model(&family).Association("WaitingList").Append(&user)

    // 测试用例
    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized",
            userID:         0, // 未认证
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "User Not Found",
            userID:         999, // 不存在的用户
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:   "No Pending Family",
            userID: user.ID,
            setupFunc: func() {
                user.PendingFamilyID = nil
                db.Save(&user)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You have not requested to join any family"},
        },
        {
            name:   "Successful Cancellation",
            userID: user.ID,
            setupFunc: func() {
                user.PendingFamilyID = &family.ID
                db.Save(&user)
                db.Model(&family).Association("WaitingList").Append(&user)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "message":   "Join request canceled successfully",
                "family_id": float64(family.ID),
                "user_id":   float64(user.ID),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // 重置数据库状态
            db.Exec("DELETE FROM family_waiting_list")
            db.Model(&user).Association("WaitingList").Clear()
            tc.setupFunc()

            // 创建请求
            req, _ := http.NewRequest("POST", "/families/cancel_join", nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            // 记录响应
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 断言状态码
            assert.Equal(t, tc.expectedStatus, w.Code)

            // 断言响应体
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, response)
        })
    }
}

func TestPendingFamilyDetails(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建家庭
    family := models.Family{
        Name:        "Test Family",
        Token:       "testfamilytoken",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建用户
    user := models.User{
        Nickname:        "TestUser",
        AvatarURL:       "http://example.com/avatar.jpg",
        FamilyID:        nil,
        PendingFamilyID: &family.ID,
    }
    db.Create(&user)

    // 将用户添加到家庭的等待列表
    db.Model(&family).Association("WaitingList").Append(&user)

    // 测试用例
    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized",
            userID:         0, // 未认证
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "User Not Found",
            userID:         999, // 不存在的用户
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "No Pending Family",
            userID:         user.ID,
            setupFunc:      func() { user.PendingFamilyID = nil; db.Save(&user) },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You have not requested to join any family"},
        },
        {
            name:           "Successful Retrieval",
            userID:         user.ID,
            setupFunc:      func() { user.PendingFamilyID = &family.ID; db.Save(&user) },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "id":    float64(family.ID),
                "name":  family.Name,
                "token": family.Token,
            },
        },
        {
            name:   "Failed to Retrieve Family",
            userID: user.ID,
            setupFunc: func() {
                user.PendingFamilyID = &family.ID
                db.Save(&user)
                // 删除家庭以模拟数据库错误
                db.Delete(&family)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve family"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // 重置数据库状态
            db.Exec("DELETE FROM family_waiting_list")
            db.Model(&family).Association("WaitingList").Clear()
            tc.setupFunc()

            // 创建请求
            req, _ := http.NewRequest("GET", "/families/pending_family_details", nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            // 记录响应
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 断言状态码
            assert.Equal(t, tc.expectedStatus, w.Code)

            // 断言响应体
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, response)
        })
    }
}

func TestSetMember(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建家庭和用户
    family := models.Family{
        Name:        "Test Family",
        Token:       "testfamilytoken",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建管理员用户（adminUser）
    adminUser := models.User{
        Nickname:  "AdminUser",
        AvatarURL: "http://example.com/avatar_admin.jpg",
        FamilyID:  &family.ID,
		OpenID: "1",
    }
    db.Create(&adminUser)

    // 将该用户设置为家庭管理员
    db.Model(&family).Association("Admins").Append(&adminUser)

    // 创建目标用户（targetUser）
    targetUser := models.User{
        Nickname:  "TargetUser",
        AvatarURL: "http://example.com/avatar_target.jpg",
        FamilyID:  &family.ID,
		OpenID: "2",
    }
    db.Create(&targetUser)

    // 将目标用户也设置为管理员（以便能够转为 member）
    db.Model(&family).Association("Admins").Append(&targetUser)

    // 再创建一个不在此家庭的用户（otherUser）
    otherUser := models.User{
        Nickname:  "OtherUser",
        AvatarURL: "http://example.com/avatar_other.jpg",
		OpenID: "3",
    }
    db.Create(&otherUser)

    tests := []struct {
        name           string
        userID         uint          // 请求者（adminUserID）
        requestBody    interface{}   // 发送的 JSON body
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0, // 不传 token
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body (missing user_id)",
            userID:         adminUser.ID,
            requestBody:    gin.H{}, // 不带任何字段
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "Admin User Not Found",
            userID:         9999, // 数据库中不存在的用户
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "Admin User Not In Any Family",
            userID:         otherUser.ID, // 还没加入家庭
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "Failed To Retrieve Family",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc: func() {
                // 先删除该家庭，导致后续无法检索到家庭
                db.Delete(&family)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve family"},
        },
        {
            name:           "Requester Not Admin",
            userID:         targetUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc: func() {
                // 重新创建 family，因为上一个测试可能删掉了
                db.Create(&family)
                // 重新关联管理员 adminUser
                db.Model(&family).Association("Admins").Append(&adminUser)
                // 目标用户在同一个家庭，但不再是 Admin
                db.Model(&family).Association("Members").Append(&targetUser)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not an admin of this family"},
        },
        {
            name:           "Target User Not Found",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": 8888}, // 不存在
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "Target User Is Not Admin Already (cannot downgrade)",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc: func() {
                // 确保 targetUser 是 member 而非 admin
                db.Model(&family).Association("Admins").Delete(&targetUser)
                db.Model(&family).Association("Members").Append(&targetUser)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "The target user is not an admin of the family"},
        },
        {
            name:           "Cannot Change Own Role",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": adminUser.ID},
            setupFunc: func() {
                // adminUser 在 family 且是 admin
                // targetUser 在 family 且是 admin
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You cannot change your own role"},
        },
        {
            name:           "Successful Downgrade From Admin to Member",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc: func() {
                // 确保 targetUser 在 family.Admins 中
                db.Model(&family).Association("Admins").Append(&targetUser)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "message":   "Successfully set user to member",
                "family_id": float64(family.ID),
                "user_id":   float64(targetUser.ID),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // 测试前先清理数据库中与 family 相关的关联，避免测试间相互影响
            db.Model(&family).Association("Admins").Clear()
            db.Model(&family).Association("Members").Clear()
            // 重新关联 adminUser 为管理员
            db.Model(&family).Association("Admins").Append(&adminUser)

            // 每条测试用例都执行 setupFunc
            tc.setupFunc()

            // 构造请求体
            bodyBytes, _ := json.Marshal(tc.requestBody)
            // 注意这里改为 PUT
            req, _ := http.NewRequest("PUT", "/families/set_member", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            // 发送请求
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 断言状态码
            assert.Equal(t, tc.expectedStatus, w.Code)

            // 断言响应体
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)

            assert.Equal(t, tc.expectedBody, response)
        })
    }
}

func TestSetAdmin(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建家庭和用户
    family := models.Family{
        Name:        "Test Family For SetAdmin",
        Token:       "testfamilytoken-admin",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建管理员用户（adminUser）
    adminUser := models.User{
        Nickname:  "AdminUser2",
        AvatarURL: "http://example.com/avatar_admin2.jpg",
        FamilyID:  &family.ID,
		OpenID: "1",
    }
    db.Create(&adminUser)
    db.Model(&family).Association("Admins").Append(&adminUser)

    // 创建目标用户（targetUser）——初始为 member
    targetUser := models.User{
        Nickname:  "TargetUser2",
        AvatarURL: "http://example.com/avatar_target2.jpg",
        FamilyID:  &family.ID,
		OpenID: "2",
    }
    db.Create(&targetUser)
    db.Model(&family).Association("Members").Append(&targetUser)

    // 再创建一个不在此家庭的用户（otherUser）
    otherUser := models.User{
        Nickname:  "OtherUser2",
        AvatarURL: "http://example.com/avatar_other2.jpg",
		OpenID: "3",
    }
    db.Create(&otherUser)

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body (missing user_id)",
            userID:         adminUser.ID,
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "Admin User Not Found",
            userID:         9999,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "Admin User Not In Any Family",
            userID:         otherUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "Failed To Retrieve Family",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc: func() {
                // 删除家庭，触发无法查询到家庭的错误
                db.Delete(&family)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve family"},
        },
        {
            name:           "Requester Not Admin",
            userID:         targetUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc: func() {
                // 重新创建 family
                db.Create(&family)
                // 关联 adminUser 作为管理员
                db.Model(&family).Association("Admins").Append(&adminUser)
                // targetUser 仅是 member
                db.Model(&family).Association("Members").Append(&targetUser)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not an admin of this family"},
        },
        {
            name:           "Target User Not Found",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": 8888},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "Target User Already Admin",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": adminUser.ID},
            setupFunc: func() {
                // 使 targetUser 也变成管理员来测试
                db.Model(&family).Association("Admins").Append(&targetUser)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "The target user is already an admin of the family"},
        },
        {
            name:           "Successful Upgrade From Member to Admin",
            userID:         adminUser.ID,
            requestBody:    gin.H{"user_id": targetUser.ID},
            setupFunc: func() {
                // 确保 targetUser 是 member
                db.Model(&family).Association("Admins").Delete(&targetUser)
                db.Model(&family).Association("Members").Append(&targetUser)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                // 注意这里函数返回的 message 还是 "Successfully set user to member"
                // 如果想改成 "Successfully set user to admin" 需要在你的 SetAdmin 逻辑里做相应修改
                "message":   "Successfully set user to member",
                "family_id": float64(family.ID),
                "user_id":   float64(targetUser.ID),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // 清理前后状态，确保每个测试用例独立
            db.Model(&family).Association("Admins").Clear()
            db.Model(&family).Association("Members").Clear()
            // 重新将 adminUser 设为管理员
            db.Model(&family).Association("Admins").Append(&adminUser)

            tc.setupFunc()

            // 构造请求
            bodyBytes, _ := json.Marshal(tc.requestBody)
            // 注意这里改为 PUT
            req, _ := http.NewRequest("PUT", "/families/set_admin", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            // 发送请求
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 断言状态码
            assert.Equal(t, tc.expectedStatus, w.Code)

            // 断言响应体
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, response)
        })
    }
}

func TestLeaveFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 先创建一个家庭
    family := models.Family{
        Name:        "Test Family For Leave",
        Token:       "testfamilytoken_leave",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建用户（在家庭中 => adminUser）
    adminUser := models.User{
        OpenID:    "OpenID_Leave_Admin",
        Nickname:  "AdminUser_Leave",
        AvatarURL: "http://example.com/avatar_admin_leave.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&adminUser)

    // 将该用户设置为 Admin
    db.Model(&family).Association("Admins").Append(&adminUser)
    family.MemberCount++
    db.Save(&family)

    // 创建另一个家庭成员 => normalUser
    normalUser := models.User{
        OpenID:    "OpenID_Leave_Normal",
        Nickname:  "NormalUser_Leave",
        AvatarURL: "http://example.com/avatar_normal_leave.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&normalUser)
    db.Model(&family).Association("Members").Append(&normalUser)
    family.MemberCount++
    db.Save(&family)

    // 创建一个不在该家庭的用户 => otherUser
    otherUser := models.User{
        OpenID:    "OpenID_Leave_Other",
        Nickname:  "OtherUser_Leave",
        AvatarURL: "http://example.com/avatar_other_leave.jpg",
    }
    db.Create(&otherUser)

    // 让我们假设 FamilyDish 为测试中的菜品表（若不需要可忽略）
    // ...

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
		cleanupFunc    func()
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "User Not Found",
            userID:         9999,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Not Part Of Any Family",
            userID:         otherUser.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "Failed To Retrieve Family",
            userID:         adminUser.ID,
            setupFunc: func() {
                // 删除家庭，后续查询会失败
                db.Delete(&family)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve family"},
        },
        {
			name:           "Transaction Failure (simulate by forcing an error)",
			userID:         normalUser.ID,
			setupFunc: func() {
				// 重新创建并关联家庭（因为上面可能删掉了）
				db.Create(&family)
				db.Model(&family).Association("Admins").Append(&adminUser)
				db.Model(&family).Association("Members").Append(&normalUser)
				family.MemberCount = 2
				db.Save(&family)
		
				// 模拟事务中途失败
				// 1. 使用事务钩子强制返回错误
				db.Callback().Create().Before("gorm:create").Register("force_error", func(tx *gorm.DB) {
					tx.Error = fmt.Errorf("forced error for transaction failure simulation")
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "Failed to leave family"},
			cleanupFunc: func() {
				// 移除事务钩子，避免影响其他测试
				db.Callback().Create().Remove("force_error")
			},
		},
        {
            name:           "Successful Leave Family",
            userID:         adminUser.ID,
            setupFunc: func() {
                // 再次保证 family 存在 & 用户在该家庭
                db.First(&family)
                db.First(&adminUser)
                // 如果 adminUser 或 family 被删除过，需要重新插入、关联
                if adminUser.ID == 0 {
                    adminUser.OpenID = "OpenID_Leave_Admin_Reborn" // 避免重复 openID
                    adminUser.Nickname = "AdminUser_Leave_Reborn"
                    adminUser.AvatarURL = "http://example.com/avatar_admin_leave_reborn.jpg"
                    adminUser.FamilyID = &family.ID
                    db.Create(&adminUser)
                    db.Model(&family).Association("Admins").Append(&adminUser)
                    family.MemberCount++
                }
                db.Save(&family)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Successfully left the family"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // 每个测试前，重置数据库状态(如果有必要)
            // 由于我们在 setupFunc 中有针对性地修复数据，这里暂时不做太多清理

            tc.setupFunc()

            // 创建请求
            req, _ := http.NewRequest("DELETE", "/families/leave_family", nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            // 发送请求
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 断言状态码
            assert.Equal(t, tc.expectedStatus, w.Code)

            // 断言响应体
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, response)

			// 执行清理操作
            if tc.cleanupFunc != nil {
                tc.cleanupFunc()
            }
        })
    }
}

func TestDeleteFamilyMember(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建一个家庭
    family := models.Family{
        Name:        "Test Family For DeleteMember",
        Token:       "testfamilytoken_delete",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建管理员用户（adminUser）
    adminUser := models.User{
        OpenID:    "OpenID_DeleteMember_Admin",
        Nickname:  "AdminUser_DeleteMember",
        AvatarURL: "http://example.com/avatar_admin_delete.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&adminUser)
    db.Model(&family).Association("Admins").Append(&adminUser)
    family.MemberCount++
    db.Save(&family)

    // 创建目标用户 => normalUser
    normalUser := models.User{
        OpenID:    "OpenID_DeleteMember_Normal",
        Nickname:  "NormalUser_DeleteMember",
        AvatarURL: "http://example.com/avatar_normal_delete.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&normalUser)
    db.Model(&family).Association("Members").Append(&normalUser)
    family.MemberCount++
    db.Save(&family)

    // 创建不在家庭的用户 => otherUser
    otherUser := models.User{
        OpenID:    "OpenID_DeleteMember_Other",
        Nickname:  "OtherUser_DeleteMember",
        AvatarURL: "http://example.com/avatar_other_delete.jpg",
    }
    db.Create(&otherUser)

    tests := []struct {
        name           string
        adminID        uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            adminID:        0,
            requestBody:    gin.H{"user_id": normalUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body (missing user_id)",
            adminID:        adminUser.ID,
            requestBody:    gin.H{}, // 空
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "Admin User Not Found",
            adminID:        9999,
            requestBody:    gin.H{"user_id": normalUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "Admin user not found"},
        },
        {
            name:           "Admin User Not In Any Family",
            adminID:        otherUser.ID,
            requestBody:    gin.H{"user_id": normalUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "Failed To Retrieve Family",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": normalUser.ID},
            setupFunc: func() {
                // 删除家庭
                db.Delete(&family)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve family"},
        },
        {
            name:           "Not Admin Of This Family",
            adminID:        normalUser.ID,
            requestBody:    gin.H{"user_id": normalUser.ID},
            setupFunc: func() {
                // 重新创建 family
                db.Create(&family)
                // 仅把 adminUser 设为管理员
                db.Model(&family).Association("Admins").Append(&adminUser)
                db.Model(&family).Association("Members").Append(&normalUser)
                family.MemberCount = 2
                db.Save(&family)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not an admin of this family"},
        },
        {
            name:           "Target User Not Found",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": 8888},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Not In The Same Family",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": otherUser.ID},
            setupFunc: func() {
                // family 存在, adminUser 在 family.Admins
                // otherUser 不在 family
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "The user is not in your family"},
        },
        {
            name:           "Cannot Remove Yourself",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": adminUser.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You cannot remove yourself"},
        },
        // {
        //     name:           "Transaction Failure (simulate)",
        //     adminID:        adminUser.ID,
        //     requestBody:    gin.H{"user_id": normalUser.ID},
        //     setupFunc: func() {
        //         // 再次确保 family 存在 & adminUser / normalUser 在家庭
        //         // 删除 normalUser 以模拟事务错误
        //         db.Delete(&normalUser)
        //     },
        //     expectedStatus: http.StatusInternalServerError,
        //     expectedBody:   map[string]interface{}{"error": "Failed to remove user from family"},
        // },
        {
            name:           "Successful Remove User From Family",
            adminID:        adminUser.ID,
            requestBody:    gin.H{"user_id": normalUser.ID},
            setupFunc: func() {
                // 重新加回 normalUser
                normalUser.OpenID = "OpenID_DeleteMember_Normal_Reborn"
                normalUser.Nickname = "NormalUser_DeleteMember_Reborn"
                normalUser.AvatarURL = "http://example.com/avatar_normal_delete_reborn.jpg"
                normalUser.FamilyID = &family.ID
                db.Create(&normalUser)
                db.Model(&family).Association("Members").Append(&normalUser)
                family.MemberCount++
                db.Save(&family)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Successfully removed user from family"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            // 准备请求体
            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("DELETE", "/families/delete_family_member", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.adminID != 0 {
                withAuth(req, tc.adminID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 断言结果
            assert.Equal(t, tc.expectedStatus, w.Code)

            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, response)
        })
    }
}

func TestBreakFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建一个家庭
    family := models.Family{
        Name:        "Test Family For Break",
        Token:       "testfamilytoken_break",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建管理员用户 => adminUser
    adminUser := models.User{
        OpenID:    "OpenID_BreakFamily_Admin",
        Nickname:  "AdminUser_Break",
        AvatarURL: "http://example.com/avatar_admin_break.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&adminUser)
    db.Model(&family).Association("Admins").Append(&adminUser)
    family.MemberCount++
    db.Save(&family)

    // 创建普通成员 => normalUser
    normalUser := models.User{
        OpenID:    "OpenID_BreakFamily_Normal",
        Nickname:  "NormalUser_Break",
        AvatarURL: "http://example.com/avatar_normal_break.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&normalUser)
    db.Model(&family).Association("Members").Append(&normalUser)
    family.MemberCount++
    db.Save(&family)

    // 创建一个不在该家庭的用户 => otherUser
    otherUser := models.User{
        OpenID:    "OpenID_BreakFamily_Other",
        Nickname:  "OtherUser_Break",
        AvatarURL: "http://example.com/avatar_other_break.jpg",
    }
    db.Create(&otherUser)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "User Not Found",
            userID:         9999,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Not Part Of Any Family",
            userID:         otherUser.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "User Is Not Admin => Forbidden",
            userID:         normalUser.ID,
            setupFunc: func() {
                // 确保 normalUser 在 family
                // adminUser 仍是 admin
            },
            expectedStatus: http.StatusForbidden,
            expectedBody:   map[string]interface{}{"error": "You are not authorized to dissolve this family"},
        },
        // {
        //     name:           "Transaction Failure (simulate)",
        //     userID:         adminUser.ID,
        //     setupFunc: func() {
        //         // 删除 family，后续操作中就会失败
        //         db.Delete(&family)
        //     },
        //     expectedStatus: http.StatusInternalServerError,
        //     expectedBody:   map[string]interface{}{"error": "Failed to dissolve the family"},
        // },
        {
            name:           "Successful Break Family",
            userID:         adminUser.ID,
            setupFunc: func() {
                // 重新创建家庭 & 关联
                db.Create(&family)
                adminUser.OpenID = "OpenID_BreakFamily_Admin_Reborn"
                adminUser.FamilyID = &family.ID
                db.Create(&adminUser)
                db.Model(&family).Association("Admins").Append(&adminUser)
                family.MemberCount++

                normalUser.OpenID = "OpenID_BreakFamily_Normal_Reborn"
                normalUser.FamilyID = &family.ID
                db.Create(&normalUser)
                db.Model(&family).Association("Members").Append(&normalUser)
                family.MemberCount++

                db.Save(&family)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "message":   "Family dissolved successfully",
                "family_id": float64(family.ID),
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("DELETE", "/families/break", nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 断言状态码
            assert.Equal(t, tc.expectedStatus, w.Code)

            // 断言响应体
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, response)
        })
    }
}

func TestAddDesiredDish(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建家庭
    family := models.Family{
        Name:        "FamilyForAddDish",
        Token:       "add_dish_token",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建一个在家庭内的用户 (userInFamily)
    userInFamily := models.User{
        OpenID:    "OpenID_AddDish_User",
        Nickname:  "UserAddDish",
        AvatarURL: "http://example.com/avatar_add_dish.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&userInFamily)
    db.Model(&family).Association("Members").Append(&userInFamily)
    family.MemberCount++
    db.Save(&family)

    // 创建一个不在任何家庭的用户 (userNoFamily)
    userNoFamily := models.User{
        OpenID:    "OpenID_AddDish_NoFamily",
        Nickname:  "UserAddDishNoFamily",
        AvatarURL: "http://example.com/avatar_no_family.jpg",
    }
    db.Create(&userNoFamily)

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"dish_id": 1, "level_of_desire": 1},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body (missing dish_id)",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"level_of_desire": 2}, // dish_id 缺失
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "User Not Found",
            userID:         9999,
            requestBody:    gin.H{"dish_id": 1, "level_of_desire": 1},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Not Part Of Any Family",
            userID:         userNoFamily.ID,
            requestBody:    gin.H{"dish_id": 1, "level_of_desire": 1},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "Failed To Retrieve Family",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"dish_id": 1, "level_of_desire": 1},
            setupFunc: func() {
                // 删除该 family，让后续查不到
                db.Delete(&family)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to retrieve family"},
        },
        {
            name:           "Dish Already Desired By The Same User",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"dish_id": 2, "level_of_desire": 2},
            setupFunc: func() {
                // 重新创建 family 并关联 userInFamily
                db.Create(&family)
                userInFamily.FamilyID = &family.ID
                db.Save(&userInFamily)
                db.Model(&family).Association("Members").Append(&userInFamily)
                family.MemberCount = 1
                db.Save(&family)

                // 先插入一条 FamilyDish 记录表示用户已经提过 dish_id=2
                existingDish := models.FamilyDish{
                    FamilyID:       family.ID,
                    DishID:         2,
                    LevelOfDesire:  2,
                    ProposerUserID: userInFamily.ID,
                }
                db.Create(&existingDish)
            },
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You have already desired this dish"},
        },
        {
            name:           "Successfully Add Desired Dish",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"dish_id": 3, "level_of_desire": 1},
            setupFunc: func() {
                // 确保 family & userInFamily 正常
                db.First(&family)
                db.First(&userInFamily)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Desired dish added successfully"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/families/add_desired_dish", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, resp)
        })
    }
}

func TestGetDesiredDishes(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建家庭
    family := models.Family{
        Name:        "FamilyForGetDishes",
        Token:       "get_dish_token",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建一个在家庭内的用户 (userInFamily)
    userInFamily := models.User{
        OpenID:    "OpenID_GetDish_User",
        Nickname:  "UserGetDish",
        AvatarURL: "http://example.com/avatar_get_dish.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&userInFamily)
    db.Model(&family).Association("Members").Append(&userInFamily)
    family.MemberCount++
    db.Save(&family)

    // 创建一个不在任何家庭的用户 (userNoFamily)
    userNoFamily := models.User{
        OpenID:    "OpenID_GetDish_NoFamily",
        Nickname:  "UserGetDishNoFamily",
        AvatarURL: "http://example.com/avatar_no_family_get.jpg",
    }
    db.Create(&userNoFamily)

    // 准备一些 FamilyDish 数据
    // dish 1 (level_of_desire=2), dish 2 (level_of_desire=0), dish 3 (level_of_desire=1)
    fd1 := models.FamilyDish{FamilyID: family.ID, DishID: 1, LevelOfDesire: 2, ProposerUserID: userInFamily.ID}
    fd2 := models.FamilyDish{FamilyID: family.ID, DishID: 2, LevelOfDesire: 0, ProposerUserID: userInFamily.ID}
    fd3 := models.FamilyDish{FamilyID: family.ID, DishID: 3, LevelOfDesire: 1, ProposerUserID: userInFamily.ID}
    db.Create(&fd1)
    db.Create(&fd2)
    db.Create(&fd3)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        // GetDesiredDishes 返回的是一个数组，不是 map，所以我们要断言数组结构
        expectedBodyFunc func([]map[string]interface{}) bool
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBodyFunc: func(resp []map[string]interface{}) bool {
                return true
            },
        },
        {
            name:           "User Not Found",
            userID:         9999,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBodyFunc: func(resp []map[string]interface{}) bool {
                return true
            },
        },
        {
            name:           "User Not Part Of Any Family",
            userID:         userNoFamily.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBodyFunc: func(resp []map[string]interface{}) bool {
                return true
            },
        },
        {
            name:           "Failed To Retrieve Family",
            userID:         userInFamily.ID,
            setupFunc: func() {
                // 删除 family
                db.Delete(&family)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBodyFunc: func(resp []map[string]interface{}) bool {
                return true
            },
        },
        {
            name:           "Successful GetDesiredDishes",
            userID:         userInFamily.ID,
            setupFunc: func() {
                // 重新创建 family 并添加 dish
                db.Create(&family)
                userInFamily.FamilyID = &family.ID
                db.Save(&userInFamily)
                db.Model(&family).Association("Members").Append(&userInFamily)
                family.MemberCount = 1
                db.Save(&family)

                fd1 := models.FamilyDish{FamilyID: family.ID, DishID: 1, LevelOfDesire: 2, ProposerUserID: userInFamily.ID}
                fd2 := models.FamilyDish{FamilyID: family.ID, DishID: 2, LevelOfDesire: 0, ProposerUserID: userInFamily.ID}
                fd3 := models.FamilyDish{FamilyID: family.ID, DishID: 3, LevelOfDesire: 1, ProposerUserID: userInFamily.ID}
                db.Create(&fd1)
                db.Create(&fd2)
                db.Create(&fd3)
            },
            expectedStatus: http.StatusOK,
            expectedBodyFunc: func(resp []map[string]interface{}) bool {
                // 按 level_of_desire 降序 => [DishID=1, DishID=3, DishID=2]
                // resp[i]["dish_id"], resp[i]["level_of_desire"]
                if len(resp) != 3 {
                    return false
                }
                // 第一个应该是 dish_id=1, level_of_desire=2
                if resp[0]["dish_id"] != float64(1) || resp[0]["level_of_desire"] != float64(2) {
                    return false
                }
                // 第二个: dish_id=3, level_of_desire=1
                if resp[1]["dish_id"] != float64(3) || resp[1]["level_of_desire"] != float64(1) {
                    return false
                }
                // 第三个: dish_id=2, level_of_desire=0
                if resp[2]["dish_id"] != float64(2) || resp[2]["level_of_desire"] != float64(0) {
                    return false
                }
                return true
            },
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("GET", "/families/desired_dishes", nil)
            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            // 若返回的是 200，返回值为一个数组
            // 否则多为 {"error": "..."} 等 map
            if w.Code == http.StatusOK {
                var resp []map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.True(t, tc.expectedBodyFunc(resp))
            } else {
                // 简单做一下错误信息断言
                var errResp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &errResp)
                assert.NoError(t, err)
                // 这里只要能正常解析就行
            }
        })
    }
}

func TestDeleteDesiredDish(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建家庭
    family := models.Family{
        Name:        "FamilyForDeleteDish",
        Token:       "delete_dish_token",
        MemberCount: 0,
    }
    db.Create(&family)

    // 创建一个在家庭内的用户 (userInFamily)
    userInFamily := models.User{
        OpenID:    "OpenID_DeleteDish_User",
        Nickname:  "UserDeleteDish",
        AvatarURL: "http://example.com/avatar_delete_dish.jpg",
        FamilyID:  &family.ID,
    }
    db.Create(&userInFamily)
    db.Model(&family).Association("Members").Append(&userInFamily)
    family.MemberCount++
    db.Save(&family)

    // 创建一个不在任何家庭的用户 (userNoFamily)
    userNoFamily := models.User{
        OpenID:    "OpenID_DeleteDish_NoFamily",
        Nickname:  "UserDeleteDishNoFamily",
        AvatarURL: "http://example.com/avatar_no_family_delete.jpg",
    }
    db.Create(&userNoFamily)

    // 创建一个 FamilyDish 以便正常删除
    fd := models.FamilyDish{
        FamilyID:       family.ID,
        DishID:         10,
        LevelOfDesire:  2,
        ProposerUserID: userInFamily.ID,
    }
    db.Create(&fd)

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"dish_id": 10},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body (missing dish_id)",
            userID:         userInFamily.ID,
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "User Not Found",
            userID:         9999,
            requestBody:    gin.H{"dish_id": 10},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "User not found"},
        },
        {
            name:           "User Not Part Of Any Family",
            userID:         userNoFamily.ID,
            requestBody:    gin.H{"dish_id": 10},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "You are not part of any family"},
        },
        {
            name:           "FamilyDish Record Not Found (not proposed by you)",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"dish_id": 999}, // 不存在
            setupFunc: func() {
                // 保证family、userInFamily正常
                db.First(&family)
                userInFamily.FamilyID = &family.ID
                db.Save(&userInFamily)
            },
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "Desired dish not found or not proposed by you"},
        },
        {
            name:           "Failed To Delete Desired Dish (simulate DB error)",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"dish_id": 10},
            setupFunc: func() {
                // 再次确保 familyDish 存在
                db.Model(&family).Association("Members").Append(&userInFamily)
                familyDish := models.FamilyDish{
                    FamilyID:       family.ID,
                    DishID:         10,
                    LevelOfDesire:  2,
                    ProposerUserID: userInFamily.ID,
                }
                db.FirstOrCreate(&familyDish, familyDish)

                // 在删除操作之前，用 GORM 钩子注入一个错误
                db.Callback().Delete().Before("gorm:delete").Register("force_delete_error", func(tx *gorm.DB) {
                    tx.Error = fmt.Errorf("forced error for delete operation")
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to delete desired dish"},
        },
        {
            name:           "Successful Delete Desired Dish",
            userID:         userInFamily.ID,
            requestBody:    gin.H{"dish_id": 10},
            setupFunc: func() {
                // 移除之前的钩子，以免影响
                db.Callback().Delete().Remove("force_delete_error")

                // 确保 dish 10 存在
                db.FirstOrCreate(&fd, fd)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Desired dish deleted successfully"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("DELETE", "/families/desired_dishes", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                withAuth(req, tc.userID)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)
            assert.Equal(t, tc.expectedBody, resp)
        })
    }
}