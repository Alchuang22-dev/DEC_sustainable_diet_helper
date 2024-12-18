// controllers/news_controller_test.go
package controllers

import (
    "bytes"
    "encoding/json"
    // "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
    // "time"

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
    if err := db.AutoMigrate(&models.User{}, &models.Draft{}, &models.DraftParagraph{}, &models.DraftImage{}); err != nil {
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
            newsGroup.POST("/upload_image", newsController.UploadImage)      // 上传单张图片
            newsGroup.POST("/create_draft", newsController.CreateDraft)     // 创建草稿
            // 其他路由...
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

// TestUploadImage_Success 测试成功上传图片
func TestUploadImage_Success(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 创建临时图片文件
    imagePath := "./test_upload_image.jpg"
    os.WriteFile(imagePath, []byte("test image content"), 0644)
    defer os.Remove(imagePath)

    // 构建 multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
    assert.NoError(t, err)

    file, err := os.Open(imagePath)
    assert.NoError(t, err)
    defer file.Close()

    _, err = io.Copy(part, file)
    assert.NoError(t, err)
    writer.Close()

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/upload_image", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusOK, w.Code)
    var response map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Image uploaded successfully", response["message"])
    assert.NotEmpty(t, response["path"])

    // 检查文件是否存在
    relativePath := response["path"].(string)
    savedPath := filepath.Join("./uploads", relativePath)
    _, err = os.Stat(savedPath)
    assert.False(t, os.IsNotExist(err))

    // 清理上传的文件
    os.Remove(savedPath)
}

// TestUploadImage_NoImage 测试未上传图片
func TestUploadImage_NoImage(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 构建 multipart form，不包含图片
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.Close()

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/upload_image", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusBadRequest, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Image file is required", response["error"])
}

// TestUploadImage_Unauthorized 测试未授权上传图片
func TestUploadImage_Unauthorized(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 构建 multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    writer.WriteField("image", "test_image_content")
    writer.Close()

    // 创建请求，不设置 Authorization 头
    req, _ := http.NewRequest("POST", "/news/upload_image", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Unauthorized", response["error"])
}

// TestCreateDraft1_Success 测试成功创建草稿
func TestCreateDraft1_Success(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 创建段落和图片描述
    paragraphs := []string{"First paragraph", "Second paragraph"}
    imageDescriptions := []string{"Image 1 description", "Image 2 description"}

    // 假设图片已经上传，并获得相对路径
    imagePaths := []string{
        "drafts/1/image1.jpg",
        "drafts/1/image2.jpg",
    }

    // 构建请求体
    requestBody := map[string]interface{}{
        "title":              "Test Draft",
        "paragraphs":         paragraphs,
        "image_descriptions": imageDescriptions,
        "image_paths":        imagePaths,
    }
    requestJSON, _ := json.Marshal(requestBody)

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(requestJSON))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusCreated, w.Code)
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Draft created successfully.", response["message"])
    assert.NotNil(t, response["draft_id"])

    // 检查数据库中是否存在草稿
    var draft models.Draft
    draftID := uint(response["draft_id"].(float64))
    result := db.Preload("Paragraphs").Preload("Images").First(&draft, draftID)
    assert.Nil(t, result.Error)
    assert.Equal(t, "Test Draft", draft.Title)
    assert.Equal(t, user.ID, draft.AuthorID)
    assert.Len(t, draft.Paragraphs, 2)
    assert.Len(t, draft.Images, 2)
    assert.Equal(t, "Image 1 description", draft.Images[0].Description)
    assert.Equal(t, "Image 2 description", draft.Images[1].Description)
}

// TestCreateDraft1_MissingTitle 测试缺少标题
func TestCreateDraft1_MissingTitle(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 创建段落和图片描述
    paragraphs := []string{"First paragraph"}
    imageDescriptions := []string{"Image 1 description"}

    // 假设图片已经上传，并获得相对路径
    imagePaths := []string{
        "drafts/1/image1.jpg",
    }

    // 构建请求体，缺少 title
    requestBody := map[string]interface{}{
        "paragraphs":         paragraphs,
        "image_descriptions": imageDescriptions,
        "image_paths":        imagePaths,
    }
    requestJSON, _ := json.Marshal(requestBody)

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(requestJSON))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusBadRequest, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Title is required", response["error"])
}

