// controllers/news_controller.go
// internal/controllers/news_controller.go
package controllers

import (
    "mime/multipart"
    // "encoding/json"
    "path/filepath"
    "net/http"
    "strings"
    "strconv"
    "errors"
    "time"
    "fmt"
    "os"
    "io"
    // "bytes"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Alchuang22-dev/DEC_sustainable_diet_helper/internal/models"
)

type NewsController struct {
    DB *gorm.DB
}

func NewNewsController(db *gorm.DB) *NewsController {
    return &NewsController{DB: db}
}

// CreateDraftRequest 定义创建草稿的请求结构
type CreateDraftRequest struct {
    Title      string               `json:"title" binding:"required"`
    Paragraphs []models.DraftParagraph `json:"paragraphs,omitempty"`
    // Images 和 Descriptions 将通过 multipart/form-data 处理
}

// CreateDraft 详细创建草稿
func (nc *NewsController) CreateDraft(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 解析 JSON 请求体
    var request struct {
        Title             string   `json:"title"`
        Paragraphs        []string `json:"paragraphs"`
        ImageDescriptions []string `json:"image_descriptions"`
        ImagePaths        []string `json:"image_paths"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // 校验必填字段
    if request.Title == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
        return
    }

    // 校验图片描述和图片路径数量是否匹配
    if len(request.ImageDescriptions) != len(request.ImagePaths) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Number of image descriptions and image paths do not match"})
        return
    }

    // 构建段落对象
    var paragraphs []models.DraftParagraph
    for _, text := range request.Paragraphs {
        paragraphs = append(paragraphs, models.DraftParagraph{
            Text: text,
        })
    }

    // 构建图片对象
    var draftImages []models.DraftImage
    for i, path := range request.ImagePaths {
        draftImages = append(draftImages, models.DraftImage{
            URL:         path,
            Description: request.ImageDescriptions[i],
        })
    }

    // 初始化草稿对象
    draft := models.Draft{
        Title:      request.Title,
        AuthorID:   userID.(uint),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        Paragraphs: paragraphs,
        Images:     draftImages,
    }

    // 插入数据库
    if err := nc.DB.Create(&draft).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create draft"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":  "Draft created successfully.",
        "draft_id": draft.ID,
    })
}

// UploadImage 处理单张图片上传，并返回图片的相对路径
func (nc *NewsController) UploadImage(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取图片文件
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
        return
    }

    // 生成保存路径
    relativePath, err := uploadImage(file, userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Image uploaded successfully",
        "path":    relativePath,
    })
}

// uploadImage 处理图片文件上传，并返回图片 URL
func uploadImage(file *multipart.FileHeader, draftID uint) (string, error) {
    // 获取基地址
    baseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
    if baseUploadPath == "" {
        baseUploadPath = "./uploads" // 如果没有设置环境变量，使用默认路径
    }

    // 打开文件
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    // 生成保存路径，确保文件名唯一
    timestamp := time.Now().UnixNano()
    filename := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename))
    relativePath := filepath.Join("drafts", fmt.Sprintf("%d", draftID), filename)
    savePath := filepath.Join(baseUploadPath, relativePath)

    // 创建目标文件
    out, err := createFile(savePath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    // 复制文件内容
    _, err = io.Copy(out, src)
    if err != nil {
        return "", err
    }

    // 返回图片的相对路径
    return relativePath, nil
}

// createFile 创建文件夹并创建目标文件
func createFile(savePath string) (*os.File, error) {
    // 创建文件夹
    dir := filepath.Dir(savePath)
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return nil, err
    }

    // 创建文件
    return os.Create(savePath)
}

// ConvertDraftToNews 处理将草稿转换为新闻的请求
func (nc *NewsController) ConvertDraftToNews(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    type ConvertDraftToNewsRequest struct {
        DraftID uint `json:"draft_id" binding:"required"`
    }

    // 绑定 JSON 请求体
    var convertRequest ConvertDraftToNewsRequest
    if err := c.ShouldBindJSON(&convertRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // 查找草稿，确保草稿属于当前用户
    var draft models.Draft
    if err := nc.DB.Preload("Paragraphs").Preload("Images").First(&draft, "id = ? AND author_id = ?", convertRequest.DraftID, userID.(uint)).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find draft"})
        return
    }

    // 初始化新闻对象
    news := models.News{
        Title:           draft.Title,
        AuthorID:        draft.AuthorID,
        UploadTime:      time.Now(),
        ViewCount:       0,
        LikeCount:       0,
        FavoriteCount:   0,
        DislikeCount:    0,
        ShareCount:      0,
        LikedByUsers:    []models.User{},
        FavoritedByUsers: []models.User{},
        DislikedByUsers: []models.User{},
        Paragraphs:      []models.Paragraph{},
        Images:          []models.NewsImage{},
    }

    // 复制段落
    for _, p := range draft.Paragraphs {
        news.Paragraphs = append(news.Paragraphs, models.Paragraph{
            Text: p.Text,
        })
    }

    // 复制图片，包括描述
    for _, img := range draft.Images {
        news.Images = append(news.Images, models.NewsImage{
            URL:         img.URL,
            Description: img.Description, // 复制描述
        })
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    // 创建新闻
    if err := tx.Create(&news).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
        return
    }

    // 删除关联的图片数据
    if err := tx.Where("draft_id = ?", draft.ID).Delete(&models.DraftImage{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete draft images"})
        return
    }
    // 删除关联的段落数据
    if err := tx.Where("draft_id = ?", draft.ID).Delete(&models.DraftParagraph{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete draft paragraphs"})
        return
    }

    // 删除草稿及其关联的段落和图片
    if err := tx.Delete(&draft).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete draft"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Draft converted to news successfully.",
        "news_id": news.ID,
    })
}

// UpdateDraft 更新草稿
func (nc *NewsController) UpdateDraft(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取草稿 ID
    draftID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
        return
    }

    // 查找草稿
    var draft models.Draft
    if err := nc.DB.Preload("Images").First(&draft, draftID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
        return
    }

    // 验证是否为草稿作者
    if draft.AuthorID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to edit this draft"})
        return
    }

    // 获取新上传的图片路径列表
    var request struct {
        Title             string   `json:"title"`
        Paragraphs        []string `json:"paragraphs"`
        ImageDescriptions []string `json:"image_descriptions"`
        ImagePaths        []string `json:"image_paths"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // 校验必填字段
    if request.Title == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
        return
    }
    if len(request.ImageDescriptions) != len(request.ImagePaths) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Number of image descriptions and image paths do not match"})
        return
    }

    // 删除旧草稿（包括旧图片的删除）
    if err := nc.DB.Transaction(func(tx *gorm.DB) error {
        // 删除关联段落
        if err := tx.Where("draft_id = ?", draft.ID).Delete(&models.DraftParagraph{}).Error; err != nil {
            return err
        }

        // 删除关联图片及本地文件（检查图片路径是否在上传列表中）
        for _, image := range draft.Images {
            filePath := filepath.Join(os.Getenv("BASE_UPLOAD_PATH"), image.URL)

            // 如果该图片的路径不在新的图片路径列表中，则删除该图片
            found := false
            for _, newImagePath := range request.ImagePaths {
                if image.URL == newImagePath {
                    found = true
                    break
                }
            }

            // 如果图片路径不在新上传的路径中，才删除它
            if !found {
                if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
                    return fmt.Errorf("failed to delete file: %s", filePath)
                }
            }
        }

        // 删除草稿图片记录
        if err := tx.Where("draft_id = ?", draft.ID).Delete(&models.DraftImage{}).Error; err != nil {
            return err
        }

        // 删除草稿本体
        if err := tx.Delete(&draft).Error; err != nil {
            return err
        }

        return nil
    }); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old draft"})
        return
    }

    // 构建新草稿并保存
    var paragraphs []models.DraftParagraph
    for _, text := range request.Paragraphs {
        paragraphs = append(paragraphs, models.DraftParagraph{
            Text: text,
        })
    }

    // 构建新上传的图片对象
    var draftImages []models.DraftImage
    for i, path := range request.ImagePaths {
        draftImages = append(draftImages, models.DraftImage{
            URL:         path,
            Description: request.ImageDescriptions[i],
        })
    }

    // 创建新草稿
    newDraft := models.Draft{
        Title:      request.Title,
        AuthorID:   userID.(uint),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        Paragraphs: paragraphs,
        Images:     draftImages,
    }

    if err := nc.DB.Create(&newDraft).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new draft"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":  "Draft updated successfully",
        "draft_id": newDraft.ID,
    })
}

