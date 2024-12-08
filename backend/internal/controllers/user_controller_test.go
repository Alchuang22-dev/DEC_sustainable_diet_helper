// controllers/user_controller_test.go
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	// "strconv"
	"testing"
	"time"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupUserTestDB 初始化内存中的 SQLite 数据库并迁移模型
func setupUserTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(&models.User{}, &models.Family{}, &models.RefreshToken{}, &models.News{}); err != nil {
		panic("failed to migrate models")
	}
	return db
}

// setupUserRouter 初始化 Gin 路由和控制器
func setupUserRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userController := NewUserController(db)

	userGroup := router.Group("/users")
	{
		userGroup.POST("/auth", userController.WeChatAuth)
		userGroup.POST("/refresh", userController.RefreshTokenHandler)
		userGroup.POST("/logout", userController.LogoutHandler)

		authGroup := userGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.PUT("/set_nickname", userController.SetNickname)
			authGroup.PUT("/set_avatar", userController.SetAvatar)
			authGroup.GET("/basic_details", userController.UserBasicDetails) // 获取基本信息
		}
	}

	return router
}

// Helper function to generate a valid JWT for testing
func generateValidJWT(userID uint) string {
	token, err := utils.GenerateAccessToken(userID)
	if err != nil {
		panic("Failed to generate valid JWT for testing")
	}
	return token
}

// Helper function to generate a valid Refresh Token for testing
func generateValidRefreshToken(db *gorm.DB, userID uint) string {
	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		panic("Failed to generate valid Refresh Token for testing")
	}

	newRefreshToken := models.RefreshToken{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
		Revoked:   false,
	}
	if err := db.Create(&newRefreshToken).Error; err != nil {
		panic("Failed to store Refresh Token for testing")
	}
	return refreshToken
}

// TestWeChatAuth_Success 测试成功注册/登录新用户
func TestWeChatAuth_Success(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 模拟微信 API 响应
	wxResponse := `{
		"openid": "test_openid",
		"session_key": "test_session_key"
	}`

	// 创建 httptest 服务器模拟微信 API
	wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, wxResponse)
	}))
	defer wxServer.Close()

	// 设置环境变量
	os.Setenv("WECHAT_API_URL", wxServer.URL)
	os.Setenv("APP_ID", "test_app_id")
	os.Setenv("APP_SECRET", "test_app_secret")
	os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
	defer func() {
		os.Unsetenv("WECHAT_API_URL")
		os.Unsetenv("APP_ID")
		os.Unsetenv("APP_SECRET")
		os.RemoveAll("./test_uploads")
	}()

	// 创建默认头像文件
	defaultAvatarPath := filepath.Join("./test_uploads", "avatars")
	os.MkdirAll(defaultAvatarPath, 0755)
	os.WriteFile(filepath.Join(defaultAvatarPath, "default.jpg"), []byte("default avatar"), 0644)

	// 构建请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("code", "test_code")
	writer.WriteField("nickname", "TestUser")
	avatarFilePath := "./test_avatar.jpg"
	os.WriteFile(avatarFilePath, []byte("this is a test image"), 0644)
	defer os.Remove(avatarFilePath)
	part, err := writer.CreateFormFile("avatar", filepath.Base(avatarFilePath))
	assert.NoError(t, err)

	file, err := os.Open(avatarFilePath)
	assert.NoError(t, err)
	defer file.Close()

	_, err = io.Copy(part, file)
	assert.NoError(t, err)
	writer.Close()

	req, _ := http.NewRequest("POST", "/users/auth", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_token"])
	assert.NotEmpty(t, response["user"])

	userData := response["user"].(map[string]interface{})
	assert.Equal(t, "TestUser", userData["nickname"])
	assert.Contains(t, userData["avatar_url"], "avatars/")

	// 检查数据库中是否存在用户
	var user models.User
	result := db.Where("open_id = ?", "test_openid").First(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, "TestUser", user.Nickname)
	assert.NotEmpty(t, user.AvatarURL)

	// 检查 Refresh Token 是否存储
	var refreshToken models.RefreshToken
	err = db.Where("token = ?", response["refresh_token"]).First(&refreshToken).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, refreshToken.UserID)
	assert.False(t, refreshToken.Revoked)
}

