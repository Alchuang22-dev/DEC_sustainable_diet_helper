// internal/controllers/food_preference_controller_test.go
package controllers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestFoodPreference(t *testing.T) (*gin.Engine, *gorm.DB) {
    // 使用 SQLite 内存数据库进行测试
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }

    // 自动迁移
    err = db.AutoMigrate(&models.User{}, &models.FoodPreference{})
    if err != nil {
        t.Fatal(err)
    }

    // 创建测试用户
    user := models.User{
        ID:       1,
        Nickname: "Test User",
    }
    db.Create(&user)

    // 设置路由
    router := gin.Default()
    fpc := &FoodPreferenceController{DB: db}

    // 添加认证中间件
    router.Use(func(c *gin.Context) {
        c.Set("user_id", uint(1))
        c.Next()
    })

    router.POST("/preferences", fpc.AddFoodPreference)
    router.DELETE("/preferences", fpc.DeleteFoodPreference)

    return router, db
}

func TestAddFoodPreference(t *testing.T) {
    router, _ := setupTestFoodPreference(t)

    // 测试添加有效的食物偏好
    requestBody := map[string]string{
        "preference_name": "highProtein",
    }
    jsonBody, _ := json.Marshal(requestBody)

    req := httptest.NewRequest("POST", "/preferences", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // 测试添加无效的食物偏好
    invalidRequestBody := map[string]string{
        "preference_name": "invalidPreference",
    }
    jsonBody, _ = json.Marshal(invalidRequestBody)

    req = httptest.NewRequest("POST", "/preferences", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteFoodPreference(t *testing.T) {
    router, db := setupTestFoodPreference(t)

    // 先添加一个食物偏好
    preference := models.FoodPreference{
        UserID: 1,
        Name:   "highProtein",
    }
    db.Create(&preference)

    // 测试删除存在的食物偏好
    requestBody := map[string]string{
        "preference_name": "highProtein",
    }
    jsonBody, _ := json.Marshal(requestBody)

    req := httptest.NewRequest("DELETE", "/preferences", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // 测试删除不存在的食物偏好
    req = httptest.NewRequest("DELETE", "/preferences", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}