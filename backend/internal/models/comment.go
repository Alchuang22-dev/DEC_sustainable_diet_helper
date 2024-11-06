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
}