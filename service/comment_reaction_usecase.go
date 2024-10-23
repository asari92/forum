package service

import (
	"forum/entities"
	"forum/repository"
)

type CommentReactionUseCase struct {
	CommentReactionRepo repository.CommentReactionRepository
}

func NewCommentReactionUseCase(commentReactionRepo repository.CommentReactionRepository) *CommentReactionUseCase {
	return &CommentReactionUseCase{CommentReactionRepo: commentReactionRepo}
}

func (u *CommentReactionUseCase) AddReaction(userID, commentID int, isLike bool) error {
	return u.CommentReactionRepo.AddReaction(userID, commentID, isLike)
}

func (u *CommentReactionUseCase) RemoveReaction(userID, commentID int) error {
	return u.CommentReactionRepo.RemoveReaction(userID, commentID)
}

func (u *CommentReactionUseCase) GetUserReaction(userID, commentID int) (*entities.CommentReaction, error) {
	return u.CommentReactionRepo.GetUserReaction(userID, commentID)
}

func (u *CommentReactionUseCase) GetLikesCount(commentID int) (int, error) {
	return u.CommentReactionRepo.GetLikesCount(commentID)
}

func (u *CommentReactionUseCase) GetDislikesCount(commentID int) (int, error) {
	return u.CommentReactionRepo.GetDislikesCount(commentID)
}
