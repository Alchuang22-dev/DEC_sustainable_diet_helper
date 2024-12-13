// // internal/controllers/family_controller_test.go
package controllers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	// "os"
// 	"testing"
// 	"time"

// 	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
// 	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
// 	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// // setupFamilyTestDB 初始化内存中的 SQLite 数据库并迁移模型
// func setupFamilyTestDB() *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect to in-memory database")
// 	}

// 	// 迁移所有相关模型
// 	err = db.AutoMigrate(&models.User{}, &models.Family{}, &models.News{})
// 	if err != nil {
// 		panic("failed to migrate models")
// 	}

// 	return db
// }

// // setupFamilyRouter 初始化 Gin 路由和控制器
// func setupFamilyRouter(db *gorm.DB) *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()
// 	fc := NewFamilyController(db)

// 	// 注册路由组
// 	familyGroup := router.Group("/families")
// 	{
// 		authGroup := familyGroup.Group("")
// 		authGroup.Use(middleware.AuthMiddleware())
// 		{
// 			authGroup.POST("/create", fc.CreateFamily)
// 			authGroup.GET("/details", fc.FamilyDetails)
// 			authGroup.GET("/search", fc.SearchFamily)
// 			authGroup.POST("/:id/join", fc.JoinFamily)
// 			authGroup.POST("/admit", fc.AdmitJoinFamily)
// 			authGroup.POST("/reject", fc.RejectJoinFamily)
// 			authGroup.POST("/cancel_join", fc.CancelJoinFamily)
// 			authGroup.GET("/pending_family_details", fc.PendingFamilyDetails)
// 			authGroup.POST("/set_member", fc.SetMember)
// 			authGroup.POST("/set_admin", fc.SetAdmin)
// 			authGroup.POST("/leave_family", fc.LeaveFamily)
// 			authGroup.POST("/delete_family_member", fc.DeleteFamilyMember)
// 			authGroup.POST("/break", fc.BreakFamily)
// 		}
// 	}

// 	return router
// }

// // generateTestJWT 生成一个简单的 JWT（用于测试）
// func generateTestJWT(userID uint) string {
// 	// 假设 utils.GenerateJWT(userID) 正常工作
// 	token, err := utils.GenerateAccessToken(userID)
// 	if err != nil {
// 		panic("failed to generate test JWT")
// 	}
// 	return token
// }

// // withAuth 在请求中添加 Authorization 头
// func withAuth(req *http.Request, userID uint) {
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(userID))
// }

// // ================ 测试 CreateFamily ================
// func TestCreateFamily_Unauthorized(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 无 token 请求
// 	req, _ := http.NewRequest("POST", "/families/create", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusUnauthorized, w.Code)
// 	assert.Contains(t, w.Body.String(), "Authorization header missing")
// }

// func TestCreateFamily_UserNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 用户 ID 999 不存在
// 	userID := uint(999)

// 	// 构建请求体
// 	bodyBytes, _ := json.Marshal(map[string]string{"name": "MyFamily"})
// 	req, _ := http.NewRequest("POST", "/families/create", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, userID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	assert.Contains(t, w.Body.String(), "User not found")
// }

// func TestCreateFamily_InvalidRequestBody(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建测试用户
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求体（缺少 name 字段）
// 	req, _ := http.NewRequest("POST", "/families/create", bytes.NewBuffer([]byte(`{}`)))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, user.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "Invalid request body")
// }

// func TestCreateFamily_Success(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建测试用户
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求体
// 	bodyBytes, _ := json.Marshal(map[string]string{"name": "MyFamily"})
// 	req, _ := http.NewRequest("POST", "/families/create", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, user.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusCreated, w.Code)
// 	var resp map[string]interface{}
// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "Family created successfully", resp["message"])

// 	// 检查家庭信息
// 	familyData, exists := resp["family"].(map[string]interface{})
// 	assert.True(t, exists)
// 	assert.Equal(t, "MyFamily", familyData["name"])
// 	assert.NotEmpty(t, familyData["family_id"])

