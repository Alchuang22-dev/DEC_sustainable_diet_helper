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

func TestDeleteNews(t *testing.T) {
    db := setupNewsTestDB()
    // DeleteNews 会用到 News, NewsImage, Paragraph, etc.
    db.AutoMigrate(&models.News{}, &models.NewsImage{}, &models.Paragraph{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.DELETE("/news/:id", newsController.DeleteNews)
    }

    // 创建两个用户 => newsOwner, otherUser
    newsOwner := models.User{
        OpenID:   "OpenID_DeleteNews_Owner",
        Nickname: "NewsOwner",
    }
    db.Create(&newsOwner)

    otherUser := models.User{
        OpenID:   "OpenID_DeleteNews_Other",
        Nickname: "NewsOtherUser",
    }
    db.Create(&otherUser)

    // 创建一条新闻 => news
    myNews := models.News{
        Title:    "NewsToDelete",
        AuthorID: newsOwner.ID,
    }
    db.Create(&myNews)

    // 创建关联图片
    newsImg := models.NewsImage{
        NewsID: myNews.ID,
        URL:    "some_news_path.jpg",
    }
    db.Create(&newsImg)

    tests := []struct {
        name           string
        userID         uint
        newsID         string
        setupFunc      func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            newsID:         "1",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authorization header missing"},
        },
        {
            name:           "Invalid News ID",
            userID:         newsOwner.ID,
            newsID:         "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   map[string]interface{}{"error": "Invalid news ID"},
        },
        {
            name:           "News Not Found",
            userID:         newsOwner.ID,
            newsID:         "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedBody:   map[string]interface{}{"error": "News not found"},
        },
        {
            name:           "No Permission (403)",
            userID:         otherUser.ID,
            newsID:         fmt.Sprintf("%d", myNews.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusForbidden,
            expectedBody:   map[string]interface{}{"error": "You do not have permission to delete this news"},
        },
        {
            name:           "Failed To Delete News (simulate DB error)",
            userID:         newsOwner.ID,
            newsID:         fmt.Sprintf("%d", myNews.ID),
            setupFunc: func() {
                db.Callback().Delete().Before("gorm:delete").Register("force_delete_news_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced delete news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "Failed to delete news"},
        },
        {
            name:           "Success Delete News",
            userID:         newsOwner.ID,
            newsID:         fmt.Sprintf("%d", myNews.ID),
            setupFunc: func() {
                db.Callback().Delete().Remove("force_delete_news_err")
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"message": "News deleted successfully."},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("DELETE", "/news/news/"+tc.newsID, nil)
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
                for k, v := range tc.expectedBody {
                    assert.Equal(t, v, resp[k])
                }
            }
        })
    }
}

func TestGetMyNews(t *testing.T) {
    db := setupNewsTestDB()
    // 需要存储 News，AutoMigrate 以免报错
    db.AutoMigrate(&models.News{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.GET("/my_news", newsController.GetMyNews)
    }

    // 创建用户: userWithNews, userNoNews
    userWithNews := models.User{
        OpenID:   "OpenID_GetMyNews_Has",
        Nickname: "UserHasNews",
    }
    db.Create(&userWithNews)

    userNoNews := models.User{
        OpenID:   "OpenID_GetMyNews_None",
        Nickname: "UserNoNews",
    }
    db.Create(&userNoNews)

    // 为 userWithNews 创建一些 news
    news1 := models.News{Title: "My News1", AuthorID: userWithNews.ID, UploadTime: time.Now().Add(-2 * time.Hour)}
    db.Create(&news1)
    news2 := models.News{Title: "My News2", AuthorID: userWithNews.ID, UploadTime: time.Now().Add(-1 * time.Hour)}
    db.Create(&news2)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedIDs:    nil,
        },
        {
            name:           "User With No News",
            userID:         userNoNews.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusOK,
            expectedIDs:    []uint{},
        },
        {
            name:   "Failed To Fetch News (simulate DB error)",
            userID: userWithNews.ID,
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_get_my_news_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced get my news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedIDs:    nil,
        },
        {
            name:   "Success Get My News",
            userID: userWithNews.ID,
            setupFunc: func() {
                db.Callback().Query().Remove("force_get_my_news_err")
            },
            expectedStatus: http.StatusOK,
            // Order("upload_time DESC") => news2 比 news1 时间更新 => [news2, news1]
            expectedIDs: []uint{news2.ID, news1.ID},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("GET", "/news/my_news", nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            if w.Code == http.StatusOK {
                var resp map[string][]uint
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, tc.expectedIDs, resp["news_ids"])
            } else if w.Code == http.StatusUnauthorized {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, "Authorization header missing", resp["error"])
            } else if w.Code == http.StatusInternalServerError {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, "Failed to fetch news", resp["error"])
            }
        })
    }
}

