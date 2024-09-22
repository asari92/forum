package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PostModelInterface interface {
	InsertPostWithCategories(title, content string, userID int, categoryIDs []int) (int, error)
	Get(id int) (*Post, error)
	GetPostsForCategory(categoryIDs []int) ([]*Post, error)
	GetUserPosts(userId int) ([]*Post, error)
	Latest() ([]*Post, error)
}

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

// func (m *PostModel) Insert(title string, content string, userID int) (int, error) {
// 	stmt := `INSERT INTO posts (title, content, user_id, created)
//     VALUES(?, ?, ?, datetime('now'))`

// 	result, err := m.DB.Exec(stmt, title, content, userID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return int(id), nil
// }

func (m *PostModel) InsertPostWithCategories(title, content string, userID int, categoryIDs []int) (int, error) {
	// Начинаем транзакцию
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	// В случае ошибки откатываем транзакцию
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Вставляем пост в базу данных
	stmt := `INSERT INTO posts (title, content, user_id, created) VALUES (?, ?, ?, datetime('now'))`
	result, err := tx.Exec(stmt, title, content, userID)
	if err != nil {
		return 0, err
	}

	// Получаем ID вставленного поста
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Вставляем категории для поста
	stmt = `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`
	for _, categoryID := range categoryIDs {
		_, err := tx.Exec(stmt, postID, categoryID)
		if err != nil {
			return 0, err
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (m *PostModel) Get(id int) (*Post, error) {
	stmt := `SELECT id, title, content, user_id, created FROM posts
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	p := &Post{}
	var created string

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	p.Created, err = time.Parse("2006-01-02 15:04:05", created)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (m *PostModel) GetPostsForCategory(categoryIDs []int) ([]*Post, error) {
	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs))

	for i, id := range categoryIDs {
		placeholders[i] = "?" // Placeholder для каждого ID категории
		args[i] = id          // Аргумент для SQL-запроса
	}

	stmt := fmt.Sprintf(`
        SELECT p.id, p.title, p.content, p.user_id, p.created 
        FROM posts p
        INNER JOIN post_categories pc ON p.id = pc.post_id
        WHERE pc.category_id IN (%s)
        GROUP BY p.id
        HAVING COUNT(DISTINCT pc.category_id) = ?
        ORDER BY p.created DESC`, strings.Join(placeholders, ", "))

	// Добавляем количество категорий в аргументы запроса
	args = append(args, len(categoryIDs))

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	var created string

	// Обрабатываем строки результата запроса.
	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &created)
		if err != nil {
			return nil, err
		}

		// Преобразуем строку даты в тип time.Time.
		p.Created, err = time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) GetUserPosts(userId int) ([]*Post, error) {
	stmt := `SELECT id, title, content, user_id, created FROM posts
	WHERE user_id = ?
    ORDER BY id DESC`

	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	var created string

	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &created)
		if err != nil {
			return nil, err
		}

		p.Created, err = time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) Latest() ([]*Post, error) {
	stmt := `SELECT id, title, content, user_id, created FROM posts
    ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	var created string

	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &created)
		if err != nil {
			return nil, err
		}

		p.Created, err = time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) Delete(id int) error {
	stmt := `DELETE FROM posts WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