// TestWeChatAuth_NoCode 测试缺少代码参数的情况
func TestWeChatAuth_NoCode(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	req, _ := http.NewRequest("POST", "/users/auth", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestWeChatAuth_WeChatAPIError 测试微信 API 调用失败
func TestWeChatAuth_WeChatAPIError(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 设置 WECHAT_API_URL 为一个不可访问的地址，模拟网络错误
	os.Setenv("WECHAT_API_URL", "http://localhost:9999") // 假设没有服务在此端口运行
	os.Setenv("APP_ID", "test_app_id")
	os.Setenv("APP_SECRET", "test_app_secret")
	os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
	defer func() {
		os.Unsetenv("WECHAT_API_URL")
		os.Unsetenv("APP_ID")
		os.Unsetenv("APP_SECRET")
		os.RemoveAll("./test_uploads")
	}()

	// 创建上传文件目录
	os.MkdirAll(filepath.Join("./test_uploads", "avatars"), 0755)

	// 构建请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("code", "test_code")
	writer.WriteField("nickname", "TestUser")
	writer.Close()

	// 创建请求
	req, _ := http.NewRequest("POST", "/users/auth", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestWeChatAuth_WeChatAPIWithErrCode 测试微信 API 返回错误码
func TestWeChatAuth_WeChatAPIWithErrCode(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 模拟微信 API 返回错误码
	wxResponse := `{
		"errcode": 40029,
		"errmsg": "invalid code"
	}`

	wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, wxResponse)
	}))
	defer wxServer.Close()

	os.Setenv("WECHAT_API_URL", wxServer.URL)
	os.Setenv("APP_ID", "test_app_id")
	os.Setenv("APP_SECRET", "test_app_secret")
	os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
	defer func() {
		os.Unsetenv("WECHAT_API_URL")
		os.Unsetenv("APP_ID")
		os.Unsetenv("APP_SECRET")
		os.RemoveAll("./test_uploads")
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("code", "invalid_code")
	writer.Close()

	req, _ := http.NewRequest("POST", "/users/auth", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestWeChatAuth_UpdateExistingUser_Nickname 测试已存在用户更新昵称
func TestWeChatAuth_UpdateExistingUser_Nickname(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 模拟微信 API 响应
    wxResponse := `{
        "openid": "test_openid",
        "session_key": "test_session_key"
    }`

    // 创建 httptest 服务器模拟微信 API
    wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, wxResponse)
    }))
    defer wxServer.Close()

    // 设置环境变量
    os.Setenv("WECHAT_API_URL", wxServer.URL)
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")
    os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
    defer func() {
        os.Unsetenv("WECHAT_API_URL")
        os.Unsetenv("APP_ID")
        os.Unsetenv("APP_SECRET")
        os.RemoveAll("./test_uploads")
    }()

    // 创建默认头像文件
    defaultAvatarPath := filepath.Join("./test_uploads", "avatars")
    os.MkdirAll(defaultAvatarPath, 0755)
    os.WriteFile(filepath.Join(defaultAvatarPath, "default.jpg"), []byte("default avatar"), 0644)

    // 创建已存在的用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "OldName",
        AvatarURL:  "avatars/old_avatar.jpg",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    db.Create(&user)

    // 保存旧头像文件
    oldAvatarPath := filepath.Join("./test_uploads", user.AvatarURL)
    os.WriteFile(oldAvatarPath, []byte("old avatar"), 0644)

    // 构建请求体，仅更新昵称
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("code", "test_code")
    writer.WriteField("nickname", "NewName")
    writer.Close()

    req, _ := http.NewRequest("POST", "/users/auth", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response["access_token"])
    assert.NotEmpty(t, response["refresh_token"])
    assert.NotEmpty(t, response["user"])

    userData := response["user"].(map[string]interface{})
    assert.Equal(t, "NewName", userData["nickname"])
    assert.Equal(t, user.AvatarURL, userData["avatar_url"]) // Avatar 未改变

    // 检查数据库中用户的昵称是否更新
    var updatedUser models.User
    result := db.Where("open_id = ?", "test_openid").First(&updatedUser)
    assert.Nil(t, result.Error)
    assert.Equal(t, "NewName", updatedUser.Nickname)
    assert.Equal(t, user.AvatarURL, updatedUser.AvatarURL)

    // 检查旧头像文件是否未被删除
    _, err = os.Stat(oldAvatarPath)
    assert.False(t, os.IsNotExist(err))
}