// DeleteDraft 删除草稿及其图片文件
func (nc *NewsController) DeleteDraft(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取草稿 ID
    draftID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
        return
    }

    // 查找草稿
    var draft models.Draft
    if err := nc.DB.Preload("Images").First(&draft, draftID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
        return
    }

    // 验证是否为草稿作者
    if draft.AuthorID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this draft"})
        return
    }

    // 删除草稿记录及其关联的段落和图片
    if err := nc.DB.Transaction(func(tx *gorm.DB) error {
        // 删除关联的段落
        if err := tx.Where("draft_id = ?", draft.ID).Delete(&models.DraftParagraph{}).Error; err != nil {
            return err
        }

        // 删除关联的图片，同时删除本地文件
        for _, image := range draft.Images {
            if err := deleteLocalFile(image.URL); err != nil {
                fmt.Printf("Failed to delete local file: %s, error: %v\n", image.URL, err)
            }
        }

        if err := tx.Where("draft_id = ?", draft.ID).Delete(&models.DraftImage{}).Error; err != nil {
            return err
        }

        // 删除草稿
        if err := tx.Delete(&draft).Error; err != nil {
            return err
        }

        return nil
    }); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete draft"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Draft deleted successfully."})
}

// DeleteNews 删除新闻
func (nc *NewsController) DeleteNews(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 查找新闻
    var news models.News
    if err := nc.DB.Preload("Images").First(&news, newsID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 验证是否为新闻作者
    if news.AuthorID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this news"})
        return
    }

    // 开始事务
    if err := nc.DB.Transaction(func(tx *gorm.DB) error {
        // 清除 LikedByUsers 关联
        if err := tx.Model(&news).Association("LikedByUsers").Clear(); err != nil {
            return err
        }
        // 清除 FavoritedByUsers 关联
        if err := tx.Model(&news).Association("FavoritedByUsers").Clear(); err != nil {
            return err
        }
        // 清除 DislikedByUsers 关联
        if err := tx.Model(&news).Association("DislikedByUsers").Clear(); err != nil {
            return err
        }
        if err := tx.Model(&news).Association("ViewedByUsers").Clear(); err != nil {
            return err
        }
        
        // 删除关联的段落
        if err := tx.Where("news_id = ?", news.ID).Delete(&models.Paragraph{}).Error; err != nil {
            return err
        }

        // 删除关联的图片
        if err := tx.Where("news_id = ?", news.ID).Delete(&models.NewsImage{}).Error; err != nil {
            return err
        }

        // **新增**：删除与该新闻关联的所有评论
        if err := tx.Where("news_id = ?", news.ID).Delete(&models.Comment{}).Error; err != nil {
            return err
        }

        // 删除本地图片文件
        for _, image := range news.Images {
            if err := deleteLocalFile(image.URL); err != nil {
                // 记录日志，但不终止事务
                fmt.Printf("Failed to delete local file: %s, error: %v\n", image.URL, err)
            }
        }

        // 删除新闻
        if err := tx.Delete(&news).Error; err != nil {
            return err
        }

        return nil
    }); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully."})
}

