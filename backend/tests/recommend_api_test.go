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

func testRecommendAPI(t *testing.T) {