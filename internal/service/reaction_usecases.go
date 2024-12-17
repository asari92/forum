package service

import (
	"fmt"

	"forum/internal/entities"
	"forum/internal/repository"
	"forum/pkg/validator"
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
	validator.Validator
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
	exists, err := uc.userRepo.Exists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	exists, err = uc.postRepo.Exists(postID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	ownerID, err := uc.postRepo.GetPostOwner(postID)
	if err != nil {
		return err
	}

	if form.Comment != "" {
		form.CheckField(validator.NotBlank(form.Comment), "comment", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Comment, validator.TextRX), "comment", "This field must contain only english or russian letters")
		if !form.Valid() {
			return entities.ErrInvalidData
		}
		err := uc.commentRepo.InsertComment(postID, userID, form.Comment)
		if err != nil {
			return err
		}
		err = uc.postReactionRepo.AddNotification(ownerID, postID, userID, "comment")
		if err != nil {
			return err
		}

	} else if form.PostIsLike != "" {
		// Преобразуем isLike в bool
		like := form.PostIsLike == "true"
		var action string
		if like {
			action = "like"
		} else {
			action = "dislike"
		}

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
			fmt.Println(action)
			err = uc.postReactionRepo.RemoveNotification(ownerID, postID, userID, action)
			if err != nil {
				return err
			}

		} else {
			if userReaction == nil {
				err = uc.postReactionRepo.AddNotification(ownerID, postID, userID, action)
				if err != nil {
					return err
				}
			}else if userReaction.IsLike != like {
				newA


			}

			

			err = uc.postReactionRepo.AddReaction(userID, postID, like)
			if err != nil {
				return err
			}
			

		}
	} else if form.CommentIsLike != "" {
		reaction := form.CommentIsLike == "true"

		exists, err := uc.commentRepo.Exists(form.CommentID)
		if err != nil {
			return err
		}
		if !exists {
			return entities.ErrNoRecord
		}

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