// deleteLocalFile 删除本地文件
func deleteLocalFile(filePath string) error {
    baseUploadPath := os.Getenv("BASE_UPLOAD_PATH")
    if baseUploadPath == "" {
        baseUploadPath = "./uploads" // 如果没有设置环境变量，使用默认路径
    }

    // 拼接完整路径
    fullPath := filepath.Join(baseUploadPath, filePath)
    if err := os.Remove(fullPath); err != nil {
        if os.IsNotExist(err) {
            // 如果文件不存在，不认为是错误
            return nil
        }
        return err
    }

    return nil
}

// GetMyNews 获取用户创建的新闻 ID 列表
func (nc *NewsController) GetMyNews(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var newsList []models.News
    if err := nc.DB.Select("id").Where("author_id = ?", userID).Order("upload_time DESC").Find(&newsList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
        return
    }

    // 提取 ID 列表
    ids := make([]uint, len(newsList))
    for i, news := range newsList {
        ids[i] = news.ID
    }

    c.JSON(http.StatusOK, gin.H{
        "news_ids": ids,
    })
}

// GetMyDrafts 获取用户创建的草稿 ID 列表
func (nc *NewsController) GetMyDrafts(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var draftList []models.Draft
    if err := nc.DB.Select("id").Where("author_id = ?", userID).Order("updated_at DESC").Find(&draftList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drafts"})
        return
    }

    // 提取 ID 列表
    ids := make([]uint, len(draftList))
    for i, draft := range draftList {
        ids[i] = draft.ID
    }

    c.JSON(http.StatusOK, gin.H{
        "draft_ids": ids,
    })
}

// PreviewNewsRequest 定义预览新闻的请求结构
type PreviewNewsRequest struct {
    IDs []uint `json:"ids" binding:"required"`
}

// PreviewNewsResponse 定义预览新闻的响应结构
type PreviewNewsResponse struct {
    Previews []PreviewNewsItem `json:"previews"`
}

// PreviewNewsItem 定义单个新闻预览的结构
type PreviewNewsItem struct {
    ID                 uint   `json:"id"`
    Title              string `json:"title"`
    FirstParagraphText string `json:"first_paragraph_text"`
    FirstImageURL      string `json:"first_image_url"`
    ImageDescription   string `json:"image_description"`      // 新增描述字段
    AuthorNickname     string `json:"author_nickname"`
    AuthorAvatarURL    string `json:"author_avatar_url"`
    LikeCount          int    `json:"like_count"`
}

// PreviewNews 处理新闻预览的请求
func (nc *NewsController) PreviewNews(c *gin.Context) {
    var req PreviewNewsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var newsList []models.News
    if err := nc.DB.Preload("Author").Preload("Paragraphs").Preload("Images").
        Where("id IN ?", req.IDs).Find(&newsList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
        return
    }

    previews := make([]PreviewNewsItem, 0, len(newsList))
    for _, news := range newsList {
        firstParagraphText := ""
        if len(news.Paragraphs) > 0 {
            firstParagraphText = news.Paragraphs[0].Text
        }

        firstImageURL := ""
        imageDescription := ""
        if len(news.Images) > 0 {
            firstImageURL = news.Images[0].URL
            imageDescription = news.Images[0].Description // 获取描述
        }

        preview := PreviewNewsItem{
            ID:                 news.ID,
            Title:              news.Title,
            FirstParagraphText: firstParagraphText,
            FirstImageURL:      firstImageURL,
            ImageDescription:   imageDescription, // 设置描述
            AuthorNickname:     news.Author.Nickname,
            AuthorAvatarURL:    news.Author.AvatarURL,
            LikeCount:          news.LikeCount,
        }
        previews = append(previews, preview)
    }

    c.JSON(http.StatusOK, PreviewNewsResponse{
        Previews: previews,
    })
}

