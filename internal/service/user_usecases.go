package service

import (
	"forum/internal/entities"
	"forum/internal/repository"
	"forum/pkg/validator"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

type userAuthForm struct {
	Username string
	Email    string
	Password string
	validator.Validator
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

func (uc *UserUseCase) NewUserAuthForm() userAuthForm {
	return userAuthForm{}
}

func (uc *UserUseCase) NewAccountPasswordUpdateForm() accountPasswordUpdateForm {
	return accountPasswordUpdateForm{}
}

func (u *UserUseCase) Insert(form *userAuthForm) error {
	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.NotHaveAnySpaces(form.Username), "username", "This field cannot have any spaces")
	form.CheckField(validator.MaxChars(form.Username, 100), "username", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotHaveAnySpaces(form.Email), "email", "This field cannot have any spaces")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.NotHaveAnySpaces(form.Password), "password", "This field cannot have any spaces")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		return entities.ErrInvalidCredentials
	}

	return u.userRepo.Insert(form.Username, form.Email, form.Password)
}

func (u *UserUseCase) Authenticate(form *userAuthForm) (int, error) {
	// Валидация данных
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		return 0, entities.ErrInvalidCredentials
	}

	return u.userRepo.Authenticate(form.Email, form.Password)
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
