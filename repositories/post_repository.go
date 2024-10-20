package repositories

import "forum/entities"

type PostModelInterface interface {
	InsertPostWithCategories(title, content string, userID int, categoryIDs []int) (int, error)
	Get(id int) (*entities.Post, error)
	GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error)
	GetUserPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error)
	GetUserLikedPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error)
	GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error)
}
