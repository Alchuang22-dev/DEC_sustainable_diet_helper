package models

import (
    "time"
)

type Family struct {
    ID        uint      `gorm:"primaryKey" json:"id"`                // 家庭唯一标识
    Name      string    `gorm:"size:100;not null" json:"name"`       // 家庭名称
    Token     string    `gorm:"size:64;unique;not null" json:"token"` // 家庭唯一 Token，用于用户加入
    Admins    []User    `gorm:"many2many:family_admins;" json:"admins"` // 管理员列表
    Members   []User    `gorm:"many2many:family_members;" json:"members"` // 普通成员列表
    CreatedAt time.Time `json:"created_at"`                         // 创建时间
    UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
}