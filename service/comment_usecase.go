package service

import (
	"errors"

	"forum/entities"
	"forum/repository"
)

type CommentUseCase struct {
	CommentRepo repository.CommentRepository
}

func NewCommentUseCase(commentRepo repository.CommentRepository) *CommentUseCase {
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
