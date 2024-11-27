// internal/controllers/family_controller.go
package controllers

import (
    // "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    // "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
	// "fmt"
)

type FamilyController struct {
    DB *gorm.DB
}

func NewFamilyController(db *gorm.DB) *FamilyController {
    return &FamilyController{DB: db}
}


// TODO
// 创建家庭
func (nc *FamilyController) CreateFamily(c *gin.Context) {
    
}

// 查看自己的家庭的信息
func (nc *FamilyController) FamilyDetails(c *gin.Context) {
    
}

// 查看搜索家庭结果
func (nc *FamilyController) SearchFamily(c *gin.Context) {
    
}

// 发送加入家庭请求
func (nc *FamilyController) JoinFamily(c *gin.Context) {
    
}

// 批准加入家庭
func (nc *FamilyController) AdmitJoinFamily(c *gin.Context) {
    
}

