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
	ID      int
	PostID  int
	UserID  int
	Content string
	Created time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func (c *CommentModel) InsertComment(postID, userID int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id,content, created)
	VALUES (?,?,?, datetime('now','localtime'))`
	_, err := c.DB.Exec(stmt, postID, userID, content)
	if err != nil {
		return err
	}

	return nil

}

func (c *CommentModel) GetComments(postID int) ([]*Comment, error) {
	stmt := `SELECT post_id, user_id, content, created FROM comments
	WHERE post_id = ?`
	comments := []*Comment{}
	rows, err := c.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		comment := &Comment{}
		var Created string
		err = rows.Scan(&comment.PostID, &comment.UserID, &comment.Content, &Created)
		if err != nil {
			return nil, err
		} else {
			comment.Created, err = time.Parse("2006-01-02 15:04:05", Created)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
