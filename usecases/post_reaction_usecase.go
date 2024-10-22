package usecases

import "forum/entities"

type PostReactionRepository interface {
	AddReaction(userID, postID int, isLike bool) error
	RemoveReaction(userID, postID int) error
	GetUserReaction(userID, postID int) (*entities.PostReaction, error)
	GetReactionsCount(postID int) (likes int, dislikes int, err error)
}

type PostReactionUseCase struct {
	PostReactionRepo PostReactionRepository
}

func NewPostReactionUseCase(postReactionRepo PostReactionRepository) *PostReactionUseCase {
	return &PostReactionUseCase{PostReactionRepo: postReactionRepo}
}

func (u *PostReactionUseCase) AddReaction(userID, postID int, isLike bool) error {
	return u.PostReactionRepo.AddReaction(userID, postID, isLike)
}

func (u *PostReactionUseCase) RemoveReaction(userID, postID int) error {
	return u.PostReactionRepo.RemoveReaction(userID, postID)
}

func (u *PostReactionUseCase) GetUserReaction(userID, postID int) (*entities.PostReaction, error) {
	return u.PostReactionRepo.GetUserReaction(userID, postID)
}

func (u *PostReactionUseCase) GetReactionsCount(postID int) (int, int, error) {
	return u.PostReactionRepo.GetReactionsCount(postID)
}
