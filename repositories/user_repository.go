package repositories

import "forum/entities"

type UserModelInterface interface {
	Insert(username, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*entities.User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}
