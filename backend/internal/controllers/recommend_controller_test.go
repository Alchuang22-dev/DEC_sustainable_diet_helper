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
	"os"
	"path/filepath"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
    "io"
    "encoding/csv"
    "fmt"
    "strconv"
)

func stringToFloat64(s string) (float64, error) {
    if s == "" {
        return 0, nil
    }
    return strconv.ParseFloat(s, 64)
}

func importFoodsData(db *gorm.DB, filename string) error {
    // 打开CSV文件
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    // 创建CSV reader
    reader := csv.NewReader(file)
    
    // 跳过表头
    if _, err := reader.Read(); err != nil {
        return fmt.Errorf("error reading headers: %v", err)
    }
    
    // 记录成功和失败的数量
    var successful, failed int

    // 逐行读取数据
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error reading row: %v", err)
            failed++
            continue
        }

        // 转换数值字段
        ghg, err := stringToFloat64(record[1])
        calories, err1 := stringToFloat64(record[3])
        protein, err2 := stringToFloat64(record[4])
        fat, err3 := stringToFloat64(record[5])
        carbs, err4 := stringToFloat64(record[6])
        sodium, err5 := stringToFloat64(record[7])
        price, err6 := stringToFloat64(record[8])
        

        if err != nil || err1 != nil || err2 != nil || err3 != nil || 
           err4 != nil || err5 != nil || err6 != nil {
            log.Printf("Error converting number fields for food %s: %v", record[0], err)
            failed++
            continue
        }

        // 创建Food记录
        food := models.Food{
            ZhFoodName:    record[0],
            GHG:           ghg,
            EnFoodName:    strings.ToLower(record[2]),  // 将英文名称转换为小写
            Calories:      calories * 10, // 对于营养的部分，我们数据中的单位是每100g，但是我们的模型中的单位是每1kg
            Protein:       protein * 10,
            Fat:           fat * 10,
            Carbohydrates: carbs * 10,
            Sodium:        sodium * 10,
            Price:         price,
            ImageUrl:      record[9],
        }

        // 保存到数据库
        if err := food.CreateFood(db); err != nil {
            log.Printf("Error saving food %s: %v", food.ZhFoodName, err)
            failed++
            continue
        }
        successful++
    }

    log.Printf("Food import completed. Successful: %d, Failed: %d", successful, failed)
    return nil
}
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

    // 导入食物数据
    projectRoot := getProjectRoot()
    filePath := filepath.Join(projectRoot, "data", "food_dataset", "foods_dataset_url.csv")

    
    err = importFoodsData(db, filePath)
    if err != nil {
        t.Fatalf("导入食物数据失败: %v", err)
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
// 创建测试用户
func createRecommendTestUser(db *gorm.DB) *models.User {
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
func TestRecommendUnauthorizedAccess(t *testing.T) {
    db := setupRecommendTestDB(t)
    router, rc := setupRecommendTestRouter(db)

    // 注册所有需要测试的路由
    router.POST("/ingredients/recommend", rc.RecommendIngredients)
    router.POST("/recipes/recommend", rc.RecommendRecipes)
    router.POST("/ingredients/set", rc.SetUserSelectedFoods)

    tests := []struct {
        name       string
        method     string
        path       string
        body       interface{}
    }{
        {
            name:   "未授权推荐食材",
            method: "POST",
            path:   "/ingredients/recommend",
            body:   IngredientRecommendRequest{},
        },
        {
            name:   "未授权推荐菜谱",
            method: "POST",
            path:   "/recipes/recommend",
            body:   RecipeRecommendAndSetUserLastSelectedFoodsRequest{},
        },
        {
            name:   "未授权设置选择的食材",
            method: "POST",
            path:   "/ingredients/set",
            body:   RecipeRecommendAndSetUserLastSelectedFoodsRequest{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            jsonBody, _ := json.Marshal(tt.body)
            req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(jsonBody))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, http.StatusUnauthorized, w.Code)
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)
            assert.Equal(t, "用户未认证", response["error"])
        })
    }
}

