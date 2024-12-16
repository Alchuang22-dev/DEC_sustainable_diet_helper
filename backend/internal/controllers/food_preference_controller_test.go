package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFoodPreferenceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the required tables
	err = db.AutoMigrate(&models.User{}, &models.FoodPreference{}, &models.DislikedFoodPreference{}, &models.Food{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Add test food data
	foods := []models.Food{
		{ZhFoodName: "猪肉", EnFoodName: "pork"},
		{ZhFoodName: "白菜", EnFoodName: "Chinese cabbage"},
	}
	db.Create(&foods)

	return db
}

func setupFoodPreferenceTestRouter(db *gorm.DB) (*gin.Engine, *FoodPreferenceController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	fpc := &FoodPreferenceController{DB: db}

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
	projectRoot := getProjectRoot()
    filePath := filepath.Join(projectRoot, "data", "food_preference", "foodPreferences-test.json")
	err = os.WriteFile(filePath, preferencesJSON, 0644)
	if err != nil {
		panic(err)
	}

	return router, fpc
}

func createFoodPreferenceTestUser(db *gorm.DB) *models.User {
	user := &models.User{
		ID:       1,
		Nickname: "TestUser",
		OpenID:   "test_open_id",
	}
	db.Create(user)
	return user
}

func TestAddFoodPreference(t *testing.T) {
	db := setupFoodPreferenceTestDB(t)
	router, fpc := setupFoodPreferenceTestRouter(db)
	user := createFoodPreferenceTestUser(db)

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
	db := setupFoodPreferenceTestDB(t)
	router, fpc := setupFoodPreferenceTestRouter(db)
	user := createFoodPreferenceTestUser(db)

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
    db := setupFoodPreferenceTestDB(t)
	router, fpc := setupFoodPreferenceTestRouter(db)
	user := createFoodPreferenceTestUser(db)

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

func TestAddDislikedFoodPreference(t *testing.T) {
	db := setupFoodPreferenceTestDB(t)
	router, fpc := setupFoodPreferenceTestRouter(db)
	user := createFoodPreferenceTestUser(db)

	// Setup the test endpoint
	router.POST("/disliked_preferences", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		fpc.AddDislikedFoodPreference(c)
	})

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Valid disliked preference",
			requestBody: map[string]interface{}{
				"food_id": 1,
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Disliked food preference added successfully",
			},
		},
		{
			name: "Duplicate disliked preference",
			requestBody: map[string]interface{}{
				"food_id": 1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Disliked preference already exists",
			},
		},
		{
			name: "Invalid food ID",
			requestBody: map[string]interface{}{
				"food_id": 999,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to add disliked preference",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/disliked_preferences", bytes.NewBuffer(jsonBody))
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

func TestDeleteDislikedFoodPreference(t *testing.T) {
	db := setupFoodPreferenceTestDB(t)
	router, fpc := setupFoodPreferenceTestRouter(db)
	user := createFoodPreferenceTestUser(db)

	// Create a test disliked preference
	dislikedPreference := models.DislikedFoodPreference{
		UserID: user.ID,
		FoodID: 1,
	}
	db.Create(&dislikedPreference)

	// Setup the test endpoint
	router.DELETE("/disliked_preferences", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		fpc.DeleteDislikedFoodPreference(c)
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
				"food_id": 1,
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Disliked food preference deleted successfully",
			},
		},
		{
			name: "Non-existent disliked preference",
			requestBody: map[string]interface{}{
				"food_id": 1,
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Disliked preference not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodDelete, "/disliked_preferences", bytes.NewBuffer(jsonBody))
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

func TestGetUserDislikedPreferences(t *testing.T) {
	db := setupFoodPreferenceTestDB(t)
	router, fpc := setupFoodPreferenceTestRouter(db)
	user := createFoodPreferenceTestUser(db)

	// Create some disliked preferences
	dislikedPreferences := []models.DislikedFoodPreference{
		{UserID: user.ID, FoodID: 1},
		{UserID: user.ID, FoodID: 2},
	}
	db.Create(&dislikedPreferences)

	// Setup the test endpoint
	router.GET("/disliked_preferences", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		fpc.GetUserDislikedPreferences(c)
	})

	// Send test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/disliked_preferences", nil)
	router.ServeHTTP(w, req)

	// Validate response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string][]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response content
	dislikedFoods := response["disliked_foods"]
	assert.Equal(t, 2, len(dislikedFoods))

	// Create expected foods mapping
	expectedFoods := []map[string]interface{}{
		{"id": float64(1), "name": "猪肉"},
		{"id": float64(2), "name": "白菜"},
	}

	// Validate returned data contains correct ID and name
	for _, food := range dislikedFoods {
		found := false
		for _, expected := range expectedFoods {
			if food["id"] == expected["id"] && food["name"] == expected["name"] {
				found = true
				break
			}
		}
		assert.True(t, found, "未找到匹配的食物: %v", food)
	}
}

func TestFoodPreferenceUnauthorizedAccess(t *testing.T) {
    db := setupFoodPreferenceTestDB(t)
    router, fpc := setupFoodPreferenceTestRouter(db)

    // 测试所有需要授权的端点
    tests := []struct {
        name       string
        method     string
        path       string
        body       map[string]interface{}
    }{
        {
            name:   "未授权添加食物偏好",
            method: http.MethodPost,
            path:   "/preferences",
            body:   map[string]interface{}{"preference_name": "highProtein"},
        },
        {
            name:   "未授权删除食物偏好",
            method: http.MethodDelete,
            path:   "/preferences",
            body:   map[string]interface{}{"preference_name": "highProtein"},
        },
        {
            name:   "未授权获取食物偏好",
            method: http.MethodGet,
            path:   "/preferences",
            body:   nil,
        },
        {
            name:   "未授权添加不喜欢的食物",
            method: http.MethodPost,
            path:   "/disliked_preferences",
            body:   map[string]interface{}{"food_id": 1},
        },
        {
            name:   "未授权删除不喜欢的食物",
            method: http.MethodDelete,
            path:   "/disliked_preferences",
            body:   map[string]interface{}{"food_id": 1},
        },
        {
            name:   "未授权获取不喜欢的食物列表",
            method: http.MethodGet,
            path:   "/disliked_preferences",
            body:   nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 设置路由处理函数
            switch tt.path {
            case "/preferences":
                switch tt.method {
                case http.MethodPost:
                    router.POST(tt.path, fpc.AddFoodPreference)
                case http.MethodDelete:
                    router.DELETE(tt.path, fpc.DeleteFoodPreference)
                case http.MethodGet:
                    router.GET(tt.path, fpc.GetUserPreferences)
                }
            case "/disliked_preferences":
                switch tt.method {
                case http.MethodPost:
                    router.POST(tt.path, fpc.AddDislikedFoodPreference)
                case http.MethodDelete:
                    router.DELETE(tt.path, fpc.DeleteDislikedFoodPreference)
                case http.MethodGet:
                    router.GET(tt.path, fpc.GetUserDislikedPreferences)
                }
            }

            // 创建请求
            var req *http.Request
            if tt.body != nil {
                jsonBody, _ := json.Marshal(tt.body)
                req, _ = http.NewRequest(tt.method, tt.path, bytes.NewBuffer(jsonBody))
                req.Header.Set("Content-Type", "application/json")
            } else {
                req, _ = http.NewRequest(tt.method, tt.path, nil)
            }

            // 发送请求并记录响应
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 验证响应
            assert.Equal(t, http.StatusUnauthorized, w.Code)
            
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, "Unauthorized", response["error"])
        })
    }
}	