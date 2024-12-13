package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"log"
	"github.com/stretchr/testify/assert"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRecommendTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移所需的表
    err = db.AutoMigrate(
        &models.User{},
        &models.Food{},
        &models.UserIngredientHistory{},    
        &models.UserRecipeHistory{},        
        &models.UserIngredientPreference{}, 
        &models.UserLastSelectedFoods{},    
        &models.Recipe{}, 
        &models.FoodPreference{},
    )
	if err != nil {
		t.Fatalf("迁移测试数据库失败: %v", err)
	}
        
    // 添加测试食材数据
    testFoods := []models.Food{
        {
            Model:      gorm.Model{ID: 1},
            ZhFoodName: "测试食材1",
            EnFoodName: "Test Food 1",
        },
        {
            Model:      gorm.Model{ID: 2},
            ZhFoodName: "测试食材2",
            EnFoodName: "Test Food 2",
        },
    }
    
    for _, food := range testFoods {
        if err := db.Create(&food).Error; err != nil {
            t.Fatalf("创建测试食材失败: %v", err)
        }
    }

	return db
}

func setupRecommendTestRouter(db *gorm.DB) (*gin.Engine, *RecommendController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rc := &RecommendController{DB: db}

	return router, rc
}

func setupRecommendTestUser(db *gorm.DB) *models.User {
	user := &models.User{
		ID:       1,
		Nickname: "TestUser",
		OpenID:   "test_open_id",
	}
	db.Create(user)
	return user
}

