package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
)

// 测试数据库配置
const (
    TestDBHost     = "localhost"
    TestDBPort     = "3306"
    TestDBUser     = "test_user"
    TestDBPassword = "test_password"
    TestDBName     = "test_sustainable_diet"
)

// setupFoodTestDB 设置测试数据库
func setupFoodTestDB(t *testing.T) *gorm.DB {
    // 使用 SQLite 内存数据库替代 MySQL
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("连接测试数据库失败: %v", err)
    }

    // 自动迁移所需的表
    err = db.AutoMigrate(
        &models.Food{},
    )
    if err != nil {
        t.Fatalf("迁移测试数据库失败: %v", err)
    }
    // 添加测试食材数据
    testFoods := []models.Food{
        {
            Model:         gorm.Model{ID: 1},
            ZhFoodName:    "苹果",
            EnFoodName:    "Apple",
            ImageUrl:      "http://example.com/apple.jpg",
            Calories:      52,
            Protein:       0.3,
            Fat:          0.2,
            Carbohydrates: 14,
            Sodium:        1,
            GHG:      0.43,
            Price:      1.0,
        },
        {
            Model:         gorm.Model{ID: 2},
            ZhFoodName:    "香蕉",
            EnFoodName:    "Banana",
            ImageUrl:      "http://example.com/banana.jpg",
            Calories:      89,
            Protein:       1.1,
            Fat:          0.3,
            Carbohydrates: 23,
            Sodium:        1,
            GHG:      0.86,
            Price:      1.0,
        },
    }
    for _, food := range testFoods {
        if err := db.Create(&food).Error; err != nil {
            t.Fatalf("添加测试食材失败: %v", err)
        }
    }

    return db
}

// setupFoodTestRouter 设置测试路由
func setupFoodTestRouter(db *gorm.DB) *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.Use(gin.Recovery())

    foodController := NewFoodController(db)

    // 注册路由
    foodGroup := router.Group("/foods")
    {
        // 需要认证的路由
        authGroup := foodGroup.Group("")
        // 使用测试用的认证中间件
        authGroup.Use(func(c *gin.Context) {
            // 模拟认证中间件，直接设置用户ID
            c.Set("user_id", uint(1))
            c.Next()
        })
        {   
            authGroup.GET("/names", foodController.GetFoodNames)
            authGroup.POST("/calculate", foodController.CalculateNutritionAndEmission)
        }
    }
    

    return router
}

// TestFoodNamesAPI 测试获取食物名称列表的 API
func TestFoodNamesAPI(t *testing.T) {
    db := setupFoodTestDB(t)

    router := setupFoodTestRouter(db)

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
                var response []models.FoodInfoResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                
                // 验证返回列表不为空
                assert.NotEmpty(t, response)
                
                // 验证第一个食物的数据完整性
                firstFood := response[0]
                assert.NotZero(t, firstFood.ID)
                assert.NotEmpty(t, firstFood.Name)
                assert.NotEmpty(t, firstFood.ImageUrl)
                
                // 验证所有条目都有完整的数据
                for _, food := range response {
                    assert.NotZero(t, food.ID)
                    assert.NotEmpty(t, food.Name)
                    assert.NotEmpty(t, food.ImageUrl)
                }
            },
        },
        {
            name:          "获取英文名称列表成功",
            language:     "en",
            expectedCode:  200,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response []models.FoodInfoResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                
                // 验证返回列表不为空
                assert.NotEmpty(t, response)
                
                // 验证第一个食物的数据完整性
                firstFood := response[0]
                assert.NotZero(t, firstFood.ID)
                assert.NotEmpty(t, firstFood.Name)
                assert.NotEmpty(t, firstFood.ImageUrl)
                
                // 验证返回的是英文名称
                assert.Contains(t, []string{"Apple", "Banana"}, firstFood.Name)
            },
        },
        {
            name:          "无效的语言参数",
            language:     "fr", // 使用不支持的语言
            expectedCode:  400,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.Nil(t, err)
                assert.Contains(t, response["error"], "Invalid language parameter")
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/foods/names?lang="+tt.language, nil)
            // 添加模拟的认证信息
            req.Header.Set("Authorization", "Bearer test_token")
            // 添加用户信息到上下文
            ctx := &gin.Context{Request: req}
            ctx.Set("user_id", uint(1))  // 使用测试用户ID
            
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
    db := setupFoodTestDB(t)
    gin.SetMode(gin.DebugMode)
    log.SetOutput(os.Stdout)
    router := setupFoodTestRouter(db)

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
                t.Logf("response: %v", response)
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

                // 使用 Equal 验证具体数值（因为现在是精确到1位小数）
                // GHG * weight * price / basePrice = 0.43 * 0.5 * 5.0 / 1.0 = 1.1
                assert.Equal(t, 1.1, response[0].Emission)
                assert.Equal(t, 26.0, response[0].Calories)
                assert.Equal(t, 0.2, response[0].Protein)
                assert.Equal(t, 0.1, response[0].Fat)
                assert.Equal(t, 7.0, response[0].Carbohydrates)
                assert.Equal(t, 0.5, response[0].Sodium)
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
