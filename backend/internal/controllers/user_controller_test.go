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
	// "github.com/DATA-DOG/go-sqlmock"
)

// // setupUserTestDBWithMock 设置带有 sqlmock 的 GORM 数据库
// func setupUserTestDBWithMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
//     db, mock, err := sqlmock.New()
//     if err != nil {
//         t.Fatalf("failed to open sqlmock database: %v", err)
//     }

//     // 初始化 GORM 使用 sqlmock 的连接
//     dialector := sqlite.New(sqlite.Config{
//         Conn: db,
//     })
//     gormDB, err := gorm.Open(dialector, &gorm.Config{})
//     if err != nil {
//         t.Fatalf("failed to open gorm DB: %v", err)
//     }

//     // 设置 GORM 初始化时的期望查询
//     mock.ExpectExec("PRAGMA foreign_keys = ON").WillReturnResult(sqlmock.NewResult(0, 0))
//     mock.ExpectExec("PRAGMA journal_mode = WAL").WillReturnResult(sqlmock.NewResult(0, 0))
//     mock.ExpectQuery("SELECT sqlite_version\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.31.1"))

//     // 设置 AutoMigrate 时的 CREATE TABLE 期望
//     mock.ExpectExec("CREATE TABLE IF NOT EXISTS `users` .*").WillReturnResult(sqlmock.NewResult(1, 1))
//     mock.ExpectExec("CREATE TABLE IF NOT EXISTS `families` .*").WillReturnResult(sqlmock.NewResult(1, 1))
//     mock.ExpectExec("CREATE TABLE IF NOT EXISTS `news` .*").WillReturnResult(sqlmock.NewResult(1, 1))

//     // 如果 GORM 创建索引，设置相应的期望
//     mock.ExpectExec("CREATE INDEX .*").WillReturnResult(sqlmock.NewResult(1, 1))

//     // 执行 AutoMigrate
//     if err := gormDB.AutoMigrate(&models.User{}, &models.Family{}, &models.News{}); err != nil {
//         t.Fatalf("failed to migrate models: %v", err)
//     }

//     return gormDB, mock
// }

func setupUserTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate User model")
	}
	return db
}

func setupUserRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
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
	os.Setenv("BASE_UPLOAD_PATH", "/tmp/test_uploads")
	defer db.Exec("DROP TABLE users")

	os.MkdirAll("/tmp/test_uploads/avatars", 0755)
	defaultAvatarPath := "/tmp/test_uploads/avatars/default.jpg"
	os.WriteFile(defaultAvatarPath, []byte("default avatar"), 0644)

	// 构建请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("code", "test_code")
	writer.WriteField("nickname", "TestUser")
	writer.Close()

	req, _ := http.NewRequest("POST", "/users/auth", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
	assert.Contains(t, userData["avatar_url"], "avatars/")

	// 检查数据库中是否存在用户
	var user models.User
	result := db.Where("open_id = ?", "test_openid").First(&user)
	assert.Nil(t, result.Error)
	assert.Equal(t, "TestUser", user.Nickname)
}