// 	// 检查数据库中的家庭记录
// 	var family models.Family
// 	err = db.Where("token = ?", familyData["family_id"]).First(&family).Error
// 	assert.Nil(t, err)
// 	assert.Equal(t, "MyFamily", family.Name)
// 	assert.Equal(t, uint(1), family.MemberCount)

// 	// 检查用户的 FamilyID
// 	var updatedUser models.User
// 	db.First(&updatedUser, user.ID)
// 	assert.NotNil(t, updatedUser.FamilyID)
// 	assert.Equal(t, family.ID, *updatedUser.FamilyID)
// }

// // ================ 测试 FamilyDetails ================
// func TestFamilyDetails_Unauthorized(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 无 token 请求
// 	req, _ := http.NewRequest("GET", "/families/details", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusUnauthorized, w.Code)
// 	assert.Contains(t, w.Body.String(), "Missing Authorization header")
// }

// func TestFamilyDetails_UserNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	userID := uint(999) // 不存在的用户 ID

// 	// 构建请求
// 	req, _ := http.NewRequest("GET", "/families/details", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(userID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	assert.Contains(t, w.Body.String(), "User not found")
// }

// func TestFamilyDetails_NoFamily(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建测试用户，无家庭
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求
// 	req, _ := http.NewRequest("GET", "/families/details", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "empty")
// }

// func TestFamilyDetails_PendingFamily(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建家庭
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	db.Create(&family)

// 	// 创建用户，申请加入家庭
// 	user := models.User{
// 		Nickname:        "TestUser",
// 		PendingFamilyID: &family.ID,
// 	}
// 	db.Create(&user)

// 	// 构建请求
// 	req, _ := http.NewRequest("GET", "/families/details", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	var resp map[string]interface{}
// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "waiting", resp["status"])
// 	assert.Equal(t, float64(family.ID), resp["id"])
// 	assert.Equal(t, "TestFamily", resp["name"])
// 	assert.Equal(t, "testtoken", resp["family_id"])
// }

// func TestFamilyDetails_Success(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建家庭
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 2,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	db.Create(&family)

// 	// 创建管理员和成员
// 	admin := models.User{
// 		Nickname:    "AdminUser",
// 		FamilyID:    &family.ID,
//         OpenID: "1",
// 	}
// 	member := models.User{
// 		Nickname:    "MemberUser",
// 		FamilyID:    &family.ID,
//         OpenID: "2",
// 	}
// 	db.Create(&admin)
// 	db.Create(&member)

// 	// 设置家庭的 Admins 和 Members 关联
// 	db.Model(&family).Association("Admins").Append(&admin)
// 	db.Model(&family).Association("Members").Append(&member)

// 	// 构建请求
// 	req, _ := http.NewRequest("GET", "/families/details", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(admin.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	var resp map[string]interface{}
// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
// 	assert.Nil(t, err)

// 	assert.Equal(t, "family", resp["status"])
// 	assert.Equal(t, float64(family.ID), resp["id"])
// 	assert.Equal(t, "TestFamily", resp["name"])
// 	assert.Equal(t, "testtoken", resp["family_id"])
// 	assert.Equal(t, float64(family.MemberCount), resp["member_count"])

// 	// 检查 admins
// 	admins, exists := resp["admins"].([]interface{})
// 	assert.True(t, exists)
// 	assert.Equal(t, 1, len(admins))
// 	adminData := admins[0].(map[string]interface{})
// 	assert.Equal(t, float64(admin.ID), adminData["id"])
// 	assert.Equal(t, "AdminUser", adminData["nickname"])
// 	assert.Equal(t, "", adminData["avatar_url"]) // 根据实际头像设置

// 	// 检查 members
// 	members, exists := resp["members"].([]interface{})
// 	assert.True(t, exists)
// 	assert.Equal(t, 1, len(members))
// 	memberData := members[0].(map[string]interface{})
// 	assert.Equal(t, float64(member.ID), memberData["id"])
// 	assert.Equal(t, "MemberUser", memberData["nickname"])
// 	assert.Equal(t, "", memberData["avatar_url"]) // 根据实际头像设置
// }