func TestGetMyDrafts(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.Draft{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.GET("/my_drafts", newsController.GetMyDrafts)
    }

    // 创建用户: userWithDrafts, userNoDrafts
    userWithDrafts := models.User{
        OpenID:   "OpenID_GetMyDrafts_Has",
        Nickname: "UserHasDrafts",
    }
    db.Create(&userWithDrafts)

    userNoDrafts := models.User{
        OpenID:   "OpenID_GetMyDrafts_None",
        Nickname: "UserNoDrafts",
    }
    db.Create(&userNoDrafts)

    // 创建一些草稿
    d1 := models.Draft{Title: "Draft1", AuthorID: userWithDrafts.ID, UpdatedAt: time.Now().Add(-2 * time.Hour)}
    db.Create(&d1)
    d2 := models.Draft{Title: "Draft2", AuthorID: userWithDrafts.ID, UpdatedAt: time.Now().Add(-1 * time.Hour)}
    db.Create(&d2)

    tests := []struct {
        name           string
        userID         uint
        setupFunc      func()
        expectedStatus int
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedIDs:    nil,
        },
        {
            name:           "User With No Drafts",
            userID:         userNoDrafts.ID,
            setupFunc:      func() {},
            expectedStatus: http.StatusOK,
            expectedIDs:    []uint{},
        },
        {
            name:   "Failed To Fetch Drafts (simulate error)",
            userID: userWithDrafts.ID,
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_get_my_drafts_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced get my drafts error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedIDs:    nil,
        },
        {
            name:   "Success Get My Drafts",
            userID: userWithDrafts.ID,
            setupFunc: func() {
                db.Callback().Query().Remove("force_get_my_drafts_err")
            },
            expectedStatus: http.StatusOK,
            // Order("updated_at DESC") => d2 比 d1 时间更新 => [d2, d1]
            expectedIDs: []uint{d2.ID, d1.ID},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            req, _ := http.NewRequest("GET", "/news/my_drafts", nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            if w.Code == http.StatusOK {
                var resp map[string][]uint
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, tc.expectedIDs, resp["draft_ids"])
            } else if w.Code == http.StatusUnauthorized {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, "Authorization header missing", resp["error"])
            } else if w.Code == http.StatusInternalServerError {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, "Failed to fetch drafts", resp["error"])
            }
        })
    }
}