// 测试设置用户选择的食材
func TestSetUserSelectedFoods(t *testing.T) {
    // 设置测试数据库和路由
    db := setupRecommendTestDB(t)
    router, rc := setupRecommendTestRouter(db)
    
    // 创建测试用户
    testUser := setupRecommendTestUser(db)


    // 设置路由 - 修改为与实际路由一致
    router.POST("/ingredients/set", func(c *gin.Context) {
        c.Set("user_id", testUser.ID)
        rc.SetUserSelectedFoods(c)
    })

    // 测试用例1: 成功设置用户选择的食材
    t.Run("成功设置用户选择的食材", func(t *testing.T) {
        w := httptest.NewRecorder()
        reqBody := `{"selected_ingredients":[1,2]}`
        req, _ := http.NewRequest("POST", "/ingredients/set", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        router.ServeHTTP(w, req)

        // 验证响应状态码
        if w.Code != http.StatusOK {
            t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
        }

        // 验证数据库中是否正确保存了记录
        var selectedFoods []models.UserLastSelectedFoods
        if err := db.Where("user_id = ?", testUser.ID).Find(&selectedFoods).Error; err != nil {
            t.Fatalf("查询用户选择的食材失败: %v", err)
        }

        if len(selectedFoods) != 2 {
            t.Errorf("期望保存2条记录，实际保存了%d条", len(selectedFoods))
        }
    })

    // 测试用例2: 空的选择食材列表
    t.Run("空的选择食材列表", func(t *testing.T) {
        w := httptest.NewRecorder()
        reqBody := `{"selected_ingredients":[]}`
        req, _ := http.NewRequest("POST", "/ingredients/set", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        router.ServeHTTP(w, req)

        // 验证响应状态码
        if w.Code != http.StatusBadRequest {
            t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
        }
    })

    // 测试用例3: 无效的食材ID
    t.Run("无效的食材ID", func(t *testing.T) {
        w := httptest.NewRecorder()
        reqBody := `{"selected_ingredients":[999]}`
        req, _ := http.NewRequest("POST", "/ingredients/set", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        router.ServeHTTP(w, req)

        // 验证响应状态码
        if w.Code != http.StatusBadRequest {
            t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
        }
    })

    // 测试用例4: 无效的请求格式
    t.Run("无效的请求格式", func(t *testing.T) {
        w := httptest.NewRecorder()
        reqBody := `{invalid json}`
        req, _ := http.NewRequest("POST", "/ingredients/set", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        router.ServeHTTP(w, req)

        // 验证响应状态码
        if w.Code != http.StatusBadRequest {
            t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
        }
    })
}

// 测试推荐食材
func TestGetRecommendedFoods(t *testing.T) {
    db := setupRecommendTestDB(t)
    router, rc := setupRecommendTestRouter(db)
    user := setupRecommendTestUser(db)

    // Add middleware first
    router.Use(func(c *gin.Context) {
        // Get auth setup from header
        if authSetup := c.GetHeader("auth_setup"); authSetup == "true" {
            c.Set("user_id", user.ID)
        }
        c.Next()
    })

    // Then register the route
    router.POST("/recommend/foods", rc.RecommendIngredients)

    tests := []struct {
        name           string
        setupAuth      bool
        requestBody    interface{}
        expectedStatus int
        checkResponse  func(*testing.T, *httptest.ResponseRecorder)
    }{
        {
            name: "成功获取推荐食材",
            setupAuth: true,
            requestBody: gin.H{
                "use_last_ingredients": false,
                "liked_ingredients": []uint{},
                "disliked_ingredients": []uint{},
            },
            expectedStatus: http.StatusOK,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response IngredientRecommendResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                assert.NotEmpty(t, response.RecommendedIngredients)
                log.Printf("response: %v", response)
            },
        },
        {
            name: "未授权访问",
            setupAuth: false,
            requestBody: gin.H{},
            expectedStatus: http.StatusUnauthorized,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                assert.Equal(t, "用户未认证", response["error"])
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.requestBody)
            req, _ := http.NewRequest(http.MethodPost, "/recommend/foods", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            // Set auth setup in header
            if tt.setupAuth {
                req.Header.Set("auth_setup", "true")
            } else {
                req.Header.Set("auth_setup", "false")
            }
            
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code, "Status code mismatch")
            if tt.checkResponse != nil {
                tt.checkResponse(t, w)
            }
        })
    }
}


// 测试推荐食谱
func TestGetRecommendedRecipes(t *testing.T) {
    db := setupRecommendTestDB(t)
    router, rc := setupRecommendTestRouter(db)
    user := setupRecommendTestUser(db)

    // 创建测试食谱数据
    recipe1 := models.Recipe{
        Name: "Test Recipe 1",
        URL: "test-url-1",
        ImageURL: "test-image-1",
        Ingredients: `{"ingredient1": 100, "ingredient2": 200}`,
    }
    recipe2 := models.Recipe{
        Name: "Test Recipe 2",
        URL: "test-url-2",
        ImageURL: "test-image-2",
        Ingredients: `{"ingredient1": 150, "ingredient2": 250}`,
    }
    db.Create(&recipe1)
    db.Create(&recipe2)

    // 添加食谱和食材的关联关系
    err := db.Exec("INSERT INTO food_recipes (recipe_id, food_id) VALUES (?, ?)", recipe1.ID, 1).Error
    assert.NoError(t, err)
    err = db.Exec("INSERT INTO food_recipes (recipe_id, food_id) VALUES (?, ?)", recipe2.ID, 2).Error
    assert.NoError(t, err)

    // Add middleware first
    router.Use(func(c *gin.Context) {
        if authSetup := c.GetHeader("auth_setup"); authSetup == "true" {
            c.Set("user_id", user.ID)
        }
        c.Next()
    })

    // Register the route
    router.POST("/recommend/recipes", rc.RecommendRecipes)

    tests := []struct {
        name           string
        setupAuth      bool
        requestBody    interface{}
        expectedStatus int
        checkResponse  func(*testing.T, *httptest.ResponseRecorder)
    }{
        {
            name: "成功获取推荐食谱",
            setupAuth: true,
            requestBody: gin.H{
                "selected_ingredients": []uint{1, 2},
                "disliked_ingredients": []uint{},
            },
            expectedStatus: http.StatusOK,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response RecipeRecommendResponse
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                assert.NotEmpty(t, response.RecommendedRecipes)
                // 验证返回的食谱数据
                for _, recipe := range response.RecommendedRecipes {
                    assert.NotEmpty(t, recipe.Name)
                    assert.NotEmpty(t, recipe.ImageURL)
                    assert.NotZero(t, recipe.RecipeID)
                    assert.NotEmpty(t, recipe.Ingredients)
                }
            },
        },
        {
            name: "未授权访问",
            setupAuth: false,
            requestBody: gin.H{
                "selected_ingredients": []uint{1},
                "disliked_ingredients": []uint{},
            },
            expectedStatus: http.StatusUnauthorized,
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                assert.Equal(t, "用户未认证", response["error"])
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.requestBody)
            req, _ := http.NewRequest(http.MethodPost, "/recommend/recipes", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            // Set auth setup in header
            if tt.setupAuth {
                req.Header.Set("auth_setup", "true")
            } else {
                req.Header.Set("auth_setup", "false")
            }
            
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)
            if tt.checkResponse != nil {
                tt.checkResponse(t, w)
            }
        })
    }
}