package repository

import (
	"database/sql"
	"errors"
	"time"

	"forum/internal/entities"
)

type CommentSqlite3 struct {
	DB *sql.DB
}

func NewCommentSqlite3(db *sql.DB) *CommentSqlite3 {
	return &CommentSqlite3{
		DB: db,
	}
}

func (r *CommentSqlite3) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM comments WHERE id = ?)"
	err := r.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (c *CommentSqlite3) InsertComment(postID, userID int, content string) (int, error) {
	stmt := `INSERT INTO comments (post_id, user_id,content, created)
	VALUES (?,?,?, datetime('now'))`
	res, err := c.DB.Exec(stmt, postID, userID, content)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *CommentSqlite3) GetUserCommentsByPosts(postId, userId int) ([]*entities.Comment, error) {
	stmt := `SELECT c.id, post_id, c.user_id, username, content, c.created, role
	FROM comments as c INNER JOIN users as u ON c.user_id = u.id
	WHERE u.id = ? AND c.post_id = ?`
	rows, err := r.DB.Query(stmt, userId, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		comment := &entities.Comment{}
		var created string

		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.UserName, &comment.Content, &created, &comment.UserRole); err != nil {
			return nil, err
		}

		commentTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}

		comment.Created = commentTime.Format(time.RFC3339)
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *CommentSqlite3) GetComments(postID int) ([]*entities.Comment, error) {
	stmt := `SELECT comments.id, post_id, username, comments.user_id, content, comments.created  
	FROM comments LEFT JOIN users ON users.id = comments.user_id
	WHERE post_id = ?
	ORDER BY comments.created DESC`

	rows, err := c.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		comment := &entities.Comment{}
		var created string
		var username sql.NullString

		if err := rows.Scan(&comment.ID, &comment.PostID, &username, &comment.UserID, &comment.Content, &created); err != nil {
			return nil, err
		}
		if username.Valid {
			comment.UserName = username.String
		} else {
			comment.UserName = "Deleted User"
		}

		commentTime, err := time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			return nil, err
		}

		comment.Created = commentTime.Format(time.RFC3339)
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *CommentSqlite3) DeleteComment(commentID int) error {
	stmt := `
	DELETE FROM comments
	WHERE id = ?
	`
	_, err := c.DB.Exec(stmt, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentSqlite3) UpdateComment(commentID int, content string) error {
	stmt := `
	UPDATE comments
	SET content = ?
	WHERE id = ?
	`
	_, err := c.DB.Exec(stmt, content, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (r *CommentSqlite3) GetComment(commentId int) (*entities.Comment, error) {
	stmt := `SELECT id, post_id, user_id, content FROM comments
	WHERE id = ?
	`

	row := r.DB.QueryRow(stmt, commentId)
	c := &entities.Comment{}

	err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return c, nil
}