// TestWeChatAuth_UpdateExistingUser_Avatar 测试已存在用户更新头像
func TestWeChatAuth_UpdateExistingUser_Avatar(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 模拟微信 API 响应
    wxResponse := `{
        "openid": "test_openid",
        "session_key": "test_session_key"
    }`

    // 创建 httptest 服务器模拟微信 API
    wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, wxResponse)
    }))
    defer wxServer.Close()

    // 设置环境变量
    os.Setenv("WECHAT_API_URL", wxServer.URL)
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")
    os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
    defer func() {
        os.Unsetenv("WECHAT_API_URL")
        os.Unsetenv("APP_ID")
        os.Unsetenv("APP_SECRET")
        os.RemoveAll("./test_uploads")
    }()

    // 创建默认头像文件
    defaultAvatarPath := filepath.Join("./test_uploads", "avatars")
    os.MkdirAll(defaultAvatarPath, 0755)
    os.WriteFile(filepath.Join(defaultAvatarPath, "default.jpg"), []byte("default avatar"), 0644)

    // 创建已存在的用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "TestUser",
        AvatarURL:  "avatars/old_avatar.jpg",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    db.Create(&user)

    // 保存旧头像文件
    oldAvatarPath := filepath.Join("./test_uploads", user.AvatarURL)
    os.WriteFile(oldAvatarPath, []byte("old avatar"), 0644)

    // 创建新头像文件
    newAvatarPath := "./test_new_avatar.jpg"
    os.WriteFile(newAvatarPath, []byte("new avatar image"), 0644)
    defer os.Remove(newAvatarPath)

    // 构建请求体，仅更新头像
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("code", "test_code")
    // 不提供昵称，保持原昵称
    part, err := writer.CreateFormFile("avatar", filepath.Base(newAvatarPath))
    assert.NoError(t, err)

    file, err := os.Open(newAvatarPath)
    assert.NoError(t, err)
    defer file.Close()

    _, err = io.Copy(part, file)
    assert.NoError(t, err)
    writer.Close()

    req, _ := http.NewRequest("POST", "/users/auth", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)
    var response map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response["access_token"])
    assert.NotEmpty(t, response["refresh_token"])
    assert.NotEmpty(t, response["user"])

    userData := response["user"].(map[string]interface{})
    assert.Equal(t, "TestUser", userData["nickname"]) // Nickname 未改变
    assert.Contains(t, userData["avatar_url"], "avatars/")

    // 检查数据库中用户的 AvatarURL 是否更新
    var updatedUser models.User
    result := db.Where("open_id = ?", "test_openid").First(&updatedUser)
    assert.Nil(t, result.Error)
    assert.Equal(t, "TestUser", updatedUser.Nickname)
    assert.NotEqual(t, user.AvatarURL, updatedUser.AvatarURL)

    // 检查旧头像文件是否被删除
    _, err = os.Stat(oldAvatarPath)
    assert.True(t, os.IsNotExist(err))

    // 检查新头像文件是否存在
    newAvatarStoredPath := filepath.Join("./test_uploads", updatedUser.AvatarURL)
    _, err = os.Stat(newAvatarStoredPath)
    assert.False(t, os.IsNotExist(err))
}

