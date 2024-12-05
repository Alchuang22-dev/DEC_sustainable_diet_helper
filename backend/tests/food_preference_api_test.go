package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the required tables
	err = db.AutoMigrate(&models.User{}, &models.FoodPreference{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupTestRouter(db *gorm.DB) (*gin.Engine, *controllers.FoodPreferenceController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	fpc := &controllers.FoodPreferenceController{DB: db}

	// Create test food preferences configuration
	err := os.MkdirAll(filepath.Join("data", "food_preference"), 0755)
	if err != nil {
		panic(err)
	}

	preferences := map[string]interface{}{
		"highProtein": true,
		"lowCH":     true,
		"vegan":       true,
	}
	preferencesJSON, _ := json.Marshal(preferences)
	err = os.WriteFile(filepath.Join("data", "food_preference", "foodPreferences.json"), preferencesJSON, 0644)
	if err != nil {
		panic(err)
	}

	return router, fpc
}

func createTestUser(db *gorm.DB) *models.User {
	user := &models.User{
		ID:       1,
		Nickname: "TestUser",
		OpenID:   "test_open_id",
	}
	db.Create(user)
	return user
}

func TestAddFoodPreference(t *testing.T) {
	db := setupTestDB(t)
	router, fpc := setupTestRouter(db)
	user := createTestUser(db)

	// Setup the test endpoint
	router.POST("/preferences", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		fpc.AddFoodPreference(c)
	})

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Valid preference",
			requestBody: map[string]interface{}{
				"preference_name": "highProtein",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Food preference added successfully",
			},
		},
		{
			name: "Invalid preference name",
			requestBody: map[string]interface{}{
				"preference_name": "invalidPreference",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid preference name",
			},
		},
		{
			name:           "Missing preference name",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/preferences", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check if response contains expected message/error
			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}
		})
	}

	// Test duplicate preference
	t.Run("Duplicate preference", func(t *testing.T) {
		jsonBody, _ := json.Marshal(map[string]interface{}{
			"preference_name": "highProtein",
		})
		req, _ := http.NewRequest(http.MethodPost, "/preferences", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Preference already exists", response["error"])
	})
}

func TestDeleteFoodPreference(t *testing.T) {
	db := setupTestDB(t)
	router, fpc := setupTestRouter(db)
	user := createTestUser(db)

	// Create a test preference
	preference := models.FoodPreference{
		UserID: user.ID,
		Name:   "highProtein",
	}
	db.Create(&preference)

	// Setup the test endpoint
	router.DELETE("/preferences", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		fpc.DeleteFoodPreference(c)
	})

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Valid deletion",
			requestBody: map[string]interface{}{
				"preference_name": "highProtein",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Food preference deleted successfully",
			},
		},
		{
			name: "Non-existent preference",
			requestBody: map[string]interface{}{
				"preference_name": "highProtein",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Preference not found",
			},
		},
		{
			name: "Invalid preference name",
			requestBody: map[string]interface{}{
				"preference_name": "invalidPreference",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid preference name",
			},
		},
		{
			name:           "Missing preference name",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodDelete, "/preferences", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check if response contains expected message/error
			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}
		})
	}
}

func TestGetUserPreferences(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
	router, fpc := setupTestRouter(db)
	user := createTestUser(db)

    // 创建一些测试偏好数据
    preferences := []models.FoodPreference{
        {
            UserID: user.ID,
            Name:   "highProtein",
        },
        {
            UserID: user.ID,
            Name:   "lowSugar",
        },
    }
    for _, pref := range preferences {
        db.Create(&pref)
    }

    // 设置测试路由
    router.GET("/preferences", func(c *gin.Context) {
        c.Set("user_id", user.ID)
        fpc.GetUserPreferences(c)
    })

    // 发送测试请求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/preferences", nil)
    router.ServeHTTP(w, req)

    // 验证响应状态码
    assert.Equal(t, http.StatusOK, w.Code)

    // 解析响应
    var response []map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    // 验证响应内容
    assert.Equal(t, 2, len(response))
    
    // 验证返回的数据格式是否正确
    for _, pref := range response {
        // 确保只包含 id 和 name 字段
        assert.Contains(t, pref, "id")
        assert.Contains(t, pref, "name")
        assert.Len(t, pref, 2)

        // 验证 name 是否为预期值之一
        name := pref["name"].(string)
        assert.Contains(t, []string{"highProtein", "lowSugar"}, name)
    }
}

func TestMain(m *testing.M) {
	// Setup
	err := os.MkdirAll(filepath.Join("data", "food_preference"), 0755)
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll(filepath.Join("data", "food_preference"))
	os.Exit(code)
}