package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"forum/entities"
)

type PostRepository struct {
	DB *sql.DB
}

func (r *PostRepository) InsertPostWithCategories(title, content string, userID int, categoryIDs []int) (int, error) {
	// Начинаем транзакцию
	tx, err := r.DB.Begin()
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

func (r *PostRepository) GetPost(postID int) (*entities.Post, error) {
	stmt := `SELECT posts.id,title,username,content, posts.created 
	FROM posts LEFT JOIN users ON posts.user_id = users.id
    WHERE posts.id = ?`

	row := r.DB.QueryRow(stmt, postID)

	p := &entities.Post{}
	var created string
	var username sql.NullString

	err := row.Scan(&p.ID, &p.Title, &username, &p.Content, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	if username.Valid {
		p.UserName = username.String
	} else {
		p.UserName = "Deleted User"
	}

	postTime, err := time.Parse("2006-01-02 15:04:05", created)
	if err != nil {
		return nil, err
	}
	p.Created = postTime.Format(time.RFC3339)

	return p, nil
}

// Получение постов по категориям с пагинацией
func (r *PostRepository) GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error) {
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

	rows, err := r.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*entities.Post{}
	var created string

	// Обрабатываем строки результата запроса.
	for rows.Next() {
		p := &entities.Post{}
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

func (r *PostRepository) GetUserPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error) {
	offset := (page - 1) * pageSize

	stmt := `SELECT id, title, content, user_id, created FROM posts
	WHERE user_id = ?
    ORDER BY id DESC
	LIMIT ? OFFSET ?`

	// запрашиваем на одну запись больше, чем pageSize
	rows, err := r.DB.Query(stmt, userId, pageSize+1, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*entities.Post{}
	var created string

	for rows.Next() {
		p := &entities.Post{}
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

func (r *PostRepository) GetUserLikedPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error) {
	offset := (page - 1) * pageSize

	stmt := `SELECT id, title, content, p.user_id, created 
	FROM posts p
	INNER JOIN post_reactions pr ON p.id = pr.post_id
	WHERE pr.user_id = ? AND pr.is_like = true
    ORDER BY id DESC
	LIMIT ? OFFSET ?`

	rows, err := r.DB.Query(stmt, userId, pageSize+1, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*entities.Post{}
	var created string

	for rows.Next() {
		p := &entities.Post{}
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

func (r *PostRepository) GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error) {
	offset := (page - 1) * pageSize // Вычисляем смещение для текущей страницы

	stmt := `SELECT id, title, content, user_id, created FROM posts
             ORDER BY created DESC
             LIMIT ? OFFSET ?`

	rows, err := r.DB.Query(stmt, pageSize+1, offset) // Лимит на одну запись больше
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*entities.Post{}
	var created string

	for rows.Next() {
		p := &entities.Post{}
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

// Удаление поста
func (r *PostRepository) DeletePost(postID int) error {
	stmt := `DELETE FROM posts WHERE id = ?`
	_, err := r.DB.Exec(stmt, postID)
	return err
}