// TestWeChatAuth_UpdateExistingUser_NicknameAndAvatar 测试已存在用户同时更新昵称和头像
func TestWeChatAuth_UpdateExistingUser_NicknameAndAvatar(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 模拟微信 API 响应
    wxResponse := `{
        "openid": "test_openid",
        "session_key": "test_session_key"
    }`

    // 创建 httptest 服务器模拟微信 API
    wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, wxResponse)
    }))
    defer wxServer.Close()

    // 设置环境变量
    os.Setenv("WECHAT_API_URL", wxServer.URL)
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")
    os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
    defer func() {
        os.Unsetenv("WECHAT_API_URL")
        os.Unsetenv("APP_ID")
        os.Unsetenv("APP_SECRET")
        os.RemoveAll("./test_uploads")
    }()

    // 创建默认头像文件
    defaultAvatarPath := filepath.Join("./test_uploads", "avatars")
    os.MkdirAll(defaultAvatarPath, 0755)
    os.WriteFile(filepath.Join(defaultAvatarPath, "default.jpg"), []byte("default avatar"), 0644)

    // 创建已存在的用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "OldName",
        AvatarURL:  "avatars/old_avatar.jpg",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    db.Create(&user)

    // 保存旧头像文件
    oldAvatarPath := filepath.Join("./test_uploads", user.AvatarURL)
    os.WriteFile(oldAvatarPath, []byte("old avatar"), 0644)

    // 创建新头像文件
    newAvatarPath := "./test_new_avatar.jpg"
    os.WriteFile(newAvatarPath, []byte("new avatar image"), 0644)
    defer os.Remove(newAvatarPath)

    // 构建请求体，更新昵称和头像
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("code", "test_code")
    writer.WriteField("nickname", "NewName")
    part, err := writer.CreateFormFile("avatar", filepath.Base(newAvatarPath))
    assert.NoError(t, err)

    file, err := os.Open(newAvatarPath)
    assert.NoError(t, err)
    defer file.Close()

    _, err = io.Copy(part, file)
    assert.NoError(t, err)
    writer.Close()

    req, _ := http.NewRequest("POST", "/users/auth", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)
    var response map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response["access_token"])
    assert.NotEmpty(t, response["refresh_token"])
    assert.NotEmpty(t, response["user"])

    userData := response["user"].(map[string]interface{})
    assert.Equal(t, "NewName", userData["nickname"])
    assert.Contains(t, userData["avatar_url"], "avatars/")

    // 检查数据库中用户的昵称和 AvatarURL 是否更新
    var updatedUser models.User
    result := db.Where("open_id = ?", "test_openid").First(&updatedUser)
    assert.Nil(t, result.Error)
    assert.Equal(t, "NewName", updatedUser.Nickname)
    assert.NotEqual(t, user.AvatarURL, updatedUser.AvatarURL)

    // 检查旧头像文件是否被删除
    _, err = os.Stat(oldAvatarPath)
    assert.True(t, os.IsNotExist(err))

    // 检查新头像文件是否存在
    newAvatarStoredPath := filepath.Join("./test_uploads", updatedUser.AvatarURL)
    _, err = os.Stat(newAvatarStoredPath)
    assert.False(t, os.IsNotExist(err))
}

