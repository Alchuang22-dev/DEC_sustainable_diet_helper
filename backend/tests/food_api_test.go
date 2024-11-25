// tests/food_api_test.go
package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

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
    // 构建测试数据库 DSN
    dsn := TestDBUser + ":" + TestDBPassword + "@tcp(" + TestDBHost + ":" + TestDBPort + ")/" + TestDBName + "?charset=utf8mb4&parseTime=True&loc=Local"

    // 连接测试数据库
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // 自动迁移
    err = db.AutoMigrate(&models.Food{})
    if err != nil {
        return nil, err
    }

    // 清空表
    db.Exec("DELETE FROM foods")

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
            return nil, err
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

func TestMain(m *testing.M) {
    // 在所有测试开始前的设置
    gin.SetMode(gin.TestMode)
    
    // 运行测试
    m.Run()
}