func TestWeChatAuth_NoCode(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	req, _ := http.NewRequest("POST", "/users/auth", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWeChatAuth_WeChatAPIError(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db)

    // 设置 WECHAT_API_URL 为一个不可访问的地址，模拟网络错误
    os.Setenv("WECHAT_API_URL", "http://localhost:9999") // 假设没有服务在此端口运行
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")
    os.Setenv("BASE_UPLOAD_PATH", "/tmp/test_uploads")
    defer db.Migrator().DropTable(&models.User{}, &models.Family{}, &models.News{}) // 修改为适当的迁移删除

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

// // TestWeChatAuth_CreateUserFail 测试创建用户失败
// func TestWeChatAuth_CreateUserFail(t *testing.T) {
//     // 使用 sqlmock 设置数据库和预期
//     gormDB, mock := setupUserTestDBWithMock(t)
//     router := setupUserRouter(gormDB)

//     // 模拟微信 API 正常返回
//     wxResponse := `{
//         "openid": "new_openid",
//         "session_key": "new_session_key"
//     }`

//     // 创建 httptest 服务器模拟微信 API
//     wxServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintln(w, wxResponse)
//     }))
//     defer wxServer.Close()

//     // 设置环境变量
//     os.Setenv("WECHAT_API_URL", wxServer.URL)
//     os.Setenv("APP_ID", "test_app_id")
//     os.Setenv("APP_SECRET", "test_app_secret")
//     os.Setenv("BASE_UPLOAD_PATH", "/tmp/test_uploads")
//     defer os.Unsetenv("WECHAT_API_URL")
//     defer os.Unsetenv("APP_ID")
//     defer os.Unsetenv("APP_SECRET")
//     defer os.Unsetenv("BASE_UPLOAD_PATH")

//     // 设置预期的 SQL 查询和插入行为
//     // 查询用户是否存在
//     mock.ExpectQuery("SELECT (.+) FROM users WHERE open_id = ?").
//         WithArgs("new_openid").
//         WillReturnRows(sqlmock.NewRows([]string{"id", "open_id", "session_key", "nickname", "family_id", "created_at", "updated_at"}))

//     // 开始事务
//     mock.ExpectBegin()

//     // 插入用户时返回错误
//     mock.ExpectExec("INSERT INTO `users`").
//         WithArgs(sqlmock.AnyArg(), "new_openid", "new_session_key", "TestUser", nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
//         WillReturnError(fmt.Errorf("insert error"))

//     // 回滚事务
//     mock.ExpectRollback()

//     // 构建请求体
//     body := &bytes.Buffer{}
//     writer := multipart.NewWriter(body)
//     writer.WriteField("code", "test_code")
//     writer.WriteField("nickname", "TestUser")
//     writer.Close()

//     // 创建请求
//     req, _ := http.NewRequest("POST", "/users/auth", body)
//     req.Header.Set("Content-Type", writer.FormDataContentType())

//     // 执行请求
//     w := httptest.NewRecorder()
//     router.ServeHTTP(w, req)

//     // 检查响应
//     assert.Equal(t, http.StatusInternalServerError, w.Code)
//     assert.Contains(t, w.Body.String(), "Failed to create user")

//     // 确认所有预期的 SQL 操作被满足
//     if err := mock.ExpectationsWereMet(); err != nil {
//         t.Errorf("there were unfulfilled expectations: %s", err)
//     }
// }

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

func TestSetNickname_NoNickname(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "OldName"}
	db.Create(&user)

	token, _ := utils.GenerateJWT(user.ID)

	bodyBytes, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetNickname_UserNotFound(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	token, _ := utils.GenerateJWT(9999) // 不存在的用户ID

	bodyBytes, _ := json.Marshal(map[string]string{"nickname": "NewName"})
	req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSetAvatar_Unauthorized(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	req, _ := http.NewRequest("PUT", "/users/set_avatar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSetAvatar_UserNotFound(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	token, _ := utils.GenerateJWT(9999) // 不存在用户

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSetAvatar_NoFile(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)

	token, _ := utils.GenerateJWT(user.ID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetAvatar_SaveFileFail(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)
	token, _ := utils.GenerateJWT(user.ID)

	os.Setenv("BASE_UPLOAD_PATH", "/not-exist-path")

	// 构造上传请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.CreateFormFile("avatar", "test.jpg")
	writer.Close()

	req, _ := http.NewRequest("PUT", "/users/set_avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 因为目录不存在，保存文件失败
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSetAvatar_UpdateDBFail(t *testing.T) {
	db := setupUserTestDB()
	router := setupUserRouter(db)

	user := models.User{OpenID: "test_openid", Nickname: "TestUser"}
	db.Create(&user)
	token, _ := utils.GenerateJWT(user.ID)

	// 创建上传文件目录
	os.Setenv("BASE_UPLOAD_PATH", "/tmp/test_uploads")
	os.MkdirAll("/tmp/test_uploads/avatars", 0755)

	// 创建模拟文件
	filePath := "./test_avatar.jpg"
	os.WriteFile(filePath, []byte("this is a test image"), 0644)
	defer os.Remove(filePath)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("avatar", filepath.Base(filePath))
	file, _ := os.Open(filePath)
	io.Copy(part, file)
	file.Close()
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