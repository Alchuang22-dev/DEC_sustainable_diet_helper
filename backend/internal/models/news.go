// models/news.go
package models

import (
    "time"
)

type NewsType string

const (
    NewsTypeVideo NewsType = "video"   // 视频新闻
    NewsTypeRegular NewsType = "regular" // 常规新闻
    NewsTypeExternal NewsType = "external" // 外部新闻，仅链接
)

// 公共新闻结构
type News struct {
    ID              uint         `gorm:"primaryKey" json:"id"`
    Title           string       `gorm:"size:255;not null" json:"title"`
    Description     string       `gorm:"type:text" json:"description"`
    UploadTime      time.Time    `json:"upload_time"`
    ViewCount       int          `json:"view_count"`
    LikeCount       int          `json:"like_count"`
    FavoriteCount   int          `json:"favorite_count"`
    DislikeCount    int          `json:"dislike_count"`
    ShareCount      int          `json:"share_count"`
    Comments        []Comment    `gorm:"foreignKey:NewsID" json:"comments"`
    NewsType        NewsType     `gorm:"size:10;not null" json:"news_type"`     // 视频新闻/常规新闻/外部新闻

    // 作者信息
    AuthorID        uint         `json:"author_id"`
    Author          User         `gorm:"foreignKey:AuthorID" json:"author"`
    AuthorName      string       `gorm:"size:255" json:"authorName"`               // 对应前端的 authorName
    AuthorAvatar    string       `gorm:"size:255" json:"authorAvatar"`             // 对应前端的 authorAvatar

    // 关联的用户列表
    LikedByUsers    []User       `gorm:"many2many:user_likes_news;"`
    FavoritedByUsers []User      `gorm:"many2many:user_favorites_news;"`
    DislikedByUsers []User       `gorm:"many2many:user_dislikes_news;"`

    // 类型特定字段
    Video           Video        `gorm:"foreignKey:NewsID" json:"video"`
    Paragraphs      []Paragraph  `gorm:"foreignKey:NewsID" json:"paragraphs"`
    Resources       []Resource   `gorm:"foreignKey:NewsID" json:"resources"`
    ExternalLink    string       `gorm:"size:255" json:"external_link"`         // 外部新闻链接

    Tags            []Tag        `gorm:"many2many:news_tags;" json:"tags"`      // 标签
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

// 标签库
type Tag struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `gorm:"size:100;unique;not null" json:"name"` // 标签名称
    News  []News `gorm:"many2many:news_tags;" json:"news"`     // 拥有该标签的新闻
}