// PreviewDraftRequest 定义预览草稿的请求结构
type PreviewDraftRequest struct {
    IDs []uint `json:"ids" binding:"required"`
}

// PreviewDraftResponse 定义预览草稿的响应结构
type PreviewDraftResponse struct {
    Previews []PreviewDraftItem `json:"previews"`
}

// PreviewDraftItem 定义单个草稿预览的结构
type PreviewDraftItem struct {
    ID                 uint   `json:"id"`
    Title              string `json:"title"`
    FirstParagraphText string `json:"first_paragraph_text"`
    FirstImageURL      string `json:"first_image_url"`
    ImageDescription   string `json:"image_description"` // 新增描述字段
}

// PreviewDrafts 处理草稿预览的请求
func (nc *NewsController) PreviewDrafts(c *gin.Context) {
    var req PreviewDraftRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var drafts []models.Draft
    if err := nc.DB.Preload("Paragraphs").Preload("Images").
        Where("id IN ?", req.IDs).Find(&drafts).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drafts"})
        return
    }

    previews := make([]PreviewDraftItem, 0, len(drafts))
    for _, draft := range drafts {
        firstParagraphText := ""
        if len(draft.Paragraphs) > 0 {
            firstParagraphText = draft.Paragraphs[0].Text
        }

        firstImageURL := ""
        imageDescription := ""
        if len(draft.Images) > 0 {
            firstImageURL = draft.Images[0].URL
            imageDescription = draft.Images[0].Description // 获取描述
        }

        preview := PreviewDraftItem{
            ID:                 draft.ID,
            Title:              draft.Title,
            FirstParagraphText: firstParagraphText,
            FirstImageURL:      firstImageURL,
            ImageDescription:   imageDescription, // 设置描述
        }
        previews = append(previews, preview)
    }

    c.JSON(http.StatusOK, PreviewDraftResponse{
        Previews: previews,
    })
}

// AuthorInfo 定义作者信息的结构
type AuthorInfo struct {
    ID         uint   `json:"id"`
    Nickname   string `json:"nickname"`
    AvatarURL  string `json:"avatar_url"`
}

// NewsDetailResponse 定义详细查看新闻的响应结构
type NewsDetailResponse struct {
    ID              uint          `json:"id"`
    Title           string        `json:"title"`
    UploadTime      time.Time     `json:"upload_time"`
    ViewCount       int           `json:"view_count"`
    LikeCount       int           `json:"like_count"`
    FavoriteCount   int           `json:"favorite_count"`
    DislikeCount    int           `json:"dislike_count"`
    ShareCount      int           `json:"share_count"`
    Author          AuthorInfo    `json:"author"`
    Paragraphs      []models.Paragraph `json:"paragraphs"`
    Images          []models.NewsImage `json:"images"`
    Comments        []models.Comment   `json:"comments"`
}

// GetNewsDetails 详细查看单个新闻
func (nc *NewsController) GetNewsDetails(c *gin.Context) {
    // 获取用户 ID
    userIDValue, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDValue.(uint)

    // 获取新闻 ID
    newsIDStr := c.Param("id")
    newsID, err := strconv.ParseUint(newsIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 查找新闻
    var news models.News
    if err := nc.DB.
        Preload("Author").
        Preload("Paragraphs").
        Preload("Images").
        First(&news, "id = ?", newsID).
        Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news details"})
        return
    }

    // 获取顶级评论
    var topLevelComments []models.Comment
    if err := nc.DB.Preload("Replies").
        Preload("Author", func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "nickname", "avatar_url")
        }).
        Where("news_id = ? AND parent_id IS NULL", newsID).
        Find(&topLevelComments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
        return
    }

    // 先查找用户所有点赞过的评论 ID
    likedCommentsMap, err := nc.fetchUserLikedCommentMap(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's liked comments"})
        return
    }

    // 构建递归评论结构,并标记 did_like
    comments := nc.buildCommentTree(topLevelComments, likedCommentsMap)

    response := NewsDetailResponse{
        ID:            news.ID,
        Title:         news.Title,
        UploadTime:    news.UploadTime,
        ViewCount:     news.ViewCount,
        LikeCount:     news.LikeCount,
        FavoriteCount: news.FavoriteCount,
        DislikeCount:  news.DislikeCount,
        ShareCount:    news.ShareCount,
        Author: AuthorInfo{
            ID:        news.Author.ID,
            Nickname:  news.Author.Nickname,
            AvatarURL: news.Author.AvatarURL,
        },
        Paragraphs: news.Paragraphs,
        Images:     news.Images,
        Comments:   comments,
    }

    c.JSON(http.StatusOK, response)
}