// // ================ 测试 SearchFamily ================
// func TestSearchFamily_Unauthorized(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 无 token 请求
// 	req, _ := http.NewRequest("GET", "/families/search?family_id=testtoken", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusUnauthorized, w.Code)
// 	assert.Contains(t, w.Body.String(), "Missing Authorization header")
// }

// func TestSearchFamily_MissingToken(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 构建请求，无 family_id
// 	req, _ := http.NewRequest("GET", "/families/search", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(1))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "Token is required")
// }

// func TestSearchFamily_FamilyNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建测试用户
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求，查找不存在的家庭
// 	req, _ := http.NewRequest("GET", "/families/search?family_id=nonexistent", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	assert.Contains(t, w.Body.String(), "Family not found")
// }

// func TestSearchFamily_Success(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建家庭
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 2,
// 	}
// 	db.Create(&family)

// 	// 创建用户
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求
// 	req, _ := http.NewRequest("GET", "/families/search?family_id=testtoken", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	var resp map[string]interface{}
// 	err := json.Unmarshal(w.Body.Bytes(), &resp)
// 	assert.Nil(t, err)

// 	assert.Equal(t, float64(family.ID), resp["id"])
// 	assert.Equal(t, "TestFamily", resp["name"])
// 	assert.Equal(t, "testtoken", resp["family_id"])
// 	assert.Equal(t, float64(family.MemberCount), resp["member_count"])
// }

// // ================ 测试 JoinFamily ================
// func TestJoinFamily_Unauthorized(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 无 token 请求
// 	req, _ := http.NewRequest("POST", "/families/1/join", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusUnauthorized, w.Code)
// 	assert.Contains(t, w.Body.String(), "Missing Authorization header")
// }

// func TestJoinFamily_InvalidFamilyID(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建测试用户
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求，家庭 ID 无效
// 	req, _ := http.NewRequest("POST", "/families/invalid/join", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "Invalid family ID")
// }

// func TestJoinFamily_FamilyNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建测试用户
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求，查找不存在的家庭
// 	req, _ := http.NewRequest("POST", "/families/999/join", nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	assert.Contains(t, w.Body.String(), "Family not found")
// }

// func TestJoinFamily_UserAlreadyInFamily(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建家庭
// 	family := models.Family{
// 		Name:        "ExistingFamily",
// 		Token:       "existingtoken",
// 		MemberCount: 1,
// 	}
// 	db.Create(&family)

// 	// 创建用户，已属于某个家庭
// 	user := models.User{
// 		Nickname:  "TestUser",
// 		FamilyID:  &family.ID,
// 	}
// 	db.Create(&user)

// 	// 构建请求
// 	req, _ := http.NewRequest("POST", fmt.Sprintf("/families/%d/join", family.ID), nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "You are already a member of a family")
// }

// func TestJoinFamily_UserAlreadyRequested(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建家庭1
// 	family1 := models.Family{
// 		Name:        "FamilyOne",
// 		Token:       "familyonetoken",
// 		MemberCount: 1,
// 	}
// 	db.Create(&family1)

// 	// 创建家庭2
// 	family2 := models.Family{
// 		Name:        "FamilyTwo",
// 		Token:       "familytwotoken",
// 		MemberCount: 1,
// 	}
// 	db.Create(&family2)

// 	// 创建用户，已申请加入家庭1
// 	user := models.User{
// 		Nickname:        "TestUser",
// 		PendingFamilyID: &family1.ID,
// 	}
// 	db.Create(&user)

// 	// 将用户添加到家庭1的等待列表
// 	db.Model(&family1).Association("WaitingList").Append(&user)

// 	// 构建请求，尝试加入家庭2
// 	req, _ := http.NewRequest("POST", fmt.Sprintf("/families/%d/join", family2.ID), nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "You have already requested to join another family")
// }

// func TestJoinFamily_Success(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建家庭
// 	family := models.Family{
// 		Name:        "NewFamily",
// 		Token:       "newfamilytoken",
// 		MemberCount: 1,
// 	}
// 	db.Create(&family)

// 	// 创建用户，无家庭
// 	user := models.User{
// 		Nickname: "TestUser",
// 	}
// 	db.Create(&user)

