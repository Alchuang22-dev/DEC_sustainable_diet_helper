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

// 设置测试数据库
func setupNutritionCarbonTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移所需的表
	err = db.AutoMigrate(
		&models.User{},
		&models.NutritionGoal{},
		&models.CarbonGoal{},
		&models.NutritionIntake{},
		&models.CarbonIntake{},
		&models.Family{},
	)
	if err != nil {
		t.Fatalf("迁移测试数据库失败: %v", err)
	}

	return db
}

// 设置测试路由
func setupNutritionCarbonTestRouter(db *gorm.DB) (*gin.Engine, *controllers.NutritionCarbonController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	nc := &controllers.NutritionCarbonController{DB: db}
	return router, nc
}

// 创建测试用户
func createNutritionCarbonTestUser(db *gorm.DB) *models.User {
	user := &models.User{
		ID:       1,
		Nickname: "TestUser",
		OpenID:   "test_open_id",
	}
	db.Create(user)
	return user
}

// 创建测试家庭
func createTestFamily(db *gorm.DB, name string) *models.Family {
    family := &models.Family{
        Name:  name,
        Token: "test_token_" + name,
    }
    db.Create(family)
    return family
}

// 创建测试用户并关联到家庭
func createTestFamilyMember(db *gorm.DB, nickname string, familyID uint) *models.User {
    user := &models.User{
        Nickname: nickname,
        OpenID:   "test_open_id_" + nickname,
        FamilyID: &familyID,
    }
    db.Create(user)

    // 建立家庭成员关系
    var family models.Family
    db.First(&family, familyID)
    
    // 使用 Association 建立多对多关系
    db.Model(&family).Association("Members").Append(user)
    
    return user
}
// 测试设置营养目标
func TestSetNutritionGoals(t *testing.T) {
	db := setupNutritionCarbonTestDB(t)
	router, nc := setupNutritionCarbonTestRouter(db)
	user := createNutritionCarbonTestUser(db)

	router.POST("/nutrition/goals", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		nc.SetNutritionGoals(c)
	})

	tests := []struct {
		name           string
		requestBody    []controllers.NutritionGoalRequest
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "有效的营养目标",
			requestBody: []controllers.NutritionGoalRequest{
				{
					Date:          time.Now(),
					Calories:      2000,
					Protein:       60,
					Fat:           70,
					Carbohydrates: 250,
					Sodium:        2000,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "目标设置成功",
			},
		},
		{
			name: "无效的营养值",
			requestBody: []controllers.NutritionGoalRequest{
				{
					Date:          time.Now(),
					Calories:      -100,
					Protein:       -10,
					Fat:           -20,
					Carbohydrates: -30,
					Sodium:        -40,
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "无效的请求数据",
			},
		},
		{
			name: "更新已存在的营养目标",
			requestBody: []controllers.NutritionGoalRequest{
				{
					Date:          time.Now(),
					Calories:      2500,
					Protein:       75,
					Fat:           85,
					Carbohydrates: 300,
					Sodium:        2200,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "目标设置成功",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/nutrition/goals", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}
		})
	}
}

// 测试获取营养目标
func TestGetNutritionGoals(t *testing.T) {
	db := setupNutritionCarbonTestDB(t)
	router, nc := setupNutritionCarbonTestRouter(db)
	user := createNutritionCarbonTestUser(db)

	// 创建测试数据
	goal := models.NutritionGoal{
		UserID:        user.ID,
		Date:          time.Now(),
		Calories:      2000,
		Protein:       60,
		Fat:           70,
		Carbohydrates: 250,
		Sodium:        2000,
	}
	db.Create(&goal)

	router.GET("/nutrition/goals", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		nc.GetNutritionGoals(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/nutrition/goals", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, data)
}

// 测试设置碳排放目标
func TestSetCarbonGoals(t *testing.T) {
	db := setupNutritionCarbonTestDB(t)
	router, nc := setupNutritionCarbonTestRouter(db)
	user := createNutritionCarbonTestUser(db)

	router.POST("/carbon/goals", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		nc.SetCarbonGoals(c)
	})

	tests := []struct {
		name           string
		requestBody    []controllers.CarbonGoalRequest
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "有效的碳排放目标",
			requestBody: []controllers.CarbonGoalRequest{
				{
					Date:     time.Now(),
					Emission: 10.5,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "目标设置成功",
			},
		},
		{
			name: "无效的碳排放值",
			requestBody: []controllers.CarbonGoalRequest{
				{
					Date:     time.Now(),
					Emission: -1,
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "无效的请求数据",
			},
		},
		{
			name: "更新已存在的碳排放目标",
			requestBody: []controllers.CarbonGoalRequest{
				{
					Date:     time.Now(),
					Emission: 12.5,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "目标设置成功",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/carbon/goals", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}
		})
	}
}

// 测试获取实际营养摄入
func TestGetActualNutrition(t *testing.T) {
	db := setupNutritionCarbonTestDB(t)
	router, nc := setupNutritionCarbonTestRouter(db)
	user := createNutritionCarbonTestUser(db)

	// 创建测试数据
	intake := models.NutritionIntake{
		UserID:        user.ID,
		Date:          time.Now(),
		MealType:      models.Breakfast,
		Calories:      500,
		Protein:       15,
		Fat:           20,
		Carbohydrates: 60,
		Sodium:        500,
	}
	db.Create(&intake)

	router.GET("/nutrition/intakes", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		nc.GetActualNutrition(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/nutrition/intakes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.NutritionIntake
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}

// 测试获取实际碳排放
func TestGetCarbonIntakes(t *testing.T) {
	db := setupNutritionCarbonTestDB(t)
	router, nc := setupNutritionCarbonTestRouter(db)
	user := createNutritionCarbonTestUser(db)

	// 创建测试数据
	intake := models.CarbonIntake{
		UserID:   user.ID,
		Date:     time.Now(),
		MealType: models.Breakfast,
		Emission: 2.5,
	}
	db.Create(&intake)

	router.GET("/carbon/intakes", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		nc.GetCarbonIntakes(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/carbon/intakes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.CarbonIntake
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
}

// 添加新的测试函数，专门测试更新功能
func TestUpdateNutritionGoal(t *testing.T) {
	db := setupNutritionCarbonTestDB(t)
	router, nc := setupNutritionCarbonTestRouter(db)
	user := createNutritionCarbonTestUser(db)

	today := time.Now().Truncate(24 * time.Hour)

	// 先创建一个初始目标
	initialGoal := models.NutritionGoal{
		UserID:        user.ID,
		Date:          today,
		Calories:      2000,
		Protein:       60,
		Fat:           70,
		Carbohydrates: 250,
		Sodium:        2000,
	}
	db.Create(&initialGoal)

	router.POST("/nutrition/goals", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		nc.SetNutritionGoals(c)
	})

	// 准备更新请求
	updateRequest := []controllers.NutritionGoalRequest{
		{
			Date:          today,
			Calories:      2500,
			Protein:       75,
			Fat:           85,
			Carbohydrates: 300,
			Sodium:        2200,
		},
	}

	// 发送更新请求
	jsonBody, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest(http.MethodPost, "/nutrition/goals", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证更新响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证数据库中的值已更新
	var updatedGoal models.NutritionGoal
	result := db.Where("user_id = ? AND date = ?", user.ID, today).First(&updatedGoal)
	assert.NoError(t, result.Error)
	assert.Equal(t, float64(2500), updatedGoal.Calories)
	assert.Equal(t, float64(75), updatedGoal.Protein)
	assert.Equal(t, float64(85), updatedGoal.Fat)
	assert.Equal(t, float64(300), updatedGoal.Carbohydrates)
	assert.Equal(t, float64(2200), updatedGoal.Sodium)
}

func TestUpdateCarbonGoal(t *testing.T) {
	db := setupNutritionCarbonTestDB(t)
	router, nc := setupNutritionCarbonTestRouter(db)
	user := createNutritionCarbonTestUser(db)

	today := time.Now().Truncate(24 * time.Hour)

	// 先创建一个初始目标
	initialGoal := models.CarbonGoal{
		UserID:   user.ID,
		Date:     today,
		Emission: 10.5,
	}
	db.Create(&initialGoal)

	router.POST("/carbon/goals", func(c *gin.Context) {
		c.Set("user_id", user.ID)
		nc.SetCarbonGoals(c)
	})

	// 准备更新请求
	updateRequest := []controllers.CarbonGoalRequest{
		{
			Date:     today,
			Emission: 12.5,
		},
	}

	// 发送更新请求
	jsonBody, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest(http.MethodPost, "/carbon/goals", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证更新响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证数据库中的值已更新
	var updatedGoal models.CarbonGoal
	result := db.Where("user_id = ? AND date = ?", user.ID, today).First(&updatedGoal)
	assert.NoError(t, result.Error)
	assert.Equal(t, float64(12.5), updatedGoal.Emission)
} 


func TestSetSharedNutritionCarbonIntake(t *testing.T) {
    db := setupNutritionCarbonTestDB(t)
    router, nc := setupNutritionCarbonTestRouter(db)

    // 创建测试家庭
    family1 := createTestFamily(db, "Family1")
    family2 := createTestFamily(db, "Family2")

    // 创建家庭成员
    user1 := createTestFamilyMember(db, "User1", family1.ID)
    user2 := createTestFamilyMember(db, "User2", family1.ID)
    user3 := createTestFamilyMember(db, "User3", family1.ID)
    user4 := createTestFamilyMember(db, "User4", family2.ID)

    // 验证关系是否正确建立
    var family1FromDB models.Family
    db.Preload("Members").First(&family1FromDB, family1.ID)
    if len(family1FromDB.Members) != 3 {
        t.Fatalf("家庭成员关系未正确建立，期望3个成员，实际有%d个", len(family1FromDB.Members))
    }

    router.POST("/shared/nutrition-carbon", func(c *gin.Context) {
        c.Set("user_id", user1.ID)
        nc.SetSharedNutritionCarbonIntake(c)
    })

    tests := []struct {
        name           string
        requestBody    controllers.SharedNutritionCarbonIntakeRequest
        expectedStatus int
        expectedError  string
    }{
        {
            name: "单人分摊-成功",
            requestBody: controllers.SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Breakfast,
                Calories:      1000,
                Protein:       30,
                Fat:           40,
                Carbohydrates: 120,
                Sodium:        1000,
                Emission:      5.0,
                UserShares: []controllers.UserShare{
                    {UserID: user1.ID, Ratio: 1.0},
                },
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "部分家庭成员分摊-成功",
            requestBody: controllers.SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Lunch,
                Calories:      1000,
                Protein:       30,
                Fat:           40,
                Carbohydrates: 120,
                Sodium:        1000,
                Emission:      5.0,
                UserShares: []controllers.UserShare{
                    {UserID: user1.ID, Ratio: 0.6},
                    {UserID: user2.ID, Ratio: 0.4},
                },
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "全部家庭成员分摊-成功",
            requestBody: controllers.SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []controllers.UserShare{
                    {UserID: user1.ID, Ratio: 0.4},
                    {UserID: user2.ID, Ratio: 0.3},
                    {UserID: user3.ID, Ratio: 0.3},
                },
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "包含其他家庭成员-失败",
            requestBody: controllers.SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []controllers.UserShare{
                    {UserID: user1.ID, Ratio: 0.5},
                    {UserID: user4.ID, Ratio: 0.5},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "用户 4 不属于同一个家庭",
        },
        {
            name: "比例总和不等于1-失败",
            requestBody: controllers.SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []controllers.UserShare{
                    {UserID: user1.ID, Ratio: 0.3},
                    {UserID: user2.ID, Ratio: 0.3},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "分摊比例之和必须等于1",
        },
        {
            name: "无效的比例值-失败",
            requestBody: controllers.SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []controllers.UserShare{
                    {UserID: user1.ID, Ratio: 1.2},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "无效的请求数据",
        },
        {
            name: "负数营养值-失败",
            requestBody: controllers.SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      -1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []controllers.UserShare{
                    {UserID: user1.ID, Ratio: 1.0},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "无效的请求数据",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            jsonBody, _ := json.Marshal(tt.requestBody)
            req, _ := http.NewRequest(http.MethodPost, "/shared/nutrition-carbon", bytes.NewBuffer(jsonBody))
            req.Header.Set("Content-Type", "application/json")
            
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)

            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)

            if tt.expectedError != "" {
                assert.Equal(t, tt.expectedError, response["error"])
            } else {
                // 验证数据库中是否正确创建了记录
                var nutritionIntakes []models.NutritionIntake
                var carbonIntakes []models.CarbonIntake
                
                err := db.Where("date = ? AND meal_type = ?", 
                    tt.requestBody.Date,
                    tt.requestBody.MealType).
                    Find(&nutritionIntakes).Error
                assert.NoError(t, err)
                
                err = db.Where("date = ? AND meal_type = ?", 
                    tt.requestBody.Date,
                    tt.requestBody.MealType).
                    Find(&carbonIntakes).Error
                assert.NoError(t, err)

                // 验证记录数量与分摊用户数量相同
                assert.Equal(t, len(tt.requestBody.UserShares), len(nutritionIntakes))
                assert.Equal(t, len(tt.requestBody.UserShares), len(carbonIntakes))
            }
        })
    }
}