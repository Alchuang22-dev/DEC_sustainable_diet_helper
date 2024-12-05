// controllers/family_controller_test.go
package controllers

import (
	"bytes"
	"encoding/json"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
)

// 设置测试数据库
func setupFamilyTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 自动迁移模型
	if err := db.AutoMigrate(&models.User{}, &models.Family{}); err != nil {
		panic("failed to migrate models")
	}
	return db
}

// 设置测试路由
func setupFamilyRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	familyController := NewFamilyController(db)

	familyGroup := router.Group("/families")
    {
        // 需要认证的路由
        authGroup := familyGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            authGroup.POST("/create", familyController.CreateFamily)
            authGroup.GET("/details", familyController.FamilyDetails)
            authGroup.GET("/search", familyController.SearchFamily)
            authGroup.POST("/:id/join", familyController.JoinFamily)
            authGroup.POST("/admit", familyController.AdmitJoinFamily)
        }
    }

	return router
}

// 测试创建家庭
func TestCreateFamily(t *testing.T) {
	db := setupFamilyTestDB()
	router := setupFamilyRouter(db)

	// 创建测试用户
	user := models.User{
		Nickname:  "TestUser",
		OpenID:    "test_openid",
		SessionKey: "test_session_key",
	}
	db.Create(&user)

	// 生成 JWT
	token, _ := utils.GenerateJWT(user.ID)

	// 构建请求体
	requestBody := map[string]string{
		"name": "TestFamily",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 创建请求
	req, _ := http.NewRequest("POST", "/families/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Family created successfully", response["message"])
	assert.NotEmpty(t, response["family"])

	familyData := response["family"].(map[string]interface{})
	assert.Equal(t, "TestFamily", familyData["name"])
	assert.NotEmpty(t, familyData["family_id"])

	// 检查数据库更新
	var createdFamily models.Family
	db.First(&createdFamily)
	assert.Equal(t, "TestFamily", createdFamily.Name)
	assert.Equal(t, 1, createdFamily.MemberCount)

	// 检查用户的 FamilyID 更新
	var updatedUser models.User
	db.First(&updatedUser, user.ID)
	assert.Equal(t, &createdFamily.ID, updatedUser.FamilyID)
}

// 测试查看家庭详情
func TestFamilyDetails(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建测试用户和家庭
    family := models.Family{Name: "TestFamily", Token: "test_token", MemberCount: 1}
    db.Create(&family)

    user := models.User{
        Nickname: "TestUser",
        OpenID:   "test_openid",
    }
    db.Create(&user)

    // 将用户设置为家庭管理员
    if err := db.Model(&family).Association("Admins").Append(&user); err != nil {
        t.Fatalf("failed to set user as admin: %v", err)
    }
    if err := db.Model(&family).Association("Members").Append(&user); err != nil {
        t.Fatalf("failed to add user to members: %v", err)
    }

    // 更新用户的 FamilyID
    user.FamilyID = &family.ID
    if err := db.Save(&user).Error; err != nil {
        t.Fatalf("failed to update user's family ID: %v", err)
    }

    // 生成 JWT
    token, _ := utils.GenerateJWT(user.ID)

    // 创建请求
    req, _ := http.NewRequest("GET", "/families/details", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    fmt.Println("Response: ", response)

    // 检查返回内容
    assert.Equal(t, float64(family.ID), response["id"])
    assert.Equal(t, "TestFamily", response["name"])
    assert.Equal(t, 1, int(response["member_count"].(float64)))

    admins := response["admins"].([]interface{})
    assert.Len(t, admins, 1)
    assert.Equal(t, "TestUser", admins[0].(map[string]interface{})["nickname"])

    members := response["members"].([]interface{})
    assert.Len(t, members, 1)
    assert.Equal(t, "TestUser", members[0].(map[string]interface{})["nickname"])
}

// 测试搜索家庭
func TestSearchFamily(t *testing.T) {
	db := setupFamilyTestDB()
	router := setupFamilyRouter(db)

	// 创建测试用户
	user := models.User{
		Nickname: "TestUser",
		OpenID:   "test_openid",
	}
	db.Create(&user)

	// 创建测试家庭
	family := models.Family{
		Name:        "TestFamily",
		Token:       "test_token",
		MemberCount: 1,
	}
	db.Create(&family)

	// 生成 JWT
	token, _ := utils.GenerateJWT(user.ID)

	// 创建请求并设置 Authorization 头
	req, _ := http.NewRequest("GET", "/families/search?family_id=test_token", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(family.ID), response["id"])
	assert.Equal(t, "TestFamily", response["name"])
	assert.Equal(t, "test_token", response["family_id"])
	assert.Equal(t, float64(1), response["member_count"])
}

func TestJoinAndAdmitFamily(t *testing.T) {
    db := setupFamilyTestDB()
    router := setupFamilyRouter(db)

    // 创建测试家庭和管理员用户
    admin := models.User{
        Nickname: "AdminUser",
        OpenID:   "admin_openid",
    }
    family := models.Family{
        Name:        "TestFamily",
        Token:       "test_token",
        MemberCount: 1,
    }
    db.Create(&admin)
    db.Create(&family)

    // 将管理员添加到家庭
    family.Admins = append(family.Admins, admin)
    admin.FamilyID = &family.ID
    db.Save(&family)
    db.Save(&admin)

    // 创建测试普通用户
    user := models.User{
        Nickname: "TestUser",
        OpenID:   "user_openid",
    }
    db.Create(&user)

    // 生成管理员和普通用户的 JWT
    adminToken, _ := utils.GenerateJWT(admin.ID)
    userToken, _ := utils.GenerateJWT(user.ID)

    // 测试普通用户发送加入家庭请求
    joinReq, _ := http.NewRequest("POST", fmt.Sprintf("/families/%d/join", family.ID), nil)
    joinReq.Header.Set("Authorization", "Bearer "+userToken)

    joinRes := httptest.NewRecorder()
    router.ServeHTTP(joinRes, joinReq)

    // 检查加入请求的响应
    assert.Equal(t, http.StatusOK, joinRes.Code)
    var joinResponse map[string]interface{}
    json.Unmarshal(joinRes.Body.Bytes(), &joinResponse)
    assert.Equal(t, "Join request sent successfully", joinResponse["message"])

    // 检查数据库中的等待列表
    var updatedFamily models.Family
    db.Preload("WaitingList").First(&updatedFamily, family.ID)
    assert.Equal(t, 1, len(updatedFamily.WaitingList))
    assert.Equal(t, user.ID, updatedFamily.WaitingList[0].ID)

    // 测试管理员批准加入请求
    admitBody := map[string]uint{"user_id": user.ID}
    admitBodyJSON, _ := json.Marshal(admitBody)
    admitReq, _ := http.NewRequest("POST", "/families/admit", bytes.NewBuffer(admitBodyJSON))
    admitReq.Header.Set("Authorization", "Bearer "+adminToken)
    admitReq.Header.Set("Content-Type", "application/json")

    admitRes := httptest.NewRecorder()
    router.ServeHTTP(admitRes, admitReq)

    // 检查批准请求的响应
    assert.Equal(t, http.StatusOK, admitRes.Code)
    var admitResponse map[string]interface{}
    json.Unmarshal(admitRes.Body.Bytes(), &admitResponse)
    assert.Equal(t, "User successfully admitted to the family", admitResponse["message"])
    assert.Equal(t, float64(family.ID), admitResponse["family_id"])
    assert.Equal(t, float64(user.ID), admitResponse["user_id"])

    // 检查数据库中的成员列表和用户状态
    db.Preload("Members").First(&updatedFamily, family.ID)
    assert.Equal(t, 1, len(updatedFamily.Members))
    assert.Equal(t, user.ID, updatedFamily.Members[0].ID)

    var updatedUser models.User
    db.First(&updatedUser, user.ID)
    assert.Equal(t, &family.ID, updatedUser.FamilyID)
    assert.Nil(t, updatedUser.PendingFamilyID)
}