// fetchUserLikedCommentMap 获取当前用户点赞过的评论，做成 map[评论ID] = true
func (nc *NewsController) fetchUserLikedCommentMap(userID uint) (map[uint]bool, error) {
    var user models.User
    // 预加载 LikedComments
    if err := nc.DB.Preload("LikedComments").First(&user, userID).Error; err != nil {
        return nil, err
    }

    result := make(map[uint]bool)
    for _, c := range user.LikedComments {
        result[c.ID] = true
    }
    return result, nil
}

// buildCommentTree 构建递归评论结构，并标记 did_like 字段
func (nc *NewsController) buildCommentTree(comments []models.Comment, likedCommentsMap map[uint]bool) []models.Comment {
    for i := range comments {
        // 判断该评论是否被当前用户点赞
        if likedCommentsMap[comments[i].ID] {
            // 可以自己给 Comment 模型新增字段 DidLike bool, 
            // 但是也可以返回到 JSON 里，灵活处理：
            comments[i].DidLike = true // 需要在 Comment 模型中加 DidLike 字段
        }

        // 继续递归处理 replies
        var replies []models.Comment
        if err := nc.DB.
            Preload("Replies").
            Preload("Author", func(db *gorm.DB) *gorm.DB {
                return db.Select("id", "nickname", "avatar_url")
            }).
            Where("parent_id = ?", comments[i].ID).
            Find(&replies).Error; err == nil {
            comments[i].Replies = nc.buildCommentTree(replies, likedCommentsMap)
        }
    }
    return comments
}

// DraftDetailResponse 定义详细查看草稿的响应结构
type DraftDetailResponse struct {
    ID          uint                `json:"id"`
    Title       string              `json:"title"`
    Author      AuthorInfo          `json:"author"`
    CreatedAt   time.Time           `json:"created_at"`
    UpdatedAt   time.Time           `json:"updated_at"`
    Paragraphs  []models.DraftParagraph `json:"paragraphs"`
    Images      []models.DraftImage `json:"images"`
}

// GetDraftDetails 详细查看单个草稿
func (nc *NewsController) GetDraftDetails(c *gin.Context) {
    // 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取草稿 ID
    draftIDStr := c.Param("id")
    draftID, err := strconv.ParseUint(draftIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
        return
    }

    var draft models.Draft
    if err := nc.DB.Preload("Author").Preload("Paragraphs").Preload("Images").
        First(&draft, "id = ? AND author_id = ?", draftID, userID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch draft details"})
        return
    }

    response := DraftDetailResponse{
        ID:       draft.ID,
        Title:    draft.Title,
        Author: AuthorInfo{
            ID:        draft.Author.ID,
            Nickname:  draft.Author.Nickname,
            AvatarURL: draft.Author.AvatarURL,
        },
        CreatedAt:  draft.CreatedAt,
        UpdatedAt:  draft.UpdatedAt,
        Paragraphs: draft.Paragraphs,
        Images:     draft.Images,
    }

    c.JSON(http.StatusOK, response)
}

