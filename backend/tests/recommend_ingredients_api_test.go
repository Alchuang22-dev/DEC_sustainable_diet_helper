// backend/tests/ingredient_api_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func setupIngredientTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the required tables
	err = db.AutoMigrate(&models.User{}, &models.UserIngredientHistory{}, &models.UserIngredientPreference{}, &models.Food{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// 添加一些测试食材数据
	testFoods := []models.Food{
		{ID: 1, EnFoodName: "Test Food 1"},
		{ID: 2, EnFoodName: "Test Food 2"},
		{ID: 3, EnFoodName: "Test Food 3"},
		{ID: 4, EnFoodName: "Test Food 4"},
	}
	
	for _, food := range testFoods {
		if err := db.Create(&food).Error; err != nil {
			t.Fatalf("Failed to create test food: %v", err)
		}
	}

	return db
}

func setupIngredientTestRouter(db *gorm.DB) (*gin.Engine, *controllers.IngredientController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	ic := &controllers.IngredientController{DB: db}
	return router, ic
}

func createIngredientTestUser(db *gorm.DB) *models.User {
    user := &models.User{
        ID: 1,
        Nickname: "TestUser",
        OpenID: "test_open_id",
    }
    db.Create(user)
    return user
}

func TestRecommendIngredients(t *testing.T) {
    db := setupIngredientTestDB(t)
    router, ic := setupIngredientTestRouter(db)
    user := createIngredientTestUser(db)
    
    // 设置路由
    authorized := router.Group("/")
    authorized.Use(func(c *gin.Context) {
        c.Set("user_id", user.ID)
        c.Next()
    })
    authorized.POST("/ingredients/recommend", ic.RecommendIngredients)

    t.Run("成功推荐食材", func(t *testing.T) {
        requestBody := map[string]interface{}{
            "use_last_ingredients": true,
            "liked_ingredients": []uint{1, 2},
            "disliked_ingredients": []uint{3, 4},
        }
        bodyBytes, _ := json.Marshal(requestBody)

        req, _ := http.NewRequest("POST", "/ingredients/recommend", bytes.NewBuffer(bodyBytes))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)

        ingredients, ok := response["recommended_ingredients"].([]interface{})
        assert.True(t, ok)
        assert.NotEmpty(t, ingredients)
        assert.LessOrEqual(t, len(ingredients), 20) // 确保推荐数量不超过20个
    })

    t.Run("清理历史记录", func(t *testing.T) {
        // 创建过期的历史记录
        expiredTime := time.Now().Add(-25 * time.Hour)
        history := models.UserIngredientHistory{
            UserID: user.ID,
            IngredientID: 1,
            SelectTime: expiredTime,
        }
        db.Create(&history)

        requestBody := map[string]interface{}{
            "use_last_ingredients": false,
            "liked_ingredients": []uint{},
            "disliked_ingredients": []uint{},
        }
        bodyBytes, _ := json.Marshal(requestBody)

        req, _ := http.NewRequest("POST", "/ingredients/recommend", bytes.NewBuffer(bodyBytes))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        // 验证过期记录是否被清理
        var count int64
        db.Model(&models.UserIngredientHistory{}).
            Where("user_id = ? AND select_time < ?", user.ID, time.Now().Add(-24*time.Hour)).
            Count(&count)
        assert.Equal(t, int64(0), count)
    })

    t.Run("考虑用户偏好", func(t *testing.T) {
        // 创建用户食材偏好
        preference := models.UserIngredientPreference{
            UserID: user.ID,
            IngredientID: 1,
            IsLike: true,
            UpdateTime: time.Now(),
        }
        db.Create(&preference)

        requestBody := map[string]interface{}{
            "use_last_ingredients": true,
            "liked_ingredients": []uint{2},
            "disliked_ingredients": []uint{3},
        }
        bodyBytes, _ := json.Marshal(requestBody)

        req, _ := http.NewRequest("POST", "/ingredients/recommend", bytes.NewBuffer(bodyBytes))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        
        ingredients, ok := response["recommended_ingredients"].([]interface{})
        if !ok {
            t.Fatalf("Failed to convert recommended_ingredients to array: %v", response["recommended_ingredients"])
        }
        
        assert.NotEmpty(t, ingredients)
    })
}
