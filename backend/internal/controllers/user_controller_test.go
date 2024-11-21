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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
)

func setupTestDB() *gorm.DB {
    // 使用 SQLite 内存数据库
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    // 自动迁移模型
    if err := db.AutoMigrate(&models.User{}); err != nil {
        panic("failed to migrate User model")
    }
    return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
    router := gin.Default()

    // 手动注册相关路由
    userController := NewUserController(db)

    userGroup := router.Group("/users")
    {
        userGroup.POST("/auth", userController.WeChatAuth)

        authGroup := userGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            authGroup.PUT("/set_nickname", userController.SetNickname)
            authGroup.PUT("/set_avatar", userController.SetAvatar)
        }
    }

    return router
}

func TestWeChatAuth(t *testing.T) {
    db := setupTestDB()
    router := setupRouter(db)

    // 模拟微信 API 响应
    wxResponse := `{
        "open_id": "test_openid",
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
	os.Setenv("BASE_UPLOAD_PATH", "/home/ubuntu/qiaoshenyu/se/backend/upload")
	defer db.Exec("DROP TABLE users")

    // 创建默认头像文件
    defaultAvatarDir := "/tmp/test_uploads/avatars"
    os.MkdirAll(defaultAvatarDir, 0755)
    defaultAvatarPath := fmt.Sprintf("%s/default.jpg", defaultAvatarDir)
    os.WriteFile(defaultAvatarPath, []byte("default avatar"), 0644)

    // 构建请求体
    requestBody := map[string]string{
        "code":     "test_code",
        "nickname": "TestUser",
    }
    jsonBody, _ := json.Marshal(requestBody)

    // 创建请求
    req, _ := http.NewRequest("POST", "/users/auth", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotEmpty(t, response["token"])
    assert.NotEmpty(t, response["user"])

    userData := response["user"].(map[string]interface{})
    assert.Equal(t, "TestUser", userData["nickname"])
    assert.NotEmpty(t, userData["avatar_url"])

    // 检查数据库中是否存在用户
    var user models.User
    result := db.Where("open_id = ?", "test_openid").First(&user)
    assert.Nil(t, result.Error)
    assert.Equal(t, "TestUser", user.Nickname)

    // 检查头像文件是否存在
    avatarPath := fmt.Sprintf("/home/ubuntu/qiaoshenyu/se/backend/upload/%s", user.AvatarURL)
    _, err := os.Stat(avatarPath)
    assert.False(t, os.IsNotExist(err))
}

func TestSetNickname(t *testing.T) {
    db := setupTestDB()
    router := setupRouter(db)

    // 创建测试用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "OldNickname",
        SessionKey: "test_session_key",
    }
    db.Create(&user)

    // 生成 JWT
    token, _ := utils.GenerateJWT(user.ID)

    // 构建请求体
    requestBody := map[string]string{
        "nickname": "NewNickname",
    }
    jsonBody, _ := json.Marshal(requestBody)

    // 创建请求
    req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "Nickname updated successfully", response["message"])
    assert.Equal(t, "NewNickname", response["nickname"])

    // 检查数据库更新
    var updatedUser models.User
    db.First(&updatedUser, user.ID)
    assert.Equal(t, "NewNickname", updatedUser.Nickname)
}

func TestSetAvatar(t *testing.T) {
	os.Setenv("BASE_UPLOAD_PATH", "/home/ubuntu/qiaoshenyu/se/backend/upload")

    db := setupTestDB()
    router := setupRouter(db)

    // 创建测试用户
    user := models.User{
        OpenID:     "test_openid",
        Nickname:   "TestUser",
        SessionKey: "test_session_key",
    }
    db.Create(&user)

    // 生成 JWT
    token, _ := utils.GenerateJWT(user.ID)

    // 创建模拟文件
    filePath := "./test_avatar.jpg"
    fileContent := []byte("this is a test image")
    os.WriteFile(filePath, fileContent, 0644)
    defer os.Remove(filePath) // 测试结束后删除文件

    // 创建表单文件
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, _ := writer.CreateFormFile("avatar", filepath.Base(filePath))
    file, _ := os.Open(filePath)
    io.Copy(part, file)
    file.Close()
    writer.Close()

    // 创建请求
    req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "Avatar updated successfully", response["message"])
    assert.NotEmpty(t, response["avatar_url"])

    // 检查数据库更新
    var updatedUser models.User
    db.First(&updatedUser, user.ID)
    assert.Equal(t, response["avatar_url"], updatedUser.AvatarURL)

	fmt.Println(updatedUser)

    // 检查文件是否存在
    BaseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
    if BaseUploadPath == "" {
        BaseUploadPath = "./uploads"
    }
    avatarPath := fmt.Sprintf("%s/%s.jpg", BaseUploadPath, updatedUser.AvatarURL)
    _, err := os.Stat(avatarPath)
    assert.False(t, os.IsNotExist(err))
}