// 	// 构建请求
// 	req, _ := http.NewRequest("POST", fmt.Sprintf("/families/%d/join", family.ID), nil)
// 	req.Header.Set("Authorization", "Bearer "+generateTestJWT(user.ID))

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "Join request sent successfully")

// 	// 检查用户的 PendingFamilyID
// 	var updatedUser models.User
// 	db.First(&updatedUser, user.ID)
// 	assert.NotNil(t, updatedUser.PendingFamilyID)
// 	assert.Equal(t, family.ID, *updatedUser.PendingFamilyID)

// 	// 检查家庭的 WaitingList
// 	var waitingUsers []models.User
// 	db.Model(&family).Association("WaitingList").Find(&waitingUsers)
// 	assert.Equal(t, 1, len(waitingUsers))
// 	assert.Equal(t, user.ID, waitingUsers[0].ID)
// }

// // ================ 测试 AdmitJoinFamily ================
// func TestAdmitJoinFamily_Unauthorized(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 无 token 请求
// 	req, _ := http.NewRequest("POST", "/families/admit", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusUnauthorized, w.Code)
// 	assert.Contains(t, w.Body.String(), "Missing Authorization header")
// }

// func TestAdmitJoinFamily_InvalidRequestBody(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
// 	}
// 	db.Create(&admin)

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "AdminFamily",
// 		Token:       "adminfamilytoken",
// 		MemberCount: 1,
// 	}
// 	db.Create(&family)
// 	db.Model(&family).Association("Admins").Append(&admin)

// 	// 构建请求体，缺少 user_id
// 	req, _ := http.NewRequest("POST", "/families/admit", bytes.NewBuffer([]byte(`{}`)))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "Invalid request body")
// }

// func TestAdmitJoinFamily_AdminUserNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	userID := uint(999) // 不存在的管理员用户

// 	// 构建请求体
// 	bodyBytes, _ := json.Marshal(map[string]uint{"user_id": 1})
// 	req, _ := http.NewRequest("POST", "/families/admit", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, userID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	assert.Contains(t, w.Body.String(), "Failed to retrieve admin user")
// }

// func TestAdmitJoinFamily_UserNotInWaitingList(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
//         OpenID: "1",
// 	}
// 	db.Create(&admin)

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "AdminFamily",
// 		Token:       "adminfamilytoken",
// 		MemberCount: 1,
// 	}
// 	db.Create(&family)
// 	db.Model(&family).Association("Admins").Append(&admin)

//     admin.FamilyID = &family.ID
//     db.Save(admin)

// 	// 创建被批准用户，不在等待列表中
// 	user := models.User{
// 		Nickname: "RegularUser",
//         OpenID: "2",
// 	}
// 	if err := db.Create(&user).Error; err != nil {
// 		t.Fatalf("failed to create user: %v", err)
// 	}

// 	// 构建请求体
// 	bodyBytes, _ := json.Marshal(map[string]uint{"user_id": user.ID})
// 	req, _ := http.NewRequest("POST", "/families/admit", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "User is not in the waiting list of your family")
// }

// func TestAdmitJoinFamily_Success(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
//         OpenID: "1",
// 	}
// 	db.Create(&admin)

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "AdminFamily",
// 		Token:       "adminfamilytoken",
// 		MemberCount: 1,
// 	}
// 	db.Create(&family)
// 	db.Model(&family).Association("Admins").Append(&admin)

//     admin.FamilyID = &family.ID
//     db.Save(admin)

// 	// 创建被批准用户，加入等待列表
// 	user := models.User{
// 		Nickname:        "RegularUser",
// 		PendingFamilyID: &family.ID,
//         OpenID: "2",
// 	}
// 	db.Create(&user)
// 	db.Model(&family).Association("WaitingList").Append(&user)

// 	// 构建请求体
// 	bodyBytes, _ := json.Marshal(map[string]uint{"user_id": user.ID})
// 	req, _ := http.NewRequest("POST", "/families/admit", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "User successfully admitted to the family")

// 	// 检查用户的 FamilyID 和 PendingFamilyID
// 	var updatedUser models.User
// 	db.First(&updatedUser, user.ID)
// 	assert.NotNil(t, updatedUser.FamilyID)
// 	assert.Nil(t, updatedUser.PendingFamilyID)