func TestPreviewNews(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.News{}, &models.NewsImage{}, &models.Paragraph{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)
    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.POST("/preview_news", newsController.PreviewNews)
    }

    // 创建测试用户
    user := models.User{
        OpenID:    "OpenID_PreviewNews_User",
        Nickname:  "PreviewNewsUser",
    }
    db.Create(&user)

    // 创建新闻数据: news1, news2
    author1 := models.User{OpenID: "OpenID_Author1", Nickname: "AuthorOne"}
    db.Create(&author1)
    author2 := models.User{OpenID: "OpenID_Author2", Nickname: "AuthorTwo"}
    db.Create(&author2)

    news1 := models.News{Title: "PreviewTitle1", AuthorID: author1.ID, LikeCount: 10}
    db.Create(&news1)
    db.Create(&models.Paragraph{NewsID: news1.ID, Text: "Paragraph1 for News1"})
    db.Create(&models.NewsImage{NewsID: news1.ID, URL: "img1.jpg", Description: "desc1"})

    news2 := models.News{Title: "PreviewTitle2", AuthorID: author2.ID, LikeCount: 5}
    db.Create(&news2)
    db.Create(&models.Paragraph{NewsID: news2.ID, Text: "Paragraph1 for News2"})
    db.Create(&models.NewsImage{NewsID: news2.ID, URL: "img2.jpg", Description: "desc2"})

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedError  string // 对错误场景做断言
        isSuccess      bool   // 是否期望成功
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"ids": []uint{news1.ID, news2.ID}},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Request Body",
            userID:         user.ID,
            requestBody:    "not_valid_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "Database Error (simulate)",
            userID:         user.ID,
            requestBody:    gin.H{"ids": []uint{news1.ID, news2.ID}},
            setupFunc: func() {
                // 模拟 DB 错误
                db.Callback().Query().Before("gorm:query").Register("force_preview_news_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced preview news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch news",
        },
        {
            name:           "Success Preview",
            userID:         user.ID,
            requestBody:    gin.H{"ids": []uint{news1.ID, news2.ID}},
            setupFunc: func() {
                db.Callback().Query().Remove("force_preview_news_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/news/preview_news", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            // Token
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            // 处理响应
            if tc.isSuccess {
                // 成功时返回 200 + 预览列表
                var resp DraftDetailResponse
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
            } else {
                // 错误时检查 error 字段
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestPreviewDrafts(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.Draft{}, &models.DraftParagraph{}, &models.DraftImage{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.POST("/preview_drafts", newsController.PreviewDrafts)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_PreviewDrafts_User",
        Nickname: "PreviewDraftsUser",
    }
    db.Create(&user)

    // 创建几个草稿
    d1 := models.Draft{Title: "DraftTitle1", AuthorID: user.ID}
    db.Create(&d1)
    db.Create(&models.DraftParagraph{DraftID: d1.ID, Text: "Draft1 Para1"})
    db.Create(&models.DraftImage{DraftID: d1.ID, URL: "draft1_img1.jpg", Description: "desc_draft1_img1"})

    d2 := models.Draft{Title: "DraftTitle2", AuthorID: user.ID}
    db.Create(&d2)
    db.Create(&models.DraftParagraph{DraftID: d2.ID, Text: "Draft2 Para1"})
    db.Create(&models.DraftImage{DraftID: d2.ID, URL: "draft2_img1.jpg", Description: "desc_draft2_img1"})

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody:    gin.H{"ids": []uint{d1.ID, d2.ID}},
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Request Body",
            userID:         user.ID,
            requestBody:    "not_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "DB Error (simulate)",
            userID:         user.ID,
            requestBody:    gin.H{"ids": []uint{d1.ID, d2.ID}},
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_preview_drafts_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced preview drafts error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch drafts",
        },
        {
            name:           "Success Preview Drafts",
            userID:         user.ID,
            requestBody:    gin.H{"ids": []uint{d1.ID, d2.ID}},
            setupFunc: func() {
                db.Callback().Query().Remove("force_preview_drafts_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/news/preview_drafts", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            if tc.isSuccess {
                var resp DraftDetailResponse
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
            } else {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestGetNewsDetails(t *testing.T) {
    db := setupNewsTestDB()
    // 需要迁移 News, Paragraph, NewsImage, Comment, User 等
    db.AutoMigrate(&models.News{}, &models.NewsImage{}, &models.Paragraph{}, &models.Comment{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.GET("/details/news/:id", newsController.GetNewsDetails)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_GetNewsDetails_User",
        Nickname: "GetNewsUser",
    }
    db.Create(&user)

    // 创建新闻
    newsItem := models.News{
        Title:       "NewsDetailsTest",
        AuthorID:    user.ID,
        UploadTime:  time.Now(),
        ViewCount:   100,
        LikeCount:   10,
        FavoriteCount: 5,
        DislikeCount: 2,
        ShareCount:  1,
    }
    db.Create(&newsItem)
    // 插入段落、图片
    db.Create(&models.Paragraph{NewsID: newsItem.ID, Text: "Paragraph1"})
    db.Create(&models.NewsImage{NewsID: newsItem.ID, URL: "news_img1.jpg", Description: "img_desc1"})

    tests := []struct {
        name           string
        userID         uint
        newsID         string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid News ID",
            userID:         user.ID,
            newsID:         "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid news ID",
        },
        {
            name:           "News Not Found",
            userID:         user.ID,
            newsID:         "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "News not found",
        },
        {
            name:           "Failed To Fetch News Detail (simulate DB error)",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_get_news_detail_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced get news detail error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch news details",
        },
        {
            name:           "Success Get News Detail",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_get_news_detail_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/news/details/news/" + tc.newsID
            req, _ := http.NewRequest("GET", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            if tc.isSuccess {
                var resp DraftDetailResponse
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, newsItem.ID, resp.ID)
            } else {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestGetDraftDetails(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.Draft{}, &models.DraftParagraph{}, &models.DraftImage{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.GET("/details/draft/:id", newsController.GetDraftDetails)
    }

    // 创建用户 => draftAuthor, otherUser
    draftAuthor := models.User{
        OpenID:   "OpenID_GetDraftDetails_Author",
        Nickname: "DraftAuthor",
    }
    db.Create(&draftAuthor)

    otherUser := models.User{
        OpenID:   "OpenID_GetDraftDetails_Other",
        Nickname: "OtherDraftUser",
    }
    db.Create(&otherUser)

    // 创建草稿 => draft
    draft := models.Draft{
        Title:    "DraftDetailsTitle",
        AuthorID: draftAuthor.ID,
        CreatedAt: time.Now().Add(-2 * time.Hour),
        UpdatedAt: time.Now().Add(-1 * time.Hour),
    }
    db.Create(&draft)
    // 段落、图片
    db.Create(&models.DraftParagraph{DraftID: draft.ID, Text: "DraftPara1"})
    db.Create(&models.DraftImage{DraftID: draft.ID, URL: "draft_img1.jpg", Description: "desc_draft_img1"})

    tests := []struct {
        name           string
        userID         uint
        draftID        string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            draftID:        fmt.Sprintf("%d", draft.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Draft ID",
            userID:         draftAuthor.ID,
            draftID:        "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid draft ID",
        },
        {
            name:           "Draft Not Found",
            userID:         draftAuthor.ID,
            draftID:        "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "Draft not found",
        },
        {
            name:           "Not Author => Not Found",
            userID:         otherUser.ID,
            draftID:        fmt.Sprintf("%d", draft.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "Draft not found",
        },
        {
            name:           "Failed To Fetch Draft (simulate DB error)",
            userID:         draftAuthor.ID,
            draftID:        fmt.Sprintf("%d", draft.ID),
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_get_draft_detail_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "drafts" {
                        tx.Error = fmt.Errorf("forced get draft detail error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch draft details",
        },
        {
            name:           "Success Get Draft Detail",
            userID:         draftAuthor.ID,
            draftID:        fmt.Sprintf("%d", draft.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_get_draft_detail_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/news/details/draft/" + tc.draftID
            req, _ := http.NewRequest("GET", url, nil)

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            if tc.isSuccess {
                var resp DraftDetailResponse
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, draft.ID, resp.ID)
                assert.Equal(t, draft.Title, resp.Title)
            } else {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestGetNewsByViewCount(t *testing.T) {
    db := setupNewsTestDB()
    // 需要迁移 News，避免操作报错
    db.AutoMigrate(&models.News{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.GET("/paginated/view_count", newsController.GetNewsByViewCount)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_GetNewsByViewCount_User",
        Nickname: "ViewCountUser",
    }
    db.Create(&user)

    // 创建一些新闻，设置不同的 view_count
    news1 := models.News{Title: "News1", ViewCount: 500}
    db.Create(&news1)
    news2 := models.News{Title: "News2", ViewCount: 1000}
    db.Create(&news2)
    news3 := models.News{Title: "News3", ViewCount: 300}
    db.Create(&news3)

    tests := []struct {
        name           string
        userID         uint
        page           string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            page:           "1",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Page Number",
            userID:         user.ID,
            page:           "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid page number",
        },
        {
            name:           "DB Error (simulate)",
            userID:         user.ID,
            page:           "1",
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_viewcount_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced viewcount error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch news",
        },
        {
            name:           "Success First Page",
            userID:         user.ID,
            page:           "1",
            setupFunc: func() {
                // 移除上一个回调
                db.Callback().Query().Remove("force_viewcount_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            // 按 view_count 降序 => news2(1000), news1(500), news3(300)
            expectedIDs: []uint{news2.ID, news1.ID, news3.ID},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/news/paginated/view_count?page=" + tc.page
            req, _ := http.NewRequest("GET", url, nil)

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            if tc.isSuccess {
                var resp map[string][]uint
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                // 校验顺序
                assert.Equal(t, tc.expectedIDs, resp["news_ids"])
            } else {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestGetNewsByLikeCount(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.News{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.GET("/paginated/like_count", newsController.GetNewsByLikeCount)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_GetNewsByLikeCount_User",
        Nickname: "LikeCountUser",
    }
    db.Create(&user)

    // 插入新闻
    n1 := models.News{Title: "News1", LikeCount: 50}
    db.Create(&n1)
    n2 := models.News{Title: "News2", LikeCount: 100}
    db.Create(&n2)
    n3 := models.News{Title: "News3", LikeCount: 10}
    db.Create(&n3)

    tests := []struct {
        name           string
        userID         uint
        page           string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized",
            userID:         0,
            page:           "1",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Page Number",
            userID:         user.ID,
            page:           "zero",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid page number",
        },
        {
            name:           "DB Error",
            userID:         user.ID,
            page:           "1",
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_likecount_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced likecount error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch news",
        },
        {
            name:           "Success (like_count desc)",
            userID:         user.ID,
            page:           "1",
            setupFunc: func() {
                db.Callback().Query().Remove("force_likecount_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            // 按 like_count 降序 => n2(100), n1(50), n3(10)
            expectedIDs: []uint{n2.ID, n1.ID, n3.ID},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/news/paginated/like_count?page=" + tc.page
            req, _ := http.NewRequest("GET", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            if tc.isSuccess {
                var resp map[string][]uint
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, tc.expectedIDs, resp["news_ids"])
            } else {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestGetNewsByUploadTime(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.News{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.GET("/paginated/upload_time", newsController.GetNewsByUploadTime)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_GetNewsByUploadTime_User",
        Nickname: "UploadTimeUser",
    }
    db.Create(&user)

    // 插入新闻: UploadTime 不同
    n1 := models.News{Title: "OldNews", UploadTime: time.Now().Add(-2 * time.Hour)}
    db.Create(&n1)
    n2 := models.News{Title: "NewerNews", UploadTime: time.Now().Add(-1 * time.Hour)}
    db.Create(&n2)
    n3 := models.News{Title: "NewestNews", UploadTime: time.Now()}
    db.Create(&n3)

    tests := []struct {
        name           string
        userID         uint
        page           string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        expectedIDs    []uint
    }{
        {
            name:           "Unauthorized",
            userID:         0,
            page:           "1",
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Page Number",
            userID:         user.ID,
            page:           "-1",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid page number",
        },
        {
            name:           "DB Error",
            userID:         user.ID,
            page:           "1",
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_uploadtime_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced upload_time error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to fetch news",
        },
        {
            name:           "Success (upload_time desc)",
            userID:         user.ID,
            page:           "1",
            setupFunc: func() {
                db.Callback().Query().Remove("force_uploadtime_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            // upload_time DESC => n3, n2, n1
            expectedIDs: []uint{n3.ID, n2.ID, n1.ID},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/news/paginated/upload_time?page=" + tc.page
            req, _ := http.NewRequest("GET", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            if tc.isSuccess {
                var resp map[string][]uint
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, tc.expectedIDs, resp["news_ids"])
            } else {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestAddComment(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.News{}, &models.Comment{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.POST("/comments", newsController.AddComment)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_AddComment_User",
        Nickname: "CommentUser",
    }
    db.Create(&user)

    // 创建新闻
    newsItem := models.News{Title: "CommentableNews"}
    db.Create(&newsItem)

    // 创建一个父评论
    parentComment := models.Comment{
        NewsID:    newsItem.ID,
        Content:   "Parent comment",
        UserID:    9999,
        IsReply:   false,
    }
    db.Create(&parentComment)

    tests := []struct {
        name           string
        userID         uint
        requestBody    interface{}
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            requestBody: gin.H{
                "news_id":  newsItem.ID,
                "content":  "Some comment",
                "is_reply": false,
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Request Body",
            userID:         user.ID,
            requestBody:    "not_json",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid request body",
        },
        {
            name:           "Invalid News ID",
            userID:         user.ID,
            requestBody: gin.H{
                "news_id":  99999,  // 不存在
                "content":  "Some comment",
                "is_reply": false,
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid news ID",
        },
        {
            name:           "Reply Without ParentID",
            userID:         user.ID,
            requestBody: gin.H{
                "news_id":  newsItem.ID,
                "content":  "Reply content",
                "is_reply": true,
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "ParentID is required for a reply",
        },
        {
            name:           "Parent Comment Not Found",
            userID:         user.ID,
            requestBody: gin.H{
                "news_id":   newsItem.ID,
                "content":   "Reply content",
                "is_reply":  true,
                "parent_id": 88888, // 不存在
            },
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Parent comment not found",
        },
        {
            name:           "Parent Comment Not Same News",
            userID:         user.ID,
            requestBody: gin.H{
                "news_id":   newsItem.ID,
                "content":   "Reply content",
                "is_reply":  true,
                "parent_id": parentComment.ID,
            },
            setupFunc: func() {
                // 将 parentComment 指向不同 news
                db.Model(&parentComment).Update("news_id", 99999)
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Parent comment does not belong to the specified news",
        },
        {
            name:           "DB Error (simulate)",
            userID:         user.ID,
            requestBody: gin.H{
                "news_id":  newsItem.ID,
                "content":  "New comment",
                "is_reply": false,
            },
            setupFunc: func() {
                db.Callback().Create().Before("gorm:create").Register("force_add_comment_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "comments" {
                        tx.Error = fmt.Errorf("forced add comment error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to add comment",
        },
        {
            name:           "Success Add Comment",
            userID:         user.ID,
            requestBody: gin.H{
                "news_id":  newsItem.ID,
                "content":  "A top-level comment",
                "is_reply": false,
            },
            setupFunc: func() {
                db.Callback().Create().Remove("force_add_comment_err")
            },
            expectedStatus: http.StatusCreated,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            bodyBytes, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest("POST", "/news/comments", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            if tc.isSuccess {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                assert.Equal(t, "Comment added successfully", resp["message"])
                // 可以进一步校验 comment 字段内容
            } else {
                var resp map[string]interface{}
                err := json.Unmarshal(w.Body.Bytes(), &resp)
                assert.NoError(t, err)
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestLikeComment(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.Comment{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.POST("/:id/comment_like", newsController.LikeComment)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_LikeComment_User",
        Nickname: "LikeCommentUser",
    }
    db.Create(&user)

    // 创建评论
    comment := models.Comment{
        NewsID:  999, // 对应的 News 不一定要存在
        Content: "Some comment",
        UserID:  12345,
    }
    db.Create(&comment)

    tests := []struct {
        name           string
        userID         uint
        commentID      string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        finalLikeCount int // 用于成功情况下的断言
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Comment ID",
            userID:         user.ID,
            commentID:      "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid comment ID",
        },
        {
            name:           "Comment Not Found",
            userID:         user.ID,
            commentID:      "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "Comment not found",
        },
        {
            name:           "Failed To Find Comment (simulate)",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_find_comment_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "comments" {
                        tx.Error = fmt.Errorf("forced find comment error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find comment",
        },
        {
            name:           "Failed To Find User",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                // 移除查 comment 时的错误
                db.Callback().Query().Remove("force_find_comment_err")
                // 模拟查 user 错误
                db.Callback().Query().Before("gorm:query").Register("force_find_user_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced find user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find user",
        },
        {
            name:           "Already Liked This Comment",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_find_user_err")

                // 让 user.LikedComments 包含这个 comment
                db.Model(&user).Association("LikedComments").Append(&comment)
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "You have already liked this comment",
        },
        {
            name:           "Failed To Like Comment (append error)",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                // 移除已经点赞关系
                db.Model(&user).Association("LikedComments").Clear()

                // 模拟 append 出错
                db.Callback().Update().Before("gorm:association").Register("force_append_err", func(tx *gorm.DB) {
                    tx.Error = fmt.Errorf("forced append error")
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to like comment",
        },
        {
            name:           "Failed To Update Comment like_count",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_append_err")

                // mock updateColumn 出错
                db.Callback().Update().Before("gorm:update").Register("force_update_like_count_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "comments" {
                        tx.Error = fmt.Errorf("forced update like_count error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update comment like_count",
        },
        {
            name:           "Success Like Comment",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_update_like_count_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            finalLikeCount: comment.LikeCount + 1,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := fmt.Sprintf("/news/%s/comment_like", tc.commentID)
            req, _ := http.NewRequest("POST", url, nil)

            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "Comment liked successfully", resp["message"])
                assert.Equal(t, float64(tc.finalLikeCount), resp["like_count"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestCancelLikeComment(t *testing.T) {
    db := setupNewsTestDB()
    // 需要迁移 Comment, User
    db.AutoMigrate(&models.Comment{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.DELETE("/:id/comment_like", newsController.CancelLikeComment)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_CancelLikeComment_User",
        Nickname: "CancelLikeCommentUser",
    }
    db.Create(&user)

    // 创建评论
    comment := models.Comment{
        NewsID:  9999,
        Content: "Liked comment",
        UserID:  123, // 不一定要存在对应 User
        LikeCount: 10,
    }
    db.Create(&comment)

    // 让 user 已经点赞了该 comment
    db.Model(&user).Association("LikedComments").Append(&comment)

    tests := []struct {
        name           string
        userID         uint
        commentID      string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        finalLikeCount int
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Comment ID",
            userID:         user.ID,
            commentID:      "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid comment ID",
        },
        {
            name:           "Comment Not Found",
            userID:         user.ID,
            commentID:      "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "Comment not found",
        },
        {
            name:           "Failed To Find Comment (simulate DB error)",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_find_comment_err_cancel", func(tx *gorm.DB) {
                    if tx.Statement.Table == "comments" {
                        tx.Error = fmt.Errorf("forced find comment error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find comment",
        },
        {
            name:           "Failed To Find User (simulate DB error)",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                // 移除前面注入的错误
                db.Callback().Query().Remove("force_find_comment_err_cancel")
                // 模拟查 user 出错
                db.Callback().Query().Before("gorm:query").Register("force_find_user_err_cancel", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced find user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find user",
        },
        {
            name:           "User Not Liked This Comment",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                // 移除上一个回调
                db.Callback().Query().Remove("force_find_user_err_cancel")
                // 先清空用户点赞，让他没点赞过
                db.Model(&user).Association("LikedComments").Clear()
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "You have not liked this comment",
        },
        {
            name:           "Failed To Update like_count",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                db.Model(&user).Association("LikedComments").Append(&comment)
                db.Callback().Update().Remove("force_cancel_like_assoc_err")

                // 注入 updateColumn 错误
                db.Callback().Update().Before("gorm:update").Register("force_update_likecount_err_cancel", func(tx *gorm.DB) {
                    if tx.Statement.Table == "comments" {
                        tx.Error = fmt.Errorf("forced update comment like_count error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update comment like_count",
        },
        {
            name:           "Success Cancel Like Comment",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                // 移除上一个错误回调
                db.Callback().Update().Remove("force_update_likecount_err_cancel")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            finalLikeCount: comment.LikeCount - 1,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := fmt.Sprintf("/news/%s/comment_like", tc.commentID)
            req, _ := http.NewRequest("DELETE", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "Comment like canceled successfully", resp["message"])
                // like_count => 原值 - 1
                assert.Equal(t, float64(tc.finalLikeCount), resp["like_count"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestDeleteComment(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.Comment{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.DELETE("/comments/:id", newsController.DeleteComment)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_DeleteComment_User",
        Nickname: "DeleteCommentUser",
    }
    db.Create(&user)

    // 创建评论
    comment := models.Comment{
        NewsID:    99999,
        Content:   "Comment to delete",
        UserID:    user.ID,
    }
    db.Create(&comment)

    // 创建其他人的评论
    otherComment := models.Comment{
        NewsID:    99999,
        Content:   "Other user's comment",
        UserID:    999, // 不同用户
    }
    db.Create(&otherComment)

    // 为演示子评论，创建子评论指向 comment
    childComment := models.Comment{
        NewsID:    99999,
        Content:   "Child comment",
        UserID:    user.ID,
        ParentID:  &comment.ID,
    }
    db.Create(&childComment)

    tests := []struct {
        name           string
        userID         uint
        commentID      string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid Comment ID",
            userID:         user.ID,
            commentID:      "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid comment ID",
        },
        {
            name:           "Comment Not Found",
            userID:         user.ID,
            commentID:      "88888",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "Comment not found",
        },
        {
            name:           "No Permission (forbidden)",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", otherComment.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusForbidden,
            expectedError:  "You do not have permission to delete this comment",
        },
        {
            name:           "Failed To Delete Comment (transaction error)",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                // 注入事务错误 / delete错误
                db.Callback().Delete().Before("gorm:delete").Register("force_delete_comment_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "comments" {
                        tx.Error = fmt.Errorf("forced delete comment error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to delete comment",
        },
        {
            name:           "Success Delete Comment",
            userID:         user.ID,
            commentID:      fmt.Sprintf("%d", comment.ID),
            setupFunc: func() {
                // 移除回调
                db.Callback().Delete().Remove("force_delete_comment_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := "/news/comments/" + tc.commentID
            req, _ := http.NewRequest("DELETE", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "Comment deleted successfully", resp["message"])
                // 可以在此处再查数据库，确认 comment 以及其子评论已被删除
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestLikeNews(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.News{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.POST("/:id/like", newsController.LikeNews)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_LikeNews_User",
        Nickname: "LikeNewsUser",
    }
    db.Create(&user)

    // 创建新闻
    newsItem := models.News{
        Title:      "Likable News",
        LikeCount:  5,
    }
    db.Create(&newsItem)

    tests := []struct {
        name           string
        userID         uint
        newsID         string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        finalLikeCount int
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid News ID",
            userID:         user.ID,
            newsID:         "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid news ID",
        },
        {
            name:           "News Not Found",
            userID:         user.ID,
            newsID:         "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "News not found",
        },
        {
            name:           "Failed To Find News (simulate DB error)",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_find_news_err_like", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced find news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find news",
        },
        {
            name:           "Failed To Find User",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_find_news_err_like")
                db.Callback().Query().Before("gorm:query").Register("force_find_user_err_like", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced find user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find user",
        },
        {
            name:           "Already Liked",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_find_user_err_like")
                // 让 user 已经点赞
                db.Model(&user).Association("LikedNews").Append(&newsItem)
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "You have already liked this news",
        },
        {
            name:           "Failed To Like News (association append error)",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                // 移除点赞关系
                db.Model(&user).Association("LikedNews").Clear()
                // 注入 append 错误
                db.Callback().Update().Before("gorm:association").Register("force_append_like_news_err", func(tx *gorm.DB) {
                    tx.Error = fmt.Errorf("forced append error")
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to like news",
        },
        {
            name:           "Failed To Update LikeCount",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_append_like_news_err")
                // mock updateColumn 出错
                db.Callback().Update().Before("gorm:update").Register("force_update_likecount_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced update like_count error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update like count",
        },
        {
            name:           "Success Like News",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_update_likecount_err")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            finalLikeCount: newsItem.LikeCount + 1,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := fmt.Sprintf("/news/%s/like", tc.newsID)
            req, _ := http.NewRequest("POST", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "News liked successfully", resp["message"])
                assert.Equal(t, float64(tc.finalLikeCount), resp["like_count"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestCancelLikeNews(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.News{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.DELETE("/:id/like", newsController.CancelLikeNews)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_CancelLikeNews_User",
        Nickname: "CancelLikeNewsUser",
    }
    db.Create(&user)

    // 创建新闻
    newsItem := models.News{
        Title:     "LikeRemovableNews",
        LikeCount: 10,
    }
    db.Create(&newsItem)

    // 让 user 已经点赞了该新闻
    db.Model(&user).Association("LikedNews").Append(&newsItem)

    tests := []struct {
        name           string
        userID         uint
        newsID         string
        setupFunc      func()
        expectedStatus int
        expectedError  string
        isSuccess      bool
        finalLikeCount int
    }{
        {
            name:           "Unauthorized",
            userID:         0,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid News ID",
            userID:         user.ID,
            newsID:         "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid news ID",
        },
        {
            name:           "News Not Found",
            userID:         user.ID,
            newsID:         "88888",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "News not found",
        },
        {
            name:           "Failed To Find News (simulate db error)",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_find_news_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced find news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find news",
        },
        {
            name:           "Failed To Find User",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_find_news_err")
                db.Callback().Query().Before("gorm:query").Register("force_find_user_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced find user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find user",
        },
        {
            name:           "User Not Liked This News",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                // 移除查 user 错误
                db.Callback().Query().Remove("force_find_user_err")
                // 先清空点赞，使用户没有点赞
                db.Model(&user).Association("LikedNews").Clear()
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "You have not liked this news",
        },
        {
            name:           "Failed To Update LikeCount",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Model(&user).Association("LikedNews").Append(&newsItem)
                db.Callback().Update().Remove("force_cancel_like_err")
                // 注入 updateColumn 错误
                db.Callback().Update().Before("gorm:update").Register("force_update_likecount_err_cancel", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced update like_count error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update like count",
        },
        {
            name:           "Success Cancel Like News",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_update_likecount_err_cancel")
            },
            expectedStatus: http.StatusOK,
            isSuccess:      true,
            finalLikeCount: newsItem.LikeCount - 1,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := fmt.Sprintf("/news/%s/like", tc.newsID)
            req, _ := http.NewRequest("DELETE", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "News like canceled successfully", resp["message"])
                assert.Equal(t, float64(tc.finalLikeCount), resp["like_count"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestFavoriteNews(t *testing.T) {
    db := setupNewsTestDB()
    // 需要迁移 News, User
    db.AutoMigrate(&models.News{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.POST("/:id/favorite", newsController.FavoriteNews)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_FavoriteNews_User",
        Nickname: "FavoriteNewsUser",
    }
    db.Create(&user)

    // 创建新闻
    newsItem := models.News{
        Title:         "Favoritable News",
        FavoriteCount: 2,
    }
    db.Create(&newsItem)

    tests := []struct {
        name            string
        userID          uint
        newsID          string
        setupFunc       func()
        expectedStatus  int
        expectedError   string
        isSuccess       bool
        finalFavoriteCT int
    }{
        {
            name:           "Unauthorized (no token)",
            userID:         0,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid News ID",
            userID:         user.ID,
            newsID:         "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid news ID",
        },
        {
            name:           "News Not Found",
            userID:         user.ID,
            newsID:         "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "News not found",
        },
        {
            name:           "Failed To Find News (simulate DB error)",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                // 在查询 news 时注入错误
                db.Callback().Query().Before("gorm:query").Register("force_find_news_err_fav", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced find news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find news",
        },
        {
            name:           "Failed To Find User",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                // 移除上一个回调
                db.Callback().Query().Remove("force_find_news_err_fav")

                // 在查询 user 时注入错误
                db.Callback().Query().Before("gorm:query").Register("force_find_user_err_fav", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced find user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find user",
        },
        {
            name:           "Already Favorited",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_find_user_err_fav")
                // 让 user 已经收藏
                db.Model(&user).Association("FavoritedNews").Append(&newsItem)
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "You have already favorited this news",
        },
        {
            name:           "Association Append Error",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                // 移除已收藏关系
                db.Model(&user).Association("FavoritedNews").Clear()

                // 注入 association 操作错误
                db.Callback().Update().Before("gorm:association").Register("force_favorite_append_err", func(tx *gorm.DB) {
                    tx.Error = fmt.Errorf("forced favorite association error")
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to favorite news",
        },
        {
            name:           "Failed To Update Favorite Count",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_favorite_append_err")

                // mock updateColumn 出错
                db.Callback().Update().Before("gorm:update").Register("force_favorite_update_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced update favorite_count error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update favorite count",
        },
        {
            name:            "Success Favorite",
            userID:          user.ID,
            newsID:          fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_favorite_update_err")
            },
            expectedStatus:  http.StatusOK,
            isSuccess:       true,
            finalFavoriteCT: newsItem.FavoriteCount + 1,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := fmt.Sprintf("/news/%s/favorite", tc.newsID)
            req, _ := http.NewRequest("POST", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "News favorited successfully", resp["message"])
                assert.Equal(t, float64(tc.finalFavoriteCT), resp["favorite_count"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}

func TestCancelFavoriteNews(t *testing.T) {
    db := setupNewsTestDB()
    db.AutoMigrate(&models.News{}, &models.User{})

    router := gin.Default()
    newsController := NewNewsController(db)

    newsGroup := router.Group("/news")
    newsGroup.Use(middleware.AuthMiddleware())
    {
        newsGroup.DELETE("/:id/favorite", newsController.CancelFavoriteNews)
    }

    // 创建用户
    user := models.User{
        OpenID:   "OpenID_CancelFavoriteNews_User",
        Nickname: "CancelFavoriteNewsUser",
    }
    db.Create(&user)

    // 创建新闻
    newsItem := models.News{
        Title:         "Cancelable Favorite News",
        FavoriteCount: 5,
    }
    db.Create(&newsItem)

    // 让 user 已经收藏了该新闻
    db.Model(&user).Association("FavoritedNews").Append(&newsItem)

    tests := []struct {
        name            string
        userID          uint
        newsID          string
        setupFunc       func()
        expectedStatus  int
        expectedError   string
        isSuccess       bool
        finalFavoriteCT int
    }{
        {
            name:           "Unauthorized",
            userID:         0,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc:      func() {},
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Authorization header missing",
        },
        {
            name:           "Invalid News ID",
            userID:         user.ID,
            newsID:         "abc",
            setupFunc:      func() {},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid news ID",
        },
        {
            name:           "News Not Found",
            userID:         user.ID,
            newsID:         "99999",
            setupFunc:      func() {},
            expectedStatus: http.StatusNotFound,
            expectedError:  "News not found",
        },
        {
            name:           "Failed To Find News",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Before("gorm:query").Register("force_find_news_err_cancel_fav", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced find news error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find news",
        },
        {
            name:           "Failed To Find User",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Query().Remove("force_find_news_err_cancel_fav")
                db.Callback().Query().Before("gorm:query").Register("force_find_user_err_cancel_fav", func(tx *gorm.DB) {
                    if tx.Statement.Table == "users" {
                        tx.Error = fmt.Errorf("forced find user error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to find user",
        },
        {
            name:           "User Not Favorited This News",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                // 移除 user 错误回调
                db.Callback().Query().Remove("force_find_user_err_cancel_fav")
                // 先清空收藏
                db.Model(&user).Association("FavoritedNews").Clear()
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "You have not favorited this news",
        },
        {
            name:           "Failed To Update FavoriteCount",
            userID:         user.ID,
            newsID:         fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Model(&user).Association("FavoritedNews").Append(&newsItem)
                db.Callback().Update().Remove("force_cancel_favorite_assoc_err")
                // 注入 updateColumn 错误
                db.Callback().Update().Before("gorm:update").Register("force_favorite_count_err", func(tx *gorm.DB) {
                    if tx.Statement.Table == "news" {
                        tx.Error = fmt.Errorf("forced favorite_count update error")
                    }
                })
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to update favorite count",
        },
        {
            name:            "Success Cancel Favorite",
            userID:          user.ID,
            newsID:          fmt.Sprintf("%d", newsItem.ID),
            setupFunc: func() {
                db.Callback().Update().Remove("force_favorite_count_err")
            },
            expectedStatus:  http.StatusOK,
            isSuccess:       true,
            finalFavoriteCT: newsItem.FavoriteCount - 1,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.setupFunc()

            url := fmt.Sprintf("/news/%s/favorite", tc.newsID)
            req, _ := http.NewRequest("DELETE", url, nil)
            if tc.userID != 0 {
                req.Header.Set("Authorization", "Bearer "+generateValidJWTNews(tc.userID))
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)

            var resp map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &resp)
            assert.NoError(t, err)

            if tc.isSuccess {
                assert.Equal(t, "News favorite canceled successfully", resp["message"])
                assert.Equal(t, float64(tc.finalFavoriteCT), resp["favorite_count"])
            } else {
                if tc.expectedError != "" {
                    assert.Equal(t, tc.expectedError, resp["error"])
                }
            }
        })
    }
}