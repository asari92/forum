package repositories

import "forum/entities"

type CommentModelInterface interface {
	InsertComment(postID, userID int, content string) error
	GetComments(postID int) ([]*entities.Comment, error)
}
