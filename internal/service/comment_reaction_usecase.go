package service

import (
	"forum/internal/entities"
	"forum/internal/repository"
)

type CommentReactionUseCase struct {
	commentReactionRepo repository.CommentReactionRepository
}

func NewCommentReactionUseCase(commentReactionRepo repository.CommentReactionRepository) *CommentReactionUseCase {
	return &CommentReactionUseCase{commentReactionRepo: commentReactionRepo}
}

func (u *CommentReactionUseCase) AddReaction(userID, commentID int, isLike bool) error {
	return u.commentReactionRepo.AddReaction(userID, commentID, isLike)
}

func (u *CommentReactionUseCase) RemoveReaction(userID, commentID int) error {
	return u.commentReactionRepo.RemoveReaction(userID, commentID)
}

func (u *CommentReactionUseCase) GetUserReaction(userID, commentID int) (*entities.CommentReaction, error) {
	return u.commentReactionRepo.GetUserReaction(userID, commentID)
}

func (u *CommentReactionUseCase) GetLikesCount(commentID int) (int, error) {
	return u.commentReactionRepo.GetLikesCount(commentID)
}

func (u *CommentReactionUseCase) GetDislikesCount(commentID int) (int, error) {
	return u.commentReactionRepo.GetDislikesCount(commentID)
}
