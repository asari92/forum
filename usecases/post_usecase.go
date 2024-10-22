package usecases

import (
	"forum/entities"
)

// Интерфейс репозитория
type PostRepository interface {
	InsertPostWithCategories(title, content string, userID int, categoryIDs []int) (int, error)
	GetPost(postID int) (*entities.Post, error)
	GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error)
	GetUserPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error)
	GetUserLikedPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error)
	GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error)
	DeletePost(postID int) error
}

// Use Case структура
type PostUseCase struct {
	PostRepo PostRepository
}

func NewPostUseCase(postRepo PostRepository) *PostUseCase {
	return &PostUseCase{
		PostRepo: postRepo,
	}
}

// Создание поста с категориями
func (uc *PostUseCase) CreatePostWithCategories(title, content string, userID int, categoryIDs []int) (int, error) {
	postID, err := uc.PostRepo.InsertPostWithCategories(title, content, userID, categoryIDs)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

// Получение поста по ID
func (uc *PostUseCase) GetPost(postID int) (*entities.Post, error) {
	return uc.PostRepo.GetPost(postID)
}

// Получение постов по категориям с пагинацией
func (uc *PostUseCase) GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error) {
	return uc.PostRepo.GetPaginatedPostsByCategory(categoryIDs, page, pageSize)
}

// Получение постов пользователя с пагинацией
func (uc *PostUseCase) GetUserPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error) {
	return uc.PostRepo.GetUserPaginatedPosts(userID, page, pageSize)
}

// Получение постов, которые пользователь лайкнул
func (uc *PostUseCase) GetUserLikedPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error) {
	return uc.PostRepo.GetUserLikedPaginatedPosts(userID, page, pageSize)
}

// Получение всех постов с пагинацией
func (uc *PostUseCase) GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error) {
	return uc.PostRepo.GetAllPaginatedPosts(page, pageSize)
}

// Удаление поста
func (uc *PostUseCase) DeletePost(postID int) error {
	return uc.PostRepo.DeletePost(postID)
}
