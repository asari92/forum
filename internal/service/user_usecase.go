package service

import (
	"errors"

	"forum/internal/entities"
	"forum/internal/repository"
	"forum/internal/validator"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

type accountPasswordUpdateForm struct {
	CurrentPassword         string
	NewPassword             string
	NewPasswordConfirmation string
	validator.Validator
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) NewAccountPasswordUpdateForm() accountPasswordUpdateForm {
	return accountPasswordUpdateForm{}
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

func (u *UserUseCase) UpdatePassword(userID int, form *accountPasswordUpdateForm) error {
	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		return entities.ErrInvalidCredentials
	}

	return u.userRepo.UpdatePassword(userID, form.CurrentPassword, form.NewPassword)
}
