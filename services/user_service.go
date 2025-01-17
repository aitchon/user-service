package services

import (
	"user-service/models"
	"user-service/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.Repo.CreateUser(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.Repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.Repo.DeleteUser(id)
}
