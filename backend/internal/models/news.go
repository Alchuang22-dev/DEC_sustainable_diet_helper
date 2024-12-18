// models/news.go
package models

import (
    "time"
)

// 公共新闻结构
type News struct {
    ID               uint           `gorm:"primaryKey" json:"id"`
    Title            string         `gorm:"size:255;not null" json:"title"`
    UploadTime       time.Time      `json:"upload_time"`
    ViewCount        int            `json:"view_count"`
    LikeCount        int            `json:"like_count"`
    FavoriteCount    int            `json:"favorite_count"`
    DislikeCount     int            `json:"dislike_count"`
    ShareCount       int            `json:"share_count"`
    Comments         []Comment      `gorm:"foreignKey:NewsID" json:"comments"`

    // 作者信息
    AuthorID         uint           `json:"author_id"`
    Author           User           `gorm:"foreignKey:AuthorID" json:"author"`

    // 关联的用户列表
    LikedByUsers     []User         `gorm:"many2many:user_likes_news;"`
    FavoritedByUsers []User         `gorm:"many2many:user_favorites_news;"`
    DislikedByUsers  []User         `gorm:"many2many:user_dislikes_news;"`

    // 内容
    Paragraphs       []Paragraph    `gorm:"foreignKey:NewsID" json:"paragraphs"`
    Images           []NewsImage    `gorm:"foreignKey:NewsID" json:"images"`
}

// 段落模型
type Paragraph struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    NewsID uint   `json:"news_id"`
    Text   string `gorm:"type:text" json:"text"`
}

// NewsImage 结构体，表示新闻中的图片
type NewsImage struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    NewsID      uint   `json:"news_id"`
    URL         string `gorm:"size:255;not null" json:"url"`
    Description string `gorm:"type:text" json:"description"` // 新增描述字段
}

// Draft 结构体，表示新闻草稿
type Draft struct {
    ID           uint             `gorm:"primaryKey" json:"id"`
    Title        string           `gorm:"size:255;not null" json:"title"`
    AuthorID     uint             `json:"author_id"`
    Author       User             `gorm:"foreignKey:AuthorID" json:"author"`
    Paragraphs   []DraftParagraph `gorm:"foreignKey:DraftID" json:"paragraphs"`
    Images       []DraftImage     `gorm:"foreignKey:DraftID" json:"images"`
    CreatedAt    time.Time        `json:"created_at"`
    UpdatedAt    time.Time        `json:"updated_at"`
}

// DraftParagraph 结构体，表示草稿中的段落
type DraftParagraph struct {
    ID      uint   `gorm:"primaryKey" json:"id"`
    DraftID uint   `json:"draft_id"`
    Text    string `gorm:"type:text" json:"text"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// DraftImage 结构体，表示草稿中的图片
type DraftImage struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    DraftID     uint   `json:"draft_id"`
    URL         string `gorm:"size:255;not null" json:"url"`
    Description string `gorm:"type:text" json:"description"` // 新增描述字段
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}