// 	// 检查家庭的成员列表
// 	var members []models.User
// 	db.Model(&family).Association("Members").Find(&members)
// 	assert.Equal(t, 1, len(members)) // 包含管理员和新成员

// 	// 检查家庭的成员计数
// 	var updatedFamily models.Family
// 	db.First(&updatedFamily, family.ID)
// 	assert.Equal(t, uint(2), updatedFamily.MemberCount)
// }

// // TestRejectJoinFamily_Unauthorized 测试未授权访问
// func TestRejectJoinFamily_Unauthorized(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 构建请求，不设置 Authorization header
// 	req, _ := http.NewRequest("POST", "/families/reject", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusUnauthorized, w.Code)
// 	assert.Contains(t, w.Body.String(), "Unauthorized")
// }

// // TestRejectJoinFamily_InvalidRequestBody 测试无效的请求体
// func TestRejectJoinFamily_InvalidRequestBody(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
// 	}
// 	if err := db.Create(&admin).Error; err != nil {
// 		t.Fatalf("failed to create admin user: %v", err)
// 	}

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}
// 	if err := db.Model(&family).Association("Admins").Append(&admin); err != nil {
// 		t.Fatalf("failed to associate admin with family: %v", err)
// 	}

// 	// 设置管理员的 FamilyID
// 	admin.FamilyID = &family.ID
// 	if err := db.Save(&admin).Error; err != nil {
// 		t.Fatalf("failed to update admin's FamilyID: %v", err)
// 	}

// 	// 构建请求体，缺少 user_id
// 	reqBody := map[string]interface{}{
// 		// "user_id" 缺失
// 	}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "Invalid request body")
// }

// // TestRejectJoinFamily_AdminUserNotFound 测试管理员用户不存在
// func TestRejectJoinFamily_AdminUserNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 使用不存在的管理员用户 ID
// 	nonExistentAdminID := uint(999)

// 	// 创建家庭和被拒绝用户
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}

// 	user := models.User{
// 		Nickname:        "WaitingUser",
// 		PendingFamilyID: &family.ID,
// 	}
// 	if err := db.Create(&user).Error; err != nil {
// 		t.Fatalf("failed to create waiting user: %v", err)
// 	}
// 	if err := db.Model(&family).Association("WaitingList").Append(&user); err != nil {
// 		t.Fatalf("failed to add user to waiting list: %v", err)
// 	}

// 	// 构建请求体
// 	reqBody := map[string]uint{"user_id": user.ID}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, nonExistentAdminID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	assert.Contains(t, w.Body.String(), "Failed to retrieve admin user")
// }

// // TestRejectJoinFamily_AdminNotPartOfAnyFamily 测试管理员不属于任何家庭
// func TestRejectJoinFamily_AdminNotPartOfAnyFamily(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建新管理员用户，不属于任何家庭
// 	newAdmin := models.User{
// 		Nickname: "NewAdminUser",
//         OpenID: "1",
// 	}
// 	if err := db.Create(&newAdmin).Error; err != nil {
// 		t.Fatalf("failed to create new admin user: %v", err)
// 	}

// 	// 创建家庭和被拒绝用户，并将被拒绝用户添加到等待列表
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}

// 	user := models.User{
// 		Nickname:        "WaitingUser",
// 		PendingFamilyID: &family.ID,
//         OpenID: "2",
// 	}
// 	if err := db.Create(&user).Error; err != nil {
// 		t.Fatalf("failed to create waiting user: %v", err)
// 	}
// 	if err := db.Model(&family).Association("WaitingList").Append(&user); err != nil {
// 		t.Fatalf("failed to add user to waiting list: %v", err)
// 	}

// 	// 构建请求体
// 	reqBody := map[string]uint{"user_id": user.ID}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, newAdmin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "You are not part of any family")
// }

// // TestRejectJoinFamily_AdminNotFamilyAdmin 测试管理员用户不是家庭管理员
// func TestRejectJoinFamily_AdminNotFamilyAdmin(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
//         OpenID: "1",
// 	}
// 	if err := db.Create(&admin).Error; err != nil {
// 		t.Fatalf("failed to create admin user: %v", err)
// 	}

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}
// 	if err := db.Model(&family).Association("Admins").Append(&admin); err != nil {
// 		t.Fatalf("failed to associate admin with family: %v", err)
// 	}