// TestCreateDraft1_MismatchedImageDescriptions 测试图片描述与图片路径数量不匹配
func TestCreateDraft1_MismatchedImageDescriptions(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 创建段落和图片描述
    paragraphs := []string{"First paragraph"}
    imageDescriptions := []string{"Image 1 description"}

    // 假设上传了两张图片，但只有一个描述
    imagePaths := []string{
        "drafts/1/image1.jpg",
        "drafts/1/image2.jpg",
    }

    // 构建请求体
    requestBody := map[string]interface{}{
        "title":              "Test Draft",
        "paragraphs":         paragraphs,
        "image_descriptions": imageDescriptions,
        "image_paths":        imagePaths,
    }
    requestJSON, _ := json.Marshal(requestBody)

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(requestJSON))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusBadRequest, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Number of image descriptions and image paths do not match", response["error"])
}

// TestCreateDraft1_InvalidRequestData 测试无效的请求数据
func TestCreateDraft1_InvalidRequestData(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 构建无效的 JSON 请求体
    invalidJSON := `{"title": "Test Draft", "paragraphs": "invalid_paragraphs"}`

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBufferString(invalidJSON))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusBadRequest, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Invalid request data", response["error"])
}

// TestCreateDraft1_Unauthorized 测试未授权创建草稿
func TestCreateDraft1_Unauthorized(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 构建请求体
    requestBody := map[string]interface{}{
        "title":              "Test Draft",
        "paragraphs":         []string{"First paragraph"},
        "image_descriptions": []string{"Image 1 description"},
        "image_paths":        []string{"drafts/1/image1.jpg"},
    }
    requestJSON, _ := json.Marshal(requestBody)

    // 创建请求，不设置 Authorization 头
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(requestJSON))
    req.Header.Set("Content-Type", "application/json")

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Unauthorized", response["error"])
}

// TestCreateDraft1_UserNotFound 测试用户不存在
func TestCreateDraft1_UserNotFound(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 生成 JWT 对应不存在的用户 ID
    fakeUserID := uint(9999)
    token := generateValidJWTNews(fakeUserID)

    // 构建请求体
    requestBody := map[string]interface{}{
        "title":              "Test Draft",
        "paragraphs":         []string{"First paragraph"},
        "image_descriptions": []string{"Image 1 description"},
        "image_paths":        []string{"drafts/9999/image1.jpg"},
    }
    requestJSON, _ := json.Marshal(requestBody)

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(requestJSON))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    // 根据您的实现，可能返回 401 Unauthorized
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Unauthorized", response["error"])
}

// TestCreateDraft1_DatabaseFailure 测试数据库保存失败
func TestCreateDraft1_DatabaseFailure(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 构建请求体
    requestBody := map[string]interface{}{
        "title":              "Test Draft",
        "paragraphs":         []string{"First paragraph"},
        "image_descriptions": []string{"Image 1 description"},
        "image_paths":        []string{"drafts/1/image1.jpg"},
    }
    requestJSON, _ := json.Marshal(requestBody)

    // 关闭数据库连接，模拟数据库错误
    sqlDB, _ := db.DB()
    sqlDB.Close()

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(requestJSON))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 检查响应
    assert.Equal(t, http.StatusInternalServerError, w.Code)
    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Failed to create draft", response["error"])
}

// TestCreateDraft1_InvalidImagePath 测试无效的图片路径
func TestCreateDraft1_InvalidImagePath(t *testing.T) {
    db := setupNewsTestDB()
    router := setupNewsRouter(db)

    // 创建用户
    user := models.User{
        Nickname:  "TestUser",
        AvatarURL: "avatars/default.jpg",
    }
    db.Create(&user)

    // 生成 JWT
    token := generateValidJWTNews(user.ID)

    // 构建请求体，提供不存在的图片路径
    requestBody := map[string]interface{}{
        "title":              "Test Draft",
        "paragraphs":         []string{"First paragraph"},
        "image_descriptions": []string{"Image 1 description"},
        "image_paths":        []string{"invalid_path/image1.jpg"},
    }
    requestJSON, _ := json.Marshal(requestBody)

    // 创建请求
    req, _ := http.NewRequest("POST", "/news/create_draft", bytes.NewBuffer(requestJSON))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 根据您的实现，可能仍然成功创建草稿，图片路径只是相对路径，不验证实际文件是否存在
    // 如果需要验证图片路径存在，可以在 CreateDraft1 中添加相应逻辑
    assert.Equal(t, http.StatusCreated, w.Code)
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Draft created successfully.", response["message"])
    assert.NotNil(t, response["draft_id"])
}