package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
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
        &models.News{},
        &models.Resource{},
        &models.Paragraph{},
        &models.Video{},
        &models.Comment{},
        &models.Food{},
        &models.Recipe{},
        &models.Family{},
        &models.FoodPreference{},
        &models.NutritionGoal{},
        &models.CarbonGoal{},
        &models.NutritionIntake{},
        &models.CarbonIntake{},
        &models.RefreshToken{},
        &models.FamilyDish{},
		&models.UserLastSelectedFoods{},
    )
	if err != nil {
		t.Fatalf("迁移测试数据库失败: %v", err)
	}

	return db
}

func setupRecommendTestRouter(db *gorm.DB) (*gin.Engine, *controllers.RecommendController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rc := &controllers.RecommendController{DB: db}

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

func setupRecommendTestFamily(db *gorm.DB) *models.Family {
	family := &models.Family{
		Name: "TestFamily",
	}
	db.Create(family)
	return family
}


// 创建测试用户并关联到家庭
func setupRecommendTestFamilyMember(db *gorm.DB, nickname string, familyID uint) *models.User {
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

func TestSetUserSelectedFoods(t *testing.T) {
    // 设置测试数据库和路由
    db := setupRecommendTestDB(t)
    router, rc := setupRecommendTestRouter(db)
    
    // 创建测试用户
    testUser := setupRecommendTestUser(db)
    
    // 创建一些测试食材
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