// 	// 设置管理员的 FamilyID
// 	admin.FamilyID = &family.ID
// 	if err := db.Save(&admin).Error; err != nil {
// 		t.Fatalf("failed to update admin's FamilyID: %v", err)
// 	}

// 	// 创建家庭成员但不是管理员
// 	member := models.User{
// 		Nickname: "FamilyMember",
// 		FamilyID: &family.ID,
//         OpenID: "2",
// 	}
// 	if err := db.Create(&member).Error; err != nil {
// 		t.Fatalf("failed to create family member: %v", err)
// 	}

// 	// 创建被拒绝用户，并添加到等待列表
// 	user := models.User{
// 		Nickname:        "WaitingUser",
// 		PendingFamilyID: &family.ID,
// 	}
// 	if err := db.Create(&user).Error; err != nil {
// 		t.Fatalf("failed to create waiting user: %v", err)
// 	}
// 	if err := db.Model(&family).Association("WaitingList").Append(&user); err != nil {
// 		t.Fatalf("failed to add user to waiting list: %v", err)
// 	}

// 	// 构建请求体
// 	reqBody := map[string]uint{"user_id": user.ID}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, member.ID) // 使用非管理员成员

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusForbidden, w.Code)
// 	assert.Contains(t, w.Body.String(), "You are not an admin of this family")
// }

// // TestRejectJoinFamily_FamilyNotFound 测试家庭不存在
// func TestRejectJoinFamily_FamilyNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
//         OpenID: "1",
// 	}
// 	if err := db.Create(&admin).Error; err != nil {
// 		t.Fatalf("failed to create admin user: %v", err)
// 	}

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}
// 	if err := db.Model(&family).Association("Admins").Append(&admin); err != nil {
// 		t.Fatalf("failed to associate admin with family: %v", err)
// 	}

// 	// 设置管理员的 FamilyID
// 	admin.FamilyID = &family.ID
// 	if err := db.Save(&admin).Error; err != nil {
// 		t.Fatalf("failed to update admin's FamilyID: %v", err)
// 	}

// 	// 创建被拒绝用户，并添加到等待列表
// 	user := models.User{
// 		Nickname:        "WaitingUser",
// 		PendingFamilyID: &family.ID,
//         OpenID: "2",
// 	}
// 	if err := db.Create(&user).Error; err != nil {
// 		t.Fatalf("failed to create waiting user: %v", err)
// 	}
// 	if err := db.Model(&family).Association("WaitingList").Append(&user); err != nil {
// 		t.Fatalf("failed to add user to waiting list: %v", err)
// 	}

// 	// 删除家庭以模拟家庭不存在
// 	if err := db.Delete(&family).Error; err != nil {
// 		t.Fatalf("failed to delete family: %v", err)
// 	}

// 	// 构建请求体
// 	reqBody := map[string]uint{"user_id": user.ID}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	assert.Contains(t, w.Body.String(), "Failed to retrieve family")
// }

// // TestRejectJoinFamily_UserNotFound 测试被拒绝用户不存在
// func TestRejectJoinFamily_UserNotFound(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
// 	}
// 	if err := db.Create(&admin).Error; err != nil {
// 		t.Fatalf("failed to create admin user: %v", err)
// 	}

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}
// 	if err := db.Model(&family).Association("Admins").Append(&admin); err != nil {
// 		t.Fatalf("failed to associate admin with family: %v", err)
// 	}

// 	// 设置管理员的 FamilyID
// 	admin.FamilyID = &family.ID
// 	if err := db.Save(&admin).Error; err != nil {
// 		t.Fatalf("failed to update admin's FamilyID: %v", err)
// 	}

// 	// 构建请求体，使用不存在的用户 ID
// 	reqBody := map[string]uint{"user_id": 999}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	assert.Contains(t, w.Body.String(), "User not found")
// }

