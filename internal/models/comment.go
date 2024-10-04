package models

import (
	"database/sql"
)

type CommentModelInterface interface {
	InsertComment(postID, userID int, content string) error
	GetComments(postID int) ([]*Comment, error)
}

type Comment struct {
	ID      int
	PostID  int
	UserID  int
	UserName string
	Content string
	Created string
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
	stmt := `SELECT post_id, username, content, comments.created
	FROM comments INNER JOIN users ON users.id = comments.user_id
	WHERE post_id = ?
	ORDER BY created DESC`
	comments := []*Comment{}
	rows, err := c.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		comment := &Comment{}
		err = rows.Scan(&comment.PostID, &comment.UserName, &comment.Content, &comment.Created)
		if err != nil {
			return nil, err
		} else {

			comments = append(comments, comment)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
