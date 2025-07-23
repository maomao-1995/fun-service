package repository

import (
	"fun-service/internal/model"
	"fun-service/pkg/database"
)

// UserRepository 数据访问层
type UserRepository struct{}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

// GetByID 根据ID查询用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	return &user, err
}
