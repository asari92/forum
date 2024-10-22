package usecases

import "forum/entities"

type CommentReactionRepository interface {
	AddReaction(userID, commentID int, isLike bool) error
	RemoveReaction(userID, commentID int) error
	GetUserReaction(userID, commentID int) (*entities.CommentReaction, error)
	GetLikesCount(commentID int) (int, error)
	GetDislikesCount(commentID int) (int, error)
}

type CommentReactionUseCase struct {
	CommentReactionRepo CommentReactionRepository
}

func NewCommentReactionUseCase(commentReactionRepo CommentReactionRepository) *CommentReactionUseCase {
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
