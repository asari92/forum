package mocks

import (
	"time"

	"forum/entities"
	"forum/repositories"
)

var mockPost = &entities.Post{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	UserID:  1,
	Created: time.Now().Format("2006-01-02 15:04:05"),
}

type PostModel struct{}

func (m *PostModel) InsertPostWithCategories(title, content string, userID int, categoryIDs []int) (int, error) {
	return 2, nil
}

func (m *PostModel) Get(id int) (*entities.Post, error) {
	switch id {
	case 1:
		return mockPost, nil
	default:
		return nil, repositories.ErrNoRecord
	}
}

func (m *PostModel) GetAllPaginatedPosts(page, pageSize int) ([]*entities.Post, error) {
	return []*entities.Post{mockPost}, nil
}

func (m *PostModel) GetPaginatedPostsByCategory(categoryIDs []int, page, pageSize int) ([]*entities.Post, error) {
	return []*entities.Post{mockPost}, nil
}

func (m *PostModel) GetUserPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error) {
	switch userId {
	case 1:
		return []*entities.Post{mockPost}, nil
	default:
		return nil, repositories.ErrNoRecord
	}
}

func (m *PostModel) GetUserLikedPaginatedPosts(userId, page, pageSize int) ([]*entities.Post, error) {
	switch userId {
	case 1:
		return []*entities.Post{mockPost}, nil
	default:
		return nil, repositories.ErrNoRecord
	}
}
