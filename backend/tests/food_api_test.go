// tests/food_api_test.go
package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "strings"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"

    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/routes"
)

// 测试数据库配置
const (
    TestDBHost     = "localhost"
    TestDBPort     = "3306"
    TestDBUser     = "test_user"
    TestDBPassword = "test_password"
    TestDBName     = "test_sustainable_diet"
)

// setupTestDB 设置测试数据库
func setupTestDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        TestDBUser, TestDBPassword, TestDBHost, TestDBPort, TestDBName)

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("连接数据库失败: %v", err)
    }

    // 先清空关联表，再清空主表
    if err := db.Exec("DELETE FROM food_recipes").Error; err != nil {
        return nil, fmt.Errorf("清空关联表失败: %v", err)
    }

    // 清空 foods 表
    if err := db.Exec("DELETE FROM foods").Error; err != nil {
        return nil, fmt.Errorf("清空食物表失败: %v", err)
    }

    // 重置自增ID
    if err := db.Exec("ALTER TABLE foods AUTO_INCREMENT = 1").Error; err != nil {
        return nil, fmt.Errorf("重置自增ID失败: %v", err)
    }

    // 插入测试数据
    testFoods := []models.Food{
        {
            ZhFoodName:    "苹果",
            EnFoodName:    "Apple",
            GHG:          0.43,
            Calories:     52,
            Protein:      0.3,
            Fat:          0.2,
            Carbohydrates: 14,
            Sodium:       1,
            Price:        2.5,
        },
        {
            ZhFoodName:    "香蕉",
            EnFoodName:    "Banana",
            GHG:          0.86,
            Calories:     89,
            Protein:      1.1,
            Fat:          0.3,
            Carbohydrates: 23,
            Sodium:       1,
            Price:        3.0,
        },
    }

    for _, food := range testFoods {
        if err := db.Create(&food).Error; err != nil {
            return nil, fmt.Errorf("插入测试数据失败: %v", err)
        }
    }

    return db, nil
}

// setupTestRouter 设置测试路由
func setupTestRouter(db *gorm.DB) *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.Use(gin.Recovery())

    // 注册路由
    routes.RegisterFoodRoutes(router, db)

    return router
}

// TestFoodNamesAPI 测试获取食物名称列表的 API
func TestFoodNamesAPI(t *testing.T) {
    // 设置测试数据库
    db, err := setupTestDB()
    if err != nil {
        t.Fatalf("设置测试数据库失败: %v", err)
    }

    // 设置测试路由
    router := setupTestRouter(db)

    // 测试用例
    tests := []struct {
        name           string
        lang          string
        expectedCode  int
        expectedLen   int
        expectedName  string
        checkResponse func(*testing.T, *httptest.ResponseRecorder)
    }{
        {
            name:          "获取中文名称",
            lang:          "zh",
            expectedCode:  200,
            expectedLen:   2,
            expectedName:  "苹果",
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response []models.FoodNameResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Len(t, response, 2)
                assert.Equal(t, "苹果", response[0].Name)
            },
        },
        {
            name:          "获取英文名称",
            lang:          "en",
            expectedCode:  200,
            expectedLen:   2,
            expectedName:  "Apple",
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response []models.FoodNameResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Len(t, response, 2)
                assert.Equal(t, "Apple", response[0].Name)
            },
        },
        {
            name:          "无效的语言参数",
            lang:          "fr",
            expectedCode:  400,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "Invalid language parameter")
            },
        },
    }

    // 运行测试用例
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/foods/names?lang="+tt.lang, nil)
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedCode, w.Code)
            if tt.checkResponse != nil {
                tt.checkResponse(t, w)
            }
        })
    }
}

// 清理测试数据
func cleanupTestDB(db *gorm.DB) error {
    return db.Exec("DELETE FROM foods").Error
}

// TestCalculateNutritionAndEmissionAPI 测试食物营养计算的 API
func TestCalculateNutritionAndEmissionAPI(t *testing.T) {
    // 设置测试数据库
    db, err := setupTestDB()
    if err != nil {
        t.Fatalf("设置测试数据库失败: %v", err)
    }

    // 验证测试数据是否正确插入
    var foods []models.Food
    if err := db.Find(&foods).Error; err != nil {
        t.Fatalf("获取测试数据失败: %v", err)
    }
    t.Logf("测试数据库中的食物数量: %d", len(foods))
    for _, food := range foods {
        t.Logf("Food ID: %d, ZhName: %s, Price: %f", food.ID, food.ZhFoodName, food.Price)
    }

    // 设置测试路由
    router := setupTestRouter(db)

    // 测试用例
    tests := []struct {
        name           string
        payload        string
        expectedCode   int
        checkResponse  func(*testing.T, *httptest.ResponseRecorder)
    }{
        {
            name: "正常计算",
            payload: `[
                {
                    "id": 1,
                    "price": 5.0,
                    "weight": 0.5
                },
                {
                    "id": 2,
                    "price": 6.0,
                    "weight": 0.3
                }
            ]`,
            expectedCode: 200,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                // 打印响应内容以便调试
                t.Logf("Response Body: %s", w.Body.String())

                var response []models.FoodCalculateResult
                err := json.Unmarshal(w.Body.Bytes(), &response)
                if err != nil {
                    t.Errorf("解析响应失败: %v", err)
                    return
                }

                if len(response) != 2 {
                    t.Errorf("Expected 2 results, got %d", len(response))
                    return
                }

                // 使用更宽松的误差范围
                const delta = 0.1

                // 验证第一个食物的计算结果
                if response[0].ID != 1 {
                    t.Errorf("Expected ID 1, got %d", response[0].ID)
                }
                // 验证计算结果时使用 InDelta
                if !assert.InDelta(t, 0.43, response[0].CO2Emission, delta) {
                    t.Errorf("CO2Emission mismatch: expected around 0.43, got %f", response[0].CO2Emission)
                }
                if !assert.InDelta(t, 26.0, response[0].Calories, delta) {
                    t.Errorf("Calories mismatch: expected 26.0, got %f", response[0].Calories)
                }
            },
        },
        {
            name: "无效的食物ID",
            payload: `[
                {
                    "id": 999,
                    "price": 5.0,
                    "weight": 0.5
                }
            ]`,
            expectedCode: 500,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "food with id 999 not found")
            },
        },
        {
            name: "无效的重量值",
            payload: `[
                {
                    "id": 1,
                    "price": 5.0,
                    "weight": -1
                }
            ]`,
            expectedCode: 400,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "weight must be positive")
            },
        },
        {
            name: "无效的价格值",
            payload: `[
                {
                    "id": 1,
                    "price": -5.0,
                    "weight": 0.5
                }
            ]`,
            expectedCode: 400,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "price must be positive")
            },
        },
        {
            name: "无效的JSON格式",
            payload: `{invalid json}`,
            expectedCode: 400,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "Invalid request format")
            },
        },
    }

    // 运行测试用例
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("POST", "/foods/calculate", 
                strings.NewReader(tt.payload))
            req.Header.Set("Content-Type", "application/json")
            router.ServeHTTP(w, req)

            t.Logf("Test '%s' response code: %d", tt.name, w.Code)
            t.Logf("Test '%s' response body: %s", tt.name, w.Body.String())

            assert.Equal(t, tt.expectedCode, w.Code)
            if tt.checkResponse != nil {
                tt.checkResponse(t, w)
            }
        })
    }
}

func TestMain(m *testing.M) {
    // 在所有测试开始前的设置
    gin.SetMode(gin.TestMode)
    
    // 运行测试
    m.Run()
}