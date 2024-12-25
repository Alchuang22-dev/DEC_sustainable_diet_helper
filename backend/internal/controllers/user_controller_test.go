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

// Helper function to generate a valid Refresh Token for testing
func generateValidRefreshTokenUser(db *gorm.DB, userID uint) string {
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

type MockUtils struct {
    GenerateAccessTokenFunc  func(userID uint) (string, error)
    GenerateRefreshTokenFunc func(userID uint) (string, error)
    CopyFileFunc             func(src, dst string) error
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

func (m *MockUtils) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
    // 不在本测试中使用，可留空或mock
    return nil, nil
}

// MockRoundTripper 用于模拟 http.Client 的 Transport，从而在测试中自定义响应
type MockRoundTripper struct {
    RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip 实现 http.RoundTripper 接口
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.RoundTripFunc(req)
}



// func TestRefreshTokenHandler(t *testing.T) {
//     db := setupUserTestDB()
//     router := setupUserRouter(db, utils.UtilsImpl{})

//     // 创建用户
//     user := models.User{
//         OpenID: "OpenID_RefreshToken_Test",
//         Nickname: "RefreshTester",
//     }
//     db.Create(&user)

//     // 生成并存储一个合法的 refresh token
//     validRefresh := generateValidRefreshTokenUser(db, user.ID)

//     tests := []struct {
//         name           string
//         requestBody    interface{}
//         setupFunc      func()
//         expectedStatus int
//         expectedError  string
//         isSuccess      bool
//     }{
//         {
//             name:           "Invalid Request Body",
//             requestBody:    "not_json",
//             setupFunc:      func() {},
//             expectedStatus: http.StatusBadRequest,
//             expectedError:  "Invalid request body",
//         },
//         {
//             name:           "Empty RefreshToken",
//             requestBody:    gin.H{},
//             setupFunc:      func() {},
//             expectedStatus: http.StatusBadRequest,
//             expectedError:  "Invalid request body",
//         },
//         {
//             name:           "ValidateToken Error",
//             requestBody:    gin.H{"refresh_token": "InvalidTokenString"},
//             setupFunc: func() {
//                 // 可以 mock utils.ValidateToken 返回错误
//             },
//             expectedStatus: http.StatusUnauthorized,
//             expectedError:  "token contains an invalid number of segments", 
//             // 这是 JWT 常见错误，可按需要改
//         },
//         {
//             name:           "Invalid Token Subject",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 // mock ValidateToken => claims.Subject 不可转换
//                 // 需要 monkey patch或自定义 UtilsInterface
//             },
//             expectedStatus: http.StatusUnauthorized,
//             expectedError:  "Invalid token subject",
//         },
//         {
//             name:           "Refresh token not found in DB",
//             requestBody:    gin.H{"refresh_token": "NotInDB"},
//             setupFunc:      func() {},
//             expectedStatus: http.StatusUnauthorized,
//             expectedError:  "Refresh token not found",
//         },
//         {
//             name:           "Refresh token is expired or revoked",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 // 将 validRefresh 对应记录 设置成 revoked 或 expiresAt<now
//                 var rt models.RefreshToken
//                 db.Where("token = ?", validRefresh).First(&rt)
//                 rt.Revoked = true
//                 db.Save(&rt)
//             },
//             expectedStatus: http.StatusUnauthorized,
//             expectedError:  "Refresh token is expired or revoked",
//         },
//         {
//             name:           "User not found for this refresh token",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 // 重新生成一个 validRefresh
//                 // 并且把 UserID 指向一个不存在的用户
//                 db.Where("token = ?", validRefresh).Delete(&models.RefreshToken{})

//                 newRT := models.RefreshToken{
//                     Token:     validRefresh,
//                     UserID:    99999, // 不存在
//                     ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
//                 }
//                 db.Create(&newRT)
//             },
//             expectedStatus: http.StatusUnauthorized,
//             expectedError:  "User not found",
//         },
//         {
//             name:           "Fail to generate new access token",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 // 重新添加一个合法记录
//                 db.Where("token = ?", validRefresh).Delete(&models.RefreshToken{})

//                 newRT := models.RefreshToken{
//                     Token:     validRefresh,
//                     UserID:    user.ID,
//                     ExpiresAt: time.Now().Add(config.JWTConfig.RefreshTokenExpiration),
//                 }
//                 db.Create(&newRT)

//                 // mock GenerateAccessToken => return error
//                 // 需要对 uc.Utils 做mock
//             },
//             expectedStatus: http.StatusInternalServerError,
//             expectedError:  "Failed to generate access token",
//         },
//         {
//             name:           "Fail to generate new refresh token",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 // mock GenerateAccessToken => success
//                 // mock GenerateRefreshToken => return error
//             },
//             expectedStatus: http.StatusInternalServerError,
//             expectedError:  "Failed to generate refresh token",
//         },
//         {
//             name:           "Failed to store new refresh token",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 // mock GenerateRefreshToken => success
//                 // mock db.Create(&newRT) => error
//                 db.Callback().Create().Before("gorm:create").Register("force_create_newRT_err", func(tx *gorm.DB) {
//                     if tx.Statement.Table == "refresh_tokens" {
//                         tx.Error = fmt.Errorf("forced create new RT error")
//                     }
//                 })
//             },
//             expectedStatus: http.StatusInternalServerError,
//             expectedError:  "Failed to store new refresh token",
//         },
//         {
//             name:           "Failed to revoke old refresh token",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 db.Callback().Create().Remove("force_create_newRT_err")
//                 db.Callback().Update().Before("gorm:update").Register("force_revoke_oldRT_err", func(tx *gorm.DB) {
//                     if tx.Statement.Table == "refresh_tokens" {
//                         tx.Error = fmt.Errorf("forced revoke RT error")
//                     }
//                 })
//             },
//             expectedStatus: http.StatusInternalServerError,
//             expectedError:  "Failed to revoke old refresh token",
//         },
//         {
//             name:           "Success Refresh Token",
//             requestBody:    gin.H{"refresh_token": validRefresh},
//             setupFunc: func() {
//                 db.Callback().Update().Remove("force_revoke_oldRT_err")
//             },
//             expectedStatus: http.StatusOK,
//             isSuccess:      true,
//         },
//     }

//     for _, tc := range tests {
//         t.Run(tc.name, func(t *testing.T) {
//             tc.setupFunc()

//             bodyBytes, _ := json.Marshal(tc.requestBody)
//             req, _ := http.NewRequest("POST", "/users/refresh", bytes.NewBuffer(bodyBytes))
//             req.Header.Set("Content-Type", "application/json")

//             w := httptest.NewRecorder()
//             router.ServeHTTP(w, req)

//             assert.Equal(t, tc.expectedStatus, w.Code)

//             var resp map[string]interface{}
//             err := json.Unmarshal(w.Body.Bytes(), &resp)
//             assert.NoError(t, err)

//             if tc.isSuccess {
//                 // 正常返回新的 access_token, refresh_token
//                 _, ok1 := resp["access_token"]
//                 _, ok2 := resp["refresh_token"]
//                 assert.True(t, ok1 && ok2)
//             } else if tc.expectedError != "" {
//                 assert.Equal(t, tc.expectedError, resp["error"])
//             }
//         })
//     }
// }