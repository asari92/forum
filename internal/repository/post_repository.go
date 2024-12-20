package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"forum/internal/entities"
)

type PostSqlite3 struct {
	DB *sql.DB
}

func NewPostSqlite3(db *sql.DB) *PostSqlite3 {
	return &PostSqlite3{
		DB: db,
	}
}

func (r *PostSqlite3) GetPostOwner(postID int) (int, error) {
	stmt := `SELECT user_id FROM posts
	WHERE id = ?
	`
	var ownerID int
	row := r.DB.QueryRow(stmt, postID)
	err := row.Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, entities.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return ownerID, nil
}

func (r *PostSqlite3) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM posts WHERE id = ?)"
	err := r.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (r *PostSqlite3) InsertPostWithCategories(title, content string, userID int, categoryIDs []int, filePaths []string) (int, error) {
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

	if len(filePaths) != 0 {
		stmt = `INSERT INTO post_images (post_id, image_url) VALUES (?,?)`
		for _, imageUrl := range filePaths {
			_, err := tx.Exec(stmt, postID, imageUrl)
			if err != nil {
				return 0, err
			}
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (r *PostSqlite3) UpdatePostWithImage(title, content string, postID int, filePaths []string) error {
	// Начинаем транзакцию
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// В случае ошибки откатываем транзакцию
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt := `UPDATE posts
	SET title = ?, content = ?
	WHERE id = ?`

	_, err = tx.Exec(stmt, title, content, postID)
	if err != nil {
		return err
	}

	if len(filePaths) != 0 {
		stmt = `INSERT INTO post_images (post_id, image_url) VALUES (?,?)`
		for _, imageUrl := range filePaths {
			_, err := tx.Exec(stmt, postID, imageUrl)
			if err != nil {
				return err
			}
		}
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *PostSqlite3) GetPost(postID int) (*entities.Post, error) {
	stmt := `SELECT posts.id,title,username, posts.user_id, content, posts.created 
	FROM posts LEFT JOIN users ON posts.user_id = users.id
    WHERE posts.id = ?`

	row := r.DB.QueryRow(stmt, postID)

	p := &entities.Post{}
	var created string
	var username sql.NullString

	err := row.Scan(&p.ID, &p.Title, &username, &p.UserID, &p.Content, &created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
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
func (r *PostSqlite3) GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error) {
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

func (r *PostSqlite3) GetUserPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error) {
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

func (r *PostSqlite3) GetUserCommentedPosts(userId, page, pageSize int) ([]*entities.Post, error) {
	offset := (page - 1) * pageSize

	stmt := `SELECT p.id, p.title, p.content, p.user_id, p.created FROM posts as p INNER JOIN comments as c ON p.id = c.post_id
	WHERE c.user_id = ?
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

func (r *PostSqlite3) GetUserLikedPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error) {
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

func (r *PostSqlite3) GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error) {
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
func (r *PostSqlite3) DeletePost(postID int) error {
	stmt := `DELETE FROM posts WHERE id = ?`
	_, err := r.DB.Exec(stmt, postID)
	return err
}

func (r *PostSqlite3) GetImagesByPost(postID int) ([]*entities.Image, error) {
	stmt := `SELECT image_url FROM post_images
	WHERE post_id = ?`
	rows, err := r.DB.Query(stmt, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
		} else {
			return nil, err
		}
	}
	images := []*entities.Image{}
	for rows.Next() {
		image := &entities.Image{}
		err := rows.Scan(&image.UrlImage)
		if err != nil {
			return nil, err
		}
		images = append(images, image)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return images, nil
}





