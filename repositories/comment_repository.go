package repositories

import "forum/entities"

type CommentRepository interface {
	InsertComment(postID, userID int, content string) error
	GetComments(postID int) ([]*entities.Comment, error)
}
