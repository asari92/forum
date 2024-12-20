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
	reportRepo          repository.ReportRepository
}

type reactionForm struct {
	Comment       string
	PostIsLike    string
	CommentIsLike string
	CommentID     int
	validator.Validator
}

type ReportsDTO struct {
	Reports       []*entities.Report
	HasNextPage   bool
	CurrentPage   int
	PaginationURL string
	Header        string
}

func NewReactionUseCase(repo *repository.Repository) *ReactionUseCase {
	return &ReactionUseCase{
		categoryRepo:        repo.CategoryRepository,
		commentRepo:         repo.CommentRepository,
		commentReactionRepo: repo.CommentReactionRepository,
		postRepo:            repo.PostRepository,
		postReactionRepo:    repo.PostReactionRepository,
		userRepo:            repo.UserRepository,
		reportRepo:          repo.ReportRepository,
	}
}

func (ruc *ReactionUseCase) NewReactionForm() reactionForm {
	return reactionForm{}
}

func (ruc *ReactionUseCase) UpdatePostReaction(userID, postID int, form *reactionForm) error {
	exists, err := ruc.userRepo.Exists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	exists, err = ruc.postRepo.Exists(postID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	ownerID, err := ruc.postRepo.GetPostOwner(postID)
	if err != nil {
		return err
	}

	if form.Comment != "" {
		form.CheckField(validator.NotBlank(form.Comment), "comment", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Comment, validator.TextRX), "comment", "This field must contain only english or russian letters")
		if !form.Valid() {
			return entities.ErrInvalidData
		}
		err := ruc.commentRepo.InsertComment(postID, userID, form.Comment)
		if err != nil {
			return err
		}
		err = ruc.postReactionRepo.AddNotification(ownerID, postID, userID, "comment")
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
		userReaction, err := ruc.postReactionRepo.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			return err
		}

		if userReaction != nil && userReaction.IsLike == like {
			err = ruc.postReactionRepo.RemoveReaction(userID, postID)
			if err != nil {
				return err
			}
			fmt.Println(action)
			err = ruc.postReactionRepo.RemoveNotification(ownerID, postID, userID, action)
			if err != nil {
				return err
			}

		} else {

			if userReaction == nil {
				err = ruc.postReactionRepo.AddNotification(ownerID, postID, userID, action)
				if err != nil {
					return err
				}
			} else if userReaction.IsLike != like {
				var oldAction string
				if userReaction.IsLike {
					oldAction = "like"
				} else {
					oldAction = "dislike"
				}

				err = ruc.postReactionRepo.UpdateNotification(ownerID, postID, userID, oldAction, action)
				if err != nil {
					return err
				}

			}

			err = ruc.postReactionRepo.AddReaction(userID, postID, like)
			if err != nil {
				return err
			}

		}
	} else if form.CommentIsLike != "" {
		reaction := form.CommentIsLike == "true"

		exists, err := ruc.commentRepo.Exists(form.CommentID)
		if err != nil {
			return err
		}
		if !exists {
			return entities.ErrNoRecord
		}

		commentReaction, err := ruc.commentReactionRepo.GetUserReaction(userID, form.CommentID)
		if err != nil {
			return err
		}

		if commentReaction != nil && reaction == commentReaction.IsLike {

			err = ruc.commentReactionRepo.RemoveReaction(userID, form.CommentID)
			if err != nil {
				return err
			}
		} else {
			err = ruc.commentReactionRepo.AddReaction(userID, form.CommentID, reaction)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ruc *ReactionUseCase) CreateReport(userID, postID int, reason string) error {
	exists, err := ruc.userRepo.Exists(userID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	exists, err = ruc.postRepo.Exists(postID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}

	err = ruc.reportRepo.CreateReport(userID, postID, reason)
	if err != nil {
		return err
	}

	return nil
}

func (ruc *ReactionUseCase) GetPostReport(postID int) (*entities.Report, error) {
	report, err := ruc.reportRepo.GetPostReport(postID)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (ruc *ReactionUseCase) GetAllPaginatedPostReportsDTO(page, pageSize int, paginationURL string) (*ReportsDTO, error) {
	// Получаем репорты для нужной страницы.
	reports, err := ruc.reportRepo.GetAllPaginatedPostReports(page, pageSize)
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли следующая страница
	hasNextPage := len(reports) > pageSize
	// Если постов больше, чем pageSize, обрезаем список до pageSize
	if hasNextPage {
		reports = reports[:pageSize]
	}

	return &ReportsDTO{
		Reports:       reports,
		HasNextPage:   hasNextPage,
		CurrentPage:   page,
		PaginationURL: paginationURL,
	}, nil
}
