// models/news.go
package models

import (
    "time"
)

// 新闻类型定义
type NewsType string

const (
    NewsTypeVideo   NewsType = "video"
    NewsTypeRegular NewsType = "regular"
)

// 新闻结构体
type News struct {
    ID              uint         `gorm:"primaryKey" json:"id"`
    Title           string       `gorm:"size:255;not null" json:"newsName"`        // 对应前端的 newsName
    Description     string       `gorm:"type:text" json:"newsinfo"`                // 对应前端的 newsinfo
    UploadTime      time.Time    `json:"time"`                                     // 对应前端的 time
    ViewCount       int          `json:"view_count"`
    ShareCount      int          `json:"shareCount"`                               // 对应前端的 shareCount
    LikeCount       int          `json:"likeCount"`                                // 对应前端的 likeCount
    FavoriteCount   int          `json:"favoriteCount"`                            // 对应前端的 favoriteCount
    DislikeCount    int          `json:"dislikeCount"`                             // 对应前端的 dislikeCount
    NewsType        NewsType     `gorm:"size:10;not null" json:"news_type"`

    // 作者信息
    AuthorID        uint         `json:"author_id"`
    Author          User         `gorm:"foreignKey:AuthorID" json:"author"`
    AuthorName      string       `gorm:"size:255" json:"authorName"`               // 对应前端的 authorName
    AuthorAvatar    string       `gorm:"size:255" json:"authorAvatar"`             // 对应前端的 authorAvatar

    // 标签和分类
    Tags            []Tag        `gorm:"many2many:news_tags;" json:"tabs"`         // 对应前端的 tabs

    // 内容相关
    NewsSrc         string       `gorm:"size:255" json:"newsSrc"`                  // 对应前端的 newsSrc
    ImgsSrc         []Image      `gorm:"foreignKey:NewsID" json:"imgsSrc"`         // 对应前端的 imgsSrc
    NewsBody        string       `gorm:"type:text" json:"newsbody"`                // 对应前端的 newsbody
    Paragraphs      []Paragraph  `gorm:"foreignKey:NewsID" json:"paragraphs"`
    Video           Video        `gorm:"foreignKey:NewsID" json:"video"`
    Resources       []Resource   `gorm:"foreignKey:NewsID" json:"resources"`

    // 关联的用户列表
    LikedByUsers     []User      `gorm:"many2many:user_likes_news;" json:"-"`
    FavoritedByUsers []User      `gorm:"many2many:user_favorites_news;" json:"-"`
    DislikedByUsers  []User      `gorm:"many2many:user_dislikes_news;" json:"-"`

    Comments        []Comment    `gorm:"foreignKey:NewsID" json:"comments"`
}

// 视频结构体
type Video struct {
    ID       uint   `gorm:"primaryKey"`
    NewsID   uint   `gorm:"uniqueIndex" json:"news_id"` // 关联到 News 表
    VideoURL string `gorm:"size:255;not null" json:"video_url"`
}

// 段落模型
type Paragraph struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    NewsID uint   `json:"news_id"`
    Text   string `gorm:"type:text" json:"text"`
}

// 标签模型
type Tag struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `gorm:"size:255;not null" json:"name"`
}

// 图片模型
type Image struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    NewsID uint   `json:"news_id"`
    URL    string `gorm:"size:255;not null" json:"url"`
}

// 用户模型（简化）
type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Name      string `gorm:"size:255" json:"name"`
    AvatarURL string `gorm:"size:255" json:"avatar_url"`
}

// 评论模型（简化）
type Comment struct {
    ID      uint   `gorm:"primaryKey" json:"id"`
    NewsID  uint   `json:"news_id"`
    UserID  uint   `json:"user_id"`
    Content string `gorm:"type:text" json:"content"`
}

// 资源模型（保持不变）
type Resource struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    NewsID uint   `json:"news_id"`
    URL    string `gorm:"size:255;not null" json:"url"`
    Type   string `gorm:"size:50" json:"type"`
}