// TestWeChatAuth_UpdateExistingUser_SaveDBError 测试更新用户时数据库保存失败
func TestWeChatAuth_UpdateExistingUser_SaveDBError(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 模拟微信 API 响应
    wxResponse := `{
        "openid": "test_openid",
        "session_key": "test_session_key"
    }`

    // 创建 httptest 服务器模拟微信 API
    wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, wxResponse)
    }))
    defer wxServer.Close()

    // 设置环境变量
    os.Setenv("WECHAT_API_URL", wxServer.URL)
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")
    os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
    defer func() {
        os.Unsetenv("WECHAT_API_URL")
        os.Unsetenv("APP_ID")
        os.Unsetenv("APP_SECRET")
        os.RemoveAll("./test_uploads")
    }()

    // 创建默认头像文件
    defaultAvatarPath := filepath.Join("./test_uploads", "avatars")
    os.MkdirAll(defaultAvatarPath, 0755)
    os.WriteFile(filepath.Join(defaultAvatarPath, "default.jpg"), []byte("default avatar"), 0644)

    // 创建已存在的用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "OldName",
        AvatarURL:  "avatars/old_avatar.jpg",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    db.Create(&user)

    // 模拟数据库保存错误，通过关闭数据库连接
    sqlDB, _ := db.DB()
    sqlDB.Close()

    // 构建请求体，仅更新昵称
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("code", "test_code")
    writer.WriteField("nickname", "NewName")
    writer.Close()

    req, _ := http.NewRequest("POST", "/users/auth", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestWeChatAuth_UpdateExistingUser_SaveAvatarError 测试上传头像时保存失败
func TestWeChatAuth_UpdateExistingUser_SaveAvatarError(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 模拟微信 API 响应
    wxResponse := `{
        "openid": "test_openid",
        "session_key": "test_session_key"
    }`

    // 创建 httptest 服务器模拟微信 API
    wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, wxResponse)
    }))
    defer wxServer.Close()

    // 设置环境变量
    os.Setenv("WECHAT_API_URL", wxServer.URL)
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")
    os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
    defer func() {
        os.Unsetenv("WECHAT_API_URL")
        os.Unsetenv("APP_ID")
        os.Unsetenv("APP_SECRET")
        os.RemoveAll("./test_uploads")
    }()

    // 创建默认头像文件
    defaultAvatarPath := filepath.Join("./test_uploads", "avatars")
    os.MkdirAll(defaultAvatarPath, 0755)
    os.WriteFile(filepath.Join(defaultAvatarPath, "default.jpg"), []byte("default avatar"), 0644)

    // 创建已存在的用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "OldName",
        AvatarURL:  "avatars/old_avatar.jpg",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    db.Create(&user)

    // 设置 BASE_UPLOAD_PATH 为不可写的路径，模拟保存文件失败
    os.Setenv("BASE_UPLOAD_PATH", "/invalid/path")
    defer os.Unsetenv("BASE_UPLOAD_PATH")

    // 构建 request body
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("code", "test_code")
    writer.WriteField("nickname", "NewName")
    avatarFilePath := "./test_avatar.jpg"
    os.WriteFile(avatarFilePath, []byte("this is a test image"), 0644)
    defer os.Remove(avatarFilePath)
    part, err := writer.CreateFormFile("avatar", filepath.Base(avatarFilePath))
    assert.NoError(t, err)

    file, err := os.Open(avatarFilePath)
    assert.NoError(t, err)
    defer file.Close()

    _, err = io.Copy(part, file)
    assert.NoError(t, err)
    writer.Close()

    req, _ := http.NewRequest("POST", "/users/auth", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestWeChatAuth_UpdateExistingUser_NoUpdate 测试已存在用户未提供昵称和头像
func TestWeChatAuth_UpdateExistingUser_NoUpdate(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 模拟微信 API 响应
    wxResponse := `{
        "openid": "test_openid",
        "session_key": "test_session_key"
    }`

    // 创建 httptest 服务器模拟微信 API
    wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, wxResponse)
    }))
    defer wxServer.Close()

    // 设置环境变量
    os.Setenv("WECHAT_API_URL", wxServer.URL)
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")
    os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
    defer func() {
        os.Unsetenv("WECHAT_API_URL")
        os.Unsetenv("APP_ID")
        os.Unsetenv("APP_SECRET")
        os.RemoveAll("./test_uploads")
    }()

    // 创建默认头像文件
    defaultAvatarPath := filepath.Join("./test_uploads", "avatars")
    os.MkdirAll(defaultAvatarPath, 0755)
    os.WriteFile(filepath.Join(defaultAvatarPath, "default.jpg"), []byte("default avatar"), 0644)

    // 创建已存在的用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "OldName",
        AvatarURL:  "avatars/old_avatar.jpg",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    db.Create(&user)

    // 保存旧头像文件
    oldAvatarPath := filepath.Join("./test_uploads", user.AvatarURL)
    os.WriteFile(oldAvatarPath, []byte("old avatar"), 0644)

    // 构建请求体，不提供昵称和头像
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("code", "test_code")
    // 不提供 nickname 和 avatar
    writer.Close()

    req, _ := http.NewRequest("POST", "/users/auth", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response["access_token"])
    assert.NotEmpty(t, response["refresh_token"])
    assert.NotEmpty(t, response["user"])

    userData := response["user"].(map[string]interface{})
    assert.Equal(t, "OldName", userData["nickname"]) // Nickname 未改变
    assert.Equal(t, user.AvatarURL, userData["avatar_url"]) // Avatar 未改变

    // 检查数据库中用户的昵称和 AvatarURL 是否未更新
    var updatedUser models.User
    result := db.Where("open_id = ?", "test_openid").First(&updatedUser)
    assert.Nil(t, result.Error)
    assert.Equal(t, "OldName", updatedUser.Nickname)
    assert.Equal(t, user.AvatarURL, updatedUser.AvatarURL)

    // 检查旧头像文件是否未被删除
    _, err = os.Stat(oldAvatarPath)
    assert.False(t, os.IsNotExist(err))
}

// TestSetNickname_Success 测试成功更新昵称
func TestSetNickname_Success(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 创建用户
	user := models.User{OpenID: "test_openid", Nickname: "OldName"}
	db.Create(&user)

	// 生成 Access Token
	token := generateValidJWT(user.ID)

	// 构建请求体
	bodyBytes, _ := json.Marshal(map[string]string{"nickname": "NewName"})
	req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Nickname updated successfully", response["message"])
	assert.Equal(t, "NewName", response["nickname"])

	// 检查数据库中用户的昵称是否更新
	var updatedUser models.User
	err = db.First(&updatedUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "NewName", updatedUser.Nickname)
}

