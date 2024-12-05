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