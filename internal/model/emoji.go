package model

import "time"

type Emoji struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:100;uniqueIndex" json:"name"`
	URL              string    `gorm:"size:255" json:"url"`
	View_count       int       `gorm:"default:0" json:"view_count"`
	Collection_count int       `gorm:"default:0" json:"collection_count"`
	Tags             []string  `gorm:"type:json" json:"tags"` // Store tags in a JSON array
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
