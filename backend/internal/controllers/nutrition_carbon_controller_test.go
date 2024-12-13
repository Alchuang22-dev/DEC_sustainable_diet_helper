package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
func setupNutritionCarbonTestRouter(db *gorm.DB) (*gin.Engine, *NutritionCarbonController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	nc := &NutritionCarbonController{DB: db}
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
// func calculateTimeRange() (time.Time, time.Time) {
//     now := time.Now()
//     today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
//     startDate := today.AddDate(0, 0, -6)
//     endDate := today.AddDate(0, 0, 1)
//     return startDate, endDate
// }

func TestCalculateTimeRange(t *testing.T) {
    startDate, endDate := calculateTimeRange()
    
    // 获��今天的零点时间
    now := time.Now()
    today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    
    // 验证开始时间是否是6天前的零点
    expectedStartDate := today.AddDate(0, 0, -6)
    assert.Equal(t, expectedStartDate, startDate)
    
    // 验证结束时间是否是明天的零点
    expectedEndDate := today.AddDate(0, 0, 1)
    assert.Equal(t, expectedEndDate, endDate)
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
		requestBody    []NutritionGoalRequest
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "有效的营养目标",
			requestBody: []NutritionGoalRequest{
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
			requestBody: []NutritionGoalRequest{
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
			requestBody: []NutritionGoalRequest{
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

    // 获取时间范围
    startDate, endDate := calculateTimeRange()
    
    // 添加日期输出
    t.Logf("测试时间范围: %v 到 %v", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

    // 创建测试数据：在不同时间点设置目标
    testGoals := []models.NutritionGoal{
        {
            UserID:        user.ID,
            Date:          startDate,
            Calories:      2000,
            Protein:       60,
            Fat:           70,
            Carbohydrates: 250,
            Sodium:        2000,
        },
        {
            UserID:        user.ID,
            Date:          endDate.AddDate(0, 0, -1), // 今天
            Calories:      2200,
            Protein:       65,
            Fat:           75,
            Carbohydrates: 275,
            Sodium:        2100,
        },
    }

    for _, goal := range testGoals {
        db.Create(&goal)
    }

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
    assert.Equal(t, 8, len(data)) // 确保返回8天的数据

    // 验证第一天和最后一天的数据
    firstDay := data[0].(map[string]interface{})
    assert.Equal(t, float64(2000), firstDay["calories"])
    
    lastDay := data[6].(map[string]interface{}) // 今天
    assert.Equal(t, float64(2200), lastDay["calories"])
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
		requestBody    []CarbonGoalRequest
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "有效的碳排放目标",
			requestBody: []CarbonGoalRequest{
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
			requestBody: []CarbonGoalRequest{
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
			requestBody: []CarbonGoalRequest{
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

// 测试获取碳排放目标
func TestGetCarbonGoals(t *testing.T) {
    db := setupNutritionCarbonTestDB(t)
    router, nc := setupNutritionCarbonTestRouter(db)
    user := createNutritionCarbonTestUser(db)

    // 获取时间范围
    startDate, endDate := calculateTimeRange()
    
    // 添加日期输出
    t.Logf("测试时间范围: %v 到 %v", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

    // 创建测试数据：在不同时间点设置目标
    testGoals := []models.CarbonGoal{
        {
            UserID:        user.ID,
            Date:          startDate,
            Emission:      10.5,
        },
        {
            UserID:        user.ID,
            Date:          endDate.AddDate(0, 0, -1), // 今天
            Emission:      12.5,
        },
    }

    for _, goal := range testGoals {
        db.Create(&goal)
    }

    router.GET("/carbon/goals", func(c *gin.Context) {
        c.Set("user_id", user.ID)
        nc.GetCarbonGoals(c)
    })

    w := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, "/carbon/goals", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    data, ok := response["data"].([]interface{})
    assert.True(t, ok)
    assert.Equal(t, 8, len(data)) // 确保返回8天的数据

    // 验证第一天和最后一天的数据
    firstDay := data[0].(map[string]interface{})
    assert.Equal(t, float64(10.5), firstDay["emission"])
    
    lastDay := data[6].(map[string]interface{}) // 今天
    assert.Equal(t, float64(12.5), lastDay["emission"])
}

// 测试获取实际营养摄入
func TestGetActualNutrition(t *testing.T) {
    db := setupNutritionCarbonTestDB(t)
    router, nc := setupNutritionCarbonTestRouter(db)
    user := createNutritionCarbonTestUser(db)

    // 获取时间范围
    startDate, endDate := calculateTimeRange()
    endDate = endDate.AddDate(0, 0, -1) // 调整为今天

    // 创建测试数据：在不同时间点和不同餐次记录摄入
    testIntakes := []models.NutritionIntake{
        {
            UserID:        user.ID,
            Date:          startDate,
            MealType:      models.Breakfast,
            Calories:      500,
            Protein:       15,
            Fat:           20,
            Carbohydrates: 60,
            Sodium:        500,
        },
        {
            UserID:        user.ID,
            Date:          endDate,
            MealType:      models.Dinner,
            Calories:      800,
            Protein:       25,
            Fat:           30,
            Carbohydrates: 90,
            Sodium:        700,
        },
    }

    for _, intake := range testIntakes {
        db.Create(&intake)
    }

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
    assert.Equal(t, 28, len(response)) // 7天 x 4餐

    // 验证第一天早餐的数据
    firstDayBreakfast := response[0]
    assert.Equal(t, float64(500), firstDayBreakfast.Calories)
    assert.Equal(t, models.Breakfast, firstDayBreakfast.MealType)

    // 验证最后一天��餐的数据
    lastDayDinner := response[26] // 第7天的晚餐索引：(7-1)*4 + 2 = 26
    assert.Equal(t, float64(800), lastDayDinner.Calories)
    assert.Equal(t, models.Dinner, lastDayDinner.MealType)
}

// 测试获取实际碳排放
func TestGetActualCarbon(t *testing.T) {
    db := setupNutritionCarbonTestDB(t)
    router, nc := setupNutritionCarbonTestRouter(db)
    user := createNutritionCarbonTestUser(db)

    // 获取时间范围
    startDate, endDate := calculateTimeRange()
    endDate = endDate.AddDate(0, 0, -1) // 调整为今天

    // 创建测试数据：在不同时间点记录碳排放
    testEmissions := []models.CarbonIntake{
        {
            UserID:        user.ID,
            Date:          startDate,
            MealType:      models.Breakfast,
            Emission:      5.0,
        },
        {
            UserID:        user.ID,
            Date:          endDate,
            MealType:      models.Dinner,
            Emission:      10.0,
        },
    }

    for _, emission := range testEmissions {
        db.Create(&emission)
    }

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
    assert.Equal(t, 28, len(response)) // 7天 x 4餐

    // 验证第一天早餐的数据
    firstDayBreakfast := response[0]
    assert.Equal(t, float64(5.0), firstDayBreakfast.Emission)
    assert.Equal(t, models.Breakfast, firstDayBreakfast.MealType)

    // 验证最后一天晚餐的数据
    lastDayDinner := response[26] // 第7天的晚餐索引：(7-1)*4 + 2 = 26
    assert.Equal(t, float64(10.0), lastDayDinner.Emission)
    assert.Equal(t, models.Dinner, lastDayDinner.MealType)
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
	updateRequest := []NutritionGoalRequest{
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
	updateRequest := []CarbonGoalRequest{
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
        requestBody    SharedNutritionCarbonIntakeRequest
        expectedStatus int
        expectedError  string
    }{
        {
            name: "单人分摊-成功",
            requestBody: SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Breakfast,
                Calories:      1000,
                Protein:       30,
                Fat:           40,
                Carbohydrates: 120,
                Sodium:        1000,
                Emission:      5.0,
                UserShares: []UserShare{
                    {UserID: user1.ID, Ratio: 1.0},
                },
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "部分家庭成员分摊-成功",
            requestBody: SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Lunch,
                Calories:      1000,
                Protein:       30,
                Fat:           40,
                Carbohydrates: 120,
                Sodium:        1000,
                Emission:      5.0,
                UserShares: []UserShare{
                    {UserID: user1.ID, Ratio: 0.6},
                    {UserID: user2.ID, Ratio: 0.4},
                },
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "全部家庭成员分摊-成功",
            requestBody: SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []UserShare{
                    {UserID: user1.ID, Ratio: 0.4},
                    {UserID: user2.ID, Ratio: 0.3},
                    {UserID: user3.ID, Ratio: 0.3},
                },
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "包含其他家庭成员-失败",
            requestBody: SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []UserShare{
                    {UserID: user1.ID, Ratio: 0.5},
                    {UserID: user4.ID, Ratio: 0.5},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "用户 4 不属于同一个家庭",
        },
        {
            name: "比例总和不等于1-失败",
            requestBody: SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []UserShare{
                    {UserID: user1.ID, Ratio: 0.3},
                    {UserID: user2.ID, Ratio: 0.3},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "分摊比例之和必须等于1",
        },
        {
            name: "比例值大于1-失败",
            requestBody: SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []UserShare{
                    {UserID: user1.ID, Ratio: 1.2},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "无效的请求数据",
        },
        {
            name: "负数比例值-失败",
            requestBody: SharedNutritionCarbonIntakeRequest{
                Date:          time.Now(),
                MealType:      models.Dinner,
                Calories:      1500,
                Protein:       45,
                Fat:           60,
                Carbohydrates: 180,
                Sodium:        1500,
                Emission:      7.5,
                UserShares: []UserShare{
                    {UserID: user1.ID, Ratio: -0.5},
                },
            },
            expectedStatus: http.StatusBadRequest,
            expectedError: "无效的请求数据",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 添加日期输出
            t.Logf("测试日期: %v", tt.requestBody.Date.Format("2006-01-02"))
            
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

func TestNutritionCarbonUnauthorizedAccess(t *testing.T) {
    db := setupNutritionCarbonTestDB(t)
    router, nc := setupNutritionCarbonTestRouter(db)

    // 测试所有需要授权的端点
    tests := []struct {
        name       string
        method     string
        path       string
        body       map[string]interface{}
    }{
        {
            name:   "未授权设置营养目标",
            method: http.MethodPost,
            path:   "/nutrition/goals",
            body: map[string]interface{}{
                "date":          time.Now(),
                "calories":      2000,
                "protein":       60,
                "fat":          70,
                "carbohydrates": 250,
                "sodium":        2000,
            },
        },
        {
            name:   "未授权获取营养目标",
            method: http.MethodGet,
            path:   "/nutrition/goals",
            body:   nil,
        },
        {
            name:   "未授权获取实际营养摄入",
            method: http.MethodGet,
            path:   "/nutrition/intakes",
            body:   nil,
        },
        {
            name:   "未授权设置碳排放目标",
            method: http.MethodPost,
            path:   "/carbon/goals",
            body: map[string]interface{}{
                "date":     time.Now(),
                "emission": 10.5,
            },
        },
        {
            name:   "未授权获取碳排放目标",
            method: http.MethodGet,
            path:   "/carbon/goals",
            body:   nil,
        },
        {
            name:   "未授权获取实际碳排放",
            method: http.MethodGet,
            path:   "/carbon/intakes",
            body:   nil,
        },
        {
            name:   "未授权设置共享营养碳排放",
            method: http.MethodPost,
            path:   "/shared/nutrition-carbon",
            body: map[string]interface{}{
                "date":          time.Now(),
                "meal_type":     "breakfast",
                "calories":      1000,
                "protein":       30,
                "fat":          40,
                "carbohydrates": 120,
                "sodium":        1000,
                "emission":      5.0,
                "user_shares": []map[string]interface{}{
                    {"user_id": 1, "ratio": 1.0},
                },
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 设置路由处理函数
            switch tt.path {
            case "/nutrition/goals":
                switch tt.method {
                case http.MethodPost:
                    router.POST(tt.path, nc.SetNutritionGoals)
                case http.MethodGet:
                    router.GET(tt.path, nc.GetNutritionGoals)
                }
            case "/nutrition/intakes":
                router.GET(tt.path, nc.GetActualNutrition)
            case "/carbon/goals":
                switch tt.method {
                case http.MethodPost:
                    router.POST(tt.path, nc.SetCarbonGoals)
                case http.MethodGet:
                    router.GET(tt.path, nc.GetCarbonGoals)
                }
            case "/carbon/intakes":
                router.GET(tt.path, nc.GetCarbonIntakes)
            case "/shared/nutrition-carbon":
                router.POST(tt.path, nc.SetSharedNutritionCarbonIntake)
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