// TestSetNickname_Unauthorized 测试未授权更新昵称
func TestSetNickname_Unauthorized(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	bodyBytes, _ := json.Marshal(map[string]string{"nickname": "NewName"})
	req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestSetNickname_NoNickname 测试缺少昵称参数
func TestSetNickname_NoNickname(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "OldName"}
	db.Create(&user)

	token := generateValidJWT(user.ID)

	bodyBytes, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSetNickname_UserNotFound 测试更新昵称时用户不存在
func TestSetNickname_UserNotFound(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	token := generateValidJWT(9999) // 不存在的用户ID

	bodyBytes, _ := json.Marshal(map[string]string{"nickname": "NewName"})
	req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestSetAvatar_Success 测试成功设置头像
func TestSetAvatar_Success(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 创建用户
	user := models.User{OpenID: "test_openid", Nickname: "TestUser", AvatarURL: "avatars/old_avatar.jpg"}
	db.Create(&user)

	// 生成 Access Token
	token := generateValidJWT(user.ID)

	// 创建上传文件目录
	os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
	defer func() {
		os.Unsetenv("BASE_UPLOAD_PATH")
		os.RemoveAll("./test_uploads")
	}()
	os.MkdirAll(filepath.Join("./test_uploads", "avatars"), 0755)

	// 创建旧头像文件
	oldAvatarPath := filepath.Join("./test_uploads", "avatars", "old_avatar.jpg")
	os.WriteFile(oldAvatarPath, []byte("old avatar"), 0644)

	// 创建新头像文件
	newAvatarPath := "./test_new_avatar.jpg"
	os.WriteFile(newAvatarPath, []byte("new avatar image"), 0644)
	defer os.Remove(newAvatarPath)

	// 构建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("avatar", filepath.Base(newAvatarPath))
	assert.NoError(t, err)

	file, err := os.Open(newAvatarPath)
	assert.NoError(t, err)
	defer file.Close()

	_, err = io.Copy(part, file)
	assert.NoError(t, err)
	writer.Close()

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Avatar updated successfully", response["message"])
	assert.Contains(t, response["avatar_url"], "avatars/")

	// 检查数据库中用户的 AvatarURL 是否更新
	var updatedUser models.User
	err = db.First(&updatedUser, user.ID).Error
	assert.NoError(t, err)
	assert.NotEqual(t, "avatars/old_avatar.jpg", updatedUser.AvatarURL)
	assert.Contains(t, updatedUser.AvatarURL, "avatars/")

	// 检查旧头像是否被删除
	_, err = os.Stat(oldAvatarPath)
	assert.True(t, os.IsNotExist(err))
}

// TestSetAvatar_Unauthorized 测试未授权设置头像
func TestSetAvatar_Unauthorized(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	req, _ := http.NewRequest("PUT", "/users/set_avatar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestSetAvatar_UserNotFound 测试设置头像时用户不存在
func TestSetAvatar_UserNotFound(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	token := generateValidJWT(9999) // 不存在的用户ID

	// 构建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.CreateFormFile("avatar", "test_avatar.jpg")
	writer.Close()

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestSetAvatar_NoFile 测试未上传文件
func TestSetAvatar_NoFile(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)

	token := generateValidJWT(user.ID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// 不上传文件
	writer.Close()

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSetAvatar_SaveFileFail 测试保存文件失败
func TestSetAvatar_SaveFileFail(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)

	token := generateValidJWT(user.ID)

	// 设置 BASE_UPLOAD_PATH 为一个不可写的路径，模拟保存失败
	os.Setenv("BASE_UPLOAD_PATH", "/invalid_path")
	defer func() {
		os.Unsetenv("BASE_UPLOAD_PATH")
	}()

	// 构建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.CreateFormFile("avatar", "test_avatar.jpg")
	writer.Close()

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestSetAvatar_UpdateDBFail 测试更新数据库失败
func TestSetAvatar_UpdateDBFail(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)

	token := generateValidJWT(user.ID)

	// 创建上传文件目录
	os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
	os.MkdirAll(filepath.Join("./test_uploads", "avatars"), 0755)
	defer func() {
		os.Unsetenv("BASE_UPLOAD_PATH")
		os.RemoveAll("./test_uploads")
	}()

	// 创建模拟文件
	filePath := "./test_avatar.jpg"
	os.WriteFile(filePath, []byte("this is a test image"), 0644)
	defer os.Remove(filePath)

	// 构建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("avatar", filepath.Base(filePath))
	assert.NoError(t, err)

	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	_, err = io.Copy(part, file)
	assert.NoError(t, err)
	writer.Close()

	// 先删除用户，让更新时报错
	db.Delete(&user)

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 用户已删除，保存更新失败
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRefreshTokenHandler_Success 测试成功刷新 Access Token
func TestRefreshTokenHandler_Success(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 创建用户
	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)

	// 生成 Refresh Token
	refreshToken := generateValidRefreshToken(db, user.ID)
	time.Sleep(time.Second)

	// 构建请求体
	reqBody := map[string]string{
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "access_token")
	assert.Contains(t, resp, "refresh_token")

	// 检查旧的 Refresh Token 是否被撤销
	var oldRT models.RefreshToken
	err = db.Where("token = ?", refreshToken).First(&oldRT).Error
	assert.NoError(t, err)
	assert.True(t, oldRT.Revoked)

	// 检查新的 Refresh Token 是否被创建
	var newRT models.RefreshToken
	err = db.Where("token = ?", resp["refresh_token"]).First(&newRT).Error
	assert.NoError(t, err)
	assert.False(t, newRT.Revoked)
	assert.Equal(t, user.ID, newRT.UserID)
}

// TestRefreshTokenHandler_InvalidToken 测试提供无效的 Refresh Token
func TestRefreshTokenHandler_InvalidToken(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 构建请求体
	reqBody := map[string]string{
		"refresh_token": "invalid_token",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRefreshTokenHandler_TokenNotFound 测试 Refresh Token 未找到
func TestRefreshTokenHandler_TokenNotFound(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 构建请求体
	reqBody := map[string]string{
		"refresh_token": "nonexistent_refresh_token",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRefreshTokenHandler_TokenRevoked 测试 Refresh Token 已被撤销
func TestRefreshTokenHandler_TokenRevoked(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 创建用户
	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)

	// 生成 Refresh Token
	refreshToken := generateValidRefreshToken(db, user.ID)

	// 撤销 Refresh Token
	var rt models.RefreshToken
	db.Where("token = ?", refreshToken).First(&rt)
	rt.Revoked = true
	db.Save(&rt)

	// 构建请求体
	reqBody := map[string]string{
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRefreshTokenHandler_TokenExpired 测试 Refresh Token 已过期
func TestRefreshTokenHandler_TokenExpired(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 创建用户
	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)

	// 生成 Refresh Token
	refreshToken := generateValidRefreshToken(db, user.ID)

	// 手动过期 Refresh Token
	var rt models.RefreshToken
	db.Where("token = ?", refreshToken).First(&rt)
	rt.ExpiresAt = time.Now().Add(-time.Hour)
	db.Save(&rt)

	// 构建请求体
	reqBody := map[string]string{
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRefreshTokenHandler_UserNotFound 测试 Refresh Token 对应的用户不存在
func TestRefreshTokenHandler_UserNotFound(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	// 生成 Refresh Token
	refreshToken, err := utils.GenerateRefreshToken(9999) // 不存在的用户ID
	assert.NoError(t, err)

	// 存储 Refresh Token
	newRT := models.RefreshToken{
		Token:     refreshToken,
		UserID:    9999,
		ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
		Revoked:   false,
	}
	db.Create(&newRT)

	// 构建请求体
	reqBody := map[string]string{
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言响应
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLogoutHandler_Success 测试成功登出
func TestLogoutHandler_Success(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 创建用户
    user := models.User{
        OpenID:    "test_openid",
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成并存储 Refresh Token
    refreshToken, err := utils.GenerateRefreshToken(user.ID)
    assert.NoError(t, err)

    newRT := models.RefreshToken{
        Token:     refreshToken,
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
        Revoked:   false,
    }
    err = db.Create(&newRT).Error
    assert.NoError(t, err)

    // 构建请求体
    reqBody := map[string]string{
        "refresh_token": refreshToken,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", "/users/logout", bytes.NewBuffer(bodyBytes))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 断言响应
    assert.Equal(t, http.StatusOK, w.Code)
    var resp map[string]string
    err = json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, "Logged out successfully", resp["message"])

    // 检查 Refresh Token 是否被撤销
    var storedRT models.RefreshToken
    err = db.Where("token = ?", refreshToken).First(&storedRT).Error
    assert.NoError(t, err)
    assert.True(t, storedRT.Revoked)
}

// TestLogoutHandler_InvalidToken 测试提供无效的 Refresh Token
func TestLogoutHandler_InvalidToken(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 构建请求体
    reqBody := map[string]string{
        "refresh_token": "invalid_token",
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", "/users/logout", bytes.NewBuffer(bodyBytes))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 断言响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLogoutHandler_TokenNotFound 测试 Refresh Token 未找到
func TestLogoutHandler_TokenNotFound(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 构建请求体
    reqBody := map[string]string{
        "refresh_token": "nonexistent_refresh_token",
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", "/users/logout", bytes.NewBuffer(bodyBytes))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 断言响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLogoutHandler_TokenRevoked 测试 Refresh Token 已被撤销
func TestLogoutHandler_TokenRevoked(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 创建用户
    user := models.User{
        OpenID:    "test_openid",
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成并存储 Refresh Token
    refreshToken, err := utils.GenerateRefreshToken(user.ID)
    assert.NoError(t, err)

    newRT := models.RefreshToken{
        Token:     refreshToken,
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
        Revoked:   true, // 已撤销
    }
    err = db.Create(&newRT).Error
    assert.NoError(t, err)

    // 构建请求体
    reqBody := map[string]string{
        "refresh_token": refreshToken,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", "/users/logout", bytes.NewBuffer(bodyBytes))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 断言响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestLogoutHandler_InvalidRequestBody 测试无效的请求体
func TestLogoutHandler_InvalidRequestBody(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 构建无效的请求体
    reqBody := "invalid_json"
    req, _ := http.NewRequest("POST", "/users/logout", bytes.NewBufferString(reqBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 断言响应
    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestLogoutHandler_UserNotFound 测试 Refresh Token 对应的用户不存在
func TestLogoutHandler_UserNotFound(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 生成 Refresh Token 对应不存在的用户
    refreshToken, err := utils.GenerateRefreshToken(9999)
    assert.NoError(t, err)

    // 存储 Refresh Token
    newRT := models.RefreshToken{
        Token:     refreshToken,
        UserID:    9999, // 不存在的用户 ID
        ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
        Revoked:   false,
    }
    db.Create(&newRT)

    // 构建请求体
    reqBody := map[string]string{
        "refresh_token": refreshToken,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    req, _ := http.NewRequest("POST", "/users/logout", bytes.NewBuffer(bodyBytes))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 根据实现，可能仍返回成功或 Unauthorized
    assert.Equal(t, http.StatusOK, w.Code)

    // 可选：检查 Refresh Token 是否被撤销
    var storedRT models.RefreshToken
    err = db.Where("token = ?", refreshToken).First(&storedRT).Error
    assert.NoError(t, err)
    assert.True(t, storedRT.Revoked)
}

// TestUserBasicDetails_Success 测试成功获取用户基本信息
func TestUserBasicDetails_Success(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 创建用户
    user := models.User{
        OpenID:    "test_openid",
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT 并模拟认证中间件
    token := generateValidJWT(user.ID)

    // 构建请求
    req, _ := http.NewRequest("GET", "/users/basic_details", nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)
    var resp map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)

    assert.Equal(t, float64(user.ID), resp["id"])
    assert.Equal(t, user.Nickname, resp["nickname"])
    assert.Equal(t, user.AvatarURL, resp["avatar_url"])
}

// TestUserBasicDetails_UserNotFound 测试用户不存在
func TestUserBasicDetails_UserNotFound(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 生成 JWT 对应不存在的用户 ID
    fakeUserID := uint(9999)
    token, err := utils.GenerateAccessToken(fakeUserID)
    assert.NoError(t, err)

    // 构建请求
    req, _ := http.NewRequest("GET", "/users/basic_details", nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusNotFound, w.Code)
    var resp map[string]string
    err = json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, "User not found", resp["error"])
}

// TestUserBasicDetails_InvalidTokenType 测试用户 ID 类型无效
func TestUserBasicDetails_InvalidTokenType(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 创建用户
    user := models.User{
        OpenID:    "test_openid",
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 手动创建一个无效的 JWT（user_id 为字符串）
    // 这里假设 utils.ValidateToken 会 parse the token and return claims with Subject as string
    // 但根据您之前的代码，似乎 Subject 是 string 类型，包含用户 ID
    // 因此，需要确保 GenerateAccessToken 符合 UserBasicDetails 的期望

    // 这里假设 Subject 是 string 的表示，可以使用 non-integer string
    // 需要根据实际情况调整

    invalidToken := "invalid_token_with_non_numeric_subject"

    // 构建请求
    req, _ := http.NewRequest("GET", "/users/basic_details", nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidToken))

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    var resp map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, "invalid token", resp["error"])
}

// TestUserBasicDetails_NoAuthHeader 测试缺少 Authorization 头
func TestUserBasicDetails_NoAuthHeader(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 构建请求
    req, _ := http.NewRequest("GET", "/users/basic_details", nil)
    // 不设置 Authorization 头

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestUserBasicDetails_InvalidAuthHeader 测试无效的 Authorization 头
func TestUserBasicDetails_InvalidAuthHeader(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 构建请求
    req, _ := http.NewRequest("GET", "/users/basic_details", nil)
    req.Header.Set("Authorization", "InvalidAuthHeader")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
