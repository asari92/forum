package service

import (
	"errors"

	"forum/entities"
	"forum/repository"
)

type UserUseCase struct {
	UserRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
	}
}

func (u *UserUseCase) Insert(username, email, password string) error {
	// Валидация данных
	if username == "" || email == "" || password == "" {
		return errors.New("fields can't be empty")
	}
	return u.UserRepo.Insert(username, email, password)
}

func (u *UserUseCase) Authenticate(email, password string) (int, error) {
	// Валидация данных
	if email == "" || password == "" {
		return 0, errors.New("fields can't be empty")
	}
	return u.UserRepo.Authenticate(email, password)
}

func (u *UserUseCase) UserExists(id int) (bool, error) {
	return u.UserRepo.Exists(id)
}

func (u *UserUseCase) GetUserByID(id int) (*entities.User, error) {
	return u.UserRepo.Get(id)
}

func (u *UserUseCase) UpdatePassword(id int, currentPassword, newPassword string) error {
	return u.UserRepo.UpdatePassword(id, currentPassword, newPassword)
}
