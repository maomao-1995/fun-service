package service

import (
	"fun-service/internal/model"
	"fun-service/internal/repository"
)

// UserService 业务逻辑层
type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: &repository.UserRepository{},
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username string, age int) (*model.User, error) {
	user := &model.User{Username: username, Age: age}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID 查询用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}
