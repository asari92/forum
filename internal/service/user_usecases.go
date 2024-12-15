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

func (u *UserUseCase) Insert(username, email, password, role string) (int, error) {
	return u.userRepo.Insert(username, email, password, role)
}

func (u *UserUseCase) Authenticate(email, password string) (*entities.User, error) {
	return u.userRepo.Authenticate(email, password)
}

func (u *UserUseCase) OauthAuthenticate(email string) (*entities.User, error) {
	return u.userRepo.OauthAuthenticate(email)
}

func (u *UserUseCase) UserExists(id int) (bool, error) {
	return u.userRepo.Exists(id)
}

func (u *UserUseCase) GetUserByID(id int) (*entities.User, error) {
	return u.userRepo.Get(id)
}

func (u *UserUseCase) UpdatePassword(userID int, form *accountPasswordUpdateForm) error {
	exists, err := u.userRepo.Exists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	return u.userRepo.UpdatePassword(userID, form.CurrentPassword, form.NewPassword)
}