// // TestRejectJoinFamily_UserNotInWaitingList 测试用户不在等待列表中
// func TestRejectJoinFamily_UserNotInWaitingList(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
//         OpenID: "1",
// 	}
// 	if err := db.Create(&admin).Error; err != nil {
// 		t.Fatalf("failed to create admin user: %v", err)
// 	}

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}
// 	if err := db.Model(&family).Association("Admins").Append(&admin); err != nil {
// 		t.Fatalf("failed to associate admin with family: %v", err)
// 	}

// 	// 设置管理员的 FamilyID
// 	admin.FamilyID = &family.ID
// 	if err := db.Save(&admin).Error; err != nil {
// 		t.Fatalf("failed to update admin's FamilyID: %v", err)
// 	}

// 	// 创建用户，不在等待列表中
// 	userNotWaiting := models.User{
// 		Nickname: "NonWaitingUser",
// 		FamilyID: &family.ID, // 已属于家庭，但不在等待列表
//         OpenID: "2",
// 	}
// 	if err := db.Create(&userNotWaiting).Error; err != nil {
// 		t.Fatalf("failed to create user not in waiting list: %v", err)
// 	}

// 	// 构建请求体
// 	reqBody := map[string]uint{"user_id": userNotWaiting.ID}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// 	assert.Contains(t, w.Body.String(), "User is not in the waiting list of your family")
// }

// // TestRejectJoinFamily_Success 测试成功拒绝用户加入家庭
// func TestRejectJoinFamily_Success(t *testing.T) {
// 	db := setupFamilyTestDB()
// 	router := setupFamilyRouter(db)

// 	// 创建管理员用户
// 	admin := models.User{
// 		Nickname: "AdminUser",
// 		OpenID: "1",
// 	}
// 	if err := db.Create(&admin).Error; err != nil {
// 		t.Fatalf("failed to create admin user: %v", err)
// 	}

// 	// 创建家庭并关联管理员
// 	family := models.Family{
// 		Name:        "TestFamily",
// 		Token:       "testtoken",
// 		MemberCount: 1,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	if err := db.Create(&family).Error; err != nil {
// 		t.Fatalf("failed to create family: %v", err)
// 	}
// 	if err := db.Model(&family).Association("Admins").Append(&admin); err != nil {
// 		t.Fatalf("failed to associate admin with family: %v", err)
// 	}

// 	// 设置管理员的 FamilyID
// 	admin.FamilyID = &family.ID
// 	if err := db.Save(&admin).Error; err != nil {
// 		t.Fatalf("failed to update admin's FamilyID: %v", err)
// 	}

// 	// 创建被拒绝用户，并添加到等待列表
// 	waitingUser := models.User{
// 		Nickname:        "WaitingUser",
// 		PendingFamilyID: &family.ID,
// 		OpenID: "2",
// 	}
// 	if err := db.Create(&waitingUser).Error; err != nil {
// 		t.Fatalf("failed to create waiting user: %v", err)
// 	}
// 	if err := db.Model(&family).Association("WaitingList").Append(&waitingUser); err != nil {
// 		t.Fatalf("failed to add user to waiting list: %v", err)
// 	}

// 	// 构建请求体
// 	reqBody := map[string]uint{"user_id": waitingUser.ID}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, _ := http.NewRequest("POST", "/families/reject", bytes.NewBuffer(bodyBytes))
// 	req.Header.Set("Content-Type", "application/json")
// 	withAuth(req, admin.ID)

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// 断言响应
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "User's join request rejected successfully")

// 	// 检查用户的 PendingFamilyID 是否被清除
// 	var updatedUser models.User
// 	if err := db.First(&updatedUser, waitingUser.ID).Error; err != nil {
// 		t.Fatalf("failed to retrieve updated user: %v", err)
// 	}
// 	assert.Nil(t, updatedUser.PendingFamilyID)

// 	// 检查用户是否从等待列表中移除
// 	var waitingUsers []models.User
// 	if err := db.Model(&family).Association("WaitingList").Find(&waitingUsers); err != nil {
// 		t.Fatalf("failed to retrieve waiting list: %v", err)
// 	}
// 	for _, wu := range waitingUsers {
// 		assert.NotEqual(t, waitingUser.ID, wu.ID)
// 	}
// }