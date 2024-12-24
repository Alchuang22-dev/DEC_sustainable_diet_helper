// models/comment.go
package models

import (
    "time"
)

type Comment struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    Content    string    `gorm:"type:text;not null" json:"content"`
    PublishTime time.Time `json:"publish_time"`
    LikeCount  int       `json:"like_count"`
    UserID     uint      `json:"user_id"`
    NewsID     uint      `json:"news_id"`
    Replies    []Comment `gorm:"foreignKey:ParentID" json:"replies"`
    ParentID   *uint     `json:"parent_id"`
    IsReply    bool      `json:"is_reply"`
    Author      User      `gorm:"foreignKey:UserID" json:"author"`  // 绑定 User 类型的 Author 字段

    LikedByUsers []User `gorm:"many2many:user_likes_comments;" json:"-"`

    // 是否已点赞（不存数据库）
    DidLike  bool  `gorm:"-" json:"did_like"`
}