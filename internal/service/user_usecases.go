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

type ModerationRequestForm struct {
	Reason string
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

func (uc *UserUseCase) NewModerationForm() ModerationRequestForm {
	return ModerationRequestForm{}
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

func (u *UserUseCase) CreateModerationRequest(userId int, form *ModerationRequestForm) error {
	exists, err := u.userRepo.ExistsModerationRequest(userId)
	if err != nil {
		return err
	}
	if exists {
		form.CheckField(false, "user", "The form has already been submitted")
		return entities.ErrFormAlreadySubmitted
	}

	form.CheckField(validator.NotBlank(form.Reason), "reason", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Reason, validator.TextRX), "reason", "This field must contain only english or russian letters")
	if !form.Valid() {
		return entities.ErrInvalidData
	}

	err = u.userRepo.InsertModerationRequest(userId, form.Reason)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetModerators() ([]*entities.User, error) {
	return u.userRepo.GetModerators()
}

func (u *UserUseCase) DeleteModerator(userId int) error {
	return u.userRepo.DeleteModerator(userId)
}

func (u *UserUseCase) GetModerationApplicants() ([]*entities.ModeratorApplicant, error) {
	return u.userRepo.ListModeratorApplicants()
}

func (u *UserUseCase) DeleteModerationRequest(userId int) error {
	return u.userRepo.DeleteModerationRequest(userId)
}

func (u *UserUseCase) ApproveModerationRequest(userId int) error {
	return u.userRepo.ApproveModeratorRequest(userId)
}
