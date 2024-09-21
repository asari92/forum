package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	CommentId int
	PostId    int
	UserId    int
	Content   string
	Created   time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func (c CommentModel) Insert(postid, userid int, content string) error {
	stmt := `INSERT INTO comments (post_id, user_id, content, created) VALUES(?,?,?,datetime('now'))`
	_, err := c.DB.Exec(stmt, postid, userid, content)
	if err != nil {
		return err
	}

	return nil

}

func (c CommentModel) GetComments(postid int) ([]*Comment, error) {
	stmt := `SELECT id, content, user_id, created 
	         FROM comments
			 WHERE post_id = ?`
	rows, err := c.DB.Query(stmt, postid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	for rows.Next() {
		comment := &Comment{}
		err = rows.Scan(&comment.CommentId, &comment.Content, &comment.UserId, &comment.Created)
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