// GetNewsByViewCount 获取按观看量降序排序的新闻 ID 列表，每页 10 条
func (nc *NewsController) GetNewsByViewCount(c *gin.Context) {
    _, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    pageStr := c.Query("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }

    var newsList []models.News
    if err := nc.DB.Select("id").
        Order("view_count DESC").
        Limit(10).
        Offset((page - 1) * 10).
        Find(&newsList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
        return
    }

    // 提取 ID 列表
    ids := make([]uint, len(newsList))
    for i, news := range newsList {
        ids[i] = news.ID
    }

    c.JSON(http.StatusOK, gin.H{
        "news_ids": ids,
    })
}

// GetNewsByLikeCount 获取按点赞量降序排序的新闻 ID 列表，每页 10 条
func (nc *NewsController) GetNewsByLikeCount(c *gin.Context) {
    _, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    pageStr := c.Query("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }

    var newsList []models.News
    if err := nc.DB.Select("id").
        Order("like_count DESC").
        Limit(10).
        Offset((page - 1) * 10).
        Find(&newsList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
        return
    }

    // 提取 ID 列表
    ids := make([]uint, len(newsList))
    for i, news := range newsList {
        ids[i] = news.ID
    }

    c.JSON(http.StatusOK, gin.H{
        "news_ids": ids,
    })
}

// GetNewsByUploadTime 获取按上传时间降序排序的新闻 ID 列表，每页 10 条
func (nc *NewsController) GetNewsByUploadTime(c *gin.Context) {
    _, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    pageStr := c.Query("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }

    var newsList []models.News
    if err := nc.DB.Select("id").
        Order("upload_time DESC").
        Limit(10).
        Offset((page - 1) * 10).
        Find(&newsList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
        return
    }

    // 提取 ID 列表
    ids := make([]uint, len(newsList))
    for i, news := range newsList {
        ids[i] = news.ID
    }

    c.JSON(http.StatusOK, gin.H{
        "news_ids": ids,
    })
}

func (nc *NewsController) AddComment(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 请求体绑定
    var commentRequest struct {
        NewsID   uint   `json:"news_id" binding:"required"`
        Content  string `json:"content" binding:"required"`
        IsReply  bool   `json:"is_reply"`
        ParentID *uint  `json:"parent_id,omitempty"`
    }
    if err := c.ShouldBindJSON(&commentRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if !commentRequest.IsReply || commentRequest.ParentID == nil || *commentRequest.ParentID == 0 {
        commentRequest.ParentID = nil
    }

    if commentRequest.IsReply && commentRequest.ParentID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ParentID is required for a reply"})
        return
    }

    // 检查新闻是否存在
    var newsExists bool
    if err := nc.DB.Model(&models.News{}).Select("count(*) > 0").Where("id = ?", commentRequest.NewsID).Scan(&newsExists).Error; err != nil || !newsExists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 校验父评论逻辑
    if commentRequest.IsReply {
        if commentRequest.ParentID == nil || *commentRequest.ParentID == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ParentID is required for a reply"})
            return
        }

        var parentComment models.Comment
        if err := nc.DB.First(&parentComment, *commentRequest.ParentID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment not found"})
            return
        }

        if parentComment.NewsID != commentRequest.NewsID {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment does not belong to the specified news"})
            return
        }
    } else {
        if commentRequest.ParentID != nil && *commentRequest.ParentID != 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ParentID must be empty for a top-level comment"})
            return
        }
    }

    // 创建评论对象
    comment := models.Comment{
        NewsID:      commentRequest.NewsID,
        Content:     commentRequest.Content,
        UserID:      userID.(uint),
        IsReply:     commentRequest.IsReply,
        ParentID:    commentRequest.ParentID,
        PublishTime: time.Now(),
        LikeCount:   0,
    }

    // 保存评论
    if err := nc.DB.Create(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
        return
    }

    var user models.User
    if err := nc.DB.Select("nickname, avatar_url").First(&user, userID.(uint)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Comment added successfully",
        "comment": gin.H{
            "id":        comment.ID,
            "news_id":   comment.NewsID,
            "content":   comment.Content,
            "is_reply":  comment.IsReply,
            "parent_id": comment.ParentID,
            "user_id":   comment.UserID,
            "publish_time": comment.PublishTime,
            "like_count":   comment.LikeCount,
            "author": gin.H{
                "nickname":   user.Nickname,
                "avatar_url": user.AvatarURL,
            },
        },
    })
}

// LikeComment 处理用户对评论点赞
func (nc *NewsController) LikeComment(c *gin.Context) {
    // 1. 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. 获取评论 ID
    commentIDStr := c.Param("id")
    commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    // 3. 开启事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 4. 查找评论
    var comment models.Comment
    if err := tx.Preload("LikedByUsers").First(&comment, commentID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find comment"})
        return
    }

    // 5. 查找用户
    var user models.User
    if err := tx.Preload("LikedComments").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 6. 检查是否已点赞
    for _, liked := range user.LikedComments {
        if liked.ID == comment.ID {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "You have already liked this comment"})
            return
        }
    }

    // 7. 建立点赞关系
    if err := tx.Model(&user).Association("LikedComments").Append(&comment); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like comment"})
        return
    }

    // 8. 增加评论点赞数
    if err := tx.Model(&comment).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment like_count"})
        return
    }

    // 9. 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    // 10. 返回成功
    c.JSON(http.StatusOK, gin.H{
        "message":    "Comment liked successfully",
        "like_count": comment.LikeCount + 1,
    })
}

// CancelLikeComment 处理用户取消对评论点赞
func (nc *NewsController) CancelLikeComment(c *gin.Context) {
    // 1. 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. 获取评论 ID
    commentIDStr := c.Param("id")
    commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    // 3. 开启事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 4. 查找评论
    var comment models.Comment
    if err := tx.First(&comment, commentID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find comment"})
        return
    }

    // 5. 查找用户
    var user models.User
    if err := tx.Preload("LikedComments").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 6. 检查是否已点赞
    liked := false
    for _, likedComment := range user.LikedComments {
        if likedComment.ID == comment.ID {
            liked = true
            break
        }
    }
    if !liked {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not liked this comment"})
        return
    }

    // 7. 删除点赞关系
    if err := tx.Model(&user).Association("LikedComments").Delete(&comment); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel like on comment"})
        return
    }

    // 8. 减少评论点赞数
    if err := tx.Model(&comment).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment like_count"})
        return
    }

    // 9. 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    // 10. 返回成功
    c.JSON(http.StatusOK, gin.H{
        "message":    "Comment like canceled successfully",
        "like_count": comment.LikeCount - 1,
    })
}

func (nc *NewsController) DeleteComment(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取评论 ID
    commentID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    // 查找评论
    var comment models.Comment
    if err := nc.DB.First(&comment, commentID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }

    // 验证用户是否为评论作者
    if comment.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this comment"})
        return
    }

    // 删除评论及其子评论
    if err := nc.DB.Transaction(func(tx *gorm.DB) error {
        // 递归删除子评论
        if err := deleteCommentRecursive(tx, commentID); err != nil {
            return err
        }
        return nil
    }); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// deleteCommentRecursive 删除评论及其所有子评论
