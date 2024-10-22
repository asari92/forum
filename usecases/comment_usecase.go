package usecases

import (
	"errors"

	"forum/entities"
)

type CommentRepository interface {
	InsertComment(postID, userID int, content string) error
	GetComments(postID int) ([]*entities.Comment, error)
}

type CommentUseCase struct {
	CommentRepo CommentRepository
}

func NewCommentUseCase(commentRepo CommentRepository) *CommentUseCase {
	return &CommentUseCase{
		CommentRepo: commentRepo,
	}
}

func (c *CommentUseCase) AddComment(postID, userID int, content string) error {
	if content == "" {
		return errors.New("content cannot be empty")
	}
	return c.CommentRepo.InsertComment(postID, userID, content)
}

func (c *CommentUseCase) GetPostComments(postID int) ([]*entities.Comment, error) {
	return c.CommentRepo.GetComments(postID)
}
