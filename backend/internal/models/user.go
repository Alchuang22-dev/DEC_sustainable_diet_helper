package models

import (
    "time"
)

type User struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Nickname    string    `gorm:"size:100;not null" json:"nickname"`
    Email       string    `gorm:"size:100;unique;not null" json:"email"`
    Password    string    `gorm:"size:255;not null" json:"-"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    LikedNews     []News `gorm:"many2many:user_likes_news;"`      // 用户点赞的新闻
    FavoritedNews []News `gorm:"many2many:user_favorites_news;"` // 用户收藏的新闻
    DislikedNews  []News `gorm:"many2many:user_dislikes_news;"`  // 用户点踩的新闻
    ViewedNews    []News `gorm:"many2many:user_viewed_news;"`    // 用户看过的新闻
}

type EmailVerification struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Code      string    `gorm:"size:6;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}