func deleteCommentRecursive(tx *gorm.DB, commentID int) error {
    // 删除所有子评论
    var childComments []models.Comment
    if err := tx.Where("parent_id = ?", commentID).Find(&childComments).Error; err != nil {
        return err
    }

    for _, child := range childComments {
        if err := deleteCommentRecursive(tx, int(child.ID)); err != nil {
            return err
        }
    }

    // 删除当前评论
    if err := tx.Delete(&models.Comment{}, commentID).Error; err != nil {
        return err
    }

    return nil
}

// LikeNews 处理用户点赞新闻的请求
func (nc *NewsController) LikeNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 查找新闻
    var news models.News
    if err := tx.First(&news, newsID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find news"})
        return
    }

    // 检查用户是否已点赞
    var user models.User
    if err := tx.Preload("LikedNews").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 如果用户已点赞，返回错误
    for _, likedNews := range user.LikedNews {
        if likedNews.ID == uint(newsID) {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "You have already liked this news"})
            return
        }
    }

    // 建立点赞关系
    if err := tx.Model(&user).Association("LikedNews").Append(&news); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like news"})
        return
    }

    // 增加点赞数
    if err := tx.Model(&news).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":     "News liked successfully",
        "like_count":  news.LikeCount + 1,
    })
}

// CancelLikeNews 处理用户取消点赞新闻的请求
func (nc *NewsController) CancelLikeNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 查找新闻
    var news models.News
    if err := tx.First(&news, newsID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find news"})
        return
    }

    // 查找用户
    var user models.User
    if err := tx.Preload("LikedNews").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 检查用户是否已点赞
    liked := false
    for _, likedNews := range user.LikedNews {
        if likedNews.ID == uint(newsID) {
            liked = true
            break
        }
    }
    if !liked {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not liked this news"})
        return
    }

    // 删除点赞关系
    if err := tx.Model(&user).Association("LikedNews").Delete(&news); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel like"})
        return
    }

    // 减少点赞数
    if err := tx.Model(&news).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":     "News like canceled successfully",
        "like_count":  news.LikeCount - 1,
    })
}

// FavoriteNews 处理用户收藏新闻的请求
func (nc *NewsController) FavoriteNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 查找新闻
    var news models.News
    if err := tx.First(&news, newsID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find news"})
        return
    }

    // 查找用户
    var user models.User
    if err := tx.Preload("FavoritedNews").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 如果用户已收藏，返回错误
    for _, favoritedNews := range user.FavoritedNews {
        if favoritedNews.ID == uint(newsID) {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "You have already favorited this news"})
            return
        }
    }

    // 建立收藏关系
    if err := tx.Model(&user).Association("FavoritedNews").Append(&news); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to favorite news"})
        return
    }

    // 增加收藏数
    if err := tx.Model(&news).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorite count"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":         "News favorited successfully",
        "favorite_count":  news.FavoriteCount + 1,
    })
}

// CancelFavoriteNews 处理用户取消收藏新闻的请求
func (nc *NewsController) CancelFavoriteNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 查找新闻
    var news models.News
    if err := tx.First(&news, newsID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find news"})
        return
    }

    // 查找用户
    var user models.User
    if err := tx.Preload("FavoritedNews").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 检查用户是否已收藏
    favorited := false
    for _, favoritedNews := range user.FavoritedNews {
        if favoritedNews.ID == uint(newsID) {
            favorited = true
            break
        }
    }
    if !favorited {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not favorited this news"})
        return
    }

    // 删除收藏关系
    if err := tx.Model(&user).Association("FavoritedNews").Delete(&news); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel favorite"})
        return
    }

    // 减少收藏数
    if err := tx.Model(&news).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorite count"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":         "News favorite canceled successfully",
        "favorite_count":  news.FavoriteCount - 1,
    })
}

// DislikeNews 处理用户点踩新闻的请求
func (nc *NewsController) DislikeNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 查找新闻
    var news models.News
    if err := tx.First(&news, newsID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find news"})
        return
    }

    // 查找用户
    var user models.User
    if err := tx.Preload("DislikedNews").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 如果用户已点踩，返回错误
    for _, dislikedNews := range user.DislikedNews {
        if dislikedNews.ID == uint(newsID) {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"error": "You have already disliked this news"})
            return
        }
    }

    // 建立点踩关系
    if err := tx.Model(&user).Association("DislikedNews").Append(&news); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dislike news"})
        return
    }

    // 增加点踩数
    if err := tx.Model(&news).UpdateColumn("dislike_count", gorm.Expr("dislike_count + ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dislike count"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":         "News disliked successfully",
        "dislike_count":   news.DislikeCount + 1,
    })
}

