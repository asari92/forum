package service

import (
	"forum/internal/entities"
	"forum/internal/repository"
)

type PostReactionUseCase struct {
	postReactionRepo repository.PostReactionRepository
}

func NewPostReactionUseCase(postReactionRepo repository.PostReactionRepository) *PostReactionUseCase {
	return &PostReactionUseCase{postReactionRepo: postReactionRepo}
}

func (u *PostReactionUseCase) AddReaction(userID, postID int, isLike bool) error {
	return u.postReactionRepo.AddReaction(userID, postID, isLike)
}

func (u *PostReactionUseCase) RemoveReaction(userID, postID int) error {
	return u.postReactionRepo.RemoveReaction(userID, postID)
}

func (u *PostReactionUseCase) GetUserReaction(userID, postID int) (*entities.PostReaction, error) {
	return u.postReactionRepo.GetUserReaction(userID, postID)
}

func (u *PostReactionUseCase) GetReactionsCount(postID int) (int, int, error) {
	return u.postReactionRepo.GetReactionsCount(postID)
}
