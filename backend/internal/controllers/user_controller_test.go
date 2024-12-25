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
	// "path/filepath"
	"time"

	// "strconv"
	"testing"
	"strings"
	"github.com/golang-jwt/jwt/v4"

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
func setupUserRouter(db *gorm.DB, utils utils.UtilsInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userController := NewUserController(db, utils)

	userGroup := router.Group("/users")
    {
        // 公共路由
        userGroup.POST("/auth", userController.WeChatAuth) // 注册
        userGroup.POST("/refresh", userController.RefreshTokenHandler) // 刷新令牌
        userGroup.POST("/logout", userController.LogoutHandler) // 登出

        // 需要认证的路由
        authGroup := userGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            authGroup.PUT("/set_nickname", userController.SetNickname) // 更新用户名
            authGroup.POST("/set_avatar", userController.SetAvatar) // 更新头像
            authGroup.GET("/basic_details", userController.UserBasicDetails) // 获取基本信息

            authGroup.GET("/liked", userController.GetMyLikedNews)
            authGroup.GET("/favorited", userController.GetMyFavoritedNews)
            authGroup.GET("/viewed", userController.GetMyViewedNews)

            authGroup.GET("/:id/profile", userController.GetUserProfile)
        }
    }

	return router
}

// Helper function to generate a valid JWT for testing
func generateValidJWTUser(userID uint) string {
	token, err := utils.GenerateAccessToken(userID)
	if err != nil {
		panic("Failed to generate valid JWT for testing")
	}
	return token
}

type MockUtils struct {
	ValidateTokenFunc        func(tokenString string) (*jwt.RegisteredClaims, error)
    GenerateAccessTokenFunc  func(userID uint) (string, error)
    GenerateRefreshTokenFunc func(userID uint) (string, error)
    CopyFileFunc             func(src, dst string) error
}

func (m *MockUtils) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
    if m.ValidateTokenFunc != nil {
        return m.ValidateTokenFunc(tokenString)
    }
    // 默认行为，模拟解析成功
    return &jwt.RegisteredClaims{Subject: "123"}, nil
}

func (m *MockUtils) GenerateAccessToken(userID uint) (string, error) {
    if m.GenerateAccessTokenFunc != nil {
        return m.GenerateAccessTokenFunc(userID)
    }
    return fmt.Sprintf("MockAccessToken_%d", userID), nil
}

func (m *MockUtils) GenerateRefreshToken(userID uint) (string, error) {
    if m.GenerateRefreshTokenFunc != nil {
        return m.GenerateRefreshTokenFunc(userID)
    }
    return fmt.Sprintf("MockRefreshToken_%d", userID), nil
}

func (m *MockUtils) CopyFile(src, dst string) error {
    if m.CopyFileFunc != nil {
        return m.CopyFileFunc(src, dst)
    }
    return nil
}

// MockRoundTripper 用于模拟 http.Client 的 Transport，从而在测试中自定义响应
type MockRoundTripper struct {
    RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip 实现 http.RoundTripper 接口
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.RoundTripFunc(req)
}

