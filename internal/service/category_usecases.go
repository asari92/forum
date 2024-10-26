package service

import (
	"errors"

	"forum/internal/entities"
	"forum/internal/repository"
)

type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: categoryRepo}
}

func (u *CategoryUseCase) Insert(name string) (int, error) {
	if name == "" {
		return 0, errors.New("category name can't be empty")
	}
	return u.categoryRepo.Insert(name)
}

func (u *CategoryUseCase) Get(categoryId int) (*entities.Category, error) {
	return u.categoryRepo.Get(categoryId)
}

func (u *CategoryUseCase) GetAll() ([]*entities.Category, error) {
	return u.categoryRepo.GetAll()
}

func (u *CategoryUseCase) Delete(categoryId int) error {
	return u.categoryRepo.Delete(categoryId)
}