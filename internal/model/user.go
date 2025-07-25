package model

import "time"

// User 对应数据库users表
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex" json:"username"`
	Birthdate time.Time `json:"birthdate"`
	Phone     string    `gorm:"size:20;uniqueIndex" json:"phone" binding:"required"`
	Email     string    `gorm:"size:100" json:"email"`
	Password  string    `gorm:"size:100" json:"password"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	UpdatedAt time.Time
}
