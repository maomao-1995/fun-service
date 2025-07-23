package model

import "gorm.io/gorm"

// User 对应数据库users表
type User struct {
	gorm.Model
	Username string `gorm:"size:50;uniqueIndex" json:"username"`
	Age      int    `json:"age"`
}
