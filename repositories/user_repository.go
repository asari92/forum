package repositories

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"forum/entities"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (m *UserRepository) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, password, created)
    VALUES(?, ?, ?, datetime('now'))`

	_, err = m.DB.Exec(stmt, username, email, string(hashedPassword))
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == sqlite3.ErrConstraint && strings.Contains(sqliteError.Error(), "users.email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserRepository) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserRepository) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserRepository) Get(id int) (*entities.User, error) {
	stmt := `SELECT id, username, email, created FROM users WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	u := &entities.User{}
	var created string

	err := row.Scan(&u.ID, &u.Username, &u.Email, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	userTime, err := time.Parse("2006-01-02 15:04:05", created)
	if err != nil {
		return nil, err
	}
	u.Created = userTime.Format(time.RFC3339)

	return u, nil
}

func (m *UserRepository) UpdatePassword(id int, currentPassword, newPassword string) error {
	var currentHashedPassword []byte

	stmt := "SELECT password FROM users WHERE id = ?"

	err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	stmt = "UPDATE users SET password = ? WHERE id = ?"

	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)
	return err
}
