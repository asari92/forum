package service

import (
	"forum/internal/entities"
	"forum/internal/repository"
	"forum/pkg/validator"
)

type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}
type CategoryForm struct {
	Name string
	validator.Validator
}

func (uc *CategoryUseCase) NewCategoryCreateForm() CategoryForm {
	return CategoryForm{}
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{categoryRepo: categoryRepo}
}

func (u *CategoryUseCase) Insert(form *CategoryForm) (int, error) {
	form.CheckField(validator.NotBlank(form.Name), "category", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Name, validator.TextRX), "category", "This field must contain only english or russian letters")
	exists, err := u.categoryRepo.ExistName(form.Name)
	if err != nil {
		return 0, err
	}
	if exists {
		form.CheckField(false, "category", "error dublicate name")

	}
	if !form.Valid() {
		return 0, entities.ErrInvalidData
	}

	return u.categoryRepo.Insert(form.Name)
}

func (u *CategoryUseCase) Get(categoryId int) (*entities.Category, error) {
	return u.categoryRepo.Get(categoryId)
}

func (u *CategoryUseCase) GetAll() ([]*entities.Category, error) {
	return u.categoryRepo.GetAll()
}

func (u *CategoryUseCase) Delete(categoryId int) error {
	exists, err := u.categoryRepo.Exists(categoryId)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrNoRecord
	}
	if categoryId == 1 {
		return entities.ErrInvalidData
	}

	return u.categoryRepo.Delete(categoryId)
}
