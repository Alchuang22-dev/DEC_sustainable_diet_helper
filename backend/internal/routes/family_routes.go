// routes/user_routes.go
package routes

import (
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/controllers"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"

    "github.com/gin-gonic/gin"
)

func RegisterFamilyRoutes(router *gin.Engine, db *gorm.DB) {
    familyController := controllers.NewFamilyController(db)

    familyGroup := router.Group("/families")
    {
        // 需要认证的路由
        authGroup := familyGroup.Group("")
        authGroup.Use(middleware.AuthMiddleware())
        {
            // 创建家庭
            authGroup.POST("/create", familyController.CreateFamily)
            // 查看自己的家庭的信息
            authGroup.GET("/details", familyController.FamilyDetails)
            // 查看搜索家庭结果
            authGroup.GET("/search", familyController.SearchFamily)
            // 发送加入家庭请求
            authGroup.POST("/:id/join", familyController.JoinFamily)
            // 批准加入家庭
            authGroup.POST("/admit", familyController.AdmitJoinFamily)
            // 拒绝加入家庭
            authGroup.POST("/reject", familyController.RejectJoinFamily)
            // 取消加入家庭
            authGroup.DELETE("/cancel_join", familyController.CancelJoinFamily)
            // 查看当前试图加入的家庭信息
            authGroup.GET("/pending_family_details", familyController.PendingFamilyDetails)
            // 更改某个家庭成员为 member
            authGroup.PUT("/set_member", familyController.SetMember)
            // 更改某个家庭成员为 admin
            authGroup.PUT("/set_admin", familyController.SetAdmin)
            // 退出家庭
            authGroup.DELETE("/leave_family", familyController.LeaveFamily)
            // 踢出家庭
            authGroup.DELETE("/delete_family_member", familyController.DeleteFamilyMember)
            // 解散家庭
            authGroup.DELETE("/break", familyController.BreakFamily)

            authGroup.POST("/add_desired_dish", familyController.AddDesiredDish)
            authGroup.GET("/desired_dishes", familyController.GetDesiredDishes)
            authGroup.DELETE("/desired_dishes", familyController.DeleteDesiredDish)
        }
    }
}