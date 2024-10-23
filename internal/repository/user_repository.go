package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"forum/internal/entities"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserSqlite3 struct {
	DB *sql.DB
}

func NewUserSqlite3(db *sql.DB) *UserSqlite3 {
	return &UserSqlite3{
		DB: db,
	}
}

func (r *UserSqlite3) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, password, created)
    VALUES(?, ?, ?, datetime('now'))`

	_, err = r.DB.Exec(stmt, username, email, string(hashedPassword))
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

func (r *UserSqlite3) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password FROM users WHERE email = ?"

	err := r.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
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

func (r *UserSqlite3) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := r.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (r *UserSqlite3) Get(id int) (*entities.User, error) {
	stmt := `SELECT id, username, email, created FROM users WHERE id = ?`

	row := r.DB.QueryRow(stmt, id)

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

func (r *UserSqlite3) UpdatePassword(id int, currentPassword, newPassword string) error {
	var currentHashedPassword []byte

	stmt := "SELECT password FROM users WHERE id = ?"

	err := r.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
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

	_, err = r.DB.Exec(stmt, string(newHashedPassword), id)
	return err
}
