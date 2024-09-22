package mocks

import (
	"time"

	"forum/internal/models"
)

var mockPost = &models.Post{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	UserID:  1,
	Created: time.Now(),
}

type PostModel struct{}

func (m *PostModel) InsertPostWithCategories(title, content string, userID int, categoryIDs []int) (int, error) {
	return 2, nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	switch id {
	case 1:
		return mockPost, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	return []*models.Post{mockPost}, nil
}

func (m *PostModel) GetPostsForCategory(categoryIDs []int) ([]*models.Post, error) {
	return nil, nil
}

func (m *PostModel) GetUserPosts(userId int) ([]*models.Post, error) {
	switch userId {
	case 1:
		return []*models.Post{mockPost}, nil
	default:
		return nil, models.ErrNoRecord
	}
}
