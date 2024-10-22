package mocks

import (
	"forum/entities"
	"forum/repositories"
)

type CategoryModel struct{}

var mockCategory = &entities.Category{
	ID:   1,
	Name: "Test",
}

func (m *CategoryModel) GetCategoriesForPost(postId int) ([]*entities.Category, error) {
	categories := []*entities.Category{}
	switch postId {
	case 1:
		categories = append(categories, mockCategory)
		return categories, nil
	default:
		return nil, repositories.ErrNoRecord
	}
}

func (m *CategoryModel) Insert(name string) (int, error) {
	return 1, nil
}

func (m *CategoryModel) Get(categoryId int) (*entities.Category, error) {
	c := &entities.Category{}

	return c, nil
}

func (m *CategoryModel) GetAll() ([]*entities.Category, error) {
	return nil, nil
}

func (m *CategoryModel) Delete(categoryId int) error {
	return nil
}

func (m *CategoryModel) DeleteCategoriesForPost(postId int) error {
	return nil
}
