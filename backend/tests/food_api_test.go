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

    // 在清空表之前，先自动创建表结构
    if err := db.AutoMigrate(
        &models.User{},
        &models.Food{},
    ); err != nil {
        return nil, fmt.Errorf("自动迁移表结构失败: %v", err)
    }

    // 先清空关联表，再清空主表
    if err := db.Exec("DELETE FROM food_recipes").Error; err != nil {
        return nil, fmt.Errorf("清空关联表失败: %v", err)
    }

    // 清空 foods 表
    if err := db.Exec("DELETE FROM foods").Error; err != nil {
        return nil, fmt.Errorf("清空食物表失败: %v", err)
    }

    // 清空 users 表
    if err := db.Exec("DELETE FROM users").Error; err != nil {
        return nil, fmt.Errorf("清空用户表失败: %v", err)
    }

    // 重置自增ID
    if err := db.Exec("ALTER TABLE foods AUTO_INCREMENT = 1").Error; err != nil {
        return nil, fmt.Errorf("重置自增ID失败: %v", err)
    }

    // 创建测试用户
    testUser := &models.User{
        ID:       1,
        Nickname: "TestUser",
        OpenID:   "test_open_id",
    }
    if err := db.Create(testUser).Error; err != nil {
        return nil, fmt.Errorf("创建测试用户失败: %v", err)
    }

    // 插入测试食物数据
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

    // 添加测试用的认证中间件
    router.Use(func(c *gin.Context) {
        // 模拟认证中间件，设置测试用户ID
        c.Set("user_id", uint(1))
        c.Next()
    })

    // 注册路由
    routes.RegisterFoodRoutes(router, db)

    return router
}

// TestFoodNamesAPI 测试获取食物名称列表的 API
func TestFoodNamesAPI(t *testing.T) {
    db, err := setupTestDB()
    if err != nil {
        t.Fatalf("设置测试数据库失败: %v", err)
    }

    router := setupTestRouter(db)

    tests := []struct {
        name           string
        language      string
        expectedCode   int
        checkResponse  func(*testing.T, *httptest.ResponseRecorder)
    }{
        {
            name:          "获取名称列表成功",
            language:     "zh",
            expectedCode:  200,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response []models.FoodNameResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                
                // 验证返回列表不为空
                assert.NotEmpty(t, response)
                
                // 验证第一个食物的数据完整性
                firstFood := response[0]
                assert.NotZero(t, firstFood.ID)
                assert.NotEmpty(t, firstFood.Name)
                
                // 验证所有条目都有完整的数据
                for _, food := range response {
                    assert.NotZero(t, food.ID)
                    assert.NotEmpty(t, food.Name)
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/foods/names?lang="+tt.language, nil)
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedCode, w.Code)
            if tt.checkResponse != nil {
                tt.checkResponse(t, w)
            }
        })
    }
}

// 测试函数
func TestCalculateFoodNutritionAndEmissionAPI(t *testing.T) {
    // 设置测试数据库
    db, err := setupTestDB()
    if err != nil {
        t.Fatalf("设置测试数据库失败: %v", err)
    }

    router := setupTestRouter(db)

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
                var response []models.FoodCalculateResult
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Len(t, response, 2)

                // 验证返回的数据结构完整性
                firstResult := response[0]
                assert.NotZero(t, firstResult.ID)
                assert.NotZero(t, firstResult.Emission)
                assert.NotZero(t, firstResult.Calories)
                assert.NotZero(t, firstResult.Protein)
                assert.NotZero(t, firstResult.Fat)
                assert.NotZero(t, firstResult.Carbohydrates)
                assert.NotZero(t, firstResult.Sodium)

                // 使用 InDelta 验证具体数值（允许0.1的误差）
                const delta = 0.1
                assert.InDelta(t, 0.43, response[0].Emission, delta)
                assert.InDelta(t, 26.0, response[0].Calories, delta)
            },
        },
        {
            name: "无效的食物ID",
            payload: `[{"id": 999, "price": 5.0, "weight": 0.5}]`,
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
            payload: `[{"id": 1, "price": 5.0, "weight": -1}]`,
            expectedCode: 400,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "Invalid weight for food ID 1: weight must be positive")
            },
        },
        {
            name: "无效的价格值",
            payload: `[{"id": 1, "price": -5.0, "weight": 0.5}]`,
            expectedCode: 400,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "Invalid price for food ID 1: price must be positive")
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
