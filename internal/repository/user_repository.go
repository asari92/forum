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

func (r *UserSqlite3) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := r.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (r *UserSqlite3) Insert(username, email, password, role string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (username, email, password, role, created)
    VALUES(?, ?, ?, ?, datetime('now'))`

	result, err := r.DB.Exec(stmt, username, email, string(hashedPassword), role)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == sqlite3.ErrConstraint {
				if strings.Contains(sqliteError.Error(), "users.email") {
					return 0, entities.ErrDuplicateEmail
				}
				if strings.Contains(sqliteError.Error(), "users.username") {
					return 0, entities.ErrDuplicateUsername
				}
			}
		}
		return 0, err
	}

	// Получаем ID вставленного поста
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

func (r *UserSqlite3) OauthAuthenticate(email string) (*entities.User, error) {
	u := &entities.User{}

	stmt := "SELECT id, role FROM users WHERE email = ?"

	err := r.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (r *UserSqlite3) Authenticate(email, password string) (*entities.User, error) {
	u := &entities.User{}

	stmt := "SELECT id, role, password FROM users WHERE email = ?"

	err := r.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Role, &u.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, entities.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	return u, nil
}

func (r *UserSqlite3) Get(id int) (*entities.User, error) {
	stmt := `SELECT id, username, email, role, created, role FROM users WHERE id = ?`

	row := r.DB.QueryRow(stmt, id)

	u := &entities.User{}
	var created string

	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &created, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
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
			return entities.ErrInvalidCredentials
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

func (r *UserSqlite3) InsertModerationRequest(userId int, reason string) error {
	stmt := `INSERT INTO moderation_requests (user_id, created, reason)
	VALUES (?, datetime('now'), ?)`
	_, err := r.DB.Exec(stmt, userId, reason)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserSqlite3) ExistsModerationRequest(userId int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM moderation_requests WHERE user_id = ?)`
	err := r.DB.QueryRow(stmt, userId).Scan(&exists)
	return exists, err
}

func (r *UserSqlite3) ListModeratorApplicants() ([]*entities.ModeratorApplicant, error) {
	stmt := `SELECT u.id, reason, m.created, username  FROM moderation_requests as m
	INNER JOIN users as u ON m.user_id = u.id
	ORDER BY m.created DESC`
	rows, err := r.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var Applicants []*entities.ModeratorApplicant
	for rows.Next() {
		applicant := &entities.ModeratorApplicant{}
		var created string

		if err := rows.Scan(&applicant.ID, &applicant.Reason, &created, &applicant.Username); err != nil {
			return nil, err
		}

		commentTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}
		

		applicant.Created = commentTime.Format(time.RFC3339)
		Applicants = append(Applicants, applicant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return Applicants, nil
}

func (r *UserSqlite3) GetModerators() ([]*entities.User, error) {
	stmt := `SELECT id, username, email, role FROM users
	WHERE role = ?`
	rows, err := r.DB.Query(stmt, entities.RoleModerator)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}

		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserSqlite3) DeleteModerator(userId int) error {
	stmt := `UPDATE users
	SET role = ?
	WHERE id = ?`
	_, err := r.DB.Exec(stmt, entities.RoleUser, userId)
	return err
}

func (r *UserSqlite3) ApproveModeratorRequest(userId int) error {
	stmt := `UPDATE users
	SET role = ?
	WHERE id = ?`
	_, err := r.DB.Exec(stmt, entities.RoleModerator, userId)
	return err
}

func (r *UserSqlite3) DeleteModerationRequest(userId int) error {
	stmt := `DELETE FROM moderation_requests
	WHERE user_id = ?`
	_, err := r.DB.Exec(stmt, userId)
	return err
}




