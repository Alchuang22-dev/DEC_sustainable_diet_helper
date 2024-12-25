// controllers/news_controller_test.go
package controllers

import (
    "bytes"
    "encoding/json"
    "fmt"
    // "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "os"
    // "path/filepath"
    "testing"
    "time"

    // "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/config"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/middleware"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/utils"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// setupNewsTestDB 初始化内存中的 SQLite 数据库并迁移模型
func setupNewsTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    if err := db.AutoMigrate(&models.User{}, &models.Draft{}, &models.DraftParagraph{}, &models.DraftImage{},
        &models.NewsImage{}, &models.Paragraph{}); err != nil {
        panic("failed to migrate models")
    }
    return db
}

// setupNewsRouter 初始化 Gin 路由和控制器
func setupNewsRouter(db *gorm.DB) *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    {
        newsGroup.Use(middleware.AuthMiddleware())
        {
            // 需要认证的路由
            authGroup := newsGroup.Group("")
            authGroup.Use(middleware.AuthMiddleware())
            {
                authGroup.POST("/upload_image", newsController.UploadImage)      // 上传单张图片
                authGroup.POST("/create_draft", newsController.CreateDraft)     // 创建草稿
                authGroup.PUT("/drafts/:id", newsController.UpdateDraft)
                authGroup.DELETE("/drafts/:id", newsController.DeleteDraft)
                authGroup.POST("/convert_draft", newsController.ConvertDraftToNews)
                authGroup.GET("/my_news", newsController.GetMyNews)
                authGroup.GET("/my_drafts", newsController.GetMyDrafts)
                authGroup.POST("/preview_news", newsController.PreviewNews)
                authGroup.POST("/preview_drafts", newsController.PreviewDrafts)
                authGroup.GET("/details/news/:id", newsController.GetNewsDetails)
                authGroup.GET("/details/draft/:id", newsController.GetDraftDetails)
                authGroup.DELETE("/:id", newsController.DeleteNews)                
                authGroup.GET("/paginated/view_count", newsController.GetNewsByViewCount) // 观看量降序
                authGroup.GET("/paginated/like_count", newsController.GetNewsByLikeCount) // 点赞量降序
                authGroup.GET("/paginated/upload_time", newsController.GetNewsByUploadTime) // 时间由旧到新
                authGroup.POST("/comments", newsController.AddComment)      // 添加评论
                authGroup.DELETE("/comments/:id", newsController.DeleteComment) // 删除评论
                authGroup.POST("/:id/like", newsController.LikeNews)            // 点赞新闻
                authGroup.DELETE("/:id/like", newsController.CancelLikeNews)   // 取消点赞新闻
                authGroup.POST("/:id/favorite", newsController.FavoriteNews)          // 收藏新闻
                authGroup.DELETE("/:id/favorite", newsController.CancelFavoriteNews)  // 取消收藏新闻
                authGroup.POST("/:id/dislike", newsController.DislikeNews)            // 点踩新闻
                authGroup.DELETE("/:id/dislike", newsController.CancelDislikeNews)   // 取消点踩新闻
                authGroup.POST("/:id/view", newsController.ViewNews) // 浏览新闻
                authGroup.GET("/:id/status", newsController.GetUserNewsStatus) // 返回用户对新闻的过往交互
                authGroup.POST("/:id/comment_like", newsController.LikeComment) // 点赞评论
                authGroup.DELETE("/:id/comment_like", newsController.CancelLikeComment) // 取消点赞评论
                authGroup.POST("/search", newsController.SearchNews)
            }
        }
    }

    return router
}

// Helper function to generate a valid JWT for testing
func generateValidJWTNews(userID uint) string {
    token, err := utils.GenerateAccessToken(userID)
    if err != nil {
        panic("Failed to generate valid JWT for testing")
    }
    return token
}

