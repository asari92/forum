package repository

import (
	"database/sql"
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

func (c *CommentSqlite3) InsertComment(postID, userID int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id,content, created)
	VALUES (?,?,?, datetime('now'))`
	_, err := c.DB.Exec(stmt, postID, userID, content)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentSqlite3) GetComments(postID int) ([]*entities.Comment, error) {
	stmt := `SELECT comments.id, post_id, username, content, comments.created  
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

		if err := rows.Scan(&comment.ID, &comment.PostID, &username, &comment.Content, &created); err != nil {
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
