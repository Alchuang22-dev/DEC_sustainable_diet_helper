package models

import (
    "time"
)

type User struct {
    ID          uint      `gorm:"primaryKey" json:"id"`           // 用户唯一标识
    Nickname    string    `gorm:"size:100;not null" json:"nickname"` // 用户昵称
    OpenID      string    `gorm:"size:64;unique;not null" json:"open_id"` // 微信 OpenID，用于标识微信用户
    SessionKey  string    `gorm:"size:64" json:"-"`              // 微信会话密钥（敏感信息，不返回给前端）
    AvatarURL   string    `gorm:"size:255" json:"avatar_url"`    // 用户头像 URL
    CreatedAt   time.Time `json:"created_at"`                   // 用户创建时间
    UpdatedAt   time.Time `json:"updated_at"`                   // 用户更新时间

    FamilyID    uint      `json:"family_id"`                    // 所属家庭 ID，唯一
    Family      *Family   `gorm:"foreignKey:FamilyID" json:"family"` // 与家庭的外键关系
    PendingFamilyID *uint     `json:"pending_family_id"`            // 正在等待批准的家庭 ID
    PendingFamily  *Family    `gorm:"foreignKey:PendingFamilyID" json:"pending_family"` // 等待批准的家庭

    // 用户与新闻的多对多关系
    LikedNews     []News `gorm:"many2many:user_likes_news;" json:"liked_news"`       // 用户点赞的新闻
    FavoritedNews []News `gorm:"many2many:user_favorites_news;" json:"favorited_news"` // 用户收藏的新闻
    DislikedNews  []News `gorm:"many2many:user_dislikes_news;" json:"disliked_news"`   // 用户点踩的新闻
    ViewedNews    []News `gorm:"many2many:user_viewed_news;" json:"viewed_news"`       // 用户看过的新闻
}