func TestCreateDraft(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建一个用户 => user
    user := models.User{
        OpenID:   "OpenID_CreateDraft_User",
        Nickname: "CreateDraftUser",
    }
    db.Create(&user)

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody: gin.H{
                "title":      "Sample Draft",
                "paragraphs": []string{"Paragraph 1", "Paragraph 2"},
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body",
            userID:         user.ID,
            requestBody:    "not_a_valid_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request data"},
        },
        {
            name:           "Empty Title",
            userID:         user.ID,
            requestBody: gin.H{
                "title":             "",
                "paragraphs":        []string{"Para 1"},
                "image_descriptions": []string{"desc1"},
                "image_paths":        []string{"path1"},
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Title is required"},
        },
        {
            name:           "Mismatch Image Descriptions & Paths",
            userID:         user.ID,
            requestBody: gin.H{
                "title":             "Draft With Mismatch",
                "paragraphs":        []string{"Para1"},
                "image_descriptions": []string{"desc1", "desc2"},
                "image_paths":        []string{"path1"}, // 数量不匹配
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Number of image descriptions and image paths do not match"},
        },
        {
            name:   "Failed To Create Draft (simulate DB error)",
            userID: user.ID,
            requestBody: gin.H{
                "title":             "Draft DB Error",
                "paragraphs":        []string{"Paragraph 1"},
                "image_descriptions": []string{},
                "image_paths":        []string{},
            },
            setupFunc: func() {
                // 模拟数据库错误
                db.Callback().Create().Before("gorm:create").Register("force_draft_create_error", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced draft create error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to create draft"},
        },
        {
            name:   "Success Create Draft",
            userID: user.ID,
            requestBody: gin.H{
                "title":             "My Draft Title",
                "paragraphs":        []string{"Paragraph 1", "Paragraph 2"},
                "image_descriptions": []string{"desc1"},
                "image_paths":        []string{"path1"},
            },
            setupFunc: func() {
                // 移除上一个模拟错误回调
                db.Callback().Create().Remove("force_draft_create_error")
            },
            expectedStatus: http.StatusCreated,
            expectedBody:   map[string]interface{}{"message": "Draft created successfully."},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            // 对成功和失败分别做检查
            if tc.expectedStatus == http.StatusCreated {
                // 只检查 message 字段即可
                assert.Equal(t, tc.expectedBody["message"], resp["message"])
                // draft_id 应该存在
                _, ok := resp["draft_id"]
                assert.True(t, ok)
            } else {
                // 检查 error 字段
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}

func TestUploadImage(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建一个用户 => user
    user := models.User{
        OpenID:   "OpenID_UploadImg_User",
        Nickname: "UploadImgUser",
    }
    db.Create(&user)

    tests := []struct {
        name           string
        userID         uint
        fileField      string
        fileName       string
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            fileField:      "image",
            fileName:       "test_image.jpg",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "No Image File Provided",
            userID:         user.ID,
            fileField:      "", // 不传图
            fileName:       "",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Image file is required"},
        },
        {
            name:           "Failed To Upload Image (simulate error)",
            userID:         user.ID,
            fileField:      "image",
            fileName:       "test_image.jpg",
            setupFunc: func() {
                // 模拟 uploadImage 出错
                // 最简单方式：在 createFile 或 io.Copy 处返回错误
                // 这里暂时仅示例，如果要细粒度 mock，需修改 `uploadImage`。
                os.Setenv("BASE_UPLOAD_PATH", "/invalid/path/that/does/not/exist")
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to upload image"},
        },
        {
            name:           "Success Upload",
            userID:         user.ID,
            fileField:      "image",
            fileName:       "success_test_image.jpg",
            setupFunc: func() {
                // 恢复到有效的上传目录
                os.Setenv("BASE_UPLOAD_PATH", "./test_uploads")
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Image uploaded successfully"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            body := &bytes.Buffer{}
            writer := multipart.NewWriter(body)

            if tc.fileField != "" {
                // 构造一个虚拟文件
                part, _ := writer.CreateFormFile(tc.fileField, tc.fileName)
                // 写点内容假装是文件数据
                _, _ = part.Write([]byte("fake_image_data"))
            }
            writer.Close()

            req, _ := http.NewRequest("POST", "/news/upload_image", body)
            req.Header.Set("Content-Type", writer.FormDataContentType())

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.expectedStatus == http.StatusOK {
                // 检查 message 字段
                assert.Equal(t, tc.expectedBody["message"], resp["message"])
                // 应该返回 path
                _, ok := resp["path"]
                assert.True(t, ok)
            } else {
                // 错误情况检查
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }

            // 如果有生成文件，需要清理
            if tc.expectedStatus == http.StatusOK {
                // 移除测试生成的文件夹
                os.RemoveAll("./test_uploads")
            }
        })
    }
}

func TestConvertDraftToNews(t *testing.T) {
    db := setupNewsTestDB()
    // 记得还要 AutoMigrate models.News, models.Paragraph, models.NewsImage，如果需要
    _ = db.AutoMigrate(&models.News{}, &models.Paragraph{}, &models.NewsImage{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    {
        newsGroup.Use(middleware.AuthMiddleware())
        {
            newsGroup.POST("/convert_draft", newsController.ConvertDraftToNews)
        }
    }

    // 创建一个用户
    user := models.User{
        OpenID:   "OpenID_ConvertDraft_User",
        Nickname: "UserConvertDraft",
    }
    db.Create(&user)

    // 插入一个草稿 => draft
    draft := models.Draft{
        Title:    "DraftToConvert",
        AuthorID: user.ID,
    }
    db.Create(&draft)

    // 写段落、图片
    db.Create(&models.DraftParagraph{DraftID: draft.ID, Text: "Draft Para 1"})
    db.Create(&models.DraftParagraph{DraftID: draft.ID, Text: "Draft Para 2"})
    db.Create(&models.DraftImage{DraftID: draft.ID, URL: "path/to/image1", Description: "desc1"})

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"draft_id": draft.ID},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Request Body",
            userID:         user.ID,
            requestBody:    "not_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request body"},
        },
        {
            name:           "Draft Not Found",
            userID:         user.ID,
            requestBody:    gin.H{"draft_id": 99999},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "Draft not found"},
        },
        {
            name:           "Failed To Find Draft (simulate error)",
            userID:         user.ID,
            requestBody:    gin.H{"draft_id": draft.ID},
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_find_draft_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced find draft error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to find draft"},
        },
        {
            name:           "Failed To Create News (simulate error)",
            userID:         user.ID,
            requestBody:    gin.H{"draft_id": draft.ID},
            setupFunc: func() {
                db.Callback().Query().Remove("force_find_draft_err")
                // 模拟创建 news 出错
                db.Callback().Create().Before("gorm:create").Register("force_create_news_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced create news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to create news"},
        },
        {
            name:           "Success Convert Draft",
            userID:         user.ID,
            requestBody:    gin.H{"draft_id": draft.ID},
            setupFunc: func() {
                db.Callback().Create().Remove("force_create_news_err")
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Draft converted to news successfully."},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/news/convert_draft", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.expectedStatus == http.StatusOK {
                // 只检查 message
                assert.Equal(t, tc.expectedBody["message"], resp["message"])
                // 还会有 news_id
                _, ok := resp["news_id"]
                assert.True(t, ok)
            } else {
                // 检查 error
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}

func TestUpdateDraft(t *testing.T) {
    db := setupNewsTestDB()
    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    {
        newsGroup.Use(middleware.AuthMiddleware())
        {
            newsGroup.PUT("/draft/:id", newsController.UpdateDraft)
        }
    }

    draftOwner := models.User{
        OpenID:   "OpenID_UpdateDraft_Owner",
        Nickname: "DraftOwner",
    }
    db.Create(&draftOwner)

    otherUser := models.User{
        OpenID:   "OpenID_UpdateDraft_Other",
        Nickname: "OtherUser",
    }
    db.Create(&otherUser)

    // 创建草稿 => ownedDraft
    ownedDraft := models.Draft{
        Title:    "DraftByOwner",
        AuthorID: draftOwner.ID,
    }
    db.Create(&ownedDraft)

    // 创建几张旧图片(后面测试删除)
    oldImage := models.DraftImage{
        DraftID: ownedDraft.ID,
        URL:     "old_path1",
    }
    db.Create(&oldImage)
    oldImage2 := models.DraftImage{
        DraftID: ownedDraft.ID,
        URL:     "old_path2",
    }
    db.Create(&oldImage2)

    tests := []struct {
        name           string
        userID         uint
        draftID        string
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            draftID:        "1",
            requestBody:    gin.H{"title": "AnyTitle"},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Draft ID",
            userID:         draftOwner.ID,
            draftID:        "abc",
            requestBody:    gin.H{"title": "NewTitle"},
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid draft ID"},
        },
        {
            name:           "Draft Not Found",
            userID:         draftOwner.ID,
            draftID:        "99999",
            requestBody:    gin.H{"title": "NewTitle"},
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "Draft not found"},
        },
        {
            name:           "No Permission (403)",
            userID:         otherUser.ID,
            draftID:        fmt.Sprintf("%d", ownedDraft.ID),
            requestBody: gin.H{
                "title":              "TryUpdateOthers",
                "paragraphs":         []string{"Para1", "Para2"},
                "image_descriptions": []string{},
                "image_paths":        []string{},
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusForbidden,
            expectedBody:   map[string]interface{}{"error": "You do not have permission to edit this draft"},
        },
        {
            name:           "Invalid Request Data",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", ownedDraft.ID),
            requestBody:    "not_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid request data"},
        },
        {
            name:           "Title Is Required",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", ownedDraft.ID),
            requestBody: gin.H{
                "title":              "",
                "paragraphs":         []string{"Paragraph1"},
                "image_descriptions": []string{},
                "image_paths":        []string{},
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Title is required"},
        },
        {
            name:           "Mismatch Image Descriptions & Paths",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", ownedDraft.ID),
            requestBody: gin.H{
                "title":              "ValidTitle",
                "paragraphs":         []string{"P1", "P2"},
                "image_descriptions": []string{"desc1", "desc2"},
                "image_paths":        []string{"path1"},
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Number of image descriptions and image paths do not match"},
        },
        {
            name:           "Failed To Delete Old Draft (simulate TX error)",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", ownedDraft.ID),
            requestBody: gin.H{
                "title":              "NewTitleAfterDeleteFail",
                "paragraphs":         []string{"P1"},
                "image_descriptions": []string{},
                "image_paths":        []string{},
            },
            setupFunc: func() {
                // 在事务中 delete(draft) 时注入错误
                db.Callback().Delete().Before("gorm:delete").Register("force_delete_old_draft_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced delete old draft error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to delete old draft"},
        },
        
        {
            name:           "Failed To Create New Draft (simulate error)",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", ownedDraft.ID),
            requestBody: gin.H{
                "title":              "NewTitleAfterCreateFail",
                "paragraphs":         []string{"P1", "P2"},
                "image_descriptions": []string{"desc1"},
                "image_paths":        []string{"path1"},
            },
            setupFunc: func() {
                // 在新建 draft 时注入错误，使用唯一回调名称
                db.Callback().Create().Before("gorm:create").Register("force_create_new_draft_err_test1", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced create new draft error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to delete old draft"},
        },
        {
            name:           "Success Update Draft",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", ownedDraft.ID),
            requestBody: gin.H{
                "title":              "UpdatedTitle",
                "paragraphs":         []string{"NewPara1", "NewPara2"},
                "image_descriptions": []string{"newDesc1", "newDesc2"},
                "image_paths":        []string{"newPath1", "old_path2"}, 
                // 这里故意保留 old_path2 以模拟只删除 old_path1
            },
            setupFunc: func() {
                // 移除错误回调，使用唯一名称
                db.Callback().Create().Remove("force_create_new_draft_err_test1")
                db.Callback().Delete().Remove("force_delete_old_draft_err")
                db.Callback().Update().Remove("force_save_pending_error")
            },
            expectedStatus: http.StatusCreated,
            expectedBody:   map[string]interface{}{"message": "Draft updated successfully"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("PUT", "/news/draft/"+tc.draftID, bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.expectedStatus == http.StatusCreated {
                assert.Equal(t, tc.expectedBody["message"], resp["message"])
                // draft_id 应该存在
                _, ok := resp["draft_id"]
                assert.True(t, ok)
            } else {
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}

func TestDeleteDraft(t *testing.T) {
    db := setupNewsTestDB()

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.DELETE("/draft/:id", newsController.DeleteDraft)
    }

    // 创建用户 => draftOwner, otherUser
    draftOwner := models.User{
        OpenID:   "OpenID_DeleteDraft_Owner",
        Nickname: "DraftOwner",
    }
    db.Create(&draftOwner)

    otherUser := models.User{
        OpenID:   "OpenID_DeleteDraft_Other",
        Nickname: "OtherDraftUser",
    }
    db.Create(&otherUser)

    // 创建一个草稿 => draft
    draft := models.Draft{
        Title:    "DraftToDelete",
        AuthorID: draftOwner.ID,
    }
    db.Create(&draft)

    // 给此草稿插入一些关联图片 => DraftImage
    img1 := models.DraftImage{DraftID: draft.ID, URL: "some_local_path.jpg"}
    db.Create(&img1)

    tests := []struct {
        name           string
        userID         uint
        draftID        string
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            draftID:        "1",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid Draft ID",
            userID:         draftOwner.ID,
            draftID:        "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid draft ID"},
        },
        {
            name:           "Draft Not Found",
            userID:         draftOwner.ID,
            draftID:        "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "Draft not found"},
        },
        {
            name:           "No Permission (403)",
            userID:         otherUser.ID,
            draftID:        fmt.Sprintf("%d", draft.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusForbidden,
            expectedBody:   map[string]interface{}{"error": "You do not have permission to delete this draft"},
        },
        {
            name:           "Failed to Delete Draft (simulate DB error)",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", draft.ID),
            setupFunc: func() {
                // 在 delete(draft) 或 delete(draftParagraph) 时注入错误
                db.Callback().Delete().Before("gorm:delete").Register("force_delete_draft_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced delete draft error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to delete draft"},
        },
        {
            name:           "Success Delete Draft",
            userID:         draftOwner.ID,
            draftID:        fmt.Sprintf("%d", draft.ID),
            setupFunc: func() {
                // 移除上一个回调
                db.Callback().Delete().Remove("force_delete_draft_err")
                // 如果需要重复测试删除，用完后还要重新创建 Draft
                // 这里仅示例一次删除操作
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "Draft deleted successfully."},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("DELETE", "/news/draft/"+tc.draftID, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.expectedStatus == http.StatusOK {
                assert.Equal(t, tc.expectedBody["message"], resp["message"])
            } else {
                // 对错误情况进行断言
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}