func TestFoodPreference(t *testing.T) {
    db := setupRecommendTestDB(t)
    router, rc := setupRecommendTestRouter(db)

    // 创建测试用户
    user := createRecommendTestUser(db)

    // 设置路由
    router.POST("/ingredients/recommend", func(c *gin.Context) {
        c.Set("user_id", user.ID)
        rc.RecommendIngredients(c)
    })

    // 测试不同的食物偏好场景
    tests := []struct {
        name           string
        requestBody    IngredientRecommendRequest
        expectedStatus int
        setupPrefs     func(*testing.T, *gorm.DB, uint)
        checkResponse  func(*testing.T, *httptest.ResponseRecorder)
    }{
        {
            name: "高蛋白饮食偏好",
            requestBody: IngredientRecommendRequest{
                UseLastIngredients: false,
                LikedIngredients:   []uint{},
                DislikedIngredients: []uint{},
            },
            expectedStatus: http.StatusOK,
            setupPrefs: func(t *testing.T, db *gorm.DB, userID uint) {
                pref := models.FoodPreference{
                    UserID: userID,
                    Name: "highProtein",
                }
                err := db.Create(&pref).Error
                assert.NoError(t, err)
            },
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                
                recommendedFoods, ok := response["recommended_ingredients"].([]interface{})
                assert.True(t, ok)
                assert.NotEmpty(t, recommendedFoods)
                
            },
        },
        {
            name: "素食偏好",
            requestBody: IngredientRecommendRequest{
                UseLastIngredients: false,
                LikedIngredients:   []uint{},
                DislikedIngredients: []uint{},
            },
            expectedStatus: http.StatusOK,
            setupPrefs: func(t *testing.T, db *gorm.DB, userID uint) {
                pref := models.FoodPreference{
                    UserID: userID,
                    Name: "vegetarian",
                }
                err := db.Create(&pref).Error
                assert.NoError(t, err)
            },
            checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
                var response map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                
                recommendedFoods, ok := response["recommended_ingredients"].([]interface{})
                assert.True(t, ok)
                assert.NotEmpty(t, recommendedFoods)
                
                // 验证推荐的食材中不包含肉类
                meatFoods := []string{"pork", "beef", "chicken", "fish"}
                for _, food := range recommendedFoods {
                    foodMap, ok := food.(map[string]interface{})
                    if !ok {
                        t.Logf("食材类型不是map")
                        continue
                    }
                    foodName, ok := foodMap["name"].(string)
                    if !ok {
                        t.Logf("食材名称类型不是string")
                        continue
                    }
                    for _, meatFood := range meatFoods {
                        assert.NotEqual(t, foodName, meatFood, "素食推荐不应包含肉类食物")
                    }
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 设置测试偏好
            if tt.setupPrefs != nil {
                tt.setupPrefs(t, db, user.ID)
            }

            // 发送请求
            jsonBody, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest("POST", "/ingredients/recommend", bytes.NewBuffer(jsonBody))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            // 验证响应
            assert.Equal(t, tt.expectedStatus, w.Code)
            if tt.checkResponse != nil {
                tt.checkResponse(t, w)
            }
        })
    }
}

// 测试食物偏好配置文件的读取
func TestLoadFoodPreferences(t *testing.T) {
    projectRoot := getProjectRoot()
    filePath := filepath.Join(projectRoot, "data", "food_preference", "foodPreferences.json")
    
    // 读取配置文件
    data, err := os.ReadFile(filePath)
    assert.NoError(t, err, "应该能够读取食物偏好配置文件")
    
    // 解析JSON
    var preferences map[string]interface{}
    err = json.Unmarshal(data, &preferences)
    assert.NoError(t, err, "应该能够解析食物偏好JSON")
    
    // 验证必要的偏好类型存在
    expectedPrefs := []string{"highProtein", "highEnergy", "lowFat", "lowCH", "lowsodium", "vegan", "vegetarian"}
    for _, pref := range expectedPrefs {
        _, exists := preferences[pref]
        assert.True(t, exists, "应该包含%s偏好类型", pref)
    }
    
    // 验证偏好的结构
    for prefName, pref := range preferences {
        prefMap := pref.(map[string]interface{})
        assert.Contains(t, prefMap, "name_en", "%s应该包含英文名称", prefName)
        assert.Contains(t, prefMap, "name_zh-Hans", "%s应该包含中文名称", prefName)
        assert.Contains(t, prefMap, "food_pos", "%s应该包含推荐食物列表", prefName)
        assert.Contains(t, prefMap, "food_neg", "%s应该包含避免食物列表", prefName)
    }
}