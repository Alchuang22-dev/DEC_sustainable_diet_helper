package models

import (
    "time"
)

type Family struct {
    ID        uint      `gorm:"primaryKey" json:"id"`                // 家庭唯一标识
    Name      string    `gorm:"size:100;not null" json:"name"`       // 家庭名称
    Token     string    `gorm:"size:64;unique;not null" json:"token"` // 家庭唯一 Token，用于用户加入
    MemberCount uint       `json:"member_count"`                        // 成员数

    Admins       []User            `gorm:"many2many:family_admins;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"admins"` // 管理员列表
    Members      []User            `gorm:"many2many:family_members;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"members"` // 普通成员列表
    WaitingList  []User            `gorm:"many2many:family_waiting_list;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"waiting_list"` // 待确认成员列表

    Dishes       []FamilyDish  `gorm:"foreignKey:FamilyID" json:"dishes"` // 想吃的菜列表

    CreatedAt time.Time `json:"created_at"`                         // 创建时间
    UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
}

type FamilyDish struct {
    ID             uint   `gorm:"primaryKey" json:"id"`                              // 主键
    FamilyID       uint      `gorm:"not null;index:idx_family_dish,unique" json:"family_id"`        // 家庭 ID
    DishID         uint      `gorm:"not null;index:idx_family_dish,unique" json:"dish_id"`          // 菜品 ID
    ProposerUserID uint      `gorm:"not null;index:idx_family_dish,unique" json:"proposer_user_id"` // 提议用户 ID
    LevelOfDesire  uint   `gorm:"not null" json:"level_of_desire"`                    // 想吃程度，1、2、3

    // 关联关系
    Family   Family `gorm:"foreignKey:FamilyID;constraint:OnDelete:CASCADE;" json:"family"`
    Proposer User   `gorm:"foreignKey:ProposerUserID;constraint:OnDelete:SET NULL;" json:"proposer"`
}

func (FamilyDish) TableName() string {
    return "family_dishes"
}