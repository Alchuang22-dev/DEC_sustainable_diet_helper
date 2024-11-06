// models/news.go
package models

import (
    "time"
)

type NewsType string

const (
    NewsTypeVideo   NewsType = "video"
    NewsTypeRegular NewsType = "regular"
)

// 公共新闻结构
type News struct {
    ID              uint         `gorm:"primaryKey" json:"id"`
    Title           string       `gorm:"size:255;not null" json:"title"`
    Description     string       `gorm:"type:text" json:"description"`
    UploadTime      time.Time    `json:"upload_time"`
    ViewCount       int          `json:"view_count"`
    Comments        []Comment    `gorm:"foreignKey:NewsID" json:"comments"`
    NewsType        NewsType     `gorm:"size:10;not null" json:"news_type"`

    // 关联的用户列表
    LikedByUsers    []User       `gorm:"many2many:user_likes_news;"`
    FavoritedByUsers []User      `gorm:"many2many:user_favorites_news;"`
    DislikedByUsers []User       `gorm:"many2many:user_dislikes_news;"`

    // 类型特定字段
    Video           Video        `gorm:"foreignKey:NewsID" json:"video"`
    Paragraphs      []Paragraph  `gorm:"foreignKey:NewsID" json:"paragraphs"`
    Resources       []Resource   `gorm:"foreignKey:NewsID" json:"resources"`
}

type Video struct {
    ID     uint   `gorm:"primaryKey"`
    NewsID uint   `gorm:"uniqueIndex" json:"news_id"` // 关联到 News 表
    VideoURL string `gorm:"size:255;not null" json:"video_url"`
}

// 段落模型
type Paragraph struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    NewsID uint   `json:"news_id"`
    Text   string `gorm:"type:text" json:"text"`
}