package service

import (
	"errors"
	"forum/internal/entities"
	"forum/internal/repository"
)

type CommentUseCase struct {
	commentRepo repository.CommentRepository
}

func NewCommentUseCase(commentRepo repository.CommentRepository) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
	}
}

func (c *CommentUseCase) AddComment(postID, userID int, content string) error {
	if content == "" {
		return errors.New("content cannot be empty")
	}
	return c.commentRepo.InsertComment(postID, userID, content)
}

func (c *CommentUseCase) GetPostComments(postID int) ([]*entities.Comment, error) {
	return c.commentRepo.GetComments(postID)
}
