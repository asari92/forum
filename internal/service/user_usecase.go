package service

import (
	"errors"
	"forum/internal/entities"
	"forum/internal/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) Insert(username, email, password string) error {
	// Валидация данных
	if username == "" || email == "" || password == "" {
		return errors.New("fields can't be empty")
	}
	return u.userRepo.Insert(username, email, password)
}

func (u *UserUseCase) Authenticate(email, password string) (int, error) {
	// Валидация данных
	if email == "" || password == "" {
		return 0, errors.New("fields can't be empty")
	}
	return u.userRepo.Authenticate(email, password)
}

func (u *UserUseCase) UserExists(id int) (bool, error) {
	return u.userRepo.Exists(id)
}

func (u *UserUseCase) GetUserByID(id int) (*entities.User, error) {
	return u.userRepo.Get(id)
}

func (u *UserUseCase) UpdatePassword(id int, currentPassword, newPassword string) error {
	return u.userRepo.UpdatePassword(id, currentPassword, newPassword)
}
