package models

import (
	"database/sql"
	"time"
)

type CommentModelInterface interface {
	InsertComment(postID, userID int, content string) error
	GetComments(postID int) ([]*Comment, error)
}

type Comment struct {
	ID       int
	PostID   int
	UserID   int
	UserName string
	Content  string
	Created  string
}

type CommentModel struct {
	DB *sql.DB
}

func (c *CommentModel) InsertComment(postID, userID int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id,content, created)
	VALUES (?,?,?, datetime('now'))`
	_, err := c.DB.Exec(stmt, postID, userID, content)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentModel) GetComments(postID int) ([]*Comment, error) {
	stmt := `SELECT post_id, username, content, comments.created  
	FROM comments INNER JOIN users ON users.id = comments.user_id
	WHERE post_id = ?
	ORDER BY comments.created DESC`

	rows, err := c.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		var created string

		if err := rows.Scan(&comment.PostID, &comment.UserName, &comment.Content, &created); err != nil {
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