func TestWeChatAuth(t *testing.T) {
    // 1. 初始化测试DB和路由
    db := setupUserTestDB()
    mockUtils := &MockUtils{
        GenerateAccessTokenFunc: func(userID uint) (string, error) {
            return fmt.Sprintf("AccessToken_ForUser_%d", userID), nil
        },
        GenerateRefreshTokenFunc: func(userID uint) (string, error) {
            return fmt.Sprintf("RefreshToken_ForUser_%d", userID), nil
        },
        CopyFileFunc: func(src, dst string) error {
            return nil
        },
    }

    router := setupUserRouter(db, mockUtils)

    // 2. 设置一些环境变量 (可根据需要修改)
    os.Setenv("WECHAT_API_URL", "http://mock-wechat-api.com/sns/jscode2session")
    os.Setenv("APP_ID", "test_app_id")
    os.Setenv("APP_SECRET", "test_app_secret")

    // 3. 定义一个 mock HTTP server or transport 来模拟微信API

    // ------- 准备一个可修改的 transport，用于模拟微信API返回 --------
    var mockWeChatTransport http.RoundTripper = &MockRoundTripper{
        RoundTripFunc: func(req *http.Request) (*http.Response, error) {
            if strings.Contains(req.URL.String(), "call_wechat_api_fail") {
                // 返回网络错误
                return nil, fmt.Errorf("forced wechat api call error")
            }

            if strings.Contains(req.URL.String(), "bad_json_response") {
                // 返回一个非JSON
                return &http.Response{
                    StatusCode: 200,
                    Body:       io.NopCloser(bytes.NewBufferString("not_json")),
                }, nil
            }

            if strings.Contains(req.URL.String(), "errcode_nonzero") {
                // 返回一个微信错误，如 errcode=40029
                respStr := `{"errcode": 40029, "errmsg": "invalid code"}` 
                return &http.Response{
                    StatusCode: 200,
                    Body:       io.NopCloser(bytes.NewBufferString(respStr)),
                }, nil
            }

            // 正常返回 => openid & session_key
            respStr := `{"openid": "OpenID_WeChatAuthTest", "session_key": "SessionKey_12345"}`
            return &http.Response{
                StatusCode: 200,
                Body:       io.NopCloser(bytes.NewBufferString(respStr)),
            }, nil
        },
    }

    oldTransport := http.DefaultTransport
    http.DefaultTransport = mockWeChatTransport
    defer func() {
        http.DefaultTransport = oldTransport
    }()

    // 4. 准备表驱动测试用例
    tests := []struct {
        name           string
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
		{
            name: "Success existing user",
            requestBody: gin.H{"code": "normal_code"},
            setupFunc: func() {
                // 移除回调
                // mock token生成都正常
                mockUtils.GenerateRefreshTokenFunc = func(userID uint) (string, error) {
                    return fmt.Sprintf("RefreshToken_ForUser_%d", userID), nil
                }
                // 预先创建一个User(已有User) => openID=OpenID_WeChatAuthTest
                userExist := models.User{
                    OpenID:     "OpenID_WeChatAuthTest",
                    SessionKey: "OldSessionKey",
                }
                db.Create(&userExist)
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
        {
            name: "Success new user",
            requestBody: gin.H{"code": "normal_code"},
            setupFunc: func() {
                // 让 openID=OpenID_WeChatAuthTest => user 不存在 => 走创建逻辑
                db.Where("open_id = ?", "OpenID_WeChatAuthTest").Delete(&models.User{}) 
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
        {
            name:           "Invalid JSON body",
            requestBody:    "not_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "Empty code field",
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "WeChat API call fail",
            requestBody:    gin.H{"code": "call_wechat_api_fail"},
            setupFunc:      func() {},
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to call WeChat API",
        },
        {
            name:           "WeChat response parse fail (bad JSON)",
            requestBody:    gin.H{"code": "bad_json_response"},
            setupFunc:      func() {},
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to parse WeChat API response",
        },
        {
            name:           "WeChat response errcode nonzero",
            requestBody:    gin.H{"code": "errcode_nonzero"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "invalid code",
        },
        {
            name:           "DB error when checking user (simulate)",
            requestBody:    gin.H{"code": "normal_code"},
            setupFunc: func() {
                // 模拟 uc.DB.Preload("RefreshTokens").Where("open_id = ?").First(&user).Error 出错
                db.Callback().Query().Before("gorm:query").Register("force_query_user_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced query user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Database error",
        },
        {
            name: "Copy default avatar fail",
            requestBody: gin.H{"code": "normal_code"},
            setupFunc: func() {
				db.Callback().Query().Remove("force_query_user_err")  // 移除上一条mock
                db.Callback().Create().Remove("force_create_user_err")
                // 在 CopyFile 里出错
                mockUtils.CopyFileFunc = func(src, dst string) error {
                    return fmt.Errorf("forced copy file error")
                }
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to store refresh token",
        },
        {
            name: "Fail to generate access token",
            requestBody: gin.H{"code": "normal_code"},
            setupFunc: func() {
                db.Callback().Update().Remove("force_update_avatar_err")
                // mock token生成
                mockUtils.GenerateAccessTokenFunc = func(userID uint) (string, error) {
                    return "", fmt.Errorf("forced access token gen error")
                }
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to generate access token",
        },
        {
            name: "Fail to generate refresh token",
            requestBody: gin.H{"code": "normal_code"},
            setupFunc: func() {
                mockUtils.GenerateAccessTokenFunc = func(userID uint) (string, error) {
                    return fmt.Sprintf("AccessToken_ForUser_%d", userID), nil
                }
                mockUtils.GenerateRefreshTokenFunc = func(userID uint) (string, error) {
                    return "", fmt.Errorf("forced refresh token gen error")
                }
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to generate refresh token",
        },
        {
            name: "Fail to store refresh token",
            requestBody: gin.H{"code": "normal_code"},
            setupFunc: func() {
                mockUtils.GenerateRefreshTokenFunc = func(userID uint) (string, error) {
                    return "RefreshToken_Abc123", nil
                }
                // mock db.Create(&newRefreshToken).Error 出错
                db.Callback().Create().Before("gorm:create").Register("force_create_refreshtoken_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "refresh_tokens" {
                        tx.Error = fmt.Errorf("forced create refresh token error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to store refresh token",
        },
    }

    // 5. 依次运行测试
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/users/auth", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                // 成功分支，可进一步校验 access_token, refresh_token, user 等字段
                _, hasAccess := resp["access_token"]
                _, hasRefresh := resp["refresh_token"]
                _, hasUser := resp["user"]
                assert.True(t, hasAccess)
                assert.True(t, hasRefresh)
                assert.True(t, hasUser)
            } else if tc.expectedError != "" {
                // 检查 error
                assert.Equal(t, tc.expectedError, resp["error"])
            }
        })
    }
}

func TestSetNickname(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db, utils.UtilsImpl{})

    // 创建一个需要鉴权的路由 => PUT /users/set_nickname
    // 在 setupUserRouter 里，authGroup.PUT("/set_nickname", userController.SetNickname)

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_SetNickname_Test",
        Nickname: "OldNickname",
    }
    db.Create(&user)

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedNew    string
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"nickname": "NewNick"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Request Body (missing nickname)",
            userID:         user.ID,
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "User Not Found",
            userID:         99999, // 不存在
            requestBody:    gin.H{"nickname": "NewNick"},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "User not found",
        },
        {
            name:           "DB Error While Updating Nickname",
            userID:         user.ID,
            requestBody:    gin.H{"nickname": "NewNick"},
            setupFunc: func() {
                // 模拟更新 user 时的错误 => Save(&user).Error
                db.Callback().Update().Before("gorm:update").Register("force_update_nickname_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced update nickname error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update nickname",
        },
        {
            name:           "Success Set Nickname",
            userID:         user.ID,
            requestBody:    gin.H{"nickname": "NewNick"},
            setupFunc: func() {
                // 移除回调
                db.Callback().Update().Remove("force_update_nickname_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            expectedNew:    "NewNick",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("PUT", "/users/set_nickname", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTUser(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                // 成功更新
                assert.Equal(t, "Nickname updated successfully", resp["message"])
                assert.Equal(t, tc.expectedNew, resp["nickname"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestSetAvatar(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db, utils.UtilsImpl{})

    // 创建用户
    user := models.User{
        OpenID:    "OpenID_SetAvatar_Test",
        Nickname:  "UserSetAvatar",
        AvatarURL:  "avatars/old_avatar.jpg",
    }
    db.Create(&user)

    tests := []struct {
        name           string
        userID         uint
        fileField      string   // form-data 中的字段名
        fileName       string   // 上传的文件名
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            fileField:      "avatar",
            fileName:       "test_avatar.jpg",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "User Not Found",
            userID:         99999,
            fileField:      "avatar",
            fileName:       "test_avatar.jpg",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "User not found",
        },
        {
            name:           "No File Provided",
            userID:         user.ID,
            fileField:      "", // 不提供 file
            fileName:       "",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Failed to retrieve file",
        },
        {
            name:           "Failed to Update DB Avatar",
            userID:         user.ID,
            fileField:      "avatar",
            fileName:       "test_avatar.jpg",
            setupFunc: func() {
                // 可以使用 GORM update callback
                db.Callback().Update().Before("gorm:update").Register("force_update_avatar_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced update avatar error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update avatar",
        },
        {
            name:           "Success Set Avatar",
            userID:         user.ID,
            fileField:      "avatar",
            fileName:       "test_avatar.jpg",
            setupFunc: func() {
                // 移除上一个callback
                db.Callback().Update().Remove("force_update_avatar_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            // 构造 multipart/form-data
            body := &bytes.Buffer{}
            writer := multipart.NewWriter(body)
            if tc.fileField != "" {
                part, _ := writer.CreateFormFile(tc.fileField, tc.fileName)
                // 写入少许字节模拟图片
                part.Write([]byte("fake_image_data"))
            }
            writer.Close()

            req, _ := http.NewRequest("POST", "/users/set_avatar", body)
            req.Method = "POST" // 在 router 里是 authGroup.POST, 这里示例函数是 SetAvatar(POST)
            req.Header.Set("Content-Type", writer.FormDataContentType())

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTUser(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                // 成功
                assert.Equal(t, "Avatar updated successfully", resp["message"])
                // resp["avatar_url"] 应该存在
                _, ok := resp["avatar_url"]
                assert.True(t, ok)
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestRefreshTokenHandler(t *testing.T) {
    db := setupUserTestDB()
    // 准备一个 mockUtils 用于 ValidateToken / GenerateAccessToken / GenerateRefreshToken
    mockUtils := &MockUtils{
        ValidateTokenFunc: func(tokenString string) (*jwt.RegisteredClaims, error) {
            // 缺省实现 => 正常解析
            // 具体可根据 tokenString 做多分支
            if tokenString == "InvalidTokenString" {
                return nil, fmt.Errorf("forced invalid token error")
            }
            if tokenString == "BadSubjectToken" {
                // Subject不合法
                return &jwt.RegisteredClaims{Subject: "non_integer_id"}, nil
            }
            // 缺省 => userID= 123
            return &jwt.RegisteredClaims{Subject: "123"}, nil
        },
        GenerateAccessTokenFunc: func(userID uint) (string, error) {
            return fmt.Sprintf("AccessToken_ForUser_%d", userID), nil
        },
        GenerateRefreshTokenFunc: func(userID uint) (string, error) {
            return fmt.Sprintf("RefreshToken_ForUser_%d", userID), nil
        },
    }
    router := setupUserRouter(db, mockUtils)

    // 创建用户
    user := models.User{
        OpenID: "OpenID_RefreshTester",
    }
    db.Create(&user)

    // 生成并存储一个合法的 old refresh token => token= "OldRefresh_123"
    oldRT := models.RefreshToken{
        Token:     "OldRefresh_123",
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24*time.Hour),
        Revoked:   false,
    }
    db.Create(&oldRT)

    tests := []struct{
        name           string
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Invalid Request Body",
            requestBody:    "not_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "Empty Refresh Token",
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "ValidateToken Error",
            requestBody:    gin.H{"refresh_token": "InvalidTokenString"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "forced invalid token error",
        },
        {
            name:           "Invalid Token Subject",
            requestBody:    gin.H{"refresh_token": "BadSubjectToken"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Invalid token subject",
        },
        {
            name:           "Refresh Token Not Found",
            requestBody:    gin.H{"refresh_token": "NotInDB"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Refresh token not found",
        },
        {
            name:           "Refresh Token Expired or Revoked",
            requestBody:    gin.H{"refresh_token": "OldRefresh_123"},
            setupFunc: func() {
                // 让 oldRT 过期或 revoked
                db.Model(&oldRT).Update("revoked", true)
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Refresh token is expired or revoked",
        },
        {
            name:           "User Not Found for this token",
            requestBody:    gin.H{"refresh_token": "OldRefresh_123"},
            setupFunc: func() {
                // 重新插回 oldRT, 并指向 不存在的UserID=999
                db.Model(&oldRT).UpdateColumns(map[string]interface{}{
                    "revoked": false,
                    "user_id": 999,
                })
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "User not found",
        },
        {
            name:           "Fail to generate new access token",
            requestBody:    gin.H{"refresh_token": "OldRefresh_123"},
            setupFunc: func() {
                // 把 userID 改回
                db.Model(&oldRT).Update("user_id", user.ID)
                // mock GenerateAccessToken => error
                mockUtils.GenerateAccessTokenFunc = func(uid uint) (string, error) {
                    return "", fmt.Errorf("forced access token error")
                }
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to generate access token",
        },
        {
            name: "Fail to generate new refresh token",
            requestBody: gin.H{"refresh_token": "OldRefresh_123"},
            setupFunc: func() {
                mockUtils.GenerateAccessTokenFunc = func(uid uint) (string, error) {
                    return "AccessToken_Success", nil
                }
                mockUtils.GenerateRefreshTokenFunc = func(uid uint) (string, error) {
                    return "", fmt.Errorf("forced refresh token error")
                }
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to generate refresh token",
        },
        {
            name:           "Fail to store new refresh token",
            requestBody:    gin.H{"refresh_token": "OldRefresh_123"},
            setupFunc: func() {
                mockUtils.GenerateRefreshTokenFunc = func(uid uint) (string, error) {
                    return "NewRefresh_ABC", nil
                }
                // mock db.Create(newRT) => error
                db.Callback().Create().Before("gorm:create").Register("force_create_newRT_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "refresh_tokens" {
                        tx.Error = fmt.Errorf("forced create new RT error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to store new refresh token",
        },
        {
            name:           "Fail to revoke old refresh token",
            requestBody:    gin.H{"refresh_token": "OldRefresh_123"},
            setupFunc: func() {
                db.Callback().Create().Remove("force_create_newRT_err")
                // mock db.Save(&storedRefreshToken).Error => revoke old token
                db.Callback().Update().Before("gorm:update").Register("force_revoke_old_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "refresh_tokens" {
                        tx.Error = fmt.Errorf("forced revoke old RT error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to revoke old refresh token",
        },
        {
            name:           "Success Refresh Token",
            requestBody:    gin.H{"refresh_token": "OldRefresh_123"},
            setupFunc: func() {
                db.Callback().Update().Remove("force_revoke_old_err")
                // mock tokens => success
                mockUtils.GenerateAccessTokenFunc = func(uid uint) (string, error) {
                    return "AccessToken_Success", nil
                }
                mockUtils.GenerateRefreshTokenFunc = func(uid uint) (string, error) {
                    return "NewRefresh_Success", nil
                }
                // 让 oldRT 不过期
                db.Model(&oldRT).UpdateColumns(map[string]interface{}{
                    "revoked": false,
                    "user_id": user.ID,
                    "expires_at": time.Now().Add(24*time.Hour),
                })
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                _, ok1 := resp["access_token"]
                _, ok2 := resp["refresh_token"]
                assert.True(t, ok1 && ok2)
            } else if tc.expectedError != "" {
                assert.Equal(t, tc.expectedError, resp["error"])
            }
        })
    }
}

func TestLogoutHandler(t *testing.T) {
    db := setupUserTestDB()
    mockUtils := &MockUtils{
        ValidateTokenFunc: func(tokenString string) (*jwt.RegisteredClaims, error) {
            // 默认情况 => subject= "123"
            if tokenString == "InvalidTokenString" {
                return nil, fmt.Errorf("invalid refresh token")
            }
            if tokenString == "BadSubjectToken" {
                return &jwt.RegisteredClaims{Subject:"non_int"}, nil
            }
            return &jwt.RegisteredClaims{Subject:"123"}, nil
        },
        // 其余方法不涉及
    }
    router := setupUserRouter(db, mockUtils)

    // 创建用户
    user := models.User{
        OpenID: "OpenID_LogoutTester",
    }
    db.Create(&user)

    // 创建 refresh token => "LogoutRefresh_123"
    rt := models.RefreshToken{
        Token:     "LogoutRefresh_123",
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24*time.Hour),
        Revoked:   false,
    }
    db.Create(&rt)

    // 再给此用户多创建几个 refresh token => 也要同时撤销
    rtExtra := models.RefreshToken{
        Token:     "ExtraRefresh_1",
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24*time.Hour),
        Revoked:   false,
    }
    db.Create(&rtExtra)

    tests := []struct {
        name           string
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Invalid Request Body",
            requestBody:    "not_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "Empty RefreshToken Field",
            requestBody:    gin.H{},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "ValidateToken Error",
            requestBody:    gin.H{"refresh_token": "InvalidTokenString"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Invalid refresh token",
        },
        {
            name:           "Invalid Token Subject",
            requestBody:    gin.H{"refresh_token": "BadSubjectToken"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Invalid token subject",
        },
        {
            name:           "Refresh Token Not Found",
            requestBody:    gin.H{"refresh_token": "NotInDB"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Refresh token not found",
        },
        {
            name:           "Refresh Token Expired Or Revoked",
            requestBody:    gin.H{"refresh_token": "LogoutRefresh_123"},
            setupFunc: func() {
                db.Model(&rt).Update("revoked", true)
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Refresh token is expired or revoked",
        },
        {
            name:           "Failed To Revoke Current Refresh Token",
            requestBody:    gin.H{"refresh_token": "LogoutRefresh_123"},
            setupFunc: func() {
                db.Model(&rt).UpdateColumns(map[string]interface{}{
                    "revoked": false,
                    "expires_at": time.Now().Add(24*time.Hour),
                })

                // mock db.Save(&storedRefreshToken).Error
                db.Callback().Update().Before("gorm:update").Register("force_revoke_current_token_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "refresh_tokens" {
                        tx.Error = fmt.Errorf("forced revoke current token error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to revoke refresh token",
        },
        {
            name:           "Success Logout",
            requestBody:    gin.H{"refresh_token": "LogoutRefresh_123"},
            setupFunc: func() {
                db.Callback().Update().Remove("force_revoke_all_tokens_err")
				db.Callback().Update().Remove("force_revoke_current_token_err")

                // 让 token 仍不过期
                db.Model(&rt).Update("revoked", false)
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/users/logout", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "Logged out successfully", resp["message"])
            } else if tc.expectedError != "" {
                assert.Equal(t, tc.expectedError, resp["error"])
            }
        })
    }
}

func TestUserBasicDetails(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db, utils.UtilsImpl{})

    // 创建用户
    user := models.User{
        OpenID:    "OpenID_BasicDetails_Test",
        Nickname:  "BasicTester",
        AvatarURL: "avatars/default.jpg",
        CreatedAt: time.Now().Add(-48 * time.Hour), // 2天前创建
    }
    db.Create(&user)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectDays     int
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "User Not Found",
            userID:         99999,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "User not found",
        },
        {
            name:           "Success Basic Details",
            userID:         user.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            expectDays:     2, // 因为CreatedAt是2天前
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("GET", "/users/basic_details", nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTUser(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                // 检查 fields
                _, ok := resp["id"]
                assert.True(t, ok)
                assert.Equal(t, float64(tc.expectDays), resp["registered_days"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestGetMyFavoritedNews(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db, utils.UtilsImpl{})

    user := models.User{
        OpenID:   "OpenID_FavoritedNews_Test",
        Nickname: "FavorTester",
    }
    db.Create(&user)

    // 创建一些新闻
    news1 := models.News{Title: "News A"}
    db.Create(&news1)
    news2 := models.News{Title: "News B"}
    db.Create(&news2)

    // 让 user 收藏 news1
    db.Model(&user).Association("FavoritedNews").Append(&news1)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "User Not Found",
            userID:         99999,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "User not found",
        },
        {
            name:           "DB Error (simulate)",
            userID:         user.ID,
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_db_err_favorited", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced db error in favorited news")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to query database",
        },
        {
            name:           "Success GetMyFavoritedNews",
            userID:         user.ID,
            setupFunc: func() {
                db.Callback().Query().Remove("force_db_err_favorited")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            expectedIDs:    []uint{news1.ID}, // user只收藏了news1
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("GET", "/users/favorited", nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTUser(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                // 返回 { "news_ids": [] }
                ids, ok := resp["news_ids"].([]interface{})
                assert.True(t, ok)
                var got []uint
                for _, idVal := range ids {
                    got = append(got, uint(idVal.(float64)))
                }
                assert.Equal(t, tc.expectedIDs, got)
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestGetMyLikedNews(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db, utils.UtilsImpl{})

    user := models.User{
        OpenID:   "OpenID_LikedNews_Test",
        Nickname: "LikeTester",
    }
    db.Create(&user)

    // 创建一些新闻
    news1 := models.News{Title: "Liked News1"}
    db.Create(&news1)
    news2 := models.News{Title: "Liked News2"}
    db.Create(&news2)

    // 让 user 点赞 news2
    db.Model(&user).Association("LikedNews").Append(&news2)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "User Not Found",
            userID:         99999,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "User not found",
        },
        {
            name:           "DB Error",
            userID:         user.ID,
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_liked_db_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced liked db error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to query database",
        },
        {
            name:           "Success LikedNews",
            userID:         user.ID,
            setupFunc: func() {
                db.Callback().Query().Remove("force_liked_db_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            expectedIDs:    []uint{news2.ID},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("GET", "/users/liked", nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTUser(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                // news_ids => [news2.ID]
                ids, ok := resp["news_ids"].([]interface{})
                assert.True(t, ok)
                var got []uint
                for _, idVal := range ids {
                    got = append(got, uint(idVal.(float64)))
                }
                assert.Equal(t, tc.expectedIDs, got)
            } else if tc.expectedError != "" {
                assert.Equal(t, tc.expectedError, resp["error"])
            }
        })
    }
}

func TestGetMyViewedNews(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db, utils.UtilsImpl{})

    user := models.User{
        OpenID:   "OpenID_ViewedNews_Test",
        Nickname: "ViewTester",
    }
    db.Create(&user)

    news1 := models.News{Title: "Viewed News1"}
    db.Create(&news1)
    news2 := models.News{Title: "Viewed News2"}
    db.Create(&news2)

    // user.ViewedNews => [news1, news2]
    db.Model(&user).Association("ViewedNews").Append(&news1, &news2)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "User Not Found",
            userID:         99999,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "User not found",
        },
        {
            name:           "DB Error",
            userID:         user.ID,
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_viewed_db_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced viewed db error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to query database",
        },
        {
            name:           "Success ViewedNews",
            userID:         user.ID,
            setupFunc: func() {
                db.Callback().Query().Remove("force_viewed_db_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            // user.ViewedNews => news1, news2
            expectedIDs: []uint{news1.ID, news2.ID},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("GET", "/users/viewed", nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTUser(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                ids, ok := resp["news_ids"].([]interface{})
                assert.True(t, ok)
                var got []uint
                for _, idVal := range ids {
                    got = append(got, uint(idVal.(float64)))
                }
                assert.Equal(t, tc.expectedIDs, got)
            } else if tc.expectedError != "" {
                assert.Equal(t, tc.expectedError, resp["error"])
            }
        })
    }
}

func TestGetUserProfile(t *testing.T) {
    db := setupUserTestDB()
    router := setupUserRouter(db, utils.UtilsImpl{})

    // 创建用户
    user := models.User{
        OpenID:    "OpenID_GetUserProfile",
        Nickname:  "ProfileTester",
        AvatarURL: "avatars/some_avatar.jpg",
    }
    db.Create(&user)

    // 创建用户的新闻
    n1 := models.News{Title: "Profile News1", AuthorID: user.ID}
    db.Create(&n1)
    n2 := models.News{Title: "Profile News2", AuthorID: user.ID}
    db.Create(&n2)

    // 创建另一个用户
    otherUser := models.User{
        OpenID:   "OpenID_OtherUser",
        Nickname: "NotProfile",
    }
    db.Create(&otherUser)

    tests := []struct {
        name           string
        paramID        string // 路径 /users/:id/profile 中的 id
        userID         uint   // 模拟的JWT userID
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedNews   []uint
    }{
        {
            name:           "Unauthorized",
            paramID:        fmt.Sprintf("%d", user.ID),
            userID:         0, // 不带有效 token
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid User ID",
            paramID:        "abc", // 无法转换成 int
            userID:         user.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid user ID",
        },
        {
            name:           "User Not Found",
            paramID:        "99999", // 数据库里无此用户
            userID:         user.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "User not found",
        },
        {
            name:           "Failed To Fetch User Data (simulate DB error)",
            paramID:        fmt.Sprintf("%d", user.ID),
            userID:         user.ID,
            setupFunc: func() {
                // 在 Model(&user).First(...) 时注入错误
                db.Callback().Query().Before("gorm:query").Register("force_fetch_user_data_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced fetch user data error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch user data",
        },
        {
            name:           "Failed To Fetch User's News",
            paramID:        fmt.Sprintf("%d", user.ID),
            userID:         user.ID,
            setupFunc: func() {
                db.Callback().Query().Remove("force_fetch_user_data_err")
                // 在查 user 写的新闻时注入错误 => Model(&models.News{})
                db.Callback().Query().Before("gorm:query").Register("force_fetch_user_news_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced fetch user's news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch user's news",
        },
        {
            name: "Success Get Profile With News",
            paramID: fmt.Sprintf("%d", user.ID),
            userID:  user.ID,
            setupFunc: func() {
                db.Callback().Query().Remove("force_fetch_user_news_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            expectedNews:   []uint{n1.ID, n2.ID}, // user 写的新闻
        },
        {
            name: "Success Get Profile With No News",
            paramID: fmt.Sprintf("%d", otherUser.ID),
            userID:  user.ID,
            setupFunc: func() {
                // otherUser 并无写新闻
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            expectedNews:   []uint(nil),
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/users/" + tc.paramID + "/profile"
            req, _ := http.NewRequest("GET", url, nil)

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTUser(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                // { "nickname":"...", "avatar_url":"...", "news":[ {...} ] }
                newsList, ok := resp["news"].([]interface{})
                assert.True(t, ok)
                var gotIDs []uint
                for _, item := range newsList {
                    m := item.(map[string]interface{})
                    gotIDs = append(gotIDs, uint(m["id"].(float64)))
                }
                assert.Equal(t, tc.expectedNews, gotIDs)
            } else if tc.expectedError != "" {
                assert.Equal(t, tc.expectedError, resp["error"])
            }
        })
    }
}