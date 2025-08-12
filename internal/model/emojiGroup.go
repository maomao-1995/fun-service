package model

import "time"

type EmojiGroup struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:100;uniqueIndex" json:"name"`
	EmojiIDs         []uint    `gorm:"type:json" json:"emoji_ids"` // Store emoji IDs in a JSON array
	View_count       int       `gorm:"default:0" json:"view_count"`
	Collection_count int       `gorm:"default:0" json:"collection_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
