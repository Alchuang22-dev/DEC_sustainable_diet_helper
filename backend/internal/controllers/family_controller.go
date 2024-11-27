// internal/controllers/family_controller.go
package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	"fmt"
)

type FamilyController struct {
    DB *gorm.DB
}

func NewFamilyController(db *gorm.DB) *FoodController {
    return &FoodController{DB: db}
}


// TODO
// 创建家庭
func (nc *NewsController) CreateFamily(c *gin.Context) {
    
}