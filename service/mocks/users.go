package mocks

import (
	"time"

	"forum/entities"
	"forum/repository"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return repository.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, repository.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (*entities.User, error) {
	if id == 1 {
		u := &entities.User{
			ID:       1,
			Username: "Alice",
			Email:    "alice@example.com",
			Created:  time.Now().Format("2006-01-02 15:04:05"),
		}

		return u, nil
	}

	return nil, repository.ErrNoRecord
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "pa$$word" {
			return repository.ErrInvalidCredentials
		}

		return nil
	}

	return repository.ErrNoRecord
}
