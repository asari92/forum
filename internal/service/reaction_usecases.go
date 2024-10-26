package service

import (
	"forum/internal/entities"
	"forum/internal/repository"
)

type ReactionUseCase struct {
	categoryRepo        repository.CategoryRepository
	commentRepo         repository.CommentRepository
	commentReactionRepo repository.CommentReactionRepository
	postRepo            repository.PostRepository
	postReactionRepo    repository.PostReactionRepository
	userRepo            repository.UserRepository
}

type reactionForm struct {
	Comment       string
	PostIsLike    string
	CommentIsLike string
	CommentID     int
}

func NewReactionUseCase(repo *repository.Repository) *ReactionUseCase {
	return &ReactionUseCase{
		categoryRepo:        repo.CategoryRepository,
		commentRepo:         repo.CommentRepository,
		commentReactionRepo: repo.CommentReactionRepository,
		postRepo:            repo.PostRepository,
		postReactionRepo:    repo.PostReactionRepository,
		userRepo:            repo.UserRepository,
	}
}

func (uc *ReactionUseCase) NewReactionForm() reactionForm {
	return reactionForm{}
}

func (uc *ReactionUseCase) UpdatePostReaction(userID, postID int, form *reactionForm) error {
	if form.Comment != "" {
		err := uc.commentRepo.InsertComment(postID, userID, form.Comment)
		if err != nil {
			return err
		}
	} else if form.PostIsLike != "" {
		// Преобразуем isLike в bool
		like := form.PostIsLike == "true"

		var userReaction *entities.PostReaction
		userReaction, err := uc.postReactionRepo.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			return err
		}

		if userReaction != nil && userReaction.IsLike == like {
			err = uc.postReactionRepo.RemoveReaction(userID, postID)
			if err != nil {
				return err
			}
		} else {
			err = uc.postReactionRepo.AddReaction(userID, postID, like)
			if err != nil {
				return err
			}
		}
	} else if form.CommentIsLike != "" {
		reaction := form.CommentIsLike == "true"

		commentReaction, err := uc.commentReactionRepo.GetUserReaction(userID, form.CommentID)
		if err != nil {
			return err
		}

		if commentReaction != nil && reaction == commentReaction.IsLike {
			err = uc.commentReactionRepo.RemoveReaction(userID, form.CommentID)
			if err != nil {
				return err
			}
		} else {
			err = uc.commentReactionRepo.AddReaction(userID, form.CommentID, reaction)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (uc *ReactionUseCase) UpdateCommentReaction(userID int, form *reactionForm) error {
	reaction := form.CommentIsLike == "true"
	var commentReaction *entities.CommentReaction

	commentReaction, err := uc.commentReactionRepo.GetUserReaction(userID, form.CommentID)
	if err != nil {
		return err
	}

	if reaction == commentReaction.IsLike {
		err = uc.commentReactionRepo.RemoveReaction(userID, form.CommentID)
		if err != nil {
			return err
		} else {
			err = uc.commentReactionRepo.AddReaction(userID, form.CommentID, reaction)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
