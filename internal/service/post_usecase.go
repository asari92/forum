package service

import (
	"errors"

	"forum/internal/entities"
	"forum/internal/repository"
)

// Use Case структура
type PostUseCase struct {
	categoryRepo        repository.CategoryRepository
	commentRepo         repository.CommentRepository
	commentReactionRepo repository.CommentReactionRepository
	postRepo            repository.PostRepository
	postReactionRepo    repository.PostReactionRepository
	userRepo            repository.UserRepository
}

type PostDTO struct {
	Post         *entities.Post
	Categories   []*entities.Category
	Likes        int
	Dislikes     int
	Comments     []*entities.Comment
	UserReaction *entities.PostReaction
}

type UserPostsDTO struct {
	User          *entities.User
	Posts         []*entities.Post
	HasNextPage   bool
	CurrentPage   int
	PaginationURL string
}

func NewPostUseCase(repo *repository.Repository) *PostUseCase {
	return &PostUseCase{
		categoryRepo:        repo.CategoryRepository,
		commentRepo:         repo.CommentRepository,
		commentReactionRepo: repo.CommentReactionRepository,
		postRepo:            repo.PostRepository,
		postReactionRepo:    repo.PostReactionRepository,
		userRepo:            repo.UserRepository,
	}
}

func (uc *PostUseCase) GetPostDTO(postID int, userID int) (*PostDTO, error) {
	post, err := uc.postRepo.GetPost(postID)
	if err != nil {
		return nil, err
	}

	categories, err := uc.categoryRepo.GetCategoriesForPost(postID)
	if err != nil {
		return nil, err
	}

	likes, dislikes, err := uc.postReactionRepo.GetReactionsCount(postID)
	if err != nil {
		return nil, err
	}

	var userReaction *entities.PostReaction
	if userID > 0 {
		userReaction, err = uc.postReactionRepo.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			return nil, err
		}
	}

	comments, err := uc.commentRepo.GetComments(postID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if userID > 0 {
			userReaction, err := uc.commentReactionRepo.GetUserReaction(userID, comment.ID)
			if err != nil {
				return nil, err
			}
			if userReaction != nil {
				if userReaction.IsLike {
					comment.UserReaction = 1
				} else {
					comment.UserReaction = -1
				}
			}
		}

		like, err := uc.commentReactionRepo.GetLikesCount(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Like = like

		dislike, err := uc.commentReactionRepo.GetDislikesCount(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Dislike = dislike
	}

	return &PostDTO{
		Post:         post,
		Categories:   categories,
		Likes:        likes,
		Dislikes:     dislikes,
		Comments:     comments,
		UserReaction: userReaction,
	}, nil
}

// Получение постов пользователя с пагинацией
func (uc *PostUseCase) GetUserPostsDTO(userID, page, pageSize int, paginationURL string) (*UserPostsDTO, error) {
	posts, err := uc.postRepo.GetUserPaginatedPosts(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.Get(userID)
	if err != nil {
		if errors.Is(err, entities.ErrNoRecord) {
			if len(posts) == 0 {
				return nil, err
			} else {
				user = &entities.User{Username: "Deleted User", ID: userID}
			}
		} else {
			return nil, err
		}
	}

	hasNextPage := len(posts) > pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	return &UserPostsDTO{
		User:          user,
		Posts:         posts,
		HasNextPage:   hasNextPage,
		CurrentPage:   page,
		PaginationURL: paginationURL,
	}, nil
}

// Получение постов, которые пользователь лайкнул
func (uc *PostUseCase) GetUserLikedPostsDTO(userID, page, pageSize int, paginationURL string) (*UserPostsDTO, error) {
	// Получаем посты, которые пользователь лайкнул
	posts, err := uc.postRepo.GetUserLikedPaginatedPosts(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Получаем пользователя
	user, err := uc.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	// Проверяем наличие следующей страницы
	hasNextPage := len(posts) > pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	return &UserPostsDTO{
		User:          user,
		Posts:         posts,
		HasNextPage:   hasNextPage,
		CurrentPage:   page,
		PaginationURL: paginationURL,
	}, nil
}

// Создание поста с категориями
func (uc *PostUseCase) CreatePostWithCategories(title, content string, userID int, categoryIDs []int) (int, error) {
	postID, err := uc.postRepo.InsertPostWithCategories(title, content, userID, categoryIDs)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

// // Получение поста по ID
// func (uc *PostUseCase) GetPost(postID int) (*entities.Post, error) {
// 	return uc.postRepo.GetPost(postID)
// }

// Получение постов по категориям с пагинацией
func (uc *PostUseCase) GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error) {
	return uc.postRepo.GetPaginatedPostsByCategory(categoryIDs, page, pageSize)
}

// func (uc *PostUseCase) GetUserPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error) {
// 	return uc.postRepo.GetUserPaginatedPosts(userID, page, pageSize)
// }

// func (uc *PostUseCase) GetUserLikedPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error) {
// 	return uc.postRepo.GetUserLikedPaginatedPosts(userID, page, pageSize)
// }

// Получение всех постов с пагинацией
func (uc *PostUseCase) GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error) {
	return uc.postRepo.GetAllPaginatedPosts(page, pageSize)
}

// Удаление поста
func (uc *PostUseCase) DeletePost(postID int) error {
	return uc.postRepo.DeletePost(postID)
}
