package service

import (
	"forum/entities"
	"forum/repository"
)

type PostReactionUseCase struct {
	PostReactionRepo repository.PostReactionRepository
}

func NewPostReactionUseCase(postReactionRepo repository.PostReactionRepository) *PostReactionUseCase {
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