// CancelDislikeNews 处理用户取消点踩新闻的请求
func (nc *NewsController) CancelDislikeNews(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 查找新闻
    var news models.News
    if err := tx.First(&news, newsID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find news"})
        return
    }

    // 查找用户
    var user models.User
    if err := tx.Preload("DislikedNews").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 检查用户是否已点踩
    disliked := false
    for _, dislikedNews := range user.DislikedNews {
        if dislikedNews.ID == uint(newsID) {
            disliked = true
            break
        }
    }
    if !disliked {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{"error": "You have not disliked this news"})
        return
    }

    // 删除点踩关系
    if err := tx.Model(&user).Association("DislikedNews").Delete(&news); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel dislike"})
        return
    }

    // 减少点踩数
    if err := tx.Model(&news).UpdateColumn("dislike_count", gorm.Expr("dislike_count - ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dislike count"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":         "News dislike canceled successfully",
        "dislike_count":   news.DislikeCount - 1,
    })
}

// 用户浏览新闻
// ViewNews 处理用户浏览新闻的请求
func (nc *NewsController) ViewNews(c *gin.Context) {
    // 从 JWT 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 开始事务
    tx := nc.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
    }()

    // 查找新闻
    var news models.News
    if err := tx.First(&news, newsID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            tx.Rollback()
            c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
            return
        }
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find news"})
        return
    }

    // 查找用户
    var user models.User
    if err := tx.Preload("ViewedNews").First(&user, userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
        return
    }

    // 检查用户是否已浏览过该新闻
    alreadyViewed := false
    for _, viewedNews := range user.ViewedNews {
        if viewedNews.ID == uint(newsID) {
            alreadyViewed = true
            break
        }
    }

    // 限制浏览记录最大数量
    const maxViewed = 200
    if !alreadyViewed {
        viewedCount := tx.Model(&user).Association("ViewedNews").Count()

        if viewedCount >= maxViewed {
            // 删除最早的浏览记录
            var oldestNews models.News
            if err := tx.Raw(`
                SELECT news.* FROM news
                JOIN user_viewed_news ON news.id = user_viewed_news.news_id
                WHERE user_viewed_news.user_id = ?
                ORDER BY user_viewed_news.created_at ASC
                LIMIT 1
            `, userID).Scan(&oldestNews).Error; err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch oldest viewed news"})
                return
            }

            if oldestNews.ID != 0 {
                if err := tx.Model(&user).Association("ViewedNews").Delete(&oldestNews); err != nil {
                    tx.Rollback()
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete oldest viewed news"})
                    return
                }
            }
        }

        // 添加新的浏览记录
        if err := tx.Model(&user).Association("ViewedNews").Append(&news); err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record news view"})
            return
        }
    }

    // 更新新闻的浏览计数
    if err := tx.Model(&news).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update view count"})
        return
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "News view recorded successfully",
    })
}

func (nc *NewsController) GetUserNewsStatus(c *gin.Context) {
    // 从 JWT 中获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 获取新闻 ID
    newsID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
        return
    }

    // 查询新闻是否存在
    var newsExists bool
    if err := nc.DB.Model(&models.News{}).Select("count(*) > 0").Where("id = ?", newsID).Scan(&newsExists).Error; err != nil || !newsExists {
        c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
        return
    }

    // 查询用户对新闻的操作状态
    var user models.User
    if err := nc.DB.Preload("LikedNews").Preload("FavoritedNews").Preload("DislikedNews").First(&user, userID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user data"})
        return
    }

    // 检查是否点赞、收藏、点踩
    liked := false
    favorited := false
    disliked := false

    for _, likedNews := range user.LikedNews {
        if likedNews.ID == uint(newsID) {
            liked = true
            break
        }
    }

    for _, favoritedNews := range user.FavoritedNews {
        if favoritedNews.ID == uint(newsID) {
            favorited = true
            break
        }
    }

    for _, dislikedNews := range user.DislikedNews {
        if dislikedNews.ID == uint(newsID) {
            disliked = true
            break
        }
    }

    // 返回状态
    c.JSON(http.StatusOK, gin.H{
        "liked":     liked,
        "favorited": favorited,
        "disliked":  disliked,
    })
}

func (nc *NewsController) SearchNews(c *gin.Context) {
    // 定义请求体结构
    var requestBody struct {
        Query string `json:"query"`
    }

    // 解析请求体
    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    query := requestBody.Query
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Query string cannot be empty"})
        return
    }

    // 拆分字符串按空格分词
    keywords := strings.Fields(query)
    if len(keywords) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Query string is invalid"})
        return
    }

    // 构建动态查询条件
    db := nc.DB.Model(&models.News{})
    for _, keyword := range keywords {
        db = db.Where("title LIKE ?", "%"+keyword+"%")
    }

    // 执行查询
    var newsList []models.News
    if err := db.Order("id DESC").Find(&newsList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search news"})
        return
    }

    // 构建返回结构，仅返回 ID 和标题（根据需要调整）
    results := make([]gin.H, len(newsList))
    for i, news := range newsList {
        results[i] = gin.H{
            "id":    news.ID,
        }
    }

    // 返回搜索结果
    c.JSON(http.StatusOK, gin.H{
        "message": "Search results retrieved successfully",
        "results": results,
    })
}