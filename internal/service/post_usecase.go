package service

import (
	"forum/internal/entities"
	"forum/internal/repository"
)

// Use Case структура
type PostUseCase struct {
	postRepo repository.PostRepository
}

func NewPostUseCase(postRepo repository.PostRepository) *PostUseCase {
	return &PostUseCase{
		postRepo: postRepo,
	}
}

// Создание поста с категориями
func (uc *PostUseCase) CreatePostWithCategories(title, content string, userID int, categoryIDs []int) (int, error) {
	postID, err := uc.postRepo.InsertPostWithCategories(title, content, userID, categoryIDs)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

// Получение поста по ID
func (uc *PostUseCase) GetPost(postID int) (*entities.Post, error) {
	return uc.postRepo.GetPost(postID)
}

// Получение постов по категориям с пагинацией
func (uc *PostUseCase) GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error) {
	return uc.postRepo.GetPaginatedPostsByCategory(categoryIDs, page, pageSize)
}

// Получение постов пользователя с пагинацией
func (uc *PostUseCase) GetUserPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error) {
	return uc.postRepo.GetUserPaginatedPosts(userID, page, pageSize)
}

// Получение постов, которые пользователь лайкнул
func (uc *PostUseCase) GetUserLikedPaginatedPosts(userID, page, pageSize int) ([]*entities.Post, error) {
	return uc.postRepo.GetUserLikedPaginatedPosts(userID, page, pageSize)
}

// Получение всех постов с пагинацией
func (uc *PostUseCase) GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error) {
	return uc.postRepo.GetAllPaginatedPosts(page, pageSize)
}

// Удаление поста
func (uc *PostUseCase) DeletePost(postID int) error {
	return uc.postRepo.DeletePost(postID)
}
