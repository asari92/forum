package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID      int
	Title   string
	Content string
	UserID  int
	Created time.Time
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title string, content string, userID int) (int, error) {
	stmt := `INSERT INTO posts (title, content, user_id, created_date)
    VALUES(?, ?, ?, datetime('now'))`

	result, err := m.DB.Exec(stmt, title, content, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PostModel) Get(id int) (*Post, error) {
	stmt := `SELECT id, title, content, user_id, created FROM posts
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Post{}
	var created string

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.UserID, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	s.Created, err = time.Parse("2006-01-02 15:04:05", created)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *PostModel) Latest() ([]*Post, error) {
	return nil, nil
}
