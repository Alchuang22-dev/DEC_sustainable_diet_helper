// models/news.go
package models

import (
    "time"
)

// 公共新闻结构
type BaseNews struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    Title         string         `gorm:"size:255;not null" json:"title"`
    Description   string         `gorm:"type:text" json:"description"`
    UploadTime    time.Time      `json:"upload_time"`
    LikeCount     int            `json:"like_count"`
    FavoriteCount int            `json:"favorite_count"`
    DislikeCount  int            `json:"dislike_count"`
    ViewCount     int            `json:"view_count"`
    Comments      []Comment      `gorm:"foreignKey:NewsID" json:"comments"`
}

// 视频新闻
type VideoNews struct {
    BaseNews
    VideoURL string `gorm:"size:255;not null" json:"video_url"`
}

// 常规新闻
type RegularNews struct {
    BaseNews
    Paragraphs []Paragraph `gorm:"foreignKey:NewsID" json:"paragraphs"`
    Resources  []Resource  `gorm:"foreignKey:NewsID" json:"resources"`
}

// 段落模型
type Paragraph struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    NewsID uint   `json:"news_id"`
    Text   string `gorm:"type:text" json:"text"`
}