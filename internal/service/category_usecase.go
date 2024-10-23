package service

import (
	"errors"

	"forum/internal/entities"
	"forum/internal/repository"
)

type CategoryUseCase struct {
	CategoryRepo repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{CategoryRepo: categoryRepo}
}

func (u *CategoryUseCase) Insert(name string) (int, error) {
	if name == "" {
		return 0, errors.New("category name can't be empty")
	}
	return u.CategoryRepo.Insert(name)
}

func (u *CategoryUseCase) Get(categoryId int) (*entities.Category, error) {
	return u.CategoryRepo.Get(categoryId)
}

func (u *CategoryUseCase) GetAll() ([]*entities.Category, error) {
	return u.CategoryRepo.GetAll()
}

func (u *CategoryUseCase) Delete(categoryId int) error {
	return u.CategoryRepo.Delete(categoryId)
}

func (u *CategoryUseCase) GetCategoriesForPost(postId int) ([]*entities.Category, error) {
	return u.CategoryRepo.GetCategoriesForPost(postId)
}

func (u *CategoryUseCase) DeleteCategoriesForPost(postId int) error {
	return u.CategoryRepo.DeleteCategoriesForPost(postId)
}
