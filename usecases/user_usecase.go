package usecases

import (
	"errors"

	"forum/entities"
)

type UserRepository interface {
	Insert(username, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*entities.User, error)
	UpdatePassword(id int, currentPassword, newPassword string) error
}

type UserUseCase struct {
	UserRepo UserRepository
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
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
