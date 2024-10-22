package repositories

import "forum/entities"

type CategoryModelRepository interface {
	Insert(name string) (int, error)
	Get(categoryId int) (*entities.Category, error)
	GetAll() ([]*entities.Category, error)
	Delete(categoryId int) error
	GetCategoriesForPost(postId int) ([]*entities.Category, error)
	DeleteCategoriesForPost(postId int) error
}
