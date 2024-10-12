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
	GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*Post, error)
	GetUserPaginatedPosts(userId, page, pageSize int) ([]*Post, error)
	GetUserLikedPaginatedPosts(userId, page, pageSize int) ([]*Post, error)
	GetAllPaginatedPosts(page, pageSize int) ([]*Post, error)
}

type Post struct {
	ID       int
	Title    string
	Content  string
	UserID   int
	UserName string
	Created  string
}

type PostModel struct {
	DB *sql.DB
}

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

func (m *PostModel) Get(postID int) (*Post, error) {

	stmt := `SELECT id,title, content, created 
	FROM posts 
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, postID)

	p := &Post{}
	var created string

	err := row.Scan(&p.ID, &p.Title, &p.Content, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	postTime, err := time.Parse("2006-01-02 15:04:05", created)
	if err != nil {
		return nil, err
	}
	p.Created = postTime.Format(time.RFC3339)

	err = m.DB.QueryRow(`SELECT username FROM posts 
	INNER JOIN users ON users.id = posts.user_id
	WHERE posts.id = ?`, postID).Scan(&p.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			p.UserName = "Deleted User"
			return p, nil
		} else {
			return nil, err
		}
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
		postTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}
		p.Created = postTime.Format(time.RFC3339)

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// Получение постов по категориям с пагинацией
func (m *PostModel) GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*Post, error) {
	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs))
	offset := (page - 1) * pageSize

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
        ORDER BY p.created DESC
		LIMIT ? OFFSET ?`, strings.Join(placeholders, ", "))

	// Добавляем в аргументы запроса: количество категорий; запрашиваем на одну запись больше, чем pageSize, чтобы проверить наличие следующей страницы; смещение
	args = append(args, len(categoryIDs), pageSize+1, offset)

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

		postTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}
		p.Created = postTime.Format(time.RFC3339)

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) GetUserPaginatedPosts(userId, page, pageSize int) ([]*Post, error) {
	offset := (page - 1) * pageSize

	stmt := `SELECT id, title, content, user_id, created FROM posts
	WHERE user_id = ?
    ORDER BY id DESC
	LIMIT ? OFFSET ?`

	// запрашиваем на одну запись больше, чем pageSize
	rows, err := m.DB.Query(stmt, userId, pageSize+1, offset)
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

		postTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}
		p.Created = postTime.Format(time.RFC3339)
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) GetUserLikedPaginatedPosts(userId, page, pageSize int) ([]*Post, error) {
	offset := (page - 1) * pageSize

	stmt := `SELECT id, title, content, p.user_id, created 
	FROM posts p
	INNER JOIN post_reactions pr ON p.id = pr.post_id
	WHERE pr.user_id = ? AND pr.is_like = true
    ORDER BY id DESC
	LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(stmt, userId, pageSize+1, offset)
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

		postTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}
		p.Created = postTime.Format(time.RFC3339)

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostModel) GetAllPaginatedPosts(page, pageSize int) ([]*Post, error) {
	offset := (page - 1) * pageSize // Вычисляем смещение для текущей страницы

	stmt := `SELECT id, title, content, user_id, created FROM posts
             ORDER BY created DESC
             LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(stmt, pageSize+1, offset) // Лимит на одну запись больше
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

		postTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}
		p.Created = postTime.Format(time.RFC3339)

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
