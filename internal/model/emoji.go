package model

import (
	"time"

	"gorm.io/datatypes"
)

type Emoji struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	AuthorUUID       string `gorm:"size:36;not null" json:"author_uuid"`
	// AuthorID         uint   `gorm:"not null" json:"author_id"`
	Name             string `gorm:"size:100;" json:"name"`
	URL              string `gorm:"size:255" json:"url"`
	View_count       int    `gorm:"default:0" json:"view_count"`
	Collection_count int    `gorm:"default:0" json:"collection_count"`
	Tags             datatypes.JSON
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
