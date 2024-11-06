// models/resource.go
package models

type ResourceType string

const (
    ResourceImage ResourceType = "image"
    ResourceVideo ResourceType = "video"
)

type Resource struct {
    ID     uint         `gorm:"primaryKey" json:"id"`
    NewsID uint         `json:"news_id"`
    Type   ResourceType `gorm:"size:10;not null" json:"type"`
    URL    string       `gorm:"size:255;not null" json:"url"`
}